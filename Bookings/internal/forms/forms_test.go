package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	isValid := form.Valid()

	if !isValid {
		t.Error("got invalid when it should have been valid")
	}
}
func TestForm_New(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	form := New(r.PostForm)
	if form == nil {
		t.Error("returned value is not of *Form type")
	}
}

func TestForm_Required(t *testing.T) {

	// r := httptest.NewRequest("POST", "/whatever", nil)
	// r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Parse the form to populate r.PostForm

	// err := r.ParseForm()
	// if err != nil {
	// 	panic(err)
	// }
	form_url := url.Values{}

	form := New(form_url)
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Form shows valid when requried field missing")
	}

	form_url = url.Values{}
	form_url.Add("a", "aa")
	form_url.Add("b", "bb")
	form_url.Add("c", "cc")

	// r = httptest.NewRequest("POST", "/whatever", strings.NewReader(form_url.Encode()))
	// r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// _ = r.ParseForm()

	form = New(form_url)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Form is not showing  valid when it should")
	}
}

func TestForm_Has(t *testing.T) {

	//create a reeust

	form_url := url.Values{}

	form := New(form_url)

	flag := form.Has("")
	if flag {
		t.Error("Error as as expecting false as This field cannot be blank but found true")
	}

	form_url = url.Values{}
	form_url.Add("a", "aa")
	form_url.Add("b", "bb")
	form_url.Add("c", "cc")

	form = New(form_url)

	flag = form.Has("c")

	if !flag {
		t.Error("Error as as expecting field to be found but not found")
	}
}

func TestForm_MinLength(t *testing.T) {

	// r := httptest.NewRequest("POST", "/whatever", nil)
	// r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	postedValues := url.Values{}
	form := New(postedValues)

	flag := form.MinLength("", 0)

	if !flag {
		t.Error("Error as as expecting false as This field cannot be blank but found true")
	}

	postedValues = url.Values{}
	postedValues.Add("a", "aa")
	postedValues.Add("b", "bb")
	postedValues.Add("c", "cc")

	form = New(postedValues)

	flag = form.MinLength("a", 2)

	if !flag {
		t.Error("Error as as expecting min lenght to be 2 but not found so")
	}

	postedValues = url.Values{}
	postedValues.Add("a", "aaa")
	postedValues.Add("b", "bbb")
	postedValues.Add("c", "cc")

	// r = httptest.NewRequest("POST", "/whatever", strings.NewReader(postedValues.Encode()))
	// r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// _ = r.ParseForm()
	form = New(postedValues)
	flag = form.MinLength("c", 2)
	//fmt.Println(form.Errors)

	if !flag {
		t.Error("Error as as expecting min lenght to be 2 but not found so")
	}

}

func TestForm_IsEmail(t *testing.T) {

	// r := httptest.NewRequest("POST", "/whatever", nil)
	// r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// _ = r.ParseForm()
	postedValues := url.Values{}
	form := New(postedValues)

	form.IsEmail("")
	//fmt.Println(form.Errors.Get(""))

	if form.Valid() {
		t.Error("FORM is INVALID but shows it valid")
	}

	form_url := url.Values{}
	form_url.Add("email", "anshuman.manfal@gmail.com")

	// r = httptest.NewRequest("POST", "/whatever", strings.NewReader(form_url.Encode()))
	// r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// _ = r.ParseForm()
	postedValues = url.Values{}
	postedValues.Add("email", "anshuman.manfal@gmail.com")
	form = New(postedValues)

	form.IsEmail("email")
	//fmt.Println(form.Errors.Get("email"))

	if !form.Valid() {
		t.Error("FORM is VALID but shows it is invalid")
	}

}
