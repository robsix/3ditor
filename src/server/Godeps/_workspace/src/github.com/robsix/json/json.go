package json

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
	"os"
)

type Json struct {
	data interface{}
}

// New returns a pointer to a new, empty `Json` object
func New() (*Json, error) {
	return FromString("{}")
}

// FromInterface returns a pointer to a new `Json` object
// after assigning `i` to its internal data
func FromInterface(i interface{}) *Json {
	return &Json{i}
}

// FromString returns a pointer to a new `Json` object
// after unmarshaling `str`
func FromString(str string) (*Json, error) {
	return FromBytes([]byte(str))
}

// FromBytes returns a pointer to a new `Json` object
// after unmarshaling `bytes`
func FromBytes(b []byte) (*Json, error) {
	return FromReader(bytes.NewReader(b))
}

// FromFile returns a pointer to a new `Json` object
// after unmarshaling the contents from `file` into it
func FromFile(file string) (*Json, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return FromBytes(data)
}

// FromReader returns a *Json by decoding from an io.Reader
func FromReader(r io.Reader) (*Json, error) {
	if r == nil {
		return FromString("null")
	}
	rc, ok := r.(io.ReadCloser)
	if !ok {
		rc = ioutil.NopCloser(r)
	}
	return FromReadCloser(rc)
}

// FromReadCloser returns a *Json by decoding from an io.ReadCloser and calls the io.ReadCloser Close method
func FromReadCloser(rc io.ReadCloser) (*Json, error) {
	if rc == nil {
		return FromString("null")
	}
	defer rc.Close()
	j := &Json{}
	dec := json.NewDecoder(rc)
	dec.UseNumber()
	err := dec.Decode(&j.data)
	return j, err
}

// ToBytes returns its marshaled data as `[]byte`
func (j *Json) ToBytes() ([]byte, error) {
	return j.MarshalJSON()
}

// ToString returns its marshaled data as `string`
func (j *Json) ToString() (string, error) {
	b, err := j.ToBytes()
	return string(b), err
}

// ToPrettyBytes returns its marshaled data as `[]byte` with indentation
func (j *Json) ToPrettyBytes() ([]byte, error) {
	return json.MarshalIndent(&j.data, "", "  ")
}

// ToPrettyString returns its marshaled data as `string` with indentation
func (j *Json) ToPrettyString() (string, error) {
	b, err := j.ToPrettyBytes()
	return string(b), err
}

// ToFile writes the Json to the `file` with permission `perm`
func (j *Json) ToFile(file string, perm os.FileMode) error {
	b, _ := j.ToBytes()
	return ioutil.WriteFile(file, b, perm)
}

// ToReader returns its marshaled data as `io.Reader`
func (j *Json) ToReader() (io.Reader, error) {
	b, err := j.ToBytes()
	r := bytes.NewReader(b)
	return r, err
}

// Implements the json.Marshaler interface.
func (j *Json) MarshalJSON() ([]byte, error) {
	return json.Marshal(&j.data)
}

// Implements the json.Unmarshaler interface.
func (j *Json) UnmarshalJSON(p []byte) (error) {
	jNew, err := FromReader(bytes.NewReader(p))
	j.data = jNew.data
	return err
}

// Get searches for the item as specified by the path.
// path can contain strings or ints to navigate through json
// objects and arrays. If the given path is not present then
// the deepest valid value is returned along with an error.
//
//   js.Get("top_level", "dict", 3, "foo")
func (j *Json) Get(path ...interface{}) (*Json, *jsonPathError) {
	tmp := j
	for i, k := range path {
		if key, ok := k.(string); ok {
			if m, err := tmp.Map(); err == nil {
				if val, ok := m[key]; ok {
					tmp = &Json{val}
				} else {
					return tmp, &jsonPathError{path[:i], path[i:]}
				}
			} else {
				return tmp, &jsonPathError{path[:i], path[i:]}
			}
		} else if index, ok := k.(int); ok {
			if a, err := tmp.Array(); err == nil {
				if index < 0 || index >= len(a) {
					return tmp, &jsonPathError{path[:i], path[i:]}
				} else {
					tmp = &Json{a[index]}
				}
			} else {
				return tmp, &jsonPathError{path[:i], path[i:]}
			}
		} else {
			return tmp, &jsonPathError{path[:i], path[i:]}
		}
	}
	return tmp, nil
}

