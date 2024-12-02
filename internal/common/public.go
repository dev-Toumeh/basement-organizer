package common

import (
	"basement/main/internal/logg"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofrs/uuid/v5"
)

func CheckIDs(mainId string, wrapperId string) (uuid.UUID, uuid.UUID, error) {
	var err error
	id := uuid.Nil
	boxId := uuid.Nil

	lengID := len(mainId)
	lengBoxID := len(wrapperId)

	if lengBoxID != 0 {
		boxId, err = uuid.FromString(wrapperId)
		if err != nil {
			logg.Errf("error while converting the boxId to type uuid: %v", err)
		}
	}

	if lengID == 0 {
		id, err = uuid.NewV4()
		if err != nil {
			return id, boxId, fmt.Errorf("error while generating the new item uuid: %w", err)
		}
		return id, boxId, nil
	} else {
		id, err = uuid.FromString(mainId)
		if err != nil {
			return id, boxId, fmt.Errorf("error while converting the itemId to type uuid: %w", err)
		}
		return id, boxId, nil
	}
}

func StringToInt(value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Error converting string to int: %v", err)
		return 0
	}
	return i
}

func StringToFloat32(value string) float32 {
	floatValue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		log.Printf("Error converting string to Float32: %v", err)
		return 0
	}
	return float32(floatValue)
}

// ParsePicture returns base64 encoded string of picture uploaded if there is any
func ParsePicture(r *http.Request) string {
	logg.Info("Parsing multipart/form-data for picture")
	// 8 MB
	var maxSize int64 = 1000 * 1000 * 8
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		logg.Err(err)
		return ""
	}

	file, header, err := r.FormFile("picture")
	if header != nil {
		logg.Debug("picture filename:", header.Filename)
	}
	if err != nil {
		logg.Info(err)
		return ""
	}

	readbytes, err := io.ReadAll(file)
	logg.Debug("picture size:", len(readbytes)/1000, "KB")
	if err != nil {
		logg.Err(err)
		return ""
	}

	return base64.StdEncoding.EncodeToString(readbytes)
}

// parseQuantity returns by default at least 1
func ParseQuantity(value string) int64 {
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 1
	}
	return intValue
}

// MergeMaps takes a slice of maps with string keys and any value types, returning a single merged map
func MergeMaps[V any](mapsList []map[string]V) map[string]V {
	merged := make(map[string]V)

	for _, m := range mapsList {
		maps.Copy(merged, m)
	}

	return merged
}

// CheckEditMode checks the "edit" parameter in the request URL
// and returns true if "edit" equals "1" or "true", otherwise false
func CheckEditMode(r *http.Request) bool {
	editValue := r.URL.Query().Get("edit")

	if editValue == "1" {
		return true
	}

	isEdit, err := strconv.ParseBool(editValue)
	if err == nil && isEdit {
		return true
	}

	return false
}

// parseIDsFromFormWithKey parses r.Form by searching all HTML input elements that start with `key` name and returns a list of valid uuid.UUIDs
//
// `r.ParseForm()` must be called before using this function!
//
// Example:
//
//	// search for all ID values that start with "delete:" key
//	// like "delete:f47ac10b-58cc-0372-8567-0e02b2c3d479"
//	toDeleteIDs := parseIDsFromFormWithKey(r.Form, "delete")
func ParseIDsFromFormWithKey(form url.Values, key string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0)
	// logg.Debug("FormValues length: " + strconv.Itoa(len(form)))
	// i := 0
	for k := range form {
		// logg.Debug("FormValue[" + strconv.Itoa(i) + "]: \"" + k + "\"")
		if strings.Contains(k, key+":") {
			// logg.Debug("\"" + k + "\"" + " contains " + "\"" + key + "\"")
			idStr := strings.Split(k, fmt.Sprintf("%s:", key))
			if len(idStr) != 2 {
				return nil, logg.NewError(fmt.Sprintf("Wrong delete key value pair: '%v'", k))
			}
			// logg.Debug("clean value \"" + idStr[1] + "\"")
			id, err := uuid.FromString(idStr[1])
			if err != nil {
				logg.Err(err)
				return nil, logg.WrapErr(err)
			}
			ids = append(ids, id)
		}
		// i++
	}
	// logg.Debugf("clean ids: %v", ids)
	return ids, nil
}
