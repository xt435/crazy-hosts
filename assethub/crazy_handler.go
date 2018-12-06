package assethub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var dataStoreCrazyHandle DataStore
var cache Cache

func InitDataStoreCrazyHandlers() {
	dataStoreCrazyHandle.session = connectToDb()
}

func InitCache() {
	cache.redisSession = RedisClient()
}

func fetchForAppUser(userId string, userName string) string {
	coll := initCrazyUserConn()
	res := ConfigAuth{}
	err := coll.Find(bson.M{"authid": userId, "authname": userName}).One(&res)
	if err != nil {
		return ""
	}
	return tokenGenerator(userId, userName)
}

func initCrazyUserConn() *mgo.Collection {
	session := MongoClient()
	return session.DB(mongod_main_db).C(mongod_coll_name_authrec)
}

func fetcherForAuthInfo(collName string, condition string) ConfigAuth {
	session, err := mgo.Dial(mongod_main_one)
	checkWithWarn(err)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	fmt.Println("dbinfo==" + mongod_main_db + "&" + mongod_coll_name_authrec)
	c := session.DB(mongod_main_db).C(mongod_coll_name_authrec)
	//	c.Insert(&ConfigAuth{AuthId: "009_doctor_who", AuthName: "doctorWho"})
	res := ConfigAuth{}
	fmt.Println("condition::" + condition + "&" + collName)
	errFin := c.Find(bson.M{collName: condition}).One(&res)
	fmt.Println(res)
	if errFin != nil {
		log.Fatal(errFin)
	}
	fmt.Println("1st_nickName::" + res.AuthName + ":" + res.AuthId)
	return res
}

func fetcherForHeadinfo(fieldName string, authName string) InteractInfo {
	session, err := mgo.Dial(mongod_main_one)
	checkWithWarn(err)
	defer session.Close()
	fmt.Println("fetcherForHeadinfo=" + fieldName + ":::" + authName)
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(mongod_main_db).C(mongod_coll_name_headinfo)
	res := InteractInfo{}
	errIi := c.Find(bson.M{fieldName: authName}).One(&res)
	if errIi != nil {
		log.Fatal(errIi)
		panic(errIi)
	}
	return res
}

func initAssetHandlerConn(isAssetOrHuman int, session *mgo.Session) *mgo.Collection {
	if isAssetOrHuman == 1 {
		collect := session.DB(mongod_truck_db).C(mongod_coll_name_assets)
		if collect == nil {
			collect := mgo.Collection{}
			collect.Name = mongod_coll_name_assets
			cInfo := mgo.CollectionInfo{}
			collect.Create(&cInfo)
		}
		return collect
	} else if isAssetOrHuman == 2 {
		collect := session.DB(mongod_truck_db).C(mongod_coll_name_humans)
		if collect == nil {
			collect := mgo.Collection{}
			collect.Name = mongod_coll_name_humans
			cInfo := mgo.CollectionInfo{}
			collect.Create(&cInfo)
		}
		return collect
	} else if isAssetOrHuman == 3 {
		collect := session.DB(mongod_truck_db).C(mongod_coll_name_bindingPool)
		if collect == nil {
			collect := mgo.Collection{}
			collect.Name = mongod_coll_name_bindingPool
			cInfo := mgo.CollectionInfo{}
			collect.Create(&cInfo)
		}
		return collect
	} else if isAssetOrHuman == 4 {
		collect := session.DB(mongod_truck_db).C(mongod_coll_name_chaintool_auth)
		if collect == nil {
			collect := mgo.Collection{}
			collect.Name = mongod_coll_name_chaintool_auth
			cInfo := mgo.CollectionInfo{}
			collect.Create(&cInfo)
		}
		return collect
	}
	return nil
}

func initChainToolUserConn(session *mgo.Session) *mgo.Collection {
	collect := session.DB(mongod_truck_db).C(mongod_coll_name_chaintool_auth)
	if collect == nil {
		collect := mgo.Collection{}
		collect.Name = mongod_coll_name_chaintool_auth
		cInfo := mgo.CollectionInfo{}
		collect.Create(&cInfo)
	}
	return collect
}

