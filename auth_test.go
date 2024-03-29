/*
Copyright AppsCode Inc. and Contributors

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

package client_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.bytebuilders.dev/client"
)

func TestClient_Signin(t *testing.T) {
	tests := map[string]struct {
		username      string
		password      string
		expectedError error
	}{
		"Credential was not provided": {
			expectedError: client.ErrNotFound,
		},
		"Invalid user": {
			username:      "does-not-exist",
			password:      "password",
			expectedError: client.ErrNotFound,
		},
		"Invalid password": {
			username:      "appscode",
			password:      "invalid",
			expectedError: client.ErrNotFound,
		},
		"Valid credentials": {
			username:      "appscode",
			password:      "password",
			expectedError: nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := client.NewClient(client.TestServerURL)
			cookies, err := c.Signin(client.SignInParams{UserName: tt.username, Password: tt.password})
			if tt.expectedError != nil {
				assert.True(t, errors.Is(err, tt.expectedError))
				return
			}
			if !assert.Nil(t, err) {
				return
			}
			if !assert.NotNil(t, cookies) {
				return
			}
		})
	}
}

func TestClient_Signout(t *testing.T) {
	t.Run("User was logged in", func(t *testing.T) {
		c := client.NewClient(client.TestServerURL)
		cookies, err := c.Signin(client.SignInParams{UserName: client.TestServerUser, Password: client.TestServrPassword})
		if !assert.Nil(t, err) {
			return
		}
		c = client.NewClient(client.TestServerURL).WithCookies(cookies)
		_, err = c.GetCurrentUser()
		if !assert.Nil(t, err) {
			return
		}
		if !assert.Nil(t, c.Signout()) {
			return
		}

		_, err = c.GetCurrentUser()
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, client.ErrUnAuthorized))
	})

	t.Run("User wasn't logged in", func(t *testing.T) {
		c := client.NewClient(client.TestServerURL).WithBasicAuth(client.TestServerUser, client.TestServrPassword)
		if !assert.Nil(t, c.Signout()) {
			return
		}
	})
}
