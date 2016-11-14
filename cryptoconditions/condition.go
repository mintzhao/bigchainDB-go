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
package cryptoconditions

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
)

const (
	// Our current implementation can only represent up to 32 bits for our bitmask
	MAX_SAFE_BITMASK = 0xfffffff

	// Feature suites supported by this implementation
	SUPPORTED_BITMASK = 0x3f

	// Max fulfillment size supported by this implementation
	MAX_FULFILLMENT_LENGTH = 65535

	// Regex for validating conditions
	// This is a generic, future-proof version of the crypto-condition regular expression.
	CONDITION_REGEX = `^cc:([1-9a-f][0-9a-f]{0,3}|0):[1-9a-f][0-9a-f]{0,15}:[a-zA-Z0-9_-]{0,86}:([1-9][0-9]{0,17}|0)$`

	// This is a stricter version based on limitations of the current implementation.
	// Specifically, we can't handle bitmasks greater than 32 bits.
	CONDITION_REGEX_STRICT = `^cc:([1-9a-f][0-9a-f]{0,3}|0):[1-9a-f][0-9a-f]{0,7}:[a-zA-Z0-9_-]{0,86}:([1-9][0-9]{0,17}|0)$`
)

/*
Crypto-condition.
A primary design goal of crypto-conditions was to keep the size of conditions
constant. Even a complex multi-signature can be represented by the same size
condition as a simple hashlock.
However, this means that a condition only carries the absolute minimum
information required. It does not tell you anything about its structure.
All that is included with a condition is the fingerprint (usually a hash of
the parts of the fulfillment that are known up-front, e.g. public keys), the
maximum fulfillment size, the set of features used and the condition type.
This information is just enough that an implementation can tell with
certainty whether it would be able to process the corresponding fulfillment.
*/
type Condition struct {
	// For simple condition types this is simply the bit representing this type.
	// For structural conditions, this is the bitwise OR of the bitmasks of
	// the condition and all its subconditions, recursively.
	Bitmask int64

	// The type is a unique integer ID assigned to each type of condition.
	TypeId int64

	hash                 string
	maxFulfillmentLength int
}

/*
A primary component of all conditions is the hash. It encodes the static
properties of the condition. This method enables the conditions to be
constant size, no matter how complex they actually are. The data used to
generate the hash consists of all the static properties of the condition
and is provided later as part of the fulfillment.
*/
func (c *Condition) Hash() string {
	return c.hash
}

/*
Validate and set the hash of this condition.
Typically conditions are generated from fulfillments and the hash is
calculated automatically. However, sometimes it may be necessary to
construct a condition URI from a known hash. This method enables that case.
*/
func (c *Condition) SetHash(h string) {
	c.hash = h
}

/*
The maximum fulfillment length is the maximum allowed length for any
fulfillment payload to fulfill this condition.
The condition defines a maximum fulfillment length which all
implementations will enforce. This allows implementations to verify that
their local maximum fulfillment size is guaranteed to accomodate any
possible fulfillment for this condition.
Otherwise an attacker could craft a fulfillment which exceeds the maximum
size of one implementation, but meets the maximum size of another, thereby
violating the fundamental property that fulfillments are either valid
everywhere or nowhere.
*/
func (c *Condition) MaxFulfillmentLength() int {
	return c.maxFulfillmentLength
}

/*
Set the maximum fulfillment length.
*/
func (c *Condition) SetMaxFulfillmentLength(length int) {
	c.maxFulfillmentLength = length
}

/*
Generate the URI form encoding of this condition.
Turns the condition into a URI containing only URL-safe characters. This
format is convenient for passing around conditions in URLs, JSON and other text-based formats.
"cc:" BASE16(TYPE_ID) ":" BASE16(BITMASK) ":" BASE64URL(HASH) ":" BASE10(MAX_FULFILLMENT_LENGTH)
*/
func (c *Condition) URI() string {
	return fmt.Sprintf("cc:%d:%d:%s:%d", c.TypeId, c.Bitmask, base64.RawURLEncoding.EncodeToString([]byte(c.hash)), c.maxFulfillmentLength)
}

func (c *Condition) JSON() string {
	conditionBytes, err := json.Marshal(map[string]interface{}{
		"type":    "condition",
		"type_id": c.TypeId,
		"bitmask": c.Bitmask,
		"hash":    base58.Encode([]byte(c.hash)),
		"max_fulfillment_length": c.maxFulfillmentLength,
	})
	if err != nil {
		return ""
	}

	return string(conditionBytes)
}

func (c *Condition) ParseJSON(data string) error {
	vals := make(map[string]interface{})
	if err := json.Unmarshal([]byte(data), vals); err != nil {
		return err
	}

	c.TypeId = vals["type_id"].(int64)
	c.Bitmask = vals["bitmask"].(int64)
	c.hash = string(base58.Decode(vals["hash"]))
	c.maxFulfillmentLength = vals["max_fulfillment_length"].(int)

	return nil
}

/*
Ensure the condition is valid according the local rules.
Checks the condition against the local bitmask (supported condition types)
and the local maximum fulfillment size.
*/
func (c *Condition) Validate() error {
	if c.Bitmask > MAX_SAFE_BITMASK {
		return fmt.Errorf("bitmask too large to be safely represented")
	}

	if c.Bitmask &^ SUPPORTED_BITMASK {
		return fmt.Errorf("condition requested unsupported feature suites")
	}

	if c.maxFulfillmentLength > MAX_FULFILLMENT_LENGTH {
		return fmt.Errorf("condition requested too large of a max fulfillment size")
	}

	return nil
}
