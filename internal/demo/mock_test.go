//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package demo_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-oapi-tools/internal/demo"
)

func TestMock(t *testing.T) {
	server := demo.StartMockServer()
	defer server.Stop(t.Context())

	client, err := demo.NewClient(server.APIURL())
	require.NoError(t, err)

	// Initial call (success)
	response, err := client.GetUsers(t.Context())
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, http.StatusOK, response.StatusCode())
	require.NotNil(t, response.JSON200)
	require.Equal(t, []string{"user1", "user2"}, *response.JSON200)

	// Follow-up calls (failure)
	response, err = client.GetUsers(t.Context())
	require.ErrorIs(t, err, demo.ErrAPIFailure)
	require.Nil(t, response)

	// Call to stopped server (failure)
	server.Stop(t.Context())
	response, err = client.GetUsers(t.Context())
	require.ErrorIs(t, err, demo.ErrClientFailure)
	require.Nil(t, response)
}
