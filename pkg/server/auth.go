/*
Copyright 2020 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/googleinterns/recomator/pkg/automation"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
)

const tokenExpiry = 24 * time.Hour

// User contains fields describing authorized user:
// Email and GoogleService with required access to Google APIs.
type User struct {
	service automation.GoogleService
	email   string
}

// AuthorizationService creates the user using authentication code and returns idToken.
// Returns authorized user for this idToken.
type AuthorizationService interface {
	// Returns redirect URL to login page used for authentication
	AuthCodeURL(options ...oauth2.AuthCodeOption) string
	// Returns idToken that should be used to authorize
	CreateUser(authCode string) (string, error)
	// Verify checks that token is valid, not expired and issued by our app and returns user email
	Verify(token string) (string, error)
	// Returns GoogleService that should be used to make requests to Google APIs
	GetUser(email string) (User, bool)
}

type authorizationService struct {
	verifier            *oidc.IDTokenVerifier
	tokenExpirationTime time.Duration
	mutex               sync.Mutex
	config              oauth2.Config
	services            map[string]automation.GoogleService // key is email of the user
}

// NewAuthorizationService creates new AuthorizationService to access GoogleAPIs
func NewAuthorizationService(config oauth2.Config) (AuthorizationService, error) {
	provider, err := oidc.NewProvider(oauth2.NoContext, "https://accounts.google.com")
	if err != nil {
		return nil, err
	}
	authService := &authorizationService{tokenExpirationTime: tokenExpiry, services: make(map[string]automation.GoogleService)}
	authService.verifier = provider.Verifier(&oidc.Config{ClientID: config.ClientID, SkipExpiryCheck: true})
	authService.config = config
	return authService, nil
}

// AuthCodeURL returns google login page url.
func (s *authorizationService) AuthCodeURL(options ...oauth2.AuthCodeOption) string {
	return s.config.AuthCodeURL("", options...)
}

// Returns idToken that should be used for authorization later.
func (s *authorizationService) CreateUser(authCode string) (string, error) {
	token, err := s.config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		return "", err
	}
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", fmt.Errorf("No valid id token where given. Casting to string failed")
	}

	service, err := automation.NewGoogleService(oauth2.NoContext, &s.config, token)
	if err != nil {
		return "", err
	}

	email, err := s.Verify(idToken)
	if err != nil {
		return "", err
	}

	s.mutex.Lock()
	s.services[email] = service
	s.mutex.Unlock()
	return idToken, nil
}

func (s *authorizationService) GetUser(email string) (User, bool) {
	s.mutex.Lock()
	service, ok := s.services[email]
	s.mutex.Unlock()
	if ok {
		return User{service, email}, true
	}
	return User{}, false
}

// Verifies idToken and returns email if everything suceeded.
func (s *authorizationService) Verify(rawToken string) (string, error) {
	idToken, err := s.verifier.Verify(oauth2.NoContext, rawToken)
	if err != nil {
		return "", err
	}
	if time.Since(time.Time(idToken.IssuedAt)) > s.tokenExpirationTime {
		return "", &googleapi.Error{Code: http.StatusUnauthorized, Message: "Token expired"}
	}
	var claims struct {
		Email string `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return "", fmt.Errorf("Extracting email failed:" + err.Error())
	}
	return claims.Email, nil
}

type idToken struct {
	Token string `json:"token"`
}

func getAuthHandler(service *SharedService) func(c *gin.Context) {
	return func(c *gin.Context) {
		authCode := c.Query("code")
		token, err := service.auth.CreateUser(authCode)
		if err != nil {
			sendError(c, err, http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, idToken{Token: token})
	}
}

func authorize(authService AuthorizationService, idToken string) (User, error) {
	email, err := authService.Verify(idToken)
	if err != nil {
		return User{}, err
	}

	user, ok := authService.GetUser(email)

	if !ok {
		return User{}, &googleapi.Error{Message: fmt.Sprintf("User with %s email not found", email),
			Code: http.StatusNotFound}
	}

	return user, nil
}

// authorizeRequest extracts the token from Authorization header in request
// and uses it to return authorized user using authService.
func authorizeRequest(authService AuthorizationService, request *http.Request) (User, error) {
	bearToken := request.Header["Authorization"]
	if len(bearToken) != 0 {

		strArr := strings.Split(bearToken[0], " ")
		if len(strArr) == 2 && strArr[0] == "Bearer" {
			return authorize(authService, strArr[1])
		}
	}
	return User{}, &googleapi.Error{Code: http.StatusBadRequest, Message: "Authorization header not in the form 'Bearer <token>'"}

}

// redirects to google for login, login_hint query parameter(user's email) might be specified for faster login.
func redirectHandler(service *SharedService) func(c *gin.Context) {
	return func(c *gin.Context) {
		email := c.Query("login_hint")
		authOptions := []oauth2.AuthCodeOption{oauth2.AccessTypeOffline, oauth2.ApprovalForce}

		if len(email) != 0 {
			authOptions = append(authOptions, oauth2.SetAuthURLParam("login_hint", email))
		}

		url := service.auth.AuthCodeURL(authOptions...)
		c.Redirect(http.StatusSeeOther, url)
		return
	}
}