func assetsHandler(jsonData string, session *mgo.Session) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("assetsHandler-inner-Error: ", err)
		}
	}()
	assetHandler := AssetObject{}
	err := json.Unmarshal([]byte(jsonData), &assetHandler)
	checkWithWarn(err)
	//hashing on ledgers
	fmt.Println("=============================================main_asset_data=============================================")
	fmt.Println(assetHandler)
	fmt.Println("=============================================main_asset_data=============================================")
	if len(assetHandler.Origin) == 0 {
		assetHandler.Origin = "default"
	}
	mainHashNum := mainHash(&assetHandler)
	assetHandler.AssetHash = mainHashNum
	vcHash := assetHandler.VirtualContract.VCHash
	ledgeHash := hashing(mainHashNum + vcHash)
	assetHandler.VirtualContract.LedgerHash = ledgeHash
	assetHandler.AssetProperties.AssetApplication.LedgerHash = ledgeHash
	innerAssets := assetHandler.InnerAssets
	if len(innerAssets) > 0 {
		for i := range innerAssets {
			hashingOnInnerAssets(&innerAssets[i])
		}
	}
	coll := initAssetHandlerConn(1, session)
	if len(assetHandler.SerialNumber) == 0 {
		assetHandler.SerialNumber = fixOnAssetSerialNumber(assetHandler.Origin)
	} else {
		dup := AssetObject{}
		errFind := coll.Find(bson.M{"serialNumber": assetHandler.SerialNumber}).One(&dup)
		if errFind == nil {
			errRemove := coll.Remove(dup)
			checkWithWarn(errRemove)
		}
	}
	finalizedData, _ := json.Marshal(assetHandler)
	// fmt.Println(assetHandler)
	insertErr := coll.Insert(assetHandler)
	if insertErr != nil {
		fmt.Println(insertErr)
	} else {
		bufForAssetSender = append(bufForAssetSender, string(finalizedData[:]))
	}
}

func hashingOnInnerAssets(innerAsset *AssetObject) {
	//hashing on ledgers
	fmt.Println(innerAsset)
	mainHashNum := mainHash(innerAsset)
	innerAsset.AssetHash = mainHashNum
	vcHash := innerAsset.VirtualContract.VCHash
	ledgeHash := hashing(mainHashNum + vcHash)
	innerAsset.VirtualContract.LedgerHash = ledgeHash
	innerAsset.AssetProperties.AssetApplication.LedgerHash = ledgeHash
	innerAssets := innerAsset.InnerAssets
	if len(innerAssets) > 0 {
		for i := range innerAssets {
			hashingOnInnerAssets(&innerAssets[i])
		}
	}
}

func mainHash(asset *AssetObject) string {
	assetHandlerObjectStr := assetObjectHandler(asset)
	hash := hashing(assetHandlerObjectStr)
	return hash
}

func assetObjectHandler(asset *AssetObject) string {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("assetObjectHandler-inner-Error: ", err)
		}
	}()
	var rtn string
	itemList := make([]string, 0)
	asset.VirtualContract.VCSerialNumber = fixOnVirtualContractSerialNumber()
	ledger := getHashVirtualContract(&asset.VirtualContract)
	asset.VirtualContract.VCHash = ledger
	itemList = append(itemList, asset.SerialNumber)
	itemList = append(itemList, asset.Origin)
	itemList = append(itemList, asset.AssetGroup)
	itemList = append(itemList, asset.AssetName)
	itemList = append(itemList, asset.AssetCategoryName)
	itemList = append(itemList, asset.AssetMainType)
	itemList = append(itemList, getHashAssetProperties(&asset.AssetProperties))
	//TODO adding up components hash. for now, it's not so important.
	for i := range itemList {
		rtn = rtn + itemList[i]
	}
	log.Println("==" + rtn)
	return rtn
}

func humansHandler(jsonData string, session *mgo.Session) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("humansHandler-inner-Error: ", err)
		}
	}()
	humanObject := HumansContext{}
	errUnmarsh := json.Unmarshal([]byte(jsonData), &humanObject)
	checkWithWarn(errUnmarsh)
	//hash on human
	innerCons := humanObject.VirtualContracts
	if innerCons != nil && len(innerCons) > 0 {
		for i := range innerCons {
			innerHumanHash := humanHasher(&innerCons[i])
			innerCons[i].HumanHash = innerHumanHash
		}
	}
	humanHash := humanHasher(&humanObject)
	humanObject.HumanHash = humanHash
	coll := initAssetHandlerConn(2, session)
	dup := HumansContext{}
	errFind := coll.Find(bson.M{"humanSerialNumber": humanObject.HumanSerialNumber}).One(&dup)
	if errFind == nil && len(dup.HumanSerialNumber) > 0 {
		errRemove := coll.Remove(dup)
		checkWithWarn(errRemove)
	}
	finalizedData, _ := json.Marshal(humanObject)
	insertErr := coll.Insert(humanObject)
	if insertErr != nil {
		fmt.Println(insertErr)
	} else {
		bufForHumanSender = append(bufForHumanSender, string(finalizedData[:]))
	}
	collUserAuth := initChainToolUserConn(session)
	dupUser := ChainToolUserAuth{}
	errUserFind := collUserAuth.Find(bson.M{"serialNumber": humanObject.HumanSerialNumber}).One(&dupUser)
	if errUserFind == nil {
		errUserRemove := collUserAuth.Remove(dupUser)
		checkWithWarn(errUserRemove)
	}
	userAuth := createChainUserWithHumanObject(&humanObject)
	insertUserErr := collUserAuth.Insert(userAuth)
	if insertUserErr != nil {
		fmt.Println(insertUserErr)
	}
}

