// Package main provides ...
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// LeagueStatus 组排状态
type LeagueStatus int

// LeagueStatusActive 活跃
const (
	LeagueStatusActive = 0
	LeagueStatusFull   = 1
	LeagueStatusEnd    = 2
)

// LeagueType 双排or四排
type LeagueType int

// LLeagueTypeDouble 双排
const (
	LeagueTypeDouble = 0
	LeagueTypeFour   = 1
)

// LeagueInvitation 邀请struct
type LeagueInvitation struct {
	ID         int          `json:"id"`
	Member1    string       `json:"member1"`
	Member2    string       `json:"member2"`
	Member3    string       `json:"member3"`
	Member4    string       `json:"member4"`
	MemberID1  int          `json:"memberid1"`
	MemberID2  int          `json:"memberid2"`
	MemberID3  int          `json:"memberid3"`
	MemberID4  int          `json:"memberid4"`
	StartTime  int          `json:"start_time"` // 开始时间
	Rule       string       `json:"rule"`       // 真格规则
	Type       LeagueType   `json:"type"`       // 双排or四排
	Status     LeagueStatus `json:"status"`     // 组队状态
	CreateDate string       `json:"create_date"`
}

var mu sync.RWMutex

// CreateLeagueInvitation 创建邀请
// id msgid
// ownerID 用户的id
// owner 用户的username
// rule 真格规则
// startTime 真格开始时间戳
func CreateLeagueInvitation(id, ownerID int, owner, rule string, startTime int64, tp LeagueType) error {
	if id == 0 || owner == "" {
		return errors.New("无效的用户")
	}

	now := time.Now().Unix()
	insert := fmt.Sprintf(`insert into league
                           (id, 
                           memberid1,
                           member1,
                           rule,
                           start_time,
                           type,
                           create_date) 
                           values('%d','%d', '%s', '%s', '%d', '%d', '%d')`,
		id, ownerID, owner, rule, startTime, tp, now)
	mu.Lock()
	_, err := DefaultDB.Exec(insert)
	mu.Unlock()
	if err != nil {
		fmt.Printf("create err:%s", err)
		return err
	}
	return nil
}

// MarkLeagueInvitation 标记组排状态
func MarkLeagueInvitation(id int, status LeagueStatus) error {
	mark := func(id int, status LeagueStatus) error {
		update := fmt.Sprintf(`update league set status = '%s' where id = %d`,
			status, id)
		mu.Lock()
		_, err := DefaultDB.Exec(update)
		mu.Unlock()
		if err != nil {
			fmt.Printf("mark err:%s", err)
			return err
		}
		return nil
	}
	return mark(id, status)
}

// DeleteLeagueInvitationMember 删除参与组排乌贼
func DeleteLeagueInvitationMember(id, userID int) error {
	invitation, err := FetchLeaugeInvitation(id)
	if err != nil {
		return errors.New("无法找到组队邀请")
	}

	del := func(id int, key, idkey string) error {
		update := fmt.Sprintf(`update league set %s = '',%s = 0 where id = %d`,
			key, idkey, id)
		mu.Lock()
		_, err := DefaultDB.Exec(update)
		mu.Unlock()
		if err != nil {
			fmt.Printf("delete err:%s", err)
			return err
		}
		in, _ := FetchLeaugeInvitation(id)
		fmt.Println(in)
		return nil
	}
	if invitation.MemberID1 == userID {
		return del(id, "member1", "memberid1")
	}
	if invitation.MemberID2 == userID {
		return del(id, "member2", "memberid2")
	}
	if invitation.MemberID3 == userID {
		return del(id, "member3", "memberid3")
	}
	if invitation.MemberID4 == userID {
		return del(id, "member4", "memberid4")
	}
	return errors.New("你这只蠢乌贼也没参与啊")
}

// AddLeagueInvitationMember 增加组排乌贼
func AddLeagueInvitationMember(id, userID int, username string) error {
	invitation, err := FetchLeaugeInvitation(id)
	if err != nil {
		return errors.New("无法找到组队邀请")
	}

	add := func(id int, key, value, idkey string, idvalue int) error {
		update := fmt.Sprintf(`update league set %s = '%s',%s = %d where id = %d`,
			key, value, idkey, idvalue, id)
		mu.Lock()
		_, err := DefaultDB.Exec(update)
		mu.Unlock()
		if err != nil {
			fmt.Printf("add err:%s", err)
			return err
		}
		return nil
	}
	if username == invitation.Member1 || username == invitation.Member2 || username == invitation.Member3 || username == invitation.Member4 {
		return errors.New("已经参与过")
	}
	if invitation.Member1 == "" {
		return add(id, "member1", username, "memberid1", userID)
	}
	if invitation.Member2 == "" {
		return add(id, "member2", username, "memberid2", userID)
	}
	if invitation.Member3 == "" {
		return add(id, "member3", username, "memberid3", userID)
	}
	if invitation.Member4 == "" {
		return add(id, "member4", username, "memberid4", userID)
	}
	return errors.New("队伍已经满员啦")
}

// FetchLeaugeInvitation 获取组队信息
func FetchLeaugeInvitation(id int) (LeagueInvitation, error) {
	ret := LeagueInvitation{}
	if id == 0 {
		return ret, errors.New("无效的用户")
	}

	fetch := fmt.Sprintf(`select * 
                          from league 
                          where id = '%d'`, id)
	mu.RLock()
	rows, err := DefaultDB.Query(fetch)
	mu.RUnlock()
	if err != nil {
		fmt.Printf("fetch err:%s", err)
		return ret, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&ret.ID,
			&ret.Member1,
			&ret.Member2,
			&ret.Member3,
			&ret.Member4,
			&ret.MemberID1,
			&ret.MemberID2,
			&ret.MemberID3,
			&ret.MemberID4,
			&ret.StartTime,
			&ret.Rule,
			&ret.Type,
			&ret.Status,
			&ret.CreateDate,
		)
		if err != nil {
			return ret, err
		}
	}
	return ret, nil
}
