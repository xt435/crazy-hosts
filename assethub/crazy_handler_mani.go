package assethub

import (
	"encoding/json"
	"fmt"

	"gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
)

var dataStoreMani DataStore
var session *mgo.Session

//InitDataStoreMani mongo init
func InitDataStoreMani() {
	dataStoreMani.session = connectToDb()
	session = dataStoreMani.session.Copy()
}

func deferToHell(mess string) {
	err := recover()
	if err != nil {
		fmt.Println(mess+"-inner-Error: ", err)
	}
}

//AssetRemover as named
func AssetRemover(serialNumber string) string {
	defer deferToHell("AssetRemover")
	assetHandler := AssetObject{}
	coll := initAssetHandlerConn(1, session)
	errFind := coll.Find(bson.M{"serialNumber": serialNumber}).One(&assetHandler)
	check := checkInnerErrAndReturnStr(errFind, "{\"TraceBackServer\" : \"there is no collection by \""+serialNumber+"}")
	if len(check) > 0 {
		return check
	}
	errRemove := coll.Remove(bson.M{"serialNumber": assetHandler.SerialNumber})
	check = checkInnerErrAndReturnStr(errRemove, "{\"TraceBackServer\" : \"failed to remove collection by \""+serialNumber+"}")
	if len(check) > 0 {
		return check
	}
	bufForRemoveAsset = append(bufForRemoveAsset, serialNumber)
	return "{\"TraceBackServer\" : \"I_AM_A_POLAR_BEAR\"}"
}

//HumanRemover as named
func HumanRemover(serialNumber string) string {
	defer deferToHell("HumanRemover")
	humanHandler := HumansContext{}
	coll := initAssetHandlerConn(2, session)
	errFind := coll.Find(bson.M{"humanSerialNumber": serialNumber}).One(&humanHandler)
	check := checkInnerErrAndReturnStr(errFind, "{\"TraceBackServer\" : \"there is no collection by \""+serialNumber+"}")
	if len(check) > 0 {
		return check
	}
	errRemove := coll.Remove(bson.M{"humanSerialNumber": humanHandler.HumanSerialNumber})
	check = checkInnerErrAndReturnStr(errRemove, "{\"TraceBackServer\" : \"failed to remove collection by \""+serialNumber+"}")
	if len(check) > 0 {
		return check
	}
	bufForRemoveHuman = append(bufForRemoveHuman, serialNumber)
	return "{\"TraceBackServer\" : \"I_AM_A_POLAR_BEAR\"}"
}

//AssetUpdate this method is depricated
func AssetUpdate(serialNumber string, jsonData string) string {
	defer deferToHell("AssetUpdate")
	assetHandler := AssetObject{}
	errJSON := json.Unmarshal([]byte(jsonData), &assetHandler)
	check := checkInnerErrAndReturnStr(errJSON, "{\"TraceBackServer\" : \"failed to parse collection data by \""+serialNumber+"}")
	if len(check) > 0 {
		return check
	}
	coll := initAssetHandlerConn(1, session)
	_, errUp := coll.Upsert(bson.M{"serialNumber": serialNumber}, assetHandler)
	check = checkInnerErrAndReturnStr(errUp, "{\"TraceBackServer\" : \"failed to update collection by \""+serialNumber+"}")
	if len(check) > 0 {
		return check
	}
	return "{\"TraceBackServer\" : \"I_AM_A_POLAR_BEAR\"}"
}

//HumanUpdate this is a depricated method
func HumanUpdate(serialNumber string, jsonData string) string {
	defer deferToHell("HumanUpdate")
	humanHandler := HumansContext{}
	errJSON := json.Unmarshal([]byte(jsonData), &humanHandler)
	check := checkInnerErrAndReturnStr(errJSON, "{\"TraceBackServer\" : \"failed to parse collection data by \""+serialNumber+"}")
	if len(check) > 0 {
		return check
	}
	coll := initAssetHandlerConn(2, session)
	_, errUp := coll.Upsert(bson.M{"serialNumber": serialNumber}, humanHandler)
	check = checkInnerErrAndReturnStr(errUp, "{\"TraceBackServer\" : \"failed to update collection by \""+serialNumber+"}")
	if len(check) > 0 {
		return check
	}
	return "{\"TraceBackServer\" : \"I_AM_A_POLAR_BEAR\"}"
}

func checkInnerErrAndReturnStr(err error, rtn string) string {
	if err != nil {
		fmt.Println(err)
		return rtn
	}
	return ""
}
