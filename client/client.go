package client

// Client represents an MQTT client.
type Client struct{}

// New creates and returns an MQTT client.
func New(opts *Options) *Client {
	return &Client{}
}
