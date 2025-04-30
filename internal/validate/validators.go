package validate

import (
	"basement/main/internal/logg"
	"fmt"
	"math"
	"net/http"
)

func (v *Validate) ValidateID(w http.ResponseWriter, field UUIDField, required bool) error {
	if field.IsEmpty() || field.IsNil() {
		if required {
			logg.Debugf("Validation failed: ID is nil or empty")
			return logg.NewError("Error happened, please come back later")
		}
		return nil
	}

	if err := field.IsValid(); err != nil {
		logg.Debugf("Validation failed for ID: \n%v", err)
		return logg.NewError("Error happened, please come back later")
	}

	return nil
}

func (v *Validate) ValidateLabel(s StringField) {
	if s.IsEmpty() {
		logg.Debugf("Validation failed: Label is empty")
		v.Messages.LabelError = "Label is required"
	} else {
		if err := s.MinLength(); err != nil {
			logg.Debugf("Validation failed for Label (MinLength): %v", err)
			v.Messages.LabelError = "Label must be at least 1 character long"
		}
		if err := s.MaxLength(); err != nil {
			logg.Debugf("Validation failed for Label (MaxLength): %v", err)
			v.Messages.LabelError = "Label must be at most 255 characters long"
		}
		if err := s.MatchesRegex(); err != nil {
			logg.Debugf("Validation failed for Label (Regex): %v", err)
			v.Messages.LabelError = "Label contains invalid characters (only letters, numbers, spaces, _, ., and - are allowed)"
		}
	}
}

func (v *Validate) ValidateDescription(s StringField) {
	if !s.IsEmpty() {
		if err := s.MinLength(); err != nil {
			logg.Debugf("Validation failed for Description (MinLength): %v", err)
			v.Messages.DescriptionError = "Description must be at least 1 character long"
		}
		if err := s.MaxLengthCustom(1000); err != nil {
			logg.Debugf("Validation failed for Description (MaxLength): %v", err)
			v.Messages.DescriptionError = "Description exceeds maximum length of 1000 characters"
		}
		if err := s.MatchesRegex(); err != nil {
			logg.Debugf("Validation failed for Description (Regex): %v", err)
			v.Messages.DescriptionError = "Description contains invalid characters (only letters, numbers, spaces, _, ., and - are allowed)"
		}
	}
}

func (v *Validate) ValidatePicture(s StringField) {
	if !s.IsEmpty() {
		if err := s.MaxLengthCustom(math.MaxInt); err != nil {
			logg.Debugf("Validation failed for Picture: %v", err)
			v.Messages.PictureError = "Picture path exceeds maximum length"
		}
		if err := s.ValidatePictureFormat(); err != nil {
			logg.Debugf("Validation failed for Picture: %v", err)
			v.Messages.PictureError = "Picture format is not acceptable. Please choose another picture"
		}
	}
}

func (v *Validate) ValidatePreviewPicture(s StringField) {
	if !s.IsEmpty() {
		if err := s.MaxLengthCustom(math.MaxInt); err != nil {
			logg.Debugf("Validation failed for Preview Picture: %v", err)
			v.Messages.PreviewPictureError = "Preview picture path exceeds maximum length"
		}
		if err := s.ValidatePictureFormat(); err != nil {
			logg.Debugf("Validation failed for Preview Picture: %v", err)
			v.Messages.PreviewPictureError = "Preview picture format is not acceptable. Please choose another picture"
		}
	}
}

func (v *Validate) ValidateQuantity(q IntField) {
	if q.err != nil {
		logg.Debugf("Validation failed for Quantity (Parse Error): %v", q.err)
		v.Messages.QuantityError = "Quantity must be a valid integer"
		return
	}

	if !q.IsEmpty() {
		if err := q.IsPositive(); err != nil {
			logg.Debugf("Validation failed for Quantity (IsPositive): %v", err)
			v.Messages.QuantityError = "Quantity must be a positive number"
			return
		}
		if err := q.MinValue(); err != nil {
			logg.Debugf("Validation failed for Quantity (MinValue): %v", err)
			v.Messages.QuantityError = fmt.Sprintf("Quantity must be at least %d", q.DefaultMinValue)
			return
		}
		if err := q.MaxValue(); err != nil {
			logg.Debugf("Validation failed for Quantity (MaxValue): %v", err)
			v.Messages.QuantityError = fmt.Sprintf("Quantity must not exceed %d", q.DefaultMaxValue)
			return
		}
	}
}

func (v *Validate) ValidateWeight(f FloatField) {
	if f.Err != nil {
		logg.Debugf("Validation failed for Weight (Parse Error): %v", f.Err)
		v.Messages.WeightError = "Weight must be a valid number"
		return
	}
	if !f.IsEmpty() {
		if err := f.IsPositive(); err != nil {
			logg.Debugf("Validation failed for Weight (IsPositive): %v", err)
			v.Messages.WeightError = "Weight must be a positive number"
			return
		}
		if err := f.MinValue(); err != nil {
			logg.Debugf("Validation failed for Weight (MinValue): %v", err)
			v.Messages.WeightError = fmt.Sprintf("Weight must be at least %.2f", f.DefaultMinValue)
			return
		}
		if err := f.MaxValue(); err != nil {
			logg.Debugf("Validation failed for Weight (MaxValue): %v", err)
			v.Messages.WeightError = fmt.Sprintf("Weight must not exceed %.2f", f.DefaultMaxValue)
			return
		}
	}
}

func (v *Validate) ValidateItem(w http.ResponseWriter, item ItemValidate) (err error) {
	if err := v.ValidateID(w, item.ID, true); err != nil {
		return err
	}
	v.ValidateLabel(item.Label)
	v.ValidateDescription(item.Description)
	v.ValidatePicture(item.Picture)
	v.ValidatePreviewPicture(item.PreviewPicture)

	v.ValidateQuantity(item.Quantity)
	v.ValidateWeight(item.Weight)

	if err := v.ValidateID(w, item.BoxID, false); err != nil {
		return err
	}
	if err := v.ValidateID(w, item.ShelfID, false); err != nil {
		return err
	}
	if err := v.ValidateID(w, item.ShelfID, false); err != nil {
		return err
	}
	return nil
}

func (v *Validate) ValidateBox(w http.ResponseWriter, box BoxValidate) (err error) {
	if err := v.ValidateID(w, box.ID, true); err != nil {
		return err
	}
	v.ValidateLabel(box.Label)
	v.ValidateDescription(box.Description)
	v.ValidatePicture(box.Picture)
	v.ValidatePreviewPicture(box.PreviewPicture)

	if err := v.ValidateID(w, box.OuterBoxID, false); err != nil {
		return err
	}
	if err := v.ValidateID(w, box.ShelfID, false); err != nil {
		return err
	}
	if err := v.ValidateID(w, box.ShelfID, false); err != nil {
		return err
	}
	return nil
}

func (v *Validate) ValidateShelf(w http.ResponseWriter, shelf ShelfValidate) (err error) {
	if err := v.ValidateID(w, shelf.ID, true); err != nil {
		return err
	}
	v.ValidateLabel(shelf.Label)
	v.ValidateDescription(shelf.Description)
	v.ValidatePicture(shelf.Picture)
	v.ValidatePreviewPicture(shelf.PreviewPicture)
	return nil
}
