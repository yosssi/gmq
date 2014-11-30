package cmd

import (
	"flag"

	"github.com/yosssi/gmq/mqtt/client"
)

// Cmd represents a command of an MQTT client.
type Cmd struct {
	// Name is a name of the command.
	Name string
	// Usage is a usage of the command.
	Usage string
	// Run runs the command.
	Run func(*client.Client, *Cmd) error
	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
}
