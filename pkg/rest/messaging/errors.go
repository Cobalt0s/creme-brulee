package messaging

import (
	"fmt"
)

type ValidationError interface {
	ToResponse() *errorResponse
	Error() string
}

type errorResponse struct {
	Error errorResponsePayload `json:"error"`
}

type errorResponsePayload struct {
	Code    string                 `json:"code"`
	Details map[string]interface{} `json:"details"`
}

type EmptyField struct {
	Name string
}

func (f EmptyField) Error() string {
	return fmt.Sprintf("%v cannot be empty", f.Name)
}

func (f EmptyField) ToResponse() *errorResponse {
	return CreateInvalidFieldError(f.Name, f.Error())
}

type ShortField struct {
	Name string
	Size int
}

func (f ShortField) Error() string {
	return fmt.Sprintf("%v cannot be less than %v", f.Name, f.Size)
}

func (f ShortField) ToResponse() *errorResponse {
	return CreateInvalidFieldError(f.Name, f.Error())
}

type LongField struct {
	Name string
	Size int
}

func (f LongField) Error() string {
	return fmt.Sprintf("%v must be larger than %v", f.Name, f.Size)
}

func (f LongField) ToResponse() *errorResponse {
	return CreateInvalidFieldError(f.Name, f.Error())
}

type UnknownEnumField struct {
	Name   string
	Values []string
}

func (f UnknownEnumField) Error() string {
	return fmt.Sprintf("%v must be one of %v", f.Name, f.Values)
}

func (f UnknownEnumField) ToResponse() *errorResponse {
	return CreateInvalidFieldError(f.Name, f.Error())
}

type InvalidField struct {
	Name   string
	Format string
}

func (i InvalidField) Error() string {
	return fmt.Sprintf("invalid field %v must be in %v format", i.Name, i.Format)
}

func (i InvalidField) ToResponse() *errorResponse {
	return CreateInvalidFieldError(i.Name, fmt.Sprintf("must be in %v format", i.Format))
}

func CreateMissingFieldError(fieldName string) *errorResponse {
	return &errorResponse{Error: errorResponsePayload{
		Code: "missing-field",
		Details: map[string]interface{}{
			"required": fieldName,
		},
	}}
}
func CreateInvalidFieldError(fieldName string, message string) *errorResponse {
	return &errorResponse{Error: errorResponsePayload{
		Code: "invalid-field",
		Details: map[string]interface{}{
			fieldName: message,
		},
	}}
}
func CreateUnsatisfiedRuleError(err error) *errorResponse {
	return &errorResponse{Error: errorResponsePayload{
		Code: "rule-unsatisfied",
		Details: map[string]interface{}{
			"rule": err.Error(),
		},
	}}
}
func CreateNotFoundError(err error) *errorResponse {
	return &errorResponse{Error: errorResponsePayload{
		Code: "not-found",
		Details: map[string]interface{}{
			"message": err.Error(),
		},
	}}
}
func CreateActionNotAllowedError(err error) *errorResponse {
	return &errorResponse{Error: errorResponsePayload{
		Code: "action-not-permitted",
		Details: map[string]interface{}{
			"message": err.Error(),
		},
	}}
}
func CreateGenericError(message string) *errorResponse {
	return &errorResponse{Error: errorResponsePayload{
		Code: "generic-error",
		Details: map[string]interface{}{
			"message": message,
		},
	}}
}
func NewValidationErrPayload(err error) error {
	return fmt.Errorf("unknown error during payload validation %v", err)
}
func NewValidationErrQP(err error) error {
	return fmt.Errorf("unknown error during queryparams validation %v", err)
}
