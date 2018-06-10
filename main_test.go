package main

import "testing"

func TestGenerate(t *testing.T) {

	err := gen("examples")
	if err != nil {
		t.Fatalf("gen() error[%s]",err)
	}

	//出力ファイルを確認




}
