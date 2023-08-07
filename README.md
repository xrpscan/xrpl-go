# xrpl-go: A Go client for the XRP Ledger
[![Go Report Card](https://goreportcard.com/badge/github.com/xrpscan/xrpl-go)](https://goreportcard.com/report/github.com/xrpscan/xrpl-go) [![GoDoc](https://pkg.go.dev/badge/github.com/xrpscan/xrpl-go?status.svg)](https://pkg.go.dev/github.com/xrpscan/xrpl-go)

## Motivation

We use Go and XRPL websocket APIs a lot a XRPScan. Unfortunately, the state of 
the Go client libraries for XRPL at the time of publishing this package is not 
ideal. This is where `xrpl-go` comes in. It provides a low level API for interacting 
with [XRPL websocket interface](https://xrpl.org/http-websocket-apis.html). This 
library aims to mirror concepts of the official JavaScript/TypeScript library 
[xrpl.js](https://github.com/XRPLF/xrpl.js).

## Reference documentation

See the [full reference documentation](https://pkg.go.dev/github.com/xrpscan/xrpl-go) 
for all packages, functions and constants.

## Features

1. Sending requests to observe ledger state using public websocket API methods
2. Subscribing to changes in the ledger (ledger, transactions, validations streams)
3. Parsing ledger data into mode convenient formats [WIP]

## rippled versions

`xrpl-go` is currently tested with rippled versions > 1.9.4. While it should
also be compatible with later versions, newer features available on XRPL mainnet
may not be available on day 0.

## Installation

```bash
go get -u github.com/xrpscan/xrpl-go
```

## Getting started

Here are some examples showing typical use:

#### Establish a new websocket connection
```go
config := xrpl.ClientConfig{
  URL: "wss://s.altnet.rippletest.net:51233",
}
client, _ := xrpl.NewClient(config)
err := client.Ping([]byte("PING"))
if err != nil {
  panic(err)
}
```

#### Send `account_info` request
```go
request := xrpl.BaseRequest{
  "command": "account_info",
  "account": "rPT1Sjq2YGrBMTttX4GZHjKu9dyfzbpAYe",
  "ledger_index": "validated",
}
response, err := client.Request(request)
if err != nil {
  fmt.Println(err)
}
fmt.Println(response)
```

#### Subscribe to a single stream
```go
client.Subscribe([]string{
  xrpl.StreamTypeLedger,
})
for {
  ledger := <-client.StreamLedger
  fmt.Println(string(ledger))
}
```

#### Subscribe to multiple streams
```go
client.Subscribe([]string{
  xrpl.StreamTypeLedger,
  xrpl.StreamTypeTransaction,
  xrpl.StreamTypeValidations,
})
for {
  select {
  case ledger := <-client.StreamLedger:
    fmt.Println(string(ledger))
  case transaction := <-client.StreamTransaction:
    fmt.Println(string(transaction))
  case validation := <-client.StreamValidation:
    fmt.Println(string(validation))
  }
}
```

## Bugs

`xrpl-go` is a work in progress. If you discover a bug or come across erratic
behavior, please [create an issue](https://github.com/xrpscan/xrpl-go/issues/new) 
and we'll do our best to address it.

## References

- [XRPL HTTP/WebSocket API methods](https://xrpl.org/public-api-methods.html)
- [XRPL WebSocket streams](https://xrpl.org/subscribe.html)
- [JavaScript/TypeScript library for interacting with the XRP Ledger](https://js.xrpl.org)