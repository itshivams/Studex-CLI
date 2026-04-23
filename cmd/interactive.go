package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
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
		"2. Send Message",
		"3. My Notifications",
		"4. My feed",
		"5. Blogs",
		"6. Search User",
		"7. Live Discussion",
		"8. Settings",
		"9. Logout",
		"10. Exit",
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
		color.Yellow("\n[WIP] Send Message...\n")
	case options[2]:
		color.Yellow("\n[WIP] My Notifications...\n")
	case options[3]:
		color.Yellow("\n[WIP] My feed...\n")
	case options[4]:
		color.Yellow("\n[WIP] Blogs...\n")
	case options[5]:
		searchUserMenu()
	case options[6]:
		color.Yellow("\n[WIP] Live Discussion...\n")
	case options[7]:
		color.Yellow("\n[WIP] Settings...\n")
	case options[8]:
		config.ClearToken()
		color.Green("\nLogged out successfully!\n")
	case options[9]:
		color.Green("Goodbye!")
		os.Exit(0)
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