// Set modifies `Json`, recursively checking/creating map keys and checking
// array indices for the supplied path, and then finally writing in the value.
// Set will only create maps where the current map[key] does not exist,
// if the key exists, even if the value is nil, a new map will not be created and an
// error wil be returned.
func (j *Json) Set(val interface{}, path ...interface{}) *jsonPathError {
	if len(path) == 0 {
		j.data = val
		return nil
	}

	tmp := j

	for i := 0; i < len(path); i++ {
		if key, ok := path[i].(string); ok {
			if m, err := tmp.Map(); err == nil {
				if i == len(path)-1 {
					m[key] = val
				} else {
					_, ok := path[i+1].(string)
					_, exists := m[key]
					if ok && !exists {
						m[key] = map[string]interface{}{}
					}
					tmp = &Json{m[key]}
				}
			} else {
				return &jsonPathError{path[:i], path[i:]}
			}
		} else if index, ok := path[i].(int); ok {
			if a, err := tmp.Array(); err == nil && index >= 0 && index < len(a) {
				if i == len(path)-1 {
					a[index] = val
				} else {
					tmp = &Json{a[index]}
				}
			} else {
				return &jsonPathError{path[:i], path[i:]}
			}
		} else {
			return &jsonPathError{path[:i], path[i:]}
		}
	}

	return nil
}

// Del modifies `Json` maps and arrays by deleting/removing the last `path` segment if it is present,
func (j *Json) Del(path ...interface{}) *jsonPathError {
	if len(path) == 0 {
		j.data = nil
		return nil
	}

	i := len(path)-1
	tmp, err := j.Get(path[:i]...)
	if err != nil {
		err.MissingPath = append(err.MissingPath, path[i])
		return err
	}

	if key, ok := path[i].(string); ok {
		if m, err := tmp.Map(); err != nil {
			return &jsonPathError{path[:i], path[i:]}
		} else {
			delete(m, key)
		}
	} else if index, ok := path[i].(int); ok {
		if a, err := tmp.Array(); err != nil {
			return &jsonPathError{path[:i], path[i:]}
		} else if index < 0 || index >= len(a) {
			return &jsonPathError{path[:i], path[i:]}
		} else {
			a, a[len(a)-1] = append(a[:index], a[index+1:]...), nil
			if i == 0 {
				j.data = a
			} else {
				tmp, _ = j.Get(path[:i-1]...)
				if key, ok := path[i-1].(string); ok {
					tmp.MustMap(nil)[key] = a //is this safe? should be 100% certainty ;)
				} else if index, ok := path[i-1].(int); ok {
					tmp.MustArray(nil)[index] = a //is this safe? should be 100% certainty ;)
				}
			}
		}
	} else {
		return &jsonPathError{path[:i], path[i:]}
	}
	return nil
}

// Interface returns the underlying data
func (j *Json) Interface(path ...interface{}) (interface{}, *jsonPathError) {
	tmp, err := j.Get(path...)
	return tmp.data, err
}

// Map type asserts to `map`
func (j *Json) Map(path ...interface{}) (map[string]interface{}, error) {
	tmp, err := j.Get(path...)
	if err != nil {
		return nil, err
	}
	if m, ok := tmp.data.(map[string]interface{}); ok {
		return m, nil
	}
	return nil, errors.New("type assertion to map[string]interface{} failed")
}

// MustMap guarantees the return of a `map[string]interface{}` (with specified default)
//
// useful when you want to iterate over map values in a succinct manner:
//		for k, v := range js.MustMap(nil) {
//			fmt.Println(k, v)
//		}
func (j *Json) MustMap(def map[string]interface{}, path ...interface{}) map[string]interface{} {
	if a, err := j.Map(path...); err == nil {
		return a
	}
	return def
}

// Array type asserts to an `array`
func (j *Json) Array(path ...interface{}) ([]interface{}, error) {
	tmp, err := j.Get(path...)
	if err != nil {
		return nil, err
	}
	if a, ok := tmp.data.([]interface{}); ok {
		return a, nil
	}
	return nil, errors.New("type assertion to []interface{} failed")
}

// MustArray guarantees the return of a `[]interface{}` (with specified default)
//
// useful when you want to iterate over array values in a succinct manner:
//		for i, v := range js.MustArray(nil) {
//			fmt.Println(i, v)
//		}
func (j *Json) MustArray(def []interface{}, path ...interface{}) []interface{} {
	if a, err := j.Array(path...); err == nil {
		return a
	}
	return def
}

