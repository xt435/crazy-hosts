package assethub

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

var dataStore DataStore

//InitDataStoreHandlers = all init.
func InitDataStoreHandlers() {
	dataStore.session = connectToDb()
}

func deferCommon(name string) {
	err := recover()
	if err != nil {
		fmt.Println(name+"-handler-Error: ", err)
	}
}

func checkPath(path string) string {
	if path != "" && path[0:1] != "/" {
		path = "/" + path
	}
	return path
}

//HandlerForFuckers as named. only fuckers will go thru.
func HandlerForFuckers(path string) {
	defer deferCommon("handlerForFuckers")
	path = checkPath(path)
	fmt.Println("default_server_port=" + default_server_port)
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/getme/fuckers"+path, Index)
	r.HandleFunc(path+novice_path, Novice)
	r.HandleFunc(path+asset_saving_path, receiverOfAssets).Methods("POST")
	r.HandleFunc(path+human_saving_path, receiverOfHumans).Methods("POST")
	r.HandleFunc(path+binding_check_path, receiveOfBindingBoundagePool).Methods("POST")
	r.HandleFunc(path+binding_check_delete_path, deleteFromBoundagePool).Methods("DELETE")
	r.HandleFunc(path+chain_tool_user, chainToolUserHandler).Methods("POST")
	r.HandleFunc(path+chain_tool_user_human, chainToolUserGetHumanObjectHandler).Methods("POST")
	r.HandleFunc(path+origin_manager, originManager).Methods("POST")
	r.HandleFunc(path+humpers_jumpers, HumpersJumpers)
	r.HandleFunc(path+recordqr, ReportToSatan).Methods("POST")
	r.HandleFunc(path+asset_remove_path, removeAsset).Methods("DELETE")
	r.HandleFunc(path+human_remove_path, removeHuman).Methods("DELETE")
	log.Fatal(http.ListenAndServe(default_server_port, r))
}

func sendToHellCommonHead(name string, sendback string, status int, w http.ResponseWriter) {
	err := recover()
	if err != nil {
		fmt.Println(name+"-handler-Error: ", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte(sendback))
		w.WriteHeader(status)
	}
}

func sendBackCommonHead(back string, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(back))
	w.WriteHeader(status)
}

func sendBackWithAuth(back string, status int, w http.ResponseWriter, auth string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set(asset_auth, auth)
	w.Write([]byte(back))
	w.WriteHeader(status)
}

func dealWithAuthToken(w http.ResponseWriter, r *http.Request) string {
	authToken := r.Header.Get(asset_auth)
	if len(authToken) == 0 {
		sendBackCommonHead("{\"TraceBackServer\":\"You r not authorized. please get token from novice req\"}", http.StatusForbidden, w)
		return "fucked"
	}
	fmt.Println("#AUTHENTICATION_ASSET_INSERT::" + authToken)
	userPass := r.Header.Get(header_identity)
	checkToken := tokenGenerator(strings.Split(userPass, "@")[0], strings.Split(userPass, "@")[1])
	if len(checkToken) != len(authToken) || checkToken != authToken {
		sendBackCommonHead("{\"TraceBackServer\":\"YouAreAnAlienYouAreNotAllowed\"}", http.StatusForbidden, w)
		return "fucked"
	}
	return authToken
}

func dealWithBodyData(w http.ResponseWriter, r *http.Request) string {
	bodyData, errBody := ioutil.ReadAll(r.Body)
	bodyDataContext := string(bodyData[:])
	fmt.Println("POSTDATA=" + bodyDataContext)
	if errBody != nil || bodyDataContext == "null" {
		fmt.Println("Error: ", errBody)
		sendBackCommonHead("{\"TraceBackServer\":\"YourRequestIsScrambled\"}", http.StatusNotAcceptable, w)
		return "fucked"
	}
	return bodyDataContext
}

func dealWithPatterned(w http.ResponseWriter, r *http.Request) (string, string) {
	auth := dealWithAuthToken(w, r)
	if strings.Compare(auth, "fucked") == 0 {
		return "fucked", "fucked"
	}
	body := dealWithBodyData(w, r)
	if strings.Compare(body, "fucked") == 0 {
		return "fucked", "fucked"
	}
	return auth, body
}

