package main

import (
	"fmt"
	"strconv"
	"time"

	// sock "./sockCommsPack"

	mgo "gopkg.in/mgo.v2"
)

func main() {

	fmt.Println("TeeHee, world. This is the truck chain lead program. My name is Go Crazy.\nI was created by necrophiliaccannibal.")
	var i = 0
	for i < 3 {
		var timeNow = time.Now().Format(time.RFC850)
		fmt.Println("TeeHee " + strconv.Itoa(i+1) + "-" + timeNow)
		time.Sleep(time.Millisecond * 500)
		i++
	}

	// go sock.MainClient("40.73.119.13", 9987) //being a client

	clientOfRedis := RedisClient()
	initDataStoreHandlers()
	initDataStoreCrazyHandlers()
	initDataStoreHostPool()
	initCache()
	initDataStoreHandlersMultiTrack()

	go initGroupingProcess(clientOfRedis)
	go grouping()
	go groupingFinal()
	go dailyChainKeyGenerate()

	go assetReceiverRunner(clientOfRedis)
	go virtualContractReceiverRunner(clientOfRedis)
	go humanReceiver(clientOfRedis)

	go assetSender(clientOfRedis)
	go humanSender(clientOfRedis)

	go syncOrigins(clientOfRedis)

	handlerForFuckers("trollschain")
}

func connectToDb() *mgo.Session {
	session, err := mgo.Dial(mongod_main_one)
	checkWithWarn(err)
	// defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	return session
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func checkWithWarn(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func testOnGoMap() {
	theMap := make([]map[string][]string, 0)
	one := make(map[string][]string)
	onelist := make([]string, 0)
	onelist = append(onelist, "tester")
	one["test"] = onelist
	for t := range one {
		for ii := range one[t] {
			fmt.Println(one[t][ii])
		}
	}
	theMap = append(theMap, one)
	for i := range theMap {
		for j := range theMap[i] {
			for l := range theMap[i][j] {
				fmt.Println(theMap[i][j][l])
			}
		}
	}
}
