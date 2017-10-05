package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
	"sync"
	"github.com/sushi86/DCTS3Bot/telegram"
	"strings"
	"github.com/sushi86/DCTS3Bot/config"
)

type client struct {
	id 			string
	nickname 	string
	game		string
	conTime		time.Time
}

var clients = make(map[string]client)
var lock = sync.RWMutex{}

func addClient(client client) {
	lock.Lock()
	defer lock.Unlock()
	clients[client.id] = client
}

func removeClient(client client) {
	lock.Lock()
	defer lock.Unlock()
	cid := client.id
	if _, ok := clients[cid]; ok {
		//bridge <- event{status: "disconnect", client: clients[cid], time: time.Now()}
		delete(clients, cid)
	}
}

func getClient(id string) client {
	lock.Lock()
	defer lock.Unlock()
	return clients[id]
}


func ConnectDc() {
	config := config.GetConfig()

	dg, err := discordgo.New(config.DiscordApiKey)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(userJoined)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	for { select {} }

}

func userJoined(s *discordgo.Session, m *discordgo.PresenceUpdate) {
	if (m.Status == "online") {
		if m.User.Username == "" {
			return
		}
		gameName := ""
		if (m.Game == nil) {
			gameName = "nix"
		} else {
			gameName = m.Game.Name
		}
		client := client{id: m.User.ID, nickname: m.User.Username, game: gameName, conTime: time.Now()}
		addClient(client)
		telegram.Send(m.User.Username + " connected (dc) playing " + gameName)
	} else if (m.Status == "offline") {
		client := getClient(m.User.ID)
		fmt.Println(client)
		if (client.id != "") {
			duration := time.Since(client.conTime)
			message := client.nickname + " disconnected " + shortDur(duration)

			telegram.Send(message)
			removeClient(client)
		}
	}
}

func shortDur(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}