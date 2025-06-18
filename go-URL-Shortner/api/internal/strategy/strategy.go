package strategy

import (
	"github.com/google/uuid"
	"go-url-shortener/internal/models"
	"errors"
)

type IDStrategy interface {
	Generate(body *models.Request) (string, error)
}

type UUIDStrategy struct{}

func (f UUIDStrategy) Generate(body *models.Request) (string, error) {
	return uuid.New().String()[:6], nil
}

type CustomStrategy struct{}

func (c CustomStrategy) Generate(body *models.Request) (string, error) {
	if body.CustomShort == "" {
		return "", errors.New("Custom short code is empty")
	}
	return body.CustomShort, nil
}
