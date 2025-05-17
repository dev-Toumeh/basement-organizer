package env

import (
	"basement/main/internal/logg"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"slices"
)

type fieldMetaData struct {
	Setter       string
	Getter       string
	Kind         reflect.Kind
	Value        string
	DefaultValue string
}

func validateInternals(applyConfigFileName string, applyConfigFuncName string) {
	developmentChecks(applyConfigFileName, applyConfigFuncName)
	ValidateOptions(configInstance)
}

// internal config checks
func developmentChecks(applyConfigFileName string, applyConfigFuncName string) {
	funcs := calledFunctionsFrom(applyConfigFileName, applyConfigFuncName)
	// logg.Infof("applied methods %v", funcs)
	allApplied := checkIfAllFieldSettersWereApplied(configInstance, funcs)
	if !allApplied {
		logg.Fatalf("Not all Config settings were applied in \"%s: %s()\"", applyConfigFileName, applyConfigFuncName)
	}
}

// collects all functions that were called in applyFunctionName.
func calledFunctionsFrom(filename string, applyFunctionName string) []string {
	// Read the file
	src, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// Parse the file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, src, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	foundAppliedSetters := []string{}

	// Find the function "applyFunctionName"
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != applyFunctionName {
			return true // Continue scanning
		}

		// Find all function calls inside
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			foundAppliedSetters = append(foundAppliedSetters, sel.Sel.Name)
			return true
		})
		// logg.Infof("foundAllMethods: %v", foundAppliedSetters)

		return false // Stop scanning other functions from this file.
	})

	return foundAppliedSetters
}

func checkIfAllFieldSettersWereApplied(c *Configuration, appliedMethods []string) bool {
	allApplied := true
	for _, m := range c.Methods() {

		// ignore methods that don't start with "Set"
		if m[:3] != "Set" {
			continue
		}

		// logg.Debugf("check method: %v", m)
		applied := slices.Contains(appliedMethods, m)
		if !applied {
			logg.Warning("Forgot to apply method \"" + m + "\" ?")
			allApplied = false
		}
	}

	return allApplied
}

func ValidateOptions(config *Configuration) (errors []error) {
	logg.Debug("validate combination of config options")
	var err error
	err = validateDefaultTableSize(config)
	if err != nil {
		errors = append(errors, err)
	}
	err = validateDBOptions(config)
	if err != nil {
		errors = append(errors, err)
	}
	return errors
}

func validateDefaultTableSize(config *Configuration) (err error) {
	// fmt.Printf("config.defaultTableSize=%d\n", config.defaultTableSize)
	validTableSize := config.defaultTableSize > 0

	if !validTableSize {
		err = logg.NewError(fmt.Sprintf("defaultTableSize must be above 0. defaultTableSize=%d", config.defaultTableSize))
	}
	return err
}

// validateDBOptions checks for consistency between different options regarding DB.
func validateDBOptions(config *Configuration) (err error) {
	invalidMemoryDB := (config.dbPath == ":memory:") && (config.useMemoryDB == false)
	invalidMemoryDB2 := (config.dbPath != ":memory:") && (config.useMemoryDB == true)
	emptyFileDBPath := (config.dbPath == "") && (config.useMemoryDB == false)

	usedOptions := fmt.Sprintf("Invalid options combination: \"dbPath=%s\" \"useMemoryDB=%t\"", config.dbPath, config.useMemoryDB)
	if emptyFileDBPath {
		msg := usedOptions + ". Can't set empty DB path without in-memory db. "
		return logg.NewError(msg)
	}
	if invalidMemoryDB {
		msg := usedOptions + ". Can't set DB path to \":memory:\". It's reserved only for in-memory db. "
		return logg.NewError(msg)
	}

	if invalidMemoryDB2 {
		msg := usedOptions + ". For using in-memory db, path must be set to \":memory:\". "
		return logg.NewError(msg)
	}
	return nil
}
