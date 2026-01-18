package port

import (
	"errors"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"sakibox/internal/voice"
)

type Entry struct {
	Port    int
	Process string
	PID     int
}

func ListPorts() ([]Entry, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("/bin/sh", "-c", "lsof -iTCP -sTCP:LISTEN -Pn")
	} else {
		cmd = exec.Command("/bin/sh", "-c", "lsof -iTCP -sTCP:LISTEN -Pn")
	}
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\n")
	entries := make([]Entry, 0)
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}
		pid, _ := strconv.Atoi(fields[1])
		name := fields[0]
		addr := fields[len(fields)-1]
		portStr := addr[strings.LastIndex(addr, ":")+1:]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}
		entries = append(entries, Entry{Port: port, Process: name, PID: pid})
	}
	return entries, nil
}

func FindPort(port int) (Entry, error) {
	entries, err := ListPorts()
	if err != nil {
		return Entry{}, err
	}
	for _, entry := range entries {
		if entry.Port == port {
			return entry, nil
		}
	}
	return Entry{}, errors.New(voice.Line("no_results"))
}

func KillByPort(port int) error {
	entry, err := FindPort(port)
	if err != nil {
		return err
	}
	if entry.PID == 0 {
		return errors.New(voice.Line("port_no_process"))
	}
	cmd := exec.Command("/bin/sh", "-c", "kill -9 "+strconv.Itoa(entry.PID))
	return cmd.Run()
}
