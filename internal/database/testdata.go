package database

import (
	"basement/main/internal/items"
	"basement/main/internal/shelves"

	"github.com/gofrs/uuid/v5"
)

var VALID_UUID_1 uuid.UUID = uuid.Must(uuid.FromString("623e4567-e89b-12d3-a456-426614174000"))
var VALID_UUID_2 uuid.UUID = uuid.Must(uuid.FromString("323e4567-e89b-12d3-a456-426614174000"))
var VALID_UUID_3 uuid.UUID = uuid.Must(uuid.FromString("423e4567-e89b-12d3-a456-426614174000"))
var VALID_UUID_4 uuid.UUID = uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174111"))

var SHELF_VALID_UUID_1 uuid.UUID = uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174000"))
var SHELF_VALID_UUID_2 uuid.UUID = uuid.Must(uuid.FromString("223e4567-e89b-12d3-a456-426614174000"))
var SHELF_VALID_UUID_3 uuid.UUID = uuid.Must(uuid.FromString("323e4567-e89b-12d3-a456-426614174000"))
var ITEM_VALID_UUID uuid.UUID = uuid.Must(uuid.FromString("133e4567-e89b-12d3-a456-426614174000"))
var VALID_UUID_NOT_EXISTING uuid.UUID = uuid.Must(uuid.FromString("033e4567-e89b-12d3-a456-426614174000"))

const VALID_BASE64_PNG = "iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAIAAAACUFjqAAAAtUlEQVR4nGJp2XGEAQb+/P49J7cgY8ZUuAgTnDUjI7vUQf3m5e3/zxyakZENFW3ZcURGQf/r52cQBGfLKhq0bD/MqKBu+ufnL4jSm5e3QxjmtuEfPnyCGn7z8na4BAMDg7quZ2mia2thMAMDA0j31TMb4XJr5s2BMKr71zIwMLAwYIDq/rWMMDaLobs7mjTEWJC6CeuYjL08+o/eU9f1RDPgsbpTxvQpjMjBAvEucrAAAgAA//+Elk5AOfCu8QAAAABJRU5ErkJggg=="
const VALID_BASE64_PREVIEW_PNG = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAAEElEQVR4nGJaOLEJEAAA//8DkwG35JmAnAAAAABJRU5ErkJggg=="
const INVALID_BASE64_PNG = "invalid base 64"

func resetTestBoxes() {
	BOX_1 = &items.Box{
		BasicInfo: items.BasicInfo{
			ID:          VALID_UUID_1,
			Label:       "box 1",
			Description: "This is the sixth box",
			Picture:     VALID_BASE64_PNG,
			QRcode:      "uvwxyzabcdefg",
		},
	}
	BOX_2 = &items.Box{
		BasicInfo: items.BasicInfo{
			ID:          VALID_UUID_2,
			Label:       "box 3",
			Description: "This is the third box",
			Picture:     VALID_BASE64_PNG,
			QRcode:      "abababababcd",
		},
	}
	BOX_3 = &items.Box{
		BasicInfo: items.BasicInfo{
			ID:          VALID_UUID_3,
			Label:       "box 4",
			Description: "This is the fourth box",
			Picture:     VALID_BASE64_PNG,
			QRcode:      "efghefghefgh",
		},
	}
	BOX_4 = &items.Box{
		BasicInfo: items.BasicInfo{
			ID:          VALID_UUID_4,
			Label:       "box 5",
			Description: "This is the fifth box",
			Picture:     VALID_BASE64_PNG,
			QRcode:      "ijklmnopqrst",
		},
	}
}

func testBoxes() []*items.Box {
	return []*items.Box{BOX_1, BOX_2, BOX_3, BOX_4}
}

var BOX_1 = &items.Box{
	BasicInfo: items.BasicInfo{
		ID:          VALID_UUID_1,
		Label:       "box 1",
		Description: "This is the sixth box",
		Picture:     VALID_BASE64_PNG,
		QRcode:      "uvwxyzabcdefg",
	},
}

var BOX_2 = &items.Box{
	BasicInfo: items.BasicInfo{
		ID:          VALID_UUID_2,
		Label:       "box 3",
		Description: "This is the third box",
		Picture:     VALID_BASE64_PNG,
		QRcode:      "abababababcd",
	},
}

var BOX_3 = &items.Box{
	BasicInfo: items.BasicInfo{
		ID:          VALID_UUID_3,
		Label:       "box 4",
		Description: "This is the fourth box",
		Picture:     VALID_BASE64_PNG,
		QRcode:      "efghefghefgh",
	},
}

var BOX_4 = &items.Box{
	BasicInfo: items.BasicInfo{
		ID:          VALID_UUID_4,
		Label:       "box 5",
		Description: "This is the fifth box",
		Picture:     VALID_BASE64_PNG,
		QRcode:      "ijklmnopqrst",
	},
}

var ITEM_1 = &items.Item{
	BasicInfo: items.BasicInfo{
		ID:          uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174000")),
		Label:       "Item 1",
		Description: "Description for item 1",
		Picture:     "base64encodedstring1",
	},
	Quantity: 10,
	Weight:   "5.5",
	QRCode:   "QRcode1",
}

