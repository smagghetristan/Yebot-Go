package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
)

func helpCommands(s disgord.Session, m *disgord.Message) {
	if strings.HasPrefix(m.Content, prefix+"help") {
		//ADMIN FIELD
		AdminField := &disgord.EmbedField{
			Name: "Admin commands:",
			Value: `**kick** : Kicks a user, you have to mention him.
		        **ban** : Bans a user, you have to mention him.
		        **ping** : Will make the bot reply "Pong !"
		        **clear** : Will make the bot clear a channel's message (clear **8**)`,
		}
		//FUN FIELD
		FunField := &disgord.EmbedField{
			Name: "Fun commands:",
			Value: `**say**: The bot will repeat what you said after the command.
						**success** title;text: The bot will create a minecraft-like success.
						**meme**: The bot will send a random meme.`,
		}
		MusicField := &disgord.EmbedField{
			Name: "Music commands:",
			Value: `**play** *youtube-url* : Will play the song from the youtube link.
		        **pause** : Will pause the current song.
		        **resume** : Will resume the current song.`,
		}
		//RANDOM FIELD
		RandomField := &disgord.EmbedField{
			Name:  "Random Commands :",
			Value: `**info** : Will return the bots info.`,
		}

		AllFields := []*disgord.EmbedField{AdminField, MusicField, FunField, RandomField}

		embed := &disgord.Embed{
			Title:       "Help Menu",
			Description: "The prefix of the bot is currently : **g!**",
			Fields:      AllFields,
			Color:       0xFFDD00,
		}

		Params := &disgord.CreateMessageParams{
			Embed: embed,
		}
		s.CreateMessage(m.ChannelID, Params)
	}

	if strings.HasPrefix(m.Content, prefix+"info") {
		ServerAmount := s.GetConnectedGuilds()
		ChannelAmount := 0
		MemberAmount := 0
		for i := 0; i < len(ServerAmount); i++ {
			g, err := s.GetGuild(m.GuildID)
			if err != nil {
				// Could not find guild.
				return
			}
			ChannelAmount += len(g.Channels)
			MemberAmount += len(g.Members)
		}
		Uptime := time.Since(startTime)

		embed := &disgord.Embed{
			Title: "Bot Statistics :",
			Description: `**Servers** : ` + strconv.Itoa(len(ServerAmount)) + `
				**Channels** : ` + strconv.Itoa(ChannelAmount) + `
				**Users** (Cumulated) : ` + strconv.Itoa(MemberAmount) + `
				**Uptime** : ` + strconv.FormatFloat(Uptime.Hours(), 'f', 0, 64) + `:` + strconv.FormatFloat(Uptime.Minutes(), 'f', 0, 64) + `:` + strconv.FormatFloat(Uptime.Seconds(), 'f', 0, 64),
			Color: 0xFFDD00,
		}
		Params := &disgord.CreateMessageParams{
			Embed: embed,
		}
		s.CreateMessage(m.ChannelID, Params)
	}
	return
}
