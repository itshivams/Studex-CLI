package cmd

import (
	"encoding/json"
	"fmt"
	"os"

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
		color.Yellow("\nSignup is currently available via the web portal or mobile app.\n")
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
		color.Yellow("\n[WIP] Fetching My Profile...\n")
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

func handleLogin() {
	qs := []*survey.Question{
		{
			Name:     "username",
			Prompt:   &survey.Input{Message: "Username or Email:"},
			Validate: survey.Required,
		},
		{
			Name:     "password",
			Prompt:   &survey.Password{Message: "Password:"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Username string
		Password string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Logging in...")
	resp, err := api.Login(answers.Username, answers.Password)
	if err != nil {
		color.Red("Login failed: %v", err)
		return
	}

	if resp.Access != "" {
		err = config.SetToken(resp.Access)
		if err != nil {
			color.Red("Could not save token: %v", err)
			return
		}
		color.Green("\nLogin successful! Welcome %s\n", resp.Username)
	} else {
		color.Red("\nLogin failed: Token not received\n")
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
		color.Cyan("\n--- User Profile ---")
		fmt.Printf("Username: %s\n", u.Username)
		fmt.Printf("Name: %s\n", u.FullName)
		fmt.Printf("Role: %s\n", u.Role)
		fmt.Printf("Status: %s\n", u.Status)
		fmt.Printf("Bio: %s\n", u.Bio)
		fmt.Printf("Location: %s\n", u.Location)
		fmt.Printf("Organization: %s\n", u.Organization)
		fmt.Printf("Followers: %d | Following: %d\n", u.FollowersCount, u.FollowingCount)
		fmt.Printf("Posts: %d | Blogs: %d\n", u.PostsCount, u.BlogsCount)

		b, _ := json.MarshalIndent(u, "", "  ")
		color.White("\nRaw Data:\n%s\n", string(b))
	}
}
