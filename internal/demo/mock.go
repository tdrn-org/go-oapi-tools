//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package demo

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/tdrn-org/go-oapi-tools/internal/demo/api"
)

// StartMockServer starts and returns a new mock instance.
//
// StartMockServer panics in case of an error. The returned server
// is listening on localhost using a dynamic port.
func StartMockServer() *Server {
	return startMockServer(&mockHandler{}, "", nil)
}

type mockHandler struct {
	called atomic.Bool
}

func (h *mockHandler) GetUsers(ctx context.Context, request api.GetUsersRequestObject) (api.GetUsersResponseObject, error) {
	initialCall := h.called.CompareAndSwap(false, true)
	if !initialCall {
		return nil, errors.New("already called")
	}
	return api.GetUsers200JSONResponse([]string{"user1", "user2"}), nil
}
