package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/itshivams/studex-cli/internal/api"
	"github.com/itshivams/studex-cli/internal/config"
)


var (
	clrTitle     = color.New(color.FgHiCyan, color.Bold)
	clrBorder    = color.New(color.FgHiBlue)
	clrSelf      = color.New(color.FgHiGreen, color.Bold)
	clrOther     = color.New(color.FgHiWhite, color.Bold)
	clrAI        = color.New(color.FgHiMagenta, color.Bold)
	clrAnon      = color.New(color.FgYellow, color.Bold)
	clrSystem    = color.New(color.FgHiBlack, color.Italic)
	clrTime      = color.New(color.FgHiBlack)
	clrPrompt    = color.New(color.FgCyan, color.Bold)
	clrError     = color.New(color.FgRed, color.Bold)
	clrSuccess   = color.New(color.FgGreen, color.Bold)
	clrHighlight = color.New(color.FgHiYellow, color.Bold)
	clrVerified  = color.New(color.FgHiCyan)
	clrInfo      = color.New(color.FgHiBlue)
)


func termWidth() int {
	return 80
}

func centerText(s string, width int) string {
	visLen := utf8.RuneCountInString(s)
	if visLen >= width {
		return s
	}
	pad := (width - visLen) / 2
	return strings.Repeat(" ", pad) + s
}

func repeatStr(s string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(s, n)
}

func boxLine(content string, width int) string {
	visLen := utf8.RuneCountInString(content)
	pad := width - visLen - 2 
	if pad < 0 {
		pad = 0
	}
	return clrBorder.Sprint("│") + content + strings.Repeat(" ", pad) + clrBorder.Sprint("│")
}

func nowHHMM() string {
	return time.Now().Format("03:04 PM")
}


func printDiscussionBanner(online int, username string, isAnon bool) {
	w := 68
	top := "╭" + repeatStr("─", w) + "╮"
	bot := "╰" + repeatStr("─", w) + "╯"

	clrBorder.Println(top)

	title := "  🚀  S T U D E X   L I V E   D I S C U S S I O N  🚀"
	clrBorder.Print("│")
	clrTitle.Printf("%-*s", w, title)
	clrBorder.Println("│")

	subtitle := fmt.Sprintf("  Connected • %d online", online)
	if online == 0 {
		subtitle = "  Connecting…"
	}
	clrBorder.Print("│")
	clrInfo.Printf("%-*s", w, subtitle)
	clrBorder.Println("│")

	userLine := "  Chatting as: "
	if isAnon {
		userLine += "👻 anonymous"
	} else {
		userLine += "👤 " + username
	}
	clrBorder.Print("│")
	clrHighlight.Printf("%-*s", w, userLine)
	clrBorder.Println("│")

	hint := "  Type /help for commands • /anon to toggle anonymous • /quit to exit"
	clrBorder.Print("│")
	clrSystem.Printf("%-*s", w, hint)
	clrBorder.Println("│")

	clrBorder.Println(bot)
}

func printHelp() {
	w := 60
	clrBorder.Println("╭" + repeatStr("─", w) + "╮")

	rows := []struct{ cmd, desc string }{
		{"/help", "Show this help"},
		{"/anon", "Toggle anonymous mode on/off"},
		{"/clear", "Clear chat history from screen"},
		{"/status", "Show connection status & online count"},
		{"@studex <q>", "Ask Studex AI a question"},
		{"/quit or /exit", "Leave the discussion room"},
	}
	for _, r := range rows {
		line := fmt.Sprintf("  %-20s  %s", r.cmd, r.desc)
		clrBorder.Print("│")
		fmt.Printf("%-*s", w, line)
		clrBorder.Println("│")
	}
	clrBorder.Println("╰" + repeatStr("─", w) + "╯")
}

func renderMessage(msg *api.ChatMessage, myUsername string, isAnon bool) {
	if msg.Type == "stats" {
		return
	}

	sender := strings.TrimSpace(msg.Username)
	content := strings.TrimSpace(msg.Content)
	ts := msg.Timestamp
	if ts == "" {
		ts = nowHHMM()
	}

	isSelf := strings.EqualFold(sender, myUsername) && !isAnon
	isAI := msg.Type == "ai" || strings.EqualFold(sender, "StudexAI")
	isAnonSender := strings.EqualFold(sender, "anonymous")

	var nameStr string
	switch {
	case isAI:
		nameStr = clrAI.Sprint("🤖 StudexAI [AI]")
	case isSelf:
		nameStr = clrSelf.Sprint("▶ You")
	case isAnonSender:
		nameStr = clrAnon.Sprint("👻 anonymous")
	default:
		verified := ""
		if msg.Verified {
			verified = clrVerified.Sprint(" ✓")
		}
		nameStr = clrOther.Sprint("● "+sender) + verified
	}

	timeStr := clrTime.Sprint("[" + ts + "]")

	if isSelf {
		clrBorder.Println("  " + repeatStr("·", 62))
		fmt.Printf("  %s %s\n", timeStr, nameStr)
		for _, line := range wrapText(content, 60) {
			clrSelf.Printf("    %s\n", line)
		}
	} else {
		fmt.Printf("  %s %s\n", nameStr, timeStr)
		for _, line := range wrapText(content, 60) {
			if isAI {
				clrAI.Printf("    %s\n", line)
			} else {
				fmt.Printf("    %s\n", line)
			}
		}
	}
}

