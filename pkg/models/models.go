package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type Position struct {
	Id_position              int
	Position_name            string
	Profil                   string
	Minimum_available_salary int
	Maximum_available_salary int
	Required_citizenship     string
	Required_education       string
	Required_hours           int
	Required_experience      string
	Required_to_relocate     bool
	Required_work_format     string
}
type Possible_candidate struct {
	Id_possible_candidate int
	Candidate_name        string
	Telegram_username     string
	Contact_number        string
	Citizenship           string
	Education             string
	Work_experience       string
	Hours                 int
	Work_format           string
	Expected_salary       int
	Ready_to_relocate     bool
	Date_of_dialog        time.Time
	Feadback              string
	Is_blocked_flag       bool
	Ready_flag            bool
	Fail_flag             bool
	Id_pos                int
}
