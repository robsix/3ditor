package golog

import(
	`encoding/json`
	`sync`
	`time`
	`os`
	`io/ioutil`
	`sort`
	`path/filepath`
)

const(
	storeGrowth = 1000
)

// Maintains log entities in memory and saves them to the file system to restore them in between app restarts.
// Since this is mainly intended for local development use it is set to purge its memory and file stores at the next
// read or write operation after UTC midnight each day
func NewLocalLog(storeDir string, printToStdOut bool, lineSpacing int) (Log, error) {
	if err := os.MkdirAll(storeDir, os.ModeDir); err != nil {
		return nil, err
	}

	mtx := sync.Mutex{}
	memStore := make([]LogEntry, 0, storeGrowth)
	var lastPurge time.Time

	getFileName := func(logId string) string {
		return filepath.Join(storeDir, logId + `.json`)
	}
	getLastPurgeFileName := func() string {
		return filepath.Join(storeDir, "lastPurge.json")
	}
	setLastPurge := func(lp time.Time) error {
		lastPurge = lp
		data, _ := json.Marshal(lastPurge)
		return ioutil.WriteFile(getLastPurgeFileName(), data, os.ModePerm)
	}
	getLastPurge := func() time.Time {
		if lastPurge.IsZero() {
			if data, err := ioutil.ReadFile(getLastPurgeFileName()); err != nil {
				setLastPurge(time.Now().UTC())
			} else {
				if err = json.Unmarshal(data, &lastPurge); err != nil {
					setLastPurge(time.Now().UTC())
				}
			}
		}
		return lastPurge
	}
	purgeIfNewDay := func(){
		defer mtx.Unlock()
		mtx.Lock()
		if getLastPurge().UTC().Day() != time.Now().UTC().Day() {
			setLastPurge(time.Now().UTC())

			files, _ := ioutil.ReadDir(storeDir)
			for _, file := range files {
				os.Remove(filepath.Join(storeDir, file.Name()))
			}

			memStore = make([]LogEntry, 0, storeGrowth)
		}
	}

	appendToMemStore := func(le LogEntry){
		if len(memStore) == cap(memStore) {
			memStore = append(make([]LogEntry, 0, len(memStore) + storeGrowth), memStore...)
		}
		memStore = append(memStore, le)
	}

	files, _ := ioutil.ReadDir(storeDir)
	for _, file := range files {
		data, _ := ioutil.ReadFile(storeDir + `/` + file.Name())
		dst := LogEntry{}
		json.Unmarshal(data, &dst)
		insertIndex := sort.Search(len(memStore), func(i int)bool{
			return memStore[i].Time.After(dst.Time)
		})
		appendToMemStore(dst)
		if insertIndex != len(memStore) - 1{
			copy(memStore[insertIndex+1:], memStore[insertIndex:])
			memStore[insertIndex] = dst
		}
	}

	put := func(le LogEntry){
		purgeIfNewDay()
		data, _ := json.Marshal(le)
		fileName := getFileName(le.LogId)

		defer mtx.Unlock()
		mtx.Lock()

		ioutil.WriteFile(fileName, data, os.ModePerm)

		appendToMemStore(le)

		insertIndex := len(memStore)

		for i := insertIndex - 1; i >= 0; i--{
			if memStore[i].Time.Before(le.Time){
				break
			}
			insertIndex--
		}

		if insertIndex != len(memStore) {
			copy(memStore[insertIndex+1:], memStore[insertIndex:])
			memStore[insertIndex] = le
		}
	}

	getById := func(logId string) (LogEntry, error) {
		purgeIfNewDay()
		defer mtx.Unlock()
		mtx.Lock()

		for _, le := range memStore {
			if le.LogId == logId {
				return le, nil
			}
		}

		return LogEntry{}, &noSuchLogEntryError{}
	}

	get := func(before time.Time, level level, limit int) ([]LogEntry, error) {
		purgeIfNewDay()
		if limit <= 0 {
			return nil, &limitNotSetError{}
		}

		defer mtx.Unlock()
		mtx.Lock()

		if len(memStore) == 0 {
			return []LogEntry{}, nil
		}

		ret := make([]LogEntry, 0, limit)

		indexToStartFrom := sort.Search(len(memStore), func(i int)bool{
			return memStore[i].Time.After(before)
		}) - 1

		for i := indexToStartFrom; i >= 0 && len(ret) < limit; i--{
			if level == ANY || memStore[i].Level == level {
				ret = append(ret, memStore[i])
			}
		}

		return ret, nil
	}

	return NewLog(put, getById, get, printToStdOut, lineSpacing), nil
}

type noSuchLogEntryError struct{
	logId string
}
func (e *noSuchLogEntryError) Error() string {return `No such LogEntry exists with id: ` + e.logId}

type limitNotSetError struct{}
func (e *limitNotSetError) Error() string {return `A limit greater than 0 must be passed to get()`}
