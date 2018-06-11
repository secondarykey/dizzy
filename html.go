package main

import (
	"fmt"
	"html/template"
)

func GenerateInputHTML(f *Field) template.HTML {

	html := ""
	switch f.Type {
	case IntField:

		html += fmt.Sprintf(divTag, f.Name, "")
		html += fmt.Sprintf(inputIntTag, f.Name, f.Name, f.Name, f.Required)
		html += fmt.Sprintf(labelTag, f.Name, f.DisplayName)
		html += fmt.Sprintf(errorTag)
		html += fmt.Sprintf("</div>")

	case FloatField:

		html += fmt.Sprintf(divTag, f.Name, "")
		html += fmt.Sprintf(inputFloatTag, f.Name, f.Name, f.Name, f.Required)
		html += fmt.Sprintf(labelTag, f.Name, f.DisplayName)
		html += fmt.Sprintf(errorTag)

	case BoolField:

		html += fmt.Sprintf(divTag, f.Name, "mdl-select")
		html += fmt.Sprintf(inputBoolTag, f.Name, f.Name, f.Name, f.Name)

		html += fmt.Sprintf(labelBoolTag, f.Name, f.Name)

		html += fmt.Sprintf(boolSelectTag, f.Name)
		html += fmt.Sprintf("</div>")

	case StringField:
		html += fmt.Sprintf(divTag, f.Name, "")
		html += fmt.Sprintf(inputStringTag, f.Name, f.Name, f.Name, f.Required)
		html += fmt.Sprintf(labelTag, f.Name, f.Name)
		html += fmt.Sprintf("</div>")

	default:
		html += fmt.Sprintf(divTag, f.Name, "")
		html += "<h7>Not Support</h7>"
	}

	html += fmt.Sprintf("</div>")
	return template.HTML(html)
}

var divTag = `<div id="div%s" class="mdl-textfield mdl-js-textfield mdl-textfield--floating-label %s">`

var inputIntTag = `<input class="mdl-textfield__input" type="text" name="%s" id="%s" pattern="-?[0-9]*?" value="{{ .Kind.%s }}" {{if eq %t true}}required{{end}}/>`
var inputFloatTag = `<input class="mdl-textfield__input" type="text" name="%s" id="%s" pattern="-?[0-9]*(\.[0-9]+)?" value="{{ .Kind.%s }}" {{if eq %t true}}required{{end}}/>`
var inputStringTag = `<input class="mdl-textfield__input" type="text" name="%s" id="%s" value="{{ .Kind.%s }}" {{if eq %t true}}required{{end}}/>`

var inputBoolTag = `
<input type="hidden" name="%s" id="%s" value="{{ .Kind.%s }}" />
<input class="mdl-textfield__input" type="text" id="display%s" readonly />
`

var boolSelectTag = `
<ul class="mdl-menu mdl-menu--bottom-left mdl-js-menu mdl-js-ripple-effect" for="display%s">
<li class="mdl-menu__item" data-id="true">True</li>
<li class="mdl-menu__item" data-id="false">False</li>
</ul>`

var labelTag = `<label class="mdl-textfield__label" for="%s">%s...</label>`
var labelBoolTag = `<label class="mdl-textfield__label" for="display%s">%s...</label>`
var errorTag = `<span class="mdl-textfield__error">Input is not a number!</span>`
