package controllers

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	// ApplicationVersion defines the running version
	ApplicationVersion string

	// BuildDate date where the application was built
	BuildDate string
)

// Version returns the current running version
func Version(context *Context, s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.MessageEmbed, error) {
	return context.ExecuteTemplate("version", map[string]interface{}{
		"version":   ApplicationVersion,
		"buildDate": BuildDate,
	})
}

// Uptime returns the current bot uptime
func Uptime(context *Context, s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.MessageEmbed, error) {
	return context.ExecuteTemplate("uptime", map[string]interface{}{
		"start":   context.Start,
		"current": time.Now(),
	})
}

// About returns information about tiger
func About(context *Context, s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.MessageEmbed, error) {
	return context.ExecuteTemplate("about", nil)
}

func uptimeMessage(day, hour, min, sec int) string {
	msg := "%d day"
	if day == 0 || day > 1 {
		msg += "s"
	}
	msg += ", %d hour"
	if hour == 0 || hour > 1 {
		msg += "s"
	}
	msg += ", %d minute"
	if min == 0 || min > 1 {
		msg += "s"
	}
	msg += ", %d second"
	if sec == 0 || sec > 1 {
		msg += "s"
	}
	return msg
}

func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}
