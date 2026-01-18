package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"sakibox/internal/history"
	"sakibox/internal/voice"
)

func showHistoryMenu(reader *bufio.Reader) error {
	for {
		printCyan("[历史命令]")
		printMagenta(voice.Line("history_intro"))
		fmt.Println("  1. 查看历史命令")
		fmt.Println("  2. 搜索历史命令")
		fmt.Println("  3. 执行历史命令")
		fmt.Println("  0. 返回主菜单")
		fmt.Printf("\n  %s", voice.Line("menu_prompt"))

		choice, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			if err := showHistoryList(); err != nil {
				return err
			}
			printMagenta(voice.Line("history_list_done"))
			if err := waitForEnter(reader); err != nil {
				return err
			}
		case "2":
			fmt.Printf("\n  %s", voice.Line("history_search_prompt"))
			keyword, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			keyword = strings.TrimSpace(keyword)
			matches, err := history.Search(keyword)
			if err != nil {
				return err
			}
			if len(matches) == 0 {
				printYellow(voice.Line("history_search_empty"))
				if err := waitForEnter(reader); err != nil {
					return err
				}
				continue
			}
			printHistoryList(matches)
			printMagenta(voice.Line("history_search_success"))
			if err := waitForEnter(reader); err != nil {
				return err
			}
		case "3":
			fmt.Printf("\n  %s", voice.Line("history_exec_prompt"))
			input, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			index, err := strconv.Atoi(strings.TrimSpace(input))
			if err != nil {
				printRed(voice.Line("invalid_index"))
				if err := waitForEnter(reader); err != nil {
					return err
				}
				continue
			}
			cmdLine, err := history.GetByIndex(index)
			if err != nil {
				printRed(err.Error())
				if err := waitForEnter(reader); err != nil {
					return err
				}
				continue
			}
			if err := executeShellCommand(cmdLine); err != nil {
				printRed(err.Error())
			} else {
				printMagenta(voice.Line("history_exec_success"))
			}
			if err := waitForEnter(reader); err != nil {
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

func showHistoryList() error {
	entries, err := history.List()
	if err != nil {
		return err
	}
	printHistoryList(entries)
	return nil
}

func printHistoryList(entries []history.Entry) {
	printWhite("\n  #  COMMAND")
	for i, entry := range entries {
		fmt.Printf("  %-3d %s\n", i+1, entry.Command)
	}
}

func executeShellCommand(command string) error {
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
