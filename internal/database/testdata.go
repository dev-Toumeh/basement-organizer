package database

import (
	"basement/main/internal/areas"
	"basement/main/internal/boxes"
	"basement/main/internal/common"
	"basement/main/internal/items"
	"basement/main/internal/shelves"

	"github.com/gofrs/uuid/v5"
)

var ITEM_VALID_UUID_1 uuid.UUID = uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174001"))
var ITEM_VALID_UUID_2 uuid.UUID = uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174002"))
var ITEM_VALID_UUID_3 uuid.UUID = uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174003"))

// var ITEM_VALID_UUID_4 uuid.UUID = uuid.Must(uuid.FromString("323e4567-e89b-12d3-a456-426614174000"))

var BOX_VALID_UUID_1 uuid.UUID = uuid.Must(uuid.FromString("223e4567-e89b-12d3-a456-426614174001"))
var BOX_VALID_UUID_2 uuid.UUID = uuid.Must(uuid.FromString("223e4567-e89b-12d3-a456-426614174002"))
var BOX_VALID_UUID_3 uuid.UUID = uuid.Must(uuid.FromString("223e4567-e89b-12d3-a456-426614174003"))
var BOX_VALID_UUID_4 uuid.UUID = uuid.Must(uuid.FromString("223e4567-e89b-12d3-a456-426614174114"))

var SHELF_VALID_UUID_1 uuid.UUID = uuid.Must(uuid.FromString("323e4567-e89b-12d3-a456-426614174001"))
var SHELF_VALID_UUID_2 uuid.UUID = uuid.Must(uuid.FromString("323e4567-e89b-12d3-a456-426614174002"))
var SHELF_VALID_UUID_3 uuid.UUID = uuid.Must(uuid.FromString("323e4567-e89b-12d3-a456-426614174003"))
var SHELF_VALID_UUID_4 uuid.UUID = uuid.Must(uuid.FromString("323e4567-e89b-12d3-a456-426614174004"))
var SHELF_VALID_UUID_5 uuid.UUID = uuid.Must(uuid.FromString("323e4567-e89b-12d3-a456-426614174005"))
var SHELF_VALID_UUID_6 uuid.UUID = uuid.Must(uuid.FromString("323e4567-e89b-12d3-a456-426614174006"))

var AREA_VALID_UUID_1 uuid.UUID = uuid.Must(uuid.FromString("523e4567-e89b-12d3-a456-426614174001"))
var AREA_VALID_UUID_2 uuid.UUID = uuid.Must(uuid.FromString("523e4567-e89b-12d3-a456-426614174002"))
var AREA_VALID_UUID_3 uuid.UUID = uuid.Must(uuid.FromString("523e4567-e89b-12d3-a456-426614174003"))
var AREA_VALID_UUID_4 uuid.UUID = uuid.Must(uuid.FromString("523e4567-e89b-12d3-a456-426614174004"))
var AREA_VALID_UUID_5 uuid.UUID = uuid.Must(uuid.FromString("523e4567-e89b-12d3-a456-426614174005"))
var AREA_VALID_UUID_6 uuid.UUID = uuid.Must(uuid.FromString("523e4567-e89b-12d3-a456-426614174006"))

var ITEM_VALID_UUID uuid.UUID = uuid.Must(uuid.FromString("433e4567-e89b-12d3-a456-426614174000"))
var VALID_UUID_NOT_EXISTING uuid.UUID = uuid.Must(uuid.FromString("033e4567-e89b-12d3-a456-426614174000"))

const VALID_BASE64_PNG = "iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAIAAAACUFjqAAAAtUlEQVR4nGJp2XGEAQb+/P49J7cgY8ZUuAgTnDUjI7vUQf3m5e3/zxyakZENFW3ZcURGQf/r52cQBGfLKhq0bD/MqKBu+ufnL4jSm5e3QxjmtuEfPnyCGn7z8na4BAMDg7quZ2mia2thMAMDA0j31TMb4XJr5s2BMKr71zIwMLAwYIDq/rWMMDaLobs7mjTEWJC6CeuYjL08+o/eU9f1RDPgsbpTxvQpjMjBAvEucrAAAgAA//+Elk5AOfCu8QAAAABJRU5ErkJggg=="
const VALID_BASE64_PREVIEW_PNG = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAAEElEQVR4nGJaOLEJEAAA//8DkwG35JmAnAAAAABJRU5ErkJggg=="
const INVALID_BASE64_PNG = "invalid base 64"

