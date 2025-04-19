package env

import (
	"basement/main/internal/logg"
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

func parseConfigFile(configFile string, config *Configuration) map[string]string {
	b, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	r := bytes.NewReader(b)
	out, _ := parseWithReader(r, config)

	return out
}

func parseWithReader(reader io.Reader, config *Configuration) (out map[string]string, errs []error) {
	b, _ := io.ReadAll(reader)

	out = make(map[string]string)
	lines := strings.Split(string(b), "\n")
	for lnr, l := range lines {
		opt, err := parseLine(l)
		if err != nil {
			errs = append(errs, err)
			if errors.Is(err, lineIsComment) || errors.Is(err, lineIsEmpty) {
				continue
			}
			logg.Warningf("ignoring line %d: %s ", lnr+1, logg.CleanLastError(err))
			continue
		}
		out[opt.Arg] = opt.Val
	}
	return out, errs
}

type option struct {
	Arg string
	Val string
}

var lineIsComment = errors.New("line is a comment")
var lineIsEmpty = errors.New("line is empty")

func parseLine(line string) (out option, err error) {
	// logg.Debugf("parsing line %d: \"%s\"", lnr, line)
	noSpaces := strings.TrimSpace(line)
	if len(noSpaces) == 0 {
		return out, logg.Errorf("%w \"%s\"", lineIsEmpty, line)
	}
	if strings.HasPrefix(noSpaces, "#") {
		return out, logg.Errorf("%w \"%s\"", lineIsComment, line)
	}
	maybecomment := strings.Split(noSpaces, "#")
	noComment := maybecomment[0]

	optval := strings.Split(noComment, "=")
	if len(optval) != 2 {
		return out, logg.NewError("option or value missing in \"" + line + "\"")
	}
	opt := strings.TrimSpace(optval[0])
	if len(opt) == 0 {
		return out, logg.NewError("len(opt) == 0 in \"" + line + "\"")
	}
	val := strings.TrimSpace(optval[1])
	if len(val) == 0 {
		return out, logg.NewError("len(val) == 0 in \"" + line + "\"")
	}
	// fmt.Printf("parsed %s: %s\n", opt, val)
	out.Arg = opt
	out.Val = val
	return out, nil
}

func applyParsedOptions(config *Configuration, parsedMap map[string]string) (errs []error) {
	// Enable logs first if present.
	val, debugLogs := parsedMap["debugLogsEnabled"]
	if debugLogs {
		if val == "true" {
			logg.EnableDebugLogger()
		}
	}
	val, infoLogs := parsedMap["infoLogsEnabled"]
	if infoLogs {
		if val == "true" {
			logg.EnableInfoLogger()
		}
	}
	logg.Debug("overriding parsed options from config file")

	for fieldName, mapValue := range parsedMap {
		fieldm, ok := config.fieldValues[fieldName]
		if !ok {
			logg.Debugf("skip fieldName=%s, not in parsed map", fieldName)
			option := reflect.ValueOf(*config).FieldByName(fieldName)
			if !option.IsValid() {
				logg.Fatalf("option \"%s\" is invalid", fieldName)
			}
			continue
		}

		logg.Debugf("fieldname: %s, fieldsetter: %s kind: %s, value:%s", fieldName, fieldm.Setter, fieldm.Kind.String(), mapValue)
		if fieldm.Setter == "" || mapValue == "" {
			logg.Debugf("skip setter=%s, value=%s", fieldm.Setter, mapValue)
			continue
		}

		switch fieldm.Kind {
		case reflect.String:
			SetUnexportedField(reflect.ValueOf(config).Elem().FieldByName(fieldName), mapValue)
			logg.Debugf("config config.dbPath=%s", config.dbPath)
			break
		case reflect.Int:
			val, err := strconv.Atoi(mapValue)
			if err != nil {
				logg.Err(err)
				logg.Fatalf("invalid int value \"%s\" in field \"%s\"", mapValue, fieldName)
			}
			SetUnexportedField(reflect.ValueOf(config).Elem().FieldByName(fieldName), val)
			break
		case reflect.Bool:
			var val bool
			switch mapValue {
			case "false":
				val = false
				break
			case "true":
				val = true
				break
			default:
				logg.Fatalf("invalid bool value \"%s\" in field \"%s\". Must be \"true\" or \"false\"", mapValue, fieldName)
			}
			SetUnexportedField(reflect.ValueOf(config).Elem().FieldByName(fieldName), val)
		default:
			logg.Debug("cont")
			continue
		}
		// logg.Debugf("Calling %s(%s)", fieldm.Setter, mapValue)
	}

	return errs
}

// From "How to access unexported struct fields" https://stackoverflow.com/a/60598827.
func GetUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

// From "How to access unexported struct fields" https://stackoverflow.com/a/60598827.
func SetUnexportedField(field reflect.Value, value interface{}) {
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(value))
}
