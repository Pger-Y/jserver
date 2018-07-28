package floader

import (
	"encoding/json"
	//"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"time"
)

func NewFloader(fname string) *Floader {
	sig := make(chan struct{}, 10)
	return &Floader{fname: fname, sig: sig}
}

type Floader struct {
	fname string
	time  time.Time
	data  map[string]interface{}
	sig   chan struct{}
}

func (fl *Floader) needsUpdate() {
	fl.sig <- struct{}{}
}

func (fl *Floader) shouldUpdate() {
	<-fl.sig
}

func (fl *Floader) loadFile() {
	data, err := ioutil.ReadFile(fl.fname)
	if err != nil {
		log.Printf("Read file[%s] error:[%s]\n", fl.fname, err.Error())
		return
	}

	fl.data = make(map[string]interface{})
	err = json.Unmarshal(data, &fl.data)
	if err != nil {
		log.Printf("Check json error[%s],content:%s\n", err.Error(), string(data))
		return
	}

	//format := "2006-01-02 15:04:05"
	nt := time.Now()
	format := time.RFC3339
	log.Printf("Update time:%s,last time:%s", nt.Format(format), fl.time.Format(format))
	fl.time = nt

	fl.needsUpdate()
}

func (fl *Floader) Sync() map[string]interface{} {
	fl.shouldUpdate()
	return fl.data
}

/*
func checkJson(data []byte)error{
	var data map[string]interface{}
	err := json.Unmarshal(fl.data,&data)
	return err
}
*/

func (fl *Floader) Run() {
	fl.loadFile()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			need_reload := false
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					need_reload = true
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("reload file:", event.Name)
					watcher.Remove(fl.fname)
					err = watcher.Add(fl.fname)
					if err != nil {
						log.Println("remove and readd file into watcher error:", err)
					} else {
						need_reload = true
					}
				}
				if need_reload {
					fl.loadFile()
				}

			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(fl.fname)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
