package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/glamour"
	"github.com/fatih/color"
	"github.com/itshivams/studex-cli/internal/api"
	"github.com/itshivams/studex-cli/internal/config"
)

func StartInteractiveMode() {
	color.Cyan(`
   _____ __            __             ________    ____
  / ___// /___  ______/ /__  _  __   / ____/ /   /  _/
  \__ \/ __/ / / / __  / _ \| |/_/  / /   / /    / /  
 ___/ / /_/ /_/ / /_/ /  __/>  <   / /___/ /____/ /   
/____/\__/\__,_/\__,_/\___/_/|_|   \____/_____/___/   
                                                      
`)
	color.White("Welcome to Studex CLI!\n\n")

	for {
		token := config.GetToken()
		if token != "" {
			loggedInMenu()
		} else {
			mainMenu()
		}
	}
}

func mainMenu() {
	options := []string{
		"1. Search User",
		"2. Login",
		"3. Signup",
		"4. Exit",
	}

	var choice string
	prompt := &survey.Select{
		Message: "Choose an option:",
		Options: options,
	}
	err := survey.AskOne(prompt, &choice)
	if err != nil {
		fmt.Println("Exited.")
		os.Exit(0)
	}

	switch choice {
	case options[0]:
		searchUserMenu()
	case options[1]:
		handleLogin()
	case options[2]:
		handleSignup()
	case options[3]:
		color.Green("Goodbye!")
		os.Exit(0)
	}
}

func loggedInMenu() {
	options := []string{
		"1. My Profile",
		"2. Messages",
		"3. My Notifications",
		"4. My Feed(Blogs)",
		"5. Search User",
		"6. Live Discussion",
		"7. Settings",
		"8. Logout",
		"9. Exit",
	}

	var choice string
	prompt := &survey.Select{
		Message: "Logged in menu:",
		Options: options,
	}
	err := survey.AskOne(prompt, &choice)
	if err != nil {
		fmt.Println("Exited.")
		os.Exit(0)
	}

	switch choice {
	case options[0]:
		color.Cyan("\nFetching My Profile...\n")
		profile, err := api.GetMyProfile()
		if err != nil {
			color.Red("Error fetching profile: %v\n", err)
		} else {
			displayProfile(profile)
		}
	case options[1]:
		messagesMenu()
	case options[2]:
		myNotificationsMenu()
	case options[3]:
		myFeedBlogsMenu()
	case options[4]:
		searchUserMenu()
	case options[5]:
		color.Yellow("\n[WIP] Live Discussion...\n")
	case options[6]:
		color.Yellow("\n[WIP] Settings...\n")
	case options[7]:
		config.ClearToken()
		color.Green("\nLogged out successfully!\n")
	case options[8]:
		color.Green("Goodbye!")
		os.Exit(0)
	}
}

func messagesMenu() {
	var recipient string
	err := survey.AskOne(&survey.Input{Message: "Enter recipient username:"}, &recipient, survey.WithValidator(survey.Required))
	if err != nil {
		return
	}

	color.Cyan("\nFetching chat history with %s...\n", recipient)
	thread, err := api.GetMessageThread(recipient)
	if err != nil {
		color.Red("Error fetching messages: %v\n", err)
		return
	}

	fmt.Printf("\n--- Chat with %s ---\n", recipient)
	if len(thread.Messages) == 0 {
		color.Yellow("No messages yet.\n")
	} else {
		for _, msg := range thread.Messages {
			timeStr := formatRelativeTime(msg.CreatedAt)
			sender := msg.From

			if msg.Type == "system" {
				color.HiBlack("[%s] System: %s\n", timeStr, msg.Text)
			} else if msg.Type == "call_log" {
				color.HiMagenta("[%s] %s: %s\n", timeStr, sender, msg.Text)
			} else {
				if msg.EncryptedData != nil {
					color.HiBlue("[%s] %s: 🔒 Encrypted Message\n", timeStr, sender)
				} else {
					if sender == recipient {
						color.HiGreen("[%s] %s: %s\n", timeStr, sender, msg.Text)
					} else {
						color.HiCyan("[%s] You: %s\n", timeStr, msg.Text)
					}
				}
			}
		}
	}
	fmt.Println("-----------------------")

	for {
		var replyText string
		prompt := &survey.Input{
			Message: "Type a message (or press enter to go back):",
		}
		err = survey.AskOne(prompt, &replyText)
		if err != nil {
			return
		}

		if replyText == "" {
			break
		}

		_, err = api.SendMessage(recipient, replyText)
		if err != nil {
			color.Red("Failed to send message: %v\n", err)
		} else {
			color.Green("Message sent to %s!\n", recipient)
		}
	}
}

