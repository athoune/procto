package java

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type AppData struct {
	Pid  int
	User string
}

func FindAppData(user string) (AppData, error) {
	path := filepath.Join(os.TempDir(), "hsperfdata_"+user)
	f, err := os.Open(path)
	if err != nil {
		return AppData{}, err
	}
	p, err := f.Readdirnames(1)
	if err != nil {
		return AppData{}, err
	}
	pid, err := strconv.Atoi(p[0])
	if err != nil {
		return AppData{}, err
	}
	return AppData{pid, user}, nil
}

func (this *AppData) Path() string {
	return filepath.Join(os.TempDir(), "hsperfdata_"+this.User, fmt.Sprintf("%d", this.Pid))
}
