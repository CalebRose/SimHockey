package managers

import "github.com/CalebRose/SimHockey/structs"

func MakeCollegePlayerMap(players []structs.CollegePlayer) map[uint]structs.CollegePlayer {
	playerMap := make(map[uint]structs.CollegePlayer)

	for _, p := range players {
		playerMap[p.ID] = p
	}

	return playerMap
}

func MakeCollegePlayerMapByTeamID(players []structs.CollegePlayer) map[uint][]structs.CollegePlayer {
	playerMap := make(map[uint][]structs.CollegePlayer)

	for _, p := range players {
		if len(playerMap[uint(p.TeamID)]) > 0 {
			playerMap[uint(p.TeamID)] = append(playerMap[uint(p.TeamID)], p)
		} else {
			playerMap[uint(p.TeamID)] = []structs.CollegePlayer{p}
		}
	}

	return playerMap
}

func MakeCollegeRecruitMapByTeamID(players []structs.Recruit) map[uint][]structs.Recruit {
	playerMap := make(map[uint][]structs.Recruit)

	for _, p := range players {
		if len(playerMap[uint(p.TeamID)]) > 0 {
			playerMap[uint(p.TeamID)] = append(playerMap[uint(p.TeamID)], p)
		} else {
			playerMap[uint(p.TeamID)] = []structs.Recruit{p}
		}
	}

	return playerMap
}

// Pro Players
func MakeProfessionalPlayerMap(players []structs.ProfessionalPlayer) map[uint]structs.ProfessionalPlayer {
	playerMap := make(map[uint]structs.ProfessionalPlayer)

	for _, p := range players {
		playerMap[p.ID] = p
	}

	return playerMap
}

func MakeProfessionalPlayerMapByTeamID(players []structs.ProfessionalPlayer) map[uint][]structs.ProfessionalPlayer {
	playerMap := make(map[uint][]structs.ProfessionalPlayer)

	for _, p := range players {
		if len(playerMap[uint(p.TeamID)]) > 0 {
			playerMap[uint(p.TeamID)] = append(playerMap[uint(p.TeamID)], p)
		} else {
			playerMap[uint(p.TeamID)] = []structs.ProfessionalPlayer{p}
		}
	}

	return playerMap
}

func MakeArenaMap(arenas []structs.Arena) map[uint]structs.Arena {
	arenaMap := make(map[uint]structs.Arena)

	for _, a := range arenas {
		arenaMap[a.ID] = a
	}

	return arenaMap
}

func MakeCollegeTeamMap(collegeTeams []structs.CollegeTeam) map[uint]structs.CollegeTeam {
	teamMap := make(map[uint]structs.CollegeTeam)
	for _, t := range collegeTeams {
		teamMap[t.ID] = t
	}
	return teamMap
}

func MakeCollegeLineupMap(lineups []structs.CollegeLineup) map[uint][]structs.CollegeLineup {
	lineupMap := make(map[uint][]structs.CollegeLineup)

	for _, l := range lineups {
		if len(lineupMap[uint(l.TeamID)]) > 0 {
			lineupMap[uint(l.TeamID)] = append(lineupMap[uint(l.TeamID)], l)
		} else {
			lineupMap[uint(l.TeamID)] = []structs.CollegeLineup{l}
		}
	}

	return lineupMap
}

func MakeIndCollegeLineupMap(lineups []structs.CollegeLineup) map[uint]structs.CollegeLineup {
	lineupMap := make(map[uint]structs.CollegeLineup)

	for _, l := range lineups {
		lineupMap[l.ID] = l
	}

	return lineupMap
}

func MakeCollegeShootoutLineupMap(lineups []structs.CollegeShootoutLineup) map[uint]structs.CollegeShootoutLineup {
	lineupMap := make(map[uint]structs.CollegeShootoutLineup)

	for _, l := range lineups {
		lineupMap[uint(l.TeamID)] = l
	}

	return lineupMap
}

func MakeProTeamMap(teams []structs.ProfessionalTeam) map[uint]structs.ProfessionalTeam {
	teamMap := make(map[uint]structs.ProfessionalTeam)
	for _, t := range teams {
		teamMap[t.ID] = t
	}
	return teamMap
}

func MakeProfessionalLineupMap(lineups []structs.ProfessionalLineup) map[uint][]structs.ProfessionalLineup {
	lineupMap := make(map[uint][]structs.ProfessionalLineup)

	for _, l := range lineups {
		if len(lineupMap[uint(l.TeamID)]) > 0 {
			lineupMap[uint(l.TeamID)] = append(lineupMap[uint(l.TeamID)], l)
		} else {
			lineupMap[uint(l.TeamID)] = []structs.ProfessionalLineup{l}
		}
	}

	return lineupMap
}

