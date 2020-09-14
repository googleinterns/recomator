package main

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
	// Returns idToken that should be used to authorize
	CreateUser(authCode string) (string, error)
	//
	Verify(token string) (string, error)
	// Returns GoogleService that should be used to make requests to Google APIs
	GetUser(email string) (User, bool)
}

type authorizationService struct {
	verifier            *oidc.IDTokenVerifier
	tokenExpirationTime time.Duration
	mutex               sync.Mutex
	services            map[string]automation.GoogleService // key is email of the user
}

// NewAuthorizationService creates new AuthorizationService to access GoogleAPIs
func NewAuthorizationService() (AuthorizationService, error) {
	provider, err := oidc.NewProvider(oauth2.NoContext, "https://accounts.google.com")
	if err != nil {
		return nil, err
	}
	authService := &authorizationService{tokenExpirationTime: tokenExpiry, services: make(map[string]automation.GoogleService)}
	authService.verifier = provider.Verifier(&oidc.Config{ClientID: config.ClientID, SkipExpiryCheck: true})
	return authService, nil
}

// Returns idToken that should be used for authorization later.
func (s *authorizationService) CreateUser(authCode string) (string, error) {
	token, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		return "", err
	}
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", fmt.Errorf("No valid id token where given. Casting to string failed")
	}

	service, err := automation.NewGoogleService(oauth2.NoContext, &config, token)
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

func getAuthHandler(service *SharedService) func(c *gin.Context) {
	return func(c *gin.Context) {
		authCode := c.Query("code")
		token, err := service.auth.CreateUser(authCode)
		if err != nil {
			sendError(c, err, http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
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
func redirectHandler(c *gin.Context) {
	email := c.Query("login_hint")
	authOptions := []oauth2.AuthCodeOption{oauth2.AccessTypeOffline, oauth2.ApprovalForce}

	if len(email) != 0 {
		authOptions = append(authOptions, oauth2.SetAuthURLParam("login_hint", email))
	}

	url := config.AuthCodeURL("", authOptions...)
	c.Redirect(http.StatusSeeOther, url)
	return
}
