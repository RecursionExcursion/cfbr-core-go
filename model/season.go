package model

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
