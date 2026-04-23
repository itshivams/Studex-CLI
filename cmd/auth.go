package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/itshivams/studex-cli/internal/api"
	"github.com/itshivams/studex-cli/internal/config"
)

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
		_ = config.SetStoredUsername(resp.Username)
		color.Green("\nLogin successful! Welcome %s\n", resp.Username)
	} else {
		color.Red("\nLogin failed: Token not received\n")
	}
}

func handleSignup() {
	color.Cyan("\n--- Studex Signup ---")

	var username string
	for {
		err := survey.AskOne(&survey.Input{Message: "Username (max 25 chars):"}, &username, survey.WithValidator(survey.Required))
		if err != nil {
			fmt.Println("Signup cancelled.")
			return
		}

		available, msg, err := api.CheckUsername(username)
		if err != nil {
			color.Red("Error checking username: %v\n", err)
			continue
		}
		if !available {
			color.Red("Username not available: %s\n", msg)
			continue
		}
		color.Green("Username is available!\n")
		break
	}

	var step1Answers struct {
		FullName    string
		Email       string
		PhoneNumber string
	}

	qs := []*survey.Question{
		{
			Name:     "FullName",
			Prompt:   &survey.Input{Message: "Full Name:"},
			Validate: survey.Required,
		},
		{
			Name:     "Email",
			Prompt:   &survey.Input{Message: "Email:"},
			Validate: survey.Required,
		},
		{
			Name:   "PhoneNumber",
			Prompt: &survey.Input{Message: "Phone Number (Optional):"},
		},
	}

	err := survey.Ask(qs, &step1Answers)
	if err != nil {
		fmt.Println("Signup cancelled.")
		return
	}

	fmt.Println("Sending OTP...")
	req1 := api.RegisterStep1Request{
		Username:    username,
		FullName:    step1Answers.FullName,
		Email:       step1Answers.Email,
		PhoneNumber: step1Answers.PhoneNumber,
	}

	resp1, err := api.RegisterStep1(req1)
	if err != nil {
		color.Red("Signup failed: %v\n", err)
		return
	}

	color.Green("OTP sent to your email!\n")

	var otp string
	err = survey.AskOne(&survey.Input{Message: "Enter 6-digit OTP:"}, &otp, survey.WithValidator(survey.Required))
	if err != nil {
		fmt.Println("Signup cancelled.")
		return
	}

	fmt.Println("Verifying OTP...")
	req2 := api.RegisterStep2Request{
		RequestId: resp1.RequestId,
		Otp:       otp,
	}

	resp2, err := api.RegisterStep2(req2)
	if err != nil {
		color.Red("OTP verification failed: %v\n", err)
		return
	}

	color.Green("Email verified successfully!\n")

	var password string
	err = survey.AskOne(&survey.Password{Message: "Create a Password (min 8 chars):"}, &password, survey.WithValidator(survey.Required))
	if err != nil {
		fmt.Println("Signup cancelled.")
		return
	}

	fmt.Println("Completing registration...")
	req3 := api.RegisterStep3Request{
		TempToken: resp2.TempToken,
		Password:  password,
	}

	resp3, err := api.RegisterStep3(req3)
	if err != nil {
		color.Red("Registration failed: %v\n", err)
		return
	}

	color.Green("\nWelcome to Studex! Account created successfully.\n")

	if resp3.Token != "" {
		err = config.SetToken(resp3.Token)
		if err != nil {
			color.Red("Could not save token: %v\n", err)
		} else {
			_ = config.SetStoredUsername(resp3.User.Username)
			color.Green("You are now logged in as %s\n", resp3.User.Username)
		}
	}
}
