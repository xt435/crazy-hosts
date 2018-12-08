package monitorspack

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"

	hub "../assethub"
)

type monitorspack struct {
	initTime string `json:"" bson:""`
	initInfo string `json:"" bson:""`
}

//MONITORS redis key for all configurations
const MONITORS string = "$#MONITORS@SERVERS"

//InternalFairHead the head runner
type InternalFairHead interface {
	checkContentAlive()
	checkAliveOnly() AliveReport
}

type conf struct {
	cmds []string
}

// type internalFairHead struct {
// 	provocate func(process InternalFairHead)
// }

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

var cacheRed *redis.Client

func redisClient(redisHost string, redisPort int, redisDB int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + strconv.Itoa(redisPort),
		Password: "",
		DB:       redisDB,
	})
	pong, err := client.Ping().Result()
	fmt.Println("*******************************REDIS CONN CHECK*******************************")
	fmt.Printf("=redis=%s == =redis=%s\n", pong, err)
	if pong == "PONG" {
		fmt.Println("REDIS CONNECTION DONE.")
	}
	fmt.Println("*******************************REDIS CONN CHECK*******************************")
	return client
}

//InitCache Please do this first.
func InitCache(host string, port int, db int) {
	cacheRed = redisClient(host, port, db)
}

// var Bmon internalFairHead
var cf conf
var initor monitorspack
var config ConfigsReader

//GetConfig this is the implement
func GetConfig() ConfigsReader {
	return config
}

//InitPool Please do this second.
func InitPool(initInfo string) {
	res := cacheRed.HGetAll(MONITORS)
	confs := res.Val()
	traverser := make([]string, 0)
	cds := []ConnectionData{}
	for i := range confs {
		traverser = append(traverser, confs[i])
		cd := ConnectionData{}
		err := json.Unmarshal([]byte(confs[i]), &cd)
		if err != nil {
			fmt.Println(err)
		} else {
			cds = append(cds, cd)
		}
	}
	config = ConfigsReader{ConnData: cds}
	cf = conf{cmds: traverser}
	initor = monitorspack{initTime: strconv.FormatInt(hub.CurrentMillis(), 10), initInfo: initInfo}
	//TODO save initor
	// Bmon = internalFairHead{}
}

//Mon is starter
func Mon(ii InternalFairHead) {
	ii.checkContentAlive()
}

func (cr ConfigsReader) checkContentAlive() {
	fmt.Printf("InternalFairs==%s", strconv.FormatInt(hub.CurrentMillis(), 10))
	conns := cr.ConnData
	for i := range conns {
		var sq = make([]string, 0)
		var rq = make([]string, 0)
		go monFeedback(rq)
		time.Sleep(time.Millisecond * 1000)
		go MainClient(conns[i].IPAddr, conns[i].Port, sq, rq)
	}
}

func (cr ConfigsReader) checkAliveOnly() AliveReport {
	return AliveReport{}
}

func monitorTimer(sq []string, novice map[string]string, gap int) {
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
		time.Sleep(time.Millisecond * 500)
	}
}
