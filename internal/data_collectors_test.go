package internal

import (
	"testing"
)

//TODO WRITE TESTSSSSSS

const year = 2025
const expectedTeamCount = 136

func TestGetSeasonData(t *testing.T) {
	szn, err := GetSeasonData(year)
	if err != nil {
		t.Error(err)
	}

	numTeamns := len(szn.Teams)

	if numTeamns != expectedTeamCount {
		t.Errorf("Collected %v teams. Expected %v", numTeamns, expectedTeamCount)
	}

}
