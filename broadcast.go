package main

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/raggaer/tiger/app/controllers"
	"github.com/raggaer/tiger/app/models"
)

func monitorServerPlayerDeaths(guild *discordgo.Guild, tick time.Duration, ctx *controllers.Context, s *discordgo.Session, event *discordgo.GuildCreate) {
	// Create event ticker
	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	// Retrieve valid death channels
	deathChannels := retrieveValidChannels(guild, ctx.Config.Discord.DeathChannels)

	// Wait for ticker channel
	for t := range ticker.C {
		deaths, err := models.GetTimeServerDeaths(ctx.DB, 10, t, t.Add(-tick))
		if err != nil {
			log.Printf("Unable to retrieve latest server deaths (death broadcast): %v \r\n", err)
			continue
		}

		// Skip when there are no deaths
		if len(deaths) <= 0 {
			continue
		}

		msg, err := ctx.ExecuteTemplate("broadcast_death", map[string]interface{}{
			"deaths": deaths,
		})
		if err != nil {
			log.Printf("Unable to execute broadcast death template: %v \r\n", err)
			continue
		}

		// Create discord message
		for _, ch := range deathChannels {
			s.ChannelMessageSendEmbed(ch, msg)
		}
	}
}

func retrieveValidChannels(guild *discordgo.Guild, channelNames []string) []string {
	validChannels := []string{}
	for _, ch := range guild.Channels {
		for _, dh := range channelNames {
			if dh == ch.Name {
				validChannels = append(validChannels, ch.ID)
			}
		}
	}
	return validChannels
}
