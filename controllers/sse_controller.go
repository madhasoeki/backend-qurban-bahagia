package controllers

import (
	"io"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	clients   = make(map[chan any]bool)
	clientsMu sync.Mutex
	Broadcast = make(chan any)
)

func init() {
	go func() {
		for msg := range Broadcast {
			clientsMu.Lock()
			for client := range clients {
				select {
				case client <- msg:
				default:
					// Drop message for slow clients to prevent blocking
				}
			}
			clientsMu.Unlock()
		}
	}()
}

func SSEHandler(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// Buffered channel prevents slow consumers from blocking the broadcaster
	clientChan := make(chan any, 16)

	clientsMu.Lock()
	clients[clientChan] = true
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, clientChan)
		close(clientChan)
		clientsMu.Unlock()
	}()

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-clientChan; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})
}

// broadcastHewanUpdate sends the updated hewan data and recalculated dashboard
// summary to all connected SSE clients. Used by all pos/hewan mutation endpoints.
func broadcastHewanUpdate(hewan any) {
	Broadcast <- gin.H{
		"action": "UPDATE_HEWAN",
		"data":   hewan,
	}

	go func() {
		summary := calculateDashboardSummary()
		Broadcast <- gin.H{
			"action": "UPDATE_DASHBOARD",
			"data":   summary,
		}
	}()
}
