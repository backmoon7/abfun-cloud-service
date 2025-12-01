package worker

import (
	"bilibili-clone/internal/danmaku/ws"
	"encoding/json"
	"log"
)

type DanmakuConsumer struct {
	hub *ws.Hub
}

func NewDanmakuConsumer(hub *ws.Hub) *DanmakuConsumer {
	return &DanmakuConsumer{hub: hub}
}

type mqDanmakuPayload struct {
	VideoID uint            `json:"video_id"`
	Payload json.RawMessage `json:"payload"`
	Data    json.RawMessage `json:"data"`
	Message json.RawMessage `json:"message"`
}

func (c *DanmakuConsumer) Handle(msg []byte) {
	var payload mqDanmakuPayload
	if err := json.Unmarshal(msg, &payload); err != nil {
		log.Printf("invalid danmaku payload: %v", err)
		return
	}

	body := payload.Payload
	if len(body) == 0 {
		switch {
		case len(payload.Data) > 0:
			body = payload.Data
		case len(payload.Message) > 0:
			body = payload.Message
		default:
			body = msg
		}
	}

	if payload.VideoID == 0 {
		log.Printf("danmaku payload missing video_id")
		return
	}

	select {
	case c.hub.Broadcast <- ws.RoomMessage{VideoID: payload.VideoID, Payload: body}:
	default:
		// Drop message if broadcast channel is saturated to protect the worker
	}
}