// Bool type asserts to `bool`
func (j *Json) Bool(path ...interface{}) (bool, error) {
	tmp, err := j.Get(path...)
	if err != nil {
		return false, err
	}
	if s, ok := tmp.data.(bool); ok {
		return s, nil
	}
	return false, errors.New("type assertion to bool failed")
}

// MustBool guarantees the return of a `bool` (with specified default)
//
// useful when you explicitly want a `bool` in a single value return context:
//     myFunc(js.MustBool(true))
func (j *Json) MustBool(def bool, path ...interface{}) bool {
	if b, err := j.Bool(path...); err == nil {
		return b
	}
	return def
}

// String type asserts to `string`
func (j *Json) String(path ...interface{}) (string, error) {
	tmp, err := j.Get(path...)
	if err != nil {
		return "", err
	}
	if s, ok := tmp.data.(string); ok {
		return s, nil
	}
	return "", errors.New("type assertion to string failed")
}

// MustString guarantees the return of a `string` (with specified default)
//
// useful when you explicitly want a `string` in a single value return context:
//     myFunc(js.MustString("my_default"))
func (j *Json) MustString(def string, path ...interface{}) string {
	if s, err := j.String(path...); err == nil {
		return s
	}
	return def
}

// StringArray type asserts to an `array` of `string`
func (j *Json) StringArray(path ...interface{}) ([]string, error) {
	arr, err := j.Array(path...)
	if err != nil {
		return nil, err
	}
	retArr := make([]string, 0, len(arr))
	for _, a := range arr {
		if s, ok := a.(string); a == nil || !ok {
			return nil, errors.New("none string value encountered")
		}else {
			retArr = append(retArr, s)
		}
	}
	return retArr, nil
}

// MustStringArray guarantees the return of a `[]string` (with specified default)
//
// useful when you want to iterate over array values in a succinct manner:
//		for i, s := range js.MustStringArray(nil) {
//			fmt.Println(i, s)
//		}
func (j *Json) MustStringArray(def []string, path ...interface{}) []string {
	if a, err := j.StringArray(path...); err == nil {
		return a
	}
	return def
}

// Int coerces into an int
func (j *Json) Int(path ...interface{}) (int, error) {
	f, err := j.Float64(path...)
	return int(f), err
}

// MustInt guarantees the return of an `int` (with specified default)
//
// useful when you explicitly want an `int` in a single value return context:
//     myFunc(js.MustInt(5150))
func (j *Json) MustInt(def int, path ...interface{}) int {
	if i, err := j.Int(path...); err == nil {
		return i
	}
	return def
}

// IntArray type asserts to an `array` of `int`
func (j *Json) IntArray(path ...interface{}) ([]int, error) {
	arr, err := j.Array(path...)
	if err != nil {
		return nil, err
	}
	retArr := make([]int, 0, len(arr))
	for _, a := range arr {
		tmp := &Json{a}
		if i, err := tmp.Int(); err != nil {
			return nil, err
		} else {
			retArr = append(retArr, i)
		}
	}
	return retArr, nil
}

// MustIntArray guarantees the return of a `[]int` (with specified default)
//
// useful when you want to iterate over array values in a succinct manner:
//		for i, s := range js.MustIntArray(nil) {
//			fmt.Println(i, s)
//		}
func (j *Json) MustIntArray(def []int, path ...interface{}) []int {
	if a, err := j.IntArray(path...); err == nil {
		return a
	}
	return def
}

// Float64 coerces into a float64
func (j *Json) Float64(path ...interface{}) (float64, error) {
	tmp, err := j.Get(path...)
	if err != nil {
		return 0, err
	}
	switch tmp.data.(type) {
	case json.Number:
		return tmp.data.(json.Number).Float64()
	case float32, float64:
		return reflect.ValueOf(tmp.data).Float(), nil
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(tmp.data).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return float64(reflect.ValueOf(tmp.data).Uint()), nil
	}
	return 0, errors.New("invalid value type")
}

// MustFloat64 guarantees the return of a `float64` (with specified default)
//
// useful when you explicitly want a `float64` in a single value return context:
//     myFunc(js.MustFloat64(5.150))
func (j *Json) MustFloat64(def float64, path ...interface{}) float64 {
	if f, err := j.Float64(path...); err == nil {
		return f
	}
	return def
}

// Float64Array type asserts to an `array` of `float64`
func (j *Json) Float64Array(path ...interface{}) ([]float64, error) {
	arr, err := j.Array(path...)
	if err != nil {
		return nil, err
	}
	retArr := make([]float64, 0, len(arr))
	for _, a := range arr {
		tmp := &Json{a}
		if f, err := tmp.Float64(); err != nil {
			return nil, err
		} else {
			retArr = append(retArr, f)
		}
	}
	return retArr, nil
}

