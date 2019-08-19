// Package main provides ...
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"image"
	"image/color"
	"image/draw"
	"image/png"

	"github.com/PangPangPangPangPang/SplatBot/faker"
	"github.com/nfnt/resize"
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
	bg = resize.Resize(380, 380*180/320, bg, resize.NearestNeighbor)
	if bg == nil {
		return
	}
	width := 400
	height := 400
	ret := image.NewRGBA(image.Rect(0, 0, 400, 400))
	// Set color for each pixel.

	cyan := color.RGBA{251, 89, 4, 0xff}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			ret.Set(x, y, cyan)
		}
	}
	bg1 := faker.Get("bg.png")
	bg1 = resize.Resize(400, 400, bg1, resize.NearestNeighbor)
	draw.Draw(ret, ret.Bounds(), bg1, image.ZP, draw.Over)
	// draw.Draw(ret, image.Rect(0, 0, 400, 400), backgroundColor, image.ZP, draw.Src)
	draw.Draw(ret, ret.Bounds().Add(image.Point{10, 90}), bg, image.ZP, draw.Over)
	for idx, weapon := range detail.Weapons {
		_, weaponName := weaponConfig(weapon)
		if faker.Exist(weaponName) {
			weaponImage := faker.Get(weaponName)
			weaponImage = resize.Resize(80, 80, weaponImage, resize.NearestNeighbor)
			// weaponPoints := weaponImage.Bounds().Min.Add(image.Point{100, 100 * idx})
			point := image.Pt(100*idx+10, 0)
			draw.Draw(ret, ret.Bounds().Add(point), weaponImage, weaponImage.Bounds().Min, draw.Over)
		}
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	addLabel(ret, 10, 350, fmt.Sprintf(`From: %s `, time.Unix(detail.StartTime, 0).In(loc).Format("Jan 2 15:04")))
	addLabel(ret, 10, 370, fmt.Sprintf(`To: %s`, time.Unix(detail.EndTime, 0).In(loc).Format("Jan 2 15:04")))
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
func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{0, 0, 0, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
