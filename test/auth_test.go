package tests

import (
	"net/http"
	"net/url"
	"testing"

	featureUrl "github.com/Sanchir01/go-shortener/internal/feature/url"
	"github.com/Sanchir01/go-shortener/internal/feature/user"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
)

const (
	host = "localhost:4200"
)

func CreateHttpExpect(t *testing.T) *httpexpect.Expect {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "/api/v1",
	}
	e := httpexpect.Default(t, u.String())
	return e
}

func Test_Auth_Register_Login(t *testing.T) {
	e := CreateHttpExpect(t)

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, true, 10)
	username := gofakeit.Username()

	e.POST("/auth/register").WithJSON(user.AuthRequest{
		Email:    email,
		Password: password,
		Username: username,
	}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		ContainsKey("status")

	e.POST("/auth/login").WithJSON(user.AuthRequest{
		Email:    email,
		Password: password,
	}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("status").
		ContainsKey("email").
		ContainsKey("username")
}

func Test_Create_Url(t *testing.T) {
	e := CreateHttpExpect(t)

	e.POST("/url/save").WithJSON(featureUrl.GetAllUrlResponse{})
}