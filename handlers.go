package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/yanzay/tbot/v2"
)

// Handle the /start command here
func (a *application) startHandler(m *tbot.Message) {
	msg := "hello"
	a.client.SendMessage(m.Chat.ID, msg, tbot.OptParseModeMarkdown)
}

// Handle the msg command here
func (a *application) msgHandler(m *tbot.Message) {
	msg := "Ты сделал что-то не так!"
	rand.Seed(time.Now().UnixNano())
	arr := strings.Split(m.Text, " или ")
	if len(arr) > 1 {
		rnd := (rand.Intn(len(arr)))
		msg = (arr[rnd])
	} else {
		msg = "Допускается минимум два варианта!"
	}
	a.client.SendMessage(m.Chat.ID, msg, tbot.OptParseModeMarkdown)
}
