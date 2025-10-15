package internal

import (
	"fmt"
	"strconv"
	"time"

	"github.com/RecursionExcursion/cfbr-core-go/model"
)

const espnSeasonDateFormat = "2006-01-02T15:04Z"
const espnQueryDateFormat = "20060102"

func CompileSeason(year int) (*model.Season, error) {
	wks, err := GetSeasonWeeks(year)
	if err != nil {
		panic(err)
	}

	evnts, err := GetSeasonEvents(wks)
	if err != nil {
		panic(err)
	}

	/* Collect teams */
	queriedMap := map[string]struct{}{}
	teamMap := map[string]model.ESPNCfbTeam{}
	failedQueries := []string{}

	i := 0
	for _, e := range evnts {
		i++
		fmt.Printf("Querying teams for event %v/%v\n", i, len(evnts))
		for _, c := range e.Competitions[0].Competitors {
			key := c.Team.Id
			if _, ok := queriedMap[key]; !ok {

				queriedMap[key] = struct{}{}
				id, err := strconv.Atoi(key)
				if err != nil {
					/* This should not happen */
					fmt.Printf("Bad competitor NAME- %v id- %v\n", c.Team.DisplayName, c.Id)
					continue
				}
				tm, err := fetchEspnTeam(id)
				if err != nil {
					fmt.Printf("Team %v(%v) notn found", c.Team.DisplayName, c.Id)
					continue
				}
				// gId, err := strconv.Atoi(tm.Team.Groups.Parent.Id)
				// if err != nil {
				// 	/* If no id then we can assume its not an FBS team */
				// 	fmt.Printf("FBS id for team %v(%v) not found in event %v\n", tm.Team.DisplayName, tm.Team.Id, e.Id)
				// 	failedQueries = append(failedQueries, fmt.Sprintf("%v-%v-%v", tm.Team.Id, tm.Team.DisplayName, e.Id))
				// 	continue
				// }
				// if gId == espnFbsGroupId {
				if _, ok := fbsConferences[tm.Team.Groups.Parent.Id]; ok {
					teamMap[key] = tm.Team
				}
				// }

			}
		}
	}

	// fbsTeams := []model.ESPNCfbTeam{}
	// nonFbsTeams := []model.ESPNCfbTeam{}

	// for _, t := range teamMap {

	// }

	// jsn, err := json.MarshalIndent(teamMap, "", " ")
	// if err != nil {
	// 	panic(err)
	// }

	// os.WriteFile("teams.json", jsn, 0644)

	fmt.Println(len(evnts))
	fmt.Println(len(teamMap))
	fmt.Println(failedQueries)

	// iterate the weeks and pull each set of data

	panic("END OF SCRIPT")

	// startDate, endDate, err := GetSeasonDateRanges(year)
	// if err != nil {
	// 	return nil, err
	// }

	// Get today's date (midnight UTC)
	// now := time.Now()
	// nowMidnightUTC := time.Date(
	// 	now.Year(), now.Month(), now.Day(),
	// 	0, 0, 0, 0,
	// 	time.UTC,
	// )

	// // //Compare it against the end of the season, do not scrape data after the current date
	// // if nowMidnightUTC.Before(endDate) {
	// // 	endDate = nowMidnightUTC
	// // }

	// scd, err := collectSeasonDates(startDate, endDate)
	// if err != nil {
	// 	return nil, err
	// }

	// // scd.FilterFbsTeams()

	// //collect games
	// gms, err := collectGames(scd)
	// if err != nil {
	// 	return nil, err
	// }

	// tms, err := collectTeamInfo(scd)
	// if err != nil {
	// 	return nil, err
	// }

	// szn := model.Season{
	// 	Year:      year,
	// 	Schedules: scd,
	// 	Games:     gms,
	// 	Teams:     tms,
	// }

	return &model.Season{}, nil
}

type Week struct {
	label     string
	startDate time.Time
	endDate   time.Time
}

func GetSeasonWeeks(year int) ([]Week, error) {
	//Query zero week 08/01/YYYY, should have no games
	s, err := fetchEspnSeason(fmt.Sprintf("%v0801", year))
	if err != nil {
		return nil, err
	}

	weeks := []Week{}
	// layout := "2006-01-02T15:04Z" // matches YYYY-MM-DDTHH:MMZ

	for _, szn := range s.Leagues[0].Calender {
		for _, wk := range szn.Entries {
			sd, err := time.Parse(espnSeasonDateFormat, wk.StartDate)
			if err != nil {
				panic(err)
			}
			ed, err := time.Parse(espnSeasonDateFormat, wk.EndDate)
			if err != nil {
				panic(err)
			}

			weeks = append(weeks, Week{
				label:     fmt.Sprintf("%v %v", szn.Label, wk.Label),
				startDate: sd,
				endDate:   ed,
			})
		}
	}
	return weeks, nil
}

func GetSeasonEvents(wks []Week) (map[string]model.SeasonEvent, error) {
	eventMap := map[string]model.SeasonEvent{}

	for i, wk := range wks {
		fmt.Printf("Fetching %v %v/%v\n", wk.label, i+1, len(wks))
		sb, err := fetchEspnScoreboard(80, wk.startDate.Format(espnQueryDateFormat), wk.endDate.Format(espnQueryDateFormat))
		if err != nil {
			panic(err)
		}
		for _, e := range sb.Events {
			eventMap[e.Id] = e
		}
	}
	return eventMap, nil
}

