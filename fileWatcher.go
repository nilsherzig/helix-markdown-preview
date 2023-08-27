package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func watchFolder(folderPath string, channel chan string) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if filepath.Ext(event.Name) == ".md" {
					if event.Has(fsnotify.Write) {
						dat, err := os.ReadFile(event.Name)
						if err != nil {
							panic(err)
						}
						channel <- string(dat)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// get all markdown files from root folder (recursion)
	err = filepath.Walk(folderPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".md" {
				fmt.Println(path, info.Size())
				// Add a path.
				err = watcher.Add(path)
				if err != nil {
					log.Fatal(err)
				}
			}
			return nil
		})

	handleErr(err)

	// Block main goroutine forever.
	<-make(chan struct{})

}
