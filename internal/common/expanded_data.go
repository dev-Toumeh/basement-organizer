package common

import (
	"errors"
	"fmt"
)

// Data structure definition
type data struct {
	ListTemplate,
	Title string
	Authenticated  bool
	User           string
	Debug          bool
	NotFound       bool
	EnvDevelopment bool
	TypeMap        map[string]any
}

// Global instance of Data
var Data = &data{
	TypeMap: make(map[string]any),
}

func (data *data) SetTitle(value string) {
	data.TypeMap["Title"] = value
}

func (data *data) GetTitle() (string, error) {
	if val, exists := data.TypeMap["Title"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'Title' does not exist")
}

func (data *data) SetAuthenticated(value bool) {
	data.TypeMap["Authenticated"] = value
}

func (data *data) GetAuthenticated() (bool, error) {
	if val, exists := data.TypeMap["Authenticated"]; exists {
		return val.(bool), nil
	}
	return false, errors.New("key 'Authenticated' does not exist")
}

func (data *data) SetUser(value string) {
	data.TypeMap["User"] = value
}

func (data *data) GetUser() (string, error) {
	if val, exists := data.TypeMap["User"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'User' does not exist")
}

func (data *data) SetDebug(value bool) {
	data.TypeMap["Debug"] = value
}

func (data *data) GetDebug() (bool, error) {
	if val, exists := data.TypeMap["Debug"]; exists {
		return val.(bool), nil
	}
	return false, errors.New("key 'Debug' does not exist")
}

func (data *data) SetNotFound(value bool) {
	data.TypeMap["NotFound"] = value
}

func (data *data) GetNotFound() (bool, error) {
	if val, exists := data.TypeMap["NotFound"]; exists {
		return val.(bool), nil
	}
	return false, errors.New("key 'NotFound' does not exist")
}

func (data *data) SetEnvDevelopment(value bool) {
	data.TypeMap["EnvDevelopment"] = value
}

func (data *data) GetEnvDevelopment() (bool, error) {
	if val, exists := data.TypeMap["EnvDevelopment"]; exists {
		return val.(bool), nil
	}
	return false, errors.New("key 'EnvDevelopment' does not exist")
}

func (data *data) SetFormID(value string) {
	data.TypeMap["FormID"] = value
}

func (data *data) GetFormID() (string, error) {
	if val, exists := data.TypeMap["FormID"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'FormID' does not exist")
}

func (data *data) SetFormHXGet(value string) {
	data.TypeMap["FormHXGet"] = value
}

func (data *data) GetFormHXGet() (string, error) {
	if val, exists := data.TypeMap["FormHXGet"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'FormHXGet' does not exist")
}

func (data *data) SetFormHXPost(value string) {
	data.TypeMap["FormHXPost"] = value
}

func (data *data) GetFormHXPost() (string, error) {
	if val, exists := data.TypeMap["FormHXPost"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'FormHXPost' does not exist")
}

func (data *data) SetFormHXTarget(value string) {
	data.TypeMap["FormHXTarget"] = value
}

func (data *data) GetFormHXTarget() (string, error) {
	if val, exists := data.TypeMap["FormHXTarget"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'FormHXTarget' does not exist")
}

func (data *data) SetRowHXGet(value string) {
	data.TypeMap["RowHXGet"] = value
}

func (data *data) GetRowHXGet() (string, error) {
	if val, exists := data.TypeMap["RowHXGet"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'RowHXGet' does not exist")
}

func (data *data) SetSearchInput(value bool) {
	data.TypeMap["SearchInput"] = value
}

func (data *data) GetSearchInput() (bool, error) {
	if val, exists := data.TypeMap["SearchInput"]; exists {
		return val.(bool), nil
	}
	return false, errors.New("key 'SearchInput' does not exist")
}

func (data *data) SetSearchInputLabel(value string) {
	data.TypeMap["SearchInputLabel"] = value
}

func (data *data) GetSearchInputLabel() (string, error) {
	if val, exists := data.TypeMap["SearchInputLabel"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'SearchInputLabel' does not exist")
}

func (data *data) SetSearchInputValue(value string) {
	data.TypeMap["SearchInputValue"] = value
}

func (data *data) GetSearchInputValue() (string, error) {
	if val, exists := data.TypeMap["SearchInputValue"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'SearchInputValue' does not exist")
}

func (data *data) SetPagination(value bool) {
	data.TypeMap["Pagination"] = value
}

func (data *data) GetPagination() (bool, error) {
	if val, exists := data.TypeMap["Pagination"]; exists {
		return val.(bool), nil
	}
	return false, errors.New("key 'Pagination' does not exist")
}

func (data *data) SetCurrentPageNumber(value int) {
	data.TypeMap["CurrentPageNumber"] = value
}

func (data *data) GetCurrentPageNumber() (int, error) {
	if val, exists := data.TypeMap["CurrentPageNumber"]; exists {
		return val.(int), nil
	}
	return 0, errors.New("key 'CurrentPageNumber' does not exist")
}

func (data *data) SetLimit(value int) {
	data.TypeMap["Limit"] = value // Store as int directly
}

func (data *data) GetLimit() (int, error) {
	if val, exists := data.TypeMap["Limit"]; exists {
		// Ensure type assertion matches storage type
		if intVal, ok := val.(int); ok {
			return intVal, nil
		}
		return 0, errors.New("key 'Limit' exists but is not an int")
	}
	return 0, errors.New("key 'Limit' does not exist")
}

func (data *data) GetLimitStr() (string, error) {
	if val, exists := data.TypeMap["Limit"]; exists {
		return fmt.Sprint(val), nil
	}
	return "", errors.New("key 'Limit' does not exist")
}

func (data *data) SetShowLimit(value bool) {
	data.TypeMap["ShowLimit"] = value
}

func (data *data) GetShowLimit() (bool, error) {
	if val, exists := data.TypeMap["ShowLimit"]; exists {
		return val.(bool), nil
	}
	return false, errors.New("key 'ShowLimit' does not exist")
}

func (data *data) SetPaginationButtons(value []PaginationButton) {
	data.TypeMap["PaginationButtons"] = value
}

func (data *data) GetPaginationButtons() ([]PaginationButton, error) {
	if val, exists := data.TypeMap["PaginationButtons"]; exists {
		return val.([]PaginationButton), nil
	}
	return nil, errors.New("key 'PaginationButtons' does not exist")
}

func (data *data) SetMoveButtonHXTarget(value string) {
	data.TypeMap["MoveButtonHXTarget"] = value
}

func (data *data) GetMoveButtonHXTarget() (string, error) {
	if val, exists := data.TypeMap["MoveButtonHXTarget"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'MoveButtonHXTarget' does not exist")
}

func (data *data) SetRows(value []ListRow) {
	data.TypeMap["Rows"] = value
}

func (data *data) GetRows() ([]ListRow, error) {
	if val, exists := data.TypeMap["Rows"]; exists {
		return val.([]ListRow), nil
	}
	return nil, errors.New("key 'Rows' does not exist")
}

func (data *data) SetRowAction(value bool) {
	data.TypeMap["RowAction"] = value
}

func (data *data) GetRowAction() (bool, error) {
	if val, exists := data.TypeMap["RowAction"]; exists {
		return val.(bool), nil
	}
	return false, errors.New("key 'RowAction' does not exist")
}

func (data *data) SetRowActionHXPost(value string) {
	data.TypeMap["RowActionHXPost"] = value
}

func (data *data) GetRowActionHXPost() (string, error) {
	if val, exists := data.TypeMap["RowActionHXPost"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'RowActionHXPost' does not exist")
}

func (data *data) SetRowActionHXPostWithID(value string) {
	data.TypeMap["RowActionHXPostWithID"] = value
}

func (data *data) GetRowActionHXPostWithID() (string, error) {
	if val, exists := data.TypeMap["RowActionHXPostWithID"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'RowActionHXPostWithID' does not exist")
}

func (data *data) SetRowActionHXPostWithIDAsQueryParam(value string) {
	data.TypeMap["RowActionHXPostWithIDAsQueryParam"] = value
}

func (data *data) GetRowActionHXPostWithIDAsQueryParam() (string, error) {
	if val, exists := data.TypeMap["RowActionHXPostWithIDAsQueryParam"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'RowActionHXPostWithIDAsQueryParam' does not exist")
}

func (data *data) SetRowActionName(value string) {
	data.TypeMap["RowActionName"] = value
}

func (data *data) GetRowActionName() (string, error) {
	if val, exists := data.TypeMap["RowActionName"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'RowActionName' does not exist")
}

func (data *data) SetRowActionHXTarget(value string) {
	data.TypeMap["RowActionHXTarget"] = value
}

func (data *data) GetRowActionHXTarget() (string, error) {
	if val, exists := data.TypeMap["RowActionHXTarget"]; exists {
		return val.(string), nil
	}
	return "", errors.New("key 'RowActionHXTarget' does not exist")
}

func (data *data) SetAdditionalDataInputs(value []DataInput) {
	data.TypeMap["AdditionalDataInputs"] = value
}

func (data *data) GetAdditionalDataInputs() ([]DataInput, error) {
	if val, exists := data.TypeMap["AdditionalDataInputs"]; exists {
		return val.([]DataInput), nil
	}
	return nil, errors.New("key 'AdditionalDataInputs' does not exist")
}

func (data *data) SetPages(value []PaginationButton) {
	data.TypeMap["Pages"] = value
}

func (data *data) GetPages() ([]PaginationButton, error) {
	if val, exists := data.TypeMap["Pages"]; exists {
		return val.([]PaginationButton), nil
	}
	return nil, errors.New("key 'Pages' does not exist")
}

func (data *data) SetNextPage(value int) {
	data.TypeMap["NextPage"] = value
}

func (data *data) GetNextPage() (int, error) {
	if val, exists := data.TypeMap["NextPage"]; exists {
		return val.(int), nil
	}
	return 0, errors.New("key 'NextPage' does not exist")
}

func (data *data) SetPrevPage(value int) {
	data.TypeMap["PrevPage"] = value
}

func (data *data) GetPrevPage() (int, error) {
	if val, exists := data.TypeMap["PrevPage"]; exists {
		return val.(int), nil
	}
	return 0, errors.New("key 'PrevPage' does not exist")
}

func (data *data) SetPageNumber(value int) {
	data.TypeMap["PageNumber"] = value
}

func (data *data) GetPageNumber() (int, error) {
	if val, exists := data.TypeMap["PageNumber"]; exists {
		return val.(int), nil
	}
	return 0, errors.New("key 'PageNumber' does not exist")
}

func (data *data) SetMove(value bool) {
	data.TypeMap["Move"] = value
}

func (data *data) GetMove() (bool, error) {
	if val, exists := data.TypeMap["Move"]; exists {
		return val.(bool), nil
	}
	return false, errors.New("key 'Move' does not exist")
}

func (data *data) SetCount(value int) {
	data.TypeMap["Count"] = value
}

func (data *data) GetCount() (int, error) {
	if val, exists := data.TypeMap["Count"]; exists {
		return val.(int), nil
	}
	return 0, errors.New("key 'Count' does not exist")
}
