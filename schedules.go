// Package main provides ...
package main

// Schedules regular/gachi/league
type Schedules struct {
	Regular []Battle `json:"regular"`
	Gachi   []Battle `json:"gachi"`
	League  []Battle `json:"league"`
}

// GameMode 4 modes
type GameMode struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// Stage stage_a/stage_b
type Stage struct {
	ID    string `json:"id"`
	Image string `json:"image"`
	Name  string `json:"name"`
}

// Rule rule
type Rule struct {
	Key           string `json:"key"`
	Name          string `json:"name"`
	MultiLineName string `json:"multiline_name"`
}

//Battle battle msg
type Battle struct {
	ID        int64    `json:"id"`
	StageA    Stage    `json:"stage_a"`
	StageB    Stage    `json:"stage_b"`
	EndTime   int64    `json:"end_time"`
	StartTime int64    `json:"start_time"`
	GameMode  GameMode `json:"game_mode"`
	Rule      Rule     `json:"rule"`
}

// type Work struct {
// ID        int64    `json:"id"`
// StageA    Stage    `json:"stage_a"`
// StageB    Stage    `json:"stage_b"`
// EndTime   int64    `json:"end_time"`
// StartTime int64    `json:"start_time"`
// GameMode  GameMode `json:"game_mode"`
// Rule      Rule     `json:"rule"`
// }
