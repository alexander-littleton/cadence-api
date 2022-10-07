package models

type Habit struct {
	Id            string   `json:"id"`
	Name          string   `json:"name"`
	UserId        string   `json:"user_id"`
	Cadence       Cadence  `json:"cadence"`
	RepeatingDays []uint16 `json:"repeating_days"`
	Streak        uint32   `json:"streak"`
}

type Cadence uint8

const (
	Day Cadence = iota
	Month
)
