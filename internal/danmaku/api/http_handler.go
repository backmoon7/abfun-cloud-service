package api

import (
	"bilibili-clone/internal/danmaku/service"
	"bilibili-clone/internal/danmaku/ws"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type DanmakuHandler struct {
	hub *ws.Hub
	svc *service.DanmakuService
}

func NewDanmakuHandler(hub *ws.Hub) *DanmakuHandler {
	return &DanmakuHandler{
		hub: hub,
		svc: service.NewDanmakuService(),
	}
}

func (h *DanmakuHandler) GetDanmakus(c *gin.Context) {
	vidStr := c.Query("video_id")
	vid, _ := strconv.Atoi(vidStr)

	danmakus, err := h.svc.GetDanmakus(uint(vid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": danmakus})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *DanmakuHandler) ServeWS(c *gin.Context) {
	vidStr := c.Query("video_id")
	vid, err := strconv.ParseUint(vidStr, 10, 64)
	if err != nil || vid == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video_id"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &ws.Client{Hub: h.hub, Conn: conn, Send: make(chan []byte, 256), VideoID: uint(vid)}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
