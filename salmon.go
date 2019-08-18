// Package main provides ...
package main

import (
	"encoding/json"
	"sync"

	"github.com/PangPangPangPangPang/SplatBot/faker"
)

// UnmarshalSalmon ..
func UnmarshalSalmon(data []byte) (Salmon, error) {
	var r Salmon
	err := json.Unmarshal(data, &r)
	return r, err
}

// Marshal ..
func (r *Salmon) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Salmon ..
type Salmon struct {
	Schedules []Schedule `json:"schedules"`
	Details   []Detail   `json:"details"`
}

// Detail ..
type Detail struct {
	Stage     SalmonStage     `json:"stage"`
	StartTime int64           `json:"start_time"`
	EndTime   int64           `json:"end_time"`
	Weapons   []WeaponElement `json:"weapons"`
}

// SalmonStage ..
type SalmonStage struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

// WeaponElement ..
type WeaponElement struct {
	CoopSpecialWeapon *SalmonStage  `json:"coop_special_weapon,omitempty"`
	ID                string        `json:"id"`
	Weapon            *WeaponWeapon `json:"weapon,omitempty"`
}

// WeaponWeapon ..
type WeaponWeapon struct {
	Image     string `json:"image"`
	Thumbnail string `json:"thumbnail"`
	Name      string `json:"name"`
	ID        string `json:"id"`
}

// Schedule ..
type Schedule struct {
	EndTime   int64 `json:"end_time"`
	StartTime int64 `json:"start_time"`
}

// DownloadSalmon ..
func DownloadSalmon(detail Detail) {
	if len(Sa.Details) < 1 {
		return
	}
	var wg sync.WaitGroup
	var url string
	var name string
	for _, weapon := range detail.Weapons {
		if weapon.Weapon != nil {
			url = "https://splatoon2.ink/assets/splatnet" + weapon.Weapon.Image
			name = weapon.Weapon.Name
		} else if weapon.CoopSpecialWeapon != nil {
			url = "https://splatoon2.ink/assets/splatnet" + weapon.CoopSpecialWeapon.Image
			name = weapon.CoopSpecialWeapon.Name
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			if !faker.Exist(url) {
				faker.Download(url, name)
				wg.Done()
			}
		}(&wg)
	}
	url = "https://splatoon2.ink/assets/splatnet" + detail.Stage.Image
	name = detail.Stage.Name
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		if !faker.Exist(url) {
			faker.Download(url, name)
			wg.Done()
		}

	}(&wg)
	wg.Wait()

	Combine(detail)
}

// Combine ..
func Combine(detail Detail) {

}
