# gomqtt/client

[![Circle CI](https://img.shields.io/circleci/project/gomqtt/client.svg)](https://circleci.com/gh/gomqtt/client)
[![Coverage Status](https://coveralls.io/repos/gomqtt/client/badge.svg?branch=master&service=github)](https://coveralls.io/github/gomqtt/client?branch=master)
[![GoDoc](https://godoc.org/github.com/gomqtt/client?status.svg)](http://godoc.org/github.com/gomqtt/client)
[![Release](https://img.shields.io/github/release/gomqtt/client.svg)](https://github.com/gomqtt/client/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/gomqtt/client)](http://goreportcard.com/report/gomqtt/client)

**This go package implements functionality for communicating with a [MQTT 3.1.1](http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/) broker.**

## Installation

Get it using go's standard toolset:

```bash
$ go get github.com/gomqtt/client
```

## Usage

### Service

```go
wait := make(chan struct{})
done := make(chan struct{})

options := NewOptions()
options.ClientID = "gomqtt/service"

s := NewService()

s.Online = func(resumed bool) {
    fmt.Println("online!")
    fmt.Printf("resumed: %v\n", resumed)
}

s.Offline = func() {
    fmt.Println("offline!")
    close(done)
}

s.Message = func(msg *packet.Message) {
    fmt.Printf("message: %s - %s\n", msg.Topic, msg.Payload)
    close(wait)
}

ClearSession("mqtt://try:try@broker.shiftr.io", "gomqtt/service")

s.Start("mqtt://try:try@broker.shiftr.io", options)

s.Subscribe("test", 0, false).Wait()

s.Publish("test", []byte("test"), 0, false)

<-wait

s.Stop(true)

<-done

// Output:
// online!
// resumed: false
// message: test - test
// offline!
```

### Client

```go
done := make(chan struct{})

c := New()

c.Callback = func(msg *packet.Message, err error) {
    if err != nil {
        panic(err)
    }

    fmt.Printf("%s: %s\n", msg.Topic, msg.Payload)
    close(done)
}

options := NewOptions()
options.ClientID = "gomqtt/client"

connectFuture, err := c.Connect("mqtt://try:try@broker.shiftr.io", options)
if err != nil {
    panic(err)
}

err = connectFuture.Wait(10 * time.Second)
if err != nil {
    panic(err)
}

subscribeFuture, err := c.Subscribe("test", 0)
if err != nil {
    panic(err)
}

err = subscribeFuture.Wait(10 * time.Second)
if err != nil {
    panic(err)
}

publishFuture, err := c.Publish("test", []byte("test"), 0, false)
if err != nil {
    panic(err)
}

err = publishFuture.Wait(10 * time.Second)
if err != nil {
    panic(err)
}

<-done

err = c.Disconnect()
if err != nil {
    panic(err)
}

// Output:
// test: test
```
