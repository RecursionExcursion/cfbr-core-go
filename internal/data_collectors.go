package internal

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/RecursionExcursion/cfbr-core-go/model"
)

const espnSeasonDateFormat = "2006-01-02T15:04Z"
const espnQueryDateFormat = "20060102"

// /* Top Level DS for Season */
// type Season struct {
// 	Year      int
// 	Schedules SeasonSchedules
// 	Games     SeasonGames
// 	Teams     SeasonTeams
// }

// type SeasonOccurences struct {
// 	GamesPlayed int
// 	Schedule    []CollectedGame
// }

// type CollectedGame struct {
// 	GameId string
// 	OppId  string
// }

// type SeasonSchedules map[string]*SeasonOccurences
// type SeasonGames map[string]ESPNCfbGame
// type SeasonTeams map[string]ESPNCfbTeam

func CompileSeason(year int) (model.Season, error) {
	s, err := getZeroDay(year)
	if err != nil {
		//TODO
		panic(err)
	}

	startDate, endDate, err := getSeasonDateRanges(s)
	if err != nil {
		//TODO
		panic(err)
	}

	scd, err := collectSeasonDates(startDate, endDate)
	if err != nil {
		//TODO
		panic(err)
	}

	scd.FilterFbsTeams()

	//collect games
	gms, err := collectGames(scd)
	if err != nil {
		//TODO
		panic(err)
	}

	tms, err := collectTeamInfo(scd)
	if err != nil {
		panic(err)
	}

	szn := model.Season{
		Year:      year,
		Schedules: scd,
		Games:     gms,
		Teams:     tms,
	}

	return szn, nil
}

func getZeroDay(year int) (model.ESPNSeason, error) {
	//0 day query 08/01
	query := fmt.Sprintf("%v0801", year)
	s, err := fetchEspnSeason(query)
	if err != nil {
		return model.ESPNSeason{}, err
	}
	return s, nil
}

func getSeasonDateRanges(s model.ESPNSeason) (start time.Time, end time.Time, err error) {
	//get regualr season dates

	if len(s.Leagues) == 0 {
		return time.Time{}, time.Time{}, errors.New("no leagues found")
	}
	for _, c := range s.Leagues[0].Calender {
		if c.Label == "Regular Season" {
			start, err = time.Parse(espnSeasonDateFormat, c.StartDate)
			if err != nil {
				return start, end, err
			}
			end, err = time.Parse(espnSeasonDateFormat, c.EndDate)
			if err != nil {
				return start, end, err
			}
		}
		if c.Label == "Postseason" {
			end, err = time.Parse(espnSeasonDateFormat, c.EndDate)
			if err != nil {
				return start, end, err
			}
		}
	}
	return start, end, nil
}

func collectSeasonDates(startDate time.Time, endDate time.Time) (model.SeasonSchedules, error) {
	tc := make(model.SeasonSchedules)
	currDate := startDate

	diff := endDate.Sub(startDate)
	days := int(diff.Hours()/24) + 1
	i := 1

	for {
		//call api
		res, err := fetchEspnSeason(currDate.Format(espnQueryDateFormat))
		if err != nil {
			return tc, err
		}
		//proccess req into map
		for _, e := range res.Events {
			match := e.Competitions[0]
			t1 := match.Competitors[0]
			t2 := match.Competitors[1]

			tc.Add(t1, t2, match)
			tc.Add(t2, t1, match)
		}
		log.Printf("Query for %v complete (%v/%v)", currDate, i, days)
		i++

		//inc date
		currDate = currDate.Add(time.Hour * 24)
		//exit
		if currDate.After(endDate) {
			break
		}
	}

	return tc, nil
}

// TODO need to batch that fetch calls for speed as order is irrelevant
func collectGames(st model.SeasonSchedules) (model.SeasonGames, error) {

	games := make(model.SeasonGames)

	gIdMap := make(map[string]struct{})

	for _, s := range st {
		for _, g := range s.Schedule {
			gIdMap[g.GameId] = struct{}{}
		}
	}

	l := len(gIdMap)
	i := 1

	for gId := range gIdMap {
		id, err := strconv.Atoi(gId)
		if err != nil {
			return games, err
		}
		gm, err := fetchEspnStats(id)
		if err != nil {
			return games, err
		}
		log.Printf("Query for Game %v complete (%v/%v)", id, i, l)
		i++

		games[gId] = gm
	}

	return games, nil

}

func collectTeamInfo(st model.SeasonSchedules) (model.SeasonTeams, error) {
	tmMap := make(model.SeasonTeams)

	l := len(st)
	index := 1

	for id := range st {
		i, err := strconv.Atoi(id)
		if err != nil {
			panic(err)
		}
		tm, err := fetchEspnTeam(i)
		if err != nil {
			panic(err)
		}
		log.Printf("Query for Team %v complete (%v/%v)", id, index, l)
		index++

		tmMap[tm.Team.Id] = tm.Team

	}
	return tmMap, nil
}
