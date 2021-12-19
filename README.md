# Go Webrtc Signaling Server
This package is used to listen for Remote SDP (Session Description) as well as information about local SDP through HTTP server.  

## Install
```go get -u github.com/aliforever/go-webrtc-signaling-server```

## Usage
There are 3 http handlers in this package:
1. `/sdp_handshake` Looks for a defined SDP listener and pass browser SDP in return of a remote SDP.
2. `/sdp_inform` Inform a defined SDP Listener and let go
3. `/sdp_store` Store SDP in storage and let go 

## Examples (TODO)
1. [examples/listener/main.go](examples/listener/main.go)
2. TODO
3. TODO