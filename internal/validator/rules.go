package validator

import (
	"shareapp/internal/data"
)

func validateEmail(v *Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(Matches(email, EmailRX), "email", "must be a valid email address")
}

func validatePasswordPlaintext(v *Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *Validator, user *data.User) {
	v.Check(user.Username != "", "username", "must be provided")
	v.Check(len(user.Username) <= 500, "username", "must not be more than 500 bytes long")

	validateEmail(v, user.Email)

	v.Check(Matches(user.Email, EmailRX), "email", "must be a valid email address")
	//v.Check(user.PasswordHash != "", "password", "must be provided")

}
