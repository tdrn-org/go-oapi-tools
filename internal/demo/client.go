//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package demo

import (
	"net/url"
)

// NewClient creates a new Client instance using the given
// connect URL as well as Client options.
func NewClient(apiURL *url.URL, options ...ClientOption) (*Client, error) {
	return newClient(apiURL, options...)
}
