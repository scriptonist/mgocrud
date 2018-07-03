package generate

import (
	"io"
	"strings"
	"text/template"
)

var funcMaps map[string]interface{}

func init() {
	funcMaps = template.FuncMap{
		"toLower": strings.ToLower,
	}
}

const createFuncTemplate = `func ({{ printf "%.1s" .StructName | toLower }}  *{{.StrucName}}) Create(s *mgo.Session,database,collection string) error{
	err := s.DB(database).C(collection).Insert({{ printf "%.1s" .StructName | toLower }})
	if err != nil{
		return err
	}
	return err
}
`

func generateCreateMethod(s *structType, f io.ReadWriter) error {
	tmpl, err := template.New("createMethod").Funcs(funcMaps).Parse(createFuncTemplate)
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, s)
	if err != nil {
		return err
	}
	return nil
}
