package history

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"sakibox/config"
	"sakibox/internal/voice"
)

type Entry struct {
	Command string
}

func List() ([]Entry, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	historyPath, err := resolveHistoryFile(cfg.HistoryFile)
	if err != nil {
		return nil, err
	}
	if historyPath == "" {
		return []Entry{}, nil
	}
	lines, err := readLines(historyPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Entry{}, nil
		}
		return nil, err
	}
	entries := make([]Entry, 0)
	for i := len(lines) - 1; i >= 0 && len(entries) < cfg.MaxHistory; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		if idx := strings.LastIndex(line, ";"); idx != -1 {
			line = line[idx+1:]
		}
		entries = append(entries, Entry{Command: line})
	}
	return entries, nil
}

func Search(keyword string) ([]Entry, error) {
	entries, err := List()
	if err != nil {
		return nil, err
	}
	matches := make([]Entry, 0)
	for _, entry := range entries {
		if strings.Contains(entry.Command, keyword) {
			matches = append(matches, entry)
		}
	}
	return matches, nil
}

func GetByIndex(index int) (string, error) {
	if index <= 0 {
		return "", errors.New(voice.Line("invalid_index"))
	}
	entries, err := List()
	if err != nil {
		return "", err
	}
	if index > len(entries) {
		return "", errors.New(voice.Line("invalid_index"))
	}
	return entries[index-1].Command, nil
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func resolveHistoryFile(primary string) (string, error) {
	if primary != "" {
		if _, err := os.Stat(primary); err == nil {
			return primary, nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
	}
	candidates, err := defaultHistoryCandidates()
	if err != nil {
		return "", err
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		} else if err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
	}
	return "", nil
}

func defaultHistoryCandidates() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return []string{
		filepath.Join(home, ".zsh_history"),
		filepath.Join(home, ".bash_history"),
		filepath.Join(home, ".config", "fish", "fish_history"),
	}, nil
}
