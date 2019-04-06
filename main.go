package main

import (
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/fatih/color"
)

var buffer = make([][]byte, 0)
var Token = "TOKEN"
var prefix = "g!"
var startTime = time.Now()

func GetAvatarURL(User *disgord.User) string {
	var AvatarURL string
	if strings.HasPrefix(*User.Avatar, "a_") {
		AvatarURL = "https://cdn.discordapp.com/avatars/" + User.ID.String() + "/" + *User.Avatar + ".gif"
	} else {
		AvatarURL = "https://cdn.discordapp.com/avatars/" + User.ID.String() + "/" + *User.Avatar + ".png"
	}
	return AvatarURL
}

func main() {
	// Create a new Discord session using the provided bot token.
	client := disgord.New(&disgord.Config{
		BotToken: Token, // debug=false
	})
	defer client.StayConnectedUntilInterrupted()

	// Register the messageCreate func as a callback for MessageCreate events.
	client.On(disgord.EvtMessageCreate, messageCreate)

	client.Ready(func() {
		GuildIDs := client.GetConnectedGuilds()
		for i := 0; i < len(GuildIDs); i++ {
			p := &disgord.RequestGuildMembersCommand{
				GuildID: GuildIDs[i],
			}
			client.Emit(disgord.CommandRequestGuildMembers, p)
		}

		botReady(client, nil)
	})
}

func botReady(session disgord.Session, evt *disgord.Ready) {
	color.Green("Bot is now running.  Press CTRL-C to exit.")
	session.UpdateStatusString("g!help for help")
}
