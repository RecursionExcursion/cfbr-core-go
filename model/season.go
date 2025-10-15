package model

import "fmt"

/* Top Level data structure for Season */
type Season struct {
	Year      int
	Schedules SeasonSchedules
	Games     SeasonGames
	Teams     SeasonTeams
}

type SeasonOccurences struct {
	GamesPlayed int
	Schedule    []CollectedGame
}

type CollectedGame struct {
	GameId string
	OppId  string
}

type SeasonSchedules map[string]*SeasonOccurences
type SeasonGames map[string]ESPNCfbGame
type SeasonTeams map[string]ESPNCfbTeam

func (tc SeasonSchedules) Add(c Competitor, opp Competitor, match SeasonCompetition) {
	so, exists := tc[c.Id]

	cg := CollectedGame{
		GameId: match.Id,
		OppId:  opp.Id,
	}

	if exists {
		so.GamesPlayed++
		so.Schedule = append(so.Schedule, cg)
	} else {
		tc[c.Id] = &SeasonOccurences{
			GamesPlayed: 1,
			Schedule:    []CollectedGame{cg},
		}
	}
}

func (tc SeasonSchedules) FilterFbsTeams() {
	/* Calculate max games played */
	maxGamesPlayed := 0
	for _, v := range tc {
		if maxGamesPlayed < v.GamesPlayed {
			maxGamesPlayed = v.GamesPlayed
		}
	}

	//Filter out teams that have not played the maxGamedplayed -2
	/* TODO this will break during postseason and the first couple weeks
	consider scrapping the entire season??? Or the logic need to be rehashed,
	will work for now, need to dig deeper into the ESPN API */
	toDelete := []string{}
	for k, v := range tc {
		if v.GamesPlayed < maxGamesPlayed-2 {
			toDelete = append(toDelete, k)
		}
	}

	for _, id := range toDelete {
		delete(tc, id)
	}

	fmt.Println(len(tc))
}
