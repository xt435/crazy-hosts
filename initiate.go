package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	sock "./monitorspack"
)

var GATE_WAY_WORD = "trollschain"

const (
	RedMonIP   = "192.168.204.145"
	RedMonPort = 6379
	RedMonDB   = 0
)

var sq = make([]string, 0)
var rq = make([]string, 0)
var wg sync.WaitGroup

func main() {

	// var isRun string
	// var testLimit int
	// fmt.Println("TeeHee, world. This is the truck chain lead program. My name is Go Crazy.\nI was created by necrophiliaccannibal.")
	// fmt.Printf("==%s\n", "Yo, man, you wanna go for a performance test run? [y]=yes [n or any]=no")
	// fmt.Scanf("%s\n", &isRun)
	// if strings.Compare(isRun, "y") == 0 || strings.Compare(isRun, "yes") == 0 {
	// 	fmt.Printf("==%s\n", "Ok, enter a number max 1000000")
	// 	_, err := fmt.Scanf("%d\n", &testLimit)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// } else {
	// 	testLimit = 10
	// }
	// if testLimit > 1000000 {
	// 	testLimit = 1000000
	// }

	// goThruGoMap(testLimit)
	// var i = 0
	// for i < 3 {
	// 	var timeNow = time.Now().Format(time.RFC850)
	// 	fmt.Println("TeeHee now we start -> " + strconv.Itoa(i+1) + "-" + timeNow)
	// 	time.Sleep(time.Millisecond * 500)
	// 	i++
	// }

	// clientOfRedis := hub.RedisClient()

	// sock.InitCache(RedMonIP, RedMonPort, RedMonDB)
	// sock.InitPool("from 127.0.0.1")
	// sock.Mon(sock.GetConfig())

	go monFeedback(rq)
	go sendCheck(sq)
	sock.SyncWG(wg)
	sock.MainClient("127.0.0.1", 9999, sq, rq)

	// hub.InitDataStoreHandlers()
	// hub.InitDataStoreCrazyHandlers()
	// hub.InitDataStoreHostPool()
	// hub.InitCache()
	// hub.InitDataStoreHandlersMultiTrack()

	// go hub.InitGroupingProcess(clientOfRedis)
	// go hub.Grouping()
	// go hub.GroupingFinal()
	// go hub.DailyChainKeyGenerate()

	// go hub.AssetReceiverRunner(clientOfRedis)
	// go hub.VirtualContractReceiverRunner(clientOfRedis)
	// go hub.HumanReceiver(clientOfRedis)

	// go hub.AssetSender(clientOfRedis)
	// go hub.HumanSender(clientOfRedis)

	// go hub.SyncOrigins(clientOfRedis)

	// hub.HandlerForFuckers(GATE_WAY_WORD)
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

func sendCheck(sq []string) {
	for {
		wg.Add(1)
		sq = append(sq, "0123456789ABCDEF")
		fmt.Printf("check fucking head==%s %d\n", sq[len(sq)-1], len(sq))
		time.Sleep(time.Millisecond * 3000)
	}
}

func goThruGoMap(limit int) {
	st := currentMilliseconds()
	startTime := time.Now()
	stStr := startTime.Format(time.RFC3339)
	theMap := make([]map[string][]string, 0)
	one := make(map[string][]string)
	var k = 0
	for k < limit {
		onelist := make([]string, 0)
		if k%2 == 0 {
			judger := rand.Intn(4)
			var j = 0
			for j < judger {
				onelist = append(onelist, "--"+time.Now().Format(time.RFC3339))
				j++
			}
		} else {
			onelist = append(onelist, "-"+time.Now().Format(time.RFC850))
		}
		one["cairo"+strconv.Itoa(k)] = onelist
		if k%5 == 0 {
			fmt.Printf("%s", "+")
		}
		k++
	}
	theMap = append(theMap, one)
	for i := range theMap {
		for j := range theMap[i] {
			for l := range theMap[i][j] {
				fmt.Printf("runner==%s\n", theMap[i][j][l])
			}
		}
	}
	ed := currentMilliseconds()
	endTime := time.Now()
	ndStr := endTime.Format(time.RFC3339)
	theMap = make([]map[string][]string, 0)
	one = make(map[string][]string)
	fmt.Printf("started-at::%s\n", stStr)
	time.Sleep(time.Millisecond * 1000)
	fmt.Printf("ended-at::%s\n", ndStr)
	time.Sleep(time.Millisecond * 1000)
	fmt.Printf("time elapse::%s (seconds)\n", strconv.FormatInt((ed-st)/1000, 10))
	time.Sleep(time.Millisecond * 2000)
}

func currentMilliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
