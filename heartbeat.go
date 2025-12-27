package xrpl

import (
	"time"
)

// Heartbeat runner to send Pings periodically. If a Pong is received, it is
// handled by handlePong handler which further extends websocket connection's
// read and write deadline into the future.
func (c *Client) heartbeat() {
	defer c.wg.Done()
	// log.Println("INF: Heartbeat started")
	ticker := time.NewTicker(c.config.HeartbeatInterval)
	for {
		select {
		case <-c.heartbeatDone:
			ticker.Stop()
			// log.Println("ERR: Heartbeat stopped")
			return
		case t := <-ticker.C:
			c.Ping([]byte(t.String()))
		}
	}
}
