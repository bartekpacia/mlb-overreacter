package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
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

// Difference is the objectiv difference in score, in favor of one
// of the teams. It can also be neutral (i.e "there is no difference").
type Difference int

const (
	// Neutral means neither home nor away is happy.
	Neutral = 0

	// GoodForHome means that home is happy and away is sad.
	GoodForHome = -1
	// BadForAway means that home is happy and away is sad.
	BadForAway = -1

	// GoodForAway means that away is happy and home is sad.
	GoodForAway = 1
	// BadForHome means that away is happy and home is sad.
	BadForHome = 1
)

var (
	fetchURL    string
	interval    time.Duration
	currentData *MatchData
)

func init() {
	flag.StringVar(&fetchURL, "fetch URL", "http://localhost:3000/data", "endpoint to fetch baseball data from")
	flag.DurationVar(&interval, "interval", 3*time.Second, "interval every which new data will be fetched")
}

func main() {
	flag.Parse()
	currentData = &MatchData{}

	finished := make(chan struct{})
	go update(finished)

	<-finished
}

func update(finished chan struct{}) {

	for {
		fmt.Printf("-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-\n")

		newData, err := Fetch(fetchURL)
		if err != nil {
			fmt.Printf("error fetching data: %v\n", err)
			break
		}

		difference := Compare(newData, currentData)

		if difference == GoodForHome {
			fmt.Printf("%s: hooray!\n", newData.TeamHome.Name)
			fmt.Printf("%s: fuck this shit, i'm out!\n", newData.TeamAway.Name)
		} else if difference == GoodForAway {
			fmt.Printf("%s: fuck this shit, i'm out!\n", newData.TeamHome.Name)
			fmt.Printf("%s: hooray!\n", newData.TeamAway.Name)
		} else {
			fmt.Printf("the difference is neutral\n")
		}

		currentData = newData

		time.Sleep(interval)
	}

	finished <- struct{}{}
}

// Fetch downloads the current match data from u.
func Fetch(u string) (*MatchData, error) {
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body := make([]byte, res.ContentLength) // Interesting: add + 1 to length and observe how unmarshalling falls apart
	read, err := res.Body.Read(body)
	if err != nil {
		fmt.Printf("read %d bytes, body's length is %d\n", read, res.ContentLength)

		if err != io.EOF {
			fmt.Printf("error reading body\n")
			return nil, err
		}
	}

	var matchData MatchData
	err = json.Unmarshal(body, &matchData)
	if err != nil {
		fmt.Printf("error unmarshalling json body\n")
		return nil, err
	}

	return &matchData, nil
}

// Compare compares data of two matches. It returns -1 if the data differs
// positively for home and returns 1 if the new data differs positively for away.
// Returns 0 when the difference in data is neutral.
// Also, "to differ positively for away" == "to differ negatively for home"
func Compare(data *MatchData, prevData *MatchData) Difference {
	fmt.Printf("scr home: %d %s, scr away: %d %s\n",
		data.ScoreHome, data.TeamHome.Name, data.ScoreAway, data.TeamAway.Name)

	if data.ScoreAway > prevData.ScoreAway {
		return GoodForAway
	}

	if data.ScoreHome > prevData.ScoreHome {
		return GoodForHome
	}

	if data.TopInning {
		// IDK, fancy Tramś baseball logic goes here.
	}

	return Neutral
}
