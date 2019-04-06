package main

import (
	"strconv"
	"strings"

	"github.com/andersfylling/disgord"
)

func adminCommands(s disgord.Session, m *disgord.Message) {
	member, err := s.GetMember(m.GuildID, m.Author.ID)
	if err != nil {
		return
	}

	if strings.HasPrefix(m.Content, prefix+"ping") {
		s.SendMsg(m.ChannelID, "Pong!")
	}

	if strings.HasPrefix(m.Content, prefix+"ban") {
		isHe := hasPermission(member, s, m.GuildID, PERM_ADMINISTRATOR)
		if isHe {
			if len(m.Mentions) != 0 {
				Params := &disgord.BanMemberParams{
					DeleteMessageDays: 1,
				}
				s.BanMember(m.GuildID, m.Mentions[0].ID, Params)
			}
		} else {
			s.SendMsg(m.ChannelID, "You're not an administrator !")
		}
	}

	if strings.HasPrefix(m.Content, prefix+"kick") {
		isHe := hasPermission(member, s, m.GuildID, PERM_ADMINISTRATOR)
		if isHe {
			if len(m.Mentions) != 0 {
				s.KickMember(m.GuildID, m.Mentions[0].ID)
			}
		} else {
			s.SendMsg(m.ChannelID, "You're not an administrator !")
		}
	}

	if strings.HasPrefix(m.Content, prefix+"clear") {
		isHe := hasPermission(member, s, m.GuildID, PERM_ADMINISTRATOR)
		if isHe {
			CommandSplit := strings.Split(m.Content, " ")
			number, err := strconv.Atoi(CommandSplit[1])
			if err != nil {
				s.SendMsg(m.ChannelID, "I had some problem reading the number you entered.")
			} else if number < 100 {
				Params := &disgord.GetMessagesParams{
					Limit: number,
				}
				MessageSlice, err := s.GetMessages(m.ChannelID, Params)
				if err != nil {
					s.SendMsg(m.ChannelID, "I had some problem deleting the messages.")
				} else {
					Params := &disgord.DeleteMessagesParams{}
					for i := 0; i < len(MessageSlice); i++ {
						Params.AddMessage(MessageSlice[i])
					}
					s.DeleteMessages(m.ChannelID, Params)
				}
			} else {
				s.SendMsg(m.ChannelID, "I cannot delete more than 100 messages.")
			}
		} else {
			s.SendMsg(m.ChannelID, "You're not an administrator !")
		}
	}
	return
}
