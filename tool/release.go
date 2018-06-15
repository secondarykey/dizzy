package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"bufio"
	"github.com/fsnotify/fsnotify"
)

func init() {
}

func main() {

	var err error
	ci := flag.Bool("ci", false, "go test,release circuit integration flag")
	flag.Parse()

	if *ci {
		err = circuit()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		//出力
		err = gen()
		if err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(0)
}

var lastUnixTime int64

func lock() bool {

	now := time.Now().Unix()
	dur := int64(1)
	if now > lastUnixTime + dur {
		return true
	}
	return false
}

func unlock() {
	lastUnixTime = time.Now().Unix()
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

				if lock() {
					err := runTest()
					if err != nil {
						done <- err
					}
					unlock()
				}

			case err := <-watcher.Errors:
				done <- err
			}
		}
	}()

	//手入力終了待ち受け
	go func() {
		stdin := bufio.NewScanner(os.Stdin)
		stdin.Scan()
		os.Exit(0)
	}()

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

func gen() error {

	//テストなどで出力してある各ファイルを削除
	file, err := os.Create("template.go")
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, err := template.ParseFiles("templates/template.tmpl")
	if err != nil {
		return err
	}

	app := getFile("templates/app.tmpl")
	dizzy := getFile("templates/dizzy.tmpl")
	edit := getFile("templates/edit.tmpl")
	errView := getFile("templates/error.tmpl")
	gen := getFile("templates/access.tmpl")
	handler := getFile("templates/handler.tmpl")
	index := getFile("templates/index.tmpl")
	layout := getFile("templates/layout.tmpl")
	top := getFile("templates/top.tmpl")
	view := getFile("templates/view.tmpl")

	dto := struct {
		Flag           bool
		Created        time.Time
		DizzyGo        string
		GenGo          string
		HandlerGo      string
		LayoutTemplate string
		TopTemplate    string
		EditTemplate   string
		ViewTemplate   string
		ErrorTemplate  string
		AppTemplate    string
		IndexTemplate  string
	}{false, time.Now(),
		string(dizzy), string(gen), string(handler),
		string(layout), string(top), string(edit), string(view), string(errView),
		string(app), string(index)}

	err = tmpl.Execute(file, dto)
	if err != nil {
		return err
	}

	fmt.Println("Generated:")
	fmt.Println("\ttemplate.go")

	//examples/templates,examples/*_gen.go,examples/*_handler,yamlを削除
	removeWork()

	fmt.Println("Complete.")

	return nil
}

func getFile(name string) []byte {
	dizzy, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	return dizzy
}

func removeWork() {

	fmt.Println("Work File remove")

	work := "examples/"

	p := work + "templates"
	//削除
	err := os.RemoveAll(p)
	if err != nil {
		fmt.Println("\t" + p)
	} else {
		fmt.Println("\tNot exists:" + p)
	}

	remove("sample_access.go")
	remove("sample_handler.go")
	remove("dizzy.go")
	remove("dizzy_app.yaml")
	remove("index.yaml")
}

func remove(f string) {
	work := "examples/"
	p := work + f
	err := os.Remove(p)
	if err != nil {
		fmt.Println("\t" + p)
	} else {
		fmt.Println("\tNot exists:" + p)
	}
}
