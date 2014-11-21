package cmd

import "flag"

// Cmd represents a command of an MQTT client.
type Cmd struct {
	// Name is a name of the command.
	Name string
	// Usage is a usage of the command.
	Usage string
	// Run runs the command.
	Run func(*Cmd) error
	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
}
