package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"online/driver"
	"online/models"
	"online/repository"
	"strings"
	"testing"

	"github.com/labstack/echo"
)

func TestSignup(t *testing.T) {
	db := driver.TestDbConnection()
	database := Database{Connection: db}
	e := echo.New()
	t.Run("missing fields", func(t *testing.T) {
		body := `{
			"username":"",
			"email":"",
			"password":"",
			"role":""
		}`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		c := e.NewContext(req, resp)
		database.Signup(c)

		if want, got := http.StatusBadRequest, resp.Result().StatusCode; want != got {
			t.Fatalf("expected: %d, got: %d", want, got)
		}
	})

	t.Run("Invalid email", func(t *testing.T) {
		body := `{
			"username":"hari",
			"email":"Hareeshgmailcom",
			"password":"12345678",
			"role":"admin"
		}`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		c := e.NewContext(req, resp)
		database.Signup(c)

		if want, got := http.StatusBadRequest, resp.Result().StatusCode; want != got {
			t.Fatalf("expected: %d, got: %d", want, got)
		}
	})

	t.Run("Invalid role", func(t *testing.T) {
		body := `{
			"username":"hari",
			"email":"hareesh@gmailcom",
			"password":"12345678",
			"role":"customer"
		}`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		c := e.NewContext(req, resp)
		database.Signup(c)

		if want, got := http.StatusBadRequest, resp.Result().StatusCode; want != got {
			t.Fatalf("expected: %d, got: %d", want, got)
		}
	})

	t.Run("Checking the length of password", func(t *testing.T) {
		body := `{
			"username":"hari",
			"email":"hareesh@gmailcom",
			"password":"1234678",
			"role":"customer"
		}`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		c := e.NewContext(req, resp)
		database.Signup(c)

		if want, got := http.StatusBadRequest, resp.Result().StatusCode; want != got {
			t.Fatalf("expected: %d, got: %d", want, got)
		}
	})

	t.Run("signup successful", func(t *testing.T) {
		body := `{
			"username":"hari",
			"email":"hareesh@gmail.com",
			"password":"12345678",
			"role":"admin"
		}`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		c := e.NewContext(req, resp)
		database.Signup(c)

		if want, got := http.StatusOK, resp.Result().StatusCode; want != got {
			t.Fatalf("expected: %d, got: %d", want, got)
		}
	})

	t.Run("User already exist", func(t *testing.T) {
		body := `{
			"username":"hari",
			"email":"hareesh@gmailcom",
			"password":"12345678",
			"role":"admin"
		}`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		c := e.NewContext(req, resp)
		database.Signup(c)

		if want, got := http.StatusBadRequest, resp.Result().StatusCode; want != got {
			t.Fatalf("expected: %d, got: %d", want, got)
		}
	})
}

func TestLogin(t *testing.T) {
	db := driver.TestDbConnection()
	database := Database{Connection: db}
	e := echo.New()
	t.Run("missing fields", func(t *testing.T) {
		body := `{
			"email":"",
			"password":""
		}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		c := e.NewContext(req, resp)
		database.Login(c)

		if want, got := http.StatusBadRequest, resp.Result().StatusCode; want != got {
			t.Fatalf("expected: %d, got: %d", want, got)
		}
	})

	t.Run("Invalid email", func(t *testing.T) {
		body := `{
			"email":"Hareeshgmailcom",
			"password":"12345678"
		}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		c := e.NewContext(req, resp)
		database.Login(c)

		if want, got := http.StatusBadRequest, resp.Result().StatusCode; want != got {
			t.Fatalf("expected: %d, got: %d", want, got)
		}
	})

	t.Run("User not found", func(t *testing.T) {
		body := `{
			"email":"bharathi@gmail.com",
			"password":"65432178"
		}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		c := e.NewContext(req, resp)
		database.Login(c)

		if want, got := http.StatusNotFound, resp.Result().StatusCode; want != got {
			t.Fatalf("expected: %d, got: %d", want, got)
		}
	})

	t.Run("Incorrect password", func(t *testing.T) {
		body := `{
			"email":"hareesh@gmail.com",
			"password":"65432178"
		}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		c := e.NewContext(req, resp)
		database.Login(c)

		if want, got := http.StatusBadRequest, resp.Result().StatusCode; want != got {
			t.Fatalf("expected: %d, got: %d", want, got)
		}
	})

	t.Run("Login successful", func(t *testing.T) {
		body := `{
			"email":"hareesh@gmail.com",
			"password":"12345678"
		}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		c := e.NewContext(req, resp)
		database.Login(c)

		if want, got := http.StatusOK, resp.Result().StatusCode; want != got {
			t.Fatalf("expected: %d, got: %d", want, got)
		}
	})
}

func TestPostProduct(t *testing.T) {
	db := driver.TestDbConnection()
	database := Database{Connection: db}
	e := echo.New()
	e.POST("/postProduct", database.PostProduct)
	t.Run("missing fields", func(t *testing.T) {
		body := `{
			"brand_name": "dell",
			"product_price": "20000",
			"ram_capacity": "2GB",
			"ram_price": "2000"
		}`
		req := httptest.NewRequest(http.MethodPost, "/admin/postProduct", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		user := models.User{UserId: 1}
		auth, err := repository.ReadTokenByUserId(db, user)
		if err != nil {
			panic(errors.New("Token not found"))
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth.Token))

		resp := httptest.NewRecorder()
		e.ServeHTTP(resp, req)
		if want, got := http.StatusNotFound, resp.Result().StatusCode; want != got {
			t.Fatalf("expected: %d, got: %d", want, got)
		}
	})
}
