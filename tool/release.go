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
	"sync/atomic"
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
		err = gen()
		if err != nil {
			log.Fatal(err)
		}
	}

}

var lockVal int64

func lock() bool {
	atomic.AddInt64(&lockVal, 1)
	if lockVal == 1 {
		return true
	}
	log.Println(lockVal)
	return false
}

func unlock() {
	atomic.AddInt64(&lockVal, -1)
}

func circuit() error {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = register(watcher, "./", ".go")
	if err != nil {
		return nil
	}
	err = register(watcher, "./templates/", ".tmpl")
	if err != nil {
		return nil
	}

	done := make(chan error)
	go func() {
		for {
			select {
			case event := <-watcher.Events:

				switch {
				case event.Op&fsnotify.Rename == fsnotify.Rename,
					event.Op&fsnotify.Create == fsnotify.Create,
					event.Op&fsnotify.Remove == fsnotify.Remove,
					event.Op&fsnotify.Write == fsnotify.Write:

					if lock() {
						if strings.Index(event.Name, ".go") != -1 {
							runTest()
						} else if strings.Index(event.Name, ".tmpl") != -1 {
							runRelease()
							runTest()
						}
					}
					unlock()
					watcher.Add(event.Name)

				}
			case err := <-watcher.Errors:
				done <- err
			}
		}
	}()

	go func() {
		stdin := bufio.NewScanner(os.Stdin)
		stdin.Scan()
		os.Exit(0)
	}()

	return <-done
}

func runTest() {
	log.Println("#### Run Test")
	cmd := exec.Command("go", "test", "-v", "-count=1", ".")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	cmd.Wait()
}

func runRelease() {
	log.Println("#### Run template generate")
	gen()
}

func register(watcher *fsnotify.Watcher, dir string, ext string) error {
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		n := info.Name()
		if n == "template.go" {
			continue
		}
		idx := strings.Index(n, ext)
		if idx != -1 {
			a := dir + n
			err = watcher.Add(a)
			if err != nil {
				return err
			}
		}
	}
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

	fmt.Println("Work remove")

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
