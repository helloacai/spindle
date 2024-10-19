package server

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/helloacai/spindle/pkg/thread"
	. "github.com/helloacai/spindle/pkg/util" // FromHex
)

const requestIDKey = "request-id"

func Run() error {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set(requestIDKey, uuid.New())
	})

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"health": "ok",
		})
	})
	r.GET("/thread/:uid/stream", ThreadStream)
	return r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

type Thread struct {
	UIDHex string `uri:"uid" binding:"required"`
}

func ThreadStream(c *gin.Context) {
	var t Thread
	if err := c.ShouldBindUri(&t); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	uid, err := FromHex(t.UIDHex)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ch, err := thread.Listen(c, uid, c.GetString(requestIDKey))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.Stream(func(w io.Writer) bool {
		if event, ok := <-ch; ok {
			c.SSEvent("event", event)
			return true
		}
		return false
	})
}
