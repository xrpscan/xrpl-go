package xrpl

import (
	"errors"
	"fmt"
	"log"
	"math"
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
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	HeartbeatInterval  time.Duration
	QueueCapacity      int
}

type Client struct {
	config              ClientConfig
	connection          *websocket.Conn
	heartbeatDone       chan bool
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

	if config.ReadTimeout < 0 ||
		config.ReadTimeout <= config.HeartbeatInterval ||
		config.ReadTimeout >= math.MaxInt32 {
		return fmt.Errorf("connection read timeout out of bounds: %d", config.ReadTimeout)
	}
	if config.WriteTimeout < 0 ||
		config.WriteTimeout <= config.HeartbeatInterval ||
		config.WriteTimeout >= math.MaxInt32 {
		return fmt.Errorf("connection write timeout out of bounds: %d", config.WriteTimeout)
	}
	if config.HeartbeatInterval < 0 ||
		config.HeartbeatInterval >= math.MaxInt32 {
		return fmt.Errorf("connection heartbeat interval out of bounds: %d", config.HeartbeatInterval)
	}

	return nil
}

func NewClient(config ClientConfig) *Client {
	if config.ReadTimeout == 0 {
		config.ReadTimeout = 20
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 20
	}
	if config.HeartbeatInterval == 0 {
		config.HeartbeatInterval = 5
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
	c.connection.SetPongHandler(c.handlePong)
	go c.handleResponse()
	go c.heartbeat()
	return c.connection, nil
}

func (c *Client) Reconnect() error {
	// Close old websocket connection
	c.Close()

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
	if err := c.connection.WriteMessage(websocket.PingMessage, message); err != nil {
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
	c.heartbeatDone <- true

	err := c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("WS write error: ", err)
		return err
	}
	err = c.connection.Close()
	if err != nil {
		log.Println("WS close error: ", err)
		return err
	}
	return nil
}
