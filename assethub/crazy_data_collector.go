package assethub

import (
	"encoding/json"
	"fmt"
	"strconv"

	bytes "bytes"
	"time"

	_redis "github.com/go-redis/redis"
)

//********************************************* ASSETS ************************************************//

/***********************************************************
main generic type as data container
display on ui is customized, unlike the 'info' parameter
which is display manditory in ui, formatted in json
************************************************************/
type GenericContext struct {
	ParamName  string `json:"paramName" bson:"paramName"`
	ParamType  string `json:"paramType" bson:"paramType"`
	ParamValue string `json:"paramValue" bson:"paramValue"`
	ShowOrNot  int    `json:"showOrNot" bson:"showOrNot"`
}

/*************************************************************
main trust bind struct. only trust bind struct actually.
vcstatus = 0, 1, -1 -> 0=ended 1=available -1=pending
time_left=endtime-starttime
cost_used=cost * cost unit (per hour or min or day or month or year)
SubContract limit to two.
VCAvail = 0:cancelled 1:available
**************************************************************/
type VirtualContract struct {
	VCSerialNumber      string             `json:"VCSerialNumber" bson:"VCSerialNumber"`
	VCHash              string             `json:"VCHash" bson:"VCHash"`
	LedgerHash          string             `json:"ledgerHash" bson:"ledgerHash"`
	VCType              string             `json:"VCType" bson:"VCType"`
	VCStartTime         int64              `json:"VCStartTime" bson:"VCStartTime"`
	VCEndTime           int64              `json:"VCEndTime" bson:"VCEndTime"`
	VCStatus            int                `json:"VCStatus" bson:"VCStatus"`
	VCCost              float32            `json:"VCCost" bson:"VCCost"`
	VCCostUnit          float32            `json:"VCCostUnit" bson:"VCCostUnit"`
	VCPrice             float32            `json:"VCPrice" bson:"VCPrice"`
	Info                string             `json:"info" bson:"info"`
	TimeOfSolidifaction int64              `json:"timeOfSolidification" bson:"timeOfSolidification"`
	CustomContext       []GenericContext   `json:"customContext" bson:"customContext"`
	SubContract         []VirtualContract  `json:"subContract" bson:"subContract"`
	VerificationInfos   []VerificationInfo `json:"verificationInfos" bson:"verificationInfos"`
}

type Calculator struct {
	CalculatorName      string            `json:"calculatorName" bson:"calculatorName"`
	Origin              string            `json:"origin" bson:"origin"`
	CalculationFormulas map[string]string `json:"calculationFormulas" bson:"calculationFormulas"`
}

/**************************************************************************************************************************
asset object. for everthing
AssetCategoryName = big type name of assets
AssetMainType = level of assets. as first level is the root,
				then inner asset level, and inner assets can have one more level of inner assets.
InnerAssets limit to 200.

struct assetAuth is for shortcut get auth token.
***************************************************************************************************************************/
type AssetObject struct {
	SerialNumber      string           `json:"serialNumber" bson:"serialNumber"`
	Origin            string           `json:"origin" bson:"origin"`
	AssetGroup        string           `json:"assetGroup" bson:"assetGroup"`
	AssetName         string           `json:"assetName" bson:"assetName"`
	AssetHash         string           `json:"assetHash" bson:"assetHash"`
	AssetCategoryName string           `json:"assetCategoryName" bson:"assetCategoryName"`
	AssetMainType     string           `json:"assetMainType" bson:"assetMainType"`
	Info              string           `json:"info" bson:"info"`
	CustomContext     []GenericContext `json:"customContext" bson:"customContext"`
	AssetProperties   AssetProperties  `json:"assetProperties" bson:"assetProperties"`
	InnerAssets       []AssetObject    `json:"innerAssets" bson:"innerAssets"`
	VirtualContract   VirtualContract  `json:"virtualContract" bson:"virtualContract"`
}

/***************************************************************
assetType will be in format of 111111 11111111 11111111 = 3F FF FF
3F being the upmost level of the categories. then FF being second, then third.
if assetBirth only have year, just make it like yyyy-00-00
date: yyyy-MM-dd  time: HH:mm:ss
AssetWorkDuration: in milliseconds log the device time-ons.
**************************************************************/
type AssetProperties struct {
	AssetType         string           `json:"assetType" bson:"assetType"`
	AssetBirthDate    string           `json:"assetBirthDate" bson:"assetBirthDate"`
	AssetBirthTime    string           `json:"assetBirthTime" bson:"assetBirthTime"`
	AssetStatus       string           `json:"assetStatus" bson:"assetStatus"`
	AssetApplyStatus  string           `json:"assetApplyStatus" bson:"assetApplyStatus"`
	AssetReleaseTime  int64            `json:"assetReleaseTime" bson:"assetReleaseTime"`
	AssetReleaseType  string           `json:"assetReleaseType" bson:"assetReleaseType"`
	AssetWorkDuration int64            `json:"assetWorkDuration" bson:"assetWorkDuration"`
	AssetApplication  AssetApplication `json:"assetApplication" bson:"assetApplication"`
}

