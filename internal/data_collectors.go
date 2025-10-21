package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/RecursionExcursion/cfbr-core-go/model"
)

const espnSeasonDateFormat = "2006-01-02T15:04Z"
const espnQueryDateFormat = "20060102"

type Week struct {
	Label     string              `json:"label"`
	StartDate time.Time           `json:"startDate"`
	EndDate   time.Time           `json:"endDate"`
	Events    []model.SeasonEvent `json:"events"`
}

func CompileSeason(year int) (*model.Season, error) {
	wks, err := GetSeasonWeeks(year)
	if err != nil {
		panic(err)
	}

	wks, err = PopulateWeekEvents(wks)
	if err != nil {
		panic(err)
	}
	teamMap, fq, err := collectTeams(wks)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(teamMap))
	fmt.Println(fq)

	ids := []string{}
	for _, e := range wks[0].Events {
		ids = append(ids, e.Id)
	}

	gms, err := CollectGameData(ids)
	if err != nil {
		panic(err)
	}

	jsn, err := json.Marshal(gms)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("gms.json", jsn, 0644)
	if err != nil {
		panic(err)
	}

	// iterate the weeks and pull each set of data

	panic("END OF SCRIPT")

	return &model.Season{}, nil
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
				Label:     fmt.Sprintf("%v %v", szn.Label, wk.Label),
				StartDate: sd,
				EndDate:   ed,
				Events:    []model.SeasonEvent{},
			})
		}
	}
	return weeks, nil
}

func PopulateWeekEvents(wks []Week) ([]Week, error) {
	retWks := []Week{}

	for i, wk := range wks {
		fmt.Printf("Fetching %v %v/%v\n", wk.Label, i+1, len(wks))
		sb, err := fetchEspnScoreboard(80, wk.StartDate.Format(espnQueryDateFormat), wk.EndDate.Format(espnQueryDateFormat))
		if err != nil {
			panic(err)
		}
		wk.Events = append(wk.Events, sb.Events...)
		retWks = append(retWks, wk)
	}
	return retWks, nil
}

func collectTeams(wks []Week) (map[string]model.ESPNCfbTeam, []string, error) {
	queriedMap := map[string]struct{}{}
	teamMap := map[string]model.ESPNCfbTeam{}
	failedQueries := []string{}

	for _, wk := range wks {
		i := 0
		for _, e := range wk.Events {
			i++
			fmt.Printf("Querying teams for %v event %v/%v\n", wk.Label, i, len(wk.Events))
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

					// Check teams group.parent id against set confrence and division ids in constants file
					if _, ok := fbsConferences[tm.Team.Groups.Parent.Id]; ok {
						teamMap[key] = tm.Team
					}

				}
			}
		}
	}
	return teamMap, failedQueries, nil
}

func CollectGameData(gIds []string) ([]model.ESPNCfbGame, error) {
	gms := []model.ESPNCfbGame{}

	for i, id := range gIds {
		fmt.Printf("Collecting game %v/%v\n", i+1, len(gIds))
		gm, err := fetchEspnStats(id)
		if err != nil {
			return nil, err
		}
		gms = append(gms, gm)
	}
	return gms, nil
}
