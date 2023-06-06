package xrpl

import (
	"encoding/json"
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
	ConnectionTimeout  time.Duration
	FeeCushion         uint32
	Key                string
	MaxFeeXRP          uint64
	Passphrase         byte
	Proxy              byte
	ProxyAuthorization byte
	Timeout            time.Duration
	QueueCapacity      int
}

type Client struct {
	config            ClientConfig
	connection        *websocket.Conn
	closed            bool
	mutex             sync.Mutex
	response          *http.Response
	StreamLedger      chan []byte
	StreamTransaction chan []byte
	StreamValidation  chan []byte
	StreamManifest    chan []byte
	StreamPeerStatus  chan []byte
	StreamConsensus   chan []byte
	StreamPathFind    chan []byte
	StreamServer      chan []byte
	StreamDefault     chan []byte
	requestQueue      map[string](chan<- BaseResponse)
	nextId            int
	err               error
}

func (config *ClientConfig) Validate() error {
	if len(config.URL) == 0 {
		return errors.New("cannot create a new connection with an empty URL")
	}

	if config.ConnectionTimeout < 0 || config.ConnectionTimeout >= math.MaxInt32 {
		return fmt.Errorf("connection timeout out of bounds: %d", config.ConnectionTimeout)
	}

	if config.Timeout < 0 || config.Timeout >= math.MaxInt32 {
		return fmt.Errorf("timeout out of bounds: %d", config.Timeout)
	}

	return nil
}

func NewClient(config ClientConfig) *Client {
	if err := config.Validate(); err != nil {
		panic(err)
	}

	if config.ConnectionTimeout == 0 {
		config.ConnectionTimeout = 60 * time.Second
	}

	if config.QueueCapacity == 0 {
		config.QueueCapacity = 128
	}

	client := &Client{
		config:            config,
		StreamLedger:      make(chan []byte, config.QueueCapacity),
		StreamTransaction: make(chan []byte, config.QueueCapacity),
		StreamValidation:  make(chan []byte, config.QueueCapacity),
		StreamManifest:    make(chan []byte, config.QueueCapacity),
		StreamPeerStatus:  make(chan []byte, config.QueueCapacity),
		StreamConsensus:   make(chan []byte, config.QueueCapacity),
		StreamPathFind:    make(chan []byte, config.QueueCapacity),
		StreamServer:      make(chan []byte, config.QueueCapacity),
		StreamDefault:     make(chan []byte, config.QueueCapacity),
		requestQueue:      make(map[string](chan<- BaseResponse)),
		nextId:            0,
	}
	c, r, err := websocket.DefaultDialer.Dial(config.URL, nil)
	if err != nil {
		client.err = err
		return nil
	}
	defer r.Body.Close()
	client.connection = c
	client.response = r
	client.connection.SetPongHandler(client.handlePong)
	go client.handleResponse()
	return client
}

func (c *Client) Ping(message []byte) error {
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

func (c *Client) Subscribe(streams []string) (BaseResponse, error) {
	req := BaseRequest{
		"command": "subscribe",
		"streams": streams,
	}
	res, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) Unsubscribe(streams []string) (BaseResponse, error) {
	req := BaseRequest{
		"command": "unsubscribe",
		"streams": streams,
	}
	res, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Send a websocket request. This method takes a BaseRequest object and automatically adds
// incremental request ID to it.
//
// Example usage:
//
//	req := BaseRequest{
//		"command": "account_info",
//		"account": "rG1QQv2nh2gr7RCZ1P8YYcBUKCCN633jCn",
//		"ledger_index": "current",
//	}
//
//	err := client.Request(req, func(){})
func (c *Client) Request(req BaseRequest) (BaseResponse, error) {
	requestId := c.NextID()
	req["id"] = requestId
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	ch := make(chan BaseResponse, 1)

	c.mutex.Lock()
	c.requestQueue[requestId] = ch
	err = c.connection.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return nil, err
	}
	c.mutex.Unlock()

	res := <-ch
	return res, nil
}

func (c *Client) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.closed = true

	err := c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("Write close error: ", err)
		return err
	}
	err = c.connection.Close()
	if err != nil {
		log.Println("Write close error: ", err)
		return err
	}
	return nil
}