/***************************************************************************************************************
any other applicable parameters e.g. speed of a car would be carried in Info and CustomContext
ledgerHash is calculated with messageDigest and is the combination of AssetHash and virtualContractHash
****************************************************************************************************************/
type AssetApplication struct {
	AssetLocation     AssetLocationInfo `json:"assetLocation" bson:"assetLocation"`
	AssetVerification VerificationInfo  `json:"assetVerification" bson:"assetVerification"`
	Info              string            `json:"info" bson:"info"`
	CustomContext     []GenericContext  `json:"customContext" bson:"customContext"`
	LedgerHash        string            `json:"ledgerHash" bson:"ledgerHash"`
}

type AssetLocationInfo struct {
	Lon      float64 `json:"lon" bson:"lon"`
	Lat      float64 `json:"lat" bson:"lat"`
	Loc      string  `json:"loc" bson:"loc"`
	TimeLast int64   `json:"timeLast" bson:"timeLast"`
}

type VerificationInfo struct {
	LastVerifiedBy      string `json:"lastVerifiedBy" bson:"lastVerifiedBy"`
	LastVerifiedTime    int64  `json:"lastVerifiedTime" bson:"lastVerifiedTime"`
	VerificationResult  string `json:"verificationResult" bson:"verificationResult"`
	VerificationContent string `json:"verificationContent" bson:"verificationContent"`
}

//********************************************* HUMANS ************************************************//

/**************************************************************************
humans
passwd being the only auth
deadOrAlive = if this human still available for the platform
canMaintain : 0=can't/user  1=can/engineer  99=can/admin
AssetAccess : the assetGroup that this human can access
HumanTrust : if score less than gauge, trust is broken
***************************************************************************/
type HumansContext struct {
	HumanSerialNumber string             `json:"humanSerialNumber" bson:"humanSerialNumber"`
	HumanHash         string             `json:"humanHash" bson:"humanHash"`
	Origin            string             `json:"origin" bson:"origin"`
	HumanName         string             `json:"humanName" bson:"humanName"`
	HumanGroup        string             `json:"humanGroup" bson:"humanGroup"`
	Passwd            string             `json:"passwd" bson:"passwd"`
	Info              string             `json:"info" bson:"info"`
	CustomContext     []GenericContext   `json:"customContext" bson:"customContext"`
	CanMaintain       int                `json:"canMaintain" bson:"canMaintain"`
	VirtualContracts  []HumansContext    `json:"virtualContracts" bson:"virtualContracts"`
	DeadOrAlive       int                `json:"deadOrAlive" bson:"deadOrAlive"`
	BirthTime         int64              `json:"birthTime" bson:"birthTime"`
	AssetAccess       []string           `json:"assetAccess" bson:"assetAccess"`
	HumanTrust        HumansTrustContext `json:"humanTrust" bson:"humanTrust"`
	VerificationInfo  VerificationInfo   `json:"verificationInfo" bson:"verificationInfo"`
}

type HumansTrustContext struct {
	Info              string           `json:"info" bson:"info"`
	CustomContext     []GenericContext `json:"customContext" bson:"customContext"`
	VirtualTrustScore int              `json:"virtualTrustScore" bson:"virtualTrustScore"`
	VirtualTrustGauge int              `json:"virtualTrustGauge" bson:"virtualTrustGauge"`
}

type BindContents struct {
	Name       string       `json:"name" bson:"name"`
	BindValues []BindValues `json:"bindValues" bson:"bindValues"`
	BindFlag   string       `json:"bindFlag" bson:"bindFlag"`
}

type BindValues struct {
	FieldName  string `json:"fieldName" bson:"fieldName"`
	FieldValue string `json:"fieldValue" bson:"fieldValue"`
}

type BindingBoundagePool struct {
	BindingSerial string         `json:"bindingSerial" bson:"bindingSerial"`
	Origin        string         `json:"origin" bson:"origin"`
	BindContent   []BindContents `json:"bindContent" bson:"bindContent"`
}

