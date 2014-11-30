package cmd

import (
	"testing"

	"github.com/yosssi/gmq/mqtt/client"
)

func Test_conn(t *testing.T) {
	conn(client.New(), Conn)
}
