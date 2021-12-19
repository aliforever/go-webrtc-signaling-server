package internal

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/pion/webrtc/v3"
)

type SDPRequest struct {
	Id   string            `json:"id"`
	SDP  string            `json:"sdp"` // BASE64
	Data map[string]string `json:"data"`
}

func (sr *SDPRequest) Validate() (err error) {
	if sr.Id == "" {
		err = errors.New("empty_id")
		return
	}

	_, err = sr.DecodeSDP()
	return
}

func (sr *SDPRequest) DecodeSDP() (sdp *webrtc.SessionDescription, err error) {
	var base64Decoded []byte
	base64Decoded, err = base64.StdEncoding.DecodeString(sr.SDP)
	if err != nil {
		return
	}

	err = json.Unmarshal(base64Decoded, &sdp)
	return
}
