// Package main provides ...
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	// "strconv"
	"time"
	// "date"

	tb "gopkg.in/tucnak/telebot.v2"
)

var S Schedules

func main() {

	InitDatabase()

	Fetch()

	poller := &tb.LongPoller{Timeout: 10 * time.Second}
	middleware := tb.NewMiddlewarePoller(poller, func(udp *tb.Update) bool {
		// fmt.Printf("%+v\n", udp.Message)
		return true
	})
	b, err := tb.NewBot(tb.Settings{
		Token:  "828734497:AAHiqTY1c0g3lVkqm8lO-h21Q_vhKkdvVyM",
		Poller: middleware,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	league(b)
	schedule(b)

	b.Start()
}
func schedule(b *tb.Bot) {

	b.Handle("/schedule", func(m *tb.Message) {
		leagues := func() string {
			var ret string = `

`
			for _, league := range S.League {
				ret += fmt.Sprintf(`
<a>%s </a><strong>%s: </strong><a href="https://splatoon2.ink/assets/splatnet%s">%s</a> / <a href="https://splatoon2.ink/assets/img%s">%s</a>
                `, time.Unix(league.StartTime, 0).Format("15:04"), league.Rule.Name, league.StageA.Image, league.StageA.Name, league.StageB.Image, league.StageB.Name)
			}
			return ret
		}()
		b.Send(m.Chat, leagues, tb.ModeHTML)
	})

}

var GlobalMsg map[int]tb.Message

func league(b *tb.Bot) {
	GlobalMsg = map[int]tb.Message{}

	confirm := tb.InlineButton{Unique: "confirm",
		Text: "参与组排",
	}
	cancel := tb.InlineButton{Unique: "cancel",
		Text: "取消组排",
	}

	// 点击加入之后刷新message
	generateAfterAdd := func(m *tb.Message) {
		invitation, err := FetchLeaugeInvitation(m.ID)
		if err != nil {
			fmt.Println(err)
			return
		}
		msg := GlobalMsg[int(m.ID)]

		markup1 := &tb.ReplyMarkup{InlineKeyboard: ([][]tb.InlineButton{{confirm, cancel}})}
		body := ""
		if invitation.Member1 != "" {
			body += fmt.Sprintf(`
<a href="tg://user?id=%d">%s</a>参与`, invitation.MemberID1, invitation.Member1)
		}
		if invitation.Member2 != "" {
			body += fmt.Sprintf(`
<a href="tg://user?id=%d">%s</a>参与`, invitation.MemberID2, invitation.Member2)
			body += `
            <strong>车已满员,等待发车!</strong>`
			b.Edit(&msg, body, tb.ModeHTML)
			return
		}
		if invitation.Member3 != "" {
			body += fmt.Sprintf(`
<a href="tg://user?id=%d">%s</a>参与`, invitation.MemberID3, invitation.Member3)
		}
		if invitation.Member4 != "" {
			body += fmt.Sprintf(`
<a href="tg://user?id=%d">%s</a>参与`, invitation.MemberID4, invitation.Member4)
			body += `
            <strong>车已满员,等待发车!</strong>`
			b.Edit(&msg, body, tb.ModeHTML)
			return
		}

		b.Edit(&msg, body, tb.ModeHTML, markup1)
		return
	}
	b.Handle(&confirm, func(m *tb.Callback) {
		err := AddLeagueInvitationMember(m.Message.ID, m.Sender.ID, m.Sender.Username)
		if err != nil {
			fmt.Println(err)
			return
		}
		generateAfterAdd(m.Message)
	})

	b.Handle("/league", func(m *tb.Message) {
		leaguebtns := [][]tb.InlineButton{}
		for _, league := range S.League {
			content := fmt.Sprintf(`%s %s: %s/%s`, time.Unix(league.StartTime, 0).Format("15:04"), league.Rule.Name, league.StageA.Name, league.StageB.Name)
			btn := tb.InlineButton{Text: content,
				Unique: "abc"}
			b.Handle(&btn, func(m *tb.Callback) {
				err := CreateLeagueInvitation(m.Message.ID, m.Sender.ID, m.Sender.Username, "aaa")
				if err != nil {
					fmt.Println(err)
					return
				}

				generateAfterAdd(m.Message)
			})
			l := []tb.InlineButton{btn}
			leaguebtns = append(leaguebtns, l)
		}
		markup := &tb.ReplyMarkup{InlineKeyboard: (leaguebtns)}
		body := fmt.Sprintf(`
<a href="tg://user?id=%d">%s %s</a>发起组排请求
        `, m.Sender.ID,
			m.Sender.FirstName,
			m.Sender.LastName)
		sm, _ := b.Send(m.Chat, body, tb.ModeHTML, markup)
		GlobalMsg[sm.ID] = *sm
	})
}

// Fetch 获取组排信息
func Fetch() {
	fetchSchedules := func() {
		resp, err := http.Get("https://splatoon2.ink/data/schedules.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		schedules := Schedules{}
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		parseErr := json.Unmarshal(body, &schedules)
		if parseErr != nil {
			fmt.Println(parseErr)
			return
		}
		S = schedules
	}

	tiker := time.NewTicker(time.Minute * 10)
	go func() {
		fetchSchedules()
		for _ = range tiker.C {
			fetchSchedules()
		}
	}()
}
