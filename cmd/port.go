package cmd

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"sakibox/internal/port"
	"sakibox/internal/voice"
)

func showPortMenu(reader *bufio.Reader) error {
	for {
		printCyan("[端口管理]")
		printMagenta(voice.Line("port_intro"))
		fmt.Println("  1. 查看所有端口")
		fmt.Println("  2. 查找指定端口")
		fmt.Println("  3. 关闭端口进程")
		fmt.Println("  0. 返回主菜单")
		fmt.Printf("\n  %s", voice.Line("menu_prompt"))

		choice, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			if err := showAllPorts(); err != nil {
				return err
			}
			printMagenta(voice.Line("port_list_done"))
			if err := waitForEnter(reader); err != nil {
				return err
			}
		case "2":
			fmt.Printf("\n  %s", voice.Line("port_find_prompt"))
			input, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			portNum, err := strconv.Atoi(strings.TrimSpace(input))
			if err != nil {
				printRed(voice.Line("port_invalid_number"))
				if err := waitForEnter(reader); err != nil {
					return err
				}
				continue
			}
			found, err := showPortByNumber(portNum)
			if err != nil {
				return err
			}
			if found {
				printMagenta(voice.Line("port_find_success"))
			}
			if err := waitForEnter(reader); err != nil {
				return err
			}
		case "3":
			fmt.Printf("\n  %s", voice.Line("port_kill_prompt"))
			input, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			portNum, err := strconv.Atoi(strings.TrimSpace(input))
			if err != nil {
				printRed(voice.Line("port_invalid_number"))
				if err := waitForEnter(reader); err != nil {
					return err
				}
				continue
			}
			fmt.Printf("  %s", voice.Line("port_kill_confirm"))
			confirm, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			confirm = strings.TrimSpace(confirm)
			if strings.ToLower(confirm) != "y" {
				printYellow(voice.Line("port_kill_cancel"))
				if err := waitForEnter(reader); err != nil {
					return err
				}
				continue
			}
			if err := port.KillByPort(portNum); err != nil {
				printRed(err.Error())
			} else {
				printGreen(voice.Line("port_kill_success"))
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

func showAllPorts() error {
	entries, err := port.ListPorts()
	if err != nil {
		return err
	}
	printWhite("\n  PORT  PROCESS       PID")
	for _, entry := range entries {
		fmt.Printf("  %-5d %-12s %d\n", entry.Port, entry.Process, entry.PID)
	}
	return nil
}

func showPortByNumber(portNum int) (bool, error) {
	entry, err := port.FindPort(portNum)
	if err != nil {
		printRed(err.Error())
		return false, nil
	}
	printWhite("\n  PORT  PROCESS       PID")
	fmt.Printf("  %-5d %-12s %d\n", entry.Port, entry.Process, entry.PID)
	return true, nil
}
