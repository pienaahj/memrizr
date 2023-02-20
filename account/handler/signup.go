package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pienaahj/memrizr/account/model"
	"github.com/pienaahj/memrizr/account/model/apperrors"
)

// signupReq is not exported, hence the lowercase name
// it is used for validation and json marshalling
// The documentation lists built-in validators and also describes things
// like validation errors and custom validators. Note that in the documentation,
// the struct tag validate is used, whereas in Gin we must use the struct tag binding.
type signupReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

// Signup handler
func (h *Handler) Signup(c *gin.Context) {
	// define a variable to which we'll bind incoming
	// json body, {email, password}
	var req signupReq

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}
	u := &model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	// use the Request Context
	ctx := c.Request.Context()
	err := h.UserService.Signup(ctx, u)

	if err != nil {
		log.Printf("Failed to sign up user: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// create token pair as strings
	tokens, err := h.TokenService.NewPairFromUser(ctx, u, "")

	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())

		// may eventually implement rollback logic here
		// meaning, if we fail to create tokens after creating a user,
		// we make sure to clear/delete the created user in the database

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"tokens": tokens,
	})
}