// saving assets data
func receiverOfAssets(w http.ResponseWriter, r *http.Request) {
	defer sendToHellCommonHead("receiverOfAssets", "{\"TraceBackServer\":\"Messy Request\"}", http.StatusNoContent, w)
	auth, body := dealWithPatterned(w, r)
	if strings.Compare(auth, "fucked") == 0 || strings.Compare(body, "fucked") == 0 {
		return
	}
	assetsHandler(body, dataStore.session.Copy())
	sendBackWithAuth("{\"TraceBackServer\":\"I_AM_A_POLAR_BEAR\"}", http.StatusOK, w, auth)
}

// saving human data
func receiverOfHumans(w http.ResponseWriter, r *http.Request) {
	defer sendToHellCommonHead("receiverOfHumans", "{\"TraceBackServer\":\"Messy Request\"}", http.StatusNoContent, w)
	auth, body := dealWithPatterned(w, r)
	if strings.Compare(auth, "fucked") == 0 || strings.Compare(body, "fucked") == 0 {
		return
	}
	humansHandler(body, dataStore.session.Copy())
	sendBackWithAuth("{\"TraceBackServer\":\"I_AM_A_POLAR_BEAR\"}", http.StatusOK, w, auth)
}

//removeAsset is as it's named
func removeAsset(w http.ResponseWriter, r *http.Request) {
	defer sendToHellCommonHead("removeAsset", "{\"TraceBackServer\":\"Messy Request\"}", http.StatusNoContent, w)
	auth, body := dealWithPatterned(w, r)
	if strings.Compare(auth, "fucked") == 0 || strings.Compare(body, "fucked") == 0 {
		return
	}
	sendBackWithAuth(AssetRemover(body), http.StatusOK, w, auth)
}

//the remove of human data
func removeHuman(w http.ResponseWriter, r *http.Request) {
	defer sendToHellCommonHead("removeHuman", "{\"TraceBackServer\":\"Messy Request\"}", http.StatusNoContent, w)
	auth, body := dealWithPatterned(w, r)
	if strings.Compare(auth, "fucked") == 0 || strings.Compare(body, "fucked") == 0 {
		return
	}
	sendBackWithAuth(HumanRemover(body), http.StatusOK, w, auth)
}

// saving binding pool data
func receiveOfBindingBoundagePool(w http.ResponseWriter, r *http.Request) {
	defer sendToHellCommonHead("receiveOfBindingBoundagePool", "{\"TraceBackServer\":\"Messy Request\"}", http.StatusNoContent, w)
	auth, body := dealWithPatterned(w, r)
	if strings.Compare(auth, "fucked") == 0 || strings.Compare(body, "fucked") == 0 {
		return
	}
	var res string
	if strings.Contains(body, "fieldName") {
		res = bindingHandler(body, dataStore.session.Copy(), 0) //content objectified.
	} else {
		res = bindingHandler(body, dataStore.session.Copy(), 1) //content flat string.
	}
	if res != "done" {
		sendBackCommonHead("{\"TraceBackServer\":\"YourRequestIsScrambled\"}", http.StatusNotAcceptable, w)
		return
	}
	sendBackWithAuth("{\"TraceBackServer\":\"I_AM_A_POLAR_BEAR\"}", http.StatusOK, w, auth)
}

// remove from binding pool.
func deleteFromBoundagePool(w http.ResponseWriter, r *http.Request) {
	defer sendToHellCommonHead("deleteFromBoundagePool", "{\"TraceBackServer\":\"Messy Request\"}", http.StatusNoContent, w)
	auth, body := dealWithPatterned(w, r)
	if strings.Compare(auth, "fucked") == 0 || strings.Compare(body, "fucked") == 0 {
		return
	}
	res := bindingDeleteHandler(body, dataStore.session.Copy())
	if res != "done" {
		sendBackCommonHead("{\"TraceBackServer\":\"YourRequestIsScrambled\"}", http.StatusNotAcceptable, w)
		return
	}
	sendBackWithAuth("{\"TraceBackServer\":\"I_AM_A_POLAR_BEAR\"}", http.StatusOK, w, auth)
}

