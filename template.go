package main

//dizzy release generated:2018-06-11 21:44:39.973654697 +0900 JST m=+0.000546971

import (
    "fmt"
    "time"
    "text/template"
)

const TemplateDirectory = "templates"

//developer mode
const UseTemplateFile = true

const AppTemplateFile = "app.tmpl"
const AppTemplate = `runtime: go
api_version: go1

# dizzy auto generated:{{generated}}

handlers:
  - url: /_dizzy/.*
    script: _go_app
    login: admin
    secure: always
  - url: /.*
    script: _go_app`

const DizzyGoTemplateFile =  "dizzy.tmpl"
const DizzyGoTemplate = `package {{.PackageName}}

//dizzy auto generated:{{generated}}

import (
	"log"
    "html/template"
    "net/http"
    "time"

	"github.com/gorilla/mux"
)

const TEMPLATE_DIR = "templates/"

func init() {

    r := mux.NewRouter()

    r.HandleFunc("/_dizzy/",view)

{{range .Kinds}}
    h{{.TypeName}} := {{.TypeName}}Handler{}
    r.HandleFunc("/_dizzy/{{.URL}}/",h{{.TypeName}}.view).Methods("GET")
    r.HandleFunc("/_dizzy/{{.URL}}/create",h{{.TypeName}}.create).Methods("GET")
    r.HandleFunc("/_dizzy/{{.URL}}/edit",h{{.TypeName}}.edit).Methods("POST")
    r.HandleFunc("/_dizzy/{{.URL}}/update",h{{.TypeName}}.update).Methods("POST")
    r.HandleFunc("/_dizzy/{{.URL}}/delete",h{{.TypeName}}.delete).Methods("POST")
{{end}}

	http.Handle("/",r)
}

func view(w http.ResponseWriter,r *http.Request)  {
    dto := struct{}{}
    parse(w,"top.tmpl",dto)
}

func convertDate(t *time.Time) string {
	if t.IsZero() {
		return "None"
	}
	jst, _ := time.LoadLocation("Asia/Tokyo")
	jt := t.In(jst)
	return jt.Format("2006/01/02 15:04:05")
}

func errorPage(w http.ResponseWriter,status int,title string,name... string) {

    w.WriteHeader(status)
    dto := struct {
        Title string
        Description []string
        Number int
    } {title,name,status}

    parse(w,"error.tmpl",dto)
    return
}


func parse(w http.ResponseWriter,tmplName string,dto interface{}) {

    funcMap := template.FuncMap{
        "convertDate":convertDate,
    }
    tmpl,err := template.New("layout").Funcs(funcMap).ParseFiles(TEMPLATE_DIR+"layout.tmpl",TEMPLATE_DIR+tmplName)
    if err != nil {
        w.WriteHeader(500)
        log.Println("Template Parse Error",err.Error())
        return
    }

    err = tmpl.Execute(w,dto)
    if err != nil {
        w.WriteHeader(500)
        log.Println("Template Execute Error",err.Error())
        return
    }
    w.WriteHeader(200)
    return
}
`

