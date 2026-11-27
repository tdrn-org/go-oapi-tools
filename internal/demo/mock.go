//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package demo

// StartMockServer starts and returns a new mock instance.
//
// StartMockServer panics in case of an error. The returned server
// is listening on localhost using a dynamic port.
func StartMockServer() *Server {
	return startMockServer(nil, "", nil)
}
