package cmd

import (
	"bufio"
	"fmt"
	"strings"

	"sakibox/internal/install"
	"sakibox/internal/voice"
)

func showInstallMenu(reader *bufio.Reader) error {
	for {
		printCyan("[安装帮助]")
		printMagenta(voice.Line("install_intro"))
		fmt.Println("  1. 一键安装Linux主流工具/依赖")
		fmt.Println("  2. 通过查找依赖/命令名进行下载")
		fmt.Println("  0. 返回主菜单")
		fmt.Printf("\n  %s", voice.Line("menu_prompt"))

		choice, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			printYellow(voice.Line("install_prepare"))
			cmdLine, err := install.InstallDefaults()
			if err != nil {
				printRed(err.Error())
				if err := waitForEnter(reader); err != nil {
					return err
				}
				continue
			}
			printYellow(voice.Linef("install_preview", cmdLine))
			printGreen(voice.Line("install_defaults_success"))
			if err := waitForEnter(reader); err != nil {
				return err
			}
		case "2":
			fmt.Printf("\n  %s", voice.Line("install_search_prompt"))
			keyword, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			printYellow(voice.Line("searching"))
			cmdLine, err := install.SearchAndInstall(strings.TrimSpace(keyword))
			if err != nil {
				printRed(err.Error())
				if err := waitForEnter(reader); err != nil {
					return err
				}
				continue
			}
			printYellow(voice.Linef("install_preview", cmdLine))
			printGreen(voice.Line("install_search_success"))
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
