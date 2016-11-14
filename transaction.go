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

import "time"

type Operation string

const (
	OP_CREATE   Operation = "CREATE"
	OP_TRANSFER Operation = "TRANSFER"
	OP_GENESIS  Operation = "GENESIS"
	VERSION     int64     = 1
)

// A Transaction is used to create and transfer assets
type Transaction struct {
	Version   int64     `json:"version,omitempty"`
	Timestamp int64     `json:"timestamp,omitempty"`
	Operation Operation `json:"operation,omitempty"`
	Asset     *Asset    `json:"asset,omitempty"`
}

// NewTransaction returns a new transaction
func NewTransaction(asset *Asset, operation Operation, version, timestamp int64) *Transaction {
	tx := new(Transaction)

	if version > 0 {
		tx.Version = version
	} else {
		tx.Version = VERSION
	}

	if timestamp > 0 {
		tx.Timestamp = timestamp
	} else {
		tx.Timestamp = time.Now().UTC().Unix()
	}

	tx.Operation = operation

	// If an asset is not defined in a `CREATE` transaction, create a default one
	if asset == nil && operation == OP_CREATE {
		asset = NewDefaultAsset()
	}
	tx.Asset = asset

	return tx
}

type TransactionInput struct {
	ConditionId   int64  `json:"cid,omitempty"`
	TransactionId string `json:"txid,omitempty"`
}

// A Fulfillment is used to spend assets locked by a Condition
type Fulfillment struct {
}

// Fulfillment shims a Cryptocondition Fulfillment for BigchainDB
func NewFulfillment()  {
	
}

// A Condition is used to lock an asset
type Condition struct {
}
