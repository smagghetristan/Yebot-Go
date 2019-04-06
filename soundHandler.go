package main

import (
	"io"
	"strings"
	"time"

	"github.com/smagghetristan/dca"

	"github.com/andersfylling/disgord"
	"github.com/rylio/ytdl"
)

type Player struct {
	id     disgord.Snowflake
	stream *dca.StreamingSession
}

type Queue struct {
	URL string
}

type Queues struct {
	GuildID disgord.Snowflake
	Queue   []Queue
}

type Parameters struct {
	InVoice bool
	vc      disgord.VoiceConnection
}

var AllQueues []Queues
var players []Player

func RemoveFromArray(s []Player, GuildID disgord.Snowflake) []Player {
	i := 0
	exist := false
	for i = 0; i < len(players); i++ {
		if players[i].id == GuildID {
			exist = true
			break
		}
	}
	if exist {
		s[i] = s[len(s)-1]
		return s[:len(s)-1]
	} else {
		return s
	}
}

func RemoveFromQueue(s []Queues, GuildID disgord.Snowflake, URL string) []Queue {
	i := 0
	k := 0
	exist := false
	for i = 0; i < len(AllQueues); i++ {
		if AllQueues[i].GuildID == GuildID {
			for k = 0; k < len(AllQueues[i].Queue); k++ {
				if AllQueues[i].Queue[k].URL == URL {
					exist = true
					break
				}
			}
			break
		}
	}
	if exist {
		s[i].Queue[k] = s[i].Queue[len(s[i].Queue)-1]
		return s[i].Queue[:len(s[i].Queue)-1]
	} else {
		return s[i].Queue
	}
}

func RemoveQueue(s []Queues, GuildID disgord.Snowflake) []Queues {
	i := 0
	exist := false
	for i = 0; i < len(AllQueues); i++ {
		if AllQueues[i].GuildID == GuildID {
			exist = true
			break
		}
	}
	if exist {
		s[i] = s[len(s)-1]
		return s[:len(s)-1]
	} else {
		return s
	}
}

func Playing(s disgord.Session, GuildID disgord.Snowflake, MessageChannelID disgord.Snowflake, ChannelID disgord.Snowflake, videoID string, InVoice Parameters) {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"

	videoInfo, err := ytdl.GetVideoInfo(`https://www.youtube.com/watch?v=` + videoID)
	if err != nil {
		// Handle the error
	}

	s.SendMsg(MessageChannelID, ":arrow_down: Downloading : "+videoInfo.Title)

	format := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
	downloadURL, err := videoInfo.GetDownloadURL(format)
	if err != nil {
		// Handle the error
	}

	encodingSession, err := dca.EncodeFile(downloadURL.String(), options)
	if err != nil {
		// Handle the error
	}
	defer encodingSession.Cleanup()

	time.Sleep(250 * time.Millisecond)

	s.SendMsg(MessageChannelID, ":musical_note: Started playing : "+videoInfo.Title)
	var vc disgord.VoiceConnection
	if InVoice.InVoice {
		vc = InVoice.vc
	} else {
		vc, err = s.VoiceConnect(GuildID, ChannelID)
		if err != nil {
			return
		}
	}
	_ = vc.StartSpeaking()
	done := make(chan error)
	stream := dca.NewStream(encodingSession, vc, done)
	stream.SetPaused(false)
	CurrentPlayer := Player{
		id:     GuildID,
		stream: stream,
	}
	players = append(players, CurrentPlayer)
	err = <-done
	if err != nil && err != io.EOF {
	}
	time.Sleep(250 * time.Millisecond)
	players = RemoveFromArray(players, GuildID)
	QueueCheck(s, GuildID, MessageChannelID, ChannelID, videoID, vc)
	return
}

func QueueCheck(s disgord.Session, GuildID disgord.Snowflake, MessageChannelID disgord.Snowflake, ChannelID disgord.Snowflake, videoID string, vc disgord.VoiceConnection) {
	i := 0
	k := 0
	exist := false
	for i = 0; i < len(AllQueues); i++ {
		if AllQueues[i].GuildID == GuildID {
			for k = 0; k < len(AllQueues[i].Queue); k++ {
				if AllQueues[i].Queue[k].URL == videoID {
					exist = true
					break
				}
			}
			break
		}
	}
	if exist {
		AllQueues[i].Queue[k] = AllQueues[i].Queue[len(AllQueues[i].Queue)-1]
		AllQueues[i].Queue = AllQueues[i].Queue[:len(AllQueues[i].Queue)-1]
		if len(AllQueues[i].Queue) != 0 {
			Params := Parameters{
				InVoice: true,
				vc:      vc,
			}
			go Playing(s, GuildID, MessageChannelID, ChannelID, AllQueues[i].Queue[0].URL, Params)
		}
	} else {
		_ = vc.Close()
	}
}

