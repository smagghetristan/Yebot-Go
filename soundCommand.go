package main

import (
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/fatih/color"
)

func soundCommand(s disgord.Session, m *disgord.Message) {
	if strings.HasPrefix(m.Content, prefix+"play") {
		// Find the guild for that channel.
		g, err := s.GetGuild(m.GuildID)
		if err != nil {
			// Could not find guild.
			return
		}
		// Look for the message sender in that guild's current voice states.
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				for _, c := range g.Channels {
					if c.ID == vs.ChannelID {
						url := strings.Replace(m.Content, prefix+"play", "", 1)
						go PlayYoutubeLink(s, g.ID, m.ChannelID, vs.ChannelID, url)
						if err != nil {
							color.Red("Error playing sound:", err)
						}
						return
					}
				}
				return
			}
		}
	}
	if strings.HasPrefix(m.Content, prefix+"stop") {
		StopPlaying(m.GuildID, m.ChannelID, s)
	}
	if strings.HasPrefix(m.Content, prefix+"pause") {
		PausePlaying(m.GuildID, m.ChannelID, s)
	}
	if strings.HasPrefix(m.Content, prefix+"resume") {
		ResumePlaying(m.GuildID, m.ChannelID, s)
	}
	return
}
