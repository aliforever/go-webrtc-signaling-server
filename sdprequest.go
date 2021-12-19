package webrtcsignalingserver

import (
	"errors"
)

type sDPRequest struct {
	Id   string            `json:"id"`
	SDP  string            `json:"sdp"` // BASE64
	Data map[string]string `json:"data"`
}

func (sr *sDPRequest) Validate() (err error) {
	if sr.Id == "" {
		err = errors.New("empty_id")
		return
	}
	return
}
