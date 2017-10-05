package ts3

import (
	"github.com/toqueteos/ts3"
	"log"
	"fmt"
	"time"
	"strings"
	"sync"

	"github.com/sushi86/DCTS3Bot/telegram"
	"github.com/sushi86/DCTS3Bot/config"
)

type client struct {
	uid 		string
	cid 		string
	nickname 	string
	conTime		time.Time
}

type event struct {
	status		string
	client		client
	time 		time.Time
}

var conn *ts3.Conn

var bridge = make(chan event)

var clients = make(map[string]client)
var lock = sync.RWMutex{}

func addClient(client client) {
	lock.Lock()
	defer lock.Unlock()
	clients[client.cid] = client
}

func removeClient(client client) {
	lock.Lock()
	defer lock.Unlock()
	cid := client.cid
	if _, ok := clients[cid]; ok {
		bridge <- event{status: "disconnect", client: clients[cid], time: time.Now()}
		delete(clients, cid)
	}
}

func getClient(id string) client {
	lock.Lock()
	defer lock.Unlock()

	return clients[id]
}

func Connect()  {
	var config = config.GetConfig()
	log.Println("Connect TS3 Bot")
	var err error
	conn, err = ts3.Dial(config.TS3Hostname, false)
	if err != nil {
		log.Fatal(err)
	}
}

// bot is a simple bot that checks version, signs in and sends a text message to
// channel#1 then exits.
func Poll() {
	log.Println("Start Polling TS3 Bot")
	defer conn.Cmd("quit")

	var config = config.GetConfig()

	conn.NotifyFunc(func (e string, s string) {
		handleEvent(e, s, conn)
		//fmt.Println(e + " " + s)
	})

	var cmds = []string{
		// Show version
		//"version",
		// Login
		"login serveradmin " + config.TS3Password,
		// Choose virtual server
		"use 1",
		// Update nickname
		`clientupdate client_nickname=Sushi\sBot`,
		// "clientlist",
		// Send message to channel with id=1
		//`sendtextmessage targetmode=2 target=1 msg=Big\sBrother\sis\snow\swatching\syou`,
		// Register to notify
		"servernotifyregister event=server",
	}

	for _, s := range cmds {
		// Send command and wait for its response
		r,_ := conn.Cmd(s)
		// Display as:
		//     > request
		//     response
		fmt.Printf("> %s\n%s", s, r)
		// Wait a bit after each command so we don't get banned. By default you
		// can issue 10 commands within 3 seconds.  More info on the
		// WHITELISTING AND BLACKLISTING section of TS3 ServerQuery Manual
		// (http://goo.gl/OpJXz).
		time.Sleep(350 * time.Millisecond)
	}
	go pingPong(conn)
	go handleSendEvent(bridge)

	for { select {} }
}

func handleSendEvent(bridge chan event) {
	e := <-bridge
	client := e.client.nickname

	if e.status == "connect" {
		if
		(client == "blaxxz") ||
			(client == "Sascha") ||
			(client == "Manni") ||
			(client == "Daniel") {
			telegram.Send(e.client.nickname + " connected (ts3)")
		}
		//telegram.SendToMe(e.client.nickname + " connected (ts3)")
	}
	if e.status == "disconnect" {
		nowTime := time.Now()
		conTime := e.client.conTime
		diff := nowTime.Sub(conTime)

		if
		(client == "blaxxz") ||
			(client == "Sascha") ||
			(client == "Manni") ||
			(client == "Daniel") {
			telegram.Send(e.client.nickname + " disconnected after " +shortDur(diff))
		}
		//telegram.SendToMe(e.client.nickname + " disconnected after " +shortDur(diff))
	}
	handleSendEvent(bridge)
}

func pingPong(conn *ts3.Conn) {
	time.Sleep(60 * time.Second)

	conn.Cmd("version")

	pingPong(conn)
}

func handleEvent(event string, response string, conn *ts3.Conn) {
	if event == "notifycliententerview" {
		//client enter
		handleConnectResponse(response)

	}
	if event == "notifyclientleftview" {
		//client leaves
		handleDisconnectResponse(response)
	}
}

func handleConnectResponse (s string) {
	arr := strings.Split(s, " ")
	cid := ""
	uid := ""
	nickname := ""
	for _, element := range arr {
		if strings.HasPrefix(element, "clid=") {
			cid = element[5:len(element)]
		}
		if strings.HasPrefix(element, "client_unique_identifier=") {
			uid = element[25:len(element)]
		}
		if strings.HasPrefix(element, "client_nickname=") {
			nickname = element[16:len(element)]
		}
	}
	client := client{cid: cid, uid: uid, nickname: nickname, conTime: time.Now()}
	addClient(client)

	bridge <- event{status: "connect", client: client, time: time.Now()}
}

func handleDisconnectResponse(s string) {
	arr := strings.Split(s, " ")
	cid := ""
	for _, element := range arr {
		if strings.HasPrefix(element, "clid=") {
			cid = element[5:len(element)]
		}
	}
	client := getClient(cid)
	removeClient(client)
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