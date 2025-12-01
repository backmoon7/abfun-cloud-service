package api

import (
	"bilibili-clone/internal/user/service"
	"net/http"
	"strings"
        "strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		svc: service.NewUserService(),
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(req.Email) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

        // Basic email format validation
        if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
                c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
                return
        }


        // Password validation
        if len(req.Password) < 6 {
                c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 6 characters"})
                return
        }
        if len(req.Password) > 100 {
                c.JSON(http.StatusBadRequest, gin.H{"error": "password too long"})
                return
        }
	if err := h.svc.Register(req.Email, req.Password, req.Nickname); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registered successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(req.Email) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	token, user, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
		},
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateUser(userID.(uint), req.Nickname, req.Avatar); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

func (h *UserHandler) GetUsersBatch(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := h.svc.GetUsersByIDs(req.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *UserHandler) GetMe(c *gin.Context) {
        userID, exists := c.Get("user_id")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
                return
        }

        user, err := h.svc.GetUserByID(userID.(uint))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "id":       user.ID,
                "email":    user.Email,
                "nickname": user.Nickname,
                "avatar":   user.Avatar,
        })
}

func (h *UserHandler) GetUsers(c *gin.Context) {
        idsStr := c.Query("ids")
        if idsStr == "" {
                c.JSON(http.StatusBadRequest, gin.H{"error": "ids parameter is required"})
                return
        }

        idStrs := strings.Split(idsStr, ",")
        var ids []uint
        for _, idStr := range idStrs {
                if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32); err == nil {
                        ids = append(ids, uint(id))
                }
        }

        users, err := h.svc.GetUsersByIDs(ids)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": users})
}