func wrapText(s string, maxWidth int) []string {
	if maxWidth <= 0 {
		maxWidth = 60
	}
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{s}
	}
	var lines []string
	current := ""
	for _, w := range words {
		if current == "" {
			current = w
		} else if utf8.RuneCountInString(current)+1+utf8.RuneCountInString(w) <= maxWidth {
			current += " " + w
		} else {
			lines = append(lines, current)
			current = w
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}


func LiveDiscussionMenu() {
	myUsername := config.GetStoredUsername()
	if myUsername == "" {
		myUsername = "anonymous"
	}
	isAnon := false
	onlineCount := 0

	fmt.Print("\033[H\033[2J")

	clrInfo.Println("\n  Connecting to Studex Discussion…")

	conn, err := api.ConnectDiscussion()
	if err != nil {
		clrError.Printf("\n  ✗ Could not connect: %v\n", err)
		clrSystem.Println("  Press Enter to go back…")
		fmt.Scanln()
		return
	}
	defer conn.Close()

	clrSuccess.Println("  ✓ Connected!\n")

	printDiscussionBanner(onlineCount, myUsername, isAnon)
	clrSystem.Println("\n  ── Chat history will appear below ──────────────────────────\n")

	var mu sync.Mutex 
	incoming := make(chan *api.ChatMessage, 64)
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			msg, err := conn.Receive()
			if err != nil {
				if err != io.EOF {
			
					select {
					case incoming <- &api.ChatMessage{Type: "_err", Content: err.Error()}:
					default:
					}
				}
				return
			}
			if msg.Type == "stats" {
				mu.Lock()
				onlineCount = msg.ActiveUsers
				mu.Unlock()
				continue
			}
			incoming <- msg
		}
	}()

	go func() {
		for msg := range incoming {
			if msg.Type == "_err" {
				clrError.Printf("\n  ⚠ Connection error: %s\n", msg.Content)
				clrSystem.Print("  » ")
				continue
			}

			renderMessage(msg, myUsername, isAnon)
			clrPrompt.Print("  » ")
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	clrPrompt.Print("  » ")

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			clrPrompt.Print("  » ")
			continue
		}

		if strings.HasPrefix(line, "/") {
			cmd := strings.ToLower(strings.Fields(line)[0])
			switch cmd {
			case "/quit", "/exit":
				clrSuccess.Println("\n  👋  Left the discussion. See you next time!\n")
				return

			case "/help":
				fmt.Println()
				printHelp()
				fmt.Println()

			case "/clear":
				fmt.Print("\033[H\033[2J")
				mu.Lock()
				cnt := onlineCount
				mu.Unlock()
				printDiscussionBanner(cnt, myUsername, isAnon)
				clrSystem.Println("\n  ── Screen cleared ──────────────────────────────────────────\n")

			case "/anon":
				isAnon = !isAnon
				if isAnon {
					clrAnon.Println("\n  👻  Anonymous mode ON – your identity is hidden.\n")
				} else {
					clrHighlight.Printf("\n  👤  Anonymous mode OFF – chatting as %s.\n\n", myUsername)
				}

			case "/status":
				mu.Lock()
				cnt := onlineCount
				mu.Unlock()
				clrInfo.Printf("\n  ● Status: Connected  •  %d users online\n", cnt)
				clrHighlight.Printf("  ● You: %s  (anon: %v)\n\n", myUsername, isAnon)

			default:
				clrError.Printf("\n  Unknown command: %s  (type /help for list)\n\n", cmd)
			}

			clrPrompt.Print("  » ")
			continue
		}

		effectiveUser := myUsername
		if isAnon {
			effectiveUser = "anonymous"
		}

		msg := api.ChatMessage{
			Username:  effectiveUser,
			Content:   line,
			Timestamp: nowHHMM(),
		}

		if err := conn.Send(msg); err != nil {
			clrError.Printf("\n  ✗ Failed to send: %v\n", err)
		}

		clrPrompt.Print("  » ")
	}

	clrSuccess.Println("\n  👋  Disconnected.\n")
}
