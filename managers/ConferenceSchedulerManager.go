package managers

import (
	"log"
	"math/rand"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GenerateCanadianHockeySchedule(ts structs.Timestamp) {
	db := dbprovider.GetInstance().GetDB()
	// ts := repository.FindTimestamp()
	chlTeams := GetAllCanadianCHLTeams()
	schedule, err := GenerateCanadaSeasonSchedule(chlTeams, ts.SeasonID, 2000)
	if err != nil {
		log.Println("Generation Failed: ", err)
		return
	}
	finalUpload := []structs.CollegeGame{}
	for _, game := range schedule {
		game.AddWeekData((ts.Season*1000 + uint(game.Round)), uint(game.Round), game.Slot)
		finalUpload = append(finalUpload, game.CollegeGame)
	}
	repository.CreateCHLGamesRecordsBatch(db, finalUpload, 100)
}

func GenerateSimCHLConferenceSchedules(ts structs.Timestamp) {
	db := dbprovider.GetInstance().GetDB()
	collegeTeams := GetAllCollegeTeams()

	collegeGamesBatch := []structs.CollegeGame{}

	conferenceMap := make(map[uint8][]structs.CollegeTeam)
	for _, team := range collegeTeams {
		if team.LeagueID > 1 {
			continue
		}
		conferenceMap[team.ConferenceID] = append(conferenceMap[team.ConferenceID], team)
	}

	ConferenceIDList := make([]uint8, len(conferenceMap))
	for conferenceID := range conferenceMap {
		ConferenceIDList = append(ConferenceIDList, conferenceID)
	}

	for _, conferenceID := range ConferenceIDList {
		// Skip independents
		if conferenceID == 7 {
			continue
		}
		teams := conferenceMap[conferenceID]
		if len(teams) == 0 {
			continue
		}

		scheduleTemplate := getScheduleTemplate(len(teams), conferenceID)

		matches := GenerateConferenceSchedule(teams, scheduleTemplate, ts, conferenceID)
		collegeGamesBatch = append(collegeGamesBatch, matches...)
	}

	repository.CreateCHLGamesRecordsBatch(db, collegeGamesBatch, 200)
}

// GenerateConferenceSchedule builds all 56 conference games (14 per team,
// home-and-away vs. every opponent) for an 8-team conference spanning weeks 9-15.
// The teams slice must contain exactly 8 teams in the desired seeding order.
func GenerateConferenceSchedule(teams []structs.CollegeTeam, entries []scheduleEntry, ts structs.Timestamp, conferenceID uint8) []structs.CollegeGame {

	var shuffledTeams []structs.CollegeTeam

	switch conferenceID {

	default:
		shuffledTeams = make([]structs.CollegeTeam, len(teams))
		copy(shuffledTeams, teams)
		rand.Shuffle(len(shuffledTeams), func(i, j int) {
			shuffledTeams[i], shuffledTeams[j] = shuffledTeams[j], shuffledTeams[i]
		})
	}

	matches := make([]structs.CollegeGame, 0, len(entries))

	for _, e := range entries {
		home := shuffledTeams[e.homeIdx]
		away := shuffledTeams[e.awayIdx]

		match := structs.CollegeGame{
			BaseGame: structs.BaseGame{
				SeasonID:      ts.SeasonID,
				Week:          int(e.week),
				WeekID:        util.GetWeekID(ts.SeasonID, e.week),
				GameDay:       e.slot,
				HomeTeamID:    home.ID,
				HomeTeam:      home.Abbreviation,
				HomeTeamCoach: home.Coach,
				ArenaID:       uint(home.ArenaID),
				Arena:         home.Arena,
				City:          home.City,
				State:         home.State,
				Country:       home.Country,
				AwayTeamID:    away.ID,
				AwayTeam:      away.Abbreviation,
				AwayTeamCoach: away.Coach,
				IsConference:  true,
			},
		}

		matches = append(matches, match)
	}

	return matches
}

func getScheduleTemplate(numOfTeams int, conferenceID uint8) []scheduleEntry {
	switch numOfTeams {
	case 8:
		return eightTeamSchedule
	case 9:
		return nineTeamSchedule
	case 10:
		return tenTeamSchedule
	case 11:
		return elevenTeamSchedule
	case 12:
		return twelveTeamSchedule
	default:
		return []scheduleEntry{}
	}
}
