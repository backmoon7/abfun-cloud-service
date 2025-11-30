package api

import (
"bilibili-clone/internal/user/service"
"net/http"

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
Username string `json:"username"`
Password string `json:"password"`
}
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

if err := h.svc.Register(req.Username, req.Password); err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, gin.H{"message": "registered successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {
var req struct {
Username string `json:"username"`
Password string `json:"password"`
}
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

token, err := h.svc.Login(req.Username, req.Password)
if err != nil {
c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, gin.H{"token": token})
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

c.JSON(http.StatusOK, gin.H{"users": users})
}
