package main

import (
	"os"
	"strings"

	"github.com/andersfylling/disgord"
)

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.

func CommandHandle(session disgord.Session, evt *disgord.MessageCreate) {
	m := evt.Message
	current, err := session.GetCurrentUser()
	if err != nil {
		//
	}
	if m.Author.ID == current.ID {
		return
	}
	if strings.HasPrefix(m.Content, prefix+"restart") {
		if m.Author.ID.String() == "144472011924570113" {
			os.Exit(0)
		}
	}
	soundCommand(session, m)
	adminCommands(session, m)
	helpCommands(session, m)
	funCommands(session, m)
}
func messageCreate(session disgord.Session, evt *disgord.MessageCreate) {
	go CommandHandle(session, evt)
	return
}
