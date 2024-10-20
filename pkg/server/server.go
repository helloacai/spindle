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
		c.Set(requestIDKey, uuid.New().String())
	})

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"health": "ok",
		})
	})
	r.GET("/thread/:uid/context/stream", StreamThreadContext)
	r.GET("/thread/:uid", GetThread)
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
			if event.Type == thread.EntryType_Debug {
				logger.Debug().Msg("skipping debug entry")
				return true
			}
			c.SSEvent("event", event)
			if event.Type != thread.EntryType_Complete {
				return true
			}

			logger.Debug().Msg("ending listener: thread complete in api")
			return false
		}
		logger.Debug().Msg("listener closed")
		return false
	})
}

func GetThread(c *gin.Context) {
	var param Thread
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	uid, err := FromHex(param.UIDHex)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	t, err := thread.Get(c, uid)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, &t)
}
