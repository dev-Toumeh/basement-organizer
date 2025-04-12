package auth

import (
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
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

// Validates the user registration input data.
// If successful, returns the validated input struct with nil.
// If validation fails, returns an empty input struct along with error.
// The function utilizes a global string array (errorMessages) to store validation error messages.
func validateRegisterInput(inputUser InputUser) (InputUser, error) {
	if !env.Development() {

		validate := validator.New(validator.WithRequiredStructEnabled())
		validate.RegisterValidation("password_strength", passwordStrengthValidator)

		if err := validate.Struct(inputUser); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				for _, validationErr := range validationErrors {
					switch validationErr.Field() {
					case "Username":
						if validationErr.Tag() == "required" {
							*errorMessages = append(*errorMessages, "The Username field is required but missing.")
						} else if validationErr.Tag() == "min" {
							*errorMessages = append(*errorMessages, "The Username must be at least 6 characters long.")
						} else if validationErr.Tag() == "max" {
							*errorMessages = append(*errorMessages, "The Username must be at most 20 characters long.")
						}
					case "Password":
						if validationErr.Tag() == "required" {
							*errorMessages = append(*errorMessages, "The Password field is required but missing.")
						} else if validationErr.Tag() == "min" {
							*errorMessages = append(*errorMessages, "The Password must be at least 8 characters long.")
						} else if validationErr.Tag() == "password_strength" {
							*errorMessages = append(*errorMessages, "The Password must contain at least one letter, one number, and one symbol.")
						}
					case "PasswordConfirm":
						if validationErr.Tag() == "required" {
							*errorMessages = append(*errorMessages, "The Password Confirm field is required but missing.")
						} else if validationErr.Tag() == "eqfield" {
							*errorMessages = append(*errorMessages, "The Password and Password Confirm fields must match.")
						}
					case "Email":
						if validationErr.Tag() == "email" {
							*errorMessages = append(*errorMessages, "The Email field must be a valid email address.")
						}
					default:
						*errorMessages = append(*errorMessages, fmt.Sprintf("Field '%s' is invalid: %s", validationErr.Field(), validationErr.Tag()))
					}
				}
			} else {
				*errorMessages = append(*errorMessages, err.Error())
			}

			//		logg.Err("User Input Validation failed")
			err := errors.New("validation failed")
			return InputUser{}, err
		} else {
			logg.Debug("User Input Validation succeeded")
			return inputUser, nil
		}
	}
	return inputUser, nil
}

// custom validate to make sure that the password has number, letters and symbols
func passwordStrengthValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check for at least one letter, one number, and one symbol
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSymbol := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	return hasLetter && hasNumber && hasSymbol
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