const (
	assetsReceiver          = "#-DATA_COLLECTOR_ASSETS"
	virtualContractReceiver = "#-DATA_COLLECTOR_VIRTUAL_CONTRACT"
	humansReceiver          = "#-DATA_COLLECTOR_HUMANS"
)

/**
to get the hash for the whole hash...
*/
func getHashAssetProperties(data *AssetProperties) string {
	var rtn bytes.Buffer
	itemList := make([]string, 0)
	itemList = append(itemList, data.AssetType)
	itemList = append(itemList, data.AssetBirthDate)
	itemList = append(itemList, data.AssetBirthTime)
	itemList = append(itemList, data.AssetStatus)
	itemList = append(itemList, data.AssetApplyStatus)
	itemList = append(itemList, strconv.Itoa(int(data.AssetReleaseTime)))
	itemList = append(itemList, data.AssetReleaseType)
	itemList = append(itemList, strconv.Itoa(int(data.AssetWorkDuration)))
	assetApplicationHash := getHashAssetApplication(&data.AssetApplication)
	itemList = append(itemList, assetApplicationHash)
	for i := range itemList {
		rtn.WriteString(itemList[i])
	}
	hash := hashing(rtn.String())
	return hash
}

/**
to get the hash for the whole hash...
*/
func getHashAssetApplication(data *AssetApplication) string {
	var rtn bytes.Buffer
	itemList := make([]string, 0)
	itemList = append(itemList, strconv.FormatFloat(data.AssetLocation.Lon, 'f', -1, 64))
	itemList = append(itemList, strconv.FormatFloat(data.AssetLocation.Lat, 'f', -1, 64))
	itemList = append(itemList, data.AssetLocation.Loc)
	itemList = append(itemList, strconv.Itoa(int(data.AssetLocation.TimeLast)))
	itemList = append(itemList, data.AssetVerification.LastVerifiedBy)
	itemList = append(itemList, strconv.Itoa(int(data.AssetVerification.LastVerifiedTime)))
	itemList = append(itemList, data.AssetVerification.VerificationResult)
	itemList = append(itemList, data.AssetVerification.VerificationContent)
	for i := range itemList {
		rtn.WriteString(itemList[i])
	}
	hash := hashing(rtn.String())
	return hash
}

func getHashGenericContexts(customContexts []GenericContext) string {
	var rtn bytes.Buffer
	itemList := make([]string, 0)
	if customContexts != nil && len(customContexts) > 0 {
		for i := range customContexts {
			itemList = append(itemList, customContexts[i].ParamName)
			itemList = append(itemList, customContexts[i].ParamType)
			itemList = append(itemList, customContexts[i].ParamValue)
			itemList = append(itemList, strconv.Itoa(customContexts[i].ShowOrNot))
		}
		for i := range itemList {
			rtn.WriteString(itemList[i])
		}
	} else {
		return ""
	}
	hash := hashing(rtn.String())
	return hash
}

func getHashVirtualContract(vc *VirtualContract) string {
	var rtn bytes.Buffer
	itemList := make([]string, 0)
	if vc == nil {
		return ""
	}
	itemList = append(itemList, vc.VCSerialNumber)
	itemList = append(itemList, vc.VCType)
	itemList = append(itemList, strconv.Itoa(int(vc.VCStartTime)))
	itemList = append(itemList, strconv.Itoa(int(vc.VCEndTime)))
	itemList = append(itemList, strconv.Itoa(vc.VCStatus))
	itemList = append(itemList, strconv.Itoa(int(vc.VCCost)))
	itemList = append(itemList, strconv.Itoa(int(vc.VCCostUnit)))
	itemList = append(itemList, strconv.Itoa(int(vc.VCPrice)))
	itemList = append(itemList, strconv.Itoa(int(vc.TimeOfSolidifaction)))
	verificationInfos := vc.VerificationInfos
	if verificationInfos != nil && len(verificationInfos) > 0 {
		for i := range verificationInfos {
			itemList = append(itemList, getHashVerificationInfo(&verificationInfos[i]))
		}
	}
	for i := range itemList {
		rtn.WriteString(itemList[i])
	}
	hash := hashing(rtn.String())
	return hash
}

func getHashVerificationInfo(vi *VerificationInfo) string {
	var rtn bytes.Buffer
	itemList := make([]string, 0)
	itemList = append(itemList, vi.LastVerifiedBy)
	itemList = append(itemList, strconv.Itoa(int(vi.LastVerifiedTime)))
	itemList = append(itemList, vi.VerificationResult)
	itemList = append(itemList, vi.VerificationContent)
	for i := range itemList {
		rtn.WriteString(itemList[i])
	}
	hash := hashing(rtn.String())
	return hash
}

