package env

import (
	"basement/main/internal/logg"
)

var isDevelopment = false
var isProduction = true

// SetDevelopment sets environment setting to development.
//
// This setting is for making development tasks easier by reducing
// checks, validation strictness, adding more logging information etc.
func SetDevelopment() {
	isProduction = false
	isDevelopment = true
	logg.InfoForceOutput(3, "Environment is development")
}

// Development returns true if development environment is set.
// Default is false.
func Development() bool {
	return isDevelopment
}

// SetProduction sets environment setting to production.
//
// This is setting is for production use.
// Enables all checks, increases strictness, removes debug logging, etc.
func SetProduction() {
	isDevelopment = false
	isProduction = true
	logg.InfoForceOutput(3, "Environment is production")
}

// Production returns true if production environment is set.
// Default is true.
func Production() bool {
	return isProduction
}

// var configurations = make(map[string]string)

// DescribeDevelopmentOverride adds a label and description for a development config change.
//
// Use this to document and track environment-specific settings when custom logic is applied. Returns an error if the label already exists.
//
// Example:
// err := DescribeDevelopmentOverride("SkipPasswordValidations", "Disables password validations (password length, password strength) for registering a user.")
// func DescribeDevelopmentOverride(label string, description string) error {
// 	val, ok := configurations[label]
// 	if ok {
// 		return errors.New(fmt.Sprintf("Environment configuration '%v' already exists. Description:'%v'\n", label, val))
// 	}
// 	configurations[label] = description
// 	return nil
// }
