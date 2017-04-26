package main

// montitor a folder for inotify events

import (
	"log"
	"time"

	"github.com/howeyc/fsnotify"
)

var (
	dir   = "./"
	label = labelType{
		Info:  "[ \x1b[33mInfo\x1b[0m ]",
		Error: "[ \x1b[31mError\x1b[0m ]",
		Fatal: "[ \x1b[31mFatal\x1b[0m ]",
	}
)

// LabelType for pretty console output
type labelType struct {
	Fatal, Error, Info string
}

// process a newly created file
func process(file string) error {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Event:
				if event.IsModify() {
					continue
				}
			case <-time.After(250 * time.Millisecond):
				done <- true
				return
			}
		}
	}()
	watcher.Watch(file)
	<-done

	// do something useful here

	return nil
}

// watch a dir for changes
// START OMIT
func watch(dir string) {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()
	done := make(chan interface{})
	go func() {
		for {
			select { // HL
			case event := <-watcher.Event: // HL
				if event.IsCreate() {
					log.Println(label.Info + " New file: " + event.Name)
					process(event.Name)
					log.Println(label.Info + " Operation successful for: " + event.Name)
				}
			case err := <-watcher.Error: // HL
				log.Println(err)
			} // HL
		}
	}()

	watcher.Watch(dir)
	log.Println(label.Info + " Startup done. Waiting for events")
	<-done
}

// END OMIT

func main() {
	watch(dir)
}
