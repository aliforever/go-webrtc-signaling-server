package webrtcsignalingserver

import (
	"errors"
	"sync"
)

type sdpStorage struct {
	listeners map[string]*Listener
	storage   map[string]*SDPClient

	// Lockers
	listenersM sync.Mutex
	storageM   sync.Mutex
}

func newSDPStorage() (ss *sdpStorage) {
	ss = &sdpStorage{
		listeners: map[string]*Listener{},
		storage:   map[string]*SDPClient{},
	}
	return
}

func (ss *sdpStorage) AddSDPListener(id string) (l *Listener, err error) {
	ss.listenersM.Lock()
	defer ss.listenersM.Unlock()

	if _, exists := ss.listeners[id]; exists {
		err = errors.New("listener_exists")
		return
	}

	l = newListener()
	ss.listeners[id] = l

	return
}

func (ss *sdpStorage) GetSDPListener(id string) (l *Listener, err error) {
	ss.listenersM.Lock()
	defer ss.listenersM.Unlock()

	var exists bool
	l, exists = ss.listeners[id]
	if !exists {
		err = errors.New("listener_does_not_exist")
		return
	}

	// Delete Listener from Storage to prevent others access the same listener
	delete(ss.listeners, id)
	return
}

func (ss *sdpStorage) AddSDPToStorage(id, sdp string, data map[string]string) (err error) {
	ss.storageM.Lock()
	defer ss.storageM.Unlock()

	if _, exists := ss.storage[id]; exists {
		err = errors.New("sdp_exists")
		return
	}

	var remoteSdp *SDPClient
	remoteSdp, err = newClientSDP(sdp, data)
	if err != nil {
		return
	}

	ss.storage[id] = remoteSdp
	return
}

func (ss *sdpStorage) GetSDPFromStorage(id string) (sdp *SDPClient, err error) {
	ss.storageM.Lock()
	defer ss.storageM.Unlock()

	var exists bool
	if sdp, exists = ss.storage[id]; !exists {
		err = errors.New("sdp_does_not_exists")
		return
	}

	return
}