func createChainUserWithHumanObject(humanObject *HumansContext) *ChainToolUserAuth {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("createChainUserWithHumanObject-inner-Error: ", err)
		}
	}()
	var custList = humanObject.CustomContext
	var phoneNumber string
	if len(custList) > 0 {
		for i := range custList {
			if custList[i].ParamName == "phoneNumber" {
				phoneNumber = custList[i].ParamValue
				break
			}
		}
		fmt.Println("make user for chaintool::" + phoneNumber)
		userAuth := ChainToolUserAuth{}
		userAuth.Origin = humanObject.Origin
		userAuth.UserName = phoneNumber
		userAuth.Password = strings.ToUpper(hashing(humanObject.HumanName +
			humanObject.HumanSerialNumber + phoneNumber)[0:8])
		userAuth.CreateTime = currentMilliseconds()
		userAuth.SerialNumber = humanObject.HumanSerialNumber
		return &userAuth
	}
	return nil
}

func humanHasher(humanObject *HumansContext) string {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("humanHasher-inner-Error: ", err)
		}
	}()
	var rtn bytes.Buffer
	itemList := make([]string, 0)
	if len(humanObject.Origin) == 0 {
		humanObject.Origin = "default"
	}
	if len(humanObject.HumanSerialNumber) == 0 {
		humanObject.HumanSerialNumber = fixOnHumanSerialNumber(humanObject.Origin, humanObject.HumanName)
		itemList = append(itemList, humanObject.HumanSerialNumber)
	}
	itemList = append(itemList, humanObject.Origin)
	itemList = append(itemList, humanObject.HumanName)
	itemList = append(itemList, humanObject.HumanGroup)
	itemList = append(itemList, humanObject.Passwd)
	itemList = append(itemList, hashing(humanObject.Info))
	itemList = append(itemList, getHashGenericContexts(humanObject.CustomContext))
	itemList = append(itemList, strconv.Itoa(humanObject.CanMaintain))
	vircons := humanObject.VirtualContracts
	if len(vircons) > 0 {
		for i := range vircons {
			itemList = append(itemList, vircons[i].HumanHash)
		}
	}
	itemList = append(itemList, strconv.Itoa(humanObject.DeadOrAlive))
	t := currentMilliseconds()
	humanObject.BirthTime = t
	itemList = append(itemList, strconv.FormatInt(humanObject.BirthTime, 10))
	assetAccess := humanObject.AssetAccess
	if len(assetAccess) > 0 {
		for i := range assetAccess {
			itemList = append(itemList, assetAccess[i])
		}
	}
	itemList = append(itemList, getHashHumanTrust(&humanObject.HumanTrust))
	itemList = append(itemList, getHashVerificationInfo(&humanObject.VerificationInfo))
	for i := range itemList {
		rtn.WriteString(itemList[i])
	}
	hash := hashing(rtn.String())
	return hash
}

func bindingHandler(jsonData string, session *mgo.Session) string {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("bindingHandler-inner-Error: ", err)
		}
	}()
	bindData := []BindingBoundagePool{}
	err := json.Unmarshal([]byte(jsonData), &bindData)
	checkWithWarn(err)
	coll := initAssetHandlerConn(3, session)
	if len(bindData) > 0 {
		for i := range bindData {
			dup := BindingBoundagePool{}
			check := bindData[i]
			tabReplace := strings.Replace(check.BindContent, "\t", "", -1)
			backReplace := strings.Replace(tabReplace, "\n", "", -1)
			check.BindContent = backReplace
			errFind := coll.Find(bson.M{"bindingSerial": check.BindingSerial, "origin": check.Origin,
				"bindingContent": check.BindContent, "bindFlag": check.BindFlag}).One(&dup)
			if errFind == nil && len(dup.BindingSerial) > 0 {
				continue
			}
			insertErr := coll.Insert(check)
			if insertErr != nil {
				fmt.Println(insertErr)
				return "insertErr"
			}
		}
	}
	return "done"
}

type OriginStruct struct {
	originName []string
}

func SyncOrigins(redisCli *redis.Client) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("syncOrigins-inner-Error: ", err)
		}
	}()
	coll := initAssetHandlerConn(4, session)
	if coll != nil {
		for {
			pipeline := []bson.M{bson.M{"$group": bson.M{"_id": "$origin"}}}
			p := coll.Pipe(pipeline)
			resp := []bson.M{}
			errPipe := p.All(&resp)
			var origins string
			if errPipe == nil {
				for item := range resp {
					ori := resp[item]["_id"].(string)
					origins = origins + (ori + "|")
				}
			}
			if len(origins) > 0 && origins[len(origins)-1:] == "|" {
				origins = origins[0 : len(origins)-1]
				redisCli.Del("#$-system-misc-")
				redisCli.Append("#$-system-misc-", origins)
			}
			time.Sleep(time.Second * 60)
		}
	}
}

func findOrigins() *OriginStruct {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("findOrigins-inner-Error: ", err)
		}
	}()
	res, err := cache.redisSession.Get("#$-system-misc-").Result()
	if err == nil {
		os := OriginStruct{}
		if strings.Contains(res, "|") {
			arr := strings.Split(res, "|")
			os.originName = arr
		} else {
			resArr := make([]string, 0)
			resArr = append(resArr, res)
			os.originName = resArr
		}
		return &os
	} else {
		return nil
	}
}
