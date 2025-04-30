package validate

import "errors"

// Map returns the validation error messages as a map[string]any,
// so it can be easily passed to templates or JSON responses.
func (v *ValidateMessages) Map() map[string]any {
	return map[string]any{
		"LabelError":          v.LabelError,
		"DescriptionError":    v.DescriptionError,
		"PictureError":        v.PictureError,
		"PreviewPictureError": v.PreviewPictureError,
		"QuantityError":       v.QuantityError,
		"WeightError":         v.WeightError,
	}
}

func (b BasicInfoValidate) Map() map[string]any {
	return map[string]any{
		"ID":             b.ID.UUID().String(),
		"Label":          b.Label.String(),
		"Description":    b.Description.String(),
		"Picture":        b.Picture.String(),
		"PreviewPicture": b.PreviewPicture.String(),
		"QRCode":         b.QRCode.String(),
	}
}

func (i ItemValidate) Map() map[string]any {
	m := i.BasicInfoValidate.Map()
	m["Quantity"] = i.Quantity.Int()
	m["Weight"] = i.Weight.Float64()
	m["BoxID"] = i.BoxID.UUID()
	m["ShelfID"] = i.ShelfID.UUID()
	m["AreaID"] = i.AreaID.UUID()
	return m
}

func (s ShelfValidate) Map() map[string]any {
	m := s.BasicInfoValidate.Map()
	m["Height"] = s.Height.Float64()
	m["Width"] = s.Width.Float64()
	m["Depth"] = s.Depth.Float64()
	m["Rows"] = s.Rows.Int()
	m["Cols"] = s.Cols.Int()
	m["AreaID"] = s.AreaID.UUID()
	return m
}

func (b BoxValidate) Map() map[string]any {
	m := b.BasicInfoValidate.Map()
	m["ShelfID"] = b.ShelfID.UUID()
	m["OuterBoxID"] = b.OuterBoxID.UUID()
	m["AreaID"] = b.AreaID.UUID()
	return m
}

func (a AreaValidate) Map() map[string]any {
	return a.BasicInfoValidate.Map()
}

// ItemFormData returns a map combining the ItemValidate fields and the validation error messages.
// It is used to render the form with the current input values and validation feedback.
func (v *Validate) ItemFormData() map[string]any {
	data := v.Item.Map()
	for key, value := range v.Messages.Map() {
		data[key] = value
	}
	return data
}

// BoxFormData returns a map combining the BoxValidate fields and the validation error messages.
// It is used to render the form with the current input values and validation feedback.
func (v *Validate) BoxFormData() map[string]any {
	data := v.Box.Map()
	for key, value := range v.Messages.Map() {
		data[key] = value
	}
	return data
}

// ShelfFormData returns a map combining the ShelfValidate fields and the validation error messages.
// It is used to render the form with the current input values and validation feedback.
func (v *Validate) ShelfFormData() map[string]any {
	data := v.Shelf.Map()
	for key, value := range v.Messages.Map() {
		data[key] = value
	}
	return data
}

// AreaFormData returns a map combining the AreaValidate fields and the validation error messages.
// It is used to render the form with the current input values and validation feedback.
func (v *Validate) AreaFormData() map[string]any {
	data := v.Area.Map()
	for key, value := range v.Messages.Map() {
		data[key] = value
	}
	return data
}

// error logic
var ValidationError = errors.New("ValidationError")

func (v *Validate) HasValidateErrors() bool {
	for _, val := range v.Messages.Map() {
		if str, ok := val.(string); ok && str != "" {
			return true
		}
	}
	return false
}

func (v Validate) Err() error {
	return ValidationError
}
