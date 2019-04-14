// Package main provides ...
package main

type Schedules struct {
	Regular []Battle `json:"regular"`
	Gachi   []Battle `json:"gachi"`
	League  []Battle `json:"league"`
}

type GameMode struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type Stage struct {
	ID    string `json:"id"`
	Image string `json:"image"`
	Name  string `json:"name"`
}

type Rule struct {
	Key           string `json:"key"`
	Name          string `json:"name"`
	MultiLineName string `json:"multiline_name"`
}

type Battle struct {
	ID        int64    `json:"id"`
	StageA    Stage    `json:"stage_a"`
	StageB    Stage    `json:"stage_b"`
	EndTime   int64    `json:"end_time"`
	StartTime int64    `json:"start_time"`
	GameMode  GameMode `json:"game_mode"`
	Rule      Rule     `json:"rule"`
}
