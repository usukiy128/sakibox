package process

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type Entry struct {
	PID  int
	Name string
	CPU  float64
	Mem  float64
}

const pageSize = 15

func List(page int) ([]Entry, error) {
	output, err := exec.Command("/bin/sh", "-c", "ps -A -o pid,comm").Output()
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
		if len(fields) < 2 {
			continue
		}
		pid, _ := strconv.Atoi(fields[0])
		name := fields[1]
		entries = append(entries, Entry{PID: pid, Name: name})
	}
	start := page * pageSize
	if start >= len(entries) {
		return []Entry{}, nil
	}
	end := start + pageSize
	if end > len(entries) {
		end = len(entries)
	}
	return entries[start:end], nil
}

func Top() ([]Entry, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("/bin/sh", "-c", "ps -A -o pid,comm,%cpu,%mem | sort -k3 -nr | head -n 11")
	} else {
		cmd = exec.Command("/bin/sh", "-c", "ps -A -o pid,comm,%cpu,%mem --sort=-%cpu | head -n 11")
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
		if len(fields) < 4 {
			continue
		}
		pid, _ := strconv.Atoi(fields[0])
		cpu, _ := strconv.ParseFloat(fields[2], 64)
		mem, _ := strconv.ParseFloat(fields[3], 64)
		entries = append(entries, Entry{PID: pid, Name: fields[1], CPU: cpu, Mem: mem})
	}
	return entries, nil
}

func Search(keyword string) ([]Entry, error) {
	output, err := exec.Command("/bin/sh", "-c", "ps -A -o pid,comm").Output()
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
		if len(fields) < 2 {
			continue
		}
		if strings.Contains(fields[1], keyword) {
			pid, _ := strconv.Atoi(fields[0])
			entries = append(entries, Entry{PID: pid, Name: fields[1]})
		}
	}
	return entries, nil
}

func Kill(pid int) error {
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("kill -9 %d", pid))
	return cmd.Run()
}
