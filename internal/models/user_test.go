package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc    string
		user    LoginUser
		wantErr string
	}{
		{
			desc: "it works",
			user: LoginUser{
				Username: "testuser",
				Password: "password",
			},
		},
		{
			desc: "username is empty",
			user: LoginUser{
				Username: "",
				Password: "password",
			},
			wantErr: "validation error: username must not be empty",
		},
		{
			desc: "password is empty",
			user: LoginUser{
				Username: "testuser",
				Password: "",
			},
			wantErr: "validation error: password must not be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.user.Validate()
			if tc.wantErr != "" {
				assert.EqualError(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
