package webrtcsignalingserver

import "github.com/pion/webrtc/v3"

type Listener struct {
	clientSDP chan *SDPClient
	serverSDP chan *SDPServer
}

func newListener() *Listener {
	return &Listener{
		clientSDP: make(chan *SDPClient),
		serverSDP: make(chan *SDPServer),
	}
}

func (l *Listener) WriteClientSDP(sdp string, data map[string]string) (err error) {
	var clientSDP *SDPClient
	clientSDP, err = newClientSDP(sdp, data)
	if err != nil {
		return
	}

	l.clientSDP <- clientSDP

	return
}

func (l *Listener) WriteServerSDP(sdp *webrtc.SessionDescription, data map[string]string) (err error) {
	var serverSDP *SDPServer
	serverSDP, err = newServerSDP(sdp, data)
	if err != nil {
		return
	}

	l.serverSDP <- serverSDP

	return
}

func (l *Listener) ReadClientSDP() (sdp *webrtc.SessionDescription, data map[string]string) {
	clientSDP := <-l.clientSDP
	sdp, data = clientSDP.sdp, clientSDP.Data()
	return
}

func (l *Listener) ReadServerSDP() (sdp *SDPServer) {
	sdp = <-l.serverSDP
	return
}
