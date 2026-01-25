package ssh

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"sakibox/internal/voice"
)

type Server struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Port     int    `json:"port"`
	Password string `json:"password,omitempty"`
}

type LogEntry struct {
	Time    string `json:"time"`
	Name    string `json:"name"`
	Host    string `json:"host"`
	Action  string `json:"action"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type Command struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

func dataPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".sakibox", "ssh.json"), nil
}

func logPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".sakibox", "ssh_logs.json"), nil
}

func commandPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".sakibox", "ssh_commands.json"), nil
}

func ensureData() error {
	path, err := dataPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent([]Server{}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func ensureLogs() error {
	path, err := logPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent([]LogEntry{}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func ensureCommands() error {
	path, err := commandPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent([]Command{}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func List() ([]Server, error) {
	if err := ensureData(); err != nil {
		return nil, err
	}
	path, err := dataPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	items := make([]Server, 0)
	if len(data) == 0 {
		return items, nil
	}
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func Add(server Server) error {
	if strings.TrimSpace(server.Name) == "" || strings.TrimSpace(server.Host) == "" || strings.TrimSpace(server.User) == "" {
		return errors.New(voice.Line("ssh_invalid_input"))
	}
	if server.Port <= 0 {
		return errors.New(voice.Line("ssh_invalid_port"))
	}
	items, err := List()
	if err != nil {
		return err
	}
	items = append(items, server)
	return save(items)
}

func Get(input string) (Server, error) {
	items, err := List()
	if err != nil {
		return Server{}, err
	}
	if index, err := strconv.Atoi(input); err == nil {
		if index <= 0 || index > len(items) {
			return Server{}, errors.New(voice.Line("invalid_index"))
		}
		return items[index-1], nil
	}
	for _, item := range items {
		if item.Name == input {
			return item, nil
		}
	}
	return Server{}, errors.New(voice.Line("ssh_not_found"))
}

func Delete(index int) error {
	items, err := List()
	if err != nil {
		return err
	}
	if index <= 0 || index > len(items) {
		return errors.New(voice.Line("invalid_index"))
	}
	items = append(items[:index-1], items[index:]...)
	return save(items)
}

func save(items []Server) error {
	path, err := dataPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func AddLog(entry LogEntry) error {
	if err := ensureLogs(); err != nil {
		return err
	}
	path, err := logPath()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	items := make([]LogEntry, 0)
	if len(data) > 0 {
		if err := json.Unmarshal(data, &items); err != nil {
			return err
		}
	}
	items = append(items, entry)
	payload, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, payload, 0644)
}

func ListLogs() ([]LogEntry, error) {
	if err := ensureLogs(); err != nil {
		return nil, err
	}
	path, err := logPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	items := make([]LogEntry, 0)
	if len(data) == 0 {
		return items, nil
	}
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func ListCommands() ([]Command, error) {
	if err := ensureCommands(); err != nil {
		return nil, err
	}
	path, err := commandPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	items := make([]Command, 0)
	if len(data) == 0 {
		return items, nil
	}
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func AddCommand(item Command) error {
	if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Command) == "" {
		return errors.New(voice.Line("ssh_cmd_invalid_input"))
	}
	items, err := ListCommands()
	if err != nil {
		return err
	}
	items = append(items, item)
	return saveCommands(items)
}

func GetCommand(input string) (Command, error) {
	items, err := ListCommands()
	if err != nil {
		return Command{}, err
	}
	if index, err := strconv.Atoi(input); err == nil {
		if index <= 0 || index > len(items) {
			return Command{}, errors.New(voice.Line("invalid_index"))
		}
		return items[index-1], nil
	}
	for _, item := range items {
		if item.Name == input {
			return item, nil
		}
	}
	return Command{}, errors.New(voice.Line("ssh_cmd_not_found"))
}

func DeleteCommand(index int) error {
	items, err := ListCommands()
	if err != nil {
		return err
	}
	if index <= 0 || index > len(items) {
		return errors.New(voice.Line("invalid_index"))
	}
	items = append(items[:index-1], items[index:]...)
	return saveCommands(items)
}

func saveCommands(items []Command) error {
	path, err := commandPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func NewLog(server Server, action string, err error) LogEntry {
	entry := LogEntry{
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Name:    server.Name,
		Host:    server.Host,
		Action:  action,
		Success: err == nil,
	}
	if err != nil {
		entry.Error = err.Error()
	}
	return entry
}
