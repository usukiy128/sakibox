package install

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"sakibox/internal/voice"
)

type Manager string

const (
	managerApt    Manager = "apt"
	managerYum    Manager = "yum"
	managerDnf    Manager = "dnf"
	managerPacman Manager = "pacman"
)

var defaultTools = []string{
	"curl",
	"wget",
	"git",
	"make",
	"unzip",
	"zip",
	"jq",
	"htop",
	"tree",
	"lsof",
}

func DetectManager() (Manager, error) {
	if runtime.GOOS != "linux" {
		return "", errors.New(voice.Line("install_linux_only"))
	}
	for _, mgr := range []Manager{managerApt, managerDnf, managerYum, managerPacman} {
		if _, err := exec.LookPath(string(mgr)); err == nil {
			return mgr, nil
		}
	}
	return "", errors.New(voice.Line("install_no_manager"))
}

func InstallDefaults() (string, error) {
	mgr, err := DetectManager()
	if err != nil {
		return "", err
	}
	cmd, err := installCommand(mgr, defaultTools)
	if err != nil {
		return "", err
	}
	return cmd, runShell(cmd)
}

func SearchAndInstall(keyword string) (string, error) {
	mgr, err := DetectManager()
	if err != nil {
		return "", err
	}
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return "", errors.New(voice.Line("install_invalid_keyword"))
	}
	pkg, err := resolvePackage(mgr, keyword)
	if err != nil {
		return "", err
	}
	cmd, err := installCommand(mgr, []string{pkg})
	if err != nil {
		return "", err
	}
	return cmd, runShell(cmd)
}

func installCommand(mgr Manager, packages []string) (string, error) {
	if len(packages) == 0 {
		return "", errors.New(voice.Line("install_no_packages"))
	}
	switch mgr {
	case managerApt:
		return fmt.Sprintf("sudo apt update && sudo apt install -y %s", strings.Join(packages, " ")), nil
	case managerYum:
		return fmt.Sprintf("sudo yum install -y %s", strings.Join(packages, " ")), nil
	case managerDnf:
		return fmt.Sprintf("sudo dnf install -y %s", strings.Join(packages, " ")), nil
	case managerPacman:
		return fmt.Sprintf("sudo pacman -Sy --noconfirm %s", strings.Join(packages, " ")), nil
	default:
		return "", errors.New(voice.Line("install_unsupported_manager"))
	}
}

func resolvePackage(mgr Manager, keyword string) (string, error) {
	switch mgr {
	case managerApt:
		return searchByApt(keyword)
	case managerYum:
		return searchByYum(keyword)
	case managerDnf:
		return searchByDnf(keyword)
	case managerPacman:
		return searchByPacman(keyword)
	default:
		return "", errors.New(voice.Line("install_unsupported_manager"))
	}
}

func searchByApt(keyword string) (string, error) {
	output, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("apt-cache search %q", keyword)).Output()
	if err != nil {
		return "", err
	}
	return pickFirstPackage(string(output))
}

func searchByYum(keyword string) (string, error) {
	output, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("yum search -q %q", keyword)).Output()
	if err != nil {
		return "", err
	}
	return pickFirstPackage(string(output))
}

func searchByDnf(keyword string) (string, error) {
	output, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("dnf search -q %q", keyword)).Output()
	if err != nil {
		return "", err
	}
	return pickFirstPackage(string(output))
}

func searchByPacman(keyword string) (string, error) {
	output, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("pacman -Ss %q", keyword)).Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		parts := strings.Split(fields[0], "/")
		if len(parts) == 2 {
			return parts[1], nil
		}
	}
	return "", errors.New(voice.Line("install_no_match"))
}

func pickFirstPackage(output string) (string, error) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		return fields[0], nil
	}
	return "", errors.New(voice.Line("install_no_match"))
}

func runShell(command string) error {
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
