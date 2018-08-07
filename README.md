# Go UDP Test
This is a simple UDP client and server implemented in go.

## Building
- Clone the repo into src folder in your GOPATH.
- CD into the folder and run `go build`

## Server Usage
To start a server run `go-udp-test -mode=server -address="addr:port"`

You will be taken to a prompt that accepts the following commands:
- `c`: List clients and their ids.
- `s`: Send a message. The format of the command is `s clientid msg`
- `q`: Quit.

In order to send a message to a client, it must have sent a message first in order for its address to be registered.

## Client Usage
To start a client run `go-udp-test -mode=client -address="addr:port"`

Simply type a message and press enter to send it to the server.