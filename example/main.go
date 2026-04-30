package main

import (
	"context"
	"fmt"
	"os"

	"github.com/carabiner-labs/protect/client"
	controlv1 "github.com/carabiner-labs/protect/gen/protect/control/v1"
)

func main() {
	ctx := context.Background()
	c, err := client.Dial(
		ctx, "unix:///var/lib/edera/protect/daemon.socket",
		client.WithInsecure(),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Dialing socket: %v", err)
		os.Exit(1)
	}
	defer c.Close()

	// Get host status:
	reply, err := c.Control.GetHostStatus(ctx, &controlv1.GetHostStatusRequest{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Calling c.Control.GetHostStatus: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Protect daemon status: %+v", reply)
}
