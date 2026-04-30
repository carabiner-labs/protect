# Edera Protect Daemon go-client

This a thin client library that allows go programs to talk to the
[Edera Protect](https://edera.dev/containers) daemon. This
client is automatically generated from the opens source
[protbuf defintions](https://github.com/edera-dev/protos) of
the daemon API.

## Simple Example

```golang

	client, err := client.Dial(ctx, "unix:///var/lib/edera/protect/daemon.socket")
	if err != nil {
        os.Exit(1)
    }
	defer c.Close()

    // Get host status:
	reply, err := c.Control.GetHostStatus(ctx, &controlv1.GetHostStatusRequest{})
    if err != nil {
        os.Exit(1)
    }

    fmt.Printf("Status: %+v", reply)
```
