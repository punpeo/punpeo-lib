package utils

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var (
	key = "5qHIk7yMbGu29SFA6"
	iv  = "Xoji5qa9"
)

func TestParseSk(t *testing.T) {
	security_key := "U1/OIwAfq3hh/YOXJsdBj2nMq1TC7StLAUXBh8+ell3SfXGdYpaS0/fo7FtHGbNkXS4WyXgLhBvIow3ohs5ACQDkPNKidtLpzX2wYUFgZgg="

	var my = SkToUser{
		UserId:     81,
		Phone:      "15820226209",
		ExpireTime: 1664450788,
	}
	user := ParseSk(security_key, key, iv)
	t.Logf("%+v\n", user)
	assert.NotNil(t, user)
	assert.Equal(t, user.UserId, uint(81))
	assert.Equal(t, user.Phone, "15820226209")
	assert.True(t, reflect.DeepEqual(&my, user))

	security_key = "xxxx"
	user = ParseSk(security_key, key, iv)
	assert.Nil(t, user)
}

func TestParseSkUrl(t *testing.T) {
	security_key := "U1%2FOIwAfq3hh%2FYOXJsdBj2nMq1TC7StLAUXBh8%2Bell3SfXGdYpaS0%2Ffo7FtHGbNkXS4WyXgLhBvIow3ohs5ACQDkPNKidtLpzX2wYUFgZgg%3D"
	var my = SkToUser{
		UserId:     81,
		Phone:      "15820226209",
		ExpireTime: 1664450788,
	}
	user := ParseSkUrl(security_key, key, iv)

	assert.NotNil(t, user)
	assert.Equal(t, user.UserId, uint(81))
	assert.Equal(t, user.Phone, "15820226209")
	assert.True(t, reflect.DeepEqual(&my, user))
}
