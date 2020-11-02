package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Team represents a baseball team data.
type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MatchData represents baseball match data at a single point in time.
type MatchData struct {
	GameID      int  `json:"game_id"`
	TeamHome    Team `json:"team_home"`
	TeamAway    Team `json:"team_away"`
	ScoreHome   int  `json:"home_score"`
	ScoreAway   int  `json:"away_score"`
	TopInning   bool `json:"top_inning"`
	Out         int  `json:"out"`
	FirstBase   bool `json:"1st_base"`
	SecondBase  bool `json:"2nd_base"`
	ThirdBase   bool `json:"3rd_base"`
	InningCount int  `json:"inning_count"`
	PitcherID   int  `json:"pitcher_id"`
	BatterID    int  `json:"batter_id"`
}

func main() {
	matchData, err := fetchData()
	if err != nil {
		log.Fatalln("error fetching data:", err)
	}

	fmt.Printf("matchData: %+v\n", matchData)
}

func fetchData() (*MatchData, error) {
	res, err := http.Get("http://localhost:3000/data")
	if err != nil {
		return nil, err
	}

	var body []byte
	_, err = res.Body.Read(body)
	if err != nil {
		return nil, err
	}

	matchData := MatchData{}
	err = json.Unmarshal(body, &matchData)

	return &matchData, nil
}
