package main

// montitor a folder for inotify events

import (
	"crypto/md5"
	"io"
	"log"
	"math"
	"os"
	"time"

	"github.com/howeyc/fsnotify"
)

const filechunk = 8192

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

// checksum a file
func checksum(file string) (sum string, err error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	info, _ := f.Stat()
	filesize := info.Size()
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))
	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)

		f.Read(buf)
		io.WriteString(hash, string(buf)) // append into the hash
	}

	return string(hash.Sum(nil)), nil
}

// START OMIT
// process a newly created file
func process(file string) error {
	//	wait until upload finishes
	sum, err := checksum(file) // HL
	if err != nil {
		log.Println(label.Error + err.Error())
		return err
	}
	for {
		time.Sleep(time.Duration(1) * time.Second) // HL
		currentSum, err := checksum(file)
		if err != nil {
			log.Println(label.Error + err.Error())
			return err
		}
		if currentSum == sum {
			break
		}
		sum = currentSum
	}
	// wait finished
	// do something useful here..
	return nil
}

// END OMIT

// watch a dir for changes
func watch(dir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	done := make(chan interface{})
	go func() {
		for {
			select {
			case event := <-watcher.Event:
				if event.IsCreate() {
					log.Println(label.Info + " New file: " + event.Name)
					process(event.Name)
					log.Println(label.Info + " Operation successful for: " + event.Name)
				}
			case err := <-watcher.Error:
				log.Println(err)
			}
		}
	}()

	watcher.Watch(dir)
	log.Println(label.Info + " Startup done. Waiting for events")
	<-done

	return nil
}

func main() {
	watch(dir)
}
