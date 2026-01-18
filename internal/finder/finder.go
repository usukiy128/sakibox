package finder

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"sakibox/config"
	"sakibox/internal/voice"
)

type Result struct {
	Path     string
	Size     string
	Modified string
	IsDir    bool
}

type ContentResult struct {
	Path    string
	Line    int
	Content string
}

func FindByName(root, keyword string, exact bool) ([]Result, error) {
	return FindByNameWithExt(root, keyword, "", exact)
}

func FindByNameWithExt(root, keyword, ext string, exact bool) ([]Result, error) {
	ext = strings.TrimSpace(ext)
	if ext != "" && !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return walkFiles(root, func(path string, info os.FileInfo) (bool, error) {
		if info.IsDir() && ext != "" {
			return false, nil
		}
		if ext != "" && !strings.EqualFold(filepath.Ext(info.Name()), ext) {
			return false, nil
		}
		if exact {
			return info.Name() == keyword, nil
		}
		return strings.Contains(info.Name(), keyword), nil
	})
}

func FindByExt(root, ext string) ([]Result, error) {
	ext = strings.TrimSpace(ext)
	if ext != "" && !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return walkFiles(root, func(path string, info os.FileInfo) (bool, error) {
		if info.IsDir() {
			return false, nil
		}
		return strings.EqualFold(filepath.Ext(info.Name()), ext), nil
	})
}

func FindByContent(root, query string) ([]ContentResult, error) {
	return FindByContentWithExt(root, query, "")
}

func FindByContentWithExt(root, query, ext string) ([]ContentResult, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	results := make([]ContentResult, 0)
	count := 0
	forEachFile(root, cfg.IgnoreDirs, func(path string, info os.FileInfo) error {
		count++
		if count%200 == 0 {
			fmt.Println("  " + voice.Linef("searching_file", path))
		}
		if ext != "" && !strings.EqualFold(filepath.Ext(info.Name()), ext) {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			text := scanner.Text()
			if strings.Contains(text, query) {
				results = append(results, ContentResult{Path: path, Line: lineNum, Content: strings.TrimSpace(text)})
			}
		}
		return nil
	})
	return results, nil
}

func FindGlobal(query, ext string) ([]ContentResult, error) {
	path, err := GlobalSearchPath()
	if err != nil {
		return nil, err
	}
	return FindByContentWithExt(path, query, ext)
}

func GlobalSearchPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home, nil
}

func FindBySize(root, condition, sizeInput string) ([]Result, error) {
	threshold, err := parseSize(sizeInput)
	if err != nil {
		return nil, err
	}
	return walkFiles(root, func(path string, info os.FileInfo) (bool, error) {
		size := info.Size()
		switch condition {
		case "1":
			return size > threshold, nil
		case "2":
			return size < threshold, nil
		case "3":
			return size == threshold, nil
		default:
			return false, errors.New(voice.Line("invalid_option"))
		}
	})
}

func FindByTime(root, condition, daysInput string) ([]Result, error) {
	days, err := strconv.Atoi(daysInput)
	if err != nil {
		return nil, errors.New(voice.Line("invalid_days"))
	}
	cutoff := time.Now().AddDate(0, 0, -days)
	return walkFiles(root, func(path string, info os.FileInfo) (bool, error) {
		mod := info.ModTime()
		switch condition {
		case "1":
			return mod.After(cutoff), nil
		case "2":
			return mod.Before(cutoff), nil
		default:
			return false, errors.New(voice.Line("invalid_option"))
		}
	})
}

func walkFiles(root string, matcher func(path string, info os.FileInfo) (bool, error)) ([]Result, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	results := make([]Result, 0)
	forEachDirThenFiles(root, cfg.IgnoreDirs, func(path string, info os.FileInfo) error {
		match, err := matcher(path, info)
		if err != nil {
			return err
		}
		if match {
			results = append(results, Result{
				Path:     path,
				Size:     formatSize(info.Size()),
				Modified: info.ModTime().Format("2006-01-02 15:04"),
				IsDir:    info.IsDir(),
			})
		}
		return nil
	})
	return results, nil
}

func forEachFile(root string, ignore []string, handle func(path string, info os.FileInfo) error) {
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") {
				if path != root {
					return filepath.SkipDir
				}
			}
			for _, skip := range ignore {
				if name == skip {
					return filepath.SkipDir
				}
			}
			return nil
		}
		_ = handle(path, info)
		return nil
	})
}

func forEachDirThenFiles(root string, ignore []string, handle func(path string, info os.FileInfo) error) {
	type searchEntry struct {
		path string
		info os.FileInfo
	}
	var dirs []searchEntry
	var files []searchEntry

	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if entry.IsDir() {
			name := entry.Name()
			if strings.HasPrefix(name, ".") {
				if path != root {
					return filepath.SkipDir
				}
			}
			for _, skip := range ignore {
				if name == skip {
					return filepath.SkipDir
				}
			}
			if path != root {
				info, err := entry.Info()
				if err == nil {
					dirs = append(dirs, searchEntry{path: path, info: info})
				}
			}
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return nil
		}
		files = append(files, searchEntry{path: path, info: info})
		return nil
	})

	for _, item := range dirs {
		_ = handle(item.path, item.info)
	}
	for _, item := range files {
		_ = handle(item.path, item.info)
	}
}

func parseSize(input string) (int64, error) {
	input = strings.TrimSpace(strings.ToUpper(input))
	if input == "" {
		return 0, errors.New(voice.Line("invalid_size"))
	}
	unit := input[len(input)-1:]
	valueStr := input
	var multiplier int64 = 1
	if unit == "K" || unit == "M" || unit == "G" {
		valueStr = input[:len(input)-1]
		switch unit {
		case "K":
			multiplier = 1024
		case "M":
			multiplier = 1024 * 1024
		case "G":
			multiplier = 1024 * 1024 * 1024
		}
	}
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, errors.New(voice.Line("invalid_size"))
	}
	return int64(value * float64(multiplier)), nil
}

func formatSize(size int64) string {
	if size > 1024*1024*1024 {
		return fmt.Sprintf("%.1fG", float64(size)/1024/1024/1024)
	}
	if size > 1024*1024 {
		return fmt.Sprintf("%.1fM", float64(size)/1024/1024)
	}
	if size > 1024 {
		return fmt.Sprintf("%.1fK", float64(size)/1024)
	}
	return fmt.Sprintf("%dB", size)
}
