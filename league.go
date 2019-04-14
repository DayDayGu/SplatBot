// Package main provides ...
package main

import (
	"errors"
	"fmt"
	"time"
)

type LeagueStatus int

const (
	LeagueStatusActive = 0
	LeagueStatusFull
	LeagueStatusEnd
)

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
	Rule       string       `json:"rule"`   // 真格规则
	Type       LeagueStatus `json:"type"`   // 双排or四排
	Status     string       `json:"status"` // 组队状态
	CreateDate string       `json:"create_date"`
}

func CreateLeagueInvitation(id, ownerID int, owner, rule string) error {
	if id == 0 || owner == "" {
		return errors.New("无效的用户")
	}

	now := time.Now().Unix()
	insert := fmt.Sprintf(`insert into league
                           (id, 
                           memberid1,
                           member1,
                           rule,
                           create_date) 
                           values('%d','%d', '%s', '%s', '%d')`,
		id, ownerID, owner, rule, now)
	_, err := DefaultDB.Exec(insert)
	if err != nil {
		return err
	}
	return nil
}

func AddLeagueInvitationMember(id, userID int, username string) error {
	invitation, err := FetchLeaugeInvitation(id)
	if err != nil {
		return errors.New("无法找到组队邀请")
	}

	add := func(id int, key, value, idkey string, idvalue int) error {
		update := fmt.Sprintf(`update league set %s = '%s',%s = %d where id = %d`,
			key, value, idkey, idvalue, id)
		_, err := DefaultDB.Exec(update)
		if err != nil {
			return err
		}
		return nil
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

func FetchLeaugeInvitation(id int) (LeagueInvitation, error) {
	ret := LeagueInvitation{}
	if id == 0 {
		return ret, errors.New("无效的用户")
	}

	fetch := fmt.Sprintf(`select * 
                          from league 
                          where id = '%d'`, id)
	rows, err := DefaultDB.Query(fetch)
	if err != nil {
		return ret, err
	}
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
