package main

import (
	"log"
	"github.com/howeyc/fsnotify"
	"strings"
	"os/exec"
	"fmt"
	"bytes"
)

func main() {
	watcher, err := fsnotify.NewWatcher()

	done := make(chan bool)

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				// through jetbrains temporary files
				if strings.Contains(ev.Name, "___") {
					break
				}

				// through myselfoddog69

				if strings.Contains(ev.Name, "main.go") {
					break
				}

				if ev.IsModify() {
					log.Printf("Run: %v", ev.Name)

					cmd := exec.Command("go", "run", ev.Name)
					var out bytes.Buffer
					var stderr bytes.Buffer
					cmd.Stdout = &out
					cmd.Stderr = &stderr
					err := cmd.Run()

					if err != nil {
						log.Println(fmt.Sprint(err) + ": " + stderr.String())
						break
					}
					log.Print(stderr.String())
					log.Print(out.String())
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch(".")

	if err != nil {
		log.Fatal(err)
	}

	// Hang so program doesn't exit
	<-done

	watcher.Close()
}
