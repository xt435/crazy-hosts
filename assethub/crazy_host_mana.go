package assethub

import (
	_ "encoding/json"
	"fmt"

	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var GAPTIME int64

//AssetsTmp as named
type AssetsTmp struct {
	SerialNumber string `json:"serialNumber" bson:"serialNumber"`
}

//HostPoolRec truckchain gate center persist.py expand_host_pool
type HostPoolRec struct {
	PublicId string `json:"publicId" bson:"publicId"`
	HostInfo string `json:"hostInfo" bson:"hostInfo"`
}

//SyncHostPool as named
func syncHostPool(client *redis.Client, session *mgo.Session) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("SyncHostPool-inner-Error: ", err)
		}
	}()
	GAPTIME, _ = strconv.ParseInt("600000", 10, 64)
	coll := session.DB(mongod_truck_db).C(mongod_coll_name_assets)
	collRec := session.DB(mongod_truck_db).C(mongod_coll_name_host_pool_rec)
	assets := make([]AssetsTmp, 0)
	recs := make([]HostPoolRec, 0)
	for {
		// errFind := coll.Find(bson.M{"$project": bson.M{"serialNumber": 1, "_id": 0}}).All(&assets)
		errFind := coll.Find(bson.M{}).Select(bson.M{"serialNumber": 1}).All(&assets)
		if errFind != nil {
			fmt.Println(errFind)
			time.Sleep(time.Millisecond * 10000)
			continue
		}
		if len(assets) > 0 {
			for i := range assets {
				rdKey := assets[i]
				hb, _ := client.Get("HB$" + rdKey.SerialNumber).Result()
				if len(hb) > 0 {
					hbNum, _ := strconv.ParseInt(hb, 10, 64)
					fmt.Printf("devi:%s  hbtime:%s\n", rdKey, time.Unix(hbNum, 0).Format(time.RFC850))
					if currentMilliseconds()-hbNum > GAPTIME {
						_, _ = client.Del("HB$" + rdKey.SerialNumber).Result()
						errFind = collRec.Find(bson.M{"publicId": rdKey.SerialNumber}).All(&recs)
						if errFind == nil {
							if len(recs) > 0 {
								for j := range recs {
									_ = collRec.Remove(bson.M{"publicId": recs[j].PublicId, "hostInfo": recs[j].HostInfo})
									ho, _ := client.Get("HOST$" + recs[j].PublicId).Result()
									if strings.Compare(ho, recs[j].HostInfo) == 0 {
										_, _ = client.Del("HOST$" + recs[j].PublicId).Result()
									}
									fmt.Printf("publicId:%s  hostInfo:%s  heartbeat cleared due to overdue time\n", recs[j].PublicId, recs[j].HostInfo)
								}
							}
						}
					}
				}
			}
		}
		time.Sleep(time.Millisecond * 60000 * 3)
	}
}

//this is the pre stage thing for ip grouping
/*
def save_heart_beat(ident, heartbeat):
    last_beat = fetch_key_val('HB$' + ident)
    if last_beat is None:
		save_key_val('HB$' + ident, str(heartbeat))
get HB$8ff7d7af11b2c16c7a2f42e12fe232e233840442
"1546830565619# 2019-01-07 10:37:38:277#1546828658278"
get HOST$8ff7d7af11b2c16c7a2f42e12fe232e233840442
"192.168.88.251:34316"
*/
func groupingByOrigin(client *redis.Client, session *mgo.Session) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("groupingByOrigin-inner-Error: ", err)
		}
	}()
	poolCacheGrp := make(map[string][]string, 0)
	allHearts, errKeys := client.Keys("HOST$*").Result()
	if errKeys != nil {
		fmt.Println(errKeys)
		return
	}
	if len(allHearts) > 0 {
		coll := session.DB(mongod_truck_db).C(mongod_coll_name_assets)
		var ori string
		for i := range allHearts {
			errFind := coll.Find(bson.M{"serialNumber": allHearts[i][5:], "$project": bson.M{"origin": 1, "_id": 0}}).One(&ori)
			if errFind != nil {
				continue
			}
			k, _ := client.Get(allHearts[i]).Result()
			if poolCacheGrp[ori] == nil || len(poolCacheGrp[ori]) == 0 {
				fmt.Println(k)
				firstBlood := make([]string, 0)
				firstBlood = append(firstBlood, allHearts[i][5:]+"#"+k)
				poolCacheGrp[ori] = firstBlood
			} else {
				poolCacheGrp[ori] = append(poolCacheGrp[ori], allHearts[i][5:]+"#"+k)
			}
		}
		//TODO make dir
		if len(poolCacheGrp) > 0 {
			hcs := make([]HostContext{}, 0)
			for k, v := range poolCacheGrp {
				if len(v) > 0 {
					for i := range v {
						hc := HostContext{}
						hc.Mask = k + strings.Split(v[i], "#")[0]
						hc.Ip = strings.Split(v[i], "#")[1]
						hcs = append(hcs, hc)
					}
				}
			}
			fmt.Println("secondBlood")
			time.Sleep(5)
			for key, val := range hcs {
				fmt.Printf("running HostContexts==%s", key)
				reduceToGrp(val)
			}
		}
	}
}
