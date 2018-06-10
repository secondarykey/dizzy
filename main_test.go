package main

import "testing"

func TestGenerate(t *testing.T) {

	example := "examples/"
	templates := example + "templates/"

	err := gen(example)
	if err != nil {
		t.Fatalf("gen() error[%s]", err)
	}

	if !exists(example + "sample_access.go") {
		t.Fatalf("not generate access file ")
	}
	if !exists(example + "sample_handler.go") {
		t.Fatalf("not generate handler file ")
	}
	if !exists(example + "dizzy.go") {
		t.Fatalf("not generate dizzy file ")
	}
	if !exists(example + "dizzy_app.yaml") {
		t.Fatalf("not generate appengine yaml")
	}
	if !exists(example + "index.yaml") {
		t.Fatalf("not generate index yaml")
	}
	if !exists(templates + "edit_KindType.tmpl") {
		t.Fatalf("not generate KindType edit template")
	}
	if !exists(templates + "edit_SampleType.tmpl") {
		t.Fatalf("not generate SampleType edit template")
	}
	if !exists(templates + "view_KindType.tmpl") {
		t.Fatalf("not generate KindType view template")
	}
	if !exists(templates + "view_SampleType.tmpl") {
		t.Fatalf("not generate SampleType view template")
	}
	if !exists(templates + "error.tmpl") {
		t.Fatalf("not generate error template")
	}
	if !exists(templates + "layout.tmpl") {
		t.Fatalf("not generate layout template")
	}
	if !exists(templates + "top.tmpl") {
		t.Fatalf("not generate top template")
	}

}
