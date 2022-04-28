// routes.go

package main

import (
	"auth-manager/api"
	"auth-manager/config"
	"auth-manager/generates"
	"auth-manager/handlers"
	"auth-manager/middleware"
	"auth-manager/user"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
)

func initializeRoutes() {
	router.Static("/static/js", "static/js")
	router.Static("/static/img", "static/img")
	router.Static("/static/css", "static/css")

	// api endpoints

	// UI endpoints
	router.GET("/", middleware.EnsureLoggedIn(), handlers.RedirectIndexHandler)
	router.GET("/login", handlers.LoginGetHandler)
	router.POST("/login", func(c *gin.Context) {
		store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if c.Request.Form == nil {
			if err := c.Request.ParseForm(); err != nil {
				http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if config.Config.LDAP.Enabled {
			// cfg := ldap.Config{
			// 	BaseDN:       config.Config.LDAP.BaseDN,
			// 	BindDN:       config.Config.LDAP.BindDN,
			// 	Port:         strconv.Itoa(config.Config.LDAP.Port),
			// 	Host:         config.Config.LDAP.Host,
			// 	BindPassword: config.Config.LDAP.BindPassword,
			// 	Filter:       config.Config.LDAP.Filter,
			// }

			// r, _ := http.NewRequest("GET", "/", nil)
			// r.SetBasicAuth(c.Request.Form.Get("username"), c.Request.Form.Get("password"))

			// usr, err := ldap.New(&cfg).Authenticate(r.Context(), r)
			// if err == nil {
			// 	username := c.Request.Form.Get("username")
			// 	password := c.Request.Form.Get("password")
			// 	givenName := usr.Extensions()[config.Config.LDAP.Keys.GivenName][0]
			// 	familyName := usr.Extensions()[config.Config.LDAP.Keys.FamilyName][0]
			// 	email := usr.Extensions()[config.Config.LDAP.Keys.Email][0]
			// 	groups := strings.Split(usr.Extensions()[config.Config.LDAP.Keys.Groups][0], ",")
			// 	roles := strings.Split(usr.Extensions()[config.Config.LDAP.Keys.Roles][0], ",")

			// 	isAvailable := user.IsUsernameAvailable(usr.ID())
			// 	if isAvailable {
			// 		user.CreateNewUser(givenName, familyName, username, password, email, roles, groups)
			// 	} else {
			// 		tempUser, _ := user.GetUserByUsername(username)
			// 		tempUser.GivenName = givenName
			// 		tempUser.FamilyName = familyName
			// 		tempUser.Email = email
			// 		tempUser.Groups = strings.Join(groups[:], ",")
			// 		tempUser.Roles = strings.Join(roles[:], ",")
			// 		user.UpdateUserContents(tempUser.ID, *tempUser)
			// 	}
			// 	u, _ := user.GetUserByUsername(username)
			// 	store.Set("LoggedInUserID", u.ID)
			// 	store.Save()

			// 	c.Writer.Header().Set("Location", "/auth")
			// 	c.Writer.WriteHeader(http.StatusFound)
			// } else {
			// 	isValid, userIdent := user.IsUserValid(c.Request.Form.Get("username"), c.Request.Form.Get("password"))
			// 	if isValid {
			// 		store.Set("LoggedInUserID", userIdent)
			// 		store.Save()

			// 		log.Printf("LOGIN QUERY PARAMS: %v", c.Request.URL.Query())
			// 		c.Writer.Header().Set("Location", "/auth")
			// 		c.Writer.WriteHeader(http.StatusFound)
			// 	} else {
			// 		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			// 			"ErrorTitle":   "Login Failed",
			// 			"ErrorMessage": "Invalid credentials provided"})
			// 	}
			// }
		} else {
			isValid, userIdent := user.IsUserValid(c.Request.Form.Get("username"), c.Request.Form.Get("password"))
			if isValid {
				store.Set("LoggedInUserID", userIdent)
				store.Save()

				c.Writer.Header().Set("Location", "/auth")
				c.Writer.WriteHeader(http.StatusFound)
			} else {
				c.HTML(http.StatusBadRequest, "login.html", gin.H{
					"ErrorTitle":   "Login Failed",
					"ErrorMessage": "Invalid credentials provided"})
			}
		}
	})
	router.GET("/auth", handlers.AuthGetHandler)

	router.GET("/.well-known/openid-configuration/", func(c *gin.Context) {
		response := gin.H{
			"issuer":                                config.Config.URL,
			"authorization_endpoint":                config.Config.URL + "/oauth/authorize",
			"token_endpoint":                        config.Config.URL + "/oauth/token",
			"userinfo_endpoint":                     config.Config.URL + "/oauth/userinfo",
			"end_session_endpoint":                  config.Config.URL + "/oauth/logout",
			"introspection_endpoint":                config.Config.URL + "/oauth/introspect",
			"jwks_uri":                              config.Config.URL + "/JWKS",
			"subject_types_supported":               []string{"public"},
			"token_endpoint_auth_methods_supported": []string{"client_secret_post"},
		}
		c.JSON(200, response)
	})

	router.GET("/JWKS", func(c *gin.Context) {
		jsonData := JWKS
		c.Data(http.StatusOK, "application/json", jsonData)
	})

	uRoutes := router.Group("/u")
	{
		uRoutes.GET("/login", middleware.EnsureNotLoggedIn(), handlers.LocalLoginHandler)
		uRoutes.POST("/login", user.PerformLogin)
		uRoutes.GET("/logout", middleware.EnsureLoggedIn(), user.Logout)
	}

	uiRoutes := router.Group("/ui")
	{
		uiRoutes.GET("/index", middleware.EnsureLoggedIn(), handlers.IndexHandler)
		uiRoutes.GET("/create", middleware.EnsureLoggedIn(), middleware.EnsureGroupAllowed("admin"), handlers.CreateHandler)
		uiRoutes.GET("/delete", middleware.EnsureLoggedIn(), middleware.EnsureGroupAllowed("admin"), handlers.DeleteHandler)
		uiRoutes.GET("/edit", middleware.EnsureLoggedIn(), handlers.EditIndexHandler)
		uiRoutes.GET("/edit/*id", middleware.EnsureLoggedIn(), middleware.EnsureGroupAllowed("admin"), handlers.EditHandler)
	}

	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/user", middleware.EnsureLoggedInAbort(), api.UserGetAllHandler)
		apiRoutes.GET("/user/:id", middleware.EnsureLoggedInAbort(), api.UserGetByIdHandler)
		apiRoutes.DELETE("/user/:id", middleware.EnsureLoggedInAbort(), api.UserDeleteHandler)
		apiRoutes.POST("/user/:id", middleware.EnsureLoggedInAbort(), api.UserUpdateHandler)
		apiRoutes.POST("/user", middleware.EnsureLoggedInAbort(), api.UserCreateHandler)
	}

	oauthRoutes := router.Group("/oauth")
	{
		oauthRoutes.POST("/introspect", func(c *gin.Context) {
			c.Request.ParseForm()

			client, err := ClientStore.GetByID(context.Background(), c.PostForm("client_id"))
			if err != nil {
				log.Println("Invalid client ID")
				c.AbortWithStatus(500)
				return
			}
			if client.GetSecret() != c.PostForm("client_secret") {
				log.Println("Invalid client secret")
				c.AbortWithStatus(500)
				return
			}
			// Parse and verify jwt access token
			token, err := jwt.ParseWithClaims(c.PostForm("token"), &generates.JWTAccessClaims{}, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("parse error")
				}
				return Privkey, nil
			})
			if err != nil {
				c.AbortWithStatus(500)
				return
			}

			claims, ok := token.Claims.(*generates.JWTAccessClaims)
			if !ok || !token.Valid {
				log.Println("invalid token")
				c.AbortWithStatus(500)
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"given_name":  claims.GivenName,
				"family_name": claims.FamilyName,
				"email":       claims.Email,
				"groups":      claims.Groups,
				"roles":       claims.Roles,
				"memberOf":    claims.Groups,
				"exp":         claims.ExpiresAt,
				"sub":         claims.Subject,
				"aud":         claims.Audience,
			})
		})

		oauthRoutes.GET("/authorize", func(c *gin.Context) {

			store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
			if err != nil {
				http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
				return
			}

			var form url.Values
			if v, ok := store.Get("ReturnUri"); ok {
				form = v.(url.Values)
			}
			c.Request.Form = form

			store.Delete("ReturnUri")
			store.Save()

			err = Srv.HandleAuthorizeRequest(c.Writer, c.Request)
			if err != nil {
				http.Error(c.Writer, err.Error(), http.StatusBadRequest)
			}
		})

		oauthRoutes.POST("/authorize", func(c *gin.Context) {
			c.Request.ParseForm()
			store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
			if err != nil {
				http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
				return
			}

			var form url.Values
			if v, ok := store.Get("ReturnUri"); ok {
				form = v.(url.Values)
			}
			c.Request.Form = form

			store.Delete("ReturnUri")
			store.Save()

			err = Srv.HandleAuthorizeRequest(c.Writer, c.Request)
			if err != nil {
				http.Error(c.Writer, err.Error(), http.StatusBadRequest)
			}
		})

		oauthRoutes.POST("/token", func(c *gin.Context) {
			c.Request.ParseForm()
			err := Srv.HandleTokenRequest(c.Writer, c.Request)
			if err != nil {
				http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			}
		})

		oauthRoutes.GET("/userinfo", api.UserInfo)
	}
}
