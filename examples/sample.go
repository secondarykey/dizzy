package examples

import "github.com/knightso/base/gae/ds"

//テストでも使用しています。

//+DIZZY(name=KindType,cache=false,content=none,limit=10,url=work)
type WorkType struct {
	Name   string `dizzy:"required=true"`
	Number int
	Type   int `datastore:",noindex" json:"version" dizzy:"type=select"`
	ds.Meta
}

//+DIZZY(name=SampleType,cache=true,content=none,limit=10,url=sample)
type SampleType struct {
	Name string `datastore:",noindex" json:"version" dizzy:"display=true"`
	ds.Meta
}

//+DIZZY(name=Test,cache=true,content=none,limit=10,url=test)
type TestKind struct {
	String string `dizzy:"display=true"`

	Int   int   `dizzy:"display=true"`
	Int8  int8  `dizzy:"display=true"`
	Int16 int16 `dizzy:"display=true"`
	Int32 int32 `dizzy:"display=true"`
	Int64 int64 `dizzy:"display=true"`

	Float32 float32 `dizzy:"display=true"`
	Float64 float64 `dizzy:"display=true"`

	Bool bool `dizzy:"display=true"`

	ds.Meta
}
