package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"os/exec"
	"bufio"
)

func main() {

	var err error
	sub := ""
	if len(os.Args) > 1 {
		sub = os.Args[1]
	}

	//ds
	//form

	switch sub {
	case "gen":
		if len(os.Args) != 3 {
			err = fmt.Errorf("usage:dizzy gen dirname")
			break
		}
		dir := os.Args[2]
		err = gen(dir)
	case "dev":
		err = dev()
	case "deploy":
		err = deploy()
	default:
		fmt.Println("usage: dizzy <command> [<args>]")
		fmt.Println("\tgen: generate datastore source(arg file)")
		fmt.Println("\thandler: generate appengine template")
		fmt.Println("\tdev: dev dev_appserver.py tool")
		fmt.Println("\tdeploy: gcloud app deploy tool")
		return
	}

	result(err)
}

func gen(dir string) error {

	//ディレクトリの場合
	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("os.Stat() error[%s]", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("second argument directory name[%s]", dir)
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("ioutil.ReadDir() error[%s]", dir)
	}

	pkg := info.Name()
	var all []*Kind
	//ファイル数回繰り返す
	for _, f := range files {

		if f.IsDir() {
			continue
		}

		file := f.Name()
		idx := strings.LastIndex(file, ".go")
		if idx == -1 {
			continue
		}

		in := filepath.Join(dir, file)

		kinds, packageName, err := createKinds(in)
		if err != nil {
			return err
		}

		//存在しなかった場合
		if len(kinds) == 0 {
			continue
		}

		all = append(all, kinds...)
		pkg = packageName

		//accessファイルを出力
		err = generateAccessFile(in, packageName, kinds)
		if err != nil {
			return err
		}

		//handlerファイルを出力
		//Goファイルを作成
		err = generateHandlerFile(in, packageName, kinds)
		if err != nil {
			return err
		}

		//テンプレートファイルを作成
		err = generateTemplateFiles(in, kinds)
		if err != nil {
			return err
		}
	}

	//アクセス用のポータルファイルを作成
	err = generateRootTemplateFiles(dir, all)
	if err != nil {
		return err
	}

	//AppEngine用のファイルを作成
	err = generateAppEngineFiles(dir, all)
	if err != nil {
		return err
	}

	//１回あればいいgoファイル(dizzy.go)の作成
	err = generateDizzyFile(dir, pkg, all)
	if err != nil {
		return err
	}

	return nil
}

func dev() error {

	//単純に「dev_appserver.py dizzy_app.yaml」を行う
	cmd := exec.Command("dev_appserver.py","dizzy_app.yaml")
	stdout,err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	cmd.Wait()
	return nil
}

func deploy() error {

	//単純に「gcloud app deploy dizzy_app.yaml dizzy_index.yaml」を行う

	return fmt.Errorf("Not implemented.")
}

//結果表示
func result(err error) {

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Success.")
	if genFiles == nil {
		fmt.Println("\tNo generate.")
	} else {
		fmt.Println("GenerateFiles : ")
		for _, elm := range genFiles {
			fmt.Println("\t" + elm)
		}
	}

	os.Exit(0)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
