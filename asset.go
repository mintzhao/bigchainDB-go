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
	"encoding/json"

	"github.com/satori/go.uuid"
)

// bigchainDB asset
type Asset struct {
	ID         string `json:"id,omitempty"`         // UUID version 4 (random) converted to a string of hex digits in standard form. Added server side.
	Divisible  bool   `json:"divisible,omitempty"`  // Whether the asset is divisible or not. Defaults to false.
	Updatable  bool   `json:"updatable,omitempty"`  // Whether the data in the asset can be updated in the future or not. Defaults to false.
	Refillable bool   `json:"refillable,omitempty"` // Whether the amount of the asset can change after its creation. Defaults to false.
	Data       []byte `json:"data,omitempty"`       // A user supplied JSON document with custom information about the asset. Defaults to null.
	Amount     int64  `json:"amount,omitempty"`     // The amount of “shares”. Only relevant if the asset is marked as divisible. Defaults to 1.
}

// NewDefaultAsset returns an asset using default value
func NewDefaultAsset() *Asset {
	return &Asset{
		ID:         uuid.NewV4().String(),
		Divisible:  false,
		Updatable:  false,
		Refillable: false,
		Data:       []byte(""),
		Amount:     1,
	}
}

// NewAsset returns a new asset all by the params, notice that if divisible is ture, amount must not be blank
func NewAsset(id string, divisible, updatable, refillable bool, metadata []byte, amount ...int64) *Asset {
	assert := new(Asset)

	assert.Data = metadata
	if id == "" {
		id = uuid.NewV4().String()
	}
	assert.ID = id
	assert.Divisible = divisible
	assert.Updatable = updatable
	assert.Refillable = refillable

	if assert.Divisible && len(amount) == 0 {
		amount = []int64{1}
	}
	if assert.Divisible {
		assert.Amount = amount[0]
	}

	return assert
}

// ParseAsset parses json format asset into Asset object
func ParseAsset(astr string) *Asset {
	a := new(Asset)
	if err := json.Unmarshal([]byte(astr), a); err != nil {
		bigchainLogger.Errorf("ParseAsset() return error: %v", err)
		return nil
	}

	return a
}

// JSON returns json format assert
func (a *Asset) JSON() string {
	assetBytes, err := json.Marshal(a)
	if err != nil {
		bigchainLogger.Errorf("a.JSON() return error: %v", err)
		return ""
	}

	return string(assetBytes)
}

// Equal returns whether two assets have same json format string
func (a *Asset) Equal(other *Asset) bool {
	return a.JSON() == other.JSON()
}