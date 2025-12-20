package pkg

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SSEConn represents a Server-Sent Events (SSE) connection.
type SSEConn struct {
	ctx      *gin.Context
	hbTicker *time.Ticker
	done     chan struct{}
}

// NewSSEConn creates a new SSEConn instance from a Gin request and a heartbeat duration.
func NewSSEConn(ctx *gin.Context, heartbeatDuration time.Duration) *SSEConn {
	return &SSEConn{
		ctx:      ctx,
		done:     make(chan struct{}),
		hbTicker: time.NewTicker(heartbeatDuration),
	}
}

// SetupHeaders sets up the necessary headers for an SSE connection.
func (s *SSEConn) SetupHeaders() {
	s.ctx.Header("Content-Type", "text/event-stream; charset=utf-8")
	s.ctx.Header("Cache-Control", "no-cache")
	s.ctx.Header("Connection", "keep-alive")
	s.ctx.Header("X-Accel-Buffering", "no")
	s.ctx.Status(http.StatusOK)
}

// SendEvent sends an SSE event with the specified event name and data.
func (s *SSEConn) SendEvent(event string, data any) error {
	s.ctx.SSEvent(event, data)
	s.ctx.Writer.Flush()
	return nil
}

// StartHeartbeats starts sending heartbeat messages at regular intervals.
func (s *SSEConn) StartHeartbeats() {
	// No heartbeat configured
	if s.hbTicker == nil {
		return
	}

	heartbeatMessage := []byte(": heartbeat\n\n")
	go func() {
		for {
			select {
			case <-s.hbTicker.C:
				if _, err := s.ctx.Writer.Write(heartbeatMessage); err != nil {
					s.Close()
					return
				}
				s.ctx.Writer.Flush()
			case <-s.ctx.Request.Context().Done():
				// client disconnected
				s.Close()
				return
			case <-s.done:
				return
			}
		}
	}()
}

// Close closes the SSE connection and stops the heartbeat ticker.
func (s *SSEConn) Close() {
	select {
	case <-s.done:
		return
	default:
		close(s.done)
	}

	// stop heartbeat ticker, if it is running
	if s.hbTicker != nil {
		s.hbTicker.Stop()
	}
}
