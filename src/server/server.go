package main

import (
	"github.com/robsix/3ditor/src/server/Godeps/_workspace/src/github.com/robsix/golog"
	"github.com/robsix/3ditor/src/server/Godeps/_workspace/src/github.com/robsix/json"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	log := golog.NewConsoleLog(0)
	wd, _ := os.Getwd()
	conf, _ := json.FromFile(filepath.Join(wd, "conf.json"))
	publicDir := filepath.Join(append([]string{wd}, conf.MustStringArray([]string{"..", "client"}, "publicDir")...)...)

	log.Info("serving static files from: ", publicDir)
	fileServer := http.FileServer(http.Dir(publicDir))
	http.Handle(`/`, fileServer)

	log.Info("server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
