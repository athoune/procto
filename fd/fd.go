package fd

import (
	"fmt"
	"os"
	"strings"
)

type Fd struct {
	Pid int
}

func NewFd(pid int) (fd *Fd, err error) {
	_, err = os.Stat(fmt.Sprintf("/proc/%d", pid))
	if err != nil {
		return nil, err
	}
	return &Fd{pid}, nil
}

func (self *Fd) CountSockets() (int, error) {
	d, err := os.Open(fmt.Sprintf("/proc/%d/fd", self.Pid))
	if err != nil {
		return 0, err
	}
	fds, err := d.Readdirnames(0)
	cpt := 0
	for _, fd := range fds {
		p, err := os.Readlink(fmt.Sprintf("/proc/%d/fd/%s", self.Pid, fd))
		if err != nil {
			return 0, err
		}
		if strings.HasPrefix(p, "socket:[") {
			cpt += 1
		}
	}
	return cpt, nil
}