func getHashHumanTrust(ht *HumansTrustContext) string {
	var rtn bytes.Buffer
	itemList := make([]string, 0)
	itemList = append(itemList, getHashGenericContexts(ht.CustomContext))
	itemList = append(itemList, hashing(ht.Info))
	itemList = append(itemList, strconv.Itoa(ht.VirtualTrustGauge))
	itemList = append(itemList, strconv.Itoa(ht.VirtualTrustScore))
	for i := range itemList {
		rtn.WriteString(itemList[i])
	}
	hash := hashing(rtn.String())
	return hash
}

func fixOnVirtualContractSerialNumber() string {
	uuid := uuidGen()
	timeStamp := time.Now().String()
	sn := hashing(uuid + timeStamp)
	return sn
}

func fixOnAssetSerialNumber(origin string) string {
	uuid := uuidGen()
	timeStamp := time.Now().String()
	sn := hashing(origin + uuid + timeStamp)
	return sn
}

func fixOnHumanSerialNumber(origin string, humanName string) string {
	uuid := uuidGen()
	timeStamp := time.Now().String()
	sn := hashing(origin + uuid + timeStamp + "whenItIsHuman")
	return sn
}

//HumanReceiver get human thru redis. not useful for version 1.
func HumanReceiver(rd *_redis.Client) {
	for {
		msg := rd.LPop(humansReceiver)
		if msg != nil && msg.Val() != "" {
			fmt.Println("msg===" + msg.Val())
		}
		time.Sleep(5)
	}
}

//AssetReceiverRunner for getting data thru redis. not useful for version 1.
func AssetReceiverRunner(rd *_redis.Client) {
	psb := rd.Subscribe(assetsReceiver)
	defer psb.Close()
	for {
		asset, err := psb.Receive()
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch msg := asset.(type) {
		case *_redis.Message:
			assetContext := msg.Payload
			fmt.Println(assetContext)
			assetObj := AssetObject{}
			json.Unmarshal([]byte(assetContext), &assetObj)
		default:
			fmt.Println("==ctrl ms==")
		}
		time.Sleep(4)
	}
}

//VirtualContractReceiverRunner for getting data thru redis. not useful for version 1.
func VirtualContractReceiverRunner(rd *_redis.Client) {
	psb := rd.Subscribe(virtualContractReceiver)
	defer psb.Close()
	for {
		virtualContract, err := psb.Receive()
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch msg := virtualContract.(type) {
		case *_redis.Message:
			virtualContext := msg.Payload
			fmt.Println(virtualContext)
			virtualObj := VirtualContract{}
			json.Unmarshal([]byte(virtualContext), &virtualObj)
		default:
			fmt.Println("==control msg==")
		}
		time.Sleep(4)
	}
}

var bufForAssetSender = make([]string, 0)
var bufForHumanSender = make([]string, 0)
var bufForRemoveAsset = make([]string, 0)
var bufForRemoveHuman = make([]string, 0)

//AssetSender pushing to backend
func AssetSender(rd *_redis.Client) {
	for {
		if len(bufForAssetSender) > 0 {
			rd.LPush(ASSET_TO_BASE, bufForAssetSender[0])
			bufForAssetSender = bufForAssetSender[1:]
		}
		time.Sleep(time.Millisecond * 2)
	}
}

//HumanSender pushing to backend
func HumanSender(rd *_redis.Client) {
	for {
		if len(bufForHumanSender) > 0 {
			rd.LPush(HUMAN_TO_BASE, bufForHumanSender[0])
			bufForHumanSender = bufForHumanSender[1:]
		}
		time.Sleep(time.Millisecond * 2)
	}
}

//AssetRemoveSender pushing to backend
func AssetRemoveSender(rd *_redis.Client) {
	for {
		if len(bufForRemoveAsset) > 0 {
			rd.LPush(ASSET_REMOVE, bufForRemoveAsset[0])
			bufForRemoveAsset = bufForRemoveAsset[1:]
		}
		time.Sleep(time.Millisecond * 2)
	}
}

//HumanRemoveSender pushing to backend
func HumanRemoveSender(rd *_redis.Client) {
	for {
		if len(bufForRemoveHuman) > 0 {
			rd.LPush(HUMAN_REMOVE, bufForRemoveHuman[0])
			bufForRemoveHuman = bufForRemoveHuman[1:]
		}
		time.Sleep(time.Millisecond * 2)
	}
}