func MakeIndProLineupMap(lineups []structs.ProfessionalLineup) map[uint]structs.ProfessionalLineup {
	lineupMap := make(map[uint]structs.ProfessionalLineup)

	for _, l := range lineups {
		lineupMap[l.ID] = l
	}

	return lineupMap
}

func MakeProfessionalShootoutLineupMap(lineups []structs.ProfessionalShootoutLineup) map[uint]structs.ProfessionalShootoutLineup {
	lineupMap := make(map[uint]structs.ProfessionalShootoutLineup)

	for _, l := range lineups {
		lineupMap[uint(l.TeamID)] = l
	}

	return lineupMap
}

func MakeCollegeStandingsMap(standings []structs.CollegeStandings) map[uint]structs.CollegeStandings {
	standingsMap := make(map[uint]structs.CollegeStandings)
	for _, stat := range standings {
		standingsMap[uint(stat.TeamID)] = stat
	}

	return standingsMap
}

func MakeCollegeStandingsMapByCoachName(standings []structs.CollegeStandings) map[string][]structs.CollegeStandings {
	standingsMap := make(map[string][]structs.CollegeStandings)
	for _, stat := range standings {
		if len(standingsMap[stat.Coach]) == 0 {
			standingsMap[stat.Coach] = []structs.CollegeStandings{}
		}
		standingsMap[stat.Coach] = append(standingsMap[stat.Coach], stat)
	}

	return standingsMap
}

func MakeProfessionalStandingsMap(standings []structs.ProfessionalStandings) map[uint]structs.ProfessionalStandings {
	standingsMap := make(map[uint]structs.ProfessionalStandings)
	for _, stat := range standings {
		standingsMap[uint(stat.TeamID)] = stat
	}

	return standingsMap
}

func MakeRecruitProfileMapByRecruitID(profiles []structs.RecruitPlayerProfile) map[uint][]structs.RecruitPlayerProfile {
	profileMap := make(map[uint][]structs.RecruitPlayerProfile)

	for _, p := range profiles {
		if len(profileMap[uint(p.RecruitID)]) > 0 {
			profileMap[uint(p.RecruitID)] = append(profileMap[uint(p.RecruitID)], p)
		} else {
			profileMap[uint(p.RecruitID)] = []structs.RecruitPlayerProfile{p}
		}
	}

	return profileMap
}

func MakeRecruitProfileMapByProfileID(profiles []structs.RecruitPlayerProfile) map[uint][]structs.RecruitPlayerProfile {
	profileMap := make(map[uint][]structs.RecruitPlayerProfile)

	for _, p := range profiles {
		if len(profileMap[uint(p.ProfileID)]) > 0 {
			profileMap[uint(p.ProfileID)] = append(profileMap[uint(p.ProfileID)], p)
		} else {
			profileMap[uint(p.ProfileID)] = []structs.RecruitPlayerProfile{p}
		}
	}

	return profileMap
}

func MakeTeamProfileMap(profiles []structs.RecruitingTeamProfile) map[uint]*structs.RecruitingTeamProfile {
	profileMap := make(map[uint]*structs.RecruitingTeamProfile)

	for _, p := range profiles {
		profileMap[uint(p.ID)] = &p
	}

	return profileMap
}

func MakeCapsheetMap(capsheets []structs.ProCapsheet) map[uint]structs.ProCapsheet {
	capsheetMap := make(map[uint]structs.ProCapsheet)

	for _, p := range capsheets {
		capsheetMap[uint(p.ID)] = p
	}

	return capsheetMap
}

func MakeContractMap(contracts []structs.ProContract) map[uint]structs.ProContract {
	contractMap := make(map[uint]structs.ProContract)

	for _, c := range contracts {
		contractMap[uint(c.PlayerID)] = c
	}

	return contractMap
}

func MakeExtensionMap(extensions []structs.ExtensionOffer) map[uint]structs.ExtensionOffer {
	contractMap := make(map[uint]structs.ExtensionOffer)

	for _, c := range extensions {
		contractMap[uint(c.PlayerID)] = c
	}

	return contractMap
}

