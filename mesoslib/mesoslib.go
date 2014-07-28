package mesoslib

import (
	"net"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
)

type MesosLib struct {
	master string
	log    *logrus.Logger
	ip     string
	port   int

	events events
}

func NewMesosLib(master string, log *logrus.Logger) *MesosLib {
	m := &MesosLib{
		log:    log,
		master: master,
		port:   9091,
		events: events{
			mesosproto.Event_REGISTERED: make(chan *mesosproto.Event),
			mesosproto.Event_OFFERS:     make(chan *mesosproto.Event),
		},
	}

	name, err := os.Hostname()
	if err != nil {
		log.Fatalf("Failed to get hostname: %+v", err)
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		log.Fatalf("Failed to get address for hostname %q: %+v", name, err)
	}

	for _, addr := range addrs {
		if m.ip == "" || !strings.HasPrefix(addr, "127") {
			m.ip = addr
		}
	}
	m.initAPI()
	return m
}