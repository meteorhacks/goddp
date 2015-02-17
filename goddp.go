package goddp

import (
	"github.com/meteorhacks/goddp/client"
	"github.com/meteorhacks/goddp/server"
)

func NewClient() client.Client {
	return client.New()
}

func NewServer() server.Server {
	return server.New()
}
