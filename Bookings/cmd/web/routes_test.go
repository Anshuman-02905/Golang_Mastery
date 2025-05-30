package main

import (
	"fmt"
	"testing"

	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/config"
	"github.com/go-chi/chi/v4"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig
	mux := routes(&app)
	switch v := mux.(type) {
	case *chi.Mux:
		//do nothing
	default:
		t.Errorf(fmt.Sprintf("type is not chi.mux but is %T,%v", v, v))
	}
}
