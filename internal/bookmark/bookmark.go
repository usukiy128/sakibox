package bookmark

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"sakibox/internal/voice"
)

type Item struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

func dataPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".sakibox", "bookmarks.json"), nil
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
	data, err := json.MarshalIndent([]Item{}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func List() ([]Item, error) {
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
	items := make([]Item, 0)
	if len(data) == 0 {
		return items, nil
	}
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func Add(name, command string) error {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(command) == "" {
		return errors.New(voice.Line("bookmark_invalid_input"))
	}
	items, err := List()
	if err != nil {
		return err
	}
	items = append(items, Item{Name: name, Command: command})
	return save(items)
}

func Get(input string) (string, error) {
	items, err := List()
	if err != nil {
		return "", err
	}
	if index, err := strconv.Atoi(input); err == nil {
		if index <= 0 || index > len(items) {
			return "", errors.New(voice.Line("invalid_index"))
		}
		return items[index-1].Command, nil
	}
	for _, item := range items {
		if item.Name == input {
			return item.Command, nil
		}
	}
	return "", errors.New(voice.Line("bookmark_not_found"))
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

func save(items []Item) error {
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