func resetTestBoxes() {
	BOX_1 = &boxes.Box{
		BasicInfo: common.BasicInfo{
			ID:          BOX_VALID_UUID_1,
			Label:       "box 1",
			Description: "This is the sixth box",
			Picture:     VALID_BASE64_PNG,
			QRCode:      "uvwxyzabcdefg",
		},
	}
	BOX_2 = &boxes.Box{
		BasicInfo: common.BasicInfo{
			ID:          BOX_VALID_UUID_2,
			Label:       "box 2",
			Description: "This is the third box",
			Picture:     VALID_BASE64_PNG,
			QRCode:      "abababababcd",
		},
	}
	BOX_3 = &boxes.Box{
		BasicInfo: common.BasicInfo{
			ID:          BOX_VALID_UUID_3,
			Label:       "box 3",
			Description: "This is the fourth box",
			Picture:     VALID_BASE64_PNG,
			QRCode:      "efghefghefgh",
		},
	}
	BOX_4 = &boxes.Box{
		BasicInfo: common.BasicInfo{
			ID:          BOX_VALID_UUID_4,
			Label:       "box 4",
			Description: "This is the fifth box",
			Picture:     VALID_BASE64_PNG,
			QRCode:      "ijklmnopqrst",
		},
	}
}

func testBoxes() []*boxes.Box {
	return []*boxes.Box{BOX_1, BOX_2, BOX_3, BOX_4}
}

var BOX_1 = &boxes.Box{
	BasicInfo: common.BasicInfo{
		ID:          BOX_VALID_UUID_1,
		Label:       "box 1",
		Description: "This is the sixth box",
		Picture:     VALID_BASE64_PNG,
		QRCode:      "uvwxyzabcdefg",
	},
}

var BOX_2 = &boxes.Box{
	BasicInfo: common.BasicInfo{
		ID:          BOX_VALID_UUID_2,
		Label:       "box 2",
		Description: "This is the third box",
		Picture:     VALID_BASE64_PNG,
		QRCode:      "abababababcd",
	},
}

var BOX_3 = &boxes.Box{
	BasicInfo: common.BasicInfo{
		ID:          BOX_VALID_UUID_3,
		Label:       "box 3",
		Description: "This is the fourth box",
		Picture:     VALID_BASE64_PNG,
		QRCode:      "efghefghefgh",
	},
}

var BOX_4 = &boxes.Box{
	BasicInfo: common.BasicInfo{
		ID:          BOX_VALID_UUID_4,
		Label:       "box 4",
		Description: "This is the fifth box",
		Picture:     VALID_BASE64_PNG,
		QRCode:      "ijklmnopqrst",
	},
}

var ITEM_1 = &items.Item{
	BasicInfo: common.BasicInfo{
		ID:          ITEM_VALID_UUID_1,
		Label:       "Item 1",
		Description: "Description for item 1",
		Picture:     VALID_BASE64_PNG,
		QRCode:      "QRCode1",
	},
	Quantity: 10,
	Weight:   "5.5",
}

var ITEM_2 = &items.Item{
	BasicInfo: common.BasicInfo{
		ID:          ITEM_VALID_UUID_2,
		Label:       "Item 2",
		Description: "Description for item 2",
		Picture:     VALID_BASE64_PNG,
		QRCode:      "QRCode2",
	},
	Quantity: 20,
	Weight:   "10.0",
}

var ITEM_3 = &items.Item{
	BasicInfo: common.BasicInfo{
		ID:          ITEM_VALID_UUID_3,
		Label:       "Item 3",
		Description: "Description for item 3",
		Picture:     VALID_BASE64_PNG,
		QRCode:      "QRCode3",
	},
	Quantity: 15,
	Weight:   "7.25",
}

