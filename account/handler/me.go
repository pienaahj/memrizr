package handler

import (
	"log"
	"net/http"

	// j "github.com/pienaahj/memrizr/account/jsoncoder"
	"github.com/gin-gonic/gin"
	"github.com/pienaahj/memrizr/account/model"
	"github.com/pienaahj/memrizr/account/model/apperrors"
)

// Me handler calls services for getting
// a user's details
func (h *Handler) Me(c *gin.Context) {
	// A *model.User will eventually be added to context in middleware
	user, exists := c.Get("user")

	// This shouldn't happen, as our middleware ought to throw an error.
	// This is an extra safety measure
	// We'll extract this logic later as it will be common to all handler
	// methods which require a valid user
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		err := apperrors.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	//  get the user id from the retrieve user from context
	//  user is an interface value, so we need to assert the correct type
	uid := user.(*model.User).UID

	// gin.Context satisfies go's context.Context interface
	u, err := h.UserService.Get(c, uid)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", uid, err)
		e := apperrors.NewNotFound("user", uid.String())

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})
}

/*	own router
// handleMe handles the Me route by calling services for getting a user's detail
func (s *Handler) handleMe(c *) http.HandlerFunc {
	// thing := prepareThing()
	newString := "Hello, it's me"
	user, exists := c.
	    // This shouldn't happen, as our middleware ought to throw an error.
    // This is an extra safety measure
    // We'll extract this logic later as it will be common to all handler
    // methods which require a valid user
    if !exists {
        log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
        err := apperrors.NewInternal()
        c.JSON(err.Status(), gin.H{
            "error": err,
        })

        return
    }

  uid := user.(*model.User).UID
	return func(w http.ResponseWriter, r *http.Request) {
		// use thing
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		j.ToJSON(newString, w)
		log.Println("sent reply to webpage...")
	}
}
*/
