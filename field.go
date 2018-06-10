package main

import (
	"strings"
	"go/ast"
	"fmt"
	"strconv"
	"log"
)

//タグ指定用
const TagPrefix = "dizzy:"

type Field struct {
	Name string
	DisplayName string
	Type int
	TypeName string
	Default string
	Display bool
	Editable bool
	Required bool
}

//引数の構造体データからフィールドに変換
func createFields(str *ast.StructType) ([]*Field,error) {

	defFields := make([]*Field,0)
	//kind からフィールドを生成
	fieldList := str.Fields
	fields := fieldList.List

	hasMeta := false
	//フィールド数繰り返す
	for _,field := range fields {
		//Type
		f,m,err := NewField(field)
		if err != nil {
			return nil,err
		}
		if f != nil {
			defFields = append(defFields,f)
		}
		if m {
			hasMeta = m
		}
	}

	if !hasMeta {
		fmt.Println("Not found ds.Meta?")
	}

	return defFields,nil
}

func NewField(field *ast.Field) (*Field,bool,error) {

	f := &Field{
		Display: false,
		Editable: true,
		Default:"",
		Required:false,
	}

	//テスト用
	if field == nil {
		return f, false, nil
	}

	expr := field.Type
	switch expr.(type) {
	case *ast.Ident:

		id := expr.(*ast.Ident)
		f.Name = field.Names[0].Name
		f.DisplayName = f.Name
		switch id.Name {
		case "string":
			f.Type = StringField
		case "int", "int8", "int16", "int32", "int64":
			f.Type = IntField
		case "float32", "float64":
			f.Type = FloatField
		case "bool":
			f.Type = BoolField
		case "[]byte":
		case "ByteString":
		}

		//デフォルト値を設定
		f.TypeName = id.Name

		//*ast.SelectorExpr -> 埋め込み構造体
	case *ast.SelectorExpr:
		se := expr.(*ast.SelectorExpr)
		if se.Sel.Name == "Meta" {
			if id, iok := se.X.(*ast.Ident); iok {
				if id.Name == "ds" {
					return nil, true, nil
				}
			}
		}
	default:
		fmt.Printf("Not support field.[%#v]\n", expr)
		return nil, false, nil
	}

	lit := field.Tag
	//タグなしはそのまま返す
	if lit == nil {
		return f, false, nil
	}

	org := lit.Value
	idx := strings.Index(org, TagPrefix)
	if idx == -1 {
		return f, false, nil
	}

	dizzy := org[idx:]
	idx = strings.Index(dizzy, "\"")
	if idx == -1 {
		return f, false, fmt.Errorf("dizzy tag error[%s]", dizzy)
	}

	line := dizzy[idx+1:]
	idx = strings.Index(line, "\"")
	if idx == -1 {
		return f, false, fmt.Errorf("dizzy tag error[%s]", line)
	}

	csv := line[:idx]

	err := setTag(f, csv)
	if err != nil {
		return f, false, err
	}
	return f, false, nil
}

func setTag(f *Field,csv string) error {

	liner,err := NewLiner(csv)
	if err != nil {
		return fmt.Errorf("field tag dizzy error[%s]",csv)
	}

	// TODO 長さ(数値時は範囲？)、選択
	for key,val := range liner {
		switch key {
		case "name":
			f.DisplayName = val
		case "default":
			f.Default = val
		case "required":
			v,err := strconv.ParseBool(val)
			if err != nil {
				continue
			}
			f.Required = v
		case "editable":
			v,err := strconv.ParseBool(val)
			if err != nil {
				continue
			}
			f.Editable = v
		case "display":
			v,err := strconv.ParseBool(val)
			if err != nil {
				continue
			}
			f.Display = v
		default:
			log.Printf("not support dizzy tag %s=[%s]",key,val)
		}
	}


	return nil
}