func resetTestItems() {
	ITEM_1 = &items.Item{
		BasicInfo: common.BasicInfo{
			ID:          ITEM_VALID_UUID_1,
			Label:       "Item 1",
			Description: "Description for item 1",
			Picture:     VALID_BASE64_PNG,
			QRCode:      "QRCode1",
		},
		Quantity: 10,
		Weight:   "5.5",
	}

	ITEM_2 = &items.Item{
		BasicInfo: common.BasicInfo{
			ID:          ITEM_VALID_UUID_2,
			Label:       "Item 2",
			Description: "Description for item 2",
			Picture:     VALID_BASE64_PNG,
			QRCode:      "QRCode2",
		},
		Quantity: 20,
		Weight:   "10.0",
	}

	ITEM_3 = &items.Item{
		BasicInfo: common.BasicInfo{
			ID:          ITEM_VALID_UUID_3,
			Label:       "Item 3",
			Description: "Description for item 3",
			Picture:     VALID_BASE64_PNG,
			QRCode:      "QRCode3",
		},
		Quantity: 15,
		Weight:   "7.25",
	}
}

func testItems() []items.Item {
	return []items.Item{*ITEM_1, *ITEM_2, *ITEM_3}
}

var SHELF_1 = &shelves.Shelf{

	BasicInfo: common.BasicInfo{
		ID:             SHELF_VALID_UUID_1,
		Label:          "Test Shelf",
		Description:    "A shelf for testing",
		Picture:        VALID_BASE64_PNG,
		PreviewPicture: "",
		QRCode:         "testqrcode",
	},
	Height: 2.0,
	Width:  1.5,
	Depth:  0.5,
	Rows:   3,
	Cols:   4,
}

var SHELF_2 = &shelves.Shelf{
	BasicInfo: common.BasicInfo{
		ID:             SHELF_VALID_UUID_2,
		Label:          "Test Shelf 2",
		Description:    "A shelf for testing",
		Picture:        VALID_BASE64_PNG,
		PreviewPicture: "",
		QRCode:         "testqrcode",
	},
	Height: 2.0,
	Width:  1.5,
	Depth:  0.5,
	Rows:   3,
	Cols:   4,
}

var SHELF_3 = &shelves.Shelf{
	BasicInfo: common.BasicInfo{
		ID:             SHELF_VALID_UUID_3,
		Label:          "Test Shelf 3",
		Description:    "A shelf for testing",
		Picture:        VALID_BASE64_PNG,
		PreviewPicture: "",
		QRCode:         "",
	},
	Height: 3.0,
	Width:  1.5,
	Depth:  0.5,
	Rows:   10,
	Cols:   10,
}

var SHELF_4 = &shelves.Shelf{
	BasicInfo: common.BasicInfo{
		ID:    SHELF_VALID_UUID_4,
		Label: "A Shelf",
	},
}

var SHELF_5 = &shelves.Shelf{
	BasicInfo: common.BasicInfo{
		ID:    SHELF_VALID_UUID_5,
		Label: "AA Shelf",
	},
}

var SHELF_6 = &shelves.Shelf{
	BasicInfo: common.BasicInfo{
		ID:    SHELF_VALID_UUID_6,
		Label: "BBB",
	},
}

func testShelves() []shelves.Shelf {
	return []shelves.Shelf{*SHELF_1, *SHELF_2, *SHELF_3, *SHELF_4, *SHELF_5, *SHELF_6}
}

func resetShelves() {
	SHELF_1 = &shelves.Shelf{
		BasicInfo: common.BasicInfo{
			ID:             SHELF_VALID_UUID_1,
			Label:          "Test Shelf",
			Description:    "A shelf for testing",
			Picture:        VALID_BASE64_PNG,
			PreviewPicture: "",
			QRCode:         "",
		},
		Height: 2.0,
		Width:  1.5,
		Depth:  0.5,
		Rows:   3,
		Cols:   4,
	}

	SHELF_2 = &shelves.Shelf{
		BasicInfo: common.BasicInfo{
			ID:             SHELF_VALID_UUID_2,
			Label:          "Test Shelf 2",
			Description:    "A shelf for testing",
			Picture:        VALID_BASE64_PNG,
			PreviewPicture: "",
			QRCode:         "testqrcode",
		},
		Height: 2.0,
		Width:  1.5,
		Depth:  0.5,
		Rows:   3,
		Cols:   4,
	}

	SHELF_3 = &shelves.Shelf{
		BasicInfo: common.BasicInfo{
			ID:             SHELF_VALID_UUID_3,
			Label:          "Test   Shelf 3",
			Description:    "A shelf for testing",
			Picture:        VALID_BASE64_PNG,
			PreviewPicture: "",
			QRCode:         "",
		},
		Height: 3.0,
		Width:  1.5,
		Depth:  0.5,
		Rows:   10,
		Cols:   10,
	}

	SHELF_4 = &shelves.Shelf{
		BasicInfo: common.BasicInfo{
			ID:    SHELF_VALID_UUID_4,
			Label: "keyword A",
		},
	}

	SHELF_5 = &shelves.Shelf{
		BasicInfo: common.BasicInfo{
			ID:    SHELF_VALID_UUID_5,
			Label: "keyword Shelf AA",
		},
	}

	SHELF_6 = &shelves.Shelf{

		BasicInfo: common.BasicInfo{
			ID:    SHELF_VALID_UUID_6,
			Label: "BBB",
		},
	}
}

