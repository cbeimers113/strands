package state

import (
	"fmt"

	"cbeimers113/strands/internal/config"
)

// Represents the local time of the simulation world
type Clock struct {
	*config.Config
	minute int
	hour   int
	day    int

	twelveHr bool
	timer    float32
}

// Create a new clock with the specified time given the in-game day length
func NewClock(cfg *config.Config, hour, minute int, twelveHr bool) *Clock {
	c := &Clock{
		Config:   cfg,
		minute:   minute,
		hour:     hour,
		twelveHr: twelveHr,
	}

	c.SetTime(hour, minute)

	return c
}

// SetTime sets the internal timer to the given hour and minute for a given day length
func (c *Clock) SetTime(hour, minute int) {
	ms := float32(60 * 1000 * (minute + hour*60))                                  // Number of ms between midnight and the selected time
	c.timer = float32(c.Simulation.DayLength*60*1000) * ms / (1000 * 60 * 60 * 24) // Scale the number of ms to the game day length
}

// Progress returns percentage of progress through the day
func (c Clock) Progress(int) float32 {
	return c.timer / float32(c.Simulation.DayLength*60*1000)
}

// Update adds ms of real time into the world's time
func (c *Clock) Update(ms float32) {
	c.timer += ms

	// 1440 minutes in a day, so set the minutes count equal to n% of 1440, where n is the progress through the day
	numMins := int(c.Progress(c.Simulation.DayLength) * 1440)
	c.minute = numMins % 60
	c.hour = int(numMins / 60)

	if c.hour >= 24 {
		c.timer = 0
		c.hour = 0
		c.day++
	}
}

// Get the current hour
func (c Clock) Hour() int {
	return c.hour
}

// Get the current minute
func (c Clock) Minute() int {
	return c.minute
}

// Get a string representation of the current time
func (c Clock) String() string {
	var (
		h    string
		m    string
		amPm string
	)

	h = fmt.Sprint(c.hour)
	if c.hour < 10 && !c.twelveHr {
		h = "0" + h
	} else if c.twelveHr {
		if c.hour > 12 {
			h = fmt.Sprint(c.hour - 12)
		} else if c.hour == 0 {
			h = "12"
		}
	}

	m = fmt.Sprint(c.minute)
	if c.minute < 10 {
		m = "0" + m
	}

	if c.twelveHr {
		amPm = "am"
		if c.hour >= 12 {
			amPm = "pm"
		}
	}

	return fmt.Sprintf("%s:%s %s, Day %d", h, m, amPm, c.day)
}
