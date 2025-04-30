package validate

type Validate struct {
	Login           LoginValidate
	Item            ItemValidate
	Box             BoxValidate
	Shelf           ShelfValidate
	Area            AreaValidate
	Messages        ValidateMessages
	ValidationError error
}

type LoginValidate struct {
	ID           UUIDField
	Username     StringField
	PasswordHash StringField
	email        StringField
}

type BasicInfoValidate struct {
	ID             UUIDField
	Label          StringField
	Description    StringField
	Picture        StringField
	PreviewPicture StringField
	QRCode         StringField
}

type ItemValidate struct {
	BasicInfoValidate
	Quantity IntField
	Weight   FloatField
	BoxID    UUIDField
	ShelfID  UUIDField
	AreaID   UUIDField
}

type BoxValidate struct {
	BasicInfoValidate
	ShelfID    UUIDField
	OuterBoxID UUIDField
	AreaID     UUIDField
}

type ShelfValidate struct {
	BasicInfoValidate
	Height FloatField
	Width  FloatField
	Depth  FloatField
	Rows   IntField
	Cols   IntField
	AreaID UUIDField
}

type AreaValidate struct {
	BasicInfoValidate
}

type ValidateMessages struct {
	LabelError          string
	DescriptionError    string
	PictureError        string
	PreviewPictureError string
	QuantityError       string
	WeightError         string
}
