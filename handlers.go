package main

import (
	"math/rand"
	"strings"
	"time"
	"fmt"
	
	"github.com/yanzay/tbot/v2"
)

// Handle the /start command here
func (a *application) startHandler(m *tbot.Message) {
	msg := "hello"
	a.client.SendMessage(m.Chat.ID, msg, tbot.OptParseModeMarkdown)
}

// Handle the msg command here
func (a *application) msgHandler(m *tbot.Message) {
	magicDeer := `
   \ /   
    :   "\))    ((
     \   //      \\/"    o
      :  \\____// |  
   *    ~/~    ~\/~
       ( Q/  _/Q  ~       .
o       /  /     ,|
    '  (~~~)__.-\ |
        \'~~    | |    *
  .      |      | |
		`
	answer := map[int]string{
		0: "Да.",
		1: "Нет.",
		2: "Это не важно.",
		3: "...",
		4: "У тебя есть проблемы серьёзней.",
		5: "Да, хотя зря.",
		6: "Никогда.",
		7: "100%",
		8: "1 из 100.",
		9: "Попробуй ещё раз.",}
	msg := "Ты сделал что-то не так!"
	rand.Seed(time.Now().UnixNano())
	arr := strings.Split(m.Text, " или ")
	if len(arr) > 1 {
		rnd := rand.Intn(len(arr))
		msg = arr[rnd]
	} else {
		rnd := rand.Intn(10)
		msg = answer[rnd]
	}
	msg = fmt.Sprintf("```< %s > %s```", msg, magicDeer)
	a.client.SendChatAction(m.Chat.ID, tbot.ActionTyping)
	tsleep := rand.Intn(2500-500) + 500
	time.Sleep(time.Duration(tsleep) * time.Millisecond)
	a.client.SendMessage(m.Chat.ID, msg, tbot.OptParseModeMarkdown)
}
