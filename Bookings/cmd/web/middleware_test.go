package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func Test_WriteToConsole(t *testing.T) {
	//Create a handler
	//Create a pipe for the println
	stdout := os.Stdout
	defer func() {
		os.Stdout = stdout
	}()
	r, w, _ := os.Pipe()

	os.Stdout = w

	called := false

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	handler := WriteToConsole(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	_ = w.Close()
	val, _ := io.ReadAll(r)
	str_val := string(val)

	if !called {
		t.Errorf("Expected Next Handler to be called but it was  not")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", rec.Code)
	}
	if !strings.Contains(str_val, "Hit the page\n") {
		t.Errorf("expected Hit the page in the output but  found %s", str_val)

	}

}

func Test_NoServer(t *testing.T) {
	// called := false
	// test_handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	called = true
	// 	w.WriteHeader(http.StatusOK)
	// })
	// handler := NoServer(test_handler)
	// req := httptest.NewRequest(http.MethodGet, "/", nil)
	// rec := httptest.NewRecorder()
	// handler.ServeHTTP(rec, req)

	// if !called {
	// 	t.Errorf("Expected Next Handler to be called but it was  not")

	// }
	var myH myHandler

	h := NoServer(&myH)
	switch v := h.(type) {
	case http.Handler:
		//do nothing
	default:
		t.Error(fmt.Sprintf("type is not handler but is %T", v))
	}

}
func Test_SessionLoad(t *testing.T) {
	// called := false
	// test_handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	called = true
	// 	w.WriteHeader(http.StatusOK)
	// })
	// handler := SessionLoad(test_handler)
	// req := httptest.NewRequest(http.MethodGet, "/", nil)
	// rec := httptest.NewRecorder()
	// handler.ServeHTTP(rec, req)

	// if !called {
	// 	t.Errorf("Expected Next Handler to be called but it was  not")

	// }
	var myH myHandler

	h := SessionLoad(&myH)
	switch v := h.(type) {
	case http.Handler:
		//do nothing
	default:
		t.Error(fmt.Sprintf("type is not handler but is %T", v))
	}
}
