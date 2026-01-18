package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"sakibox/internal/bookmark"
	"sakibox/internal/voice"
)

func showBookmarkMenu(reader *bufio.Reader) error {
	for {
		printCyan("[命令收藏夹]")
		printMagenta(voice.Line("bookmark_intro"))
		fmt.Println("  1. 查看收藏")
		fmt.Println("  2. 添加收藏")
		fmt.Println("  3. 执行收藏")
		fmt.Println("  4. 删除收藏")
		fmt.Println("  0. 返回主菜单")
		fmt.Printf("\n  %s", voice.Line("menu_prompt"))

		choice, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			items, err := bookmark.List()
			if err != nil {
				return err
			}
			printWhite("\n  #  NAME         COMMAND")
			for i, item := range items {
				fmt.Printf("  %-3d %-12s %s\n", i+1, item.Name, item.Command)
			}
			printMagenta(voice.Line("bookmark_list_done"))
			if err := waitForEnter(reader); err != nil {
				return err
			}
		case "2":
			fmt.Printf("\n  %s", voice.Line("bookmark_add_name_prompt"))
			name, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			fmt.Printf("  %s", voice.Line("bookmark_add_cmd_prompt"))
			cmdLine, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			if err := bookmark.Add(strings.TrimSpace(name), strings.TrimSpace(cmdLine)); err != nil {
				printRed(err.Error())
			} else {
				printGreen(voice.Line("bookmark_add_success"))
			}
			if err := waitForEnter(reader); err != nil {
				return err
			}
		case "3":
			fmt.Printf("\n  %s", voice.Line("bookmark_exec_prompt"))
			input, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			cmdLine, err := bookmark.Get(strings.TrimSpace(input))
			if err != nil {
				printRed(err.Error())
				if err := waitForEnter(reader); err != nil {
					return err
				}
				continue
			}
			cmd := exec.Command("/bin/sh", "-c", cmdLine)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				printRed(err.Error())
			} else {
				printMagenta(voice.Line("bookmark_exec_success"))
			}
			if err := waitForEnter(reader); err != nil {
				return err
			}
		case "4":
			fmt.Printf("\n  %s", voice.Line("bookmark_delete_prompt"))
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
			fmt.Printf("  %s", voice.Line("bookmark_delete_confirm"))
			confirm, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
				printYellow(voice.Line("bookmark_delete_cancel"))
				if err := waitForEnter(reader); err != nil {
					return err
				}
				continue
			}
			if err := bookmark.Delete(index); err != nil {
				printRed(err.Error())
			} else {
				printGreen(voice.Line("bookmark_delete_success"))
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
