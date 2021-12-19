package webrtcsignalingserver

import (
	"bytes"
	"crypto/md5"
	"testing"

	"github.com/pion/webrtc/v3"
)

func TestSDPRequest_DecodeSDP(t *testing.T) {
	type fields struct {
		Id   string
		SDP  string
		Data map[string]string
	}

	sdp1Base64 := "eyJ0eXBlIjoib2ZmZXIiLCJzZHAiOiJ2PTBcclxubz1tb3ppbGxhLi4uVEhJU19JU19TRFBBUlRBLTk1LjAuMSA3MjAzNTUzNjgwODM5MjQ1MDkwIDAgSU4gSVA0IDAuMC4wLjBcclxucz0tXHJcbnQ9MCAwXHJcbmE9c2VuZHJlY3ZcclxuYT1maW5nZXJwcmludDpzaGEtMjU2IDc2OjVFOjI0OjI5OjBFOkUwOjREOjMwOjcyOjJBOjQ1OjdGOkNBOjUxOjg5OjMzOkJBOkUwOkE0Ojk5OkIwOjEzOjVCOjA4OkY5OjAwOkI0OjRCOjVGOjM0OjIwOkQzXHJcbmE9Z3JvdXA6QlVORExFIDBcclxuYT1pY2Utb3B0aW9uczp0cmlja2xlXHJcbmE9bXNpZC1zZW1hbnRpYzpXTVMgKlxyXG5tPXZpZGVvIDUyMDg2IFVEUC9UTFMvUlRQL1NBVlBGIDEyMCAxMjQgMTIxIDEyNSAxMjYgMTI3IDk3IDk4XHJcbmM9SU4gSVA0IDUuMTIxLjgwLjYzXHJcbmE9Y2FuZGlkYXRlOjAgMSBVRFAgMjEyMjE4NzAwNyAxOTIuMTY4LjEuMTA3IDUyNTY4IHR5cCBob3N0XHJcbmE9Y2FuZGlkYXRlOjIgMSBVRFAgMjEyMjI1MjU0MyAxNzIuMjguMjA4LjEgNTI1NjkgdHlwIGhvc3RcclxuYT1jYW5kaWRhdGU6NCAxIFRDUCAyMTA1NDU4OTQzIDE5Mi4xNjguMS4xMDcgOSB0eXAgaG9zdCB0Y3B0eXBlIGFjdGl2ZVxyXG5hPWNhbmRpZGF0ZTo1IDEgVENQIDIxMDU1MjQ0NzkgMTcyLjI4LjIwOC4xIDkgdHlwIGhvc3QgdGNwdHlwZSBhY3RpdmVcclxuYT1jYW5kaWRhdGU6MCAyIFVEUCAyMTIyMTg3MDA2IDE5Mi4xNjguMS4xMDcgNTI1NzAgdHlwIGhvc3RcclxuYT1jYW5kaWRhdGU6MiAyIFVEUCAyMTIyMjUyNTQyIDE3Mi4yOC4yMDguMSA1MjU3MSB0eXAgaG9zdFxyXG5hPWNhbmRpZGF0ZTo0IDIgVENQIDIxMDU0NTg5NDIgMTkyLjE2OC4xLjEwNyA5IHR5cCBob3N0IHRjcHR5cGUgYWN0aXZlXHJcbmE9Y2FuZGlkYXRlOjUgMiBUQ1AgMjEwNTUyNDQ3OCAxNzIuMjguMjA4LjEgOSB0eXAgaG9zdCB0Y3B0eXBlIGFjdGl2ZVxyXG5hPWNhbmRpZGF0ZToxIDEgVURQIDE2ODU5ODczMjcgNS4xMjEuODAuNjMgNTIwODYgdHlwIHNyZmx4IHJhZGRyIDE5Mi4xNjguMS4xMDcgcnBvcnQgNTI1NjhcclxuYT1jYW5kaWRhdGU6MSAyIFVEUCAxNjg1OTg3MzI2IDUuMTIxLjgwLjYzIDUyMDg3IHR5cCBzcmZseCByYWRkciAxOTIuMTY4LjEuMTA3IHJwb3J0IDUyNTcwXHJcbmE9c2VuZHJlY3ZcclxuYT1lbmQtb2YtY2FuZGlkYXRlc1xyXG5hPWV4dG1hcDozIHVybjppZXRmOnBhcmFtczpydHAtaGRyZXh0OnNkZXM6bWlkXHJcbmE9ZXh0bWFwOjQgaHR0cDovL3d3dy53ZWJydGMub3JnL2V4cGVyaW1lbnRzL3J0cC1oZHJleHQvYWJzLXNlbmQtdGltZVxyXG5hPWV4dG1hcDo1IHVybjppZXRmOnBhcmFtczpydHAtaGRyZXh0OnRvZmZzZXRcclxuYT1leHRtYXA6Ni9yZWN2b25seSBodHRwOi8vd3d3LndlYnJ0Yy5vcmcvZXhwZXJpbWVudHMvcnRwLWhkcmV4dC9wbGF5b3V0LWRlbGF5XHJcbmE9ZXh0bWFwOjcgaHR0cDovL3d3dy5pZXRmLm9yZy9pZC9kcmFmdC1ob2xtZXItcm1jYXQtdHJhbnNwb3J0LXdpZGUtY2MtZXh0ZW5zaW9ucy0wMVxyXG5hPWZtdHA6MTI2IHByb2ZpbGUtbGV2ZWwtaWQ9NDJlMDFmO2xldmVsLWFzeW1tZXRyeS1hbGxvd2VkPTE7cGFja2V0aXphdGlvbi1tb2RlPTFcclxuYT1mbXRwOjk3IHByb2ZpbGUtbGV2ZWwtaWQ9NDJlMDFmO2xldmVsLWFzeW1tZXRyeS1hbGxvd2VkPTFcclxuYT1mbXRwOjEyMCBtYXgtZnM9MTIyODg7bWF4LWZyPTYwXHJcbmE9Zm10cDoxMjQgYXB0PTEyMFxyXG5hPWZtdHA6MTIxIG1heC1mcz0xMjI4ODttYXgtZnI9NjBcclxuYT1mbXRwOjEyNSBhcHQ9MTIxXHJcbmE9Zm10cDoxMjcgYXB0PTEyNlxyXG5hPWZtdHA6OTggYXB0PTk3XHJcbmE9aWNlLXB3ZDo3NWUxOTYwMGVjZDA1NDM3NDFkY2Y3ZWVhOWU0Y2VlM1xyXG5hPWljZS11ZnJhZzphNmM4NzQxYlxyXG5hPW1pZDowXHJcbmE9bXNpZDp7NGZlMGE5ZmUtMjYyNS00NTcwLTg2N2QtOTJiZDU0N2ZlNDdlfSB7MzZiM2UzNGItNmE3OC00NTRmLTk4M2EtZWUyMzFmNmE4OWI0fVxyXG5hPXJ0Y3A6NTIwODcgSU4gSVA0IDUuMTIxLjgwLjYzXHJcbmE9cnRjcC1mYjoxMjAgbmFja1xyXG5hPXJ0Y3AtZmI6MTIwIG5hY2sgcGxpXHJcbmE9cnRjcC1mYjoxMjAgY2NtIGZpclxyXG5hPXJ0Y3AtZmI6MTIwIGdvb2ctcmVtYlxyXG5hPXJ0Y3AtZmI6MTIwIHRyYW5zcG9ydC1jY1xyXG5hPXJ0Y3AtZmI6MTIxIG5hY2tcclxuYT1ydGNwLWZiOjEyMSBuYWNrIHBsaVxyXG5hPXJ0Y3AtZmI6MTIxIGNjbSBmaXJcclxuYT1ydGNwLWZiOjEyMSBnb29nLXJlbWJcclxuYT1ydGNwLWZiOjEyMSB0cmFuc3BvcnQtY2NcclxuYT1ydGNwLWZiOjEyNiBuYWNrXHJcbmE9cnRjcC1mYjoxMjYgbmFjayBwbGlcclxuYT1ydGNwLWZiOjEyNiBjY20gZmlyXHJcbmE9cnRjcC1mYjoxMjYgZ29vZy1yZW1iXHJcbmE9cnRjcC1mYjoxMjYgdHJhbnNwb3J0LWNjXHJcbmE9cnRjcC1mYjo5NyBuYWNrXHJcbmE9cnRjcC1mYjo5NyBuYWNrIHBsaVxyXG5hPXJ0Y3AtZmI6OTcgY2NtIGZpclxyXG5hPXJ0Y3AtZmI6OTcgZ29vZy1yZW1iXHJcbmE9cnRjcC1mYjo5NyB0cmFuc3BvcnQtY2NcclxuYT1ydGNwLW11eFxyXG5hPXJ0Y3AtcnNpemVcclxuYT1ydHBtYXA6MTIwIFZQOC85MDAwMFxyXG5hPXJ0cG1hcDoxMjQgcnR4LzkwMDAwXHJcbmE9cnRwbWFwOjEyMSBWUDkvOTAwMDBcclxuYT1ydHBtYXA6MTI1IHJ0eC85MDAwMFxyXG5hPXJ0cG1hcDoxMjYgSDI2NC85MDAwMFxyXG5hPXJ0cG1hcDoxMjcgcnR4LzkwMDAwXHJcbmE9cnRwbWFwOjk3IEgyNjQvOTAwMDBcclxuYT1ydHBtYXA6OTggcnR4LzkwMDAwXHJcbmE9c2V0dXA6YWN0cGFzc1xyXG5hPXNzcmM6MTk3MjYyNTcxNSBjbmFtZTp7OWMwYTI2YzMtZDI4My00MmNhLTk3NjEtNzQ2YmI3ZGM3Y2ZlfVxyXG5hPXNzcmM6Njk3MjY1OTU3IGNuYW1lOns5YzBhMjZjMy1kMjgzLTQyY2EtOTc2MS03NDZiYjdkYzdjZmV9XHJcbmE9c3NyYy1ncm91cDpGSUQgMTk3MjYyNTcxNSA2OTcyNjU5NTdcclxuIn0="

	tests := []struct {
		name    string
		fields  fields
		wantSdp webrtc.SessionDescription
		wantErr bool
	}{
		{
			name: "sdp1",
			fields: fields{
				Id:   "publisher",
				SDP:  sdp1Base64,
				Data: nil,
			},
			wantSdp: webrtc.SessionDescription{
				Type: webrtc.SDPTypeOffer,
				SDP: `v=0
o=mozilla...THIS_IS_SDPARTA-95.0.1 7203553680839245090 0 IN IP4 0.0.0.0
s=-
t=0 0
a=sendrecv
a=fingerprint:sha-256 76:5E:24:29:0E:E0:4D:30:72:2A:45:7F:CA:51:89:33:BA:E0:A4:99:B0:13:5B:08:F9:00:B4:4B:5F:34:20:D3
a=group:BUNDLE 0
a=ice-options:trickle
a=msid-semantic:WMS *
m=video 52086 UDP/TLS/RTP/SAVPF 120 124 121 125 126 127 97 98
c=IN IP4 5.121.80.63
a=candidate:0 1 UDP 2122187007 192.168.1.107 52568 typ host
a=candidate:2 1 UDP 2122252543 172.28.208.1 52569 typ host
a=candidate:4 1 TCP 2105458943 192.168.1.107 9 typ host tcptype active
a=candidate:5 1 TCP 2105524479 172.28.208.1 9 typ host tcptype active
a=candidate:0 2 UDP 2122187006 192.168.1.107 52570 typ host
a=candidate:2 2 UDP 2122252542 172.28.208.1 52571 typ host
a=candidate:4 2 TCP 2105458942 192.168.1.107 9 typ host tcptype active
a=candidate:5 2 TCP 2105524478 172.28.208.1 9 typ host tcptype active
a=candidate:1 1 UDP 1685987327 5.121.80.63 52086 typ srflx raddr 192.168.1.107 rport 52568
a=candidate:1 2 UDP 1685987326 5.121.80.63 52087 typ srflx raddr 192.168.1.107 rport 52570
a=sendrecv
a=end-of-candidates
a=extmap:3 urn:ietf:params:rtp-hdrext:sdes:mid
a=extmap:4 http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time
a=extmap:5 urn:ietf:params:rtp-hdrext:toffset
a=extmap:6/recvonly http://www.webrtc.org/experiments/rtp-hdrext/playout-delay
a=extmap:7 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01
a=fmtp:126 profile-level-id=42e01f;level-asymmetry-allowed=1;packetization-mode=1
a=fmtp:97 profile-level-id=42e01f;level-asymmetry-allowed=1
a=fmtp:120 max-fs=12288;max-fr=60
a=fmtp:124 apt=120
a=fmtp:121 max-fs=12288;max-fr=60
a=fmtp:125 apt=121
a=fmtp:127 apt=126
a=fmtp:98 apt=97
a=ice-pwd:75e19600ecd0543741dcf7eea9e4cee3
a=ice-ufrag:a6c8741b
a=mid:0
a=msid:{4fe0a9fe-2625-4570-867d-92bd547fe47e} {36b3e34b-6a78-454f-983a-ee231f6a89b4}
a=rtcp:52087 IN IP4 5.121.80.63
a=rtcp-fb:120 nack
a=rtcp-fb:120 nack pli
a=rtcp-fb:120 ccm fir
a=rtcp-fb:120 goog-remb
a=rtcp-fb:120 transport-cc
a=rtcp-fb:121 nack
a=rtcp-fb:121 nack pli
a=rtcp-fb:121 ccm fir
a=rtcp-fb:121 goog-remb
a=rtcp-fb:121 transport-cc
a=rtcp-fb:126 nack
a=rtcp-fb:126 nack pli
a=rtcp-fb:126 ccm fir
a=rtcp-fb:126 goog-remb
a=rtcp-fb:126 transport-cc
a=rtcp-fb:97 nack
a=rtcp-fb:97 nack pli
a=rtcp-fb:97 ccm fir
a=rtcp-fb:97 goog-remb
a=rtcp-fb:97 transport-cc
a=rtcp-mux
a=rtcp-rsize
a=rtpmap:120 VP8/90000
a=rtpmap:124 rtx/90000
a=rtpmap:121 VP9/90000
a=rtpmap:125 rtx/90000
a=rtpmap:126 H264/90000
a=rtpmap:127 rtx/90000
a=rtpmap:97 H264/90000
a=rtpmap:98 rtx/90000
a=setup:actpass
a=ssrc:1972625715 cname:{9c0a26c3-d283-42ca-9761-746bb7dc7cfe}
a=ssrc:697265957 cname:{9c0a26c3-d283-42ca-9761-746bb7dc7cfe}
a=ssrc-group:FID 1972625715 697265957
`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := &sDPRequest{
				Id:   tt.fields.Id,
				SDP:  tt.fields.SDP,
				Data: tt.fields.Data,
			}
			gotSdp, err := sr.DecodeSDP()
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeSDP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			md5Got := md5.Sum([]byte(gotSdp.SDP))
			md5Want := md5.Sum([]byte(gotSdp.SDP))
			if !bytes.Equal(md5Got[:], md5Want[:]) {
				t.Errorf("DecodeSDP() gotSdp = %v, want %v", gotSdp, tt.wantSdp)
			}
		})
	}
}
