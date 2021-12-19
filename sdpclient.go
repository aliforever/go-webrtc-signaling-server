package webrtcsignalingserver

import (
	"encoding/base64"
	"encoding/json"

	"github.com/pion/webrtc/v3"
)

type SDPClient struct {
	b64  string
	sdp  *webrtc.SessionDescription
	data map[string]string
}

func (sc *SDPClient) SDP() *webrtc.SessionDescription {
	return sc.sdp
}

func (sc *SDPClient) Data() map[string]string {
	return sc.data
}

func newClientSDP(sdpBase64Str string, data map[string]string) (sdp *SDPClient, err error) {
	var webrtcSDP *webrtc.SessionDescription
	webrtcSDP, err = DecodeBase64StringToWebrtcSDP(sdpBase64Str)
	if err != nil {
		return
	}

	sdp = &SDPClient{
		b64:  sdpBase64Str,
		data: data,
		sdp:  webrtcSDP,
	}
	return
}

func EncodeWebrtcSdpToBase64(sdp *webrtc.SessionDescription) (encodedBase64 string, err error) {
	var jsonSDP []byte
	jsonSDP, err = json.Marshal(sdp)
	if err != nil {
		return
	}

	encodedBase64 = base64.StdEncoding.EncodeToString(jsonSDP)
	return
}