const EditTemplateFile =  "edit.tmpl"
const EditTemplate = `{{ "{{define \"title\"}}" }}
{{ .Kind.KindName }} Edit[dizzy auto generated:{{generated}}
{{ "{{end}}" }}

{{ "{{define \"content\"}}" }}
<form method="post" action="/_dizzy/{{.Kind.URL}}/update" onsubmit="return confirm('realy?');">

<input type="hidden" name="Version" value="{{ "{{ .Kind.Version }}" }}">

{{ "{{ if not .Kind.Key }}" }}
<input type="hidden" name="keyId" value="">
{{ "{{ else }}" }}
<input type="hidden" name="keyId" value="{{"{{.Kind.Key.StringID}}"}}">
{{ "{{ end }}" }}

  <table class="mdl-data-table mdl-js-data-table mdl-shadow--2dp">
    <tr>
        <th class="mdl-data-table__cell--non-numeric">
{{ "{{ if not .Kind.Key }}" }}
        New {{ .Kind.KindName }}
{{ "{{ else }}" }}
        Key : {{ "{{ .Kind.Key.StringID }}" }}
{{ "{{ end }}" }}
        </th>
    </tr>
    <tbody>

{{ range .Kind.Fields }}

      <tr>
        <td style="text-align:center;">

{{ generateInputHTML . }}

        </td>
      </tr>

{{ end }}
      <tr>
        <td>
          <button type="submit" class="mdl-button mdl-js-button mdl-button--raised mdl-button--icon mdl-button--primary">
            <i class="material-icons">save</i>
          </button>
        </td>
      </tr>
    </tbody>
  </table>
</form>


<script>
var list = ""
var valTag = ""
var displayTag = ""

{{ range .Kind.Fields }}

{{ if eq .Editable true }}
{{ if eq .TypeName "bool" }}

list = document.querySelectorAll('#div{{ .Name }}.mdl-select > ul > li');
valTag = document.querySelector('#{{ .Name }}');
displayTag = document.querySelector('#display{{ .Name }}');


for (var i=0; i< list.length; i++) {

  var li = list[i];
  li.addEventListener('click', function(e) {
    valTag.setAttribute('value', e.target.getAttribute("data-id"));
    displayTag.setAttribute('value', e.target.textContent);
  });

  //default value
  if ( li.getAttribute("data-id") == "{{"{{"}} .Kind.{{.Name}} {{"}}"}}" ) {
    valTag.setAttribute('value', li.getAttribute("data-id"));
    displayTag.setAttribute('value', li.textContent);
  }

}

{{ end }}

{{ end }}

{{ end }}

</script>

{{ "{{ end }}" }}

`

const ErrorTemplateFile =  "error.tmpl"
const ErrorTemplate = `{{ "{{define \"title\"}}" }}

{{ "{{ .Title }}[{{ .Number }}]" }}

{{ "{{end}}" }}

{{ "{{define \"content\"}}" }}
<ul>
{{ "{{ range .Description }}" }}
    <li>{{ "{{.}}" }}</li>
{{ "{{ end }}" }}
</ul>

<h7> dizzy auto generated:{{generated}} </h7>
{{ "{{ end }}" }}
`

