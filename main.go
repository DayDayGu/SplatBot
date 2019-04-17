// Package main provides ...
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"strings"
	// "date"

	tb "gopkg.in/tucnak/telebot.v2"
)

// S Splat轮转状态
var S Schedules

func main() {
	// 初始化用于splat数据库
	InitDatabase()

	// 获取Splat轮转状态
	Fetch()

	// 初始化bot
	poller := &tb.LongPoller{Timeout: 10 * time.Second}
	middleware := tb.NewMiddlewarePoller(poller, func(udp *tb.Update) bool {
		// fmt.Printf("%+v\n", udp.Message)
		return true
	})
	b, err := tb.NewBot(tb.Settings{
		// Modify Token here.
		Token:  "828734497:AAHncNJ0penMnYm0mqfaELVSdaunVACBtF4",
		Poller: middleware,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	// league命令
	league(b)
	// schedule命令
	schedule(b)

	// 启动bot
	b.Start()
}

func schedule(b *tb.Bot) {

	b.Handle("/schedule", func(m *tb.Message) {
		leagues := func() string {
			ret := `

`
			for _, league := range S.League {
				ret += fmt.Sprintf(`
<a>%s </a><strong>%s: </strong><a href="https://splatoon2.ink/assets/splatnet%s">%s</a> / <a href="https://splatoon2.ink/assets/img%s">%s</a>
                `, time.Unix(league.StartTime, 0).Format("15:04"),
					league.Rule.Name,
					league.StageA.Image,
					league.StageA.Name,
					league.StageB.Image,
					league.StageB.Name)
			}
			return ret
		}()
		b.Send(m.Chat, leagues, tb.ModeHTML)
	})

}

func league(b *tb.Bot) {

	confirm := tb.InlineButton{Unique: "confirm",
		Text: "参与组排",
	}

	cancel := tb.InlineButton{Unique: "cancel",
		Text: "取消组排",
	}

	// 点击选择组排按钮的handler
	b.Handle(&confirm, func(m *tb.Callback) {
		err := AddLeagueInvitationMember(m.Message.ID, m.Sender.ID, m.Sender.Username)
		if err != nil {
			fmt.Println(err)
			return
		}
		updateAddMessage(b, m.Message, [][]tb.InlineButton{{confirm, cancel}})
	})

	b.Handle(&cancel, func(m *tb.Callback) {
		err := DeleteLeagueInvitationMember(m.Message.ID, m.Sender.ID)
		if err != nil {
			fmt.Println(err)
			return
		}
		updateAddMessage(b, m.Message, [][]tb.InlineButton{{confirm, cancel}})
	})
	b.Handle("/league", func(m *tb.Message) {
		leaguebtns := [][]tb.InlineButton{}
		for k, league := range S.League {
			content := fmt.Sprintf(`%s %s`, time.Unix(league.StartTime, 0).Format("15:04"), league.Rule.Name)
			data := fmt.Sprintf("%s&&%d", league.Rule.Name, league.StartTime)
			btn := tb.InlineButton{Text: content,
				Data:   data,
				Unique: strconv.Itoa(k)}
			b.Handle(&btn, func(m *tb.Callback) {
				data := strings.Split(m.Data, "&&")
				startTime, _ := strconv.Atoi(data[1])
				err := CreateLeagueInvitation(m.Message.ID, m.Sender.ID, m.Sender.Username, data[0], startTime)
				if err != nil {
					fmt.Println(err)
					return
				}
				Notify(b, m.Message, int64(startTime))
				updateAddMessage(b, m.Message, [][]tb.InlineButton{{confirm, cancel}})

			})
			l := []tb.InlineButton{btn}
			leaguebtns = append(leaguebtns, l)
		}
		markup := &tb.ReplyMarkup{InlineKeyboard: (leaguebtns)}
		body := fmt.Sprintf(`
<a href="tg://user?id=%d">%s %s</a>你好，点击可以发起组排哦
        `, m.Sender.ID,
			m.Sender.FirstName,
			m.Sender.LastName)
		b.Send(m.Chat, body, tb.ModeHTML, markup)
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
		for range tiker.C {
			fetchSchedules()
		}
	}()
}

// Notify 增加到组排时间的通知
func Notify(b *tb.Bot, m *tb.Message, des int64) {
	current := time.Now().Unix()
	duration := time.Duration((des - current))

	fmt.Printf("duration is %d\n", duration)

	go func(ch <-chan time.Time) {

		<-ch
		MarkLeagueInvitation(m.ID, LeagueStatusEnd)
		invitation, _ := FetchLeaugeInvitation(m.ID)
		if invitation.Member1 == "" && invitation.Member2 == "" && invitation.Member3 == "" && invitation.Member4 == "" {
			return
		}

		body := fmt.Sprintf(`
        <a>当前模式：</a><strong>%s</strong>`, invitation.Rule)
		if invitation.Member1 != "" {
			body += fmt.Sprintf(`
<a href="tg://user?id=%d">@%s </a>`, invitation.MemberID1, invitation.Member1)
		}
		if invitation.Member2 != "" {
			body += fmt.Sprintf(`
<a href="tg://user?id=%d">@%s </a>`, invitation.MemberID2, invitation.Member2)
		}
		if invitation.Member3 != "" {
			body += fmt.Sprintf(`
<a href="tg://user?id=%d">@%s </a>`, invitation.MemberID3, invitation.Member3)
		}
		if invitation.Member4 != "" {
			body += fmt.Sprintf(`
<a href="tg://user?id=%d">@%s </a>`, invitation.MemberID4, invitation.Member4)
		}
		body += `
组排啦！鸽子？不存在的!
        `
		b.Send(m.Chat, body, tb.ModeHTML)
	}(time.After(time.Second * duration))
}

// 点击加入之后刷新message
func updateAddMessage(b *tb.Bot, m *tb.Message, btns [][]tb.InlineButton) {
	invitation, err := FetchLeaugeInvitation(m.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	msg := *m

	markup1 := &tb.ReplyMarkup{InlineKeyboard: btns}
	body := fmt.Sprintf(`<a>组排模式：</a><strong>%s</strong>
<a>组排时间：</a><strong>%s</strong>
<a>参与乌贼：</a>`, invitation.Rule, time.Unix(int64(invitation.StartTime), 0).Format("15:04"))
	if invitation.Member1 != "" {
		body += fmt.Sprintf(`
<a href="tg://user?id=%d">%s</a>`, invitation.MemberID1, invitation.Member1)
	}
	if invitation.Member2 != "" {
		body += fmt.Sprintf(`
<a href="tg://user?id=%d">%s</a>`, invitation.MemberID2, invitation.Member2)
	}
	if invitation.Member3 != "" {
		body += fmt.Sprintf(`
<a href="tg://user?id=%d">%s</a>`, invitation.MemberID3, invitation.Member3)
	}
	if invitation.Member4 != "" {
		body += fmt.Sprintf(`
<a href="tg://user?id=%d">%s</a>`, invitation.MemberID4, invitation.Member4)
		body += `
            <strong>车已满员,等待发车!</strong>`
		// 若满4人，则组排成功
		b.Edit(&msg, body, tb.ModeHTML)
		return
	}

	b.Edit(&msg, body, tb.ModeHTML, markup1)
	return
}
