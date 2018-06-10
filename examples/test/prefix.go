package test

import "fmt"

//+OTHERCOMMENT
type NoComment struct {
}

//+OTHERCOMMENT
func FunctionComment() {
	fmt.Println("Hello datastore generator")
}

func CommentNone() {
	fmt.Println("Hello datastore generator")
}
