// Copyright 2020 The Matrix.org Foundation C.I.C.
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

package api

import (
	"context"
	"encoding/json"
	"strings"
	"time"
)

type KeyInternalAPI interface {
	PerformUploadKeys(ctx context.Context, req *PerformUploadKeysRequest, res *PerformUploadKeysResponse)
	PerformClaimKeys(ctx context.Context, req *PerformClaimKeysRequest, res *PerformClaimKeysResponse)
	QueryKeys(ctx context.Context, req *QueryKeysRequest, res *QueryKeysResponse)
}

// KeyError is returned if there was a problem performing/querying the server
type KeyError struct {
	Err string
}

func (k *KeyError) Error() string {
	return k.Err
}

// DeviceKeys represents a set of device keys for a single device
// https://matrix.org/docs/spec/client_server/r0.6.1#post-matrix-client-r0-keys-upload
type DeviceKeys struct {
	// The user who owns this device
	UserID string
	// The device ID of this device
	DeviceID string
	// The raw device key JSON
	KeyJSON []byte
}

// OneTimeKeys represents a set of one-time keys for a single device
// https://matrix.org/docs/spec/client_server/r0.6.1#post-matrix-client-r0-keys-upload
type OneTimeKeys struct {
	// The user who owns this device
	UserID string
	// The device ID of this device
	DeviceID string
	// A map of algorithm:key_id => key JSON
	KeyJSON map[string]json.RawMessage
}

// Split a key in KeyJSON into algorithm and key ID
func (k *OneTimeKeys) Split(keyIDWithAlgo string) (algo string, keyID string) {
	segments := strings.Split(keyIDWithAlgo, ":")
	return segments[0], segments[1]
}

// OneTimeKeysCount represents the counts of one-time keys for a single device
type OneTimeKeysCount struct {
	// The user who owns this device
	UserID string
	// The device ID of this device
	DeviceID string
	// algorithm to count e.g:
	// {
	//   "curve25519": 10,
	//   "signed_curve25519": 20
	// }
	KeyCount map[string]int
}

// PerformUploadKeysRequest is the request to PerformUploadKeys
type PerformUploadKeysRequest struct {
	DeviceKeys  []DeviceKeys
	OneTimeKeys []OneTimeKeys
}

// PerformUploadKeysResponse is the response to PerformUploadKeys
type PerformUploadKeysResponse struct {
	// A fatal error when processing e.g database failures
	Error *KeyError
	// A map of user_id -> device_id -> Error for tracking failures.
	KeyErrors        map[string]map[string]*KeyError
	OneTimeKeyCounts []OneTimeKeysCount
}

// KeyError sets a key error field on KeyErrors
func (r *PerformUploadKeysResponse) KeyError(userID, deviceID string, err *KeyError) {
	if r.KeyErrors[userID] == nil {
		r.KeyErrors[userID] = make(map[string]*KeyError)
	}
	r.KeyErrors[userID][deviceID] = err
}

type PerformClaimKeysRequest struct {
}

type PerformClaimKeysResponse struct {
	Error *KeyError
}

type QueryKeysRequest struct {
	// Maps user IDs to a list of devices
	UserToDevices map[string][]string
	Timeout       time.Duration
}

type QueryKeysResponse struct {
	// Map of remote server domain to error JSON
	Failures map[string]interface{}
	// Map of user_id to device_id to device_key
	DeviceKeys map[string]map[string]json.RawMessage
	// Set if there was a fatal error processing this query
	Error *KeyError
}
