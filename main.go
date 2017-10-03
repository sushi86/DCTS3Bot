package main

import (
	"fmt"
	"github.com/sushi86/DCTS3Bot/discord"
	"github.com/sushi86/DCTS3Bot/telegram"
	"github.com/sushi86/DCTS3Bot/ts3"
)

func main() {
	fmt.Println("Bot is starting...")

	go discord.ConnectDc()

	telegram.Connect()
	go telegram.Poll()

	ts3.Connect()
	go ts3.Poll()

	for {select {}}
}
