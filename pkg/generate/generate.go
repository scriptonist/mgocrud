package generate

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Opts specifies the options that generate accepts
type Opts struct {
	Filename       string `bson:"filename" json:"filename,omitempty"`
	GenAnotherFile bool   `json:"gen_another_file,omitempty"`
}

// StructType encompasses a go struct's name and it's corresponding node
type structType struct {
	PackageName string
	Name        string
	Node        *ast.StructType
}

// Generate CRUD function for the file given
func Generate(opts *Opts) error {
	b, err := ioutil.ReadFile(opts.Filename)
	if err != nil {
		return err
	}
	// var workingDir = filepath.Dir(opts.Filename)
	absPath, err := filepath.Abs(opts.Filename)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, opts.Filename, b, parser.ParseComments)
	if err != nil {
		return err
	}
	// ast.Print(fset, file)
	structs := collectStructs(file)
	for _, val := range structs {
		var buf bytes.Buffer
		generateCRUD(val, &buf)
		formatedContent, err := format.Source(buf.Bytes())
		if err != nil {
			return err
		}
		f, err := os.OpenFile(absPath, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}

		_, err = f.Write(formatedContent)
		if err != nil {
			return err
		}

		// err = ioutil.WriteFile(filepath.Join(workingDir, val.Name+"CRUD.go"), formatedContent, 0666)
		// if err != nil {
		// 	log.Println(err)
		// 	return err
		// }
	}

	return nil
}

func collectStructs(node ast.Node) map[token.Pos]*structType {
	structs := make(map[token.Pos]*structType, 0)
	var structName string
	var packagename string
	collectStructs := func(n ast.Node) bool {
		switch ntype := n.(type) {
		case *ast.File:
			packagename = ntype.Name.Name
			return true
		case *ast.TypeSpec:
			structName = ntype.Name.Name
			s, ok := ntype.Type.(*ast.StructType)
			if !ok {
				return true
			}
			structs[s.Pos()] = &structType{
				Name:        structName,
				Node:        s,
				PackageName: packagename,
			}
		default:
			// fmt.Printf("%T\n", ntype)
			return true
		}
		return true
	}
	ast.Inspect(node, collectStructs)
	return structs
}

// ProcessStruct processes a struct to generate the CRUD function for the given struct
func generateCRUD(s *structType, w io.ReadWriter) (string, error) {
	// err := generateFileStub(s, w)
	// if err != nil {
	// 	log.Println(err)
	// 	return "", err
	// }

	err := generateFuncs(s, w)
	if err != nil {
		log.Println(err)
	}
	return "", nil
}
