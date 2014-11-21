package client

// Client represents an MQTT client.
type Client struct{}

// New creates and returns an MQTT client.
func New() *Client {
	return &Client{}
}
