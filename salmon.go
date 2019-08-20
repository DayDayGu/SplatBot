// Package main provides ...
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"golang.org/x/image/font"

	"image"
	"image/color"
	"image/draw"
	"image/png"

	"github.com/PangPangPangPangPang/SplatBot/faker"
	"github.com/golang/freetype"
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

type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
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
	width := 660
	height := 660
	ret := image.NewRGBA(image.Rect(0, 0, width, height))

	// Set color for each pixel.
	bgColor := color.RGBA{251, 89, 4, 0xff}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			ret.Set(x, y, bgColor)
		}
	}

	bgImg := faker.GetResrc("bg.png")
	bgImg = resize.Resize(uint(width), uint(height), bgImg, resize.NearestNeighbor)
	draw.Draw(ret, ret.Bounds(), bgImg, image.ZP, draw.Over)

	stage := faker.Get(detail.Stage.Name)
	// stage = resize.Resize(640, 360, stage, resize.NearestNeighbor)
	stage = circleMask(stage, 20)
	if stage == nil {
		return
	}
	draw.Draw(ret, ret.Bounds().Add(image.Point{10, 150}), stage, image.ZP, draw.Over)

	for idx, weapon := range detail.Weapons {
		_, weaponName := weaponConfig(weapon)
		if faker.Exist(weaponName) {
			weaponImage := faker.Get(weaponName)
			weaponImage = resize.Resize(150, 150, weaponImage, resize.NearestNeighbor)
			point := image.Pt(150*idx+30, 0)
			draw.Draw(ret, ret.Bounds().Add(point), weaponImage, weaponImage.Bounds().Min, draw.Over)
		}
	}

	loc, _ := time.LoadLocation("Asia/Shanghai")
	fontRender(ret, 10, height-90, fmt.Sprintf(`From: %s `, time.Unix(detail.StartTime, 0).In(loc).Format("Jan 2 15:04 Mon")))
	fontRender(ret, 10, height-30, fmt.Sprintf(`To: %s`, time.Unix(detail.EndTime, 0).In(loc).Format("Jan 2 15:04 Mon")))

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

func circleMask(img image.Image, d int) image.Image {

	ret := image.NewRGBA(img.Bounds())
	X := img.Bounds().Max.X
	Y := img.Bounds().Max.Y
	for x := 0; x < X; x++ {
		for y := 0; y < Y; y++ {
			if x < d && y < d && distance(x, y, d, d) > float64(d) {
				ret.Set(x, y, color.RGBA{255, 255, 255, 0})
			} else if X-x < d && y < d && distance(X-x, y, d, d) > float64(d) {
				ret.Set(x, y, color.RGBA{255, 255, 255, 0})
			} else if x < d && Y-y < d && distance(x, Y-y, d, d) > float64(d) {
				ret.Set(x, y, color.RGBA{255, 255, 255, 0})
			} else if X-x < d && Y-y < d && distance(X-x, Y-y, d, d) > float64(d) {
				ret.Set(x, y, color.RGBA{255, 255, 255, 0})
			} else {
				ret.Set(x, y, img.At(x, y))

			}

		}
	}

	return ret
}

func distance(x1, y1, x2, y2 int) float64 {
	first := math.Pow(float64(x2-x1), 2)
	second := math.Pow(float64(y2-y1), 2)
	return math.Sqrt(first + second)
}

var (
	dpi = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	// fontfile = flag.String("fontfile", "/Users/zhiliao/Downloads/ffffonts/simsun.ttf", "filename of the ttf font")
	fontfile = flag.String("fontfile", "./resource/BenMoZhuHei-2.ttf", "BenMoZhuHei")
	hinting  = flag.String("hinting", "none", "none | full")
	size     = flag.Float64("size", 40, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
)

func fontRender(jpg *image.RGBA, x int, y int, text string) {
	flag.Parse()
	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetClip(jpg.Bounds())
	c.SetDst(jpg)
	fg := image.Opaque
	c.SetSrc(fg)

	switch *hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	pt := freetype.Pt(x, y)
	_, err = c.DrawString(text, pt)
	if err != nil {
		log.Println(err)
		return
	}
	pt.Y += c.PointToFixed(*size * *spacing)

}
