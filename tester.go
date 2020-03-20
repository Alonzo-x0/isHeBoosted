package main

import (
	"net/http"
	"fmt"
	"log"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"github.com/joho/godotenv"
	"os"
)

type matchINFO struct {
	SeasonID              int   `json:"seasonId"`
	QueueID               int   `json:"queueId"`
	GameID                int64 `json:"gameId"`
	ParticipantIdentities []struct {
		Player struct {
			CurrentPlatformID string `json:"currentPlatformId"`
			SummonerName      string `json:"summonerName"`
			MatchHistoryURI   string `json:"matchHistoryUri"`
			PlatformID        string `json:"platformId"`
			CurrentAccountID  string `json:"currentAccountId"`
			ProfileIcon       int    `json:"profileIcon"`
			SummonerID        string `json:"summonerId"`
			AccountID         string `json:"accountId"`
		} `json:"player"`
		ParticipantID int `json:"participantId"`
	} `json:"participantIdentities"`
	GameVersion string `json:"gameVersion"`
	PlatformID  string `json:"platformId"`
	GameMode    string `json:"gameMode"`
	MapID       int    `json:"mapId"`
	GameType    string `json:"gameType"`
	Teams       []struct {
		FirstDragon bool `json:"firstDragon"`
		Bans        []struct {
			PickTurn   int `json:"pickTurn"`
			ChampionID int `json:"championId"`
		} `json:"bans"`
		FirstInhibitor       bool   `json:"firstInhibitor"`
		Win                  string `json:"win"`
		FirstRiftHerald      bool   `json:"firstRiftHerald"`
		FirstBaron           bool   `json:"firstBaron"`
		BaronKills           int    `json:"baronKills"`
		RiftHeraldKills      int    `json:"riftHeraldKills"`
		FirstBlood           bool   `json:"firstBlood"`
		TeamID               int    `json:"teamId"`
		FirstTower           bool   `json:"firstTower"`
		VilemawKills         int    `json:"vilemawKills"`
		InhibitorKills       int    `json:"inhibitorKills"`
		TowerKills           int    `json:"towerKills"`
		DominionVictoryScore int    `json:"dominionVictoryScore"`
		DragonKills          int    `json:"dragonKills"`
	} `json:"teams"`
	GameDuration int   `json:"gameDuration"`
	GameCreation int64 `json:"gameCreation"`
}


type Summoner struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	Puuid         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

	
type matchHistory struct {
	Matches []struct {
		Lane       string `json:"lane"`
		GameID     int64  `json:"gameId"`
		Champion   int    `json:"champion"`
		PlatformID string `json:"platformId"`
		Timestamp  int64  `json:"timestamp"`
		Queue      int    `json:"queue"`
		Role       string `json:"role"`
		Season     int    `json:"season"`
	} `json:"matches"`
	EndIndex   int `json:"endIndex"`
	StartIndex int `json:"startIndex"`
	TotalGames int `json:"totalGames"`
}

	
type errSpec struct {
	Status struct {
		Message    string `json:"message"`
		StatusCode int    `json:"status_code"`
	} `json:"status"`
}

	
func urlRequest(url string) []byte{
	//simple request function returns response body
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return body
}
func dumpMap(space string, m map[string]interface{}) {
	for k, v := range m {
		if mv, ok := v.(map[string]interface{}); ok {
			fmt.Printf("{ \"%v\": \n", k)
			dumpMap(space+"\t", mv)
			fmt.Printf("}\n")
		} else {
			fmt.Printf("%v %v : %v\n", space, k, v)
		}
	}
}


func getSummoner(name string, key string) (string, string){
	//id=encryptedID accountid=encryptedaccountid
	var url = "https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/" + name + "?api_key=" + key
	var summoner Summoner
	//jsonMap := make(map[string]interface{})
	err := json.Unmarshal(urlRequest(url), &summoner)
	if err != nil {
		log.Println(err)
	}
	var enID = summoner.ID
	var accID = summoner.AccountID
	return enID, accID
}


func freqCount(list []string) map[string]int{

	freqBind := make(map[string]int)

	for _, item := range list {

		_, exist := freqBind[item]

		if exist {
			freqBind[item] += 1
		}else {
			freqBind[item] = 1
		}
	}
	return freqBind
}



func winRatio(accID string, key string) []int64{
	var IDs []int64
	var ratio []string
	var url = "https://na1.api.riotgames.com/lol/match/v4/matchlists/by-account/" + accID + "?queue=420&endIndex=20&api_key=" + key
	var history matchHistory

	err := json.Unmarshal(urlRequest(url), &history)
	if err != nil {
		log.Println(err)
	}
	for i:=0; i <=19; i++ {
		IDs = append(IDs, history.Matches[i].GameID)		
	}
	for _, killme := range IDs {
		ratio = append(ratio, checkMatch(killme, key))
	}
	
	fmt.Println(freqCount(ratio))
	for cat, quant := range freqCount(ratio) {
		fmt.Println(cat, quant)
	}

	return IDs
	}

func checkMatch(gameId int64, key string) string{
	var url = "https://na1.api.riotgames.com/lol/match/v4/matches/" + strconv.FormatInt(gameId, 10) + "?api_key=" + key
	var matchinfo matchINFO
	err := json.Unmarshal(urlRequest(url), &matchinfo)
	if err != nil {
		log.Println(err)
	}


	return matchinfo.Teams[0].Win
}



func getMatchID(enID string, key string) (string) {
	var url = "https://na1.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/" + enID + "?api_key=" + key
	
	response := make(map[string]interface{})
	err := json.Unmarshal(urlRequest(url), &response)
	if err != nil {
		log.Println(err)
	}

	if response["status"] != nil{
		
		return response["status"].(map[string]interface{})["message"].(string)
	} else {
		
		return strconv.FormatFloat(response["gameId"].(float64), 'f', -1, 64)
	return "you should not be seeing this."
	} 
}


func main() {
	
	//summoner := Summoner{}
	getJson("https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/incompatible?api_key=RGAPI-fasef-070944a6a31f")
	//println(summoner.Id)
}

	
