package auth

import (
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"context"
	"fmt"
	"net/http"
)

func UpdateHandler(db AuthDatabase) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			updateUser(w, r, db)
			return
		} else {
			generateUpdateFormWithData(w, r)
		}
		w.Header().Add("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		//logg.Debug(w, "Method:'", r.Method, "' not allowed")
	}
}

func updateUser(w http.ResponseWriter, r *http.Request, db AuthDatabase) {

	//  1. Validate the input
	_, id := UserSessionData(r)
	inputFromPost := InputUser{
		Id:              id,
		Username:        r.PostFormValue(USERNAME),
		Password:        r.PostFormValue(PASSWORD),
		PasswordConfirm: r.PostFormValue(PASSWORD_CONFIRM),
		Email:           r.PostFormValue(EMAIL),
	}

	validInputUser, err := validateRegisterInput(inputFromPost)
	if err != nil {
		RenderValidateErrorMessages(w, inputFromPost)
		return
	}

	// 2. Put the data into struct from type user
	newUser, err := user(validInputUser)
	if err != nil {
		logg.Err(err)
		templates.RenderErrorNotification(w, FAILED_MESSAGE)
	}

	// 3. Check if the username exist
	ctx := context.TODO()
	exist := db.UserExist(ctx, newUser.Username)
	if exist {
		err = db.UpdateUser(ctx, newUser)
		if err != nil {
			logg.Err(err)
			templates.RenderErrorNotification(w, FAILED_MESSAGE)
		}
	} else {
		logg.Err("some thing wrong happened while checking the user exist")
		templates.RenderErrorNotification(w, FAILED_MESSAGE)
	}

	// 4. Create the new Record
	if err != nil {
		logg.Debug(err)
		templates.RenderErrorNotification(w, FAILED_MESSAGE)
		return
	}

	// https://htmx.org/headers/hx-location/
	w.Header().Add("HX-Location", "/")
	templates.RenderSuccessNotification(w, fmt.Sprintf("the user %s was updated successfully", newUser.Username))
	http.Redirect(w, r, "/", http.StatusOK)
	logg.Debugf("User %s updated successfully:", newUser.Username)
	return
}

func generateUpdateFormWithData(w http.ResponseWriter, r *http.Request) {
	username, userid := UserSessionData(r)

	inputUser := InputUser{
		Id:       userid,
		Username: username,
	}
	templateUpdateUserData := UserResponseData{
		InputUserData: inputUser,
		ErrorMessages: *errorMessages,
	}
	err := templates.Render(w, templates.TEMPLATE_UPDATE_FORM, templateUpdateUserData)
	if err != nil {
		logg.Debug(err)
		templates.RenderErrorNotification(w, FAILED_MESSAGE)
	}
	*errorMessages = []string{}
}
