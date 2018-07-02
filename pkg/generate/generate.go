package generate

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"strings"

	"github.com/fatih/structtag"
)

// Opts specifies the options that generate accepts
type Opts struct {
	Filename       string `bson:"filename" json:"filename,omitempty"`
	GenAnotherFile bool   `json:"gen_another_file,omitempty"`
}

// StructType encompasses a go struct's name and it's corresponding node
type structType struct {
	name string
	node *ast.StructType
}

// Generate CRUD function for the file given
func Generate(opts *Opts) error {
	b, err := ioutil.ReadFile(opts.Filename)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, opts.Filename, b, parser.ParseComments)
	if err != nil {
		return err
	}
	ast.Inspect(file, func(n ast.Node) bool {
		s, ok := n.(*ast.StructType)
		if !ok {
			return true
		}
		generateCRUD(&structType{
			node: s,
		})
		return true
	})

	return nil
}

// ProcessStruct processes a struct to generate the CRUD function for the given struct
func generateCRUD(s *structType) (string, error) {
	collectionName := s.name
	// A Map containing the key and type of the key
	documentKeys := make(map[string]string)
	// Extract fields and extract bson tags
	for _, field := range s.node.Fields.List {
		if field.Tag != nil {
			// Remove backticks at two ends
			ns := strings.TrimPrefix(strings.TrimSuffix(field.Tag.Value, "`"), "`")

			tags, err := structtag.Parse(ns)
			if err != nil {
				panic(err)
			}
			bsontag, err := tags.Get("bson")
			if err != nil {
				log.Println(err)
			} else {

			}

		}
	}
	return "", nil
}
