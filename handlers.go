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
	msg := "Что может *сакральный олень?*\n\nОтветит на вопрос:\nЗадай вопрос и ты получишь ответ." +
	"\nНа вопрос ответом должны быть да/нет.\n\nВыберет из нескольких вариантов:" +
	"\nРаздели варианты союзом *или*.\nМинимум - 2 варианта, максимума нет.\nНе забудь про *пробелы*, пример:" +
	"\nЛечь спать *или* дочитать мангу?\n" +
	"\nПредскажет будущее:\nДля начала выбери свой знак зодиака.\nОтправь его в чат на русском языке." +
	"\nДалее введи соответствующую команду:\n*/today* - предсказание на сегодня\n*/tomorrow* - предсказание на завтра" +
	//сделаю "\n*/daily* - ежедневные предсказания" + 
	"\n\nВнимание:\n*Сакрального оленя* нельзя тревожить," + 
	"\nзадавая тот же вопрос несколько раз.\nТакже нельзя задавать любые\nвопросы связанные с *оленем*.\n" +
	"\nВ случае нарушения правил выше\nты не получишь достоверного ответа."
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
		"дева":     "virgo",
		"весы":     "libra",
		"скорпион": "scorpio",
		"стрелец":  "sagittarius",
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
	msg = "ok"
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
		msg = "err"
	}
	defer conn.Close(context.Background())
	}
	if strings.ToLower(strings.TrimRight(m.Text, " .!")) == "спасибо"{
		msg = "Пожалуйста"
	}
	msg = fmt.Sprintf("```\n< %s > %s```", msg, magicDeer)
	a.client.SendChatAction(m.Chat.ID, tbot.ActionTyping)
	tsleep := rand.Intn(1000-200) + 200
	time.Sleep(time.Duration(tsleep) * time.Millisecond)
	a.client.SendMessage(m.Chat.ID, msg, tbot.OptParseModeMarkdown)
}
