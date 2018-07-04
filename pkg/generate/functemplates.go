package generate

import (
	"io"
	"log"
	"strings"
	"text/template"
)

var funcMaps map[string]interface{}

type FileStub struct {
	PackageName string
	Imports     []string
}

var requiredImports = []string{"github.com/globalsign/mgo"}

func init() {
	funcMaps = template.FuncMap{
		"toLower": strings.ToLower,
	}
}

const fileStubTemplate = `package {{ .PackageName }}
	import(
		{{ range $ind,$pkg := .Imports}}
			"{{$pkg}}"
		{{end}}
	)
`

const createFuncTemplate = `func ({{ printf "%.1s" .Name | toLower }}  *{{.Name}}) Create(mgosession *mgo.Session,database,collection string) error{
	err := mgosession.DB(database).C(collection).Insert({{ printf "%.1s" .Name | toLower }})
	if err != nil{
		return err
	}
	return err
}
`

const readFuncTemplate = `func ({{ printf "%.1s" .Name | toLower }}  *{{.Name}}) Read(mgosession *mgo.Session,database,collection string,selector bson.M) (*[]{{ .Name }},error){
	var results []{{.Name}}
	err := mgosession.DB(database).C(collection).Find(selector).All(&results)
	if err != nil{
		return nil,err
	}
	return &results,nil
}
`

const updateFuncTemplate = `func ({{ printf "%.1s" .Name | toLower }}  *{{.Name}}) Update(mgosession *mgo.Session,database,collection string,selector,change bson.M) error {
	err := mgosession.DB(database).C(collection).Update(selector,bson.M{"$set":change})
	if err != nil{
		return err
	}
	return nil
}
`
const deleteFuncTemplate = `func ({{ printf "%.1s" .Name | toLower }}  *{{.Name}}) Delete(mgosession *mgo.Session,database,collection string,selector bson.M) (*mgo.ChangeInfo,error){
	info,err := mgosession.DB(database).C(collection).RemoveAll(selector)
	if err != nil{
		return nil,err
	}
	return info,nil
}
`

func generateFileStub(s *structType, f io.ReadWriter) error {
	filestub := FileStub{
		PackageName: s.PackageName,
		Imports:     requiredImports,
	}
	tmpl, err := template.New("createfilestub").Parse(fileStubTemplate)
	if err != nil {
		return err
	}
	err = tmpl.Execute(f, filestub)
	if err != nil {
		return err
	}

	return nil
}

var funcTemplates = []string{createFuncTemplate, readFuncTemplate, updateFuncTemplate, deleteFuncTemplate}

func generateFuncs(s *structType, f io.ReadWriter) error {
	for _, tstring := range funcTemplates {
		tmpl, err := template.New("Func").Funcs(funcMaps).Parse(tstring)
		if err != nil {
			log.Println(err)
			return err
		}
		err = tmpl.Execute(f, s)
		if err != nil {
			log.Println(err)
			return err
		}

	}
	return nil
}
