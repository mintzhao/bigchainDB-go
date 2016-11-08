/*
Copyright Mojing Inc. 2016 All Rights Reserved.
Written by mint.zhao.chiu@gmail.com. github.com: https://www.github.com/mintzhao

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package bigchainDB

import (
	"errors"

	"github.com/op/go-logging"
	"github.com/parnurzeal/gorequest"
)

const (
	DEFAULT_NODE = "http://localhost:9984/api/v1"
)

var (
	ErrInvalidSignKey   = errors.New("invalid sign key")
	ErrInvalidVerifyKey = errors.New("invalid verify key")

	bigchainLogger = logging.MustGetLogger("bigchainDB")
)

type BigchainDBDriver struct {
	// BigchainDB nodes to connect to. Currently, the full URL must be given, if not, using DEFAULT_NODE
	node string

	// The base58 encoded public key for the ED25519 curve to bind this driver with.
	verifyKey string

	// The base58 encoded private key for the ED25519 curve to bind this driver with.
	signKey string

	// http requests
	requests chan *gorequest.SuperAgent
}

// NewBigchainDBDriver returns a driver
func NewBigchainDBDriver(node string, cacheLen int) *BigchainDBDriver {
	if node == "" {
		node = DEFAULT_NODE
	}

	driver := &BigchainDBDriver{
		node:     node,
		requests: make(chan *gorequest.Request, cacheLen),
	}
	go driver.send()

	return driver
}

// send sends requests to bigchain node
func (d *BigchainDBDriver) send() {
	for {
		select {
		case req := <-d.requests:
			req.End()
		}
	}
}

// Issue a transaction to create an asset.
func (d *BigchainDBDriver) Create(asset *Asset, verifyKey, signKey string) error {
	if verifyKey == "" {
		verifyKey = d.verifyKey
	}
	if signKey == "" {
		signKey = d.signKey
	}

	if asset == nil {
		asset = NewDefaultAsset()
	}
}