// chain tool user(for various logins) saving
func chainToolUserHandler(w http.ResponseWriter, r *http.Request) {
	defer sendToHellCommonHead("chainToolUserHandler", "{\"TraceBackServer\":\"Messy Request\"}", http.StatusNoContent, w)
	bodyData, errBody := ioutil.ReadAll(r.Body)
	bodyDataContext := string(bodyData[:])
	fmt.Println("POSTDATA=" + bodyDataContext)
	if errBody != nil || bodyDataContext == "null" || bodyDataContext == "" {
		fmt.Println("Error: ", errBody)
		sendBackCommonHead("{\"TraceBackServer\":\"YourRequestIsScrambled\"}", http.StatusNotAcceptable, w)
		return
	}
	user := ChainToolUserAuthMask{}
	errParse := json.Unmarshal([]byte(bodyDataContext), &user)
	checkWithWarn(errParse)
	resAuth := chainToolAuthentication(&user)
	if resAuth == "{\"auth\": \"I_AM_A_POLAR_BEAR\"}" {
		genToken := tokenGeneratorTimely(user.Origin, user.UserName+"@"+user.Password)
		sendBackWithAuth("{\"TraceBackServer\":\"I_AM_A_POLAR_BEAR\", \"token\":\""+genToken+"\"}", http.StatusOK, w, genToken)
	} else {
		sendBackCommonHead(resAuth, http.StatusBadRequest, w)
	}
}

func chainToolUserGetHumanObjectHandler(w http.ResponseWriter, r *http.Request) {
	defer sendToHellCommonHead("chainToolUserGetHumanObjectHandler", "{\"TraceBackServer\":\"Messy Request\"}", http.StatusNoContent, w)
	bodyData, errBody := ioutil.ReadAll(r.Body)
	bodyDataContext := string(bodyData[:])
	fmt.Println("POSTDATA=" + bodyDataContext)
	if errBody != nil || bodyDataContext == "null" || bodyDataContext == "" {
		fmt.Println("bodyData Error: ", errBody)
		sendBackCommonHead("{\"TraceBackServer\":\"YourRequestIsScrambled\"}", http.StatusNotAcceptable, w)
		return
	}
	user := ChainToolUserAuthMask{}
	errParse := json.Unmarshal([]byte(bodyDataContext), &user)
	checkWithWarn(errParse)
	s := dataStore.session.Copy()
	collUserAuth := s.DB(mongod_truck_db).C(mongod_coll_name_chaintool_auth)
	userAuth := ChainToolUserAuth{}
	errUserAuth := collUserAuth.Find(bson.M{"origin": user.Origin, "userName": user.UserName, "password": user.Password}).One(&userAuth)
	if errUserAuth != nil {
		fmt.Println("errUserAuth Error: ", errUserAuth)
		sendBackCommonHead("{\"TraceBackServer\":\"requesting user not found\"}", http.StatusNotAcceptable, w)
		return
	}
	collHuman := s.DB(mongod_truck_db).C(mongod_coll_name_humans)
	humanFound := HumansContext{}
	errHumanFound := collHuman.Find(bson.M{"humanSerialNumber": userAuth.SerialNumber}).One(&humanFound)
	if errHumanFound != nil {
		sendBackCommonHead("{\"TraceBackServer\":\"no basic info found\"}", http.StatusNotAcceptable, w)
		return
	}
	b, errByte := json.Marshal(humanFound)
	if errByte != nil {
		panic(errByte)
	}
	fmt.Println("BACKING==" + string(b[:]))
	iAmAPollarBear := "{\"TraceBackServer\":\"I_AM_A_POLAR_BEAR\", \"humanContext\":" + string(b[:]) + "}"
	sendBackCommonHead(iAmAPollarBear, http.StatusOK, w)
}

//Index feeds back timing and current version of app.
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("version of current:" + version_name_of_app)
	fmt.Fprintln(w, version_name_of_app, html.EscapeString(" AT:"+time.Now().String()))
}

