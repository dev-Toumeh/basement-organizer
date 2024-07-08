package main

import (
	"basement/main/internal/templates"
	"fmt"
	"go/format"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	TEMPLATES_DIR            string = "../internal/templates/"
	TEMPLATES_CONSTANTS_PATH string = "../internal/templates/constants.go"
	TARGET_TEMPLATES_DIR     string = templates.TEMPLATE_DIR
)

// GenerateTemplateNames automatically generates cons variables from parsed template definitions.//
// File definitions are parsed from "TEMPLATES_DIR".
//
// Creates "constants.go" file defined in "TEMPLATES_CONSTRANTS_PATH".
func main() {
	template, paths, err := templates.ParseDirectory(TEMPLATES_DIR)
	if err != nil {
		panic(err)
	}
	err = GenerateTemplateNames(template, paths)
	if err != nil {
		panic(err)
	}
}

// GenerateTemplateNames automatically generates cons variables from parsed template definitions.//
// File definitions are parsed from "TEMPLATES_DIR".
//
// Creates "constants.go" file defined in "TEMPLATES_CONSTRANTS_PATH".
func GenerateTemplateNames(template *template.Template, parsedPaths []string) error {
	definedTemplates := template.DefinedTemplates()
	templates2 := strings.TrimPrefix(definedTemplates, "; defined templates are: ")
	templates3 := strings.Split(templates2, ", ")
	if len(templates3) == 0 {
		panic("templates3 has 0 length")
	}

	cleanTemplateNames := make([]string, len(templates3))
	cleanTemplateConstVariables := make([]string, len(templates3))
	parsedPathsMap := map[string]string{}
	// log.Println("Base paths:")
	for _, v := range parsedPaths {
		parsedPathsMap[filepath.Base(v)] = v
		// log.Println(filepath.Base(v))
	}
	// log.Println(parsedPathsMap)
	cleanTemplateDefinitionSourceFiles := map[string]string{}

	for i, s := range templates3 {
		noSpaceQuotes := strings.TrimFunc(s, func(r rune) bool {
			switch r {
			case '"':
				return true
			case ' ':
				return true
			default:
				return false
			}
		})
		cleanTemplateNames[i] = noSpaceQuotes
		fileSource := template.Lookup(noSpaceQuotes).Tree.ParseName
		mappedPath, ok := parsedPathsMap[filepath.Base(fileSource)]
		cleanPath := filepath.ToSlash(mappedPath)
		if ok {
			cleanTemplateDefinitionSourceFiles[noSpaceQuotes] = cleanPath
			// log.Println(filepath.ToSlash(mappedPath))
		}
		// os.PathSeparator
		// filepath.ListSeparator
		// cleanTemplateDefinitionSourceFiles[noSpaceQuotes] = template.Lookup(noSpaceQuotes).Tree.ParseName

		withUnderscores := strings.Map(func(r rune) rune {
			switch r {
			case '-':
				return '_'
			case '.':
				return '_'
			default:
				return r
			}

		}, noSpaceQuotes)
		upperCase := strings.ToUpper(withUnderscores)
		// log.Println(filepath.ToSlash(filepath.Clean(templates.TEMPLATE_DIR)))
		// log.Println(cleanPath)
		targetTmplDir := filepath.ToSlash(filepath.Clean(TARGET_TEMPLATES_DIR))
		// b, a, _ := strings.Cut(targetTmplDir, cleanPath)
		// log.Println(a, b)
		c := strings.Split(cleanPath, targetTmplDir)
		// log.Println(c[1])
		// log.Println(strings.Index(cleanPath, targetTmplDir))
		// d, err := filepath.Rel(targetTmplDir,filepath.ToSlash( c[1]))
		targetPath := filepath.ToSlash(filepath.Join(targetTmplDir, c[1]))
		// log.Println(d)
		log.Println(targetPath)
		// filepath.re
		// if err != nil {
		// 	panic(err)
		// }
		// log.Println("asdfasdfk:", p)
		constVariable := fmt.Sprintf(`const %s string = %s // "%s"`, upperCase, s, targetPath)

		// log.Println(constVariable)
		cleanTemplateConstVariables[i] = constVariable
	}
	// log.Println(parsedPaths)
	// log.Println()
	log.Println()
	// p, err := filepath.Rel(".", "internal/templates/")
	// // filepath.re
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println("asdfasdfk:", p)
	log.Println()
	sort.Strings(cleanTemplateConstVariables)
	// strings.HasSuffix()

	file, err := os.OpenFile(TEMPLATES_CONSTANTS_PATH, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	// file, err := os.OpenFile(TEMPLATES_CONSTANTS_PATH, os.O_RDWR|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	// var line int
	fstats, err := file.Stat()
	if err != nil {
		panic(err)
	}
	// bufio.NewReadWriter(file, file)

	// var buf [fstats.Size()]byte
	buf := make([]byte, fstats.Size())
	// bytes.NewBuffer(buf)
	eof, err := file.Read(buf)
	if err != nil {
		panic(err)
	}
	log.Println("eof:", eof)
	// log.Println("file contents:\n", string(buf))

	newFileContentBegin := "// THIS FILE IS AUTO GENERATED!\npackage templates\n\n"
	newFileContentConstVariables := strings.Join(cleanTemplateConstVariables, "\n")
	// newFileContentEnd := ")"
	newFileContentEnd := ""
	newFileContent := newFileContentBegin + newFileContentConstVariables + newFileContentEnd
	newFileContentFormatted, err := format.Source([]byte(newFileContent))
	if err != nil {
		panic(err)
	}
	log.Printf("New file content to be written in '%s': %s\n", TEMPLATES_CONSTANTS_PATH, string(newFileContentFormatted))

	_, err = file.WriteString(string(newFileContentFormatted))
	if err != nil {
		panic(err)
	}
	return nil
}
