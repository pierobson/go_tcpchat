# go_tcpchat

Really simple Go tcp chat client/server that works over the local network. I hardcoded my local ip for now but I'll fix that in a bit. Just getting used to Go as I've never used it before.

## Installation
Just `git clone https://github.com/pierobson/go_tcpchat.git` in your Go workspace to download.

## Running the Client/Server
You're going to want to run the server once and the client as many times as you would like (currently 64 concurrent users max).

### With Compilation
`go build` in the client and server directories.
`./go_tcpclient` or `./go_tcpserver` to run either.

### Without Compilation
`go run main.go` in each directory.