func searchUserMenu() {
	options := []string{
		"1. By username",
		"2. By name",
		"3. Cancel",
	}

	var choice string
	prompt := &survey.Select{
		Message: "Search User:",
		Options: options,
	}
	err := survey.AskOne(prompt, &choice)
	if err != nil {
		return
	}

	switch choice {
	case options[0]:
		var username string
		survey.AskOne(&survey.Input{Message: "Enter username:"}, &username, survey.WithValidator(survey.Required))
		performSearch(username, true)
	case options[1]:
		var name string
		survey.AskOne(&survey.Input{Message: "Enter name:"}, &name, survey.WithValidator(survey.Required))
		performSearch(name, false)
	case options[2]:
		return
	}
}

func performSearch(query string, byUsername bool) {
	fmt.Printf("Searching for %s...\n", query)
	var users []api.UserProfile
	var err error

	if byUsername {
		users, err = api.SearchUser(query)
	} else {
		users, err = api.SearchUserByName(query)
	}

	if err != nil {
		color.Red("Search API error: %v", err)
		return
	}

	if len(users) == 0 {
		color.Yellow("No users found.\n")
		return
	}

	for _, u := range users {
		displayProfile(&u)
	}
}

func displayProfile(u *api.UserProfile) {
	color.Cyan("\n--- User Profile ---")
	fmt.Printf("Username: %s\n", u.Username)
	if u.Email != "" {
		fmt.Printf("Email: %s\n", u.Email)
	}
	fmt.Printf("Name: %s\n", u.FullName)
	if u.Gender != "" {
		fmt.Printf("Gender: %s\n", u.Gender)
	}
	fmt.Printf("Role: %s\n", u.Role)
	fmt.Printf("Status: %s\n", u.Status)
	fmt.Printf("Bio: %s\n", u.Bio)
	fmt.Printf("Location: %s\n", u.Location)
	if u.LastSeen != "" {
		fmt.Printf("Last Seen: %s\n", formatRelativeTime(u.LastSeen))
	}
	fmt.Printf("Organization: %s\n", u.Organization)
	fmt.Printf("Followers: %d | Following: %d\n", u.FollowersCount, u.FollowingCount)
	fmt.Printf("Posts: %d | Blogs: %d\n", u.PostsCount, u.BlogsCount)

	studexURL := fmt.Sprintf("https://studex.itshivam.in/profile/%s", u.Username)
	fmt.Printf("\nLinks:\n")
	fmt.Printf("- \033]8;;%s\033\\Studex Profile\033]8;;\033\\\n", studexURL)
	if u.Website != "" {
		fmt.Printf("- \033]8;;%s\033\\Website\033]8;;\033\\\n", u.Website)
	}
	if u.Linkedin != "" {
		fmt.Printf("- \033]8;;%s\033\\LinkedIn\033]8;;\033\\\n", u.Linkedin)
	}
	if u.Github != "" {
		fmt.Printf("- \033]8;;%s\033\\GitHub\033]8;;\033\\\n", u.Github)
	}
	if u.Instagram != "" {
		fmt.Printf("- \033]8;;%s\033\\Instagram\033]8;;\033\\\n", u.Instagram)
	}
}

func formatRelativeTime(timeStr string) string {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return timeStr
	}
	d := time.Since(t)

	if d < time.Minute {
		return "just now"
	} else if d < time.Hour {
		mins := int(d.Minutes())
		return fmt.Sprintf("%d m ago", mins)
	} else if d < 24*time.Hour {
		hrs := int(d.Hours())
		return fmt.Sprintf("%d hr ago", hrs)
	} else if d < 30*24*time.Hour {
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	} else if d < 365*24*time.Hour {
		months := int(d.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}
	years := int(d.Hours() / 24 / 365)
	if years == 1 {
		return "1 year ago"
	}
	return fmt.Sprintf("%d years ago", years)
}

