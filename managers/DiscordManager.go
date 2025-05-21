package managers

import (
	"log"
	"strconv"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetCHLTeamDataForDiscord(id string) structs.CollegeTeamResponseData {
	ts := GetTimestamp()
	seasonId := strconv.Itoa(int(ts.SeasonID))

	team := repository.FindCollegeTeamRecord(id)
	standings := repository.FindAllCollegeStandings(seasonId, "", id)
	matches := repository.FindCollegeGames(repository.GamesClauses{SeasonID: seasonId, TeamID: id, IsPreseason: ts.IsPreseason})
	wins := 0
	losses := 0
	confWins := 0
	confLosses := 0
	otWins := 0
	otLosses := 0
	soWins := 0
	soLosses := 0
	matchList := []structs.CollegeGame{}

	for _, m := range matches {
		if m.Week > int(ts.Week) {
			break
		}
		gameNotRan := (m.GameDay == "A" && !ts.GamesARan) ||
			(m.GameDay == "B" && !ts.GamesBRan) ||
			(m.GameDay == "C" && !ts.GamesCRan) ||
			(m.GameDay == "D" && !ts.GamesDRan)

		earlierWeek := m.Week < int(ts.Week)

		if ((strconv.Itoa(int(m.HomeTeamID)) == id && m.HomeTeamWin) ||
			(strconv.Itoa(int(m.AwayTeamID)) == id && m.AwayTeamWin)) && (earlierWeek || !gameNotRan) {
			wins += 1
			if m.IsConference {
				confWins += 1
			}
			if m.IsOvertime {
				otWins += 1
			}
			if m.IsShootout {
				soWins += 1
			}
		} else if ((strconv.Itoa(int(m.HomeTeamID)) == id && m.AwayTeamWin) ||
			(strconv.Itoa(int(m.AwayTeamID)) == id && m.HomeTeamWin)) && (earlierWeek || !gameNotRan) {
			losses += 1
			if m.IsConference {
				confLosses += 1
			}
			if m.IsOvertime {
				otLosses += 1
			}
			if m.IsShootout {
				soLosses += 1
			}
		}
		if gameNotRan {
			m.HideScore()
		}
		if m.Week == int(ts.Week) {
			matchList = append(matchList, m)
		}
	}
	standing := standings[0]
	standing.MaskGames(uint8(wins), uint8(losses), uint8(confWins), uint8(confLosses), uint8(otWins), uint8(otLosses), uint8(soWins), uint8(soLosses))

	return structs.CollegeTeamResponseData{
		TeamData:        team,
		TeamStandings:   standing,
		UpcomingMatches: matchList,
	}
}

func GetPHLTeamDataForDiscord(id string) structs.ProTeamResponseData {
	ts := GetTimestamp()
	seasonId := strconv.Itoa(int(ts.SeasonID))

	team := repository.FindProTeamRecord(id)
	standings := repository.FindAllProfessionalStandings(seasonId, "", id)
	matches := repository.FindProfessionalGames(repository.GamesClauses{SeasonID: seasonId, TeamID: id, IsPreseason: ts.IsPreseason})
	wins := 0
	losses := 0
	confWins := 0
	confLosses := 0
	otWins := 0
	otLosses := 0
	soWins := 0
	soLosses := 0
	matchList := []structs.ProfessionalGame{}

	for _, m := range matches {
		if m.Week > int(ts.Week) {
			break
		}
		gameNotRan := (m.GameDay == "A" && !ts.GamesARan) ||
			(m.GameDay == "B" && !ts.GamesBRan) ||
			(m.GameDay == "C" && !ts.GamesCRan) ||
			(m.GameDay == "D" && !ts.GamesDRan)

		earlierWeek := m.Week < int(ts.Week)

		if ((strconv.Itoa(int(m.HomeTeamID)) == id && m.HomeTeamWin) ||
			(strconv.Itoa(int(m.AwayTeamID)) == id && m.AwayTeamWin)) && (earlierWeek || !gameNotRan) {
			wins += 1
			if m.IsConference {
				confWins += 1
			}
			if m.IsOvertime {
				otWins += 1
			}
			if m.IsShootout {
				soWins += 1
			}
		} else if ((strconv.Itoa(int(m.HomeTeamID)) == id && m.AwayTeamWin) ||
			(strconv.Itoa(int(m.AwayTeamID)) == id && m.HomeTeamWin)) && (earlierWeek || !gameNotRan) {
			losses += 1
			if m.IsConference {
				confLosses += 1
			}
			if m.IsOvertime {
				otLosses += 1
			}
			if m.IsShootout {
				soLosses += 1
			}
		}
		if gameNotRan {
			m.HideScore()
		}
		if m.Week == int(ts.Week) {
			matchList = append(matchList, m)
		}
	}
	standing := standings[0]
	standing.MaskGames(uint8(wins), uint8(losses), uint8(confWins), uint8(confLosses), uint8(otWins), uint8(otLosses), uint8(soWins), uint8(soLosses))

	return structs.ProTeamResponseData{
		TeamData:        team,
		TeamStandings:   standing,
		UpcomingMatches: matchList,
	}
}

func GetCollegePlayerViaDiscord(id string) structs.DiscordPlayer {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	seasonID := strconv.Itoa(int(ts.SeasonID))
	var player structs.CollegePlayer
	_, gt := ts.GetCurrentGameType(true)
	db.Preload("SeasonStats", "season_id = ? AND game_type = ?", seasonID, gt).Where("id = ?", id).Find(&player)

	return structs.DiscordPlayer{
		FirstName:          player.FirstName,
		LastName:           player.LastName,
		Position:           player.Position,
		Archetype:          player.Archetype,
		TeamID:             player.TeamID,
		Team:               player.Team,
		Height:             player.Height,
		Weight:             player.Weight,
		Stars:              player.Stars,
		Age:                player.Age,
		Overall:            util.GetLetterGrade(int(player.Overall), player.Year),
		Agility:            util.GetLetterGrade(int(player.Agility), player.Year),
		Faceoffs:           util.GetLetterGrade(int(player.Faceoffs), player.Year),
		LongShotAccuracy:   util.GetLetterGrade(int(player.LongShotAccuracy), player.Year),
		LongShotPower:      util.GetLetterGrade(int(player.LongShotPower), player.Year),
		CloseShotAccuracy:  util.GetLetterGrade(int(player.CloseShotAccuracy), player.Year),
		CloseShotPower:     util.GetLetterGrade(int(player.CloseShotPower), player.Year),
		Passing:            util.GetLetterGrade(int(player.Passing), player.Year),
		PuckHandling:       util.GetLetterGrade(int(player.PuckHandling), player.Year),
		Strength:           util.GetLetterGrade(int(player.Strength), player.Year),
		BodyChecking:       util.GetLetterGrade(int(player.BodyChecking), player.Year),
		StickChecking:      util.GetLetterGrade(int(player.StickChecking), player.Year),
		ShotBlocking:       util.GetLetterGrade(int(player.ShotBlocking), player.Year),
		Goalkeeping:        util.GetLetterGrade(int(player.Goalkeeping), player.Year),
		GoalieVision:       util.GetLetterGrade(int(player.GoalieVision), player.Year),
		HighSchool:         player.HighSchool,
		City:               player.City,
		State:              player.State,
		Country:            player.Country,
		OriginalTeamID:     player.OriginalTeamID,
		OriginalTeam:       player.OriginalTeam,
		PreviousTeamID:     player.PreviousTeamID,
		PreviousTeam:       player.PreviousTeam,
		Competitiveness:    util.GetCompetitivenessLabel(int(player.Competitiveness)),
		TeamLoyalty:        util.GetTeamLoyaltyLabel(int(player.TeamLoyalty)),
		PlaytimePreference: util.GetPlaytimePreferenceLabel(int(player.PlaytimePreference)),
		PlayerMorale:       player.PlayerMorale,
		Personality:        player.Personality,
		RelativeID:         player.RelativeID,
		RelativeType:       player.RelativeType,
		Notes:              player.Notes,
		CollegeStats:       player.SeasonStats,
		Stamina:            util.GetPotentialGrade(int(player.Stamina)),
		InjuryRating:       util.GetPotentialGrade(int(player.InjuryRating)),
		IsRedshirt:         player.IsRedshirt,
		IsRedshirting:      player.IsRedshirting,
		Year:               uint8(player.Year),
		PlayerID:           player.ID,
	}
}

func GetCollegePlayerByNameViaDiscord(firstName, lastName, teamID string) structs.DiscordPlayer {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	_, gt := ts.GetCurrentGameType(true)

	seasonID := strconv.Itoa(int(ts.SeasonID))
	var player structs.CollegePlayer

	db.Preload("SeasonStats", "season_id = ? AND game_type = ?", seasonID, gt).
		Where("first_name = ? AND last_name = ? and team_id = ?", firstName, lastName, teamID).
		Find(&player)

	return structs.DiscordPlayer{
		FirstName:          player.FirstName,
		LastName:           player.LastName,
		Position:           player.Position,
		Archetype:          player.Archetype,
		TeamID:             player.TeamID,
		Team:               player.Team,
		Height:             player.Height,
		Weight:             player.Weight,
		Stars:              player.Stars,
		Age:                player.Age,
		Overall:            util.GetLetterGrade(int(player.Overall), player.Year),
		Agility:            util.GetLetterGrade(int(player.Agility), player.Year),
		Faceoffs:           util.GetLetterGrade(int(player.Faceoffs), player.Year),
		LongShotAccuracy:   util.GetLetterGrade(int(player.LongShotAccuracy), player.Year),
		LongShotPower:      util.GetLetterGrade(int(player.LongShotPower), player.Year),
		CloseShotAccuracy:  util.GetLetterGrade(int(player.CloseShotAccuracy), player.Year),
		CloseShotPower:     util.GetLetterGrade(int(player.CloseShotPower), player.Year),
		Passing:            util.GetLetterGrade(int(player.Passing), player.Year),
		PuckHandling:       util.GetLetterGrade(int(player.PuckHandling), player.Year),
		Strength:           util.GetLetterGrade(int(player.Strength), player.Year),
		BodyChecking:       util.GetLetterGrade(int(player.BodyChecking), player.Year),
		StickChecking:      util.GetLetterGrade(int(player.StickChecking), player.Year),
		ShotBlocking:       util.GetLetterGrade(int(player.ShotBlocking), player.Year),
		Goalkeeping:        util.GetLetterGrade(int(player.Goalkeeping), player.Year),
		GoalieVision:       util.GetLetterGrade(int(player.GoalieVision), player.Year),
		Stamina:            util.GetPotentialGrade(int(player.Stamina)),
		InjuryRating:       util.GetPotentialGrade(int(player.InjuryRating)),
		IsRedshirt:         player.IsRedshirt,
		IsRedshirting:      player.IsRedshirting,
		Year:               uint8(player.Year),
		PlayerID:           player.ID,
		HighSchool:         player.HighSchool,
		City:               player.City,
		State:              player.State,
		Country:            player.Country,
		OriginalTeamID:     player.OriginalTeamID,
		OriginalTeam:       player.OriginalTeam,
		PreviousTeamID:     player.PreviousTeamID,
		PreviousTeam:       player.PreviousTeam,
		Competitiveness:    util.GetCompetitivenessLabel(int(player.Competitiveness)),
		TeamLoyalty:        util.GetTeamLoyaltyLabel(int(player.TeamLoyalty)),
		PlaytimePreference: util.GetPlaytimePreferenceLabel(int(player.PlaytimePreference)),
		PlayerMorale:       player.PlayerMorale,
		Personality:        player.Personality,
		RelativeID:         player.RelativeID,
		RelativeType:       player.RelativeType,
		Notes:              player.Notes,
		CollegeStats:       player.SeasonStats,
	}
}

func GetProPlayerViaDiscord(id string) structs.ProDiscordPlayer {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	seasonID := strconv.Itoa(int(ts.SeasonID))
	var player structs.ProfessionalPlayer
	_, gt := ts.GetCurrentGameType(false)
	db.Preload("SeasonStats", "season_id = ? AND game_type = ?", seasonID, gt).Where("id = ?", id).Find(&player)

	return structs.ProDiscordPlayer{
		FirstName:             player.FirstName,
		LastName:              player.LastName,
		Position:              player.Position,
		Archetype:             player.Archetype,
		TeamID:                player.TeamID,
		Team:                  player.Team,
		Height:                player.Height,
		Weight:                player.Weight,
		Stars:                 player.Stars,
		Age:                   player.Age,
		Overall:               player.Overall,
		Agility:               player.Agility,
		Faceoffs:              player.Faceoffs,
		LongShotAccuracy:      player.LongShotAccuracy,
		LongShotPower:         player.LongShotPower,
		CloseShotAccuracy:     player.CloseShotAccuracy,
		CloseShotPower:        player.CloseShotPower,
		Passing:               player.Passing,
		PuckHandling:          player.PuckHandling,
		Strength:              player.Strength,
		BodyChecking:          player.BodyChecking,
		StickChecking:         player.StickChecking,
		ShotBlocking:          player.ShotBlocking,
		Goalkeeping:           player.Goalkeeping,
		GoalieVision:          player.GoalieVision,
		HighSchool:            player.HighSchool,
		City:                  player.City,
		State:                 player.State,
		Country:               player.Country,
		OriginalTeamID:        player.OriginalTeamID,
		OriginalTeam:          player.OriginalTeam,
		PreviousTeamID:        player.PreviousTeamID,
		PreviousTeam:          player.PreviousTeam,
		PlayerMorale:          player.PlayerMorale,
		Personality:           player.Personality,
		RelativeID:            player.RelativeID,
		RelativeType:          player.RelativeType,
		Notes:                 player.Notes,
		ProStats:              player.SeasonStats,
		IsFreeAgent:           player.IsFreeAgent,
		MarketPreference:      util.GetFAMarketPrefLabel(int(player.MarketPreference)),
		CompetitivePreference: util.GetFACompetitivePrefLabel(int(player.CompetitivePreference)),
		FinancialPreference:   util.GetFAFinancialPrefLabel(int(player.FinancialPreference)),
		Stamina:               util.GetPotentialGrade(int(player.Stamina)),
		InjuryRating:          util.GetPotentialGrade(int(player.InjuryRating)),
		Year:                  uint8(player.Year),
		PlayerID:              player.ID,
	}
}

func GetProPlayerByNameViaDiscord(firstName, lastName, teamID string) structs.ProDiscordPlayer {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	_, gt := ts.GetCurrentGameType(false)

	seasonID := strconv.Itoa(int(ts.SeasonID))
	var player structs.ProfessionalPlayer

	db.Preload("SeasonStats", "season_id = ? AND game_type = ?", seasonID, gt).
		Where("first_name = ? AND last_name = ? and team_id = ?", firstName, lastName, teamID).
		Find(&player)

	return structs.ProDiscordPlayer{
		FirstName:             player.FirstName,
		LastName:              player.LastName,
		Position:              player.Position,
		Archetype:             player.Archetype,
		TeamID:                player.TeamID,
		Team:                  player.Team,
		Height:                player.Height,
		Weight:                player.Weight,
		Stars:                 player.Stars,
		Age:                   player.Age,
		Overall:               player.Overall,
		Agility:               player.Agility,
		Faceoffs:              player.Faceoffs,
		LongShotAccuracy:      player.LongShotAccuracy,
		LongShotPower:         player.LongShotPower,
		CloseShotAccuracy:     player.CloseShotAccuracy,
		CloseShotPower:        player.CloseShotPower,
		Passing:               player.Passing,
		PuckHandling:          player.PuckHandling,
		Strength:              player.Strength,
		BodyChecking:          player.BodyChecking,
		StickChecking:         player.StickChecking,
		ShotBlocking:          player.ShotBlocking,
		Goalkeeping:           player.Goalkeeping,
		GoalieVision:          player.GoalieVision,
		HighSchool:            player.HighSchool,
		City:                  player.City,
		State:                 player.State,
		Country:               player.Country,
		OriginalTeamID:        player.OriginalTeamID,
		OriginalTeam:          player.OriginalTeam,
		PreviousTeamID:        player.PreviousTeamID,
		PreviousTeam:          player.PreviousTeam,
		PlayerMorale:          player.PlayerMorale,
		Personality:           player.Personality,
		RelativeID:            player.RelativeID,
		RelativeType:          player.RelativeType,
		Notes:                 player.Notes,
		ProStats:              player.SeasonStats,
		IsFreeAgent:           player.IsFreeAgent,
		MarketPreference:      util.GetFAMarketPrefLabel(int(player.MarketPreference)),
		CompetitivePreference: util.GetFACompetitivePrefLabel(int(player.CompetitivePreference)),
		FinancialPreference:   util.GetFAFinancialPrefLabel(int(player.FinancialPreference)),
		Stamina:               util.GetPotentialGrade(int(player.Stamina)),
		InjuryRating:          util.GetPotentialGrade(int(player.InjuryRating)),
		Year:                  uint8(player.Year),
		PlayerID:              player.ID,
	}
}

func CompareCHLTeams(t1ID, t2ID string) structs.TeamComparisonModel {
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	teamOneChan := make(chan structs.CollegeTeam)
	teamTwoChan := make(chan structs.CollegeTeam)

	go func() {
		t1 := GetCollegeTeamByTeamID(t1ID)
		teamOneChan <- t1
	}()

	teamOne := <-teamOneChan
	close(teamOneChan)

	go func() {
		t2 := GetCollegeTeamByTeamID(t2ID)
		teamTwoChan <- t2
	}()

	teamTwo := <-teamTwoChan
	close(teamTwoChan)

	allTeamOneGames := repository.FindCollegeGames(repository.GamesClauses{SeasonID: seasonID, TeamID: t1ID, IsPreseason: false})

	t1Wins := 0
	t1Losses := 0
	t1OTWins := 0
	t1OTLosses := 0
	t1SOWins := 0
	t1SOLosses := 0
	t1Streak := 0
	t1CurrentStreak := 0
	t1LargestMarginSeason := 0
	t1LargestMarginDiff := 0
	t1LargestMarginScore := ""
	t2Wins := 0
	t2Losses := 0
	t2OTWins := 0
	t2OTLosses := 0
	t2SOWins := 0
	t2SOLosses := 0
	t2Streak := 0
	t2CurrentStreak := 0
	latestWin := ""
	t2LargestMarginSeason := 0
	t2LargestMarginDiff := 0
	t2LargestMarginScore := ""

	for _, game := range allTeamOneGames {
		gameNotRan := (game.GameDay == "A" && !ts.GamesARan) ||
			(game.GameDay == "B" && !ts.GamesBRan) ||
			(game.GameDay == "C" && !ts.GamesCRan) ||
			(game.GameDay == "D" && !ts.GamesDRan)
		if !game.GameComplete ||
			gameNotRan {
			continue
		}
		doComparison := (game.HomeTeamID == teamOne.ID && game.AwayTeamID == teamTwo.ID) ||
			(game.HomeTeamID == teamTwo.ID && game.AwayTeamID == teamOne.ID)

		if !doComparison {
			continue
		}
		homeTeamTeamOne := game.HomeTeamID == teamOne.ID
		if homeTeamTeamOne {
			if game.HomeTeamWin {
				t1Wins += 1
				t1CurrentStreak += 1
				if game.IsOvertime {
					t1OTWins += 1
				}
				if game.IsShootout {
					t1SOWins += 1
				}
				latestWin = game.HomeTeam
				diff := game.HomeTeamScore - game.AwayTeamScore
				if diff > uint(t1LargestMarginDiff) {
					t1LargestMarginDiff = int(diff)
					t1LargestMarginSeason = int(game.SeasonID) + 2024
					t1LargestMarginScore = "" + strconv.Itoa(int(game.HomeTeamScore)) + "-" + strconv.Itoa(int(game.AwayTeamScore))
				}
			} else {
				t1Streak = t1CurrentStreak
				t1CurrentStreak = 0
				t1Losses += 1
				if game.IsOvertime {
					t1OTLosses += 1
				}
				if game.IsShootout {
					t1SOLosses += 1
				}
			}
		} else {
			if game.HomeTeamWin {
				t2Wins += 1
				t2CurrentStreak += 1
				if game.IsOvertime {
					t2OTWins += 1
				}
				if game.IsShootout {
					t2SOWins += 1
				}
				latestWin = game.HomeTeam
				diff := game.HomeTeamScore - game.AwayTeamScore
				if diff > uint(t2LargestMarginDiff) {
					t2LargestMarginDiff = int(diff)
					t2LargestMarginSeason = int(game.SeasonID) + 2024
					t2LargestMarginScore = "" + strconv.Itoa(int(game.HomeTeamScore)) + "-" + strconv.Itoa(int(game.AwayTeamScore))
				}
			} else {
				t2Streak = t2CurrentStreak
				t2CurrentStreak = 0
				t2Losses += 1
				if game.IsOvertime {
					t2OTLosses += 1
				}
				if game.IsShootout {
					t2SOLosses += 1
				}
			}
		}

		awayTeamTeamOne := game.AwayTeamID == teamOne.ID
		if awayTeamTeamOne {
			if game.AwayTeamWin {
				t1Wins += 1
				t1CurrentStreak += 1
				if game.IsOvertime {
					t1OTWins += 1
				}
				if game.IsShootout {
					t1SOWins += 1
				}
				latestWin = game.AwayTeam
				diff := game.AwayTeamScore - game.HomeTeamScore
				if diff > uint(t1LargestMarginDiff) {
					t1LargestMarginDiff = int(diff)
					t1LargestMarginSeason = int(game.SeasonID) + 2024
					t1LargestMarginScore = "" + strconv.Itoa(int(game.AwayTeamScore)) + "-" + strconv.Itoa(int(game.HomeTeamScore))
				}
			} else {
				t1Streak = t1CurrentStreak
				t1CurrentStreak = 0
				t1Losses += 1
				if game.IsOvertime {
					t1OTLosses += 1
				}
				if game.IsShootout {
					t1SOLosses += 1
				}
			}
		} else {
			if game.AwayTeamWin {
				t2Wins += 1
				t2CurrentStreak += 1
				if game.IsOvertime {
					t2OTWins += 1
				}
				if game.IsShootout {
					t2SOWins += 1
				}
				latestWin = game.AwayTeam
				diff := game.AwayTeamScore - game.HomeTeamScore
				if diff > uint(t2LargestMarginDiff) {
					t2LargestMarginDiff = int(diff)
					t2LargestMarginSeason = int(game.SeasonID) + 2024
					t2LargestMarginScore = "" + strconv.Itoa(int(game.AwayTeamScore)) + "-" + strconv.Itoa(int(game.HomeTeamScore))
				}
			} else {
				t2Streak = t2CurrentStreak
				t2CurrentStreak = 0
				t2Losses += 1
				if game.IsOvertime {
					t2OTLosses += 1
				}
				if game.IsShootout {
					t2SOLosses += 1
				}
			}
		}
	}

	if t1CurrentStreak > 0 && t1CurrentStreak > t1Streak {
		t1Streak = t1CurrentStreak
	}
	if t2CurrentStreak > 0 && t2CurrentStreak > t2Streak {
		t2Streak = t2CurrentStreak
	}

	currentStreak := 0
	currentStreak = max(t1CurrentStreak, t2CurrentStreak)

	return structs.TeamComparisonModel{
		TeamOneID:       teamOne.ID,
		TeamOne:         teamOne.TeamName,
		TeamOneWins:     uint(t1Wins),
		TeamOneLosses:   uint(t1Losses),
		TeamOneOTWins:   uint(t1OTWins),
		TeamOneOTLosses: uint(t1OTLosses),
		TeamOneSOWins:   uint(t1SOWins),
		TeamOneSOLosses: uint(t1SOLosses),
		TeamOneStreak:   uint(t1Streak),
		TeamOneMSeason:  t1LargestMarginSeason,
		TeamOneMScore:   t1LargestMarginScore,
		TeamTwoID:       teamTwo.ID,
		TeamTwo:         teamTwo.TeamName,
		TeamTwoWins:     uint(t2Wins),
		TeamTwoLosses:   uint(t2Losses),
		TeamTwoOTWins:   uint(t2OTWins),
		TeamTwoOTLosses: uint(t2OTLosses),
		TeamTwoSOWins:   uint(t2SOWins),
		TeamTwoSOLosses: uint(t2SOLosses),
		TeamTwoStreak:   uint(t2Streak),
		TeamTwoMSeason:  t2LargestMarginSeason,
		TeamTwoMScore:   t2LargestMarginScore,
		CurrentStreak:   uint(currentStreak),
		LatestWin:       latestWin,
	}
}

func ComparePHLTeams(t1ID, t2ID string) structs.TeamComparisonModel {
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	teamOneChan := make(chan structs.ProfessionalTeam)
	teamTwoChan := make(chan structs.ProfessionalTeam)

	go func() {
		t1 := GetProTeamByTeamID(t1ID)
		teamOneChan <- t1
	}()

	teamOne := <-teamOneChan
	close(teamOneChan)

	go func() {
		t2 := GetProTeamByTeamID(t2ID)
		teamTwoChan <- t2
	}()

	teamTwo := <-teamTwoChan
	close(teamTwoChan)

	allTeamOneGames := repository.FindProfessionalGames(repository.GamesClauses{SeasonID: seasonID, TeamID: t1ID, IsPreseason: false})

	t1Wins := 0
	t1Losses := 0
	t1OTWins := 0
	t1OTLosses := 0
	t1SOWins := 0
	t1SOLosses := 0
	t1Streak := 0
	t1CurrentStreak := 0
	t1LargestMarginSeason := 0
	t1LargestMarginDiff := 0
	t1LargestMarginScore := ""
	t2Wins := 0
	t2Losses := 0
	t2OTWins := 0
	t2OTLosses := 0
	t2SOWins := 0
	t2SOLosses := 0
	t2Streak := 0
	t2CurrentStreak := 0
	latestWin := ""
	t2LargestMarginSeason := 0
	t2LargestMarginDiff := 0
	t2LargestMarginScore := ""

	for _, game := range allTeamOneGames {
		gameNotRan := (game.GameDay == "A" && !ts.GamesARan) ||
			(game.GameDay == "B" && !ts.GamesBRan) ||
			(game.GameDay == "C" && !ts.GamesCRan) ||
			(game.GameDay == "D" && !ts.GamesDRan)
		if !game.GameComplete ||
			gameNotRan {
			continue
		}
		doComparison := (game.HomeTeamID == teamOne.ID && game.AwayTeamID == teamTwo.ID) ||
			(game.HomeTeamID == teamTwo.ID && game.AwayTeamID == teamOne.ID)

		if !doComparison {
			continue
		}
		homeTeamTeamOne := game.HomeTeamID == teamOne.ID
		if homeTeamTeamOne {
			if game.HomeTeamWin {
				t1Wins += 1
				t1CurrentStreak += 1
				if game.IsOvertime {
					t1OTWins += 1
				}
				if game.IsShootout {
					t1SOWins += 1
				}
				latestWin = game.HomeTeam
				diff := game.HomeTeamScore - game.AwayTeamScore
				if diff > uint(t1LargestMarginDiff) {
					t1LargestMarginDiff = int(diff)
					t1LargestMarginSeason = int(game.SeasonID) + 2024
					t1LargestMarginScore = "" + strconv.Itoa(int(game.HomeTeamScore)) + "-" + strconv.Itoa(int(game.AwayTeamScore))
				}
			} else {
				t1Streak = t1CurrentStreak
				t1CurrentStreak = 0
				t1Losses += 1
				if game.IsOvertime {
					t1OTLosses += 1
				}
				if game.IsShootout {
					t1SOLosses += 1
				}
			}
		} else {
			if game.HomeTeamWin {
				t2Wins += 1
				t2CurrentStreak += 1
				if game.IsOvertime {
					t2OTWins += 1
				}
				if game.IsShootout {
					t2SOWins += 1
				}
				latestWin = game.HomeTeam
				diff := game.HomeTeamScore - game.AwayTeamScore
				if diff > uint(t2LargestMarginDiff) {
					t2LargestMarginDiff = int(diff)
					t2LargestMarginSeason = int(game.SeasonID) + 2024
					t2LargestMarginScore = "" + strconv.Itoa(int(game.HomeTeamScore)) + "-" + strconv.Itoa(int(game.AwayTeamScore))
				}
			} else {
				t2Streak = t2CurrentStreak
				t2CurrentStreak = 0
				t2Losses += 1
				if game.IsOvertime {
					t2OTLosses += 1
				}
				if game.IsShootout {
					t2SOLosses += 1
				}
			}
		}

		awayTeamTeamOne := game.AwayTeamID == teamOne.ID
		if awayTeamTeamOne {
			if game.AwayTeamWin {
				t1Wins += 1
				t1CurrentStreak += 1
				if game.IsOvertime {
					t1OTWins += 1
				}
				if game.IsShootout {
					t1SOWins += 1
				}
				latestWin = game.AwayTeam
				diff := game.AwayTeamScore - game.HomeTeamScore
				if diff > uint(t1LargestMarginDiff) {
					t1LargestMarginDiff = int(diff)
					t1LargestMarginSeason = int(game.SeasonID) + 2024
					t1LargestMarginScore = "" + strconv.Itoa(int(game.AwayTeamScore)) + "-" + strconv.Itoa(int(game.HomeTeamScore))
				}
			} else {
				t1Streak = t1CurrentStreak
				t1CurrentStreak = 0
				t1Losses += 1
				if game.IsOvertime {
					t1OTLosses += 1
				}
				if game.IsShootout {
					t1SOLosses += 1
				}
			}
		} else {
			if game.AwayTeamWin {
				t2Wins += 1
				t2CurrentStreak += 1
				if game.IsOvertime {
					t2OTWins += 1
				}
				if game.IsShootout {
					t2SOWins += 1
				}
				latestWin = game.AwayTeam
				diff := game.AwayTeamScore - game.HomeTeamScore
				if diff > uint(t2LargestMarginDiff) {
					t2LargestMarginDiff = int(diff)
					t2LargestMarginSeason = int(game.SeasonID) + 2024
					t2LargestMarginScore = "" + strconv.Itoa(int(game.AwayTeamScore)) + "-" + strconv.Itoa(int(game.HomeTeamScore))
				}
			} else {
				t2Streak = t2CurrentStreak
				t2CurrentStreak = 0
				t2Losses += 1
				if game.IsOvertime {
					t2OTLosses += 1
				}
				if game.IsShootout {
					t2SOLosses += 1
				}
			}
		}
	}

	if t1CurrentStreak > 0 && t1CurrentStreak > t1Streak {
		t1Streak = t1CurrentStreak
	}
	if t2CurrentStreak > 0 && t2CurrentStreak > t2Streak {
		t2Streak = t2CurrentStreak
	}

	currentStreak := 0
	currentStreak = max(t1CurrentStreak, t2CurrentStreak)

	return structs.TeamComparisonModel{
		TeamOneID:       teamOne.ID,
		TeamOne:         teamOne.TeamName,
		TeamOneWins:     uint(t1Wins),
		TeamOneLosses:   uint(t1Losses),
		TeamOneOTWins:   uint(t1OTWins),
		TeamOneOTLosses: uint(t1OTLosses),
		TeamOneSOWins:   uint(t1SOWins),
		TeamOneSOLosses: uint(t1SOLosses),
		TeamOneStreak:   uint(t1Streak),
		TeamOneMSeason:  t1LargestMarginSeason,
		TeamOneMScore:   t1LargestMarginScore,
		TeamTwoID:       teamTwo.ID,
		TeamTwo:         teamTwo.TeamName,
		TeamTwoWins:     uint(t2Wins),
		TeamTwoLosses:   uint(t2Losses),
		TeamTwoOTWins:   uint(t2OTWins),
		TeamTwoOTLosses: uint(t2OTLosses),
		TeamTwoSOWins:   uint(t2SOWins),
		TeamTwoSOLosses: uint(t2SOLosses),
		TeamTwoStreak:   uint(t2Streak),
		TeamTwoMSeason:  t2LargestMarginSeason,
		TeamTwoMScore:   t2LargestMarginScore,
		CurrentStreak:   uint(currentStreak),
		LatestWin:       latestWin,
	}
}

func AssignDiscordIDToCHLTeam(tID, dID string) {
	db := dbprovider.GetInstance().GetDB()

	team := GetCollegeTeamByTeamID(tID)

	team.AssignDiscordID(dID)

	repository.SaveCollegeTeamRecord(db, team)
}

func AssignDiscordIDToPHLTeam(tID, dID string) {
	db := dbprovider.GetInstance().GetDB()

	team := GetProTeamByTeamID(tID)

	team.AssignDiscordID(dID)

	repository.SaveProTeamRecord(db, team)
}

func GetCollegeRecruitViaDiscord(id string) structs.Croot {
	db := dbprovider.GetInstance().GetDB()

	var recruit structs.Recruit

	err := db.Preload("RecruitPlayerProfiles").Where("id = ?", id).Find(&recruit).Error
	if err != nil {
		log.Fatalln(err)
	}

	var croot structs.Croot

	croot.Map(recruit)

	return croot
}