// func GetSeasonDateRanges(year int) (start time.Time, end time.Time, err error) {
// 	s, err := getZeroDay(year)
// 	if err != nil {
// 		return time.Time{}, time.Time{}, err
// 	}
// 	//get regular season dates

// 	if len(s.Leagues) == 0 {
// 		return time.Time{}, time.Time{}, errors.New("no leagues found")
// 	}
// 	for _, c := range s.Leagues[0].Calender {
// 		if c.Label == "Regular Season" {
// 			start, err = time.Parse(espnSeasonDateFormat, c.StartDate)
// 			if err != nil {
// 				return start, end, err
// 			}
// 			end, err = time.Parse(espnSeasonDateFormat, c.EndDate)
// 			if err != nil {
// 				return start, end, err
// 			}
// 		}
// 		if c.Label == "Postseason" {
// 			end, err = time.Parse(espnSeasonDateFormat, c.EndDate)
// 			if err != nil {
// 				return start, end, err
// 			}
// 		}
// 	}
// 	return start, end, nil
// }

// func getZeroDay(year int) (model.ESPNScoreboard, error) {
// 	//0 day query 08/01
// 	query := fmt.Sprintf("%v", year)
// 	s, err := fetchEspnSeason(query)
// 	if err != nil {
// 		return model.ESPNScoreboard{}, err
// 	}
// 	return s, nil
// }

// func collectSeasonDates(startDate time.Time, endDate time.Time) (model.SeasonSchedules, error) {
// 	tc := make(model.SeasonSchedules)
// 	currDate := startDate

// 	diff := endDate.Sub(startDate)
// 	days := int(diff.Hours()/24) + 1
// 	i := 1

// 	for {
// 		//call api
// 		res, err := fetchEspnSeason(currDate.Format(espnQueryDateFormat))
// 		if err != nil {
// 			fmt.Println(err)
// 			return tc, err
// 		}
// 		//proccess req into map
// 		for _, e := range res.Events {
// 			match := e.Competitions[0]
// 			t1 := match.Competitors[0]
// 			t2 := match.Competitors[1]

// 			tc.Add(t1, t2, match)
// 			tc.Add(t2, t1, match)
// 		}
// 		log.Printf("Query for %v complete (%v/%v)", currDate, i, days)
// 		i++

// 		//inc date
// 		currDate = currDate.Add(time.Hour * 24)
// 		//exit
// 		if currDate.After(endDate) {
// 			break
// 		}
// 	}

// 	return tc, nil
// }

// // TODO need to batch fetch calls for speed as order is irrelevant
// func collectGames(st model.SeasonSchedules) (model.SeasonGames, error) {

// 	games := make(model.SeasonGames)

// 	gIdMap := make(map[string]struct{})

// 	for _, s := range st {
// 		for _, g := range s.Schedule {
// 			gIdMap[g.GameId] = struct{}{}
// 		}
// 	}

// 	l := len(gIdMap)
// 	i := 1

// 	for gId := range gIdMap {
// 		id, err := strconv.Atoi(gId)
// 		if err != nil {
// 			return games, err
// 		}
// 		gm, err := fetchEspnStats(id)
// 		if err != nil {
// 			return games, err
// 		}
// 		log.Printf("Query for Game %v complete (%v/%v)", id, i, l)
// 		i++

// 		games[gId] = gm
// 	}

// 	return games, nil

// }

// func collectTeamInfo(st model.SeasonSchedules) (model.SeasonTeams, error) {
// 	tmMap := make(model.SeasonTeams)

// 	l := len(st)
// 	index := 1

// 	for id := range st {
// 		i, err := strconv.Atoi(id)
// 		if err != nil {
// 			panic(err)
// 		}
// 		tm, err := fetchEspnTeam(i)
// 		if err != nil {
// 			panic(err)
// 		}
// 		log.Printf("Query for Team %v complete (%v/%v)", id, index, l)
// 		index++

// 		tmMap[tm.Team.Id] = tm.Team

// 	}
// 	return tmMap, nil
// }

// func filterFbsTeams(tc model.SeasonSchedules) {
// 	/* Calculate max games played */
// 	maxGamesPlayed := 0
// 	for _, v := range tc {
// 		if maxGamesPlayed < v.GamesPlayed {
// 			maxGamesPlayed = v.GamesPlayed
// 		}
// 	}

// 	//Filter out teams that have not played the maxGamedplayed -2
// 	/* TODO this will break during postseason and the first couple weeks
// 	consider scrapping the entire season??? Or the logic need to be rehashed,
// 	will work for now, need to dig deeper into the ESPN API */
// 	toDelete := []string{}
// 	for k, v := range tc {
// 		if v.GamesPlayed < maxGamesPlayed-2 {
// 			toDelete = append(toDelete, k)
// 		}
// 	}

// 	for _, id := range toDelete {
// 		delete(tc, id)
// 	}

// 	fmt.Println(len(tc))
// }
