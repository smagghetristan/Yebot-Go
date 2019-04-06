package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
)

type MemeStruct struct {
	PostLink  string `json:"postLink"`
	Subreddit string `json:"subreddit"`
	Title     string `json:"title"`
	Url       string `json:"url"`
}

func funCommands(s disgord.Session, m *disgord.Message) {
	if strings.HasPrefix(m.Content, prefix+"say ") {
		message := strings.Replace(m.Content, prefix+"say ", "", 1)
		_, err := s.SendMsg(m.ChannelID, message)
		if err != nil {
			return
		}
		err = s.DeleteMessage(m.ChannelID, m.ID)
		if err != nil {
			return
		}
	}

	if strings.HasPrefix(m.Content, prefix+"meme") {
		response, err := http.Get("https://meme-api.herokuapp.com/gimme")
		if err != nil {
			//
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			JSON := MemeStruct{}
			err := json.Unmarshal(data, &JSON)
			if err != nil {
				//
			}
			AvatarURL := GetAvatarURL(m.Author)
			Author := &disgord.EmbedAuthor{
				Name:    m.Author.Username,
				IconURL: AvatarURL,
			}
			Image := &disgord.EmbedImage{
				URL: JSON.Url,
			}
			embed := &disgord.Embed{
				Author: Author,
				Image:  Image,
				Color:  0xFFDD00,
			}
			Params := &disgord.CreateMessageParams{
				Embed: embed,
			}
			_, err = s.CreateMessage(m.ChannelID, Params)
			if err != nil {
				return
			}
			err = s.DeleteMessage(m.ChannelID, m.ID)
			if err != nil {
				return
			}
		}
	}

	if strings.HasPrefix(m.Content, prefix+"success ") {
		message := strings.Replace(m.Content, prefix+"success ", "", 1)
		title := strings.Split(message, ";")[0]
		text := strings.Split(message, ";")[1]
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1).Intn(39)
		url := "http://www.minecraftachievement.net/achievement/a.php?i=" + strconv.Itoa(r1) + "&h=" + url.QueryEscape(title) + "&t=" + url.QueryEscape(text)
		AvatarURL := GetAvatarURL(m.Author)
		Author := &disgord.EmbedAuthor{
			Name:    m.Author.Username,
			IconURL: AvatarURL,
		}
		Image := &disgord.EmbedImage{
			URL: url,
		}
		embed := &disgord.Embed{
			Author: Author,
			Image:  Image,
			Color:  0xFFDD00,
		}
		Params := &disgord.CreateMessageParams{
			Embed: embed,
		}
		_, err := s.CreateMessage(m.ChannelID, Params)
		if err != nil {
			return
		}
		err = s.DeleteMessage(m.ChannelID, m.ID)
		if err != nil {
			return
		}
	}

	return
}
