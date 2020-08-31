package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/googleinterns/recomator/pkg/automation"
	"golang.org/x/oauth2"
)

// User contains fields describing authorized user:
// Email and GoogleService with required access to Google APIs.
type User struct {
	service automation.GoogleService
	email   string
}

// AuthorizationService authorizes the user using authentication code
type AuthorizationService interface {
	Authorize(authCode string) (User, error)
}

// authorizeRequest extracts the authorization code from Authorization header in request
// and uses it to return authorize user using authService.
func authorizeRequest(authService AuthorizationService, request *gin.Request) (User, err) {
	authCode := c.Request.Header["Authorization"]
	name := c.Query("name")
	if len(authCode) == 0 {
		return fmt.Errorf("Authorization code not specified")
	}
	return authService.Authorize(authCode[0])
}

func authHandler(c *gin.Context) {
	email := c.Query("login_hint")
	var authOptions []oauth2.AuthCodeOption

	if len(email) != 0 {
		authOptions = append(authOptions, oauth2.SetAuthURLParam("login_hint", email))
	}

	url := config.AuthCodeURL(c.Request.URL.Path, authOptions...)
	c.Redirect(http.StatusSeeOther, url)
	return
}
