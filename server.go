package webrtc_signaling_server

import (
	"net/http"

	"github.com/aliforever/go-httpjson"
)

type SignalingServer struct {
	storage *sdpStorage
}

func NewSignalingServer() (ss *SignalingServer) {
	ss = &SignalingServer{storage: newSDPStorage()}
	return
}

func (ss *SignalingServer) Listen(address string) (server *http.Server, err error) {
	m := http.NewServeMux()
	m.HandleFunc("/inform-SDP", ss.informSDPListenerHandler)
	m.HandleFunc("/await-SDP", ss.getSDPFromListenerHandler)

	server = &http.Server{Addr: address, Handler: m}
	err = server.ListenAndServe()

	return
}

func (ss *SignalingServer) AddSDPListener(id string) (sdpListener chan *SDP, err error) {
	sdpListener, err = ss.storage.AddSDPListener(id)
	return
}

func (ss *SignalingServer) informSDPListenerHandler(writer http.ResponseWriter, request *http.Request) {
	var sar *sdpAddRequest

	err := httpjson.ParseRequest(writer, request, &sar)
	if err != nil {
		return
	}

	if sar.Id == "" {
		httpjson.BadRequest(writer, "empty_id")
		return
	}

	if sar.Sdp == "" {
		httpjson.BadRequest(writer, "empty_sdp")
		return
	}

	err = ss.storage.InformSDPListener(sar.Id, sar.Sdp, sar.Data)
	if err != nil {
		httpjson.BadRequest(writer, err.Error())
		return
	}

	httpjson.Ok(writer, "success")
	return
}

func (ss *SignalingServer) getSDPFromListenerHandler(writer http.ResponseWriter, request *http.Request) {
	var sar *sdpAwaitRequest

	err := httpjson.ParseRequest(writer, request, &sar)
	if err != nil {
		return
	}

	if sar.Id == "" {
		httpjson.BadRequest(writer, "empty_id")
		return
	}

	var sdpListener chan *SDP
	sdpListener, err = ss.storage.GetSDPListener(sar.Id)
	if err != nil {
		httpjson.BadRequest(writer, err.Error())
		return
	}

	httpjson.Ok(writer, newSDPAwaitResponse(<-sdpListener))
	return
}
