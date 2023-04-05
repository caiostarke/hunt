package http

import (
	"fmt"
	"hunt/internal/repository"
	db "hunt/internal/repository/mysql"
	"hunt/pkg/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-mysql/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct {
	repo *db.Repository
}

func New(repo *db.Repository) *Handler {
	return &Handler{repo}
}

func (h *Handler) Login(c *gin.Context) {
	var json struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if json.Email == "" || json.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "fields cannot be empty",
		})
		return
	}

	u, err := h.repo.GetUserByEmail(c.Request.Context(), json.Email)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	u.Password.Plaintext = json.Password
	if err := u.Password.CheckPasswordHash(); err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email or Password are wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Logged IN",
	})

}

func (h *Handler) CreateUser(c *gin.Context) {
	var userJSON model.UserJSON
	if err := c.ShouldBindJSON(&userJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if userJSON.Email == "" || userJSON.Name == "" || userJSON.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "fields cannot be empty",
		})
		return
	}

	var user model.User

	user.ID = uuid.New()
	user.Level = 0
	user.CreatedAt = time.Now()
	user.ProfilePicture = userJSON.ProfilePicture
	user.Name = userJSON.Name
	user.Email = userJSON.Email

	if err := user.Password.Set(userJSON.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "try again or later",
		})
		return
	}

	if err := h.repo.CreateUser(c.Request.Context(), &user); err != nil {
		// erros package is responsible for handling common mysql errors and this function makes type assertion
		_, err = errors.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userID": user.ID,
	})
}

// StartMatch receives a JWT Token (testing only with userID), validate it and then add a user into a queue
func (h *Handler) StartMatch(c *gin.Context) {
	var json model.UserID

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no ID provided",
		})
		return
	}

	if json.ID.String() == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "empty ID",
		})
		return
	}

	u, err := h.repo.GetUser(c.Request.Context(), json.ID)
	if err != nil && err == repository.ErrNotFound {
		c.Status(http.StatusNotFound)
		return
	} else if err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Print(err)
		return
	}

	u.ID = json.ID

	logger, _ := zap.NewProduction()
	queueID := model.SearchQueue(u, logger)

	c.JSON(http.StatusOK, gin.H{
		"queue": queueID,
	})
}

func (h *Handler) QueueMonitor(c *gin.Context) {
	aQ := model.AvailableQueues
	availableQueuesLen := len(aQ)
	usersInAvailableQueues := model.GetUsersFromAvailableQueues()

	sQ := len(model.QueuesStarted)

	c.JSON(http.StatusOK, gin.H{
		"Available Queues":          availableQueuesLen,
		"Users in Available Queues": usersInAvailableQueues,
		"StartedQueues":             sQ,
	})
}
