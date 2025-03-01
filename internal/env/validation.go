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
	"strings"
)

// Init returns false if some Get or Set methods are missing from struct.
func (config *Configuration) Init() bool {
	cfg := *config
	t := reflect.TypeOf(cfg)
	fmt.Println(t.Kind())
	if t.Kind() != reflect.Struct {
		panic("invalid config struct")
	}
	configStructName := t.Name()

	v := reflect.ValueOf(cfg)
	vv, ok := v.Interface().(Configuration)
	if ok {
		vvvt := reflect.TypeOf(&vv)
		// logg.Infof("method nums: \"%d\"", vvvt.NumMethod())

		excludeMethods := strings.Split(ignoreConfigMethodChecks, ",")
		// vvvv := reflect.ValueOf(&vv)
		for i := 0; i < vvvt.NumMethod(); i++ {
			methodName := vvvt.Method(i).Name
			// logg.Infof("method: \"%s\"=\"%v\"", methodName, vvvv.Method(i))
			if slices.Contains(excludeMethods, methodName) {
				continue
			}
			// logg.Debug(methodName)
			config.methods = append(config.methods, methodName)
		}
	}

	hasMissing := false
	excludeFields := strings.Split(ignoreConfigFieldChecks, ",")
	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		// logg.Infof("field: \"%s\"=\"%v\"", fieldName, v.Field(i))
		if slices.Contains(excludeFields, fieldName) {
			continue
		}

		expectedPublicFieldSetter := "Set" + strings.ToUpper(fieldName[:1]) + fieldName[1:]
		expectedPublicFieldGetter := strings.ToUpper(fieldName[:1]) + fieldName[1:]

		if !slices.Contains(config.methods, expectedPublicFieldSetter) {
			hasMissing = true
			logg.Err("Missing method \"" + expectedPublicFieldSetter + "\" in \"type " + configStructName + " struct {...}\". Implement this method.")
		}
		if !slices.Contains(config.methods, expectedPublicFieldGetter) {
			hasMissing = true
			logg.Err("Missing method \"" + expectedPublicFieldGetter + "\" in \"type " + configStructName + " struct {...}\". Implement this method.")
		}

		config.fields = append(config.fields, fieldName)
	}

	return hasMissing
}

func checkIfAllSettingsWereApplied(c *Configuration, appliedMethods []string) bool {
	allApplied := true
	for _, m := range c.Methods() {

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

	// Find the function "ApplyConfig()"
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != applyFunctionName {
			return true // Continue scanning
		}

		// Find all function calls inside "ApplyConfig()"
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
