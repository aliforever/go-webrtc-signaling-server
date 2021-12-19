package internal

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"sync"

	"github.com/pion/webrtc/v3"
)

type SdpStorage struct {
	listeners map[string]chan *SDP
	storage   map[string]*SDP
	callbacks map[string]func(remoteSDP *SDP) (localSDP *webrtc.SessionDescription, data map[string]string)

	// Lockers
	listenersM sync.Mutex
	storageM   sync.Mutex
	callbacksM sync.Mutex
}

func NewSDPStorage() (ss *SdpStorage) {
	ss = &SdpStorage{
		listeners: map[string]chan *SDP{},
	}
	return
}

func (ss *SdpStorage) AddSDPCallback(id string, callback func(remoteSDP *SDP) (localSDP *webrtc.SessionDescription, data map[string]string)) (err error) {
	ss.callbacksM.Lock()
	defer ss.callbacksM.Unlock()

	if _, exists := ss.callbacks[id]; exists {
		err = errors.New("id_exists")
		return
	}

	ss.callbacks[id] = callback
	return
}

func (ss *SdpStorage) GetLocalSDPFromCallback(id string, remoteSDP string, remoteSDPData map[string]string) (s *SDP, err error) {
	ss.callbacksM.Lock()
	defer ss.callbacksM.Unlock()

	callback, exists := ss.callbacks[id]
	if !exists {
		err = errors.New("id_exists")
		return
	}

	// TODO: defer delete(ss.callback[id]) - Decide to Remove Item from Callback or Not

	localSDP, data := callback(NewRemoteSDP(remoteSDP, remoteSDPData))

	var sdpJson []byte
	sdpJson, err = json.Marshal(localSDP)
	if err != nil {
		return
	}

	var sdpBase64 string
	sdpBase64 = base64.StdEncoding.EncodeToString(sdpJson)

	s = NewLocalSDP(sdpBase64, data)
	return
}

func (ss *SdpStorage) AddSDPListener(id string) (sdpListener chan *SDP, err error) {
	ss.listenersM.Lock()
	defer ss.listenersM.Unlock()

	if _, exists := ss.listeners[id]; exists {
		err = errors.New("listener_exists")
		return
	}

	ss.listeners[id] = make(chan *SDP)
	sdpListener = ss.listeners[id]

	return
}

func (ss *SdpStorage) InformLocalSDPListener(id, sdp string, data map[string]string) (err error) {
	ss.listenersM.Lock()
	defer ss.listenersM.Unlock()

	listener, exists := ss.listeners[id]
	if !exists {
		err = errors.New("listener_does_not_exist")
		return
	}

	listener <- NewRemoteSDP(sdp, data)

	close(listener)
	delete(ss.listeners, id)

	return
}

func (ss *SdpStorage) AddSDPToStorage(id, sdp string, data map[string]string) (err error) {
	ss.storageM.Lock()
	defer ss.storageM.Unlock()

	if _, exists := ss.storage[id]; exists {
		err = errors.New("sdp_exists")
		return
	}

	ss.storage[id] = NewRemoteSDP(sdp, data)
	return
}

func (ss *SdpStorage) GetSDPFromStorage(id string) (sdp *SDP, err error) {
	ss.storageM.Lock()
	defer ss.storageM.Unlock()

	var exists bool
	if sdp, exists = ss.storage[id]; !exists {
		err = errors.New("sdp_does_not_exists")
		return
	}

	return
}
