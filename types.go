package main

type Content int

const (
	ContentNone Content = iota
	ContentBinary
	ContentText
)

func (c Content) Value(val string) Content {

	if val == "binary" {
		return ContentBinary
	} else if val == "text" {
		return ContentText
	}

	return ContentNone
}

type FieldType int

const (
	IntField   = 10 //input integer
	FloatField = 20 //input integer
	BoolField  = 30
	StringField = 50 //input text

	//LongStringField //textarea
	//appengine.BlobKey
	//appengine.GeoPoint
)
