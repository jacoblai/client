// Copyright (c) 2014 The gomqtt Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"fmt"
	"time"

	"github.com/gomqtt/packet"
)

func ExampleClient() {
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
}

func ExampleService() {
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
}
