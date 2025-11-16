package managers

import (
	"sort"
	"strconv"
	"strings"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func HandlePostSeasonMigration() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	if ts.Week < 22 || !ts.CollegeSeasonOver || !ts.NHLSeasonOver {
		return
	}
	HandleTeamProfileValues(ts)
	HandleRecruitingTeamProfileReset()
	ts.MoveUpSeason()
	GenerateStandingsForNewSeason(ts)

	repository.SaveTimestamp(ts, db)
}

func HandleTeamProfileValues(ts structs.Timestamp) structs.Timestamp {
	db := dbprovider.GetInstance().GetDB()
	collegeTeams := GetAllCollegeTeams()
	collegeStandings := GetAllCollegeStandingsBySeasonID(strconv.Itoa(int(ts.SeasonID)))
	allCollegeStandings := repository.FindAllCollegeStandings(repository.StandingsQuery{})
	collegeStandingsMap := MakeCollegeStandingsMap(collegeStandings)
	coachRecords := MakeCollegeStandingsMapByCoachName(allCollegeStandings)
	allCollegeGames := repository.FindCollegeGames(repository.GamesClauses{SeasonID: strconv.Itoa(int(ts.SeasonID))})
	collegeGamesMap := MakeCollegeGameMapByTeamID(allCollegeGames)

	// Prestige Map Modifier. 1 == Did not make past first round of post season, 2 == frozen four, 3 == national championship
	conferencePrestigeMap := make(map[uint]uint)

	for _, t := range collegeTeams {
		if t.ID > 72 && ts.SeasonID == 1 {
			continue
		}

		standings := collegeStandingsMap[t.ID]

		// Program Development Rating
		programDevelopment := t.ProgramPrestige
		winPct := float64(standings.TotalWins) + float64(standings.TotalOTWins)/float64(standings.TotalWins+standings.TotalLosses+standings.TotalOTWins+standings.TotalOTLosses)
		programDevelopmentMod := 0
		if winPct < 0.25 {
			programDevelopmentMod += -3
		} else if winPct < 0.5 {
			programDevelopmentMod += -1
		} else if winPct >= 0.75 {
			programDevelopmentMod += 1
		} else if winPct >= 1 {
			programDevelopmentMod += 2
		}

		if strings.Contains(standings.PostSeasonStatus, "Quarterfinals") {
			programDevelopmentMod += 1
		}

		if strings.Contains(standings.PostSeasonStatus, "Frozen Four") || strings.Contains(standings.PostSeasonStatus, "National Champion Runner-Up") {
			programDevelopmentMod += 2
		} else if strings.Contains(standings.PostSeasonStatus, "National Champion") {
			programDevelopmentMod += 3
		}

		programDevelopment += uint8(programDevelopmentMod)
		if programDevelopment < 1 {
			programDevelopment = 1
		} else if programDevelopment > 9 {
			programDevelopment = 9
		}

		// Professional Development Rating
		// Do this when migrating players from draft to API. Skip.

		// Coach Rating
		// This is the historical w/l record of the coach in question.
		coachRating := 0
		coach := t.Coach
		coachRecords := coachRecords[coach]
		coachWins := 0
		coachOTWins := 0
		coachLosses := 0
		coachOTLosses := 0
		for _, cr := range coachRecords {
			coachWins += int(cr.TotalWins)
			coachOTWins += int(cr.TotalOTWins)
			coachLosses += int(cr.TotalLosses)
			coachOTLosses += int(cr.TotalOTLosses)
		}
		totalGames := coachWins + coachOTWins + coachLosses + coachOTLosses
		winPercentage := float64(coachWins+coachOTWins) / float64(totalGames)
		// Normalize within a 1-9 range
		coachRating = int(winPercentage*8) + 1

		// Season Momentum -- Normalize win percentage to 1-9 scale
		seasonMomentum := int(winPct*8) + 1

		// Conference Prestige
		teamGames := collegeGamesMap[t.ID]
		for _, g := range teamGames {
			if !g.IsPlayoffGame {
				continue
			}
			if strings.Contains(g.GameTitle, "Round of 16") {
				conferencePrestigeMap[uint(t.ConferenceID)] = 1
			} else if strings.Contains(g.GameTitle, "Frozen Four") {
				conferencePrestigeMap[uint(t.ConferenceID)] = 2
			} else if strings.Contains(g.GameTitle, "National Championship") {
				conferencePrestigeMap[uint(t.ConferenceID)] = 3
			}
		}

		t.UpdateRecordRatings(programDevelopment, uint8(coachRating), uint8(seasonMomentum))

		repository.SaveCollegeTeamRecord(db, t)
	}

	for _, t := range collegeTeams {
		if t.ConferenceID == 7 {
			continue // Independents are not in a conference and thus cannot have prestige impacted.
		}
		conferencePrestigeStatus := conferencePrestigeMap[uint(t.ConferenceID)]
		conferenceMod := 0

		switch conferencePrestigeStatus {
		case 1:
			conferenceMod -= 1
		case 2:
			conferenceMod += 1
		case 3:
			conferenceMod += 2
		}

		newConferencePrestige := t.ConferencePrestige + uint8(conferenceMod)
		t.UpdateConferencePrestige(newConferencePrestige)
		repository.SaveCollegeTeamRecord(db, t)
	}

	return ts
}

func HandleRecruitingTeamProfileReset() {
	db := dbprovider.GetInstance().GetDB()

	teamProfiles := repository.FindTeamRecruitingProfiles(false)

	// sort team profiles by composite score
	sort.Slice(teamProfiles, func(i, j int) bool {
		return teamProfiles[i].CompositeScore > teamProfiles[j].CompositeScore
	})

	for idx, tp := range teamProfiles {
		tp.AssignHistoricRank(idx + 1)
	}

	for _, tp := range teamProfiles {
		tp.ResetWeeklyPoints(0, true)
		tp.ResetScholarshipCount()
		tp.ResetStarCount()

		// Save
		repository.SaveTeamProfileRecord(db, tp)
	}

}

func GenerateStandingsForNewSeason(ts structs.Timestamp) structs.Timestamp {
	db := dbprovider.GetInstance().GetDB()

	collegeTeams := repository.FindAllCollegeTeams(repository.TeamClauses{})
	proTeams := repository.FindAllProTeams(repository.TeamClauses{})
	collegeStandingsUpsert := []structs.CollegeStandings{}
	proStandingsUpsert := []structs.ProfessionalStandings{}

	for _, team := range collegeTeams {
		standings := structs.CollegeStandings{
			BaseStandings: structs.BaseStandings{
				TeamID:       team.ID,
				TeamName:     team.TeamName,
				SeasonID:     ts.SeasonID,
				Season:       ts.Season,
				LeagueID:     uint(team.LeagueID),
				ConferenceID: uint(team.ConferenceID),
			},
		}

		collegeStandingsUpsert = append(collegeStandingsUpsert, standings)
	}

	for _, team := range proTeams {
		standings := structs.ProfessionalStandings{
			BaseStandings: structs.BaseStandings{
				TeamID:       team.ID,
				TeamName:     team.TeamName,
				SeasonID:     ts.SeasonID,
				Season:       ts.Season,
				LeagueID:     1,
				ConferenceID: uint(team.ConferenceID),
			},
			DivisionID: uint(team.DivisionID),
		}

		proStandingsUpsert = append(proStandingsUpsert, standings)
	}

	repository.CreateCollegeStandingsRecordsBatch(db, collegeStandingsUpsert, 50)
	repository.CreateProStandingsRecordsBatch(db, proStandingsUpsert, 50)

	return ts
}