const AccessTemplateFile = "access.tmpl"
const AccessTemplate = `package {{.PackageName}}

//dizzy auto generated:{{generated}}

import (
    "fmt"
    "net/http"
    "strconv"

	"github.com/satori/go.uuid"
	"github.com/knightso/base/gae/ds"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

func init() {
{{ range .Kinds }}
    ds.CacheKinds["{{ .KindName }}"] = {{ .Cache }}
{{ end }}
}

{{ range .Kinds }}

const KIND_{{ .KindName }} = "{{ .KindName }}"

func (s *{{.TypeName}}) GenerateKey(r *http.Request) error {
	uid, err := uuid.NewV4()
	if err != nil {
	    return err
	}
	id := uid.String()
	key := s.CreateKey(r,id)
	s.SetKey(key)

	return nil
}

func (s *{{.TypeName}}) CreateKey(r *http.Request,id string) *datastore.Key {
	c := appengine.NewContext(r)
   	return datastore.NewKey(c, KIND_{{ .KindName }}, id, 0, nil)
}

func (s *{{.TypeName}}) Select(r *http.Request,p int) ([]{{.TypeName}},error) {
    //UpdatedAt,Deletedによる検索
   	rtn := make([]{{.TypeName}},0)

    c := appengine.NewContext(r)
    cursor := ""

    q := datastore.NewQuery(KIND_{{.KindName}}).Order("- UpdatedAt").Filter("Deleted=",false)
    if  p > 0 {
    	item, err := memcache.Get(c, s.getCursorName(p))
    	if err == nil {
    		cursor = string(item.Value)
    	}
    	q = q.Limit({{.Limit}})
    	if cursor != "" {
    		cur, err := datastore.DecodeCursor(cursor)
    		if err == nil {
    			q  = q.Start(cur)
    		}
    	}
    }

   	t := q.Run(c)
   	for {
   		var d {{.TypeName}}
   		key, err := t.Next(&d)
   		if err == datastore.Done {
   			break
   		}
   	    //TODO どうするか？
       	if _,ok := err.(*datastore.ErrFieldMismatch) ; !ok {
       		if err != nil {
       			return nil,err
       		}
       	}
   		d.SetKey(key)
   		rtn = append(rtn, d)
   	}

   	if p > 0 {
   		cur,err := t.Cursor()
   		if err != nil {
   			return nil,err
   		}

   		err = memcache.Set(c,&memcache.Item{
   			Key:s.getCursorName(p+1),
   			Value:[]byte(cur.String()),
   		})
   		if err != nil {
   			return nil,err
   		}
   	}

   	return rtn, nil
}

func (s *{{.TypeName}}) getCursorName(p int) string {
	return "{{ .KindName }}_cursor_"+strconv.Itoa(p)
}

func (s *{{.TypeName}}) SelectById(r *http.Request,id string) error {
    //データを検索
    key := s.CreateKey(r,id)
	c := appengine.NewContext(r)
	//自身に設定する
    err := ds.Get(c,key,s)
    //TODO どうするか？
	if _,ok := err.(*datastore.ErrFieldMismatch) ; !ok {
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *{{.TypeName}}) SelectByIdWithVersion(r *http.Request,id string,version int) error {
    //データを検索
    key := s.CreateKey(r,id)
	c := appengine.NewContext(r)
	//自身に設定する
    err := ds.GetWithVersion(c,key,version,s)
    //TODO どうするか？
	if _,ok := err.(*datastore.ErrFieldMismatch) ; !ok {
		if err != nil {
			return err
		}
	}
	return nil
}


func (s *{{.TypeName}}) Put(r *http.Request) error {
	c := appengine.NewContext(r)
	return ds.Put(c,s)
}

func (s *{{.TypeName}}) Delete(r *http.Request) error {
	var err error
	c := appengine.NewContext(r)
	err = ds.Delete(c,s.GetKey())
	if err != nil {
		return err
	}
{{ if ne .Content 0 }}
{{ end }}
  return nil
}

func (s *{{.TypeName}}) DeleteLogical(r *http.Request) error {

    var err error
    s.Deleted = true
    err = s.Put(r)
    if err != nil {
        return err
    }

{{ if ne .Content 0}}
{{ end }}

  return nil
}

func (s *{{.TypeName}}) SetDefault() (err error) {


{{ range .Fields }}
    {{ if eq .Display true }}

    {{ end }}
{{ end }}

    return nil
}

func (s *{{.TypeName}}) Validate(r *http.Request) (err error) {

    var val string
    var bit int
    var typeName string

    //declared and not used
    if false {
        fmt.Printf("declared and not used.bit[%d]typeName[%s]",bit,typeName)
    }

{{ range .Fields }}

    val = r.FormValue("{{.Name}}")
    typeName = "{{.TypeName}}"

    {{ if eq .Editable true }}

        //int value
        {{ if eq .Type 10 }}

    if len(typeName) > 3 {
        bit,err = strconv.Atoi(typeName[3:])
        if err != nil {
            return err
        }
    } else {
        bit = 64
    }

    if v,err := strconv.ParseInt(val,10,bit) ; err != nil {
    return fmt.Errorf("[%s] value parse error[%s] [%s]","{{.Name}}",err,val)
    } else {
        s.{{ .Name }} = {{.TypeName}}(v)
    }

        //float value
        {{ else if eq .Type 20 }}

    if len(typeName) > 5 {
        bit,err = strconv.Atoi(typeName[5:])
        if err != nil {
            return err
        }
    } else {
        return fmt.Errorf("float type error[%s]",typeName)
    }

    if v,err := strconv.ParseFloat(val,64) ; err != nil {
        return fmt.Errorf("[%s] value parse error[%s] [%s]","{{.Name}}",err,val)
    } else {
        s.{{ .Name }} = {{.TypeName}}(v)
    }

        //bool value
        {{ else if eq .Type 30 }}

    if v,err := strconv.ParseBool(val) ; err != nil {
        return fmt.Errorf("[%s] value parse error[%s] [%s]","{{.Name}}",err,val)
    } else {
        s.{{ .Name }} = v
    }

        {{ else if eq .Type 50 }}

    s.{{ .Name }} = val

        {{ end }}


    {{ end }}

{{ end }}
    return nil
}

{{ if ne .Content 0 }}

//コンテンツデータが存在する場合
type {{.TypeName}}Data struct {
	key     *datastore.Key
	Content []byte
}

//子のキー作成
func (d *{{ .TypeName }}Data) GetKey() *datastore.Key {
	return d.key
}

func (d *{{ .TypeName }}Data) SetKey(k *datastore.Key) {
	d.key = k
}

//同時にコンテンツを生成する場合
func (s *{{.TypeName}}) Put(content []byte) error {
    //
    s.Put()
    //ページャに対する検索
}

{{ end }}

{{ end }}

`

