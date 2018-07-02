package generate

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"log"
	"sort"
)

// Opts specifies the options that generate accepts
type Opts struct {
	Filename       string
	GenAnotherFile bool
}

// Generate CRUD function for the file given
func Generate(opts *Opts) error {
	var s scanner.Scanner

	fcontent, err := ioutil.ReadFile(opts.Filename)
	if err != nil {
		log.Println(err)
		return err
	}
	counts := make(map[string]int)
	fset := token.NewFileSet()
	file := fset.AddFile(opts.Filename, fset.Base(), len(fcontent))
	s.Init(file, fcontent, nil, scanner.ScanComments)
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		if tok == token.IDENT {
			counts[lit]++
			fmt.Println(lit)
		}

	}
	type pair struct {
		s string
		n int
	}
	pairs := make([]pair, 0, len(counts))
	for s, n := range counts {
		pairs = append(pairs, pair{s, n})
	}

	sort.Slice(pairs, func(i, j int) bool { return pairs[i].n > pairs[j].n })
	for i := 0; i < len(pairs); i++ {
		fmt.Println(pairs[i].s, pairs[i].n)
	}
	return nil
}
