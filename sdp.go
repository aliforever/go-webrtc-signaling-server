package webrtcsignalingserver

import (
	"encoding/base64"
	"encoding/json"

	"github.com/pion/webrtc/v3"
)

type SDPServer struct {
	sdp  *webrtc.SessionDescription
	b64  string
	Data map[string]string `json:"data,omitempty"`
}

func (s *SDPServer) SDP() (sd *webrtc.SessionDescription) {
	return s.sdp
}

func DecodeBase64StringToWebrtcSDP(sdpBase64Str string) (sdp *webrtc.SessionDescription, err error) {
	var decodedBase64 []byte
	decodedBase64, err = base64.StdEncoding.DecodeString(sdpBase64Str)
	if err != nil {
		return
	}

	var webrtcSDP *webrtc.SessionDescription
	err = json.Unmarshal(decodedBase64, &webrtcSDP)
	if err != nil {
		return
	}

	sdp = webrtcSDP
	return
}

func newServerSDP(sdp *webrtc.SessionDescription, data map[string]string) (serverSDP *SDPServer, err error) {
	var sdpBase64 string
	sdpBase64, err = EncodeWebrtcSdpToBase64(sdp)
	if err != nil {
		return
	}

	serverSDP = &SDPServer{
		sdp:  sdp,
		b64:  sdpBase64,
		Data: data,
	}

	return
}
