// Package main provides ...
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"errors"
	"regexp"
	"strings"

	// "date"

	faker "SplatBot/faker"
	// "github.com/PangPangPangPangPang/SplatBot/faker"
	tb "gopkg.in/tucnak/telebot.v2"
)

// S Splat轮转状态
var S Schedules

// Sa 打工状态
var Sa Salmon

var SalmonMap map[string][]byte = make(map[string][]byte)

var Bot *tb.Bot

func main() {
	ClearTmpPath()
	// 初始化用于splat数据库
	InitDatabase()

	// 获取Splat轮转状态
	Fetch()
	// DownloadSalmon(Sa.Details[0])

	// 初始化bot
	poller := &tb.LongPoller{Timeout: 10 * time.Second}
	middleware := tb.NewMiddlewarePoller(poller, func(udp *tb.Update) bool {
		// go commonQuery(udp)
		return true
	})
	b, err := tb.NewBot(tb.Settings{
		// Modify Token here.
		Token:  os.Getenv("SPLAT_BOT_TOKEN"),
		Poller: middleware,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	Bot = b

	// league命令
	league(b)
	// schedule命令
	schedule(b)
	// gachi命令
	gachi(b)
	// salmon命令
	salmon(b)
	salmonRaw(b)

	b.Handle(tb.OnText, func(m *tb.Message) {
		commonQuery(m)
	})

	// 启动bot
	b.Start()
}
func gachi(b *tb.Bot) {
	b.Handle("/gachi", func(m *tb.Message) {
		leagues := func() string {
			ret := `
<strong>Gachi</strong>
`
			loc, _ := time.LoadLocation("Asia/Shanghai")
			for idx, gachi := range S.Gachi {
				if idx > 5 {
					break
				}
				ret += fmt.Sprintf(`
<a>%s </a><strong>%s: </strong><a href="https://splatoon2.ink/assets/splatnet%s">%s</a> / <a href="https://splatoon2.ink/assets/splatnet%s">%s</a>
                `, time.Unix(gachi.StartTime, 0).In(loc).Format("15:04"),
					gachi.Rule.Name,
					gachi.StageA.Image,
					gachi.StageA.Name,
					gachi.StageB.Image,
					gachi.StageB.Name)
			}
			return ret
		}()
		b.Send(m.Chat, leagues, tb.ModeHTML)
	})

}
func salmonRaw(b *tb.Bot) {
	b.Handle("/salmonraw", func(m *tb.Message) {
		salmons := func() string {
			ret := `
<strong>Salmon</strong>
`
			loc, _ := time.LoadLocation("Asia/Shanghai")
			ret += fmt.Sprintf(`
<strong>Duration:</strong>
<strong>from:</strong>%s
<strong>to:</strong>%s
            `, time.Unix(Sa.Details[0].StartTime, 0).In(loc).Format("Jan 2 15:04"), time.Unix(Sa.Details[0].EndTime, 0).In(loc).Format("Jan 2 15:04"))
			ret += `
<strong>Weapons:</strong>`
			for _, weapon := range Sa.Details[0].Weapons {
				if weapon.CoopSpecialWeapon != nil {
					ret += fmt.Sprintf(`
%s`, weapon.CoopSpecialWeapon.Name)
				} else if weapon.Weapon != nil {
					ret += fmt.Sprintf(`
%s
            `, weapon.Weapon.Name)
				}
			}
			if Sa.Details[0].Stage.Name != "" {
				ret += fmt.Sprintf(`

<strong>Stage:</strong>
<a href="https://splatoon2.ink/assets/splatnet%s">%s</a>
                `, Sa.Details[0].Stage.Image, Sa.Details[0].Stage.Name)
			}
			return ret
		}()
		b.Send(m.Chat, salmons, tb.ModeHTML)
	})
}

func salmon(b *tb.Bot) {
	b.Handle("/salmon", func(m *tb.Message) {
		show := func() {
			path := fmt.Sprintf("%s%d", faker.TempPath(), Sa.Details[0].StartTime)
			var photo *tb.Photo
			if SalmonMap[path] != nil {
				err := json.Unmarshal(SalmonMap[path], &photo)
				if err != nil {
					return
				}
			} else {
				file := tb.FromDisk(path)
				photo = &tb.Photo{File: file}
			}
			b.Send(m.Chat, photo)

			bytes, _ := json.Marshal(photo)
			SalmonMap[path] = bytes
		}
		if faker.Exist(fmt.Sprintf("%d", Sa.Details[0].StartTime)) {
			go func() {
				show()
			}()
		} else if len(Sa.Details) > 0 {
			// 拼接打工内容图片
			go func() {
				DownloadSalmon(Sa.Details)
				show()
			}()
		}
	})

}

func schedule(b *tb.Bot) {
	b.Handle("/schedule", func(m *tb.Message) {
		leagues := func() string {
			ret := `
<strong>League</strong>
`
			loc, _ := time.LoadLocation("Asia/Shanghai")
			for idx, league := range S.League {
				if idx > 5 {
					break
				}
				ret += fmt.Sprintf(`
<a>%s </a><strong>%s: </strong><a href="https://splatoon2.ink/assets/splatnet%s">%s</a> / <a href="https://splatoon2.ink/assets/splatnet%s">%s</a>
                `, time.Unix(league.StartTime, 0).In(loc).Format("15:04"),
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
		Text: "取消参与",
	}

	two := tb.InlineButton{Unique: "two",
		Text: "立即双排"}

	four := tb.InlineButton{Unique: "four",
		Text: "立即四排"}

	rightnow := tb.InlineButton{Unique: "rightnow",
		Text: "通知发车"}
	// 点击选择组排按钮的handler
	b.Handle(&confirm, func(m *tb.Callback) {
		err := AddLeagueInvitationMember(m.Message.ID, m.Sender.ID, m.Sender.Username)
		if err != nil {
			fmt.Println(err)
			return
		}
		updateAddMessage(b, m.Message, [][]tb.InlineButton{{confirm, cancel, rightnow}})
	})

	b.Handle(&cancel, func(m *tb.Callback) {
		err := DeleteLeagueInvitationMember(m.Message.ID, m.Sender.ID)
		if err != nil {
			fmt.Println(err)
			return
		}
		updateAddMessage(b, m.Message, [][]tb.InlineButton{{confirm, cancel, rightnow}})
	})

	b.Handle(&two, func(m *tb.Callback) {
		cur, _ := currentLeague()
		err := CreateLeagueInvitation(m.Message.ID, m.Sender.ID, m.Sender.Username, cur.Rule.Name, cur.StartTime, LeagueTypeDouble)
		if err != nil {
			fmt.Println(err)
			return
		}
		updateAddMessage(b, m.Message, [][]tb.InlineButton{{confirm, cancel, rightnow}})

	})
	b.Handle(&four, func(m *tb.Callback) {
		cur, _ := currentLeague()
		err := CreateLeagueInvitation(m.Message.ID, m.Sender.ID, m.Sender.Username, cur.Rule.Name, cur.StartTime, LeagueTypeFour)
		if err != nil {
			fmt.Println(err)
			return
		}
		updateAddMessage(b, m.Message, [][]tb.InlineButton{{confirm, cancel, rightnow}})
	})

	b.Handle(&rightnow, func(m *tb.Callback) {
		invitation, _ := FetchLeaugeInvitation(m.Message.ID)
		if invitation.Member1 == "" && invitation.Member2 == "" && invitation.Member3 == "" && invitation.Member4 == "" {
			return
		}
		if invitation.Member1 != m.Sender.Username && invitation.Member2 != m.Sender.Username && invitation.Member3 != m.Sender.Username && invitation.Member4 != m.Sender.Username {
			return
		}

		if invitation.Status == LeagueStatusEnd {
			return
		}

		MarkLeagueInvitation(m.Message.ID, LeagueStatusEnd)

		sendNotify(b, m.Message, invitation)

	})
	loc, _ := time.LoadLocation("Asia/Shanghai")
	b.Handle("/league", func(m *tb.Message) {
		leaguebtns := [][]tb.InlineButton{{two, four}}

		for k, league := range futureLeague() {
			content := fmt.Sprintf(`%s %s`, time.Unix(league.StartTime, 0).In(loc).Format("15:04"), league.Rule.Name)
			data := fmt.Sprintf("%s&&%d", league.Rule.Name, league.StartTime)
			btn := tb.InlineButton{Text: content,
				Data:   data,
				Unique: strconv.Itoa(k)}
			b.Handle(&btn, func(m *tb.Callback) {
				data := strings.Split(m.Data, "&&")
				startTime, _ := strconv.Atoi(data[1])
				err := CreateLeagueInvitation(m.Message.ID, m.Sender.ID, m.Sender.Username, data[0], int64(startTime), LeagueTypeFour)
				if err != nil {
					fmt.Println(err)
					return
				}
				Notify(b, m.Message, int64(startTime))
				updateAddMessage(b, m.Message, [][]tb.InlineButton{{confirm, cancel, rightnow}})

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

func commonQuery(msg *tb.Message) {
	r, err := regexp.MatchString(".*拉稀.*[B|b]ot", msg.Text)
	if err == nil && r {
		path := fmt.Sprintf("%s%s", faker.ResourcePath(), "animation.mp4")
		file := tb.FromDisk(path)
		video := &tb.Video{File: file}
		time.Sleep(300 * time.Millisecond)
		Bot.Reply(msg, video)
	}
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

	fetchSalmon := func() {
		resp, err := http.Get("https://splatoon2.ink/data/coop-schedules.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		ret, err := UnmarshalSalmon(body)
		if err != nil {
			fmt.Println(err)
			return
		}
		Sa = ret
	}

	fetchSalmon()
	fetchSchedules()

	tiker := time.NewTicker(time.Minute * 10)
	go func() {
		for range tiker.C {
			fetchSalmon()
			fetchSchedules()
		}
	}()
}

func currentLeague() (Battle, error) {
	leagues := S.League
	for _, league := range leagues {
		if time.Now().Unix() > league.StartTime && (time.Now().Unix()-league.StartTime < int64(time.Hour/time.Second*2)) {
			return league, nil
		}
	}
	return Battle{}, errors.New("Can not find current league battle")
}

func futureLeague() []Battle {
	leagues := S.League
	ret := []Battle{}
	for _, league := range leagues {
		if time.Now().Unix() < league.StartTime {
			ret = append(ret, league)
		}
	}
	return ret
}

// Notify 增加到组排时间的通知
func Notify(b *tb.Bot, m *tb.Message, des int64) {
	current := time.Now().Unix()
	duration := time.Duration((des - current))

	fmt.Printf("duration is %d\n", duration)

	go func(ch <-chan time.Time) {

		<-ch
		invitation, _ := FetchLeaugeInvitation(m.ID)
		if invitation.Member1 == "" && invitation.Member2 == "" && invitation.Member3 == "" && invitation.Member4 == "" {
			return
		}

		if invitation.Status == LeagueStatusEnd {
			return
		}

		MarkLeagueInvitation(m.ID, LeagueStatusEnd)

		sendNotify(b, m, invitation)
	}(time.After(time.Second * duration))
}

func sendNotify(b *tb.Bot, m *tb.Message, invitation LeagueInvitation) {
	body := fmt.Sprintf(`
<a>当前模式：</a><strong>%s</strong>`, invitation.Rule)
	var count int // 统计组排人数
	if invitation.Member1 != "" {
		count++
		body += fmt.Sprintf(`
<a href="tg://user?id=%d">@%s </a>`, invitation.MemberID1, invitation.Member1)
	}
	if invitation.Member2 != "" {
		count++
		body += fmt.Sprintf(`
<a href="tg://user?id=%d">@%s </a>`, invitation.MemberID2, invitation.Member2)
	}
	if invitation.Member3 != "" {
		count++
		body += fmt.Sprintf(`
<a href="tg://user?id=%d">@%s </a>`, invitation.MemberID3, invitation.Member3)
	}
	if invitation.Member4 != "" {
		count++
		body += fmt.Sprintf(`
<a href="tg://user?id=%d">@%s </a>`, invitation.MemberID4, invitation.Member4)
	}
	if count == 1 {
		body += `
新的一年祝您身体健康，万事如意！
`
	} else {
		body += `
组排啦！鸽子？不存在的!
`
	}
	b.Send(m.Chat, body, tb.ModeHTML)
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
	loc, _ := time.LoadLocation("Asia/Shanghai")
	startTime := time.Unix(int64(invitation.StartTime), 0).In(loc).Format("15:04")
	body := fmt.Sprintf(`<a>组排模式：</a><strong>%s</strong>
<a>组排时间：</a><strong>%s</strong>
<a>参与乌贼：</a>`, invitation.Rule, startTime)

	var count int // 统计组排人数
	if invitation.Member1 != "" {
		count++
		body += fmt.Sprintf(`
<a href="tg://user?id=%d">%s</a>`, invitation.MemberID1, invitation.Member1)
	}
	if invitation.Member2 != "" {
		count++
		body += fmt.Sprintf(`
<a href="tg://user?id=%d">%s</a>`, invitation.MemberID2, invitation.Member2)
	}
	if invitation.Member3 != "" {
		count++
		body += fmt.Sprintf(`
<a href="tg://user?id=%d">%s</a>`, invitation.MemberID3, invitation.Member3)
	}
	if invitation.Member4 != "" {
		count++
		body += fmt.Sprintf(`
<a href="tg://user?id=%d">%s</a>`, invitation.MemberID4, invitation.Member4)
	}

	// 满员
	if (invitation.Type == LeagueTypeDouble && count == 2) || (invitation.Type == LeagueTypeFour && count == 4) {
		body += `
<strong>车已满员,等待发车!</strong>`
		// 满员且需要立即发车
		if int64(invitation.StartTime) <= time.Now().Unix() {
			b.Edit(&msg, body, tb.ModeHTML)
			sendNotify(b, &msg, invitation)
			return
		}
	}

	b.Edit(&msg, body, tb.ModeHTML, markup1)
}
