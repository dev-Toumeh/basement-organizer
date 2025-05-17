package auth

import (
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/sessions"
)

const (
	REGISTER_FAILED_MESSAGE string = "register failed"
	COOKIE_NAME             string = "mycookie"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key           = []byte("super-secret-key")
	store         = sessions.NewCookieStore(key)
	errorMessages = &[]string{}
)

type AuthDatabase interface {
	CreateNewUser(ctx context.Context, username string, passwordhash string) error
	UserByField(ctx context.Context, field string, value string) (User, error)
	UpdateUser(ctx context.Context, user User) error
	UserExist(ctx context.Context, username string) bool
	ErrorExist() error
}

type User struct {
	Id           uuid.UUID
	Username     string
	PasswordHash string
	email        string
}

// use this struct only for validation
type InputUser struct {
	Id              string `validate:"-"`
	Username        string `validate:"required,min=6,max=20"`
	Password        string `validate:"required,min=8,password_strength"`
	PasswordConfirm string `validate:"required,eqfield=Password"`
	Email           string `validate:"omitempty,email"`
}

// use this struct for generating responses with data error/update user
type UserResponseData struct {
	InputUserData InputUser `json:"input_user_data"`
	ErrorMessages []string  `json:"error_messages"`
}

const (
	ID               string = "id"
	USERNAME         string = "username"
	PASSWORD         string = "password"
	PASSWORD_CONFIRM string = "password-confirm"
	PASSWORDHASH     string = "passwordhash"
	EMAIL            string = "email"
	FAILED_MESSAGE   string = "we was not able to create the new User please comeback later"
)

// this function will check the type of the request
// if it is from type post it will register the user otherwise it will generate the register template
func RegisterHandler(db AuthDatabase) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			registerUser(w, r, db)
			return
		}
		if r.Method == http.MethodGet {
			generateRegisterPage(w, r)
			return
		}
		w.Header().Add("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		logg.Debug(w, "Method:'", r.Method, "' not allowed")
	}
}

func generateRegisterPage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := Authenticated(r)
	user, _ := UserSessionData(r)
	data := templates.PageTemplate{
		Title:         "Register",
		Authenticated: authenticated,
		User:          user,
	}

	if err := templates.Render(w, templates.TEMPLATE_REGISTER_PAGE, data); err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}

func registerUser(w http.ResponseWriter, r *http.Request, db AuthDatabase) {

	//  1. Validate the input
	inputFromPost := InputUser{
		Username:        r.PostFormValue(USERNAME),
		Password:        r.PostFormValue(PASSWORD),
		PasswordConfirm: r.PostFormValue(PASSWORD_CONFIRM),
		Email:           r.PostFormValue(EMAIL),
	}
	logg.Debugf("register input values %v", inputFromPost)

	// validInputUser, err := validateRegisterInput(inputFromPost)
	// if err != nil {
	// 	RenderValidateErrorMessages(w, inputFromPost)
	// 	return
	// }

	// 2. Put the data into struct from type user
	newUser, err := user(inputFromPost)
	if err != nil {
		logg.Err(err)
		templates.RenderErrorNotification(w, FAILED_MESSAGE)
	}

	// 3. Check if the username exist
	ctx := context.TODO()
	exist := db.UserExist(ctx, newUser.Username)
	if exist {
		message := fmt.Sprintf(`User already exists: "%s"`, newUser.Username)
		logg.Debug(message)
		*errorMessages = append(*errorMessages, message)
		RenderValidateErrorMessages(w, inputFromPost)
		return

	} else {
		err = db.CreateNewUser(ctx, newUser.Username, newUser.PasswordHash)
		// 4. Create the new Record
		if err != nil {
			logg.Debug(err)
			templates.RenderErrorNotification(w, FAILED_MESSAGE)
			return
		}

		server.RedirectWithSuccessNotification(w, "/login", fmt.Sprintf("the user %s was created successfully", newUser.Username))
		logg.Debugf("User %s registered successfully:", newUser.Username)
		return
	}
}

// return filled struct from type user
// generate new id in case of register new user
func user(inputUser InputUser) (User, error) {
	var id uuid.UUID
	var err error

	idStr := inputUser.Id
	if idStr != "" {
		id, err = uuid.FromString(idStr)
		if err != nil {
			logg.Errf("error parsing UUID from request: %v", err)
			return User{}, err
		}
	} else {
		id = uuid.Must(uuid.NewV4())
	}
	hashedPassword, err := hashPassword(inputUser.Password)
	if err != nil {
		logg.Err(err)
		return User{}, err
	}
	newUser := User{
		Id:           id,
		Username:     inputUser.Username,
		PasswordHash: hashedPassword,
	}

	return newUser, nil
}

// render register validate error Messages
func RenderValidateErrorMessages(w io.Writer, inputUser InputUser) {
	templateError := UserResponseData{InputUserData: inputUser, ErrorMessages: *errorMessages}
	err := templates.Render(w, templates.TEMPLATE_REGISTER_FORM, templateError)
	if err != nil {
		logg.Debug(err)
		templates.RenderErrorNotification(w, FAILED_MESSAGE)
	}
	*errorMessages = []string{}
}
