package main

import (
	"go/ast"
	"go/token"
	"go/parser"
	"strings"
	"fmt"
	"strconv"
)

const GeneratePrefix = "+DIZZY"

type Kind struct {
	TypeName string
	KindName string
	URL      string
	Cache    bool
	Content  Content
	Limit    int
	Owner   *ast.StructType
	Fields   []*Field
}

func createKinds(src string) ([]*Kind, string, error) {

	fset := token.NewFileSet()
	//構文木を取得
	f, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
	if err != nil {
		return nil, "", err
	}
	packageName := f.Name.Name

	//コメントからPREFIXのついているstructを見つける
	cm := ast.NewCommentMap(fset, f, f.Comments)
	//戻り値を生成
	rtn := make([]*Kind, 0)
	for n, cgs := range cm {
		for _, cg := range cgs {
			comment := cg.Text()
			//対象のコメントの場合
			if strings.Index(comment, GeneratePrefix) == 0 {
				//ちゃんとした指定だった場合
				t := getStruct(n)
				if t != nil {
					//Kind情報を生成
					kind,err := NewKind(comment, t)
					if err != nil {
						return nil,packageName,err
					}
					if kind != nil {
						rtn = append(rtn, kind)
					} else {
						//TODO ERROR
					}
				} else {
					return nil, packageName, fmt.Errorf("prefix target error")
				}
			}
		}
	}

	return rtn, packageName, nil
}

func getStruct(n ast.Node) *ast.TypeSpec {

	//Declの確認
	s, ok := n.(*ast.GenDecl)
	//指定がおかしい
	if !ok {
		return nil
	}

	if s.Specs == nil || len(s.Specs) != 1 {
		//2以上のときって何？
		return nil
	}

	spec := s.Specs[0]
	typeSpec, ok := spec.(*ast.TypeSpec)
	if !ok {
		return nil
	}

	workType := typeSpec.Type
	_, ok = workType.(*ast.StructType)
	if !ok {
		return nil
	}

	//TODO HasKey実装かを確認したい

	return typeSpec
}


func NewKind(line string, t *ast.TypeSpec) (*Kind,error) {

	name := "TestName"
	//nil is test value
	if t != nil {
		name = t.Name.Name
	}
	//default
	defName := name
	defURL := name
	defCache := true
	defContent := ContentNone
	defLimit := 10

	var liner Liner
	var err error
	//のみの指定
	if line == GeneratePrefix {
		liner, err = NewLiner("")
	} else {
		//+DIZZY(name=kindName,cache=true,content=none,limit=10)
		sIdx := strings.Index(line, "(")
		eIdx := strings.LastIndex(line, ")")
		//カッコがどっちにもあること
		if sIdx == -1 || eIdx == -1 {
			return nil, fmt.Errorf("dizzy Comment required()")
		}
		argsBuf := line[sIdx+1 : eIdx]

		liner, err = NewLiner(argsBuf)
	}
	if err != nil {
		return nil, err
	}

	for key, val := range liner {
		switch key {
		case "name":
			defName = val
		case "url":
			defURL = val
		case "cache":
			cache, err := strconv.ParseBool(val)
			if err == nil {
				defCache = cache
			} else {
				fmt.Printf("cache parse error(bool value)[%s]\n", val)
			}
		case "content":
			ct := ContentNone
			defContent = ct.Value(val)
		case "limit":
			limit, err := strconv.ParseInt(val, 10, 64)
			if err == nil {
				defLimit = int(limit)
			} else {
				fmt.Printf("limit parse error(int value)[%s]\n", val)
			}
		default:
			//TODO no support error?
		}
	}

	kind := &Kind{
		TypeName: name,
		KindName: defName,
		URL:      defURL,
		Cache:    defCache,
		Content:  defContent,
		Limit:    defLimit,
	}

	// nil is test value
	if t == nil {
		return kind, nil
	}

	str := t.Type.(*ast.StructType)
	kind.Owner = str

	fields, err := createFields(str)
	if err != nil {
		return nil, err
	}

	kind.Fields = fields
	return kind,nil
}