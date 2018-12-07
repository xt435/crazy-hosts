package monitorspack

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	hub "../assethub"
)

//InternalFairHead the head orders
type InternalFairHead interface {
	getFromPool()
	CheckContentAlive()
	CheckAliveOnly() bool
}

//AliveReport when getting status from serv
type AliveReport struct {
	Status     string `json:"status" bson:"status"`
	IPInfo     string `json:"ipInfo" bson:"ipInfo"`
	InsTime    int64  `json:"insTime" bson:"insTime"`
	ServerName string `json:"serverName" bson:"serverName"`
}

//SampleData saves some sample data
type SampleData struct {
	Content []interface{} `json:"content" bson:"content"`
	InsTime int64         `json:"insTime" bson:"insTime"`
}

//ConnectionData the configs
type ConnectionData struct {
	IPAddr string            `json:"ipAddr" bson:"ipAddr"`
	Port   int               `json:"port" bson:"port"`
	Gap    int64             `json:"gap" bson:"gap"`
	Novice map[string]string `json:"novice" bson:"novice"`
}

//ConfigsReader config
type ConfigsReader struct {
	ConnData []ConnectionData `json:"connData" bson:"connData"`
}

var config ConfigsReader

func getFromPool(file string) {
	configTemp := ConfigsReader{}
	configFile, err := os.Open(file)
	defer configFile.Close()
	hub.Check(err)
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&configTemp)
	config = configTemp
}

//CheckContentAlive alive check
func (conn ConnectionData) CheckContentAlive() {
	fmt.Printf("InternalFairs==%s", strconv.FormatInt(hub.CurrentMillis(), 10))
	conns := config.ConnData
	for i := range conns {
		var sq = make([]string, 0)
		var rq = make([]string, 0)
		go monFeedback(rq)
		time.Sleep(time.Millisecond * 1000)
		go MainClient(conns[i].IPAddr, conns[i].Port, sq, rq)
	}
}

//MonitorTimer start
func MonitorTimer(sq []string, novice map[string]string, gap int) {
	for {
		if len(sq) > 0 {
			sq = append(sq, novice["fatal"])
		}
		time.Sleep(time.Millisecond * 3000)
	}
}

func monFeedback(rq []string) {
	var r string
	for {
		if len(rq) > 0 {
			r = rq[0]
			fmt.Printf("FB==%s", r)
			rq = rq[1:]
		}
		time.Sleep(time.Millisecond * 1000)
	}
}
