package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// a custom form struct
type Form struct {
	url.Values
	Errors errors
}

// Returns true if no errors otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// initia lize a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {

			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// if form field is in post and not empty
func (f *Form) Has(field string) bool {

	//x := r.Form.Get(field)
	x := f.Get(field)

	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}

// Min Length checks for string minimum length
func (f *Form) MinLength(field string, length int) bool {
	//x := r.Form.Get(field)
	x := f.Get(field)

	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This filed must be at least %d chatecters long", length))
		return false
	}
	return true

}

func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid Email Address")
	}
}
