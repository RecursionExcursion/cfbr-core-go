package cfbrcore

import (
	"github.com/RecursionExcursion/cfbr-core-go/internal"
	"github.com/RecursionExcursion/cfbr-core-go/model"
)

func CollectSeason(year int) (internal.SeasonData, error) {
	szn, err := internal.GetSeasonData(year)
	if err != nil {
		return internal.SeasonData{}, err
	}
	return szn, nil
}

func CollectGames(ids []string) ([]model.ESPNCfbGame, error) {
	gd, err := internal.CollectGameData(ids)
	if err != nil {
		return nil, err
	}
	return gd, nil
}
