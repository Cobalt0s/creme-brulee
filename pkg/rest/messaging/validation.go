package messaging

import (
	"net/mail"
)

func ValidEmail(propertyName, data string) error {
	_, err := mail.ParseAddress(data)
	if err != nil {
		return InvalidEmailFormat{Name: propertyName}
	}
	return nil
}

type InvalidEmailFormat struct {
	Name string
}

func (f InvalidEmailFormat) Error() string {
	return "invalid email format"
}

func (f InvalidEmailFormat) ToResponse() *errorResponse {
	return CreateInvalidFieldError(f.Name, f.Error())
}
