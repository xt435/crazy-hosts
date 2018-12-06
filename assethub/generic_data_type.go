package assethub

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"encoding/json"
)

type GenericDataMock struct {
	PublicId   string            `json:"publicId"`
	SerialCode string            `json:"serialCode"`
	DataTree   []DataTreeContext `json:"dataTree"`
	CreateAt   string            `json:"createAt"`
}

type DataTreeContext struct {
	ParamName string            `json:"paramName"`
	ParamType string            `json:"paramType"`
	ParamVal  map[string]string `json:"paramVal"`
}

func checkBasicEntryData(dataJson string) (*GenericDataMock, string) {
	fmt.Println("dataJson==" + dataJson)
	if dataJson == "" {
		fmt.Println("input data is nil")
		return nil, "data not available"
	}
	res := GenericDataMock{}
	err := json.Unmarshal([]byte(dataJson), &res)
	if err != nil {
		return nil, "json format error"
	}
	return &res, "success"
}

func dataPushIn(mocker *GenericDataMock) *ResultOfPersist {
	//	if nil == mocker.PublicId || nil == mocker.SerialCode {
	//		log.Fatal("mocker is empty at " + time.Now().String())
	//	}
	session, err := mgo.Dial(mongod_main_one)
	checkErr(err)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(mongod_main_db).C(mongod_coll_name_generic_data)
	mocker.CreateAt = time.Now().String()
	err = c.Insert(mocker)
	checkErr(err)
	err = c.Find(bson.M{"publicid": mocker.PublicId}).One(&mocker)
	checkErr(err)
	fmt.Println(mocker.PublicId + "==all done." + mocker.CreateAt)
	return &ResultOfPersist{ResultOfPers: data_insert_ok, TimeLog: time.Now().String()}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}
}
