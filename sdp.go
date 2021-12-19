package internal

type SDP struct {
	SDP        string            `json:"sdp"`
	RemoteData map[string]string `json:"-"`
	LocalData  map[string]string `json:"data,omitempty"`
}

func NewRemoteSDP(sdpStr string, data map[string]string) *SDP {
	return &SDP{
		SDP:        sdpStr,
		RemoteData: data,
	}
}

func NewLocalSDP(sdpStr string, data map[string]string) *SDP {
	return &SDP{
		SDP:       sdpStr,
		LocalData: data,
	}
}
