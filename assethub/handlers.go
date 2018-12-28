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

func InitDataStoreHandlers() {
	dataStore.session = connectToDb()
}

func HandlerForFuckers(path string) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("handlerForFuckers-handler-Error: ", err)
		}
	}()
	if path != "" {
		path = "/" + path
	}
	fmt.Println("default_server_port=" + default_server_port)
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/getme/fuckers"+path, Index)
	r.HandleFunc(path+"/novice/{reqid}/{reqname}", Novice)
	r.HandleFunc(path+asset_saving_path, receiverOfAssets).Methods("POST")
	r.HandleFunc(path+human_saving_path, receiverOfHumans).Methods("POST")
	r.HandleFunc(path+binding_check_path, receiveOfBindingBoundagePool).Methods("POST")
	r.HandleFunc(path+chain_tool_user, chainToolUserHandler).Methods("POST")
	r.HandleFunc(path+chain_tool_user_human, chainToolUserGetHumanObjectHandler).Methods("POST")
	r.HandleFunc(path+origin_manager, originManager).Methods("POST")
	r.HandleFunc(path+humpers_jumpers, HumpersJumpers)
	r.HandleFunc(path+recordqr, ReportToSatan).Methods("POST")
	log.Fatal(http.ListenAndServe(default_server_port, r))
}

func receiverOfAssets(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("receiverOfAssets-handler-Error: ", err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Write([]byte("{\"TraceBackServer\":\"SomethingWrongWithTheData\"}"))
			w.WriteHeader(http.StatusNoContent)
		}
	}()
	authToken := r.Header.Get(asset_auth)
	if len(authToken) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{\"TraceBackServer\":\"You r not authorized. please get token from novice req\"}"))
		w.WriteHeader(http.StatusForbidden)
		return
	}
	fmt.Println("#AUTHENTICATION_ASSET_INSERT::" + authToken)
	bodyData, errBody := ioutil.ReadAll(r.Body)
	userPass := r.Header.Get(header_identity)
	bodyDataContext := string(bodyData[:])
	fmt.Println("POSTDATA=" + bodyDataContext)
	if errBody != nil || bodyDataContext == "null" {
		fmt.Println("Error: ", errBody)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{\"TraceBackServer\":\"YourRequestIsScrambled\"}"))
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	checkToken := tokenGenerator(strings.Split(userPass, "@")[0], strings.Split(userPass, "@")[1])
	if len(checkToken) != len(authToken) || checkToken != authToken {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{\"TraceBackServer\":\"YouAreAnAlienYouAreNotAllowed\"}"))
		w.WriteHeader(http.StatusForbidden)
		return
	}
	assetsHandler(bodyDataContext, dataStore.session.Copy())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set(asset_auth, authToken)
	iAmAPollarBear := "{\"TraceBackServer\":\"I_AM_A_POLAR_BEAR\"}"
	w.Write([]byte(iAmAPollarBear))
	w.WriteHeader(http.StatusOK)
}

func receiverOfHumans(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("receiverOfHumans-handler-Error: ", err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Write([]byte("{\"TraceBackServer\":\"SomethingWrongWithTheData\"}"))
			w.WriteHeader(http.StatusNoContent)
		}
	}()
	authToken := r.Header.Get(asset_auth)
	if len(authToken) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{\"TraceBackServer\":\"You r not authorized. please get token from novice req\"}"))
		w.WriteHeader(http.StatusForbidden)
		return
	}
	userPass := r.Header.Get(header_identity) //this is in the form of userId@username
	fmt.Println("#AUTHENTICATION_ASSET_INSERT::" + authToken)
	bodyData, errBody := ioutil.ReadAll(r.Body)
	bodyDataContext := string(bodyData[:])
	fmt.Println("POSTDATA=" + bodyDataContext)
	if errBody != nil || bodyDataContext == "null" {
		fmt.Println("Error: ", errBody)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{\"TraceBackServer\":\"YourRequestIsScrambled\"}"))
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	checkToken := tokenGenerator(strings.Split(userPass, "@")[0], strings.Split(userPass, "@")[1])
	if len(checkToken) != len(authToken) || checkToken != authToken {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{\"TraceBackServer\":\"YouAreAnAlienYouAreNotAllowed\"}"))
		w.WriteHeader(http.StatusForbidden)
		return
	}
	humansHandler(bodyDataContext, dataStore.session.Copy())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set(asset_auth, authToken)
	iAmAPollarBear := "{\"TraceBackServer\":\"I_AM_A_POLAR_BEAR\"}"
	w.Write([]byte(iAmAPollarBear))
	w.WriteHeader(http.StatusOK)
}

func receiveOfBindingBoundagePool(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("receiveOfBindingBoundagePool-handler-Error: ", err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Write([]byte("{\"TraceBackServer\":\"SomethingWrongWithTheData\"}"))
			w.WriteHeader(http.StatusNoContent)
		}
	}()
	authToken := r.Header.Get(asset_auth)
	if len(authToken) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{\"TraceBackServer\":\"You r not authorized. please get token from novice req\"}"))
		w.WriteHeader(http.StatusForbidden)
		return
	}
	userPass := r.Header.Get(header_identity) //this is in the form of userId@username
	fmt.Println("#AUTHENTICATION_ASSET_INSERT::" + authToken)
	bodyData, errBody := ioutil.ReadAll(r.Body)
	bodyDataContext := string(bodyData[:])
	fmt.Println("POSTDATA=" + bodyDataContext)
	if errBody != nil || len(bodyDataContext) <= 0 {
		fmt.Println("Error: ", errBody)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{\"TraceBackServer\":\"YourRequestIsScrambled\"}"))
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	checkToken := tokenGenerator(strings.Split(userPass, "@")[0], strings.Split(userPass, "@")[1])
	if len(checkToken) != len(authToken) || checkToken != authToken {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{\"TraceBackServer\":\"YouAreAnAlienYouAreNotAllowed\"}"))
		w.WriteHeader(http.StatusForbidden)
		return
	}
	res := bindingHandler(bodyDataContext, dataStore.session.Copy())
	if res != "done" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{\"TraceBackServer\":\"YourRequestIsScrambled\"}"))
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set(asset_auth, authToken)
	iAmAPollarBear := "{\"TraceBackServer\":\"I_AM_A_POLAR_BEAR\"}"
	w.Write([]byte(iAmAPollarBear))
	w.WriteHeader(http.StatusOK)
}

func chainToolUserHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("chainToolUserHandler-handler-Error: ", err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Write([]byte("{\"TraceBackServer\":\"SomethingWrongWithTheData\"}"))
			w.WriteHeader(http.StatusNoContent)
		}
	}()
	bodyData, errBody := ioutil.ReadAll(r.Body)
	bodyDataContext := string(bodyData[:])
	fmt.Println("POSTDATA=" + bodyDataContext)
	if errBody != nil || bodyDataContext == "null" || bodyDataContext == "" {
		fmt.Println("Error: ", errBody)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{\"TraceBackServer\":\"YourRequestIsScrambled\"}"))
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	user := ChainToolUserAuthMask{}
	errParse := json.Unmarshal([]byte(bodyDataContext), &user)
	checkWithWarn(errParse)
	resAuth := chainToolAuthentication(&user)
	if resAuth == "{\"auth\": \"I_AM_A_POLAR_BEAR\"}" {
		genToken := tokenGeneratorTimely(user.Origin, user.UserName+"@"+user.Password)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		iAmAPollarBear := "{\"TraceBackServer\":\"I_AM_A_POLAR_BEAR\", \"token\":\"" + genToken + "\"}"
		w.Write([]byte(iAmAPollarBear))
		w.WriteHeader(http.StatusOK)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte(resAuth))
		w.WriteHeader(http.StatusBadRequest)
	}
}

func chainToolUserGetHumanObjectHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("chainToolUserGetHumanObjectHandler-handler-Error: ", err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Write([]byte("{\"TraceBackServer\":\"cannot get the human info. due to req info mistakes maybe.\"}"))
			w.WriteHeader(http.StatusNoContent)
		}
	}()
	bodyData, errBody := ioutil.ReadAll(r.Body)
	bodyDataContext := string(bodyData[:])
	fmt.Println("POSTDATA=" + bodyDataContext)
	if errBody != nil || bodyDataContext == "null" || bodyDataContext == "" {
		fmt.Println("bodyData Error: ", errBody)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{\"TraceBackServer\":\"request not found\"}"))
		w.WriteHeader(http.StatusNotAcceptable)
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
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{\"TraceBackServer\":\"requesting user not found\"}"))
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	collHuman := s.DB(mongod_truck_db).C(mongod_coll_name_humans)
	humanFound := HumansContext{}
	errHumanFound := collHuman.Find(bson.M{"humanSerialNumber": userAuth.SerialNumber}).One(&humanFound)
	if errHumanFound != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{\"TraceBackServer\":\"no basic info found\"}"))
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	b, errByte := json.Marshal(humanFound)
	if errByte != nil {
		panic(errByte)
	}
	fmt.Println("BACKING==" + string(b[:]))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	iAmAPollarBear := "{\"TraceBackServer\":\"I_AM_A_POLAR_BEAR\", \"humanContext\":" + string(b[:]) + "}"
	w.Write([]byte(iAmAPollarBear))
	w.WriteHeader(http.StatusOK)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("version of current:" + version_name_of_app)
	fmt.Fprintln(w, version_name_of_app, html.EscapeString(" AT:"+time.Now().String()))
}

func Novice(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("Novice-handler-Error: ", err)
		}
	}()
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
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		iAmAPollarBear := "{\"TraceBackServer\": \"" + string(b[:]) + " AT: " + time.Now().String() + "\"}"
		w.Write([]byte(iAmAPollarBear))
		w.WriteHeader(http.StatusOK)
	} else {
		token := fetchForAppUser(auth.AuthId, auth.AuthName)
		if len(token) == 0 {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			youAreAWitch := "{\"TraceBackServer\":\"YOU_ARE_A_WITCH_-_NOT_ALLOWED_TO_CAST\"}"
			w.Write([]byte(youAreAWitch))
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Header().Set("auth", token)
			iAmAPollarBear := "{\"TraceBackServer\":\"I_AM_A_POLAR_BEAR\", \"token\":\"" + token + "\"}"
			w.Write([]byte(iAmAPollarBear))
			w.WriteHeader(http.StatusOK)
		}
	}
}

func originManager(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("originManager-handler-Error: ", err)
		}
	}()
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
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				iAmAPollarBear := "{\"TraceBackServer\":\"I_AM_A_POLAR_BEAR\", \"origins\":\"" + flat + "\"}"
				w.Write([]byte(iAmAPollarBear))
				w.WriteHeader(http.StatusOK)
			}
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

/*
r := mux.NewRouter()
r.Host("{subdomain}.domain.com").
  Path("/articles/{category}/{id:[0-9]+}").
  Queries("filter", "{filter}").
  HandlerFunc(ArticleHandler).
  Name("article")

// url.String() will be "http://news.domain.com/articles/technology/42?filter=gorilla"
url, err := r.Get("article").URL("subdomain", "news",
                                 "category", "technology",
                                 "id", "42",
                                 "filter", "gorilla")
*/
