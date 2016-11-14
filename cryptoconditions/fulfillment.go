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
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	FULFILLMENT_REGEX = `^cf:([1-9a-f][0-9a-f]{0,3}|0):[a-zA-Z0-9_-]*$`
)

type Fulfillment struct {
	TypeId  int64 // the type ID of this fulfillment
	Bitmask int64 // the bitmask of this fulfillment, for simple fulfillment types this is simply the bit representing this type, for meta-fulfillments, these are the bits representing the types of the subconditions.

}

// FromURI create a fulfillment from a URI.
// This function will parse a fulfillment URI and construct a corresponding Fulfillment object.
func FromURI(serializedFulfillment string) (*Fulfillment, error) {
	if serializedFulfillment == "" {
		return nil, fmt.Errorf("empty serialized fulfillment string")
	}

	pieces := strings.Split(serializedFulfillment, ":")
	if len(pieces) != 3 || pieces[0] != "cf" {
		return nil, fmt.Errorf("serialized fulfillment must start with 'cf'")
	}

	if matched, err := regexp.MatchString(FULFILLMENT_REGEX, serializedFulfillment); !matched || err != nil {
		return fmt.Errorf("invalid fulfillment format")
	}

	typeId, err := strconv.ParseInt(pieces[1], 16, 64)
	if err != nil {
		return nil, err
	}
	payload, err := base64.URLEncoding.DecodeString(pieces[2])
	if err != nil {
		return nil, err
	}

}
