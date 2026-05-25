//go:build windows
// +build windows

package platform

import (
	"xprinter/internal/server"

	"golang.org/x/sys/windows/svc"
)

type service struct{}

func (m *service) Execute(
	args []string,
	r <-chan svc.ChangeRequest,
	s chan<- svc.Status,
) (bool, uint32) {

	s <- svc.Status{State: svc.StartPending}

	go server.Run()

	s <- svc.Status{
		State:   svc.Running,
		Accepts: svc.AcceptStop | svc.AcceptShutdown,
	}

	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Stop, svc.Shutdown:
				s <- svc.Status{State: svc.StopPending}
				return false, 0
			}
		}
	}
}

func RunWindowsService() {
	svc.Run("xprinter", &service{})
}
