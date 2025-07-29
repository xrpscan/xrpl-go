package xrpl

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func (c *Client) handlePong(message string) error {
	// log.Println("PONG:", message)
	c.connection.SetReadDeadline(time.Now().Add(c.config.ReadTimeout * time.Second))
	c.connection.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout * time.Second))
	return nil
}

func (c *Client) handleResponse() error {
	for {
		select {
		case <-c.handlerDone:
			return nil
		default:
			if c.closed {
				return nil
			}
		}

		messageType, message, err := c.connection.ReadMessage()
		if err != nil {
			log.Println("WS read error:", err)
			c.Reconnect()
			return nil
		}

		switch messageType {
		case websocket.CloseMessage:
			log.Println("WS websocket.CloseMessage received")
			return nil
		case websocket.TextMessage:
			c.resolveStream(message)
		case websocket.BinaryMessage:
		default:
		}
	}
}

func (c *Client) resolveStream(message []byte) {
	var m BaseResponse
	if err := json.Unmarshal(message, &m); err != nil {
		log.Println("json.Unmarshal error: ", err)
	}

	switch m["type"] {
	case StreamResponseType(StreamTypeLedger):
		c.StreamLedger <- message

	case StreamResponseType(StreamTypeTransaction):
		c.StreamTransaction <- message

	case StreamResponseType(StreamTypeValidations):
		c.StreamValidation <- message

	case StreamResponseType(StreamTypeManifests):
		c.StreamManifest <- message

	case StreamResponseType(StreamTypePeerStatus):
		c.StreamPeerStatus <- message

	case StreamResponseType(StreamTypeConsensus):
		c.StreamConsensus <- message

	case StreamResponseType(StreamTypePathFind):
		c.StreamPathFind <- message

	case StreamResponseType(StreamTypeServer):
		c.StreamServer <- message

	case StreamResponseType(StreamTypeResponse):
		requestId := fmt.Sprintf("%v", m["id"])
		c.mutex.Lock()
		ch, ok := c.requestQueue[requestId]
		if ok {
			ch <- m
			delete(c.requestQueue, requestId)
			close(ch)
		}
		c.mutex.Unlock()

	default:
		c.StreamDefault <- message
	}
}
