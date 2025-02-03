package common

import (
	"basement/main/internal/logg"
	"fmt"
)

const (
	THING_ITEM = iota
	THING_BOX
	THING_SHELF
	THING_AREA
)

// ValidThing validates string and returns const int (THING_ITEM, THING_BOX, THING_SHELF, THING_AREA).
// On error logs warning and returns default thing "0".
func MustValidThing(thing string) int {
	t, err := ValidThing(thing)
	if err != nil {
		logg.Err(err)
		logg.Warningf(`using THING_ITEM="%i"`, thing, THING_ITEM)
		return THING_ITEM
	}
	return t
}

// ValidThing validates string and returns const int (THING_ITEM, THING_BOX, THING_SHELF, THING_AREA).
func ValidThing(thing string) (int, error) {
	switch thing {
	case "item":
		return THING_ITEM, nil
	case "box":
		return THING_BOX, nil
	case "shelf":
		return THING_SHELF, nil
	case "area":
		return THING_AREA, nil
	default:
		return 0, logg.NewError(`thing "` + thing + `" is not valid`)
	}
}

// ValidThingString validates const int (THING_ITEM, THING_BOX, THING_SHELF, THING_AREA) and returns string.
func ValidThingString(thing int) (out string, err error) {
	switch thing {
	case THING_ITEM:
		return "item", err
	case THING_BOX:
		return "box", err
	case THING_SHELF:
		return "shelf", err
	case THING_AREA:
		return "area", err
	default:
		return out, logg.NewError(fmt.Sprintf(`thing "%d" is not valid, using "item"`, thing))
	}
}
