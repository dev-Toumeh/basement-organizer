package common

import (
	"basement/main/internal/logg"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

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
		log.Printf("Error converting string to int64: %v", err)
		return 0
	}
	return i
}

func StringToFloat32(value string) float32 {
	floatValue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		log.Printf("Error converting string to int64: %v", err)
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
		logg.Err(err)
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
