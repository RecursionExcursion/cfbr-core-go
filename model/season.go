package model

/* Top Level DS for Season */
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
	toDelete := []string{}

	for k, v := range tc {
		// most fbs teams play 12+ games, 10 gives it a nice buffer (134 teams in 2024)
		if v.GamesPlayed < 10 {
			toDelete = append(toDelete, k)
		}
	}

	for _, id := range toDelete {
		delete(tc, id)
	}

	/* At this point *most teams will be filtered but.....
	* the geniuses over at ESPN include future fbs addtions
	* so we need to cross ref the scheduled and ensure the majority of games
	* are not paycheck games (fbs vs fcs)
	 */

	toDelete = []string{}
	for k, v := range tc {
		fbsGames := 0
		for _, g := range v.Schedule {
			_, exists := tc[g.OppId]
			if exists {
				fbsGames++
			}
		}
		fbsRatio := float32(fbsGames) / float32(v.GamesPlayed)

		// 50% games are played against fbs teams, this number is negotiable
		if fbsRatio < .5 {
			toDelete = append(toDelete, k)
		}

	}

	for _, id := range toDelete {
		delete(tc, id)
	}
}
