package cmd

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/fatih/color"

	"sakibox/config"
	"sakibox/internal/finder"
	"sakibox/internal/voice"
)

func showFinderMenu(reader *bufio.Reader) error {
	for {
		printCyan("[文件查找]")
		printMagenta(voice.Line("finder_intro"))
		fmt.Println("  1. 按文件名查找")
		fmt.Println("  2. 按扩展名查找")
		fmt.Println("  3. 按内容查找")
		fmt.Println("  4. 按大小查找")
		fmt.Println("  5. 按修改时间查找")
		fmt.Println("  6. 全局检索")
		fmt.Println("  0. 返回主菜单")
		fmt.Printf("\n  %s", voice.Line("menu_prompt"))

		choice, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			if err := findByName(reader); err != nil {
				return err
			}
		case "2":
			if err := findByExt(reader); err != nil {
				return err
			}
		case "3":
			if err := findByContent(reader); err != nil {
				return err
			}
		case "4":
			if err := findBySize(reader); err != nil {
				return err
			}
		case "5":
			if err := findByTime(reader); err != nil {
				return err
			}
		case "6":
			if err := findGlobal(reader); err != nil {
				return err
			}
		case "0":
			return nil
		default:
			printRed(voice.Line("invalid_option"))
			if err := waitForEnter(reader); err != nil {
				return err
			}
		}
	}
}

func findByName(reader *bufio.Reader) error {
	path, err := promptPath(reader)
	if err != nil {
		return err
	}
	fmt.Printf("  %s", voice.Line("finder_name_prompt"))
	keyword, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		printRed(voice.Line("invalid_keyword"))
		return waitForEnter(reader)
	}
	fmt.Printf("  %s", voice.Line("finder_match_prompt"))
	matchType, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	exact := strings.TrimSpace(matchType) == "2"
	printYellow(voice.Line("searching"))
	results, err := finder.FindByName(path, keyword, exact)
	if err != nil {
		return err
	}
	return showFinderResults(reader, results)
}

func findByExt(reader *bufio.Reader) error {
	path, err := promptPath(reader)
	if err != nil {
		return err
	}
	fmt.Printf("  %s", voice.Line("finder_ext_prompt"))
	ext, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	printYellow(voice.Line("searching"))
	results, err := finder.FindByExt(path, strings.TrimSpace(ext))
	if err != nil {
		return err
	}
	return showFinderResults(reader, results)
}

func findByContent(reader *bufio.Reader) error {
	path, err := promptPath(reader)
	if err != nil {
		return err
	}
	fmt.Printf("  %s", voice.Line("finder_content_prompt"))
	query, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	printYellow(voice.Line("searching"))
	results, err := finder.FindByContent(path, strings.TrimSpace(query))
	if err != nil {
		return err
	}
	if len(results) == 0 {
		printYellow(voice.Line("no_results"))
		return waitForEnter(reader)
	}
	printWhite("\n  FILE                             LINE  CONTENT")
	for _, item := range results {
		printBlue(fmt.Sprintf("  %-32s", item.Path))
		printWhite(fmt.Sprintf("  %-4d %s", item.Line, item.Content))
	}
	printMagenta(voice.Line("finder_content_success"))
	return waitForEnter(reader)
}

func findBySize(reader *bufio.Reader) error {
	path, err := promptPath(reader)
	if err != nil {
		return err
	}
	fmt.Println(voice.Line("finder_size_condition"))
	fmt.Printf("  %s", voice.Line("finder_condition_input"))
	cond, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	cond = strings.TrimSpace(cond)
	fmt.Printf("  %s", voice.Line("finder_size_prompt"))
	sizeInput, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	printYellow(voice.Line("searching"))
	results, err := finder.FindBySize(path, cond, strings.TrimSpace(sizeInput))
	if err != nil {
		return err
	}
	return showFinderResults(reader, results)
}

func findByTime(reader *bufio.Reader) error {
	path, err := promptPath(reader)
	if err != nil {
		return err
	}
	fmt.Println(voice.Line("finder_time_condition"))
	fmt.Printf("  %s", voice.Line("finder_condition_input"))
	cond, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Printf("  %s", voice.Line("finder_days_prompt"))
	daysInput, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	printYellow(voice.Line("searching"))
	results, err := finder.FindByTime(path, strings.TrimSpace(cond), strings.TrimSpace(daysInput))
	if err != nil {
		return err
	}
	return showFinderResults(reader, results)
}

func findGlobal(reader *bufio.Reader) error {
	fmt.Printf("  %s", voice.Line("finder_global_prompt"))
	query, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	query = strings.TrimSpace(query)
	if query == "" {
		printRed(voice.Line("invalid_keyword"))
		return waitForEnter(reader)
	}
	fmt.Printf("  %s", voice.Line("finder_match_prompt"))
	matchType, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	exact := strings.TrimSpace(matchType) == "2"
	fmt.Printf("  %s", voice.Line("finder_global_type_prompt"))
	extInput, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	extInput = strings.TrimSpace(extInput)
	if extInput != "" && !strings.HasPrefix(extInput, ".") {
		extInput = "." + extInput
	}
	printYellow(voice.Line("searching_once"))
	searchPath, err := finder.GlobalSearchPath()
	if err != nil {
		return err
	}
	printYellow(fmt.Sprintf("%s%s", voice.Line("searching_dir"), searchPath))

	results, err := finder.FindByNameWithExt(searchPath, query, extInput, exact)
	if err != nil {
		return err
	}
	if len(results) == 0 {
		printYellow(voice.Line("global_no_results"))
		return waitForEnter(reader)
	}

	printMagenta(fmt.Sprintf("\n  %s", voice.Linef("finder_results_count", len(results))))
	printWhite("\n  PATH                             SIZE    MODIFIED")
	for _, item := range results {
		printFinderPath(item)
		printWhite(fmt.Sprintf("  %-7s %s", item.Size, item.Modified))
	}
	printMagenta(fmt.Sprintf("%s%s", voice.Line("global_found"), searchPath))
	return waitForEnter(reader)
}

func promptPath(reader *bufio.Reader) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}
	fmt.Printf("  %s", voice.Line("finder_path_prompt"))
	path, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	path = strings.TrimSpace(path)
	if path == "" {
		return cfg.DefaultSearchPath, nil
	}
	return path, nil
}

func showFinderResults(reader *bufio.Reader, results []finder.Result) error {
	if len(results) == 0 {
		printYellow(voice.Line("no_results"))
		return waitForEnter(reader)
	}
	printMagenta(fmt.Sprintf("\n  %s", voice.Linef("finder_results_count", len(results))))
	printWhite("\n  PATH                             SIZE    MODIFIED")
	for _, item := range results {
		printFinderPath(item)
		printWhite(fmt.Sprintf("  %-7s %s", item.Size, item.Modified))
	}
	printMagenta(voice.Line("finder_results_done"))
	return waitForEnter(reader)
}

func printFinderPath(item finder.Result) {
	path := fmt.Sprintf("  %-32s", item.Path)
	if item.IsDir {
		color.New(color.FgYellow).Print(path)
		return
	}
	printBlue(path)
}