// MustFloat64Array guarantees the return of a `[]float64` (with specified default)
//
// useful when you want to iterate over array values in a succinct manner:
//		for i, s := range js.MustFloat64Array(nil) {
//			fmt.Println(i, s)
//		}
func (j *Json) MustFloat64Array(def []float64, path ...interface{}) []float64 {
	if a, err := j.Float64Array(path...); err == nil {
		return a
	}
	return def
}

// Int64 coerces into an int64
func (j *Json) Int64(path ...interface{}) (int64, error) {
	tmp, err := j.Get(path...)
	if err != nil {
		return 0, err
	}
	switch tmp.data.(type) {
	case json.Number:
		return tmp.data.(json.Number).Int64()
	case float32, float64:
		return int64(reflect.ValueOf(tmp.data).Float()), nil
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(tmp.data).Int(), nil
	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(tmp.data).Uint()), nil
	}
	return 0, errors.New("invalid value type")
}

// MustInt64 guarantees the return of an `int64` (with specified default)
//
// useful when you explicitly want an `int64` in a single value return context:
//     myFunc(js.MustInt64(5150))
func (j *Json) MustInt64(def int64, path ...interface{}) int64 {
	if i, err := j.Int64(path...); err == nil {
		return i
	}
	return def
}

// Int64Array type asserts to an `array` of `int64`
func (j *Json) Int64Array(path ...interface{}) ([]int64, error) {
	arr, err := j.Array(path...)
	if err != nil {
		return nil, err
	}
	retArr := make([]int64, 0, len(arr))
	for _, a := range arr {
		tmp := &Json{a}
		if i, err := tmp.Int64(); err != nil {
			return nil, err
		} else {
			retArr = append(retArr, i)
		}
	}
	return retArr, nil
}

// MustInt64Array guarantees the return of a `[]int64` (with specified default)
//
// useful when you want to iterate over array values in a succinct manner:
//		for i, s := range js.MustInt64Array(nil) {
//			fmt.Println(i, s)
//		}
func (j *Json) MustInt64Array(def []int64, path ... interface{}) []int64 {
	if a, err := j.Int64Array(path...); err == nil {
		return a
	}
	return def
}

// Uint64 coerces into an uint64
func (j *Json) Uint64(path ...interface{}) (uint64, error) {
	tmp, err := j.Get(path...)
	if err != nil {
		return 0, err
	}
	switch tmp.data.(type) {
	case json.Number:
		return strconv.ParseUint(tmp.data.(json.Number).String(), 10, 64)
	case float32, float64:
		return uint64(reflect.ValueOf(tmp.data).Float()), nil
	case int, int8, int16, int32, int64:
		return uint64(reflect.ValueOf(tmp.data).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(tmp.data).Uint(), nil
	}
	return 0, errors.New("invalid value type")
}

// MustUInt64 guarantees the return of an `uint64` (with specified default)
//
// useful when you explicitly want an `uint64` in a single value return context:
//     myFunc(js.MustUint64(5150))
func (j *Json) MustUint64(def uint64, path ...interface{}) uint64 {
	if i, err := j.Uint64(path...); err == nil {
		return i
	}
	return def
}

// Uint64Array type asserts to an `array` of `uint64`
func (j *Json) Uint64Array(path ...interface{}) ([]uint64, error) {
	arr, err := j.Array(path...)
	if err != nil {
		return nil, err
	}
	retArr := make([]uint64, 0, len(arr))
	for _, a := range arr {
		tmp := &Json{a}
		if u, err := tmp.Uint64(); err != nil {
			return nil, err
		} else {
			retArr = append(retArr, u)
		}
	}
	return retArr, nil
}

// MustUint64Array guarantees the return of a `[]uint64` (with specified default)
//
// useful when you want to iterate over array values in a succinct manner:
//		for i, s := range js.MustUint64Array(nil) {
//			fmt.Println(i, s)
//		}
func (j *Json) MustUint64Array(def []uint64, path ...interface{}) []uint64 {
	if a, err := j.Uint64Array(path...); err == nil {
		return a
	}
	return def
}

type jsonPathError struct {
	FoundPath   []interface{}
	MissingPath []interface{}
}

func (e *jsonPathError) Error() string {
	return fmt.Sprintf("found: %v missing: %v", e.FoundPath, e.MissingPath)
}
