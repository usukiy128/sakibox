package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"sakibox/config"
	"sakibox/internal/voice"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sakibox",
	Short: "sakibox - terminal toolbox",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.EnsureConfig(); err != nil {
			return err
		}
		return showMainMenu()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.New(color.FgRed).Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func showMainMenu() error {
	reader := bufio.NewReader(os.Stdin)

	for {
		printLogo()
		printMagenta(voice.Line("greeting"))
		fmt.Println()
		fmt.Println("  1. 端口管理")
		fmt.Println("  2. 历史命令")
		fmt.Println("  3. 命令收藏夹")
		fmt.Println("  4. 进程监控")
		fmt.Println("  5. 文件查找")
		fmt.Println("  6. 安装帮助")
		fmt.Println("  7. SSH 工具")
		fmt.Println("  8. 更新 sakibox")
		fmt.Println("  0. 退出")
		fmt.Printf("\n  %s", voice.Line("main_prompt"))

		choice, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			if err := showPortMenu(reader); err != nil {
				return err
			}
		case "2":
			if err := showHistoryMenu(reader); err != nil {
				return err
			}
		case "3":
			if err := showBookmarkMenu(reader); err != nil {
				return err
			}
		case "4":
			if err := showProcessMenu(reader); err != nil {
				return err
			}
		case "5":
			if err := showFinderMenu(reader); err != nil {
				return err
			}
		case "6":
			if err := showInstallMenu(reader); err != nil {
				return err
			}
		case "7":
			if err := showSSHMenu(reader); err != nil {
				return err
			}
		case "8":
			printMagenta(voice.Line("update_intro"))
			if err := runUpdate(); err != nil {
				printRed(voice.Line("update_failed"))
			} else {
				printGreen(voice.Line("update_done"))
			}
			if err := waitForEnter(reader); err != nil {
				return err
			}
		case "0":
			printMagenta(voice.Line("exit"))
			return nil
		default:
			printRed(voice.Line("invalid_option"))
			if err := waitForEnter(reader); err != nil {
				return err
			}
		}
	}
}

func printLogo() {
	logo := "_    _ _               \n          | |  (_) |              \n ___  __ _| | ___| |__   _____  __\n/ __|/ _` | |/ / | '_ \\ / _ \\ \\/ /\n\\__ \\ (_| |   <| | |_) | (_) >  < \n|___/\\__,_|_|\\_\\_|_.__/ \\___/_/\\_\\ v0.0.5"
	printCyan(logo)
	fmt.Println()
}

func waitForEnter(reader *bufio.Reader) error {
	printWhite(voice.Line("press_enter"))
	_, err := reader.ReadString('\n')
	return err
}

func printMagenta(msg string) {
	color.New(color.FgMagenta).Println(msg)
}

func printGreen(msg string) {
	color.New(color.FgGreen).Println(msg)
}

func printRed(msg string) {
	color.New(color.FgRed).Println(msg)
}

func printYellow(msg string) {
	color.New(color.FgYellow).Println(msg)
}

func printCyan(msg string) {
	color.New(color.FgCyan).Println(msg)
}

func printWhite(msg string) {
	color.New(color.FgWhite).Println(msg)
}

func printBlue(msg string) {
	color.New(color.FgBlue).Println(msg)
}

func runUpdate() error {
	cmd := exec.Command("/bin/sh", "-c", "curl -fsSL https://raw.githubusercontent.com/usukiy128/sakibox/main/install.sh | sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
