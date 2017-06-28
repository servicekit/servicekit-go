package health

import (
	"fmt"
	"net/http"

	"github.com/servicekit/servicekit-go/logger"
)

type ServiceState string

const (
	ServiceStateOK     ServiceState = "OK"
	ServiceStateBusy   ServiceState = "Busy"
	ServiceStateIdling ServiceState = "Idle"

	ServiceStateUnavailable ServiceState = "Unavailable"
)

type health struct {
	host string
	port int
	path string

	state  ServiceState
	reason string
	c      chan ServiceState

	log *logger.Logger
}

func NewHealth(host string, port int, path string, log *logger.Logger) *health {
	h := &health{
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

func (h *health) start() {
	for {
		s := <-h.c
		oldState := h.state
		h.state = s

		if oldState != h.state {
			h.log.Infof("health: state changed. %v -> %v", oldState, h.state)
		}
	}

}

func (h *health) GetChan() chan<- ServiceState {
	return h.c
}

func (h *health) handler(w http.ResponseWriter, req *http.Request) {
	if h.state == ServiceStateUnavailable {
		w.WriteHeader(500)
	}
}

func (h *health) serve() {
	http.HandleFunc(h.path, h.handler)
	http.ListenAndServe(fmt.Sprintf("%s:%d", h.host, h.port), nil)
}
