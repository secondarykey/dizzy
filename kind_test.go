package main

import (
	"testing"
	"fmt"
)

func TestCreateKind(t *testing.T) {

	var kinds []*Kind
	var err error
	var p string

	kinds, p, err = createKinds("./examples/test/notexists.go")
	if err == nil {
		t.Fatalf("Error createKinds() : file not exists ")
	}
	t.Logf("LOG:error not exist message:[%v]", err)

	kinds, p, err = createKinds("./examples/test/prefix.go")
	if err != nil {
		t.Errorf("Error createKinds() : no comment is not error")
	}
	if kinds == nil {
		t.Errorf("Error createKinds() : kinds not nil")
	}

	if len(kinds) != 0 {
		t.Errorf("Error createKinds() : not error")
	}

	if p != "test" {
		t.Errorf("Error createKinds() : package name[%s]", p)
	}

	kinds, p, err = createKinds("./examples/test/nopackage.go.tmp")
	if err == nil {
		t.Fatalf("Error createKinds() : no package name")
	}
	t.Logf("LOG:error nopackage message:[%v]", err)

	//GenDecl以外

	//type以外

	//struct以外

	//正常系

	//２つの場合の読み込み

}


func TestNewKind(t *testing.T) {

	var kind *Kind
	var err error
	t.Logf("LOG:Default Value")
	//全未指定(デフォルト)の場合
	kind,err = NewKind("+DIZZY", nil)
	defaultKind(t, kind, "+DIZZY")
	kind,err = NewKind("+DIZZY()", nil)
	defaultKind(t, kind, "+DIZZY()")

	kind,err = NewKind("+DIZZY(", nil)
	if err == nil {
		t.Errorf("parse error")
	}
	kind,err = NewKind("+DIZZY)", nil)
	if err == nil {
		t.Errorf("parse error")
	}
	//other value
	kind,err = NewKind("+DIZZY(test=test,xxx=true)", nil)
	defaultKind(t, kind, "+DIZZY(test=test,xxx=true)")

	t.Logf("LOG:KindName(name) Value")
	//nameテスト
	kind,err = NewKind("+DIZZY(name=SampleName)", nil)
	if kind.KindName != "SampleName" {
		t.Errorf("Error:name KindName[%s]", kind.KindName)
	}
	//csv
	kind,err = NewKind("+DIZZY(csv=head,name=SampleName)", nil)
	if kind.KindName != "SampleName" {
		t.Errorf("Error:name csv KindName[%s]", kind.KindName)
	}
	//nothing
	kind,err = NewKind("+DIZZY(csv=head)", nil)
	if kind.KindName != "TestName" {
		t.Errorf("Error:name csv KindName[%s]", kind.KindName)
	}

	//cacheテスト
	t.Logf("LOG:Cache(cache) Value")
	//false
	kind,err = NewKind("+DIZZY(cache=false)", nil)
	if kind.Cache != false {
		t.Errorf("Error:cache false Cache[%v]", kind.Cache)
	}
	//true
	kind,err = NewKind("+DIZZY(cache=true)", nil)
	if kind.Cache != true {
		t.Errorf("Error:cache false Cache[%v]", kind.Cache)
	}
	//other
	kind,err = NewKind("+DIZZY(cache=other)", nil)
	if kind.Cache != true {
		t.Errorf("Error:cache other Cache[%v]", kind.Cache)
	}
	//csv
	kind ,err= NewKind("+DIZZY(csv=head,cache=false)", nil)
	if kind.Cache != false {
		t.Errorf("Error:cache csv Cache[%v]", kind.Cache)
	}
	//nothing
	kind,err = NewKind("+DIZZY(csv=head)", nil)
	if kind.Cache != true {
		t.Errorf("Error:cache nothing Cache[%v]", kind.Cache)
	}

	t.Logf("LOG:Content(content) Value")
	//contentテスト
	kind,err = NewKind("+DIZZY(content=text)", nil)
	if kind.Content != ContentText {
		t.Errorf("Error:content text Content[%d]", kind.Content)
	}
	kind,err = NewKind("+DIZZY(content=binary)", nil)
	if kind.Content != ContentBinary {
		t.Errorf("Error:content binary Content[%d]", kind.Content)
	}
	kind,err = NewKind("+DIZZY(content=other)", nil)
	if kind.Content != ContentNone {
		t.Errorf("Error:content other Content[%d]", kind.Content)
	}
	kind,err = NewKind("+DIZZY(csv=head,content=binary)", nil)
	if kind.Content != ContentBinary {
		t.Errorf("Error:content csv Content[%d]", kind.Content)
	}
	kind,err = NewKind("+DIZZY(csv=head)", nil)
	if kind.Content != ContentNone {
		t.Errorf("Error:content nothing Content[%d]", kind.Content)
	}

	t.Logf("LOG:Limit(limit) Value")
	//limitテスト
	kind,err = NewKind("+DIZZY(limit=20)", nil)
	if kind.Limit != 20 {
		t.Errorf("Error:limit 20 Limit[%d]", kind.Limit)
	}
	kind,err = NewKind("+DIZZY(limit=20.0)", nil)
	if kind.Limit != 10 {
		t.Errorf("Error:limit 20.0 Limit[%d]", kind.Limit)
	}
	kind,err = NewKind("+DIZZY(csv=head,limit=30)", nil)
	if kind.Limit != 30 {
		t.Errorf("Error:limit csv Limit[%d]", kind.Limit)
	}
	kind,err = NewKind("+DIZZY(csv=head)", nil)
	if kind.Limit != 10 {
		t.Errorf("Error:limit nothing Limit[%d]", kind.Limit)
	}

	//type spec pattern

	fmt.Println(err)
	//parser Test

}

func defaultKind(t *testing.T, kind *Kind, prefix string) {
	if kind.TypeName != "TestName" {
		t.Errorf("Error:%s default TypeName[%s]", prefix, kind.TypeName)
	}
	if kind.KindName != "TestName" {
		t.Errorf("Error:%s default KindName[%s]", prefix, kind.KindName)
	}
	if kind.Cache != true {
		t.Errorf("Error:%s default Cache[%v]", prefix, kind.Cache)
	}
	if kind.Content != ContentNone {
		t.Errorf("Error:%s default Content[%d]", prefix, kind.Content)
	}
	//kind.Struct != nil ()
	if kind.Limit != 10 {
		t.Errorf("Error:%s default Limit[%d]", prefix, kind.Limit)
	}
}