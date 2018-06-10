package examples

import "github.com/knightso/base/gae/ds"

//+DIZZY(name=KindType,cache=false,content=none,limit=10,url=work)
type WorkType struct {
	Name string `dizzy:"required=true"`
	Number  int
	Type    int `datastore:",noindex" json:"version" dizzy:"type=select"`
	ds.Meta
}

//+DIZZY(name=SampleType,cache=true,content=none,limit=10,url=sample)
type SampleType struct {
	Name string `datastore:",noindex" json:"version" dizzy:"display=true"`
	ds.Meta
}