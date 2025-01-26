package common

import (
	"basement/main/internal/auth"
	"basement/main/internal/logg"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

// Data structure definition
type Data struct {
	ListTemplate,
	Title string
	Authenticated  bool
	User           string
	Debug          bool
	NotFound       bool
	EnvDevelopment bool
	RequestOrigin  string
	TypeMap        map[string]any
}

func InitData(r *http.Request) (data Data) {
	data = Data{
		TypeMap: make(map[string]any),
	}
	data.SetPageNumber(ParsePageNumber(r))
	data.SetSearchInputValue(SearchString(r))
	data.SetLimit(ParseLimit(r))

	user, _ := auth.UserSessionData(r)
	authenticated, _ := auth.Authenticated(r)
	data.SetUser(user)
	data.SetAuthenticated(authenticated)

	return data
}

// init main Data for non Paginated Templates
func InitData2(r *http.Request) (data *Data) {
	data = &Data{
		TypeMap: make(map[string]any),
	}
	data.SetEdit(CheckEditMode(r))
	user, _ := auth.UserSessionData(r)
	authenticated, _ := auth.Authenticated(r)
	data.SetUser(user)
	data.SetAuthenticated(authenticated)

	return data
}

func (data *Data) SetTitle(value string) {
	data.TypeMap["Title"] = value
}

func (data *Data) GetTitle() (string, error) {
	if val, exists := data.TypeMap["Title"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'Title' does not exist")
}

func (data *Data) SetAuthenticated(value bool) {
	data.Authenticated = value
	data.TypeMap["Authenticated"] = value
}

func (data *Data) GetAuthenticated() (bool, error) {
	if val, exists := data.TypeMap["Authenticated"]; exists {
		return val.(bool), nil
	}
	return false, errors.New("key 'Authenticated' does not exist")
}

func (data *Data) SetUser(value string) {
	data.TypeMap["User"] = value
}

func (data *Data) GetUser() string {
	if val, exists := data.TypeMap["User"]; exists {
		return val.(string)
	}
	return ""
}

func (data *Data) SetDebug(value bool) {
	data.TypeMap["Debug"] = value
}

func (data *Data) GetDebug() (bool, error) {
	if val, exists := data.TypeMap["Debug"]; exists {
		return val.(bool), nil
	}
	return false, errors.New("key 'Debug' does not exist")
}

func (data *Data) SetNotFound(value bool) {
	data.TypeMap["NotFound"] = value
}

func (data *Data) GetNotFound() (bool, error) {
	if val, exists := data.TypeMap["NotFound"]; exists {
		return val.(bool), nil
	}
	return false, errors.New("key 'NotFound' does not exist")
}

func (data *Data) SetEnvDevelopment(value bool) {
	data.EnvDevelopment = value
	data.TypeMap["EnvDevelopment"] = value
}

func (data *Data) GetEnvDevelopment() (bool, error) {
	if val, exists := data.TypeMap["EnvDevelopment"]; exists {
		return val.(bool), nil
	}
	return false, errors.New("key 'EnvDevelopment' does not exist")
}

func (data *Data) SetFormID(value string) {
	data.TypeMap["FormID"] = value
}

func (data *Data) GetFormID() (string, error) {
	if val, exists := data.TypeMap["FormID"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'FormID' does not exist")
}

func (data *Data) SetFormHXGet(value string) {
	data.TypeMap["FormHXGet"] = value
}

func (data *Data) GetFormHXGet() (string, error) {
	if val, exists := data.TypeMap["FormHXGet"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'FormHXGet' does not exist")
}

func (data *Data) SetFormHXPost(value string) {
	data.TypeMap["FormHXPost"] = value
}

func (data *Data) GetFormHXPost() (string, error) {
	if val, exists := data.TypeMap["FormHXPost"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'FormHXPost' does not exist")
}

func (data *Data) SetFormHXTarget(value string) {
	data.TypeMap["FormHXTarget"] = value
}

func (data *Data) GetFormHXTarget() (string, error) {
	if val, exists := data.TypeMap["FormHXTarget"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'FormHXTarget' does not exist")
}

func (data *Data) SetRowHXGet(value string) {
	data.TypeMap["RowHXGet"] = value
}

func (data *Data) GetRowHXGet() (string, error) {
	if val, exists := data.TypeMap["RowHXGet"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'RowHXGet' does not exist")
}

func (data *Data) SetSearchInput(value bool) {
	data.TypeMap["SearchInput"] = value
}

func (data *Data) GetSearchInput() (bool, error) {
	if val, exists := data.TypeMap["SearchInput"]; exists {
		return val.(bool), nil
	}
	return false, errors.New("key 'SearchInput' does not exist")
}

func (data *Data) SetSearchInputLabel(value string) {
	data.TypeMap["SearchInputLabel"] = value
}

func (data *Data) GetSearchInputLabel() (string, error) {
	if val, exists := data.TypeMap["SearchInputLabel"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'SearchInputLabel' does not exist")
}

func (data *Data) SetSearchInputValue(value string) {
	data.TypeMap["SearchInputValue"] = value
}

func (data *Data) GetSearchInputValue() string {
	if val, exists := data.TypeMap["SearchInputValue"]; exists {
		return val.(string)
	}
	return ""
}

func (data *Data) SetPagination(value bool) {
	data.TypeMap["Pagination"] = value
}

func (data *Data) GetPagination() (bool, error) {
	if val, exists := data.TypeMap["Pagination"]; exists {
		return val.(bool), nil
	}
	return false, errors.New("key 'Pagination' does not exist")
}

func (data *Data) SetCurrentPageNumber(value int) {
	data.TypeMap["CurrentPageNumber"] = value
}

func (data *Data) GetCurrentPageNumber() (int, error) {
	if val, exists := data.TypeMap["CurrentPageNumber"]; exists {
		return val.(int), nil
	}
	return 0, errors.New("key 'CurrentPageNumber' does not exist")
}

func (data *Data) SetLimit(value int) {
	data.TypeMap["Limit"] = value // Store as int directly
}

func (data *Data) GetLimit() int {
	if val, exists := data.TypeMap["Limit"]; exists {
		return val.(int)
	}
	return 10
}

func (data *Data) GetLimitStr() (string, error) {
	if val, exists := data.TypeMap["Limit"]; exists {
		return fmt.Sprint(val), nil
	}
	return "", errors.New("key 'Limit' does not exist")
}

func (data *Data) SetShowLimit(value bool) {
	data.TypeMap["ShowLimit"] = value
}

func (data *Data) GetShowLimit() (bool, error) {
	if val, exists := data.TypeMap["ShowLimit"]; exists {
		return val.(bool), nil
	}
	return false, errors.New("key 'ShowLimit' does not exist")
}

func (data *Data) SetPaginationButtons(value []PaginationButton) {
	data.TypeMap["PaginationButtons"] = value
}

func (data *Data) GetPaginationButtons() ([]PaginationButton, error) {
	if val, exists := data.TypeMap["PaginationButtons"]; exists {
		return val.([]PaginationButton), nil
	}
	return nil, errors.New("key 'PaginationButtons' does not exist")
}

func (data *Data) SetMoveButtonHXTarget(value string) {
	data.TypeMap["MoveButtonHXTarget"] = value
}

func (data *Data) GetMoveButtonHXTarget() (string, error) {
	if val, exists := data.TypeMap["MoveButtonHXTarget"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'MoveButtonHXTarget' does not exist")
}

func (data *Data) SetRows(value []ListRow) {
	data.TypeMap["Rows"] = value
}

func (data *Data) GetRows() ([]ListRow, error) {
	if val, exists := data.TypeMap["Rows"]; exists {
		return val.([]ListRow), nil
	}
	return nil, errors.New("key 'Rows' does not exist")
}

func (data *Data) SetRowAction(value bool) {
	data.TypeMap["RowAction"] = value
}

func (data *Data) GetRowAction() (bool, error) {
	if val, exists := data.TypeMap["RowAction"]; exists {
		return val.(bool), nil
	}
	return false, errors.New("key 'RowAction' does not exist")
}

func (data *Data) SetRowActionType(value string) {
	data.TypeMap["RowActionType"] = value
}

func (data *Data) GetRowActionType() string {
	if val, exists := data.TypeMap["RowAction"]; exists {
		return val.(string)
	}
	return ""
}

func (data *Data) SetRowActionHXPost(value string) {
	data.TypeMap["RowActionHXPost"] = value
}

func (data *Data) GetRowActionHXPost() (string, error) {
	if val, exists := data.TypeMap["RowActionHXPost"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'RowActionHXPost' does not exist")
}

func (data *Data) SetRowActionHXPostWithID(value string) {
	data.TypeMap["RowActionHXPostWithID"] = value
}

func (data *Data) GetRowActionHXPostWithID() (string, error) {
	if val, exists := data.TypeMap["RowActionHXPostWithID"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'RowActionHXPostWithID' does not exist")
}

func (data *Data) SetRowActionHXPostWithIDAsQueryParam(value string) {
	data.TypeMap["RowActionHXPostWithIDAsQueryParam"] = value
}

func (data *Data) GetRowActionHXPostWithIDAsQueryParam() (string, error) {
	if val, exists := data.TypeMap["RowActionHXPostWithIDAsQueryParam"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'RowActionHXPostWithIDAsQueryParam' does not exist")
}

func (data *Data) SetRowActionName(value string) {
	data.TypeMap["RowActionName"] = value
}

func (data *Data) GetRowActionName() (string, error) {
	if val, exists := data.TypeMap["RowActionName"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'RowActionName' does not exist")
}

func (data *Data) SetRowActionHXTarget(value string) {
	data.TypeMap["RowActionHXTarget"] = value
}

func (data *Data) GetRowActionHXTarget() (string, error) {
	if val, exists := data.TypeMap["RowActionHXTarget"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'RowActionHXTarget' does not exist")
}

func (data *Data) SetAdditionalDataInputs(value []DataInput) {
	data.TypeMap["AdditionalDataInputs"] = value
}

func (data *Data) GetAdditionalDataInputs() ([]DataInput, error) {
	if val, exists := data.TypeMap["AdditionalDataInputs"]; exists {
		return val.([]DataInput), nil
	}
	return nil, errors.New("key 'AdditionalDataInputs' does not exist")
}

func (data *Data) SetPages(value []PaginationButton) {
	data.TypeMap["Pages"] = value
}

func (data *Data) GetPages() ([]PaginationButton, error) {
	if val, exists := data.TypeMap["Pages"]; exists {
		return val.([]PaginationButton), nil
	}
	return nil, errors.New("key 'Pages' does not exist")
}

func (data *Data) SetNextPage(value int) {
	data.TypeMap["NextPage"] = value
}

func (data *Data) GetNextPage() (int, error) {
	if val, exists := data.TypeMap["NextPage"]; exists {
		return val.(int), nil
	}
	return 0, errors.New("key 'NextPage' does not exist")
}

func (data *Data) SetPrevPage(value int) {
	data.TypeMap["PrevPage"] = value
}

func (data *Data) GetPrevPage() (int, error) {
	if val, exists := data.TypeMap["PrevPage"]; exists {
		return val.(int), nil
	}
	return 0, errors.New("key 'PrevPage' does not exist")
}

func (data *Data) SetPageNumber(value int) {
	data.TypeMap["PageNumber"] = value
}

func (data *Data) GetPageNumber() int {
	if val, exists := data.TypeMap["PageNumber"]; exists {
		return val.(int)
	}
	return 1
}

func (data *Data) SetMove(value bool) {
	data.TypeMap["Move"] = value
}

func (data *Data) GetMove() (bool, error) {
	if val, exists := data.TypeMap["Move"]; exists {
		return val.(bool), nil
	}
	return false, errors.New("key 'Move' does not exist")
}

func (data *Data) SetCount(value int) {
	data.TypeMap["Count"] = value
}

func (data *Data) GetCount() int {
	if val, exists := data.TypeMap["Count"]; exists {
		return val.(int)
	}
	return 0
}

func (data *Data) SetPlaceHolder(placeHolder bool) {
	data.TypeMap["PlaceHolder"] = placeHolder
}

func (data *Data) GetPlaceHolder() bool {
	if _, exists := data.TypeMap["PlaceHolder"]; exists {
		return true
	}
	return false
}

func (data *Data) SetRequestOrigin(value string) {
	data.TypeMap["RequestOrigin"] = value
}

func (data *Data) GetOriginRequest() string {
	if val, exists := data.TypeMap["RequestOrigin"]; exists {
		return val.(string)
	}
	return ""
}

// Assign the thing data (ID, Label, etc..)to the dataTypeMap
func (data *Data) SetDetailesData(value map[string]any) {
	maps := []map[string]any{data.TypeMap, value}
	data.TypeMap = MergeMaps(maps)
}

// retrieves Item value for the template
func (data *Data) GetItem() map[string]interface{} {
	if raw, exists := data.TypeMap["Item"]; exists {
		if val, ok := raw.(map[string]interface{}); ok {
			return val
		}
	}
	return nil
}

// Set the Edit state
func (data *Data) SetEdit(value bool) {
	data.TypeMap["Edit"] = value
}

// Get the Edit state
func (data *Data) GetEdit() bool {
	if raw, exists := data.TypeMap["Edit"]; exists {
		if val, ok := raw.(bool); ok {
			return val
		}
	}
	return false
}

// Get the Edit state
func (data *Data) GetTypeMap() map[string]interface{} {
	return data.TypeMap
}

// check if the Box is available while previewing the Item
func (data *Data) IsBoxAvailable() bool {
	item := data.GetItem()
	if item == nil {
		logg.Infof("Item is not set. Please set an Item before checking if the Box is available.")
		return false
	}

	boxID, exists := item["BoxID"]
	if !exists {
		logg.Infof("BoxID key is missing in the Item.")
		return false
	}

	if uuidValue, ok := boxID.(uuid.UUID); ok {
		if uuidValue == uuid.Nil {
			return false
		}
		return true
	}

	// If BoxID is a string, parse it as a UUID
	if boxIDStr, ok := boxID.(string); ok {
		parsedUUID, err := uuid.FromString(boxIDStr)
		if err != nil || parsedUUID == uuid.Nil {
			logg.Infof("BoxID is either invalid or uuid.Nil.")
			return false
		}
		return true
	}

	// If BoxID is of an unexpected type
	logg.Infof("BoxID is of an unsupported type.")
	return false
}

// check if the Shelf is available while previewing the Item
func (data *Data) IsShelfAvailable() bool {
	item := data.GetItem()
	if item == nil {
		logg.Warning("Item is not set. Please set an Item before checking if the Shelf is available.")
		return false
	}

	ShelfID, exists := item["ShelfID"]
	if !exists {
		logg.Warning("ShelfID key is missing in the Item.")
		return false
	}

	if uuidValue, ok := ShelfID.(uuid.UUID); ok {
		if uuidValue == uuid.Nil {
			return false
		}
		return true
	}

	// If ShelfID is a string, parse it as a UUID
	if shelfIDStr, ok := ShelfID.(string); ok {
		parsedUUID, err := uuid.FromString(shelfIDStr)
		if err != nil || parsedUUID == uuid.Nil {
			logg.Warning("ShelfID is either invalid or uuid.Nil.")
			return false
		}
		return true
	}

	// If ShelfID is of an unexpected type
	logg.Warning("ShelfID is of an unsupported type.")
	return false
}

// check if the Area is available while previewing the Item
func (data *Data) IsAreaAvailable() bool {
	item := data.GetItem()
	if item == nil {
		logg.Warning("Item is not set. Please set an Item before checking if the Area is available.")
		return false
	}

	AreaID, exists := item["AreaID"]
	if !exists {
		logg.Warning("AreaID key is missing in the Item.")
		return false
	}

	if uuidValue, ok := AreaID.(uuid.UUID); ok {
		if uuidValue == uuid.Nil {
			return false
		}
		return true
	}

	// If AreaID is a string, parse it as a UUID
	if AreaIDStr, ok := AreaID.(string); ok {
		parsedUUID, err := uuid.FromString(AreaIDStr)
		if err != nil || parsedUUID == uuid.Nil {
			logg.Warning("AreaID is either invalid or uuid.Nil.")
			return false
		}
		return true
	}

	// If AreaID is of an unexpected type
	logg.Warning("AreaID is of an unsupported type.")
	return false
}

// SetListRowTemplateOptions sets the ListRowTemplateOptions in the TypeMap.
func (data *Data) SetListRowTemplateOptions(value ListRowTemplateOptions) {
	data.TypeMap["ListRowTemplateOptions"] = value.Map()
}

// GetListRowTemplateOptions retrieves the ListRowTemplateOptions from the TypeMap.
func (data *Data) GetListRowTemplateOptions() ListRowTemplateOptions {
	if val, exists := data.TypeMap["ListRowTemplateOptions"]; exists {
		if optionsMap, ok := val.(map[string]interface{}); ok {
			return ListRowTemplateOptions{
				RowHXGet:                          optionsMap["RowHXGet"].(string),
				RowAction:                         optionsMap["RowAction"].(bool),
				RowActionType:                     optionsMap["RowActionType"].(string),
				RowActionHXPost:                   optionsMap["RowActionHXPost"].(string),
				RowActionHXPostWithID:             optionsMap["RowActionHXPostWithID"].(string),
				RowActionHXPostWithIDAsQueryParam: optionsMap["RowActionHXPostWithIDAsQueryParam"].(string),
				RowActionName:                     optionsMap["RowActionName"].(string),
				RowActionHXTarget:                 optionsMap["RowActionHXTarget"].(string),
				HideMoveCol:                       optionsMap["HideMoveCol"].(bool),
			}
		}
	}
	return ListRowTemplateOptions{}
}
