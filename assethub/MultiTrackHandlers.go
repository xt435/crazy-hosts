package assethub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

var dataStoreMultiTrack DataStore

func InitDataStoreHandlersMultiTrack() {
	dataStoreMultiTrack.session = connectToDb()
}

const (
	SITE_NAME = "http://www.linde-china.com"
)

func HumpersJumpers(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("Error: ", err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Write([]byte("{\"HumpAndJump\":\"I_AM_GROOTS_COUSIN_FLOOT\"}"))
			w.WriteHeader(http.StatusAccepted)
		}
	}()
	vars := mux.Vars(r)
	pubId := vars["pid"]
	if pubId[0:2] == "AZ" {
		http.Redirect(w, r, SITE_NAME, http.StatusSeeOther)
	}
}

func ReportToSatan(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("Error: ", err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Write([]byte("{\"HumpAndJump\":\"I_AM_GROOTS_COUSIN_FLOOT\"}"))
			w.WriteHeader(http.StatusAccepted)
		}
	}()
	bodyData, errBody := ioutil.ReadAll(r.Body)
	if CheckErrWithRespond(w, errBody) {
		bodyDataContext := string(bodyData[:])
		fmt.Println("POSTDATA=" + bodyDataContext)
		req := make(map[string]string, 0)
		errReq := json.Unmarshal([]byte(bodyDataContext), &req)
		check := CheckErrWithRespond(w, errReq)
		if check {
			RecordHandler(req)
		}
	}
}

func RecordHandler(data map[string]string) {
	sess := dataStoreMultiTrack.session.Copy()
	coll := sess.DB(mongod_truck_db).C(mongod_coll_name_qrgen_rec)
	dataArr := strings.Split(data["gens"], "_")
	if len(dataArr) > 0 {
		for i := range dataArr {
			coll.Insert(bson.M{"qrcode": dataArr[i]})
		}
	}
}

func CheckErrWithRespond(w http.ResponseWriter, err error) bool {
	if err != nil {
		fmt.Println("Error: ", err)
		RespondWithFailure(w)
		return false
	}
	return true
}

func RespondWithFailure(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{\"HumpAndJump\":\"I_AM_GROOTS_COUSIN_FLOOT\"}"))
	w.WriteHeader(http.StatusAccepted)
}
