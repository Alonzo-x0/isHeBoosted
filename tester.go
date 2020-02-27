package main

import (
	"net/http"
	"fmt"
	"log"
	"io/ioutil"
	"encoding/json"
)

type Summoner struct {
	//id = Encrypted summoner ID
	//accountID = Encrypted account ID
	Id string
	accountId string
	puuid string
	Name string
	profileIconId string
	revisionDate string
	summonerLevel string
}

func getJson(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	parseJson := string(body)

	fmt.Println(parseJson)
	
	var champ Summoner
	awef := json.Unmarshal(body, &champ)
	if awef != nil {
		log.Println(awef)
	}
	fmt.Println(champ.Name)
}
	

func main() {
	
	//summoner := Summoner{}
	getJson("https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/incompatible?api_key=RGAPI-fasef-070944a6a31f")
	//println(summoner.Id)
}
