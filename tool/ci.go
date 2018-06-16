package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"bufio"
	"github.com/fsnotify/fsnotify"
)

func init() {
}

func main() {

	err := circuit()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

var lastUnixTime int64

func lock() bool {

	now := time.Now().Unix()
	dur := int64(2)
	if now > lastUnixTime + dur {
		return true
	}

	fmt.Println(now)
	return false
}

func unlock() {
	now := time.Now()
	fmt.Println(now)
	lastUnixTime = now.Unix()
}

func circuit() error {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	watcher.Add("./")
	if err != nil {
		return nil
	}

	done := make(chan error)
	ci := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				switch {
				case event.Op&fsnotify.Rename == fsnotify.Rename:
				case event.Op&fsnotify.Create == fsnotify.Create:
				case event.Op&fsnotify.Remove == fsnotify.Remove:
				case event.Op&fsnotify.Chmod == fsnotify.Chmod:
				case event.Op&fsnotify.Write == fsnotify.Write:
				default:
				}


				fname := event.Name
				if fname == "template.go" {
					continue
				}

				f,err := os.Stat(fname)
				if err != nil {
					continue
				}

				if f.IsDir() {
					continue
				}

				ci <- true

			case err := <-watcher.Errors:
				done <- err
			}
		}
	}()

	//手入力終了待ち受け
	go func() {
		stdin := bufio.NewScanner(os.Stdin)

		for {
			stdin.Scan()
			if stdin.Text() == "quit" {
				os.Exit(0)
			} else {
				ci <- true
			}
		}

	}()

	go func() {
		for <-ci {
			if lock() {
				err := runTest()
				if err != nil {
					done <- err
				}
				unlock()
			}
		}
	}()

	//一回処理
	ci <- true

	return <-done
}

func runTest() error {

	log.Println("\x1b[36m######################################################\x1b[0m")
	log.Println("\x1b[36m# Test [enter -> quit]\x1b[0m")
	log.Println("\x1b[36m######################################################\x1b[0m")

	cmd := exec.Command("go", "test", "-v",  ".")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return err
	}
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		txt := scanner.Text()
		color := 37

		if strings.Index(txt, "FAIL") != -1 {
			color = 31
		} else if strings.Index(txt,"PASS") != -1 ||
			 strings.Index(txt,"ok") != -1 {
			color = 32
		} else if strings.Index(txt, "RUN") != -1 {
			color = 35
		}

		fmt.Printf("\x1b[%dm%s\x1b[0m\n",color,txt)
	}
	cmd.Wait()

	return nil
}

