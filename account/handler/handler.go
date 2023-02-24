package handler

import (
	"net/http"
	"time"

	"github.com/pienaahj/memrizr/account/handler/middleware"
	"github.com/pienaahj/memrizr/account/model"
	"github.com/pienaahj/memrizr/account/model/apperrors"

	// rt "github.com/pienaahj/memrizr/account/router"
	"github.com/gin-gonic/gin"
)

// Handler struct is an object representing the service and holds all
// Shared dependacies - these should be fields of the structure
type Handler struct {
	//logger *someLogger
	// db     *someDatabase
	// router *someRouter
	// R*rt.Router
	// email  EmailSender
	// user handler
	UserService model.UserService
	// token service
	TokenService model.TokenService
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	// R           *rt.Router
	R               *gin.Engine
	UserService     model.UserService
	TokenService    model.TokenService
	BaseURL         string
	TimeoutDuration time.Duration
}

//	NewRouter creates a new Router initializes the handler with required injected services along with http routes
//
// Does not return as it deals directly with a reference to the router engine
// func NewHandler(c *Config) *Handler {
func NewHandler(c *Config) {
	// Create a handler (which will later have injected services)
	h := &Handler{
		// R:           c.R,
		UserService:  c.UserService,
		TokenService: c.TokenService,
	}
	// Create a group, or base url for all routes
	g := c.R.Group(c.BaseURL)

	if gin.Mode() != gin.TestMode {
		g.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))
		g.GET("/me", middleware.AuthUser(h.TokenService), h.Me)
		g.POST("/signout", middleware.AuthUser(h.TokenService), h.Signout)
		g.PUT("/details", middleware.AuthUser(h.TokenService), h.Details)
	} else {
		g.GET("/me", h.Me)
		g.POST("/signout", h.Signout)
		g.PUT("/details", h.Details)
	}

	g.POST("/signup", h.Signup)
	g.POST("/signin", h.Signin)
	g.POST("/tokens", h.Tokens)
	g.POST("/image", h.Image)
	g.DELETE("/image", h.DeleteImage)
}

// Image handler
func (h *Handler) Image(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's image",
	})
}

// DeleteImage handler
func (h *Handler) DeleteImage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's deleteImage",
	})
}

// Create an account group
// g := c.R.Group("/api/account")

// g.GET("/", func(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{
// 		"hello": "space persons",
// 	})
// })
//
// return h
// }
/*    own router
func (hd *Handler) Routes() {
	//  Define the URL group
	group := os.Getenv("../ACCOUNT_API_URL")
	fmt.Println("Environment: ", group)
	hd.R.Route(http.MethodGet, group, hd.handleAccount())
	hd.R.Route(http.MethodGet, group+"/me", hd.handleMe())
	hd.R.Route(http.MethodPost, group+"/signup", hd.handleSignup())
	hd.R.Route(http.MethodPost, group+"/signin", hd.handleSignin())
	hd.R.Route(http.MethodPost, group+"/signout", hd.handleSignout())
	hd.R.Route(http.MethodPost, group+"/tokens", hd.handleTokens())
	hd.R.Route(http.MethodDelete, group+"/images", hd.handleImages())
	hd.R.Route(http.MethodPut, group+"/details", hd.handleDetails())

}
*/

/*   own handlers
// handleAccount handles route /api/account
func (hd *Handler) handleAccount() http.HandlerFunc {
	log.Println("handling route /api/account...")
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			return
		}
		// define the reply
		type reply struct {
			Hello string `json:"hello"`
		}
		// create the reply value
		resp := reply{
			"peoples",
		}
		// encode to json
		// add the header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
		log.Println("sent reply to webpage...")
	}

}

// handleMe handles the Signup route
func (s *Handler) handleSignup() http.HandlerFunc {
	// thing := prepareThing()
	newString := "Hello, it's signup"
	return func(w http.ResponseWriter, r *http.Request) {
		// use thing
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		j.ToJSON(newString, w)
		log.Println("sent reply to webpage...")
	}
}

// handleMe handles the Signin route
func (s *Handler) handleSignin() http.HandlerFunc {
	// thing := prepareThing()
	newString := "Hello, it's signin"
	return func(w http.ResponseWriter, r *http.Request) {
		// use thing
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		j.ToJSON(newString, w)
		log.Println("sent reply to webpage...")
	}
}

// handleMe handles the Signout route
func (s *Handler) handleSignout() http.HandlerFunc {
	/// thing := prepareThing()
	newString := "Hello, it's signout"
	return func(w http.ResponseWriter, r *http.Request) {
		// use thing
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		j.ToJSON(newString, w)
		log.Println("sent reply to webpage...")
	}
}

// handleMe handles the Tokens route
func (s *Handler) handleTokens() http.HandlerFunc {
	// thing := prepareThing()
	newString := "Hello, it's tokens"
	return func(w http.ResponseWriter, r *http.Request) {
		// use thing
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		j.ToJSON(newString, w)
		log.Println("sent reply to webpage...")
	}
}

// handleMe handles the Images route
func (s *Handler) handleImages() http.HandlerFunc {
	// thing := prepareThing()
	newString := "Hello, it's images"
	return func(w http.ResponseWriter, r *http.Request) {
		// use thing
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		j.ToJSON(newString, w)
		log.Println("sent reply to webpage...")
	}
}

// handleMe handles the SDetails route
func (s *Handler) handleDetails() http.HandlerFunc {
	// thing := prepareThing()
	newString := "Hello, it's details"
	return func(w http.ResponseWriter, r *http.Request) {
		// use thing
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		j.ToJSON(newString, w)
		log.Println("sent reply to webpage...")
	}
}
*/
