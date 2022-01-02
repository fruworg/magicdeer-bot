package main

import (
	"math/rand"
	"net/http"
	"context"
	"strings"
	"time"
	"fmt"
	"log"
	"os"
	
	"github.com/PuerkitoBio/goquery"
	"github.com/yanzay/tbot/v2"
	"github.com/jackc/pgx/v4"
)

//deer by asciiart.eu
//flower by eng50232@leonis.nus.sg
var magicDeer = `
 \ /   .              *
* :       ))    ((
   \     // (") \\   '      .
    :    \\_\)/_//  
 .   \  ~/~  ' ~\~\
       ( Q/  _/Q  ~     o
o       /  /     ,|
    '  (~~~)__.-\ |
        \'~~    | |   *
  .      |      | |
		`

// Handle the /start command here
func (a *application) startHandler(m *tbot.Message) {
	msg := "Для гороскопа введи свой знак зодиака на русском.\nДля выбора нескольких вариантов:"+
	"\nГоланг или Руби или Питон?\n\nДля ответа на вопорс: Я крутой?"
	a.client.SendMessage(m.Chat.ID, msg, tbot.OptParseModeMarkdown)
}

// Handle the msg command here
func (a *application) msgHandler(m *tbot.Message) {
	msg := "Ты сделал что-то не так"
	signs := map[string]string{
		"овен":     "aries",
		"телец":    "taurus",
		"близнецы": "gemini",
		"рак":      "cancer",
		"лев":      "leo",
		"дева":     "virgio",
		"весы":     "libra",
		"скорпион": "scorpio",
		"стрелец":  "saggitarius",
		"козерог":  "capricorn",
		"водолей":  "aquarius",
		"рыбы":     "pisces"}
	
	if signs[strings.ToLower(m.Text)] != "" {
		day := "tod"
		res, err := http.Get("https://ignio.com/r/daily/" + day + "/" + signs[strings.ToLower(m.Text)] + ".html")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		doc.Find(`div[style="margin: 20px 0;"]`).Each(func(i int, s *goquery.Selection) {
			msg = fmt.Sprintf("Гороскоп для тебя, %s: \n%s", m.Text, strings.TrimSpace(s.Text()))
		})
	} else {
		answer := map[int]string{
			0: "Да",
			1: "Нет",
			2: "Это не важно",
			3: "...",
			4: "У тебя есть проблемы серьёзней",
			5: "Да, хотя зря",
			6: "Никогда",
			7: "100%",
			8: "1 из 100",
			9: "Попробуй ещё раз"}
		rand.Seed(time.Now().UnixNano())
		arr := strings.Split(m.Text, " или ")
		if len(arr) > 1 {
			msg = "Ты не оставил мне выбора"
			for i := 0; i < len(arr)-1; i++ {
				if strings.TrimRight(arr[i], "?") != strings.TrimRight(arr[i+1], "?") {
					rnd := rand.Intn(len(arr))
					msg = strings.TrimRight(arr[rnd], "?")
					continue
				}
			}
		} else {
			rnd := rand.Intn(10)
			msg = answer[rnd]
		}
	}
	if m.Text == "тест"{
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var name string
	var weight int64
	err = conn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	msg = fmt.Sprintf("%s; %s", name, weight)
	}
	msg = fmt.Sprintf("```\n< %s > %s```", msg, magicDeer)
	a.client.SendChatAction(m.Chat.ID, tbot.ActionTyping)
	tsleep := rand.Intn(2500-500) + 500
	time.Sleep(time.Duration(tsleep) * time.Millisecond)
	a.client.SendMessage(m.Chat.ID, msg, tbot.OptParseModeMarkdown)
}
