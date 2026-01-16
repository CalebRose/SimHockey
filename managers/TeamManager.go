package managers

import (
	"sort"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetCollegeTeamByTeamID(teamID string) structs.CollegeTeam {
	return repository.FindCollegeTeamRecord(teamID)
}

func GetAllCollegeTeams() []structs.CollegeTeam {
	return repository.FindAllCollegeTeams(repository.TeamClauses{LeagueID: "1"})
}

func GetAllCanadianCHLTeams() []structs.CollegeTeam {
	return repository.FindAllCollegeTeams(repository.TeamClauses{LeagueID: "2"})
}

func GetCollegeTeamMap() map[uint]structs.CollegeTeam {
	teams := repository.FindAllCollegeTeams(repository.TeamClauses{LeagueID: "1"})
	return MakeCollegeTeamMap(teams)
}

func GetProTeamByTeamID(teamID string) structs.ProfessionalTeam {
	return repository.FindProTeamRecord(teamID)
}

func GetAllProfessionalTeams() []structs.ProfessionalTeam {
	return repository.FindAllProTeams(repository.TeamClauses{LeagueID: "1"})
}

func GetProTeamMap() map[uint]structs.ProfessionalTeam {
	teams := repository.FindAllProTeams(repository.TeamClauses{LeagueID: "1"})
	return MakeProTeamMap(teams)
}

// Helper functions for team grade calculations

// calculateOffenseGrade calculates the offense grade based on top 3 centers and top 6 forwards
func calculateOffenseGrade(players []structs.BasePlayer) float64 {
	var centers []structs.BasePlayer
	var forwards []structs.BasePlayer

	// Separate players by position
	for _, player := range players {
		switch player.Position {
		case "C":
			centers = append(centers, player)
		case "F":
			forwards = append(forwards, player)
		}
	}

	// Sort by overall rating (descending)
	sort.Slice(centers, func(i, j int) bool {
		return centers[i].Overall > centers[j].Overall
	})

	sort.Slice(forwards, func(i, j int) bool {
		return forwards[i].Overall > forwards[j].Overall
	})

	// Get top 3 centers and top 6 forwards
	topCenters := getTopPlayers(centers, 3)
	topForwards := getTopPlayers(forwards, 6)

	// Calculate average offensive rating
	var totalOffenseRating float64
	var playerCount int

	for _, center := range topCenters {
		offenseRating := calculatePlayerOffenseRating(center)
		totalOffenseRating += offenseRating
		playerCount++
	}

	for _, forward := range topForwards {
		offenseRating := calculatePlayerOffenseRating(forward)
		totalOffenseRating += offenseRating
		playerCount++
	}

	if playerCount == 0 {
		return 0
	}

	return totalOffenseRating / float64(playerCount)
}

// calculateDefenseGrade calculates the defense grade based on top 4 defenders
func calculateDefenseGrade(players []structs.BasePlayer) float64 {
	var defenders []structs.BasePlayer

	// Get all defenders
	for _, player := range players {
		if player.Position == "D" {
			defenders = append(defenders, player)
		}
	}

	// Sort by overall rating (descending)
	sort.Slice(defenders, func(i, j int) bool {
		return defenders[i].Overall > defenders[j].Overall
	})

	// Get top 4 defenders
	topDefenders := getTopPlayers(defenders, 4)

	// Calculate average defensive rating
	var totalDefenseRating float64
	for _, defender := range topDefenders {
		defenseRating := calculatePlayerDefenseRating(defender)
		totalDefenseRating += defenseRating
	}

	if len(topDefenders) == 0 {
		return 0
	}

	return totalDefenseRating / float64(len(topDefenders))
}

// calculateGoalieGrade calculates the goalie grade based on top 2 goalies
func calculateGoalieGrade(players []structs.BasePlayer) float64 {
	var goalies []structs.BasePlayer

	// Get all goalies
	for _, player := range players {
		if player.Position == "G" {
			goalies = append(goalies, player)
		}
	}

	// Sort by overall rating (descending)
	sort.Slice(goalies, func(i, j int) bool {
		return goalies[i].Overall > goalies[j].Overall
	})

	// Get top 2 goalies
	topGoalies := getTopPlayers(goalies, 2)

	// Calculate average goalie rating
	var totalGoalieRating float64
	for _, goalie := range topGoalies {
		goalieRating := calculatePlayerGoalieRating(goalie)
		totalGoalieRating += goalieRating
	}

	if len(topGoalies) == 0 {
		return 0
	}

	return totalGoalieRating / float64(len(topGoalies))
}

// calculateOverallGrade calculates overall grade with 40% offense, 40% defense, 20% goalie
func calculateOverallGrade(offenseGrade, defenseGrade, goalieGrade float64) float64 {
	return (offenseGrade * 0.4) + (defenseGrade * 0.4) + (goalieGrade * 0.2)
}

// calculatePlayerOffenseRating calculates a player's offensive contribution
func calculatePlayerOffenseRating(player structs.BasePlayer) float64 {
	// Weight offensive attributes: shooting, passing, puck handling, agility
	offenseRating := float64(player.CloseShotAccuracy)*0.25 +
		float64(player.CloseShotPower)*0.15 +
		float64(player.LongShotAccuracy)*0.20 +
		float64(player.LongShotPower)*0.15 +
		float64(player.Passing)*0.15 +
		float64(player.PuckHandling)*0.10

	// Add faceoff bonus for centers
	if player.Position == "C" {
		offenseRating = offenseRating*0.9 + float64(player.Faceoffs)*0.1
	}

	return offenseRating
}

// calculatePlayerDefenseRating calculates a player's defensive contribution
func calculatePlayerDefenseRating(player structs.BasePlayer) float64 {
	// Weight defensive attributes: checking, shot blocking, strength, agility
	return float64(player.BodyChecking)*0.3 +
		float64(player.StickChecking)*0.3 +
		float64(player.ShotBlocking)*0.25 +
		float64(player.Strength)*0.10 +
		float64(player.Agility)*0.05
}

// calculatePlayerGoalieRating calculates a goalie's contribution
func calculatePlayerGoalieRating(player structs.BasePlayer) float64 {
	// Weight goalie attributes: goalkeeping, vision, rebound control
	return float64(player.Goalkeeping)*0.6 +
		float64(player.GoalieVision)*0.25 +
		float64(player.GoalieReboundControl)*0.15
}

// getTopPlayers returns the top N players from a sorted list
func getTopPlayers(players []structs.BasePlayer, count int) []structs.BasePlayer {
	if len(players) <= count {
		return players
	}
	return players[:count]
}

// calculateMeanAndStdDev calculates mean and standard deviation for a specific grade type
func calculateMeanAndStdDev(teamGrades []structs.TeamGrade, gradeType string) (float64, float64) {
	if len(teamGrades) == 0 {
		return 0, 0
	}

	var sum float64
	for _, tg := range teamGrades {
		switch gradeType {
		case "offense":
			sum += tg.OffenseGradeNumber
		case "defense":
			sum += tg.DefenseGradeNumber
		case "goalie":
			sum += tg.GoalieGradeNumber
		case "overall":
			sum += tg.OverallGradeNumber
		}
	}

	mean := sum / float64(len(teamGrades))

	// Calculate variance
	var variance float64
	for _, tg := range teamGrades {
		var value float64
		switch gradeType {
		case "offense":
			value = tg.OffenseGradeNumber
		case "defense":
			value = tg.DefenseGradeNumber
		case "goalie":
			value = tg.GoalieGradeNumber
		case "overall":
			value = tg.OverallGradeNumber
		}
		diff := value - mean
		variance += diff * diff
	}

	variance /= float64(len(teamGrades))
	stdDev := 0.0
	if variance > 0 {
		// Simple square root approximation using Newton's method
		stdDev = variance / 2
		for i := 0; i < 10; i++ {
			stdDev = (stdDev + variance/stdDev) / 2
		}
	}

	return mean, stdDev
}

// assignLetterGrade assigns letter grade based on standard deviations from mean
func assignLetterGrade(value, mean, stdDev float64) string {
	if stdDev == 0 {
		return "C"
	}

	z := (value - mean) / stdDev

	switch {
	case z >= 2.0:
		return "A+"
	case z >= 1.75:
		return "A"
	case z >= 1.5:
		return "A-"
	case z >= 1.25:
		return "B+"
	case z >= 1.0:
		return "B"
	case z >= 0.75:
		return "B-"
	case z >= 0.5:
		return "C+"
	case z >= -0.5:
		return "C"
	case z >= -0.75:
		return "C-"
	case z >= -1.0:
		return "D+"
	case z >= -1.5:
		return "D"
	case z >= -2.0:
		return "D-"
	default:
		return "F"
	}
}

func GenerateCHLTeamLetterGrades() {
	db := dbprovider.GetInstance().GetDB()

	teams := repository.FindAllCollegeTeams(repository.TeamClauses{LeagueID: "1"})
	collegePlayers := repository.FindAllCollegePlayers(repository.PlayerQuery{LeagueID: "1"})
	collegeRosterMap := MakeCollegePlayerMapByTeamID(collegePlayers)

	// Calculate team grades for all teams
	var teamGrades []structs.TeamGrade
	for _, team := range teams {
		players := collegeRosterMap[team.ID]

		basePlayers := MakeBasePlayerList(players, []structs.ProfessionalPlayer{})

		// Calculate numeric grades for this team
		offenseGrade := calculateOffenseGrade(basePlayers)
		defenseGrade := calculateDefenseGrade(basePlayers)
		goalieGrade := calculateGoalieGrade(basePlayers)
		overallGrade := calculateOverallGrade(offenseGrade, defenseGrade, goalieGrade)

		// Create team grade struct
		tg := structs.TeamGrade{}
		tg.SetOffenseGradeNumber(offenseGrade)
		tg.SetDefenseGradeNumber(defenseGrade)
		tg.SetGoalieGradeNumber(goalieGrade)
		tg.SetOverallGradeNumber(overallGrade)

		teamGrades = append(teamGrades, tg)
	}

	// Calculate statistical distribution for letter grade assignment
	offenseMean, offenseStdDev := calculateMeanAndStdDev(teamGrades, "offense")
	defenseMean, defenseStdDev := calculateMeanAndStdDev(teamGrades, "defense")
	goalieMean, goalieStdDev := calculateMeanAndStdDev(teamGrades, "goalie")
	overallMean, overallStdDev := calculateMeanAndStdDev(teamGrades, "overall")

	// Assign letter grades and update teams
	for i, team := range teams {
		teamGrade := &teamGrades[i]

		// Assign letter grades based on standard deviation from mean
		offenseLetter := assignLetterGrade(teamGrade.OffenseGradeNumber, offenseMean, offenseStdDev)
		defenseLetter := assignLetterGrade(teamGrade.DefenseGradeNumber, defenseMean, defenseStdDev)
		goalieLetter := assignLetterGrade(teamGrade.GoalieGradeNumber, goalieMean, goalieStdDev)
		overallLetter := assignLetterGrade(teamGrade.OverallGradeNumber, overallMean, overallStdDev)

		teamGrade.SetOffenseGradeLetter(offenseLetter)
		teamGrade.SetDefenseGradeLetter(defenseLetter)
		teamGrade.SetGoalieGradeLetter(goalieLetter)
		teamGrade.SetOverallGradeLetter(overallLetter)

		// Update the team with letter grades
		teamToUpdate := team
		teamToUpdate.AssignLetterGrades(overallLetter, offenseLetter, defenseLetter, goalieLetter)

		repository.SaveCollegeTeamRecord(db, teamToUpdate)
	}
}

func GeneratePHLTeamLetterGrades() {
	db := dbprovider.GetInstance().GetDB()

	teams := repository.FindAllProTeams(repository.TeamClauses{})
	proPlayers := repository.FindAllProPlayers(repository.PlayerQuery{})
	proRosterMap := MakeProfessionalPlayerMapByTeamID(proPlayers)

	// Calculate team grades for all teams
	var teamGrades []structs.TeamGrade
	for _, team := range teams {
		players := proRosterMap[team.ID]

		basePlayers := MakeBasePlayerList([]structs.CollegePlayer{}, players)

		// Calculate numeric grades for this team
		offenseGrade := calculateOffenseGrade(basePlayers)
		defenseGrade := calculateDefenseGrade(basePlayers)
		goalieGrade := calculateGoalieGrade(basePlayers)
		overallGrade := calculateOverallGrade(offenseGrade, defenseGrade, goalieGrade)

		// Create team grade struct
		tg := structs.TeamGrade{}
		tg.SetOffenseGradeNumber(offenseGrade)
		tg.SetDefenseGradeNumber(defenseGrade)
		tg.SetGoalieGradeNumber(goalieGrade)
		tg.SetOverallGradeNumber(overallGrade)

		teamGrades = append(teamGrades, tg)
	}

	// Calculate statistical distribution for letter grade assignment
	offenseMean, offenseStdDev := calculateMeanAndStdDev(teamGrades, "offense")
	defenseMean, defenseStdDev := calculateMeanAndStdDev(teamGrades, "defense")
	goalieMean, goalieStdDev := calculateMeanAndStdDev(teamGrades, "goalie")
	overallMean, overallStdDev := calculateMeanAndStdDev(teamGrades, "overall")

	// Assign letter grades and update teams
	for i, team := range teams {
		teamGrade := &teamGrades[i]

		// Assign letter grades based on standard deviation from mean
		offenseLetter := assignLetterGrade(teamGrade.OffenseGradeNumber, offenseMean, offenseStdDev)
		defenseLetter := assignLetterGrade(teamGrade.DefenseGradeNumber, defenseMean, defenseStdDev)
		goalieLetter := assignLetterGrade(teamGrade.GoalieGradeNumber, goalieMean, goalieStdDev)
		overallLetter := assignLetterGrade(teamGrade.OverallGradeNumber, overallMean, overallStdDev)

		teamGrade.SetOffenseGradeLetter(offenseLetter)
		teamGrade.SetDefenseGradeLetter(defenseLetter)
		teamGrade.SetGoalieGradeLetter(goalieLetter)
		teamGrade.SetOverallGradeLetter(overallLetter)

		// Update the team with letter grades
		teamToUpdate := team
		teamToUpdate.AssignLetterGrades(overallLetter, offenseLetter, defenseLetter, goalieLetter)

		repository.SaveProTeamRecord(db, teamToUpdate)
	}
}
