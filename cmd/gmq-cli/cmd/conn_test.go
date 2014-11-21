package cmd

import (
	"testing"

	"github.com/yosssi/gmq/client"
)

func Test_conn(t *testing.T) {
	conn(client.New(), Conn)
}
