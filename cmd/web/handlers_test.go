package main

import (
	"net/http"
	"net/url"
	"testing"
	"thienel/lets-go/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/anippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	validCSRFToken := extractCSRFToken(t, body)

	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = `<form action="/user/signup" method="POST" novalidate>`
	)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  string
	}{
		{
			name:         "Valid submission",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusSeeOther,
		},
		{
			name:         "Invalid CSRF Token",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "wrongToken",
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Empty name",
			userName:     "",
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty email",
			userName:     validName,
			userEmail:    "",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid email",
			userName:     validName,
			userEmail:    "bob@example.",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Short password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "pa$$",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Duplicate email",
			userName:     validName,
			userEmail:    "dupe@example.com",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}
		})
	}
}

func TestSnippetCreate(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		code, header, _ := ts.get(t, "/snippet/create")

		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, header.Get("Location"), "/user/login")
	})

	validCSRFToken := ""

	t.Run("Authenticated", func(t *testing.T) {
		_, _, body := ts.get(t, "/user/login")
		validCSRFToken = extractCSRFToken(t, body)

		form := url.Values{}
		form.Add("email", "alice@example.com")
		form.Add("password", "pa$$word")
		form.Add("csrf_token", validCSRFToken)

		ts.postForm(t, "/user/login", form)

		code, _, body := ts.get(t, "/snippet/create")
		assert.Equal(t, code, http.StatusOK)
		assert.StringContains(t, body, `<form action="/snippet/create" method="POST">`)
	})

	const (
		validTitle   = "Valid title"
		validContent = "Valid content"
		validExpires = "7"
	)

	tests := []struct {
		name           string
		snippetTitle   string
		snippetContent string
		snippetExpires string
		csrfToken      string
		wantCode       int
		wantLocation   string
	}{
		{
			name:           "Valid submission",
			snippetTitle:   validTitle,
			snippetContent: validContent,
			snippetExpires: validExpires,
			csrfToken:      validCSRFToken,
			wantCode:       http.StatusSeeOther,
			wantLocation:   "/snippet/view",
		},
		{
			name:           "Invalid CSRF Token",
			snippetTitle:   validTitle,
			snippetContent: validContent,
			snippetExpires: validExpires,
			csrfToken:      "Invalid token",
			wantCode:       http.StatusSeeOther,
		},
		{
			name:           "Title blank",
			snippetTitle:   "",
			snippetContent: validContent,
			snippetExpires: validExpires,
			csrfToken:      validCSRFToken,
			wantCode:       http.StatusUnprocessableEntity,
		},
		{
			name:           "Title exceeds limit characters",
			snippetTitle:   "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789" + "1",
			snippetContent: validContent,
			snippetExpires: validExpires,
			csrfToken:      validCSRFToken,
			wantCode:       http.StatusUnprocessableEntity,
		},
		{
			name:           "Content blank",
			snippetTitle:   validTitle,
			snippetContent: "",
			snippetExpires: validExpires,
			csrfToken:      validCSRFToken,
			wantCode:       http.StatusUnprocessableEntity,
		},
		{
			name:           "Invalid expires",
			snippetTitle:   validTitle,
			snippetContent: validContent,
			snippetExpires: "0",
			csrfToken:      validCSRFToken,
			wantCode:       http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("title", tt.snippetTitle)
			form.Add("content", tt.snippetContent)
			form.Add("expires", tt.snippetExpires)
			form.Add("csrf_token", validCSRFToken)

			code, header, _ := ts.postForm(t, "/snippet/create", form)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantLocation != "" {
				assert.StringContains(t, header.Get("Location"), tt.wantLocation)
			}
		})
	}
}
