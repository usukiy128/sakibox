package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"sakibox/internal/ssh"
	"sakibox/internal/voice"
)

func showSSHMenu(reader *bufio.Reader) error {
	for {
		printCyan("[SSH 工具]")
		printMagenta(voice.Line("ssh_intro"))
		fmt.Println("  1. 查看服务器")
		fmt.Println("  2. 添加服务器")
		fmt.Println("  3. 快速连接")
		fmt.Println("  4. 删除服务器")
		fmt.Println("  5. 快速命令")
		fmt.Println("  6. 连接日志")
		fmt.Println("  0. 返回主菜单")
		fmt.Printf("\n  %s", voice.Line("menu_prompt"))

		choice, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			if err := showSSHServers(reader); err != nil {
				return err
			}
		case "2":
			if err := addSSHServer(reader); err != nil {
				return err
			}
		case "3":
			if err := connectSSHServer(reader); err != nil {
				return err
			}
		case "4":
			if err := deleteSSHServer(reader); err != nil {
				return err
			}
		case "5":
			if err := showSSHCommandsMenu(reader); err != nil {
				return err
			}
		case "6":
			if err := showSSHLogs(reader); err != nil {
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

func showSSHServers(reader *bufio.Reader) error {
	servers, err := ssh.List()
	if err != nil {
		return err
	}
	if len(servers) == 0 {
		printYellow(voice.Line("ssh_empty"))
		return waitForEnter(reader)
	}
	printWhite("\n  #  NAME         HOST               USER      PORT  PASS")
	for i, item := range servers {
		pass := "no"
		if strings.TrimSpace(item.Password) != "" {
			pass = "yes"
		}
		fmt.Printf("  %-3d %-12s %-18s %-8s %-4d %s\n", i+1, item.Name, item.Host, item.User, item.Port, pass)
	}
	printMagenta(voice.Line("ssh_list_done"))
	return waitForEnter(reader)
}

func addSSHServer(reader *bufio.Reader) error {
	fmt.Printf("\n  %s", voice.Line("ssh_add_name_prompt"))
	name, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Printf("  %s", voice.Line("ssh_add_host_prompt"))
	host, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Printf("  %s", voice.Line("ssh_add_user_prompt"))
	user, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Printf("  %s", voice.Line("ssh_add_port_prompt"))
	portInput, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	port, err := strconv.Atoi(strings.TrimSpace(portInput))
	if err != nil {
		printRed(voice.Line("ssh_invalid_port"))
		return waitForEnter(reader)
	}
	fmt.Printf("  %s", voice.Line("ssh_add_password_prompt"))
	password, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	entry := ssh.Server{
		Name:     strings.TrimSpace(name),
		Host:     strings.TrimSpace(host),
		User:     strings.TrimSpace(user),
		Port:     port,
		Password: strings.TrimSpace(password),
	}
	if err := ssh.Add(entry); err != nil {
		printRed(err.Error())
	} else {
		printGreen(voice.Line("ssh_add_success"))
	}
	return waitForEnter(reader)
}

func connectSSHServer(reader *bufio.Reader) error {
	fmt.Printf("\n  %s", voice.Line("ssh_connect_prompt"))
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	server, err := ssh.Get(strings.TrimSpace(input))
	if err != nil {
		printRed(err.Error())
		return waitForEnter(reader)
	}
	printYellow(voice.Line("ssh_connecting"))
	if strings.TrimSpace(server.Password) == "" {
		printYellow(voice.Line("ssh_pass_missing"))
	} else if !ensureSSHPass(reader) {
		_ = ssh.AddLog(ssh.NewLog(server, "connect", fmt.Errorf("sshpass not installed")))
		return nil
	}
	command := buildSSHCommand(server, "")
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	_ = ssh.AddLog(ssh.NewLog(server, "connect", err))
	if err != nil {
		printRed(err.Error())
	} else {
		printMagenta(voice.Line("ssh_connect_success"))
	}
	return waitForEnter(reader)
}

func deleteSSHServer(reader *bufio.Reader) error {
	fmt.Printf("\n  %s", voice.Line("ssh_delete_prompt"))
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		printRed(voice.Line("invalid_index"))
		return waitForEnter(reader)
	}
	fmt.Printf("  %s", voice.Line("ssh_delete_confirm"))
	confirm, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
		printYellow(voice.Line("ssh_delete_cancel"))
		return waitForEnter(reader)
	}
	if err := ssh.Delete(index); err != nil {
		printRed(err.Error())
	} else {
		printGreen(voice.Line("ssh_delete_success"))
	}
	return waitForEnter(reader)
}

func showSSHCommandsMenu(reader *bufio.Reader) error {
	for {
		printCyan("[快速命令]")
		printMagenta(voice.Line("ssh_cmd_intro"))
		fmt.Println("  1. 查看命令")
		fmt.Println("  2. 添加命令")
		fmt.Println("  3. 执行命令")
		fmt.Println("  4. 删除命令")
		fmt.Println("  0. 返回上级")
		fmt.Printf("\n  %s", voice.Line("menu_prompt"))

		choice, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			if err := showSSHCommands(reader); err != nil {
				return err
			}
		case "2":
			if err := addSSHCommand(reader); err != nil {
				return err
			}
		case "3":
			if err := executeSSHCommand(reader); err != nil {
				return err
			}
		case "4":
			if err := deleteSSHCommand(reader); err != nil {
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

func showSSHCommands(reader *bufio.Reader) error {
	items, err := ssh.ListCommands()
	if err != nil {
		return err
	}
	if len(items) == 0 {
		printYellow(voice.Line("ssh_cmd_empty"))
		return waitForEnter(reader)
	}
	printWhite("\n  #  NAME         COMMAND")
	for i, item := range items {
		fmt.Printf("  %-3d %-12s %s\n", i+1, item.Name, item.Command)
	}
	printMagenta(voice.Line("ssh_cmd_list_done"))
	return waitForEnter(reader)
}

func addSSHCommand(reader *bufio.Reader) error {
	fmt.Printf("\n  %s", voice.Line("ssh_cmd_add_name_prompt"))
	name, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Printf("  %s", voice.Line("ssh_cmd_add_cmd_prompt"))
	cmdLine, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	item := ssh.Command{
		Name:    strings.TrimSpace(name),
		Command: strings.TrimSpace(cmdLine),
	}
	if err := ssh.AddCommand(item); err != nil {
		printRed(err.Error())
	} else {
		printGreen(voice.Line("ssh_cmd_add_success"))
	}
	return waitForEnter(reader)
}

func executeSSHCommand(reader *bufio.Reader) error {
	fmt.Printf("\n  %s", voice.Line("ssh_cmd_exec_prompt"))
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	item, err := ssh.GetCommand(strings.TrimSpace(input))
	if err != nil {
		printRed(err.Error())
		return waitForEnter(reader)
	}
	fmt.Printf("  %s", voice.Line("ssh_cmd_server_prompt"))
	serverInput, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	server, err := ssh.Get(strings.TrimSpace(serverInput))
	if err != nil {
		printRed(err.Error())
		return waitForEnter(reader)
	}
	if strings.TrimSpace(server.Password) != "" && !ensureSSHPass(reader) {
		_ = ssh.AddLog(ssh.NewLog(server, fmt.Sprintf("cmd:%s", item.Name), fmt.Errorf("sshpass not installed")))
		return nil
	}
	command := buildSSHCommand(server, item.Command)
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	_ = ssh.AddLog(ssh.NewLog(server, fmt.Sprintf("cmd:%s", item.Name), err))
	if err != nil {
		printRed(err.Error())
	} else {
		printMagenta(voice.Line("ssh_cmd_exec_success"))
	}
	return waitForEnter(reader)
}

func deleteSSHCommand(reader *bufio.Reader) error {
	fmt.Printf("\n  %s", voice.Line("ssh_cmd_delete_prompt"))
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		printRed(voice.Line("invalid_index"))
		return waitForEnter(reader)
	}
	fmt.Printf("  %s", voice.Line("ssh_cmd_delete_confirm"))
	confirm, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
		printYellow(voice.Line("ssh_cmd_delete_cancel"))
		return waitForEnter(reader)
	}
	if err := ssh.DeleteCommand(index); err != nil {
		printRed(err.Error())
	} else {
		printGreen(voice.Line("ssh_cmd_delete_success"))
	}
	return waitForEnter(reader)
}

func buildSSHCommand(server ssh.Server, remoteCmd string) string {
	base := fmt.Sprintf("ssh %s@%s -p %d", server.User, server.Host, server.Port)
	if strings.TrimSpace(server.Password) == "" {
		if strings.TrimSpace(remoteCmd) == "" {
			return base
		}
		return fmt.Sprintf("%s %s", base, remoteCmd)
	}
	command := fmt.Sprintf("sshpass -p %q %s", server.Password, base)
	if strings.TrimSpace(remoteCmd) == "" {
		return command
	}
	return fmt.Sprintf("%s %s", command, remoteCmd)
}

func ensureSSHPass(reader *bufio.Reader) bool {
	if _, err := exec.LookPath("sshpass"); err == nil {
		return true
	}
	printYellow(voice.Line("sshpass_missing"))
	printYellow(voice.Line("sshpass_install_hint"))
	_ = waitForEnter(reader)
	return false
}

func showSSHLogs(reader *bufio.Reader) error {
	logs, err := ssh.ListLogs()
	if err != nil {
		return err
	}
	if len(logs) == 0 {
		printYellow(voice.Line("ssh_log_empty"))
		return waitForEnter(reader)
	}
	printWhite("\n  TIME                NAME         HOST               ACTION        RESULT")
	for _, entry := range logs {
		result := "OK"
		if !entry.Success {
			result = "FAIL"
		}
		fmt.Printf("  %-19s %-12s %-18s %-12s %s\n", entry.Time, entry.Name, entry.Host, entry.Action, result)
	}
	printMagenta(voice.Line("ssh_log_done"))
	return waitForEnter(reader)
}
