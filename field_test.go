package main

import (
	"go/ast"
	"testing"
	"go/token"
	"go/parser"
)

func TestCreateFields(t *testing.T) {

	var f []*Field
	var err error
	//件数
	str := getStringStruct(t,StringTestStruct)
	if str == nil {
		t.Fatal("getTestStruct() Error.")
		return
	}

	f,err = createFields(str)
	//非サポートデータ
	if err != nil {
		t.Fatal("TestStruct Fields Error.")
	}

	if f[0].Name != "TestField1" {
		t.Fatalf("TestStruct Name Fields Error.[%s]",f[0].Name)
	}

	if f[0].Type != IntField {
		t.Fatalf("TestStruct Type Error.[%d]",f[0].Type)
	}

	if f[0].TypeName != "int" {
		t.Fatalf("TestStruct TypeName Field Error.[%s]",f[0].TypeName)
	}

	if len(f) != 2 {
		t.Fatalf("TestStruct Field length Error.[%d]",len(f))
	}

	str = getStringStruct(t,StringTagStruct)
	f,err = createFields(str)

	if len(f) != 8 {
		t.Fatalf("TagStruct Field length Error.[%d]",len(f))
	}

}

func TestNewFields(t *testing.T) {

	//フィールド情報のテスト

	//名称
	//タイプ番号
	//型名

	//タグの設定ができるかをチェック（中身はsetTag）

	//無視するフィールドのときにnil

	//ds.Metaを渡した時にtrue

	t.Fatal("Not Implements.")
}

func TestSetTag(t *testing.T) {

	var err error
	var f *Field

	f = &Field{}
	err = setTag(f,"test")
	if err == nil {
		t.Errorf("csv line error")
	}
	t.Log("LOG:" + err.Error())

	err = setTag(f,"test=test,csverror")
	if err == nil {
		t.Errorf("csv line error")
	}
	t.Log("LOG:" + err.Error())

	err = setTag(f,"test=2")
	if err != nil {
		t.Errorf("mismatch flag is error nil")
	}

	//test default value
	defaultSetTagValue(f,t)

	err = setTag(f,"name=test")
	if err != nil {
		t.Errorf("mismatch flag is error nil")
	}
	if f.DisplayName != "test" {
		t.Errorf("name value update error[%s]",f.DisplayName)
	}
	err = setTag(f,"args=aaa,name=test5")
	if f.DisplayName != "test5" {
		t.Errorf("name value update error[%s]",f.DisplayName)
	}


	err = setTag(f,"type=list")
	if f.EditType != "list" {
		t.Errorf("type value update error[%s]",f.EditType)
	}
	err = setTag(f,"args=aaa,type=list2")
	if f.EditType != "list2" {
		t.Errorf("type value update error[%s]",f.EditType)
	}


	err = setTag(f,"default=aaa")
	if f.Default != "aaa" {
		t.Errorf("default value update error[%s]",f.Default)
	}

	err = setTag(f,"key=value,default=bbb")
	if f.Default != "bbb" {
		t.Errorf("default value update error[%s]",f.Default)
	}

	err = setTag(f,"required=aaaa")
	if f.Required != false {
		t.Errorf("required value unknown not update error[%v]",f.Required)
	}

	err = setTag(f,"required=true")
	if f.Required != true {
		t.Errorf("required value update error[%v]",f.Required)
	}

	err = setTag(f,"csv=111,required=false")
	if f.Required != false {
		t.Errorf("required value csv update error[%v]",f.Required)
	}

	err = setTag(f,"csv=111,required=aaaa")
	if f.Required != false {
		t.Errorf("required value unknown not update error[%v]",f.Required)
	}

	err = setTag(f,"editable=aaaa")
	if f.Editable != false {
		t.Errorf("editable value unknown not update error[%v]",f.Editable)
	}

	err = setTag(f,"editable=true")
	if f.Editable != true {
		t.Errorf("editable value update error[%v]",f.Editable)
	}

	err = setTag(f,"csv=111,editable=true")
	if f.Editable != true {
		t.Errorf("editable value csv update error[%v]",f.Editable)
	}

	err = setTag(f,"csv=111,editable=aaaa")
	if f.Editable != true {
		t.Errorf("editable value unknown not update error[%v]",f.Editable)
	}

	err = setTag(f,"display=aaaa")
	if f.Display != false {
		t.Errorf("display value unknown not update error[%v]",f.Display)
	}

	err = setTag(f,"display=true")
	if f.Display != true {
		t.Errorf("display value update error[%v]",f.Display)
	}

	err = setTag(f,"csv=111,display=true")
	if f.Display != true {
		t.Errorf("display value csv update error[%v]",f.Display)
	}

	err = setTag(f,"csv=111,display=aaaa")
	if f.Display != true {
		t.Errorf("display value unknown not update error[%v]",f.Display)
	}


}

func defaultSetTagValue(f *Field,t *testing.T) {
	if f.DisplayName != "" || f.Default != "" ||
		f.Required != false || f.Editable != false ||
		f.Display != false || f.EditType != "" {
		t.Errorf("no parameter default error")
	}
}

func getStringStruct(t *testing.T,ss string) *ast.StructType {
	fset := token.NewFileSet()
	//構文木を取得
	f, err := parser.ParseFile(fset, "", ss, 0)
	if err != nil {
		t.Errorf("initialize test struct")
		return nil
	}

	decls := f.Decls

	for _,decl := range decls {
		if gen,ok := decl.(*ast.GenDecl); !ok {
			continue
		} else if gen.Specs == nil || len(gen.Specs) <= 0 {
			continue
		} else if typ,tOK := gen.Specs[0].(*ast.TypeSpec); !tOK {
			continue
		} else if rtn,sOK := typ.Type.(*ast.StructType); sOK {
			return rtn
		}
	}
	return nil
}

const StringTestStruct =`
package test
import(
)

type TestStruct struct {
	TestField1 int
	TestField2 string
}

`
const StringTagStruct =`
package test
import(
)

type TagStruct struct {
	TestField1 int     ` + "`" + `datastore:"" json:"" ` + "`" + `
	TestField2 string  ` + "`" + `datastore:"" json:"" dizzy:""` + "`" +  `
	TestField3 int     ` + "`" + `datastore:"" json:"" ` + "`" + `
	TestField4 string  ` + "`" + `datastore:"" json:"" dizzy:""` + "`" + `
	TestField5 int     ` + "`" + `datastore:"" json:"" ` + "`" + `
	TestField6 string  ` + "`" + `datastore:"" json:"" dizzy:""` + "`" + `
	TestField7 int     ` + "`" + `datastore:"" json:"" ` + "`" + `
	TestField8 string  ` + "`" + `datastore:"" json:"" dizzy:""` + "`" + `
}
`