const HandlerGoTemplateFile =  "handler.tmpl"
const HandlerGoTemplate = `package {{.PackageName}}

//dizzy auto generated:{{generated}}

import (
    "net/http"
    "strconv"
)


{{range .Kinds}}
type {{.TypeName}}Handler struct { }

func (h {{.TypeName}}Handler) view(w http.ResponseWriter,r *http.Request)  {
    dao := {{.TypeName}}{}
    p := 1
  	q := r.URL.Query()
   	pageBuf := q.Get("page")
   	if pageBuf != "" {
   		page,err := strconv.Atoi(pageBuf)
   		if err == nil {
   			p = page
   		}
   	}
	kinds,err := dao.Select(r,p)
	if err != nil {
		errorPage(w,500,"Datastore select error",err.Error())
		return
	}

    max := false
    if kinds == nil || len(kinds) == 0 {
        max = true
    }

	dto := struct{
		List []{{.TypeName}}
		Prev int
		Next int
		Max bool
	} {kinds,p-1,p+1,max}
	parse(w,"view_{{.KindName}}.tmpl",dto)
}

func (h {{.TypeName}}Handler) create(w http.ResponseWriter,r *http.Request)  {
	//empty
	obj := &{{.TypeName}} {}

	dto := struct{
		Kind *{{.TypeName}}
	} {obj}

    //デフォルト設定
    obj.SetDefault()

	parse(w,"edit_{{.KindName}}.tmpl",dto)
}

func (h {{.TypeName}}Handler) edit(w http.ResponseWriter,r *http.Request)  {

	obj := &{{.TypeName}} {}
    keyId := r.FormValue("keyId")

    err := obj.SelectById(r,keyId)
    if err != nil {
        errorPage(w,500,"Select Error",keyId)
        return
	}

	dto := struct{
		Kind *{{.TypeName}}
	} {obj}

	parse(w,"edit_{{.KindName}}.tmpl",dto)
}

func (h {{.TypeName}}Handler) update(w http.ResponseWriter,r *http.Request)  {

	obj := &{{.TypeName}} {}
    keyId := r.FormValue("keyId")

    var err error

    if keyId == "" {
        err = obj.GenerateKey(r)
	    if err != nil {
            errorPage(w,500,"Datastore key set error",err.Error())
            return
	    }
    } else {
        verBuf := r.FormValue("Version")
        ver,err := strconv.Atoi(verBuf)
	    if err != nil {
            errorPage(w,500,"Version parse Error",err.Error())
            return
	    }

        err = obj.SelectByIdWithVersion(r,keyId,ver)
	    if err != nil {
            errorPage(w,500,"Datastore select Error",err.Error())
            return
	    }
    }

    //入力値の検証と設定
	if err = obj.Validate(r); err != nil {
        errorPage(w,400,"Validate error",err.Error())
        return
	}

    err = obj.Put(r)
    if err != nil {
        errorPage(w,500,"Datastore put error",err.Error())
        return
    }

    //項目の設定
	dto := struct{
		Kind *{{.TypeName}}
	} {obj}
	parse(w,"edit_{{.KindName}}.tmpl",dto)
}

func (h {{.TypeName}}Handler) delete(w http.ResponseWriter,r *http.Request)  {
    keyId := r.FormValue("keyId")
	obj := &{{.TypeName}} {}

    verBuf := r.FormValue("Version")
    ver,err := strconv.Atoi(verBuf)
    if err != nil {
       errorPage(w,500,"Version parse Error",err.Error())
       return
    }

    err = obj.SelectByIdWithVersion(r,keyId,ver)
    if err != nil {
        errorPage(w,500,"Datastore select error",err.Error())
        return
    }

    //物理削除の場合
	//err = obj.Delete(r)
	err = obj.DeleteLogical(r)
    if err != nil {
        errorPage(w,500,"Datastore delete error",err.Error())
        return
    }

	h.view(w,r)
}
{{end}}




`

