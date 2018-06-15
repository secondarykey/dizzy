package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"
	"time"

)

func init() {
}

func main() {

	//出力
	err := gen()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
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
