package webrtcsignalingserver

import (
	"net/http"

	"github.com/aliforever/go-httpjson"
)

type SignalingServer struct {
	storage *sdpStorage
}

func New() (ss *SignalingServer) {
	ss = &SignalingServer{storage: newSDPStorage()}
	return
}

func (ss *SignalingServer) Listen(address string, handler *http.ServeMux) (server *http.Server, err error) {
	m := handler
	if m == nil {
		m = http.NewServeMux()
	}

	m.HandleFunc("/sdp_handshake", ss.sdpHandShakerHandler)
	m.HandleFunc("/sdp_inform", ss.sdpInformListenerHandler)
	m.HandleFunc("/sdp_store", ss.sdpInformListenerHandler)

	server = &http.Server{Addr: address, Handler: m}
	err = server.ListenAndServe()

	return
}

func (ss *SignalingServer) AddSDPListener(id string) (l *Listener, err error) {
	l, err = ss.storage.AddSDPListener(id)
	return
}

func (ss *SignalingServer) sdpHandShakerHandler(writer http.ResponseWriter, request *http.Request) {
	var sar *sDPRequest

	err := httpjson.ParseRequest(writer, request, &sar)
	if err != nil {
		return
	}

	err = sar.Validate()
	if err != nil {
		httpjson.BadRequest(writer, err.Error())
		return
	}

	listener, err := ss.storage.GetSDPListener(sar.Id)
	if err != nil {
		httpjson.BadRequest(writer, err.Error())
		return
	}

	err = listener.WriteClientSDP(sar.SDP, sar.Data)
	if err != nil {
		httpjson.BadRequest(writer, err.Error())
		return
	}

	httpjson.Ok(writer, listener.ReadServerSDP())
}

func (ss *SignalingServer) sdpInformListenerHandler(writer http.ResponseWriter, request *http.Request) {
	var sar *sDPRequest

	err := httpjson.ParseRequest(writer, request, &sar)
	if err != nil {
		return
	}

	err = sar.Validate()
	if err != nil {
		httpjson.BadRequest(writer, err.Error())
		return
	}

	var l *Listener
	l, err = ss.storage.GetSDPListener(sar.Id)
	if err != nil {
		httpjson.BadRequest(writer, err.Error())
		return
	}

	err = l.WriteClientSDP(sar.SDP, sar.Data)
	if err != nil {
		httpjson.BadRequest(writer, err.Error())
		return
	}

	httpjson.Ok(writer, "success")
}

func (ss *SignalingServer) sdpStoreHandler(writer http.ResponseWriter, request *http.Request) {
	var sar *sDPRequest

	err := httpjson.ParseRequest(writer, request, &sar)
	if err != nil {
		return
	}

	err = sar.Validate()
	if err != nil {
		httpjson.BadRequest(writer, err.Error())
		return
	}

	err = ss.storage.AddSDPToStorage(sar.Id, sar.SDP, sar.Data)
	if err != nil {
		httpjson.BadRequest(writer, err.Error())
		return
	}
}