const IndexTemplateFile = "index.tmpl"
const IndexTemplate = `indexes:

# dizzy auto generated:{{generated}}

{{ range .Kinds }}
- kind: {{ .KindName }}
  properties:
  - name: Deleted
  - name: UpdatedAt
    direction: desc
{{ end }}

`

const LayoutTemplateFile =  "layout.tmpl"
const LayoutTemplate = `{{ "{{define \"layout\"}}" }}
<!doctype html>
<html lang="en">
  <head>
    <!-- dizzy auto generated:{{generated}} -->
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="description" content="A front-end template that helps you build fast, modern mobile web apps.">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, minimum-scale=1.0">
    <title>Dizzy</title>

    <!-- Add to homescreen for Chrome on Android -->
    <meta name="mobile-web-app-capable" content="yes">
    <link rel="icon" sizes="192x192" href="/manage/images/android-desktop.png">

    <!-- Add to homescreen for Safari on iOS -->
    <meta name="apple-mobile-web-app-capable" content="yes">
    <meta name="apple-mobile-web-app-status-bar-style" content="black">
    <meta name="apple-mobile-web-app-title" content="Material Design Lite">
    <link rel="apple-touch-icon-precomposed" href="/manage/images/ios-desktop.png">

    <!-- Tile icon for Win8 (144x144 + tile color) -->
    <meta name="msapplication-TileImage" content="/manage/images/touch/ms-touch-icon-144x144-precomposed.png">
    <meta name="msapplication-TileColor" content="#3372DF">

    <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto:regular,bold,italic,thin,light,bolditalic,black,medium&amp;lang=en">
    <link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons">
    <link rel="stylesheet" href="https://code.getmdl.io/1.3.0/material.indigo-red.min.css">

    <style>
    table {
        min-width: 400px;
        table-layout: fixed;
    }
    table td {
        word-break: break-all;
        overflow-wrap : break-word;
    }
    .layout .mdl-layout__header .mdl-layout__drawer-button {
      color: rgba(0, 0, 0, 0.54);
    }

    .drawer {
      border: none;
    }
    /* iOS Safari specific workaround */
    .drawer .mdl-menu__container {
      z-index: -1;
    }
    .drawer .navigation {
      z-index: -2;
    }
    /* END iOS Safari specific workaround */
    .drawer .mdl-menu .mdl-menu__item {
      display: -webkit-flex;
      display: -ms-flexbox;
      display: flex;
      -webkit-align-items: center;
          -ms-flex-align: center;
              align-items: center;
    }

    .drawer-header {
      box-sizing: border-box;
      display: -webkit-flex;
      display: -ms-flexbox;
      display: flex;
      -webkit-flex-direction: column;
          -ms-flex-direction: column;
              flex-direction: column;
      -webkit-justify-content: flex-end;
          -ms-flex-pack: end;
              justify-content: flex-end;
      padding: 16px;
      height: 51px;
    }

    .navigation {
      -webkit-flex-grow: 1;
          -ms-flex-positive: 1;
              flex-grow: 1;
    }

    .layout .navigation .mdl-navigation__link {
      display: -webkit-flex !important;
      display: -ms-flexbox !important;
      display: flex !important;
      -webkit-flex-direction: row;
          -ms-flex-direction: row;
              flex-direction: row;
      -webkit-align-items: center;
          -ms-flex-align: center;
              align-items: center;
      color: rgba(255, 255, 255, 0.56);
      font-weight: 500;
    }
    .layout .navigation .mdl-navigation__link:hover {
      background-color: #00BCD4;
      color: #37474F;
    }
    .navigation .mdl-navigation__link .material-icons {
      font-size: 24px;
      color: rgba(255, 255, 255, 0.56);
      margin-right: 32px;
    }

    #add-content {
      position: fixed;
      display: block;
      right: 0;
      bottom: 0;
      margin-right: 40px;
      margin-bottom: 40px;
      z-index: 900;
    }
    </style>

  </head>

  <body>

    <div class="layout mdl-layout mdl-js-layout mdl-layout--fixed-drawer mdl-layout--fixed-header">

      <header class="header mdl-layout__header mdl-color--grey-100 mdl-color-text--grey-600">
        <div class="mdl-layout__header-row">
          <span class="mdl-layout-title"></span>
          <div class="mdl-layout-spacer">
{{"{{template \"title\" .}}"}}
          </div>
          <div class="mdl-textfield mdl-js-textfield mdl-textfield--expandable">
          </div>
        </div>
      </header>

      <div class="drawer mdl-layout__drawer mdl-color--blue-grey-900 mdl-color-text--blue-grey-50">

        <header class="drawer-header">
            <span> Dizzy </span>
        </header>

        <nav class="navigation mdl-navigation mdl-color--blue-grey-800">

          <a class="mdl-navigation__link" href="/_dizzy/"><i class="mdl-color-text--blue-grey-400 material-icons" role="presentation">home</i>Home</a>
{{ range .Kinds }}
          <a class="mdl-navigation__link" href="/_dizzy/{{.URL}}/"><i class="mdl-color-text--blue-grey-400 material-icons" role="presentation">folder</i>{{.KindName}}</a>
{{ end }}
        </nav>
      </div>

      <main class="mdl-layout__content mdl-color--grey-100">

{{ "{{template \"content\" .}}" }}

      </main>

    </div>

    <script src="https://code.getmdl.io/1.3.0/material.min.js"></script>

  </body>
</html>
{{ "{{ end }}" }}

`

