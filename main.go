package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bluele/slack"
)

var (
	channelName string
	msg         string
)

func getToken() (string, error) {
	envToken := os.Getenv("SLACK_TOKEN")
	if envToken != "" {
		return envToken, nil
	}

	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		return "", errors.New("could not lookup home directory by environment variable")
	}

	tokenPath := fmt.Sprintf("%s/.slack_token", homeDir)
	fd, err := os.Open(tokenPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New("could not find .slack_token file")
		} else {
			return "", errors.New("issue opening .slack_token file in known location")
		}
	}

	// Slack tokens are 51 characters long but I wanted to leave some room for future changes
	// Plus 100 bytes isn't a huge amount of memory
	tokenBytes := make([]byte, 100, 100)

	size, err := fd.Read(tokenBytes)
	if err != nil {
		return "", errors.New("issue reading from .slack_token file")
	}
	if size <= 0 {
		return "", errors.New("no token read from file")
	}

	return string(tokenBytes[:size]), nil
}

func init() {
	flag.StringVar(&channelName, "channel", "", "Name of the channel to send the message to without a '#'")
	flag.StringVar(&msg, "message", "", "Message to be sent to the Slack channel")
	flag.Parse()
}

func main() {
	var err error

	token, err := getToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error looking up token: %s\n", err.Error())
		os.Exit(1)
	}

	sl := slack.New(strings.TrimSpace(token))

	if channelName == "" {
		fmt.Fprintf(os.Stderr, "No channel provided. Please re-run and supply a channel name\n")
		os.Exit(1)
	}
	channel, err := sl.FindChannelByName(channelName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding channel: %s\n", err.Error())
		os.Exit(1)
	}

	if msg == "" {
		fmt.Fprintf(os.Stderr, "No message provided. Please re-run and supply a message to post\n")
		os.Exit(1)
	}
	err = sl.ChatPostMessage(channel.Id, msg, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending message to channel\n")
		os.Exit(1)
	}

	fmt.Println("Success")
}
