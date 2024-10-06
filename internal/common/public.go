package common

import (
	"basement/main/internal/logg"
	"fmt"
	"log"
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
