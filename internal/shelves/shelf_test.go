package shelves

import (
	"basement/main/internal/common"
	"errors"
	"net/http"
	"net/http/httptest"

	// "net/url"
	// "strconv"
	// "strings"
	"testing"

	"github.com/gofrs/uuid/v5"
	// "github.com/stretchr/testify/assert"
)

var shelf1 *Shelf = &Shelf{
	BasicInfo: common.BasicInfo{
		ID:             uuid.Must(uuid.FromString("111e4567-e89b-12d3-a456-426614174000")),
		Label:          "Storage Shelf 1",
		Description:    "This is the first dummy shelf",
		Picture:        "base64PictureData1",
		PreviewPicture: "base64PreviewPictureData1",
		QRCode:         "QR1234ABC",
	},
	Height: 250.0,
	Width:  120.0,
	Depth:  60.0,
	Rows:   4,
	Cols:   3,
	AreaId: uuid.Must(uuid.FromString("222e4567-e89b-12d3-a456-426614174001")),
}

var shelf2 *Shelf = &Shelf{
	BasicInfo: common.BasicInfo{
		ID:             uuid.Must(uuid.FromString("333e4567-e89b-12d3-a456-426614174002")),
		Label:          "Storage Shelf 2",
		Description:    "This is the second dummy shelf",
		Picture:        "base64PictureData2",
		PreviewPicture: "base64PreviewPictureData2",
		QRCode:         "QR5678XYZ",
	},
	Height: 300.0,
	Width:  150.0,
	Depth:  70.0,
	Rows:   5,
	Cols:   4,
	AreaId: uuid.Must(uuid.FromString("444e4567-e89b-12d3-a456-426614174003")),
}

type ShelfDatabaseError struct{}

func (db *ShelfDatabaseError) CreateShelf(shelf *Shelf) error {
	return errors.New("unable to create shelf")
}

func (db *ShelfDatabaseError) Shelf(id uuid.UUID) (*Shelf, error) {
	return nil, errors.New("shelf not found")
}

func (db *ShelfDatabaseError) UpdateShelf(shelf *Shelf) error {
	return errors.New("unable to update shelf")
}

func (db *ShelfDatabaseError) DeleteShelf(id uuid.UUID) error {
	return errors.New("unable to delete shelf")
}

func (db *ShelfDatabaseError) SearchShelves(page int, rows int, query string) ([]*common.ListRow, error) {
	return nil, errors.New("unable to delete shelf")
}

// ShelfDatabaseSuccess implements ShelfDB interface without errors for success testing.
type ShelfDatabaseSuccess struct{}

func (db *ShelfDatabaseSuccess) CreateShelf(shelf *Shelf) error {
	return nil
}

func (db *ShelfDatabaseSuccess) Shelf(id uuid.UUID) (*Shelf, error) {
	return shelf1, nil
}

func (db *ShelfDatabaseSuccess) UpdateShelf(shelf *Shelf) error {
	return nil
}

func (db *ShelfDatabaseSuccess) DeleteShelf(id uuid.UUID) error {
	return nil
}

func (db *ShelfDatabaseSuccess) SearchShelves(page int, rows int, query string) ([]*common.ListRow, error) {
	return nil, nil
}

// TestTypeFromRequest tests the typeFromRequest function with different scenarios.
func TestTypeFromRequest(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		want     string
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "Valid type 'add'",
			url:     "/?type=add",
			want:    "Add",
			wantErr: false,
		},
		{
			name:    "Valid type 'move'",
			url:     "/?type=move",
			want:    "Move",
			wantErr: false,
		},
		{
			name:    "Valid type 'search'",
			url:     "/?type=search",
			want:    "Search",
			wantErr: false,
		},
		{
			name:     "Invalid type 'delete'",
			url:      "/?type=delete",
			want:     "",
			wantErr:  true,
			errorMsg: "unexpected type: delete, while preparing the search Template",
		},
		{
			name:    "Empty type",
			url:     "/?type=",
			want:    "",
			wantErr: false,
		},
		{
			name:    "No type parameter",
			url:     "/",
			want:    "",
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			got, err := typeFromRequest(req)
			if (err != nil) != tc.wantErr {
				t.Errorf("typeFromRequest() error = %v, wantErr %v", err, tc.wantErr)
			}
			if err != nil && tc.errorMsg != "" && err.Error() != tc.errorMsg {
				t.Errorf("typeFromRequest() error = %v, wantErr message %v", err, tc.errorMsg)
			}
			if got != tc.want {
				t.Errorf("typeFromRequest() = %v, want %v", got, tc.want)
			}
		})
	}
}

// func TestShelvesHandler(t *testing.T) {
// 	dbErr := &ShelfDatabaseError{} // Ensure this implements ShelfDB
// 	mux := http.NewServeMux()
// 	mux.Handle("/api/v1/create/shelf", ShelfHandler(dbErr))
//
// 	type handlerInput struct {
// 		R *http.Request
// 		W httptest.ResponseRecorder
// 	}
//
// 	testCases := []struct {
// 		name               string
// 		input              handlerInput
// 		expectedStatusCode int
// 	}{
// 		{
// 			name: "Test Get Request",
// 			input: handlerInput{
// 				R: httptest.NewRequest(http.MethodGet, "/api/v1/create/shelf", nil),
// 				W: *httptest.NewRecorder(),
// 			},
// 			expectedStatusCode: http.StatusPermanentRedirect,
// 		},
// 		{
// 			name: "Test Valid POST Request",
// 			input: func() handlerInput {
// 				formData := url.Values{}
// 				formData.Set("id", shelf1.ID.String())
// 				formData.Set("area_id", shelf1.AreaId.String())
// 				formData.Set("label", shelf1.Label)
// 				formData.Set("description", shelf1.Description)
// 				formData.Set("qrcode", shelf1.QRCode)
// 				formData.Set("height", strconv.FormatFloat(float64(shelf1.Height), 'f', 2, 32))
// 				formData.Set("width", strconv.FormatFloat(float64(shelf1.Width), 'f', 2, 32))
// 				formData.Set("depth", strconv.FormatFloat(float64(shelf1.Depth), 'f', 2, 32))
// 				formData.Set("rows", strconv.Itoa(shelf1.Rows))
// 				formData.Set("cols", strconv.Itoa(shelf1.Cols))
//
// 				encodedFormData := formData.Encode()
//
// 				req := httptest.NewRequest(http.MethodPost, "/api/v1/create/shelf", strings.NewReader(encodedFormData))
// 				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//
// 				rec := httptest.NewRecorder()
//
// 				return handlerInput{
// 					R: req,
// 					W: *rec, // Use the value type
// 				}
// 			}(),
// 			expectedStatusCode: http.StatusOK,
// 		},
// 	}
//
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Pass &tc.input.W to ServeHTTP
// 			mux.ServeHTTP(&tc.input.W, tc.input.R)
// 			result := tc.input.W.Result()
// 			assert.Equal(t, tc.expectedStatusCode, result.StatusCode, "URL: "+tc.input.R.URL.String())
// 		})
// 	}
// }
