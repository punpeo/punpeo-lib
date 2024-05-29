package rmq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRmq(t *testing.T) {
	MqConnect("amqp://admin:你的密码@127.0.0.1:5672/", []string{"Dsh_Message_Exchange"})
	r := IntoQueue("Dsh_Message_Exchange", "Queue_2003", "test")
	assert.Equal(t, true, r, "success")
}