//Novice is auth
func Novice(w http.ResponseWriter, r *http.Request) {
	defer sendToHellCommonHead("Novice", "{\"TraceBackServer\":\"Messy Request\"}", http.StatusNoContent, w)
	vars := mux.Vars(r)
	auth := ConfigAuth{AuthId: vars["reqid"], AuthName: vars["reqname"]}
	fmt.Println("id:" + auth.AuthId + "_____name:" + auth.AuthName)
	if auth.AuthId == "007" && auth.AuthName == "JamesBond" { //broadcast test account
		fmt.Fprintln(w, "answer:", html.UnescapeString("{\"result\":\"allRite\"}"))
		interfaceInfo := forgeFeedbackOnInterfaceInfo(&auth)
		infos := InteractInfos{interfaceInfo}
		b, errors := json.Marshal(infos)
		if errors != nil {
			panic(errors)
		}
		fmt.Println(string(b[:]))
		if err := json.NewEncoder(w).Encode(infos); err != nil {
			panic(err)
		}
		iAmAPollarBear := "{\"TraceBackServer\": \"" + string(b[:]) + " AT: " + time.Now().String() + "\"}"
		sendBackCommonHead(iAmAPollarBear, http.StatusOK, w)
	} else {
		token := fetchForAppUser(auth.AuthId, auth.AuthName)
		fmt.Println("TOKENGEN==" + token)
		if len(token) == 0 {
			sendBackCommonHead("{\"TraceBackServer\":\"YOU_ARE_A_WITCH_-_NOT_ALLOWED_TO_CAST\"}", http.StatusForbidden, w)
		} else {
			sendBackWithAuth("{\"TraceBackServer\":\"I_AM_A_POLAR_BEAR\", \"token\":\""+token+"\"}", http.StatusOK, w, token)
		}
	}
}

func forgeFeedbackOnInterfaceInfo(auth *ConfigAuth) InteractInfo {
	fmt.Println("authentication:" + auth.AuthName)
	confAuth := fetcherForAuthInfo("authid", auth.AuthId+"_"+auth.AuthName)
	fmt.Println(confAuth.AuthId + "::" + confAuth.AuthName)
	fieldName := "authname"
	return fetcherForHeadinfo(fieldName, auth.AuthName)
}

//login seletor
func originManager(w http.ResponseWriter, r *http.Request) {
	defer sendToHellCommonHead("originManager", "{\"TraceBackServer\":\"Messy Request\"}", http.StatusNoContent, w)
	reqStr := dealWithReq(r)
	if reqStr != "nil" {
		req := make(map[string]string)
		errParse := json.Unmarshal([]byte(reqStr), &req)
		checkWithWarn(errParse)
		if req["order"] == "origins" {
			os := findOrigins()
			if os != nil {
				var flat string
				for i := range os.originName {
					flat = flat + os.originName[i] + "_"
				}
				if flat[len(flat)-1:] == "_" {
					flat = flat[0 : len(flat)-1]
				}
				sendBackCommonHead("{\"TraceBackServer\":\"I_AM_A_POLAR_BEAR\", \"origins\":\""+flat+"\"}", http.StatusOK, w)
			}
		}
	}
}

func chainToolAuthentication(auth *ChainToolUserAuthMask) string {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("chainToolAuthentication-handler-Error: ", err)
		}
	}()
	fmt.Println("auth::" + auth.UserName + "_" + auth.Password + "_" + auth.Origin)
	s := dataStore.session.Copy()
	coll := s.DB(mongod_truck_db).C(mongod_coll_name_chaintool_auth)
	yes := ChainToolUserAuth{}
	err := coll.Find(bson.M{"origin": auth.Origin, "userName": auth.UserName, "password": auth.Password}).One(&yes)
	if err == nil && len(yes.UserName) > 0 {
		withLogin := ChainToolUserAuth{}
		withLogin.LoginTime = time.Now().UnixNano() / int64(time.Millisecond)
		withLogin.CreateTime = yes.CreateTime
		withLogin.Origin = yes.Origin
		withLogin.Password = yes.Password
		withLogin.SerialNumber = yes.SerialNumber
		withLogin.UserName = yes.UserName
		err = coll.Update(yes, withLogin)
		checkWithWarn(err)
		return "{\"auth\": \"I_AM_A_POLAR_BEAR\"}"
	}
	return "{\"auth\": \"HANABI_ENCORE\"}"
}

func dealWithReq(r *http.Request) string {
	bodyData, errBody := ioutil.ReadAll(r.Body)
	bodyDataContext := string(bodyData[:])
	fmt.Println("POSTDATA=" + bodyDataContext)
	if errBody != nil || bodyDataContext == "null" || bodyDataContext == "" {
		fmt.Println("Error: ", errBody)
		return "Nil"
	}
	return bodyDataContext
}

/*  OPTIONS on routing::
r := mux.NewRouter()
r.Host("{subdomain}.domain.com").Path("/articles/{category}/{id:[0-9]+}").Queries("filter", "{filter}").HandlerFunc(ArticleHandler).Name("article")
// url.String() will be "http://news.domain.com/articles/technology/42?filter=gorilla"
url, err := r.Get("article").URL("subdomain", "news", "category", "technology", "id", "42", "filter", "gorilla")
*/