func myNotificationsMenu() {
	color.Cyan("\nFetching My Notifications...\n")
	notifications, err := api.GetNotifications()
	if err != nil {
		color.Red("Error fetching notifications: %v\n", err)
		return
	}

	if len(notifications) == 0 {
		color.Yellow("You have no notifications.\n")
		return
	}

	color.HiCyan("\n=== My Notifications ===")
	for i, notif := range notifications {
		timeStr := formatRelativeTime(notif.CreatedAt)
		
		statusStr := color.YellowString("NEW")
		if notif.Read {
			statusStr = color.HiBlackString("READ")
		}

		fmt.Printf("\n[%d] %s [%s]\n", i+1, color.HiWhiteString(notif.Title), statusStr)
		color.HiBlack("    %s\n", notif.Message)
		color.HiBlack("    Received: %s\n", timeStr)
	}
	color.HiCyan("\n========================\n")
}

func myFeedBlogsMenu() {
	color.Cyan("\nFetching My Feed...\n")
	res, err := api.GetBlogs()
	if err != nil {
		color.Red("Error fetching blogs: %v\n", err)
		return
	}

	if res == nil || len(res.Items) == 0 {
		color.Yellow("No blogs found in your feed.\n")
		return
	}

	blogs := res.Items
	totalBlogs := len(blogs)
	limit := 5
	page := 1
	totalPages := (totalBlogs + limit - 1) / limit

	for {
		start := (page - 1) * limit
		end := start + limit
		if end > totalBlogs {
			end = totalBlogs
		}

		if start >= totalBlogs {
			color.Yellow("No more blogs to show.\n")
			break
		}

		color.HiCyan("\n=== My Feed (Blogs) - Page %d of %d (Total: %d) ===", page, totalPages, totalBlogs)
		for i := start; i < end; i++ {
			b := blogs[i]
			timeStr := formatRelativeTime(b.CreatedAt)

			fmt.Printf("\n[%d] %s\n", i+1, color.HiWhiteString(b.Title))
			color.HiMagenta("    By %s (@%s) • %s • %d min read\n", b.Author.FullName, b.Author.Username, timeStr, b.ReadTime)
			color.HiBlack("    %s\n", b.Excerpt)

			if len(b.Tags) > 0 {
				color.HiBlue("    Tags: %v\n", b.Tags)
			}

			color.HiGreen("    👍 %d | 💬 %d | 👁️ %d\n", b.LikesCount, b.CommentsCount, b.Views)
		}
		color.HiCyan("========================================================\n")

		var options []string
		for i := start; i < end; i++ {
			options = append(options, fmt.Sprintf("Read: %s", blogs[i].Title))
		}
		if end < totalBlogs {
			options = append(options, "Show more")
		}
		options = append(options, "Go back")

		var nextAction string
		prompt := &survey.Select{
			Message: "What would you like to do?",
			Options: options,
		}
		err = survey.AskOne(prompt, &nextAction)
		if err != nil || nextAction == "Go back" {
			break
		}

		if nextAction == "Show more" {
			page++
			continue
		}

		var selectedSlug string
		for i := start; i < end; i++ {
			if nextAction == fmt.Sprintf("Read: %s", blogs[i].Title) {
				selectedSlug = blogs[i].Slug
				break
			}
		}

		if selectedSlug != "" {
			viewBlog(selectedSlug)
		}
	}
}

func viewBlog(slug string) {
	color.Cyan("\nFetching blog details...\n")
	res, err := api.GetBlogView(slug)
	if err != nil {
		color.Red("Error fetching blog: %v\n", err)
		return
	}

	b := res.Blog
	timeStr := formatRelativeTime(b.CreatedAt)

	color.HiCyan("\n========================================================\n")
	color.HiWhite("%s\n", b.Title)
	color.HiMagenta("By %s (@%s) • %s • %d min read\n", b.Author.FullName, b.Author.Username, timeStr, b.ReadTime)

	if len(b.Tags) > 0 {
		color.HiBlue("Tags: %v\n", b.Tags)
	}

	likeCount := b.LikesCount
	if b.LikedBy != nil {
		likeCount = len(b.LikedBy)
	}

	color.HiGreen("👍 %d | 💬 %d | 👁️ %d\n", likeCount, b.CommentsCount, b.Views)
	color.HiCyan("--------------------------------------------------------\n")

	if b.Markdown != "" {
		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(100),
		)
		if err == nil {
			rendered, err := renderer.Render(b.Markdown)
			if err == nil {
				fmt.Print(rendered)
			} else {
				fmt.Println(b.Markdown)
			}
		} else {
			fmt.Println(b.Markdown)
		}
	}

	color.HiCyan("\n========================================================\n")

	var dummy string
	survey.AskOne(&survey.Input{Message: "Press Enter to go back..."}, &dummy)
}