func MakeCollegeGameMapByTeamID(games []structs.CollegeGame) map[uint][]structs.CollegeGame {
	gameMap := make(map[uint][]structs.CollegeGame)

	for _, p := range games {
		if _, exists := gameMap[uint(p.HomeTeamID)]; !exists {
			gameMap[uint(p.HomeTeamID)] = []structs.CollegeGame{}
		}
		gameMap[uint(p.HomeTeamID)] = append(gameMap[uint(p.HomeTeamID)], p)
		if _, exists := gameMap[uint(p.AwayTeamID)]; !exists {
			gameMap[uint(p.AwayTeamID)] = []structs.CollegeGame{}
		}
		gameMap[uint(p.AwayTeamID)] = append(gameMap[uint(p.AwayTeamID)], p)
	}
	return gameMap
}

func MakeCollegeGameMap(games []structs.CollegeGame) map[uint]structs.CollegeGame {
	gameMap := make(map[uint]structs.CollegeGame)

	for _, p := range games {
		gameMap[p.ID] = p
	}

	return gameMap
}

func MakeProGameMap(games []structs.ProfessionalGame) map[uint]structs.ProfessionalGame {
	gameMap := make(map[uint]structs.ProfessionalGame)

	for _, p := range games {
		gameMap[p.ID] = p
	}

	return gameMap
}

func MakeCollegePlayerSeasonStatMap(stats []structs.CollegePlayerSeasonStats) map[uint]structs.CollegePlayerSeasonStats {
	seasonStatMap := make(map[uint]structs.CollegePlayerSeasonStats)
	for _, stat := range stats {
		seasonStatMap[stat.PlayerID] = stat
	}

	return seasonStatMap
}

func MakeProPlayerSeasonStatMap(stats []structs.ProfessionalPlayerSeasonStats) map[uint]structs.ProfessionalPlayerSeasonStats {
	seasonStatMap := make(map[uint]structs.ProfessionalPlayerSeasonStats)

	for _, stat := range stats {
		seasonStatMap[stat.PlayerID] = stat
	}

	return seasonStatMap
}

func MakeHistoricProPlayerSeasonStatMap(stats []structs.ProfessionalPlayerSeasonStats) map[uint][]structs.ProfessionalPlayerSeasonStats {
	seasonStatMap := make(map[uint][]structs.ProfessionalPlayerSeasonStats)

	for _, stat := range stats {
		if len(seasonStatMap[stat.PlayerID]) == 0 {
			seasonStatMap[stat.PlayerID] = []structs.ProfessionalPlayerSeasonStats{}
		}
		seasonStatMap[stat.PlayerID] = append(seasonStatMap[stat.PlayerID], stat)
	}

	return seasonStatMap
}

func MakeCollegeTeamSeasonStatMap(stats []structs.CollegeTeamSeasonStats) map[uint]structs.CollegeTeamSeasonStats {
	seasonStatMap := make(map[uint]structs.CollegeTeamSeasonStats)
	for _, stat := range stats {
		seasonStatMap[stat.TeamID] = stat
	}

	return seasonStatMap
}

func MakeProTeamSeasonStatMap(stats []structs.ProfessionalTeamSeasonStats) map[uint]structs.ProfessionalTeamSeasonStats {
	seasonStatMap := make(map[uint]structs.ProfessionalTeamSeasonStats)

	for _, stat := range stats {
		seasonStatMap[stat.TeamID] = stat
	}

	return seasonStatMap
}

func MakeFreeAgencyOfferMap(offers []structs.FreeAgencyOffer) map[uint][]structs.FreeAgencyOffer {
	offerMap := make(map[uint][]structs.FreeAgencyOffer)

	for _, offer := range offers {
		if len(offerMap[offer.PlayerID]) > 0 {
			offerMap[offer.PlayerID] = append(offerMap[uint(offer.PlayerID)], offer)
		} else {
			offerMap[offer.PlayerID] = []structs.FreeAgencyOffer{offer}
		}
	}

	return offerMap
}

func MakeFreeAgencyOfferMapByTeamID(offers []structs.FreeAgencyOffer) map[uint][]structs.FreeAgencyOffer {
	offerMap := make(map[uint][]structs.FreeAgencyOffer)

	for _, offer := range offers {
		if len(offerMap[offer.TeamID]) > 0 {
			offerMap[offer.TeamID] = append(offerMap[uint(offer.TeamID)], offer)
		} else {
			offerMap[offer.TeamID] = []structs.FreeAgencyOffer{offer}
		}
	}

	return offerMap
}

func MakeTradePreferencesMap(prefs []structs.TradePreferences) map[uint]structs.TradePreferences {
	prefMap := make(map[uint]structs.TradePreferences)

	for _, pref := range prefs {
		prefMap[pref.TeamID] = pref
	}

	return prefMap
}

