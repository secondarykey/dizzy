package main

import (
	"testing"
)

func TestLiner(t *testing.T) {

	var l Liner
	var err error
	l,err  = NewLiner("test1=1,test2=2")
	if err != nil {
		t.Fatal("NewLiner Error")
	}
	if v,ok := l["test1"] ; !ok {
		t.Fatal("NewLiner exists Error(test1)")
	}else {
		if v != "1" {
			t.Fatalf("NewLiner value Error[%s]",v)
		}
	}

	if v,ok := l["test2"] ; !ok {
		t.Fatal("NewLiner exists Error(test2)")
	}else {
		if v != "2" {
			t.Fatalf("NewLiner value Error(test2)[%s]",v)
		}
	}

	if _,ok := l["test3"] ; ok {
		t.Fatal("NewLiner not found Error(test3)")
	}

	idx := 0
	for key := range l {
		switch key {
		case "test1","test2":
		default :
			t.Fatalf("NewLiner analysis Error[%s]",key)
		}
		idx++
	}
	if idx != 2 {
		t.Fatalf("NewLiner count Error[2 != %d]",idx)
	}

	//空文字アクセス
	l,err = NewLiner("")
	if err != nil {
		t.Fatalf("NewLiner empty string Error[%s]",err)
	}

	if l == nil {
		t.Fatalf("NewLiner empty string nil Error[%s]",l)
	}

	if v,ok := l["nothing"] ; ok {
		t.Fatalf("NewLiner empty string exist Error[%s]",v)
	}

	//形式エラー
	l,err  = NewLiner("test")
	if err == nil {
		t.Fatalf("NewLiner empty format Error[%s]", l)
	}
	l,err  = NewLiner("prefix=4,test")
	if err == nil {
		t.Fatalf("NewLiner empty format Error[%s]", l)
	}

	l,err  = NewLiner("prefix=4,prefix=3")
	if err == nil {
		t.Fatalf("NewLiner ducaplicate Error[%s]", l)
	}


}
