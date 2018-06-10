package main

import (
	"fmt"
	"html/template"
)

func GenerateInputHTML(f *Field) template.HTML {

	switch f.Type {
	case IntField:
		html := fmt.Sprintf(inputIntTag, f.Name, f.Name, f.Name, f.Required)
		html += fmt.Sprintf(labelTag, f.Name, f.DisplayName)
		html += fmt.Sprintf(errorTag)
		return template.HTML(html)
	case FloatField:

		html := fmt.Sprintf(inputFloatTag, f.Name, f.Name, f.Name, f.Required)
		html += fmt.Sprintf(labelTag, f.Name, f.DisplayName)
		html += fmt.Sprintf(errorTag)
		return template.HTML(html)

	case BoolField:
		html := fmt.Sprintf(inputStringTag, f.Name, f.Name, f.Name, f.Required)
		html += fmt.Sprintf(labelTag, f.Name, f.Name)
		return template.HTML(html)

	case StringField:
		html := fmt.Sprintf(inputStringTag, f.Name, f.Name, f.Name, f.Required)
		html += fmt.Sprintf(labelTag, f.Name, f.Name)
		return template.HTML(html)

	}

	return "<h7>Not Support</h7>"
}

var inputIntTag = `<input class="mdl-textfield__input" type="text" name="%s" id="%s" pattern="-?[0-9]*?" value="{{ .Kind.%s }}" {{if eq %t true}}required{{end}}/>`
var inputFloatTag = `<input class="mdl-textfield__input" type="text" name="%s" id="%s" pattern="-?[0-9]*(\.[0-9]+)?" value="{{ .Kind.%s }}" {{if eq %t true}}required{{end}}/>`
var inputStringTag = `<input class="mdl-textfield__input" type="text" name="%s" id="%s" value="{{ .Kind.%s }}" {{if eq %t true}}required{{end}}/>`
var labelTag = `<label class="mdl-textfield__label" for="%s">%s...</label>`
var errorTag = `<span class="mdl-textfield__error">Input is not a number!</span>`
