# GMQ - Pure Go MQTT Client

[![wercker status](https://app.wercker.com/status/3e5533f6f9aa61384eb2dd2d8f102cfd/m "wercker status")](https://app.wercker.com/project/bykey/3e5533f6f9aa61384eb2dd2d8f102cfd)
[![Build status](https://ci.appveyor.com/api/projects/status/7gigy6i4tknxh9x3?svg=true)](https://ci.appveyor.com/project/yosssi/gmq)
[![Coverage Status](https://img.shields.io/coveralls/yosssi/gmq.svg)](https://coveralls.io/r/yosssi/gmq?branch=master)
[![GoDoc](https://godoc.org/github.com/yosssi/gmq?status.svg)](https://godoc.org/github.com/yosssi/gmq)
[![Gitter](https://badges.gitter.im/Join Chat.svg)](https://gitter.im/yosssi/gmq?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

## Overview

GMQ is a pure Go MQTT client. This library is compatible with [MQTT Version 3.1.1](http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html).

## Installation

```sh
$ go get -u github.com/yosssi/gmq/...
```

## MQTT Client Library

### Example

```go
package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

func main() {
	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	// Create an MQTT Client.
	cli := client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})

	// Terminate the Client.
	defer cli.Terminate()

	// Connect to the MQTT Server.
	err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  "iot.eclipse.org:1883",
		ClientID: []byte("example-client"),
	})
	if err != nil {
		panic(err)
	}

	// Subscribe to topics.
	err = cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte("foo"),
				QoS:         mqtt.QoS0,
				// Define the processing of the message handler.
				Handler: func(topicName, message []byte) {
					fmt.Println(string(topicName), string(message))
				},
			},
			&client.SubReq{
				TopicFilter: []byte("bar/#"),
				QoS:         mqtt.QoS1,
				Handler: func(topicName, message []byte) {
					fmt.Println(string(topicName), string(message))
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// Publish messages.
	err = cli.Publish(&client.PublishOptions{
		QoS:       mqtt.QoS0,
		TopicName: []byte("bar/baz"),
		Message:   []byte("testMessage"),
	})
	if err != nil {
		panic(err)
	}

	// Unsubscribe from topics.
	err = cli.Unsubscribe(&client.UnsubscribeOptions{
		TopicFilters: [][]byte{
			[]byte("foo"),
		},
	})
	if err != nil {
		panic(err)
	}

	// Wait for receiving a signal.
	<-sigc

	// Disconnect the network connection.
	if err := cli.Disconnect(); err != nil {
		panic(err)
	}
}
```

### APIs

See [GoDoc](https://godoc.org/github.com/yosssi/gmq).

## MQTT Client Command Line Application