var AREA_1 = &areas.Area{
	BasicInfo: common.BasicInfo{
		ID:             AREA_VALID_UUID_1,
		Label:          "Test Area",
		Description:    "A shelf for testing",
		Picture:        VALID_BASE64_PNG,
		PreviewPicture: "",
		QRCode:         "testqrcode",
	},
}

var AREA_2 = &areas.Area{
	BasicInfo: common.BasicInfo{
		ID:             AREA_VALID_UUID_2,
		Label:          "Test Area 2",
		Description:    "A shelf for testing",
		Picture:        VALID_BASE64_PNG,
		PreviewPicture: "",
		QRCode:         "testqrcode",
	},
}

var AREA_3 = &areas.Area{
	BasicInfo: common.BasicInfo{
		ID:             AREA_VALID_UUID_3,
		Label:          "Test Area 3",
		Description:    "A shelf for testing",
		Picture:        VALID_BASE64_PNG,
		PreviewPicture: "",
		QRCode:         "",
	},
}

var AREA_4 = &areas.Area{
	BasicInfo: common.BasicInfo{
		ID:    AREA_VALID_UUID_4,
		Label: "A Area",
	},
}

var AREA_5 = &areas.Area{
	BasicInfo: common.BasicInfo{
		ID:    AREA_VALID_UUID_5,
		Label: "AA Area",
	},
}

var AREA_6 = &areas.Area{
	BasicInfo: common.BasicInfo{
		ID:    AREA_VALID_UUID_6,
		Label: "BBB",
	},
}

func testAreas() []areas.Area {
	return []areas.Area{*AREA_1, *AREA_2, *AREA_3, *AREA_4, *AREA_5, *AREA_6}
}

func resetAreas() {
	AREA_1 = &areas.Area{
		BasicInfo: common.BasicInfo{
			ID:             AREA_VALID_UUID_1,
			Label:          "Test Area",
			Description:    "An area for testing",
			Picture:        VALID_BASE64_PNG,
			PreviewPicture: "",
			QRCode:         "",
		},
	}

	AREA_2 = &areas.Area{
		BasicInfo: common.BasicInfo{
			ID:             AREA_VALID_UUID_2,
			Label:          "Test Area 2",
			Description:    "An area for testing",
			Picture:        VALID_BASE64_PNG,
			PreviewPicture: "",
			QRCode:         "testqrcode",
		},
	}

	AREA_3 = &areas.Area{
		BasicInfo: common.BasicInfo{
			ID:             AREA_VALID_UUID_3,
			Label:          "Test   Area 3",
			Description:    "An area for testing",
			Picture:        VALID_BASE64_PNG,
			PreviewPicture: "",
			QRCode:         "",
		},
	}

	AREA_4 = &areas.Area{
		BasicInfo: common.BasicInfo{
			ID:    AREA_VALID_UUID_4,
			Label: "keyword A",
		},
	}

	AREA_5 = &areas.Area{
		BasicInfo: common.BasicInfo{
			ID:    AREA_VALID_UUID_5,
			Label: "keyword Area AA",
		},
	}

	AREA_6 = &areas.Area{

		BasicInfo: common.BasicInfo{
			ID:    AREA_VALID_UUID_6,
			Label: "BBB",
		},
	}
}
