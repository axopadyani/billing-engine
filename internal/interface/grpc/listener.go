package grpc

import (
	"net"
	"os"
)

// InitListener initializes and returns a TCP network listener.
// It sets up the listener on a specified PORT environment variable, defaulting to 8080.
//
// The function does not take any parameters.
//
// Returns:
//   - net.Listener: A network listener that can be used to accept incoming connections.
//   - error: An error if the listener could not be created, or nil if successful.
func InitListener() (net.Listener, error) {
	addr := ":8080"
	if os.Getenv("PORT") != "" {
		addr = ":" + os.Getenv("PORT")
	}

	return net.Listen("tcp", addr)
}
