package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const KindAccessSuffix = "_access.go"
const KindHandlerSuffix = "_handler.go"
const DizzyGlobalFile = "dizzy.go"

var genFiles []string

func addGeneratedFile(name string) {
	genFiles = append(genFiles, name)
}

//ファイル出力
//.goファイルはfmtを行う
func generate(dst, tmplName string, dto interface{}) error {

	tmpl, err := createGenerateTemplate(tmplName)
	if err != nil {
		return err
	}

	gofmt := false
	// go files go fmt generate
	if strings.LastIndex(dst, ".go") != -1 {
		gofmt = true
	}

	var writer io.Writer
	//Goファイルの場合,バッファで一度取得
	if gofmt {
		writer = &bytes.Buffer{}
	} else {
		closer, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer closer.Close()
		writer = closer
	}

	err = tmpl.Execute(writer, dto)
	if err != nil {
		return err
	}

	//Goファイルの場合
	if byt, ok := writer.(*bytes.Buffer); ok {
		dstByt, err := format.Source(byt.Bytes())
		if err != nil {
			return err
		}

		newW, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer newW.Close()

		_, err = newW.Write(dstByt)
		if err != nil {
			return err
		}
	}

	//出力ファイルに追加
	addGeneratedFile(dst)
	return nil
}

func generateAccessFile(src string, packageName string, kinds []*Kind) error {

	idx := strings.LastIndex(src, ".go")
	if idx == -1 {
		return fmt.Errorf("[%s] is not go file", src)
	}

	dst := src[:idx] + KindAccessSuffix
	//取得したstructから出力ソースを生成
	dto := struct {
		PackageName string
		Kinds       []*Kind
	}{packageName, kinds}

	//generate _access.go
	err := generate(dst, AccessTemplateFile, dto)
	if err != nil {
		return err
	}
	return nil
}

func generateHandlerFile(src, pname string, kinds []*Kind) error {

	idx := strings.LastIndex(src, ".go")
	if idx == -1 {
		return fmt.Errorf("Not Found go file[%s]\n", src)
	}

	dstGo := src[:idx] + KindHandlerSuffix

	dto := struct {
		PackageName string
		Kinds       []*Kind
	}{pname, kinds}
	err := generate(dstGo, HandlerGoTemplateFile, dto)
	if err != nil {
		return err
	}
	return nil
}

func generateTemplateFiles(src string, kinds []*Kind) error {

	dir := "templates"

	pwd := filepath.Dir(src)
	mkdir := pwd + "/" + dir

	if _, err := os.Stat(mkdir); os.IsNotExist(err) {
		err := os.Mkdir(mkdir, 0777)
		if err != nil {
			return err
		}
	}

	var err error
	for _, elm := range kinds {

		dto := struct {
			Kind    *Kind
			Created time.Time
		}{elm, time.Now()}

		viewName := "view_" + elm.KindName + ".tmpl"
		err = generate(mkdir+"/"+viewName, ViewTemplateFile, dto)
		if err != nil {
			return err
		}

		editName := "edit_" + elm.KindName + ".tmpl"
		err = generate(mkdir+"/"+editName, EditTemplateFile, dto)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateRootTemplateFiles(dst string, kinds []*Kind) error {

	var err error
	dto := struct {
		Kinds []*Kind
	}{kinds}

	dst += "/templates"

	err = generate(dst+"/layout.tmpl", LayoutTemplateFile, dto)
	if err != nil {
		return err
	}

	err = generate(dst+"/top.tmpl", TopTemplateFile, dto)
	if err != nil {
		return err
	}

	err = generate(dst+"/error.tmpl", ErrorTemplateFile, dto)
	if err != nil {
		return err
	}
	return nil
}

func generateAppEngineFiles(dir string, kinds []*Kind) error {

	dto := struct {
		Kinds []*Kind
	}{kinds}
	err := generate(dir+"/dizzy_app.yaml", AppTemplateFile, dto)
	if err != nil {
		return err
	}

	err = generate(dir+"/index.yaml", IndexTemplateFile, dto)
	return err
}

func generateDizzyFile(dir, pname string, kinds []*Kind) error {

	dto := struct {
		PackageName string
		Kinds       []*Kind
	}{pname, kinds}

	err := generate(dir+"/"+DizzyGlobalFile, DizzyGoTemplateFile, dto)

	return err
}
