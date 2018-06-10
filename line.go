package main

import (
	"strings"
	"fmt"
)

// name=value 形式のcsv変換
type Liner map[string]string

func NewLiner (l string) (Liner,error) {

	lmap := make(map[string]string)
	if l == "" {
		return lmap,nil
	}

	csv := strings.Split(l,",")
	for _,elm := range csv {
		args := strings.Split(elm,"=")
		if len(args) != 2 {
			return nil,fmt.Errorf("format Error[%s]",elm)
		}
		name := args[0]
		val := args[1]

		_,ok := lmap[name]
		if ok {
			return nil,fmt.Errorf("ducaplicate key exists [%s][%s]",name,elm)
		}
		lmap[name] = val
	}

	return lmap,nil
}
