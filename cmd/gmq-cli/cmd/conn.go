package cmd

// Conn is a command which sends a connection request to the server.
var Conn = &Cmd{
	Name:  "conn",
	Usage: "send a connection request to the server",
	Run:   conn,
}

// conn sends a connection request to the server.
func conn(c *Cmd) error {
	return nil
}
