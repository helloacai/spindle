package server

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/helloacai/spindle/pkg/log"
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
	r.GET("/thread/:uid/context/stream", StreamThreadContext)
	//r.PUT("/thread/:uid/context", PutThreadContext)
	return r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

type Thread struct {
	UIDHex string `uri:"uid" binding:"required"`
}

func StreamThreadContext(c *gin.Context) {
	c.Header("access-control-allow-origin", "*")
	requestID := c.GetString(requestIDKey)

	var t Thread
	if err := c.ShouldBindUri(&t); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	logger := log.With().Str("request_id", requestID).Str("thread", t.UIDHex).Logger()

	logger.Debug().Msg("stream thread")

	uid, err := FromHex(t.UIDHex)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ch, err := thread.Listen(c, uid, requestID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.Stream(func(w io.Writer) bool {
		if event, ok := <-ch; ok {
			c.SSEvent("event", event)
			logger.Debug().Msg("closing listener: thread complete in api")
			return event.Type != thread.EntryType_Complete
		}
		return false
	})
}

//type Entry struct {
//	Type       thread.EntryType `json:"type"`
//	Originator string           `json:"originator"`
//	Message    string           `json:"message"`
//}
//
//func PutThreadContext(c *gin.Context) {
//	var t Thread
//	if err := c.ShouldBindUri(&t); err != nil {
//		c.JSON(400, gin.H{"error": err.Error()})
//		return
//	}
//
//	uid, err := FromHex(t.UIDHex)
//	if err != nil {
//		c.JSON(400, gin.H{"error": err.Error()})
//		return
//	}
//
//	var json Entry
//	if err := c.ShouldBindJSON(&json); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	originator, err := FromHex(json.Originator)
//	if err != nil {
//		c.JSON(400, gin.H{"error": err.Error()})
//		return
//	}
//
//	if err := thread.Append(uid, json.Type, originator, json.Message); err != nil {
//		c.JSON(400, gin.H{"error": err.Error()})
//		return
//	}
//}
