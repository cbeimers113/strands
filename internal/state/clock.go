package state

import (
	"fmt"

	"cbeimers113/strands/internal/config"
)

// Represents the local time of the simulation world
type Clock struct {
	*config.Config `json:"-"`

	Minute int `json:"-"`
	Hour   int `json:"-"`
	Day    int `json:"day"`

	TwelveHr bool    `json:"twelve_hr"`
	Timer    float32 `json:"timer"`
}

// Create a new clock with the specified time given the in-game day length
func NewClock(cfg *config.Config, hour, minute int, twelveHr bool) *Clock {
	c := &Clock{
		Config:   cfg,
		Minute:   minute,
		Hour:     hour,
		TwelveHr: twelveHr,
	}

	c.SetTime(hour, minute)

	return c
}

// SetTime sets the internal timer to the given hour and minute for a given day length
func (c *Clock) SetTime(hour, minute int) {
	ms := float32(60 * 1000 * (minute + hour*60))                                  // Number of ms between midnight and the selected time
	c.Timer = float32(c.Simulation.DayLength*60*1000) * ms / (1000 * 60 * 60 * 24) // Scale the number of ms to the game day length
}

// Progress returns percentage of progress through the day
func (c Clock) Progress(int) float32 {
	return c.Timer / float32(c.Simulation.DayLength*60*1000)
}

// Update adds ms of real time into the world's time
func (c *Clock) Update(ms float32) {
	c.Timer += ms

	// 1440 minutes in a day, so set the minutes count equal to n% of 1440, where n is the progress through the day
	numMins := int(c.Progress(c.Simulation.DayLength) * 1440)
	c.Minute = numMins % 60
	c.Hour = int(numMins / 60)

	if c.Hour >= 24 {
		c.Timer = 0
		c.Hour = 0
		c.Day++
	}
}

// Get a string representation of the current time
func (c Clock) String() string {
	var (
		h    string
		m    string
		amPm string
	)

	h = fmt.Sprint(c.Hour)
	if c.Hour < 10 && !c.TwelveHr {
		h = "0" + h
	} else if c.TwelveHr {
		if c.Hour > 12 {
			h = fmt.Sprint(c.Hour - 12)
		} else if c.Hour == 0 {
			h = "12"
		}
	}

	m = fmt.Sprint(c.Minute)
	if c.Minute < 10 {
		m = "0" + m
	}

	if c.TwelveHr {
		amPm = "am"
		if c.Hour >= 12 {
			amPm = "pm"
		}
	}

	return fmt.Sprintf("%s:%s %s, Day %d", h, m, amPm, c.Day)
}
