package validate_test

import (
	"basement/main/internal/validate"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateID_RequiredEmpty(t *testing.T) {
	v := validate.Validate{}
	rr := httptest.NewRecorder()
	emptyUUID := validate.UUIDField{}
	err := v.ValidateID(rr, emptyUUID, true)
	assert.Error(t, err)
}

func TestValidateID_ValidUUID(t *testing.T) {
	v := validate.Validate{}
	rr := httptest.NewRecorder()
	validUUID := validate.NewUUIDField(uuid.New().String())
	err := v.ValidateID(rr, validUUID, true)
	assert.NoError(t, err)
}

func TestValidateLabel_Empty(t *testing.T) {
	v := validate.Validate{}
	field := validate.NewStringField("")
	v.ValidateLabel(field)
	assert.Equal(t, "Label is required", v.Messages.LabelError)
}

func TestValidateQuantity_InvalidParse(t *testing.T) {
	v := validate.Validate{}
	field := validate.NewIntField("abc")
	v.ValidateQuantity(field)
	assert.Equal(t, "Quantity must be a valid integer", v.Messages.QuantityError)
}

func TestValidateWeight_InvalidParse(t *testing.T) {
	v := validate.Validate{}
	field := validate.NewFloatField("abc")
	v.ValidateWeight(field)
	assert.Equal(t, "Weight must be a valid number", v.Messages.WeightError)
}

func TestValidateDescription_Empty(t *testing.T) {
	v := validate.Validate{}
	field := validate.NewStringField("")
	v.ValidateDescription(field)
	assert.Equal(t, "", v.Messages.DescriptionError)
}

func TestValidateDescription_Invalid(t *testing.T) {
	v := validate.Validate{}
	field := validate.NewStringField("@")
	v.ValidateDescription(field)
	assert.Contains(t, v.Messages.DescriptionError, "invalid characters")
}

func TestValidatePicture_InvalidFormat(t *testing.T) {
	v := validate.Validate{}
	field := validate.NewStringField("text/plain")
	v.ValidatePicture(field)
	assert.Equal(t, "Picture format is not acceptable. Please choose another picture", v.Messages.PictureError)
}

func TestValidatePreviewPicture_InvalidFormat(t *testing.T) {
	v := validate.Validate{}
	field := validate.NewStringField("text/html")
	v.ValidatePreviewPicture(field)
	assert.Equal(t, "Preview picture format is not acceptable. Please choose another picture", v.Messages.PreviewPictureError)
}

func TestValidateItem_Valid(t *testing.T) {
	v := validate.Validate{}
	rr := httptest.NewRecorder()
	item := validate.ItemValidate{
		BasicInfoValidate: validate.BasicInfoValidate{
			ID:             validate.NewUUIDField(uuid.New().String()),
			Label:          validate.NewStringField("Item 1"),
			Description:    validate.NewStringField("Valid desc"),
			Picture:        validate.NewStringField("image/png"),
			PreviewPicture: validate.NewStringField("image/jpeg"),
		},
		Quantity: validate.NewIntField("10"),
		Weight:   validate.NewFloatField("2.5"),
		BoxID:    validate.NewUUIDField(uuid.New().String()),
		ShelfID:  validate.NewUUIDField(uuid.New().String()),
	}
	err := v.ValidateItem(rr, item)
	assert.NoError(t, err)
}

func TestValidateBox_Valid(t *testing.T) {
	v := validate.Validate{}
	rr := httptest.NewRecorder()
	box := validate.BoxValidate{
		BasicInfoValidate: validate.BasicInfoValidate{
			ID:             validate.NewUUIDField(uuid.New().String()),
			Label:          validate.NewStringField("Box A"),
			Description:    validate.NewStringField("desc"),
			Picture:        validate.NewStringField("image/png"),
			PreviewPicture: validate.NewStringField("image/jpeg"),
		},
		ShelfID:    validate.NewUUIDField(uuid.New().String()),
		OuterBoxID: validate.NewUUIDField(uuid.New().String()),
	}
	err := v.ValidateBox(rr, box)
	assert.NoError(t, err)
}

func TestValidateShelf_Valid(t *testing.T) {
	v := validate.Validate{}
	rr := httptest.NewRecorder()
	shelf := validate.ShelfValidate{
		BasicInfoValidate: validate.BasicInfoValidate{
			ID:             validate.NewUUIDField(uuid.New().String()),
			Label:          validate.NewStringField("Shelf X"),
			Description:    validate.NewStringField("desc"),
			Picture:        validate.NewStringField("image/png"),
			PreviewPicture: validate.NewStringField("image/jpeg"),
		},
	}
	err := v.ValidateShelf(rr, shelf)
	assert.NoError(t, err)
}
