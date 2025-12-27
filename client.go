package xrpl

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ClientConfig struct {
	URL                string
	Authorization      string
	Certificate        string
	FeeCushion         uint32
	Key                string
	MaxFeeXRP          uint64
	Passphrase         byte
	Proxy              byte
	ProxyAuthorization byte
	ReadTimeout        time.Duration // Default is 60 seconds
	WriteTimeout       time.Duration // Default is 60 seconds
	HeartbeatInterval  time.Duration // Default is 5 seconds
	QueueCapacity      int           // Default is 128
}

type Client struct {
	config              ClientConfig
	connection          *websocket.Conn
	heartbeatDone       chan bool
	handlerDone         chan bool
	closed              bool
	mutex               sync.Mutex
	response            *http.Response
	StreamLedger        chan []byte
	StreamTransaction   chan []byte
	StreamValidation    chan []byte
	StreamManifest      chan []byte
	StreamPeerStatus    chan []byte
	StreamConsensus     chan []byte
	StreamPathFind      chan []byte
	StreamServer        chan []byte
	StreamDefault       chan []byte
	StreamSubscriptions map[string]bool
	requestQueue        map[string](chan<- BaseResponse)
	nextId              int
	err                 error
}

func (config *ClientConfig) Validate() error {
	if len(config.URL) == 0 {
		return errors.New("cannot create a new connection with an empty URL")
	}

	if config.ReadTimeout < 0*time.Second || config.ReadTimeout >= 1*time.Hour {
		return fmt.Errorf("connection read timeout out of bounds: %d", config.ReadTimeout)
	}
	if config.WriteTimeout < 0*time.Second || config.WriteTimeout >= 1*time.Hour {
		return fmt.Errorf("connection write timeout out of bounds: %d", config.WriteTimeout)
	}
	if config.HeartbeatInterval < 0*time.Second || config.HeartbeatInterval >= 1*time.Hour {
		return fmt.Errorf("connection heartbeat interval out of bounds: %d", config.HeartbeatInterval)
	}

	return nil
}

func NewClient(config ClientConfig) *Client {
	if config.ReadTimeout == 0*time.Second {
		config.ReadTimeout = 60 * time.Second
	}
	if config.WriteTimeout == 0*time.Second {
		config.WriteTimeout = 60 * time.Second
	}
	if config.HeartbeatInterval == 0*time.Second {
		config.HeartbeatInterval = 5 * time.Second
	}

	if config.QueueCapacity == 0 {
		config.QueueCapacity = 128
	}

	if err := config.Validate(); err != nil {
		panic(err)
	}

	client := &Client{
		config:              config,
		heartbeatDone:       make(chan bool),
		handlerDone:         make(chan bool),
		StreamLedger:        make(chan []byte, config.QueueCapacity),
		StreamTransaction:   make(chan []byte, config.QueueCapacity),
		StreamValidation:    make(chan []byte, config.QueueCapacity),
		StreamManifest:      make(chan []byte, config.QueueCapacity),
		StreamPeerStatus:    make(chan []byte, config.QueueCapacity),
		StreamConsensus:     make(chan []byte, config.QueueCapacity),
		StreamPathFind:      make(chan []byte, config.QueueCapacity),
		StreamServer:        make(chan []byte, config.QueueCapacity),
		StreamDefault:       make(chan []byte, config.QueueCapacity),
		StreamSubscriptions: make(map[string]bool),
		requestQueue:        make(map[string](chan<- BaseResponse)),
		nextId:              0,
	}

	_, err := client.NewConnection()
	if err != nil {
		log.Println("WS connection error:", client.config.URL, err)
	}
	return client
}

func (c *Client) NewConnection() (*websocket.Conn, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	conn, r, err := websocket.DefaultDialer.Dial(c.config.URL, nil)
	if err != nil {
		c.err = err
		return nil, err
	}
	defer r.Body.Close()
	c.connection = conn
	c.response = r
	c.closed = false

	// Set connection handlers and heartbeat
	c.connection.SetReadDeadline(time.Now().Add(c.config.ReadTimeout))
	c.connection.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout))
	c.connection.SetPongHandler(c.handlePong)
	go c.handleResponse()
	go c.heartbeat()
	return c.connection, nil
}

func (c *Client) Reconnect() error {
	// Close old websocket connection
	c.Close()

	// Recreate stream channels
	c.mutex.Lock()
	c.StreamLedger = make(chan []byte, c.config.QueueCapacity)
	c.StreamTransaction = make(chan []byte, c.config.QueueCapacity)
	c.StreamValidation = make(chan []byte, c.config.QueueCapacity)
	c.StreamManifest = make(chan []byte, c.config.QueueCapacity)
	c.StreamPeerStatus = make(chan []byte, c.config.QueueCapacity)
	c.StreamConsensus = make(chan []byte, c.config.QueueCapacity)
	c.StreamPathFind = make(chan []byte, c.config.QueueCapacity)
	c.StreamServer = make(chan []byte, c.config.QueueCapacity)
	c.StreamDefault = make(chan []byte, c.config.QueueCapacity)
	c.mutex.Unlock()

	// Create a new websocket connection
	_, err := c.NewConnection()
	if err != nil {
		log.Println("WS reconnection error:", c.config.URL, err)
		return err
	}

	// Re-subscribe xrpl streams
	_, err = c.Subscribe(c.Subscriptions())
	if err != nil {
		log.Println("WS stream subscription error:", err)
	}
	return nil
}

func (c *Client) Ping(message []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// log.Println("PING:", string(message))
	newDeadline := time.Now().Add(c.config.WriteTimeout)
	if err := c.connection.WriteControl(websocket.PingMessage, message, newDeadline); err != nil {
		return err
	}
	return nil
}

// Returns incremental ID that may be used as request ID for websocket requests
func (c *Client) NextID() string {
	c.mutex.Lock()
	c.nextId++
	c.mutex.Unlock()
	return strconv.Itoa(c.nextId)
}

func (c *Client) Subscriptions() []string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	subs := make([]string, 0, len(c.StreamSubscriptions))
	for k := range c.StreamSubscriptions {
		subs = append(subs, k)
	}
	return subs
}

func (c *Client) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.closed = true

	// Signal both goroutines to stop
	select {
	case c.heartbeatDone <- true:
	default:
	}

	select {
	case c.handlerDone <- true:
	default:
	}

	// Close all stream channels to prevent blocking
	close(c.StreamLedger)
	close(c.StreamTransaction)
	close(c.StreamValidation)
	close(c.StreamManifest)
	close(c.StreamPeerStatus)
	close(c.StreamConsensus)
	close(c.StreamPathFind)
	close(c.StreamServer)
	close(c.StreamDefault)

	err := c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("WS write error:", err)
		return err
	}
	err = c.connection.Close()
	if err != nil {
		log.Println("WS close error:", err)
		return err
	}
	return nil
}
