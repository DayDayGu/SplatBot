// Package main provides ...
package main

import "encoding/json"

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
