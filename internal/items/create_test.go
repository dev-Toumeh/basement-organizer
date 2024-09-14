package items

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid/v5"
)

func TestCheckIDs(t *testing.T) {
	t.Run("Valid IDs provided", func(t *testing.T) {
		validID, _ := uuid.FromString("f47ac10b-58cc-0372-8567-0e02b2c3d479")
		validBoxID, _ := uuid.FromString("e2f234e7-5d59-0985-4f88-5ebb7cc5f31f")
		data := fmt.Sprintf("%s=%s&%s=%s", ID, validID, BOX_ID, validBoxID)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		id, boxID, err := checkIDs(req)

		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
		if id != validID {
			t.Errorf("Expected ID %s, but got: %s", validID, id)
		}
		if boxID != validBoxID {
			t.Errorf("Expected BoxID %s, but got: %s", validBoxID, boxID)
		}
	})

	t.Run("No IDs provided", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)

		id, boxID, err := checkIDs(req)

		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
		if id == uuid.Nil {
			t.Error("Expected a generated ID, but got Nil")
		}
		if boxID != uuid.Nil {
			t.Errorf("Expected BoxID is Nil, but got %s", boxID)
		}
	})

	t.Run("Invalid IDs provided", func(t *testing.T) {
		data := fmt.Sprintf("%s=invalid&%s=alsoinvalid", ID, BOX_ID)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		_, _, err := checkIDs(req)

		if err == nil {
			t.Error("Expected an error for invalid IDs, but got none")
		}
	})

	t.Run("One ID missing", func(t *testing.T) {
		validID, _ := uuid.FromString("f47ac10b-58cc-0372-8567-0e02b2c3d479")
		data := fmt.Sprintf("%s=%s", ID, validID)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		id, boxId, err := checkIDs(req)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if id == uuid.Nil {
			t.Errorf("Expected a valid ID, but got uuid.Nil")
		}

		if boxId != uuid.Nil {
			t.Errorf("Expected boxId to be uuid.Nil, but got %v", boxId)
		}
	})

}
