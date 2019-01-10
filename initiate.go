package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/alexflint/go-arg"

	hub "./assethub"
)

//GATEWAYWORD is actually base url
var GATEWAYWORD = "trollschain"

var sq = make(chan string)
var rq = make(chan string)

// cliArgs defines the configuration that the CLI
// expects. By using a struct we can very easily
// aggregate them into an object and check what are
// the expected types.
// If we need to mock this later it's just a matter
// of reusing the struct.
type cliArgs struct {
	Port int `arg:"-p,help:port to listen to"`
}

var (
	// args is a reference to an instantiation of
	// the configuration that the CLI expects but
	// with some values set.
	// By setting some values in advance we provide
	// default values that the user might provide
	// or not.
	args = &cliArgs{
		Port: 8896,
	}
)

func preCheck() {
	var isRun string
	var testLimit int
	fmt.Println("TeeHee, world. This is the truck chain lead program. My name is Go Crazy.\nI was created by necrophiliaccannibal.")
	fmt.Printf("==%s\n", "Yo, man, you wanna go for a performance test run? [y]=yes [n or any]=no")
	fmt.Scanf("%s\n", &isRun)
	if strings.Compare(isRun, "y") == 0 || strings.Compare(isRun, "yes") == 0 {
		fmt.Printf("==%s\n", "Ok, enter a number max 1000000")
		_, err := fmt.Scanf("%d\n", &testLimit)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		testLimit = 10
	}
	if testLimit > 1000000 {
		testLimit = 1000000
	}

	goThruGoMap(testLimit)
	var i = 0
	for i < 3 {
		var timeNow = time.Now().Format(time.RFC850)
		fmt.Println("TeeHee now we start -> " + strconv.Itoa(i+1) + "-" + timeNow)
		time.Sleep(time.Millisecond * 500)
		i++
	}

	arg.MustParse(args)
	fmt.Printf("must-port=%d\n", args.Port)
}

func main() {

	// preCheck()

	clientOfRedis := hub.RedisClient()
	hub.InitDataStoreHandlers()
	hub.InitDataStoreCrazyHandlers()
	hub.InitDataStoreHostPool()
	hub.InitDataStoreMani()
	hub.InitCache()
	hub.InitDataStoreHandlersMultiTrack()

	// go hub.InitGroupingProcess(clientOfRedis)
	// go hub.Grouping()
	// go hub.GroupingFinal()
	// go hub.DailyChainKeyGenerate()

	// not useful for this version::
	// go hub.AssetReceiverRunner(clientOfRedis)
	// go hub.VirtualContractReceiverRunner(clientOfRedis)
	// go hub.HumanReceiver(clientOfRedis)

	go hub.AssetSender(clientOfRedis)
	go hub.HumanSender(clientOfRedis)
	go hub.AssetRemoveSender(clientOfRedis)
	go hub.HumanRemoveSender(clientOfRedis)

	go hub.SyncOrigins(clientOfRedis)

	hub.HandlerForFuckers(GATEWAYWORD)
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
