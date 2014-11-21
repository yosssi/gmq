package cmd

import "flag"

// Cmd represents a command of an MQTT client.
type Cmd struct {
	// Run runs the command.
	Run func(*Cmd) error
	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
}
