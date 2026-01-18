package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"sakibox/internal/process"
	"sakibox/internal/voice"

	"golang.org/x/term"
)

func showProcessMenu(reader *bufio.Reader) error {
	for {
		printCyan("[进程监控]")
		printMagenta(voice.Line("process_intro"))
		fmt.Println("  1. 查看所有进程(实时)")
		fmt.Println("  2. 查看资源占用TOP10(实时)")
		fmt.Println("  3. 搜索进程")
		fmt.Println("  4. 杀死进程")
		fmt.Println("  0. 返回主菜单")
		fmt.Printf("\n  %s", voice.Line("menu_prompt"))

		choice, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			if err := showAllProcesses(reader); err != nil {
				return err
			}
		case "2":
			if err := showTopProcessesLive(reader); err != nil {
				return err
			}
		case "3":
			fmt.Printf("\n  %s", voice.Line("process_search_prompt"))
			name, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			entries, err := process.Search(strings.TrimSpace(name))
			if err != nil {
				return err
			}
			if len(entries) == 0 {
				printYellow(voice.Line("process_search_empty"))
				if err := waitForEnter(reader); err != nil {
					return err
				}
				continue
			}
			printProcessList(entries)
			printMagenta(voice.Line("process_search_success"))
			if err := waitForEnter(reader); err != nil {
				return err
			}
		case "4":
			fmt.Printf("\n  %s", voice.Line("process_kill_prompt"))
			input, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			pid, err := strconv.Atoi(strings.TrimSpace(input))
			if err != nil {
				printRed(voice.Line("process_invalid_pid"))
				if err := waitForEnter(reader); err != nil {
					return err
				}
				continue
			}
			fmt.Printf("  %s", voice.Line("process_kill_confirm"))
			confirm, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
				printYellow(voice.Line("process_kill_cancel"))
				if err := waitForEnter(reader); err != nil {
					return err
				}
				continue
			}
			if err := process.Kill(pid); err != nil {
				printRed(err.Error())
			} else {
				printGreen(voice.Line("process_kill_success"))
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

func showAllProcesses(reader *bufio.Reader) error {
	printMagenta(voice.Line("process_live_hint"))
	if err := enableRawMode(); err != nil {
		return err
	}
	defer disableRawMode()

	for {
		entries, err := process.List(0)
		if err != nil {
			return err
		}

		cols, rows, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil || cols <= 0 || rows <= 0 {
			cols, rows = 80, 24
		}

		frame := buildProcessFrame(entries, rows, cols, voice.Line("process_live_prompt"))
		_, _ = os.Stdout.WriteString(frame)

		if waitForQuit(time.Second) {
			disableRawMode()
			printMagenta(voice.Line("process_list_done"))
			return waitForEnter(reader)
		}
	}
}

func showTopProcessesLive(reader *bufio.Reader) error {
	for {
		entries, err := process.Top()
		if err != nil {
			return err
		}
		printWhite("\n  PID   NAME          CPU%   MEM%")
		for _, entry := range entries {
			fmt.Printf("  %-5d %-12s %-6.1f %-6.1f\n", entry.PID, entry.Name, entry.CPU, entry.Mem)
		}
		printMagenta(voice.Line("process_top_hint"))
		fmt.Printf("\n  %s", voice.Line("process_top_prompt"))
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		switch strings.TrimSpace(input) {
		case "0":
			return nil
		case "r", "R":
			continue
		default:
			printRed(voice.Line("invalid_option"))
		}
	}
}

func printProcessList(entries []process.Entry) {
	printWhite("\n  PID   NAME")
	for _, entry := range entries {
		fmt.Printf("  %-5d %s\n", entry.PID, entry.Name)
	}
}

func buildProcessFrame(entries []process.Entry, rows, cols int, footer string) string {
	maxRows := rows - 4
	if maxRows < 5 {
		maxRows = 5
	}
	if len(entries) > maxRows {
		entries = entries[:maxRows]
	}

	var frame strings.Builder
	frame.WriteString("\033[2J\033[H")
	frame.WriteString("\r" + fmt.Sprintf("%-5s %s", "PID", "NAME") + "\n")

	for i := 0; i < maxRows; i++ {
		if i < len(entries) {
			line := "\r" + fmt.Sprintf("%-5d %s", entries[i].PID, entries[i].Name)
			frame.WriteString(padLine(line, cols))
		} else {
			frame.WriteString(padLine("\r", cols))
		}
	}

	frame.WriteString(padLine("\r", cols))
	frame.WriteString(padLine("\r"+footer, cols))
	return frame.String()
}

func padLine(line string, cols int) string {
	if cols <= 0 {
		return line + "\n"
	}
	normalized := strings.TrimPrefix(line, "\r")
	if len(normalized) > cols {
		normalized = normalized[:cols]
	}
	if len(normalized) < cols {
		normalized = normalized + strings.Repeat(" ", cols-len(normalized))
	}
	return "\r" + normalized + "\n"
}

func waitForQuit(timeout time.Duration) bool {
	buf := make([]byte, 1)
	expire := time.Now().Add(timeout)
	for time.Now().Before(expire) {
		_ = os.Stdin.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
		n, err := os.Stdin.Read(buf)
		if errors.Is(err, syscall.EAGAIN) || errors.Is(err, syscall.EWOULDBLOCK) {
			continue
		}
		if errors.Is(err, os.ErrDeadlineExceeded) || errors.Is(err, syscall.ETIMEDOUT) {
			continue
		}
		if err != nil || n == 0 {
			if err != nil {
				return false
			}
			continue
		}
		if buf[0] == 'q' || buf[0] == 'Q' {
			return true
		}
	}
	return false
}

var originalState *term.State

func enableRawMode() error {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nil
	}
	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	originalState = state
	if err := syscall.SetNonblock(int(os.Stdin.Fd()), true); err != nil {
		return err
	}
	return nil
}

func disableRawMode() {
	if originalState == nil {
		return
	}
	_ = syscall.SetNonblock(int(os.Stdin.Fd()), false)
	_ = term.Restore(int(os.Stdin.Fd()), originalState)
	originalState = nil
}

func clearScreen() {
	fmt.Fprint(os.Stdout, "\033[H\033[2J")
}