func QueueManagement(s disgord.Session, GuildID disgord.Snowflake, MessageChannelID disgord.Snowflake, ChannelID disgord.Snowflake, videoID string) {
	i := 0
	exist := false
	for i = 0; i < len(AllQueues); i++ {
		if AllQueues[i].GuildID == GuildID {
			exist = true
			break
		}
	}
	if exist {
		Song := Queue{
			URL: videoID,
		}
		AllQueues[i].Queue = append(AllQueues[i].Queue, Song)
	} else {
		Song := Queue{
			URL: videoID,
		}
		var songs []Queue
		songs = append(songs, Song)
		Queue := Queues{
			GuildID: GuildID,
			Queue:   songs,
		}
		AllQueues = append(AllQueues, Queue)
	}
	Params := Parameters{
		InVoice: false,
	}
	Playing(s, GuildID, MessageChannelID, ChannelID, videoID, Params)
}

func AddToQueue(Session disgord.Session, GuildID disgord.Snowflake, MessageChannelID disgord.Snowflake, ChannelID disgord.Snowflake, URL string) {
	NotId := strings.Split(URL, "&")
	NotId = strings.Split(NotId[0], "=")
	videoID := NotId[1]
	i := 0
	exist := false
	for i = 0; i < len(AllQueues); i++ {
		if AllQueues[i].GuildID == GuildID {
			exist = true
			break
		}
	}
	if exist {
		Song := Queue{
			URL: videoID,
		}
		AllQueues[i].Queue = append(AllQueues[i].Queue, Song)
		videoInfo, err := ytdl.GetVideoInfo(`https://www.youtube.com/watch?v=` + videoID)
		if err != nil {
			// Handle the error
		}

		Session.SendMsg(MessageChannelID, ":arrow_down: Adding to queue : "+videoInfo.Title)
	} else {
		Song := Queue{
			URL: videoID,
		}
		var songs []Queue
		songs = append(songs, Song)
		Queue := Queues{
			GuildID: GuildID,
			Queue:   songs,
		}
		AllQueues = append(AllQueues, Queue)
		Params := Parameters{
			InVoice: false,
		}
		go Playing(Session, GuildID, MessageChannelID, ChannelID, videoID, Params)
	}
}

func PlayYoutubeLink(Session disgord.Session, GuildID disgord.Snowflake, MessageChannelID disgord.Snowflake, ChannelID disgord.Snowflake, URL string) {
	i := 0
	for i = 0; i < len(players); i++ {
		if players[i].id == GuildID {
			AddToQueue(Session, GuildID, MessageChannelID, ChannelID, URL)
			return
		}
	}
	NotId := strings.Split(URL, "&")
	NotId = strings.Split(NotId[0], "=")
	videoID := NotId[1]
	go QueueManagement(Session, GuildID, MessageChannelID, ChannelID, videoID)
}

func StopPlaying(GuildID disgord.Snowflake, ChannelID disgord.Snowflake, Session disgord.Session) {
	i := 0
	exist := false
	for i = 0; i < len(players); i++ {
		if players[i].id == GuildID {
			exist = true
			break
		}
	}
	if exist {
		Session.SendMsg(ChannelID, ":no_entry: Stopped the current song.")
		players[i].stream.Stop()
	}
}

func PausePlaying(GuildID disgord.Snowflake, ChannelID disgord.Snowflake, Session disgord.Session) {
	i := 0
	exist := false
	for i = 0; i < len(players); i++ {
		if players[i].id == GuildID {
			exist = true
			break
		}
	}
	if exist {
		Session.SendMsg(ChannelID, ":pause_button: Paused the current song.")
		players[i].stream.SetPaused(true)
	}
}

func ResumePlaying(GuildID disgord.Snowflake, ChannelID disgord.Snowflake, Session disgord.Session) {
	i := 0
	exist := false
	for i = 0; i < len(players); i++ {
		if players[i].id == GuildID {
			exist = true
			break
		}
	}
	if exist {
		Session.SendMsg(ChannelID, ":musical_note: Resuming the current song.")
		players[i].stream.SetPaused(false)
	}
}
