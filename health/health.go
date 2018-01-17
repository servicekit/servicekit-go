package health

import (
	"fmt"
	"net/http"

	"github.com/servicekit/servicekit-go/logger"
)

// ServiceState represents the state of service
type ServiceState string

const (
	// ServiceStateOK represents the state of service that is OK
	ServiceStateOK ServiceState = "OK"
	// ServiceStateBusy represents the state of service that is Busy
	ServiceStateBusy ServiceState = "Busy"
	// ServiceStateIdling represents the state of service that is Idle
	ServiceStateIdling ServiceState = "Idle"

	// ServiceStateUnavailable represents the state of service that is Unavailable
	ServiceStateUnavailable ServiceState = "Unavailable"
)

// Health represents the state info of service
type Health struct {
	host string
	port int
	path string

	state  ServiceState
	reason string
	c      chan ServiceState

	log *logger.Logger
}

// NewHealth returns a Health
func NewHealth(host string, port int, path string, log *logger.Logger) *Health {
	h := &Health{
		host: host,
		port: port,
		path: path,

		state: ServiceStateUnavailable,
		c:     make(chan ServiceState),

		log: log,
	}

	go h.start()
	go h.serve()

	log.Info("health: health started")

	return h
}

// start update the state of service periodically
func (h *Health) start() {
	for {
		s := <-h.c
		oldState := h.state
		h.state = s

		if oldState != h.state {
			h.log.Infof("health: state changed. %v -> %v", oldState, h.state)
		}
	}

}

// GetChan returns a write-only channel that you can pass new state to it
func (h *Health) GetChan() chan<- ServiceState {
	return h.c
}

// handler is a http hander
func (h *Health) handler(w http.ResponseWriter, req *http.Request) {
	if h.state == ServiceStateUnavailable {
		w.WriteHeader(500)
	}
}

// serve serve a http server
func (h *Health) serve() {
	http.HandleFunc(h.path, h.handler)
	http.ListenAndServe(fmt.Sprintf("%s:%d", h.host, h.port), nil)
}