func MakeTradeProposalMap(proposals []structs.TradeProposal) map[uint][]structs.TradeProposal {
	proposalMap := make(map[uint][]structs.TradeProposal)

	for _, proposal := range proposals {
		if len(proposalMap[proposal.TeamID]) > 0 {
			proposalMap[proposal.TeamID] = append(proposalMap[uint(proposal.TeamID)], proposal)
		} else {
			proposalMap[proposal.TeamID] = []structs.TradeProposal{proposal}
		}
	}

	return proposalMap
}

func MakeCollegeGameplanMap(gameplans []structs.CollegeGameplan) map[uint]structs.CollegeGameplan {
	gameplanMap := make(map[uint]structs.CollegeGameplan)

	for _, gameplan := range gameplans {
		gameplanMap[gameplan.TeamID] = gameplan
	}

	return gameplanMap
}

func MakeProGameplanMap(gameplans []structs.ProGameplan) map[uint]structs.ProGameplan {
	gameplanMap := make(map[uint]structs.ProGameplan)

	for _, gameplan := range gameplans {
		gameplanMap[gameplan.TeamID] = gameplan
	}

	return gameplanMap
}

func MakeCollegePlayerGameStatsMap(stats []structs.CollegePlayerGameStats) map[uint][]structs.CollegePlayerGameStats {
	statsMap := make(map[uint][]structs.CollegePlayerGameStats)

	for _, stat := range stats {
		if len(statsMap[stat.PlayerID]) > 0 {
			statsMap[stat.PlayerID] = append(statsMap[uint(stat.PlayerID)], stat)
		} else {
			statsMap[stat.PlayerID] = []structs.CollegePlayerGameStats{stat}
		}
	}

	return statsMap
}

func MakeProPlayerGameStatsMap(stats []structs.ProfessionalPlayerGameStats) map[uint][]structs.ProfessionalPlayerGameStats {
	statsMap := make(map[uint][]structs.ProfessionalPlayerGameStats)

	for _, stat := range stats {
		if len(statsMap[stat.PlayerID]) > 0 {
			statsMap[stat.PlayerID] = append(statsMap[uint(stat.PlayerID)], stat)
		} else {
			statsMap[stat.PlayerID] = []structs.ProfessionalPlayerGameStats{stat}
		}
	}

	return statsMap
}

func MakePortalProfileMapByTeamID(profiles []structs.TransferPortalProfile) map[uint][]structs.TransferPortalProfile {
	playerMap := make(map[uint][]structs.TransferPortalProfile)

	for _, p := range profiles {
		if p.ProfileID == 0 {
			continue
		}
		if len(playerMap[uint(p.ProfileID)]) > 0 {
			playerMap[uint(p.ProfileID)] = append(playerMap[uint(p.ProfileID)], p)
		} else {
			playerMap[uint(p.ProfileID)] = []structs.TransferPortalProfile{p}
		}
	}

	return playerMap
}

func MakeCollegePromiseMap(promises []structs.CollegePromise) map[uint]structs.CollegePromise {
	playerMap := make(map[uint]structs.CollegePromise)

	for _, p := range promises {
		playerMap[p.ID] = p
	}

	return playerMap
}

func MakePromiseMapByTeamID(profiles []structs.CollegePromise) map[uint][]structs.CollegePromise {
	playerMap := make(map[uint][]structs.CollegePromise)

	for _, p := range profiles {
		if p.TeamID == 0 || !p.IsActive {
			continue
		}
		if len(playerMap[uint(p.TeamID)]) > 0 {
			playerMap[uint(p.TeamID)] = append(playerMap[uint(p.TeamID)], p)
		} else {
			playerMap[uint(p.TeamID)] = []structs.CollegePromise{p}
		}
	}

	return playerMap
}

func MakePromiseMapByPlayerIDByTeam(promises []structs.CollegePromise) map[uint]structs.CollegePromise {
	playerMap := make(map[uint]structs.CollegePromise)

	for _, p := range promises {
		playerMap[p.CollegePlayerID] = p
	}

	return playerMap
}

func MakePortalProfileMapByPlayerID(profiles []structs.TransferPortalProfile) map[uint][]structs.TransferPortalProfile {
	playerMap := make(map[uint][]structs.TransferPortalProfile)

	for _, p := range profiles {
		if p.ProfileID == 0 {
			continue
		}
		if len(playerMap[uint(p.CollegePlayerID)]) > 0 {
			playerMap[uint(p.CollegePlayerID)] = append(playerMap[uint(p.CollegePlayerID)], p)
		} else {
			playerMap[uint(p.CollegePlayerID)] = []structs.TransferPortalProfile{p}
		}
	}

	return playerMap
}
