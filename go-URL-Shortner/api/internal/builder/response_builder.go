package builder

import (
	"go-url-shortener/internal/models"
	"time"
)

type ResponseBuilder struct {
	resp models.Response
}

func NewReponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{
		resp: models.Response{},
	}
}

func (rb *ResponseBuilder) SetURL(url string) *ResponseBuilder {
	rb.resp.URL = url
	return rb
}
func (rb *ResponseBuilder) SetCustomSort(short string) *ResponseBuilder {
	rb.resp.CustomShort = short
	return rb
}

func (rb *ResponseBuilder) SetExpiry(exp time.Duration) *ResponseBuilder {
	rb.resp.Expiry = exp
	return rb
}

func (rb *ResponseBuilder) SetRateLimit(remaining int, reset time.Duration) *ResponseBuilder {
	rb.resp.XRateRemaining = remaining
	rb.resp.XRateLimitReset = reset
	return rb
}

func (rb *ResponseBuilder) Build() models.Response {
	return rb.resp
}
