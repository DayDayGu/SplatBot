// Package main provides ...
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"image"
	"image/draw"
	"image/png"

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
	for _, weapon := range detail.Weapons {
		weaponURL, weaponName := weaponConfig(weapon)
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			if !faker.Exist(weaponURL) {
				faker.Download(weaponURL, weaponName)
				wg.Done()
			}
		}(&wg)
	}
	url := "https://splatoon2.ink/assets/splatnet" + detail.Stage.Image
	name := detail.Stage.Name
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

func weaponConfig(weapon WeaponElement) (url string, name string) {
	var weaponURL string
	var weaponName string
	if weapon.Weapon != nil {
		weaponURL = "https://splatoon2.ink/assets/splatnet" + weapon.Weapon.Image
		weaponName = weapon.Weapon.Name
	} else if weapon.CoopSpecialWeapon != nil {
		weaponURL = "https://splatoon2.ink/assets/splatnet" + weapon.CoopSpecialWeapon.Image
		weaponName = weapon.CoopSpecialWeapon.Name
	}
	return weaponURL, weaponName
}

// Combine ..
func Combine(detail Detail) {
	bg := faker.Get(detail.Stage.Name)
	if bg == nil {
		return
	}
	bounds := bg.Bounds()
	ret := image.NewRGBA(bounds)
	draw.Draw(ret, bounds, bg, image.ZP, draw.Src)
	for idx, weapon := range detail.Weapons {
		_, weaponName := weaponConfig(weapon)
		if faker.Exist(weaponName) {
			weaponImage := faker.Get(weaponName)
			weaponBounds := weaponImage.Bounds().Min.Add(image.Point{-100 * idx, 0})
			draw.Draw(ret, bounds, weaponImage, weaponBounds, draw.Src)
		}
	}
	path, _ := os.Getwd()
	fp := fmt.Sprintf(`%s/tmp/%d`, path, detail.StartTime)
	f, err := os.Create(fp)
	if err != nil {
		fmt.Println(err)
		fmt.Println(fp)
	}
	defer f.Close()
	png.Encode(f, ret)

}
