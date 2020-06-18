package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os/exec"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	err = watcher.Add("./docs")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Fatal("Something went wrong with watcher.Events %s", err)
				}

				fmt.Print(event.Op)

				if event.Op == fsnotify.Create {
					out, err := exec.Command("soffice", "--convert-to", "pdf", "./docs/menu.doc").Output()

					if err != nil {
						fmt.Printf("Error", out)
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

	<-done
}