var ITEM_2 = &items.Item{
	BasicInfo: items.BasicInfo{
		ID:          uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174001")),
		Label:       "Item 2",
		Description: "Description for item 2",
		Picture:     "base64encodedstring2",
	},
	Quantity: 20,
	Weight:   "10.0",
	QRCode:   "QRcode2",
}

var ITEM_3 = &items.Item{
	BasicInfo: items.BasicInfo{
		ID:          uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174002")),
		Label:       "Item 3",
		Description: "Description for item 3",
		Picture:     "base64encodedstring3",
	},
	Quantity: 15,
	Weight:   "7.25",
	QRCode:   "QRcode3",
}

func resetTestItems() {
	ITEM_1 = &items.Item{
		BasicInfo: items.BasicInfo{
			ID:          uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174000")),
			Label:       "Item 1",
			Description: "Description for item 1",
			Picture:     "base64encodedstring1",
		},
		Quantity: 10,
		Weight:   "5.5",
		QRCode:   "QRcode1",
	}

	ITEM_2 = &items.Item{
		BasicInfo: items.BasicInfo{
			ID:          uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174001")),
			Label:       "Item 2",
			Description: "Description for item 2",
			Picture:     "base64encodedstring2",
		},
		Quantity: 20,
		Weight:   "10.0",
		QRCode:   "QRcode2",
	}

	ITEM_3 = &items.Item{
		BasicInfo: items.BasicInfo{
			ID:          uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174002")),
			Label:       "Item 3",
			Description: "Description for item 3",
			Picture:     "base64encodedstring3",
		},
		Quantity: 15,
		Weight:   "7.25",
		QRCode:   "QRcode3",
	}
}

func testItems() []items.Item {
	return []items.Item{*ITEM_1, *ITEM_2, *ITEM_3}
}

var SHELF_1 = &shelves.Shelf{
	ID:             SHELF_VALID_UUID_1,
	Label:          "Test Shelf",
	Description:    "A shelf for testing",
	Picture:        VALID_BASE64_PNG,
	PreviewPicture: "",
	QRcode:         "testqrcode",
	Height:         2.0,
	Width:          1.5,
	Depth:          0.5,
	Rows:           3,
	Cols:           4,
}

var SHELF_2 = &shelves.Shelf{
	ID:             SHELF_VALID_UUID_2,
	Label:          "Test Shelf 2",
	Description:    "A shelf for testing",
	Picture:        VALID_BASE64_PNG,
	PreviewPicture: "",
	QRcode:         "testqrcode",
	Height:         2.0,
	Width:          1.5,
	Depth:          0.5,
	Rows:           3,
	Cols:           4,
}

var SHELF_3 = &shelves.Shelf{
	ID:             SHELF_VALID_UUID_3,
	Label:          "Test Shelf 3",
	Description:    "A shelf for testing",
	Picture:        VALID_BASE64_PNG,
	PreviewPicture: "",
	QRcode:         "",
	Height:         3.0,
	Width:          1.5,
	Depth:          0.5,
	Rows:           10,
	Cols:           10,
}

var SHELF_4 = &shelves.Shelf{
	ID:    uuid.Must(uuid.NewV4()),
	Label: "A Shelf",
}

var SHELF_5 = &shelves.Shelf{
	ID:    uuid.Must(uuid.NewV4()),
	Label: "AA Shelf",
}

var SHELF_6 = &shelves.Shelf{
	ID:    uuid.Must(uuid.NewV4()),
	Label: "BBB",
}

func testShelves() []shelves.Shelf {
	return []shelves.Shelf{*SHELF_1, *SHELF_2, *SHELF_3, *SHELF_4, *SHELF_5, *SHELF_6}
}

func resetShelves() {
	SHELF_1 = &shelves.Shelf{
		ID:             SHELF_VALID_UUID_1,
		Label:          "Test Shelf",
		Description:    "A shelf for testing",
		Picture:        VALID_BASE64_PNG,
		PreviewPicture: "",
		QRcode:         "",
		Height:         2.0,
		Width:          1.5,
		Depth:          0.5,
		Rows:           3,
		Cols:           4,
	}

	SHELF_2 = &shelves.Shelf{
		ID:             SHELF_VALID_UUID_2,
		Label:          "Test Shelf 2",
		Description:    "A shelf for testing",
		Picture:        VALID_BASE64_PNG,
		PreviewPicture: "",
		QRcode:         "testqrcode",
		Height:         2.0,
		Width:          1.5,
		Depth:          0.5,
		Rows:           3,
		Cols:           4,
	}
	SHELF_3 = &shelves.Shelf{
		ID:             SHELF_VALID_UUID_3,
		Label:          "Test   Shelf 3",
		Description:    "A shelf for testing",
		Picture:        VALID_BASE64_PNG,
		PreviewPicture: "",
		QRcode:         "",
		Height:         3.0,
		Width:          1.5,
		Depth:          0.5,
		Rows:           10,
		Cols:           10,
	}

	SHELF_4 = &shelves.Shelf{
		ID:    uuid.Must(uuid.NewV4()),
		Label: "keyword A",
	}

	SHELF_5 = &shelves.Shelf{
		ID:    uuid.Must(uuid.NewV4()),
		Label: "keyword Shelf AA",
	}

	SHELF_6 = &shelves.Shelf{
		ID:    uuid.Must(uuid.NewV4()),
		Label: "BBB",
	}
}