const TopTemplateFile = "top.tmpl"
const TopTemplate = `{{ "{{ define \"title\" }}" }}
Home[dizzy auto generated:{{generated}}]
{{ "{{ end }}" }}

{{ "{{ define \"content\" }}" }}

<h4>
     datastore is lovely.and easy access.
</h4>

18 Jun 2018.</br>
<a href="https://cloud.google.com/appengine/docs/standard/go/datastore/reference">The datastore package</a>

<pre>

support 

- signed integers (int, int8, int16, int32 and int64),
- bool,
- string,
- float32 and float64,

not supported

- []byte (up to 1 megabyte in length),
- ByteString,
- appengine.BlobKey,
- appengine.GeoPoint,

I do not plan support it

- *Key,
- time.Time (stored with microsecond precision),
- any type whose underlying type is one of the above predeclared types,
- structs whose fields are all valid value types,
- slices of any of the above.

</pre>

{{ "{{ end }}" }}

`

const ViewTemplateFile =  "view.tmpl"
const ViewTemplate = `{{ "{{define \"title\"}} "}}
{{ .Kind.KindName }} View[dizzy auto generated:{{generated}}]
{{ "{{end}}" }}

{{ "{{define \"content\"}}" }}
<table class="mdl-data-table mdl-js-data-table mdl-shadow--2dp">
  <thead>
    <tr>

      <th class="mdl-data-table__cell--non-numeric">
      {{ "{{if ne .Prev 0}}" }}
              <a href="/_dizzy/{{.Kind.URL}}/?page={{"{{.Prev}}"}}" title="Older">
                <button class="mdl-button mdl-button--icon  mdl-button--primary">
                   <i class="material-icons md-dark">arrow_back</i>
                </button>
              </a>
      {{ "{{end}}" }}
      </th>
{{ range .Kind.Fields }}

{{ if eq .Display true }}
        <th
{{if ne .Type 10}}
        {{ "class=\"mdl-data-table__cell--non-numeric\"" }}
{{end}}
        >
      {{ .Name }} </th>
{{ end }}

{{ end }}

      <th class="mdl-data-table__cell--non-numeric">Create/Update</th>

      <th>
      {{ "{{if ne .Max true}}" }}
              <a href="/_dizzy/{{.Kind.URL}}/?page={{"{{.Next}}"}}" title="Newer">
                <button class="mdl-button mdl-button--icon  mdl-button--primary">
                  <i class="material-icons md-dark" role="presentation">arrow_forward</i>
                </button>
              </a>
      {{ "{{ end }}" }}
      </th>
    </tr>
  </thead>
  <tbody>

{{ "{{ range .List}}" }}
    <tr>
        <td class="mdl-data-table__cell--non-numeric">
        <form method="post" action="/_dizzy/{{.Kind.URL}}/edit">
          <input type="hidden" name="keyId" value="{{"{{.Key.StringID}}"}}" />
          <button class="mdl-button mdl-js-button mdl-button--raised mdl-js-ripple-effect mdl-button--icon mdl-button--primary">
            <i class="material-icons">edit</i>
          </button>
        </form>
      </td>

{{ range .Kind.Fields }}
{{ if eq .Display true }}
      <td class="mdl-data-table__cell--non-numeric">
      {{"{{"}} .{{ .Name }} {{"}}"}}
      </td>
{{ end }}
{{ end }}

      {{ "<td class=\"mdl-data-table__cell--non-numeric\">{{convertDate .CreatedAt}}<br>{{convertDate .UpdatedAt}}</td>" }}

      <td>
        <form method="post" action="/_dizzy/{{.Kind.URL}}/delete" onsubmit="return confirm('realy?');">
          <input type="hidden" name="keyId" value="{{"{{.Key.StringID}}"}}" />
          <input type="hidden" name="Version" value="{{"{{.Version}}"}}" />
          <button class="mdl-button mdl-js-button mdl-button--raised mdl-js-ripple-effect mdl-button--icon mdl-button--accent">
            <i class="material-icons">remove</i>
          </button>
        </form>
      </td>

    </tr>
{{ "{{ end }}" }}

  </tbody>
</table>

<!-- view template -->
<a href="/_dizzy/{{.Kind.URL}}/create">
<button id="add-content" class="mdl-button mdl-js-button mdl-button--fab mdl-button--primary">
  <i class="material-icons">add</i>
</button>
</a>

<script>
</script>

{{ "{{ end }}" }}

`

