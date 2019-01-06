package assethub

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
	mgo "gopkg.in/mgo.v2"
)

const (
	version_name_of_app = "crazy.app._ver=0.0.1"

	WORDS_OF_CHOICE     = "TheFashionableWorldDismayedByTheMurderOfTheHonourableRonaldAdair"
	default_server_port = ":8093"

	mongod_main_one = "192.168.204.147:27017" //this one is for my local
	// mongod_main_one = "127.0.0.1:27017" //this one is for alpha version

	mongod_main_db                  = "go_crazy_lemons"
	mongod_coll_name_headinfo       = "headinfos"
	mongod_coll_name_authrec        = "auth_rec"
	mongod_coll_name_generic_data   = "generic_data"
	mongod_coll_name_chaintool_auth = "chaintool_users"

	mongod_truck_db                     = "truck-lift-forks"
	mongod_coll_name_host_pool          = "Host_Pool"
	mongod_coll_name_host_pool_rec      = "host_pool"
	mongod_coll_name_pool_of_host_final = "Host_Pool_Final"
	mongod_coll_name_key_chain          = "KeyChain"
	mongod_coll_name_assets             = "Asset_Slices"
	mongod_coll_name_humans             = "Human_Slices"
	mongod_coll_name_bindingPool        = "BindingBoundagePool"
	mongod_coll_name_qrgen_rec          = "QRGen_Record"

	redis_host = "192.168.204.147" //this one is for local
	// redis_host = "127.0.0.1" //this one is for alpha version
	redis_port = 6379
	redis_db   = 0

	data_insert_ok = "inserted"

	//ROUTES FOR data_entry.main::
	data_entry_main = "/products-that-you-have/entry/"

	basic_data_request_name = "basicdata"

	asset_auth            = "assetknock"
	header_identity       = "identity"
	asset_saving_path     = "/assetent/"
	asset_remove_path     = "/removeasset/"
	human_saving_path     = "/humanent/"
	human_remove_path     = "/removehuman/"
	binding_check_path    = "/bc/"
	chain_tool_user       = "/ctv/"
	chain_tool_user_human = "/ctvh/"
	origin_manager        = "/ori/"

	//multitrack
	humpers_jumpers = "/h1/pid={pid}"
	recordqr        = "/qrrec/"

	//for calculation of host grouping::
	pool_max = 20

	ASSET_TO_BASE = "#ASSET-DATA@POSTGRES"
	HUMAN_TO_BASE = "#HUMAN-DATA@POSTGRES"
	ASSET_REMOVE  = "#ASSET-MINUS@POSTGRES"
	HUMAN_REMOVE  = "#HUMAN-MINUS@POSTGRES"
)

func uuidGenAlt() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	fmt.Println(uuid)
	return uuid
}

func uuidGen() string {
	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}
	return u.String()
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

func Check(err error) {
	if err != nil {
		fmt.Println("==not finished with panic==")
		panic(err)
	}
}

func checkWithWarn(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

type DataStore struct {
	session *mgo.Session
}

type Cache struct {
	redisSession *redis.Client
}

/*
	for user entrance.
*/
type ConfigAuth struct {
	AuthId   string `bson:"authId" json:"authId"`
	AuthName string `bson:"authName" json:"authName"`
}

/*
	for chaintool user authentication only.
*/
type ChainToolUserAuth struct {
	Origin       string `bson:"origin" json:"origin"`
	UserName     string `bson:"userName" json:"userName"`
	Password     string `bson:"password" json:"password"`
	LoginTime    int64  `bson:"loginTime" json:"loginTime"`
	CreateTime   int64  `bson:"createTime" json:"createTime"`
	SerialNumber string `bson:"serialNumber" json:"serialNumber"`
}

/*
	for http json on chaintool user auth.
*/
type ChainToolUserAuthMask struct {
	Origin   string `bson:"origin" json:"origin"`
	UserName string `bson:"userName" json:"userName"`
	Password string `bson:"password" json:"password"`
}

/*
	headcheck for admin auth users.
*/
type InteractInfo struct {
	AuthName        string            `json:"authName"`
	AuthNick        string            `json:"authNick"`
	PublicInterface []string          `json:"publicInterface"`
	HostPools       map[string]string `json:"hostPools"`
}
type InteractInfos []InteractInfo

type ResultOfPersist struct {
	ResultOfPers string `json:"resultOfPers"`
	TimeLog      string `json:"timeLog"`
}

type DailyChainKey struct {
	Date string `bson:"date" json:"date"`
	Key  string `bson:"key" json:"key"`
}

/**************
mgo connection
***************/
func MongoClient() *mgo.Session {
	session, err := mgo.Dial(mongod_main_one)
	if err != nil {
		panic(err)
	}
	// defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	return session
}

/**************
redis client
***************/
func RedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redis_host + ":" + strconv.Itoa(redis_port),
		Password: "",
		DB:       redis_db,
	})
	pong, err := client.Ping().Result()
	fmt.Println("*******************************REDIS CONN CHECK*******************************")
	fmt.Printf("=redis=%s == =redis=%s\n", pong, err)
	if pong == "PONG" {
		fmt.Println("REDIS CONNECTION DONE.")
	}
	fmt.Println("*******************************REDIS CONN CHECK*******************************")
	return client
}

func hashing(rtn string) string {
	hasher := md5.New()
	io.WriteString(hasher, rtn)
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash
}

func tokenGenerator(idStr string, nameStr string) string {
	dailyKeyChainStr := dailyKeyChain(idStr)
	dailyHash := hashing(dailyKeyChainStr)
	idHash := hashing(idStr)
	var decorateName string
	if len(nameStr) > 0 {
		decorateName = hashing("#$" + nameStr + "$#")
	} else {
		decorateName = hashing("NO_NAME")
	}
	return decorateName + dailyHash + idHash
}

func tokenGeneratorTimely(idStr string, nameStr string) string {
	dailyKeyChainStr := dailyKeyChain(idStr)
	dailyHash := hashing(dailyKeyChainStr)
	idHash := hashing(idStr)
	timer := time.Now().UnixNano() / int64(time.Millisecond)
	var decorateName string
	if len(nameStr) > 0 {
		decorateName = hashing("#$" + nameStr + "$#" + strconv.Itoa(int(timer)))
	} else {
		decorateName = hashing("NO_NAME" + strconv.Itoa(int(timer)))
	}
	return decorateName + dailyHash + idHash
}

func dailyKeyChain(ref string) string {
	tm := time.Now().Format(time.RFC3339)
	fmt.Println("TIMEFORCHAINKEY==" + tm)
	tm = strings.Split(tm, "T")[0]
	upperpart := ""
	judge, _ := strconv.Atoi(tm[9:])
	if judge >= 0 && judge < 3 {
		upperpart += "cadaver"
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
	if len(ref) > 0 {
		upperpart += ref
	}
	return upperpart
}

func currentMilliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func CurrentMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func errorFeedbackCommon(err error, w http.ResponseWriter, answer string) bool {
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte("{\"TraceBackServer\":\"" + answer + "\"}"))
		w.WriteHeader(http.StatusNotAcceptable)
		return true
	}
	return false
}

func commonWriteBack(w http.ResponseWriter, returns []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(returns)
	w.WriteHeader(http.StatusAccepted)
}

func deferErrorFeedbackRest(w http.ResponseWriter, returns []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(returns)
	w.WriteHeader(http.StatusNoContent)
}
