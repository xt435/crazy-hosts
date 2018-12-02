package main

/*
for redis:
go get -u github.com/go-redis/redis

the map-reduce instruction:
https://appliedgo.net/mapreduce/
*/
import (
	"encoding/json"
	"fmt"
	"math/rand"

	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var dataStoreHostPool DataStore
var session *mgo.Session

func initDataStoreHostPool() {
	dataStoreHostPool.session = connectToDb()
	session = dataStoreHostPool.session.Copy()
}

var buf = make([]string, 0)

var totalNumber int

func initGroupingProcess(client *redis.Client) {
	for {
		if totalNumber == 0 {
			ids, leng := getAllSerials(client)
			if leng > 0 {
				manageThePool(ids, client)
			}
		}
		time.Sleep(time.Millisecond * 1000 * 30)
	}
}

var runNumber int

func grouping() {
	for {
		if len(buf) > 0 {
			hostFlating(buf[0])
			buf = buf[1:]
			runNumber++
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func groupingFinal() {
	for {
		if runNumber > 0 && totalNumber > 0 && runNumber >= totalNumber {
			fmt.Println("group final process==" + time.Now().Format(time.RFC850))
			reduceToGroup()
			time.Sleep(time.Millisecond * 1000 * 1)
			runNumber = 0
			totalNumber = 0
		}
		time.Sleep(time.Millisecond * 2)
	}
}

/********************************************************************************************/
//	#######################################################################################
//	# the ip will be mapped as four elements [0][1][2][3] each element represents one 	  #
//	# section of the address domain. then saved to key-map in redis. key=id               #
//	#######################################################################################
/********************************************************************************************/
type HostContext struct {
	Mask string `json:"mask", bson:"mask"`
	Ip   string `json:"ip", bson:"ip"`
}

type HostMan struct {
	Ip    string   `json:"ip", bson:"ip"`
	Hosts []string `json:"hosts", bson:"hosts"`
}

func reduceToGroup() {
	coll := session.DB(mongod_truck_db).C(mongod_coll_name_host_pool)
	collHostMan := session.DB(mongod_truck_db).C(mongod_coll_name_pool_of_host_final)

	contexts := []HostContext{}
	err := coll.Find(bson.M{"mask": bson.M{"$ne": "null"}}).All(&contexts)
	checkWithWarn(err)
	if err == nil {
		tmpHostMan := HostMan{}
		for i := range contexts {
			con := contexts[i]
			errFound := collHostMan.Find(bson.M{"ip": con.Ip}).One(&tmpHostMan)
			if errFound == nil {
				continue
			}
			foundCon := []HostContext{}
			errIn := coll.Find(bson.M{"mask": con.Mask}).All(&foundCon)
			checkWithWarn(errIn)
			if errIn == nil {
				hostsArr := make([]string, 0)
				if len(foundCon) < pool_max {
					for inner := range foundCon {
						found := foundCon[inner]
						if found.Ip != "" {
							if found.Ip != "" {
								hostsArr = append(hostsArr, found.Ip)
							}
						}
					}
					for j := range hostsArr {
						host := hostsArr[j]
						if &host != nil && host != "" {
							hostMan := HostMan{Ip: host, Hosts: hostsArr}
							htmp := HostMan{}
							errTmp := collHostMan.Find(bson.M{"ip": con.Ip}).One(&htmp)
							if errTmp != nil {
								collHostMan.Insert(hostMan)
								fmt.Println("host::" + host + "==")
							} else {
								errTmp := collHostMan.Find(bson.M{"ip": con.Ip}).One(&htmp)
								if errTmp == nil {
									err_remove := collHostMan.Remove(&htmp)
									checkErr(err_remove)
									if err_remove == nil {
										collHostMan.Insert(hostMan)
										fmt.Println("host::" + host + "==")
									}
								}
							}
						}
					}
					hostsArr = hostsArr[:cap(hostsArr)]
					continue
				}
				for f := range foundCon {
					found := foundCon[f]
					if found.Ip != "" {
						hostsArr = append(hostsArr, found.Ip)
					}
					if len(hostsArr) == pool_max {
						for j := range hostsArr {
							host := hostsArr[j]
							if &host != nil && host != "" {
								hostMan := HostMan{Ip: host, Hosts: hostsArr}
								htmp := HostMan{}
								errTmp := collHostMan.Find(bson.M{"ip": con.Ip}).One(&htmp)
								if errTmp == nil {
									err_remove := collHostMan.Remove(&htmp)
									checkErr(err_remove)
									if err_remove == nil {
										collHostMan.Insert(hostMan)
										fmt.Println("host::" + host + "==")
									}
								} else {
									collHostMan.Insert(hostMan)
									fmt.Println("host::" + host + "==")
								}
							}
						}
						hostsArr = hostsArr[:cap(hostsArr)]
					}
				}
				if len(foundCon)%pool_max > 0 {
					count := len(foundCon) - len(foundCon)%pool_max
					for inner := count; inner < len(foundCon); inner++ {
						found := foundCon[inner]
						if found.Ip != "" {
							hostsArr = append(hostsArr, found.Ip)
						}
					}
					for j := range hostsArr {
						host := hostsArr[j]
						if &host != nil && host != "" {
							hostMan := HostMan{Ip: host, Hosts: hostsArr}
							htmp := HostMan{}
							errTmp := collHostMan.Find(bson.M{"ip": con.Ip}).One(&htmp)
							if errTmp == nil {
								err_remove := collHostMan.Remove(&htmp)
								checkErr(err_remove)
								if err_remove == nil {
									collHostMan.Insert(hostMan)
									fmt.Println("host::" + host + "==")
								}
							} else {
								collHostMan.Insert(hostMan)
								fmt.Println("host::" + host + "==")
							}
						}
					}
				}
			}
		}
	}
	_, errDrop := coll.RemoveAll(bson.D{})
	checkWithWarn(errDrop)
	time.Sleep(time.Millisecond * 1000 * 3)
}

func hostFlating(host string) {
	coll := session.DB(mongod_truck_db).C(mongod_coll_name_host_pool)

	maindict := strings.Split(buf[0], "_")
	dict := maindict[1] //being the ip
	dict = strings.Split(dict, ":")[0]
	if dict == "" {
		fmt.Println("found dict none val")
	}
	dictArr := strings.Split(dict, ".")
	threeKey := dictArr[0] + dictArr[1] + dictArr[2]
	twoKey := dictArr[0] + dictArr[1]
	oneKey := dictArr[0]
	fitToMap(threeKey, dict, coll)
	fitToMap(twoKey, dict, coll)
	fitToMap(oneKey, dict, coll)
}

func fitToMap(key string, dict string, coll *mgo.Collection) {
	ipGroupInsertNew(dict, key, coll)
}

func ipGroupInsertNew(dict string, key string, coll *mgo.Collection) {
	hostContextStruct := &HostContext{Mask: key, Ip: dict}
	err := coll.Insert(hostContextStruct)
	check(err)
}

func manageThePool(keyNameForHost []string, client *redis.Client) {
	for i := range keyNameForHost {
		val, errFetch := client.Get("HOST$" + keyNameForHost[i]).Result()
		if errFetch != nil {
			continue
		}
		if &val == nil || val == "" {
			continue
		}
		buf = append(buf, keyNameForHost[i]+"_"+val)
		totalNumber++
		fmt.Println("BUF==" + strconv.Itoa(i+1) + "==" + keyNameForHost[i] + "_" + val)
	}
}

func convBasicInfoToJson(ba *AssetObject) string {
	baJson, _ := json.Marshal(ba)
	return string(baJson)
}

func getAllSerials(redisCli *redis.Client) ([]string, int) {
	coll := session.DB(mongod_truck_db).C(mongod_coll_name_assets)
	bsMnts := []AssetObject{}
	errorForFind := coll.Find(nil).All(&bsMnts)
	check(errorForFind)
	fmt.Println("hostPoolCalculation-baseMnt size==" + strconv.Itoa(len(bsMnts)))
	rnt := []string{}
	var count = 0
	if len(bsMnts) > 0 {
		fmt.Println("device-ID-Pool::" + strconv.Itoa(len(bsMnts)))
		for i := range bsMnts {
			rnt = append(rnt, bsMnts[i].SerialNumber)
			// mockHosts(redisCli, "mocked"+strconv.Itoa(i))
			count++
		}
	}
	return rnt, len(rnt)
}

/*
	daily chain public key generator
*/
func dailyChainKeyGenerate() {
	for {
		tm := time.Now().Format(time.RFC3339)
		fmt.Println("TIMEFORCHAINKEY==" + tm)
		tm = strings.Split(tm, "T")[0]
		fmt.Println("db==" + mongod_truck_db + "+coll==" + mongod_coll_name_key_chain)
		coll := session.DB(mongod_truck_db).C(mongod_coll_name_key_chain)
		upperpart := ""
		judge, _ := strconv.Atoi(tm[9:])
		if judge >= 0 && judge < 3 {
			upperpart += "dexter"
		} else if judge >= 3 && judge < 6 {
			upperpart += "eliminate"
		} else if judge >= 6 && judge < 8 {
			upperpart += "foreign"
		} else if judge >= 8 && judge < 10 {
			upperpart += "goremegadon"
		} else {
			upperpart += "hallucinate"
		}
		yearVal, _ := strconv.Atoi(tm[0:4])
		monthVal, _ := strconv.Atoi(tm[5:7])
		dayVal, _ := strconv.Atoi(tm[8:])
		upperpart += strconv.Itoa(yearVal+monthVal+dayVal) + "jamesbond"
		fmt.Println("UPPERPART==" + upperpart)
		dailyChainKey := &DailyChainKey{Date: tm, Key: upperpart}
		check := &DailyChainKey{}
		errCheck := coll.Find(bson.M{"date": dailyChainKey.Date}).One(&check)
		if errCheck != nil {
			coll.Insert(dailyChainKey)
		}
		time.Sleep(time.Hour * 1 / 6)
	}
}

/**
this method is for test only.
*/
func mockHosts(redisCli *redis.Client, id string) {
	ipGen := ""
	judger := rand.Intn(4)
	if judger == 1 {
		ipGen = strconv.Itoa(210) + "." + strconv.Itoa(rand.Intn(5)) + "." +
			strconv.Itoa(rand.Intn(10)) + "." + strconv.Itoa(rand.Intn(254)) + ":" + strconv.Itoa(rand.Intn(65534))
	} else if judger == 2 {
		ipGen = strconv.Itoa(59) + "." + strconv.Itoa(rand.Intn(5)) + "." +
			strconv.Itoa(rand.Intn(10)) + "." + strconv.Itoa(rand.Intn(254)) + ":" + strconv.Itoa(rand.Intn(65534))
	} else if judger == 3 {
		ipGen = strconv.Itoa(198) + "." + strconv.Itoa(rand.Intn(5)) + "." +
			strconv.Itoa(rand.Intn(10)) + "." + strconv.Itoa(rand.Intn(254)) + ":" + strconv.Itoa(rand.Intn(65534))
	} else {
		ipGen = strconv.Itoa(34) + "." + strconv.Itoa(rand.Intn(5)) + "." +
			strconv.Itoa(rand.Intn(10)) + "." + strconv.Itoa(rand.Intn(254)) + ":" + strconv.Itoa(rand.Intn(65534))
	}
	redisCli.Set("HOST$"+id, ipGen, -1)
}