var TemplateStringMap map[string]string

func init() {

    TemplateStringMap = make(map[string]string)

    TemplateStringMap[AppTemplateFile] = AppTemplate
    TemplateStringMap[DizzyGoTemplateFile] = DizzyGoTemplate
    TemplateStringMap[EditTemplateFile] = EditTemplate
    TemplateStringMap[ErrorTemplateFile] = ErrorTemplate
    TemplateStringMap[AccessTemplateFile] = AccessTemplate
    TemplateStringMap[HandlerGoTemplateFile] = HandlerGoTemplate
    TemplateStringMap[IndexTemplateFile] = IndexTemplate
    TemplateStringMap[LayoutTemplateFile] = LayoutTemplate
    TemplateStringMap[TopTemplateFile] = TopTemplate
    TemplateStringMap[ViewTemplateFile] = ViewTemplate
}

func generated() string {
	t := time.Now()
    if t.IsZero() {
        return "None"
    }
    return t.Format(time.RFC3339)
}

func createGenerateTemplate(name string) (*template.Template,error ) {

    funcMap := template.FuncMap{
        "generateInputHTML" : GenerateInputHTML,
        "generated" : generated,
    }
    var tmpl *template.Template
    var err error
    if UseTemplateFile {
        f := TemplateDirectory + "/" + name
        tmpl, err = template.New(name).Funcs(funcMap).ParseFiles(f)
    } else {
        val,ok := TemplateStringMap[name]
        if !ok {
            return nil, fmt.Errorf("TemplateStringMap not exists[%s]\n",name)
        }
        tmpl, err = template.New(name).Funcs(funcMap).Parse(val)
    }

    return tmpl,err
}