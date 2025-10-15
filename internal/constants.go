package internal

import "github.com/RecursionExcursion/go-toolkit/core"

/*
// TODO old api endpoints, OBSOLETE
const baseRoute = "https://api.collegefootballdata.com"
const teams = "/teams"       //?year=<year>"
const games = "/games"       //?division=<division>&year=<year>&seasonType=<type>" //fbs?
const stats = "/games/teams" //?year=<year>&week=<week>&seasonType=<type>""

// seasonTypes
const regularSeason = 2
const postseason = 3

// classifications
const fbs = "fbs"
const fcs = "fcs"
const ii = "ii"
const iii = "iii"

var classes = []string{
	fbs, fcs, ii, iii,
} */

/* CFBR batching
 * cfbr only makes 18 req but gets ratelimited pretty quickly, 10 works but is not stable (yet?)
 * 5 seems safe for now until a more robust ratelmiting logic is impl
 */
const batchSize = 5

var BatchRunner = core.RunBatchSizeClosure(batchSize)

/*
//TODO Obsolete as ranking has been moved to client
var trackedStatCategories = []string{
	"totalYards",
} */

var TotalYardsStatKey = "totalYards"

const espnFbsGroupId = 80

/* ESPN Routes */
const espnBase = "https://site.api.espn.com/apis/site/v2/sports/football/college-football"
const espnGroups = "/groups"
const espnSeason = "/scoreboard" //dates=2024 or dates=20240921
const espnTeams = "/teams"       //</teamid>
const espnGame = "/summary"      //?event=<eventId>

/* Group Keys */
type GroupName = struct {
	name     string
	children []string
}

var D1 = GroupName{
	name: "NCAA Division I",
	children: []string{
		"FBS (I-A)",
		"FCS (I-AA)",
	},
}

var D2_3 = GroupName{
	name: "Division II/III",
	children: []string{
		"NCAA Division II",
		"NCAA Division III",
	},
}

/* TODO: DELETE place holder for flow
 * Groups (Teams) -> Season (Games) -> Games
 * The actual teams endpoint may be  moot but we will see what we need from it.
 *
 *
 */

var fbsConferences = map[string]string{
	"1":   "ACC",
	"4":   "Big 12",
	"5":   "Big Ten",
	"8":   "SEC",
	"9":   "Pac-12 (remnants)",
	"12":  "Conference USA",
	"15":  "MAC",
	"17":  "Mountain West",
	"18":  "Independents",
	"37":  "Sun Belt",
	"151": "AAC",
	"80":  "FBS",
}
