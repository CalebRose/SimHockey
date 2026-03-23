package managers

import (
	"sort"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetAllCollegePollsByWeekIDAndSeasonID(weekID, seasonID string) []structs.CollegePollSubmission {
	return repository.FindAllCollegePolls(weekID, seasonID, "")
}

func GetPollSubmissionBySubmissionID(id string) structs.CollegePollSubmission {
	return repository.FindCollegePollSubmission(id, "", "", "")
}

func GetPollSubmissionByUsernameWeekAndSeason(username string) structs.CollegePollSubmission {
	ts := GetTimestamp()
	weekID := strconv.Itoa(int(ts.WeekID + 1))
	seasonID := strconv.Itoa(int(ts.SeasonID))

	submission := repository.FindCollegePollSubmission("", weekID, seasonID, username)

	return submission
}

func SyncCollegePollSubmissionForCurrentWeek(week, weekID, seasonID uint) {
	db := dbprovider.GetInstance().GetDB()

	weekIDStr := strconv.Itoa(int(weekID))
	seasonIDStr := strconv.Itoa(int(seasonID))
	standingsMap := GetCollegeStandingsMap(seasonIDStr)

	submissions := GetAllCollegePollsByWeekIDAndSeasonID(weekIDStr, seasonIDStr)

	allCollegeTeams := GetAllCollegeTeams()

	voteMap := make(map[uint]*structs.TeamVote)

	for _, t := range allCollegeTeams {
		voteMap[t.ID] = &structs.TeamVote{TeamID: t.ID, Team: t.TeamName}
	}

	for _, s := range submissions {
		if week > 3 {
			// Invalid check
			s1standings := standingsMap[s.Rank1ID]
			if s1standings.TotalWins == 0 {
				continue
			}
			s2standings := standingsMap[s.Rank2ID]
			if s2standings.TotalWins == 0 {
				continue
			}
			s3standings := standingsMap[s.Rank3ID]
			if s3standings.TotalWins == 0 {
				continue
			}
			s4standings := standingsMap[s.Rank4ID]
			if s4standings.TotalWins == 0 {
				continue
			}
			s5standings := standingsMap[s.Rank5ID]
			if s5standings.TotalWins == 0 {
				continue
			}
			s6standings := standingsMap[s.Rank6ID]
			if s6standings.TotalWins == 0 {
				continue
			}
			s7standings := standingsMap[s.Rank7ID]
			if s7standings.TotalWins == 0 {
				continue
			}
			s8standings := standingsMap[s.Rank8ID]
			if s8standings.TotalWins == 0 {
				continue
			}
			s9standings := standingsMap[s.Rank9ID]
			if s9standings.TotalWins == 0 {
				continue
			}
			s10standings := standingsMap[s.Rank10ID]
			if s10standings.TotalWins == 0 {
				continue
			}
			s11standings := standingsMap[s.Rank11ID]
			if s11standings.TotalWins == 0 {
				continue
			}
			s12standings := standingsMap[s.Rank12ID]
			if s12standings.TotalWins == 0 {
				continue
			}
			s13standings := standingsMap[s.Rank13ID]
			if s13standings.TotalWins == 0 {
				continue
			}
			s14standings := standingsMap[s.Rank14ID]
			if s14standings.TotalWins == 0 {
				continue
			}
			s15standings := standingsMap[s.Rank15ID]
			if s15standings.TotalWins == 0 {
				continue
			}
			s16standings := standingsMap[s.Rank16ID]
			if s16standings.TotalWins == 0 {
				continue
			}
			s17standings := standingsMap[s.Rank17ID]
			if s17standings.TotalWins == 0 {
				continue
			}
			s18standings := standingsMap[s.Rank18ID]
			if s18standings.TotalWins == 0 {
				continue
			}
			s19standings := standingsMap[s.Rank19ID]
			if s19standings.TotalWins == 0 {
				continue
			}
			s20standings := standingsMap[s.Rank20ID]
			if s20standings.TotalWins == 0 {
				continue
			}
		}
		voteMap[s.Rank1ID].AddVotes(1)
		voteMap[s.Rank2ID].AddVotes(2)
		voteMap[s.Rank3ID].AddVotes(3)
		voteMap[s.Rank4ID].AddVotes(4)
		voteMap[s.Rank5ID].AddVotes(5)
		voteMap[s.Rank6ID].AddVotes(6)
		voteMap[s.Rank7ID].AddVotes(7)
		voteMap[s.Rank8ID].AddVotes(8)
		voteMap[s.Rank9ID].AddVotes(9)
		voteMap[s.Rank10ID].AddVotes(10)
		voteMap[s.Rank11ID].AddVotes(11)
		voteMap[s.Rank12ID].AddVotes(12)
		voteMap[s.Rank13ID].AddVotes(13)
		voteMap[s.Rank14ID].AddVotes(14)
		voteMap[s.Rank15ID].AddVotes(15)
		voteMap[s.Rank16ID].AddVotes(16)
		voteMap[s.Rank17ID].AddVotes(17)
		voteMap[s.Rank18ID].AddVotes(18)
		voteMap[s.Rank19ID].AddVotes(19)
		voteMap[s.Rank20ID].AddVotes(20)
	}

	allVotes := []structs.TeamVote{}

	for _, t := range allCollegeTeams {
		v := voteMap[t.ID]
		if v.TotalVotes == 0 {
			continue
		}
		newVoteObj := structs.TeamVote{TeamID: v.TeamID, Team: v.Team, TotalVotes: v.TotalVotes, Number1Votes: v.Number1Votes}

		allVotes = append(allVotes, newVoteObj)
	}

	sort.Slice(allVotes, func(i, j int) bool {
		return allVotes[i].TotalVotes > allVotes[j].TotalVotes
	})

	officialPoll := structs.CollegePollOfficial{
		WeekID:   weekID,
		Week:     week,
		SeasonID: seasonID,
	}
	count := 0
	for idx, v := range allVotes {
		if count > 19 {
			break
		}

		count += 1
		officialPoll.AssignRank(idx, v)
		// Get Standings
		teamID := strconv.Itoa(int(v.TeamID))
		teamStandings := standingsMap[v.TeamID]
		rank := idx + 1
		teamStandings.AssignRank(rank)
		repository.SaveCollegeStandingsRecord(teamStandings, db)

		if week > 17 {
			continue
		}

		matches := GetCollegeGamesByTeamIDAndSeasonID(teamID, seasonIDStr, false)
		for _, m := range matches {
			if m.Week < int(week) && !m.IsPreseason {
				continue
			}
			if m.Week > int(week) && !m.IsPreseason {
				continue
			}
			m.AssignRank(v.TeamID, uint(rank))
			repository.SaveCollegeGameRecord(m, db)
		}
	}

	repository.CreateCollegePollRecord(db, officialPoll)
}

func CreatePoll(dto structs.CollegePollSubmission) structs.CollegePollSubmission {
	db := dbprovider.GetInstance().GetDB()
	existingPoll := GetPollSubmissionBySubmissionID(strconv.Itoa(int(dto.ID)))

	if existingPoll.ID > 0 {
		dto.AssignID(existingPoll.ID)
		repository.SaveCollegePollSubmissionRecord(dto, db)
	} else {
		repository.CreateCollegePollSubmissionRecord(db, dto)
	}

	return dto
}

func GetOfficialPollBySeasonID(seasonID string) []structs.CollegePollOfficial {
	return repository.FindCollegePollOfficial(seasonID)
}

func GeneratePreseasonPoll(ts structs.Timestamp) {
	db := dbprovider.GetInstance().GetDB()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	collegeStandings := repository.FindAllCollegeStandings(repository.StandingsQuery{SeasonID: seasonID})
	preseasonGames := repository.FindCollegeGames(repository.GamesClauses{IsPreseason: true, SeasonID: seasonID})

	collegeStandingsMap := MakeCollegeStandingsMap(collegeStandings)

	// Apply preseason game results to temporarily adjust wins/losses used for ranking.
	// CollegeStandings is a value type in the map, so write-back is required.
	for _, game := range preseasonGames {
		homeStandings := collegeStandingsMap[game.HomeTeamID]
		awayStandings := collegeStandingsMap[game.AwayTeamID]

		homeStandings.UpdateStandings(game.BaseGame)
		awayStandings.UpdateStandings(game.BaseGame)

		collegeStandingsMap[game.HomeTeamID] = homeStandings
		collegeStandingsMap[game.AwayTeamID] = awayStandings
	}

	// Build a vote list from the temporarily adjusted standings. Only include
	// college teams (TeamID <= 74) that have a preseason rank assigned.
	preseasonPoll := []structs.TeamVote{}
	for _, standings := range collegeStandingsMap {
		// Calculate a ranking based on the standings records and their preseason record
		if standings.TeamID > 74 {
			continue
		}
		if standings.PreseasonRank == 0 {
			continue
		}

		// Performance score from preseason games + inverted expert rank (rank 1 = 25 pts).
		performanceScore := int(standings.TotalWins)*2 - int(standings.TotalLosses)
		expertScore := int(26 - standings.PreseasonRank)
		totalVotes := performanceScore + expertScore
		if totalVotes < 0 {
			totalVotes = 0
		}
		preseasonPoll = append(preseasonPoll, structs.TeamVote{
			TeamID:     standings.TeamID,
			Team:       standings.TeamName,
			TotalVotes: uint(totalVotes),
		})
	}

	sort.Slice(preseasonPoll, func(i, j int) bool {
		return preseasonPoll[i].TotalVotes > preseasonPoll[j].TotalVotes
	})

	officialPoll := structs.CollegePollOfficial{
		WeekID:   ts.WeekID,
		Week:     ts.Week,
		SeasonID: ts.SeasonID,
	}

	// Use the original (unmodified) standings when saving so that temporary
	// preseason W/L totals are not persisted to the database.
	originalStandingsMap := make(map[uint]structs.CollegeStandings)
	for _, s := range collegeStandings {
		originalStandingsMap[s.TeamID] = s
	}

	count := 0
	for idx, v := range preseasonPoll {
		if count >= 20 {
			break
		}
		count++
		officialPoll.AssignRank(idx, v)

		teamStandings := originalStandingsMap[v.TeamID]
		rank := idx + 1
		teamStandings.AssignRank(rank)
		repository.SaveCollegeStandingsRecord(teamStandings, db)
	}

	repository.CreateCollegePollRecord(db, officialPoll)
}
