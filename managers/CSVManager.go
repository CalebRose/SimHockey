package managers

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/engine"
	"github.com/CalebRose/SimHockey/structs"
)

func WritePlayByPlayCSVFile(playByPlays []structs.PbP, filename string, collegePlayerMap map[uint]structs.CollegePlayer, collegeTeamMap map[uint]structs.CollegeTeam) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filename, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Period", "TimeOnClock", "Time Consumed", "Zone", "Event", "Outcome", "Penalty Called", "Severity", "Fight?", "HTS", "ATS", "PossessingTeam", "Notes"})
	// Iterate through play by play data to generate []string

	for _, play := range playByPlays {
		periodStr := strconv.Itoa(int(play.Period))
		timeOnClock := FormatTimeToClock(play.TimeOnClock)
		timeConsumed := strconv.Itoa(int(play.SecondsConsumed))
		event := util.ReturnStringFromPBPID(play.EventID)
		outcome := util.ReturnStringFromPBPID(play.Outcome)
		hts := strconv.Itoa(int(play.HomeTeamScore))
		ats := strconv.Itoa(int(play.AwayTeamScore))
		possessingTeam := collegeTeamMap[uint(play.TeamID)]
		zone := getZoneLabel(play.ZoneID)
		abbr := possessingTeam.Abbreviation
		penalty := getPenaltyByID(uint(play.PenaltyID))
		severity := getSeverityByID(play.Severity)
		isFight := "No"
		if play.IsFight {
			isFight = "Yes"
		}

		result := generateCollegeResultsString(play, event, outcome, collegePlayerMap, possessingTeam)
		writer.Write([]string{
			periodStr,
			timeOnClock,
			timeConsumed,
			zone,
			event,
			outcome,
			penalty,
			severity,
			isFight,
			hts,
			ats,
			abbr,
			result,
		})
	}
	return err
}

func WriteProPlayByPlayCSVFile(playByPlays []structs.PbP, filename string, playerMap map[uint]structs.ProfessionalPlayer, collegeTeamMap map[uint]structs.ProfessionalTeam) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filename, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Period", "TimeOnClock", "Time Consumed", "Zone", "Event", "Outcome", "Penalty Called", "Severity", "Fight?", "HTS", "ATS", "PossessingTeam", "Notes"})
	// Iterate through play by play data to generate []string

	for _, play := range playByPlays {
		periodStr := strconv.Itoa(int(play.Period))
		timeOnClock := FormatTimeToClock(play.TimeOnClock)
		timeConsumed := strconv.Itoa(int(play.SecondsConsumed))
		event := util.ReturnStringFromPBPID(play.EventID)
		outcome := util.ReturnStringFromPBPID(play.Outcome)
		hts := strconv.Itoa(int(play.HomeTeamScore))
		ats := strconv.Itoa(int(play.AwayTeamScore))
		possessingTeam := collegeTeamMap[uint(play.TeamID)]
		zone := getZoneLabel(play.ZoneID)
		abbr := possessingTeam.Abbreviation
		penalty := getPenaltyByID(uint(play.PenaltyID))
		severity := getSeverityByID(play.Severity)
		isFight := "No"
		if play.IsFight {
			isFight = "Yes"
		}

		result := generateProResultsString(play, event, outcome, playerMap, possessingTeam)
		writer.Write([]string{
			periodStr,
			timeOnClock,
			timeConsumed,
			zone,
			event,
			outcome,
			penalty,
			severity,
			isFight,
			hts,
			ats,
			abbr,
			result,
		})
	}
	return err
}

func WriteBoxScoreFile(r engine.GameState, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filename, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Team", "1", "2", "3", "OT", "T"})
	hts := r.HomeTeamStats
	writer.Write([]string{r.HomeTeam, strconv.Itoa(int(hts.Period1Score)), strconv.Itoa(int(hts.Period2Score)), strconv.Itoa(int(hts.Period3Score)), strconv.Itoa(int(hts.OTScore)), strconv.Itoa(int(hts.Points))})
	ats := r.AwayTeamStats
	writer.Write([]string{r.AwayTeam, strconv.Itoa(int(ats.Period1Score)), strconv.Itoa(int(ats.Period2Score)), strconv.Itoa(int(ats.Period3Score)), strconv.Itoa(int(ats.OTScore)), strconv.Itoa(int(ats.Points))})
	writer.Write([]string{})
	writer.Write([]string{"Home Team"})
	writer.Write([]string{"Position", "Name", "G", "A", "P", "+/-", "PIM", "TOI", "PPG", "S", "BLK", "BCHK", "STCHK", "FO"})
	hpb := r.HomeStrategy
	for _, line := range hpb.Forwards {
		players := line.Players
		for _, p := range players {
			s := p.Stats
			writer.Write([]string{p.Position, p.FirstName + " " + p.LastName, strconv.Itoa(int(s.Goals)), strconv.Itoa(int(s.Assists)), strconv.Itoa(int(s.Points)), strconv.Itoa(int(s.PlusMinus)), FormatTimeToClock(s.PenaltyMinutes), FormatTimeToClock(uint16(s.TimeOnIce)), strconv.Itoa(int(s.PowerPlayGoals)), strconv.Itoa(int(s.Shots)), strconv.Itoa(int(s.ShotsBlocked)), strconv.Itoa(int(s.BodyChecks)), strconv.Itoa(int(s.StickChecks)), strconv.Itoa(int(s.FaceOffsWon))})
		}
	}
	for _, line := range hpb.Defenders {
		players := line.Players
		for _, p := range players {
			s := p.Stats
			writer.Write([]string{p.Position, p.FirstName + " " + p.LastName, strconv.Itoa(int(s.Goals)), strconv.Itoa(int(s.Assists)), strconv.Itoa(int(s.Points)), strconv.Itoa(int(s.PlusMinus)), FormatTimeToClock(s.PenaltyMinutes), FormatTimeToClock(uint16(s.TimeOnIce)), strconv.Itoa(int(s.PowerPlayGoals)), strconv.Itoa(int(s.Shots)), strconv.Itoa(int(s.ShotsBlocked)), strconv.Itoa(int(s.BodyChecks)), strconv.Itoa(int(s.StickChecks)), strconv.Itoa(int(s.FaceOffsWon))})
		}
	}
	writer.Write([]string{"Position", "Name", "SA", "SV", "GA", "SV%", "TOI"})
	for _, line := range hpb.Goalies {
		players := line.Players
		for _, p := range players {
			s := p.Stats
			writer.Write([]string{p.Position, p.FirstName + " " + p.LastName, strconv.Itoa(int(s.ShotsAgainst)), strconv.Itoa(int(s.Saves)), strconv.Itoa(int(s.GoalsAgainst)), strconv.Itoa(int(s.SavePercentage)), FormatTimeToClock(uint16(s.TimeOnIce))})
		}
	}
	// Iterate through play by play data to generate []string
	writer.Write([]string{})
	writer.Write([]string{"Away Team"})
	writer.Write([]string{"Position", "Name", "G", "A", "P", "+/-", "PIM", "TOI", "PPG", "S", "BLK", "BCHK", "STCHK", "FO"})
	apb := r.AwayStrategy
	for _, line := range apb.Forwards {
		players := line.Players
		for _, p := range players {
			s := p.Stats
			writer.Write([]string{p.Position, p.FirstName + " " + p.LastName, strconv.Itoa(int(s.Goals)), strconv.Itoa(int(s.Assists)), strconv.Itoa(int(s.Points)), strconv.Itoa(int(s.PlusMinus)), FormatTimeToClock(s.PenaltyMinutes), FormatTimeToClock(uint16(s.TimeOnIce)), strconv.Itoa(int(s.PowerPlayGoals)), strconv.Itoa(int(s.Shots)), strconv.Itoa(int(s.ShotsBlocked)), strconv.Itoa(int(s.BodyChecks)), strconv.Itoa(int(s.StickChecks)), strconv.Itoa(int(s.FaceOffsWon))})
		}
	}
	for _, line := range apb.Defenders {
		players := line.Players
		for _, p := range players {
			s := p.Stats
			writer.Write([]string{p.Position, p.FirstName + " " + p.LastName, strconv.Itoa(int(s.Goals)), strconv.Itoa(int(s.Assists)), strconv.Itoa(int(s.Points)), strconv.Itoa(int(s.PlusMinus)), FormatTimeToClock(s.PenaltyMinutes), FormatTimeToClock(uint16(s.TimeOnIce)), strconv.Itoa(int(s.PowerPlayGoals)), strconv.Itoa(int(s.Shots)), strconv.Itoa(int(s.ShotsBlocked)), strconv.Itoa(int(s.BodyChecks)), strconv.Itoa(int(s.StickChecks)), strconv.Itoa(int(s.FaceOffsWon))})
		}
	}
	writer.Write([]string{"Position", "Name", "SA", "SV", "GA", "SV%", "TOI"})
	for _, line := range apb.Goalies {
		players := line.Players
		for _, p := range players {
			s := p.Stats
			writer.Write([]string{p.Position, p.FirstName + " " + p.LastName, strconv.Itoa(int(s.ShotsAgainst)), strconv.Itoa(int(s.Saves)), strconv.Itoa(int(s.GoalsAgainst)), strconv.Itoa(int(s.SavePercentage)), FormatTimeToClock(uint16(s.TimeOnIce))})
		}
	}
	return err
}

func FormatTimeToClock(timeInSeconds uint16) string {
	minutes := timeInSeconds / 60
	seconds := timeInSeconds % 60
	formatted := fmt.Sprintf("%02d:%02d", minutes, seconds)
	return formatted
}

func generateCollegeResultsString(play structs.PbP, event, outcome string, playerMap map[uint]structs.CollegePlayer, possessingTeam structs.CollegeTeam) string {
	puckCarrier := playerMap[play.PuckCarrierID]
	puckCarrierLabel := getPlayerLabel(puckCarrier.BasePlayer)
	receivingPlayer := playerMap[play.PassedPlayerID]
	receivingPlayerLabel := getPlayerLabel(receivingPlayer.BasePlayer)
	assistingPlayer := playerMap[play.AssistingPlayerID]
	assistingPlayerLabel := getPlayerLabel(assistingPlayer.BasePlayer)
	defendingPlayer := playerMap[play.DefenderID]
	defendingPlayerLabel := getPlayerLabel(defendingPlayer.BasePlayer)
	goalie := playerMap[play.GoalieID]
	goalieLabel := getPlayerLabel(goalie.BasePlayer)
	statement := ""
	nextZoneLabel := getZoneLabel(play.NextZoneID)
	teamLabel := possessingTeam.TeamName
	// First Segment
	if event == Faceoff {
		if outcome == "Home Faceoff Win" {
			statement = puckCarrierLabel + " wins the faceoff! "
		} else if outcome == util.GoalieHold {
			statement = puckCarrierLabel + " holds onto the puck, and it's going to a faceoff."
		} else {
			statement = receivingPlayerLabel + " wins the faceoff! "
		}
		// Mention receiving player
		statement += assistingPlayerLabel + " receives the puck on the faceoff."
	} else if event == PhysDefenseCheck {
		if outcome == DefenseTakesPuck {
			statement = defendingPlayerLabel + " bodies " + puckCarrierLabel + " right into the boards and snatches the puck away!"
		} else if outcome == CarrierKeepsPuck {
			statement = defendingPlayerLabel + " attempts to body right into " + puckCarrierLabel + ", but " + puckCarrierLabel + " maneuvers effortlessly within the zone!"
		}
	} else if event == DexDefenseCheck {
		if outcome == DefenseTakesPuck {
			statement = defendingPlayerLabel + " with a bit of stick-play swipes the puck right from under " + puckCarrierLabel + "!"
		} else if outcome == CarrierKeepsPuck {
			statement = defendingPlayerLabel + " attempts to swipe the puck from " + puckCarrierLabel + ", but his stick is batted away!"
		}
	} else if event == PassCheck {
		if outcome == InterceptedPass {
			statement = defendingPlayerLabel + " intercepts the pass right from " + puckCarrierLabel + "!"
		} else if outcome == ReceivedPass {
			statement = puckCarrierLabel + " finds " + receivingPlayerLabel + " and makes the pass!"
		}
	} else if event == AgilityCheck {
		if outcome == DefenseStopAgility {
			statement = defendingPlayerLabel + " with a bit of stick-play swipes the puck right from under " + puckCarrierLabel + "!"
		} else if outcome == OffenseMovesUp {
			statement = puckCarrierLabel + " moves the puck up to the " + nextZoneLabel + "."
		}
	} else if event == WristshotCheck {
		statement = puckCarrierLabel + " attempts a long shot on goal..."
		if outcome == ShotBlocked {
			statement += " and the shot is blocked by " + defendingPlayerLabel + "!"
		} else if outcome == GoalieSave {
			statement += " and the shot is SAVED by " + goalieLabel + "!"
		} else if outcome == GoalieReboundOutcome {
			//
		} else if outcome == ShotOnGoal {
			statement += " and he scores! That's a point for " + teamLabel + "!"
		}
	} else if event == SlapshotCheck {
		statement = puckCarrierLabel + " attempts a slapshot on goal..."
		if outcome == ShotBlocked {
			statement += " and the shot is blocked by " + defendingPlayerLabel + "!"
		} else if outcome == GoalieSave {
			statement += " and the shot is SAVED by " + goalieLabel + "!"
		} else if outcome == GoalieReboundOutcome {
			//
		} else if outcome == ShotOnGoal {
			statement += " and he scores! That's a point for " + teamLabel + "!"
		} else if outcome == PenaltyCheck {
			penalty := getPenaltyByID(uint(play.PenaltyID))
			severity := getSeverityByID(play.Severity)
			penaltyMinutes := "two"
			if play.Severity > 1 {
				penaltyMinutes = "five"
			}
			statement += " and a penalty is called! " + defendingPlayerLabel + " has been called for a " + severity + " " + penalty + " on " + puckCarrierLabel + ". This will lead into a faceoff. Power play for " + penaltyMinutes + " minutes."
		}
	} else if event == PenaltyCheck {
		penalty := getPenaltyByID(uint(play.PenaltyID))
		severity := getSeverityByID(play.Severity)
		penaltyMinutes := "two"
		if play.Severity > 1 {
			penaltyMinutes = "five"
		}
		if play.IsFight {
			statement = "There's a fight on center ice! " + defendingPlayerLabel + " and " + goalieLabel + " are right at with the fisticuffs. Refs are breaking up the fight. Both players are out for " + penaltyMinutes + " minutes. Resetting play with a faceoff. "
		} else {
			statement = "Penalty called! " + defendingPlayerLabel + " has been called for " + severity + " " + penalty + " on " + puckCarrierLabel + ". Power play for " + penaltyMinutes + " minutes."
		}
	}

	return statement
}

func generateProResultsString(play structs.PbP, event, outcome string, playerMap map[uint]structs.ProfessionalPlayer, possessingTeam structs.ProfessionalTeam) string {
	puckCarrier := playerMap[play.PuckCarrierID]
	puckCarrierLabel := getPlayerLabel(puckCarrier.BasePlayer)
	receivingPlayer := playerMap[play.PassedPlayerID]
	receivingPlayerLabel := getPlayerLabel(receivingPlayer.BasePlayer)
	assistingPlayer := playerMap[play.AssistingPlayerID]
	assistingPlayerLabel := getPlayerLabel(assistingPlayer.BasePlayer)
	defendingPlayer := playerMap[play.DefenderID]
	defendingPlayerLabel := getPlayerLabel(defendingPlayer.BasePlayer)
	goalie := playerMap[play.GoalieID]
	goalieLabel := getPlayerLabel(goalie.BasePlayer)
	statement := ""
	nextZoneLabel := getZoneLabel(play.NextZoneID)
	teamLabel := possessingTeam.TeamName
	// First Segment
	if event == Faceoff {
		if outcome == "Home Faceoff Win" {
			statement = puckCarrierLabel + " wins the faceoff! "
		} else if outcome == util.GoalieHold {
			statement = puckCarrierLabel + " holds onto the puck, and it's going to a faceoff."
		} else {
			statement = receivingPlayerLabel + " wins the faceoff! "
		}
		// Mention receiving player
		statement += assistingPlayerLabel + " receives the puck on the faceoff."
	} else if event == PhysDefenseCheck {
		if outcome == DefenseTakesPuck {
			statement = defendingPlayerLabel + " bodies " + puckCarrierLabel + " right into the boards and snatches the puck away!"
		} else if outcome == CarrierKeepsPuck {
			statement = defendingPlayerLabel + " attempts to body right into " + puckCarrierLabel + ", but " + puckCarrierLabel + " maneuvers effortlessly within the zone!"
		}
	} else if event == DexDefenseCheck {
		if outcome == DefenseTakesPuck {
			statement = defendingPlayerLabel + " with a bit of stick-play swipes the puck right from under " + puckCarrierLabel + "!"
		} else if outcome == CarrierKeepsPuck {
			statement = defendingPlayerLabel + " attempts to swipe the puck from " + puckCarrierLabel + ", but his stick is batted away!"
		}
	} else if event == PassCheck {
		if outcome == InterceptedPass {
			statement = defendingPlayerLabel + " intercepts the pass right from " + puckCarrierLabel + "!"
		} else if outcome == ReceivedPass {
			statement = puckCarrierLabel + " finds " + receivingPlayerLabel + " and makes the pass!"
		}
	} else if event == AgilityCheck {
		if outcome == DefenseStopAgility {
			statement = defendingPlayerLabel + " with a bit of stick-play swipes the puck right from under " + puckCarrierLabel + "!"
		} else if outcome == OffenseMovesUp {
			statement = puckCarrierLabel + " moves the puck up to the " + nextZoneLabel + "."
		}
	} else if event == WristshotCheck {
		statement = puckCarrierLabel + " attempts a long shot on goal..."
		if outcome == ShotBlocked {
			statement += " and the shot is blocked by " + defendingPlayerLabel + "!"
		} else if outcome == GoalieSave {
			statement += " and the shot is SAVED by " + goalieLabel + "!"
		} else if outcome == GoalieReboundOutcome {
			//
		} else if outcome == ShotOnGoal {
			statement += " and he scores! That's a point for " + teamLabel + "!"
		}
	} else if event == SlapshotCheck {
		statement = puckCarrierLabel + " attempts a slapshot on goal..."
		if outcome == ShotBlocked {
			statement += " and the shot is blocked by " + defendingPlayerLabel + "!"
		} else if outcome == GoalieSave {
			statement += " and the shot is SAVED by " + goalieLabel + "!"
		} else if outcome == GoalieReboundOutcome {
			//
		} else if outcome == ShotOnGoal {
			statement += " and he scores! That's a point for " + teamLabel + "!"
		} else if outcome == PenaltyCheck {
			penalty := getPenaltyByID(uint(play.PenaltyID))
			severity := getSeverityByID(play.Severity)
			penaltyMinutes := "two"
			if play.Severity > 1 {
				penaltyMinutes = "five"
			}
			statement += " and a penalty is called! " + defendingPlayerLabel + " has been called for a " + severity + " " + penalty + " on " + puckCarrierLabel + ". This will lead into a faceoff. Power play for " + penaltyMinutes + " minutes."
		}
	} else if event == PenaltyCheck {
		penalty := getPenaltyByID(uint(play.PenaltyID))
		severity := getSeverityByID(play.Severity)
		penaltyMinutes := "two"
		if play.Severity > 1 {
			penaltyMinutes = "five"
		}
		if play.IsFight {
			statement = "There's a fight on center ice! " + defendingPlayerLabel + " and " + goalieLabel + " are right at with the fisticuffs. Refs are breaking up the fight. Both players are out for " + penaltyMinutes + " minutes. Resetting play with a faceoff. "
		} else {
			statement = "Penalty called! " + defendingPlayerLabel + " has been called for " + severity + " " + penalty + " on " + puckCarrierLabel + ". Power play for " + penaltyMinutes + " minutes."
		}
	}

	return statement
}

func getPlayerLabel(player structs.BasePlayer) string {
	if len(player.FirstName) == 0 {
		return ""
	}
	return player.Team + " " + player.Position + " " + player.FirstName + " " + player.LastName
}

func getZoneLabel(zoneID uint8) string {
	if zoneID == 0 {
		return ""
	}
	if zoneID == HomeGoalZoneID {
		return HomeGoal
	}
	if zoneID == HomeZoneID {
		return HomeZone
	}
	if zoneID == NeutralZoneID {
		return NeutralZone
	}
	if zoneID == AwayZoneID {
		return AwayZone
	}
	if zoneID == AwayGoalZoneID {
		return AwayGoal
	}
	return ""
}

func getPenaltyByID(penaltyType uint) string {
	var penaltyMap = map[uint]string{
		1:  "Aggressor Penalty",
		2:  "Attempt to Injure",
		3:  "Biting",
		4:  "Boarding",
		5:  "Boarding",
		6:  "Stabbing",
		7:  "Charging",
		8:  "Charging",
		9:  "Checking from Behind",
		10: "Checking from Behind",
		11: "Clipping",
		12: "Clipping",
		13: "Cross Checking",
		14: "Cross Checking",
		15: "Delay of Game",
		16: "Diving",
		17: "Elbowing",
		18: "Elbowing",
		19: "Eye-Gouging",
		20: "Fighting",
		21: "Goaltender Interference",
		22: "Headbutting",
		23: "High-sticking",
		24: "High-sticking",
		25: "Holding",
		26: "Hooking",
		27: "Hooking",
		28: "Kicking",
		29: "Kicking",
		30: "Kneeing",
		31: "Kneeing",
		32: "Roughing",
		33: "Roughing",
		34: "Slashing",
		35: "Slashing",
		36: "Slew footing",
		37: "Slew footing",
		38: "Throwing the stick",
		39: "Too many men on the ice",
		40: "Tripping",
		41: "Tripping",
		42: "Unsportsmanlike conduct",
	}
	return penaltyMap[penaltyType]
}

func getSeverityByID(sevId uint8) string {
	var severityMap = map[uint8]string{
		1: "Minor Penalty",
		2: "Major Penalty",
		3: "Game Misconduct",
		4: "Match Penalty",
	}
	return severityMap[sevId]
}

func WriteProPlayersExport(w http.ResponseWriter, players []structs.ProfessionalPlayer) {
	ts := GetTimestamp()
	w.Header().Set("Content-Disposition", "attachment;filename="+strconv.Itoa(int(ts.Season))+"_pro_players.csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	writer := csv.NewWriter(w)

	writer.Write([]string{"ID", "First Name", "Last Name", "Position", "Archetype",
		"Height", "Weight", "City", "Region", "Country", "Stars", "Age", "Overall",
		util.Agility, util.Faceoffs, util.LongShotAccuracy, util.LongShotPower, util.CloseShotAccuracy,
		util.CloseShotPower, util.Passing, util.PuckHandling, util.Strength, util.BodyChecking, util.StickChecking,
		util.ShotBlocking, util.Goalkeeping, util.GoalieVision, "Stamina", "Injury Rating", "Agility Pot.", "Faceoffs Pot.", "Long Shot Accuracy Pot.", "Long Shot Power Pot.",
		"Close Shot Accuracy Pot.", "Close Shot Power Pot.", "Passing Pot.", "Puck Handling Pot.",
		"Strength Pot.", "Body Checking Pot.", "Stick Checking Pot.", "Shot Blocking Pot.", "Goalkeeping Pot.", "Goalie Vision Pot."})

	for _, p := range players {
		idStr := strconv.Itoa(int(p.ID))

		playerRow := []string{
			idStr, p.FirstName, p.LastName, p.Position, p.Archetype, strconv.Itoa(int(p.Height)), strconv.Itoa(int(p.Weight)), p.City, p.State, p.Country,
			strconv.Itoa(int(p.Stars)), strconv.Itoa(int(p.Age)), strconv.Itoa(int(p.Overall)), strconv.Itoa(int(p.Agility)), strconv.Itoa(int(p.Faceoffs)), strconv.Itoa(int(p.LongShotAccuracy)),
			strconv.Itoa(int(p.LongShotPower)), strconv.Itoa(int(p.CloseShotAccuracy)), strconv.Itoa(int(p.CloseShotPower)), strconv.Itoa(int(p.Passing)), strconv.Itoa(int(p.PuckHandling)), strconv.Itoa(int(p.Strength)),
			strconv.Itoa(int(p.BodyChecking)), strconv.Itoa(int(p.StickChecking)), strconv.Itoa(int(p.ShotBlocking)), strconv.Itoa(int(p.Goalkeeping)), strconv.Itoa(int(p.GoalieVision)), util.GetPotentialGrade(int(p.Stamina)),
			util.GetPotentialGrade(int(p.InjuryRating)), util.GetPotentialGrade(int(p.AgilityPotential)), util.GetPotentialGrade(int(p.FaceoffsPotential)), util.GetPotentialGrade(int(p.LongShotAccuracyPotential)), util.GetPotentialGrade(int(p.LongShotPowerPotential)),
			util.GetPotentialGrade(int(p.CloseShotAccuracyPotential)), util.GetPotentialGrade(int(p.CloseShotPowerPotential)), util.GetPotentialGrade(int(p.PassingPotential)), util.GetPotentialGrade(int(p.PuckHandlingPotential)),
			util.GetPotentialGrade(int(p.StrengthPotential)), util.GetPotentialGrade(int(p.BodyCheckingPotential)), util.GetPotentialGrade(int(p.StickCheckingPotential)), util.GetPotentialGrade(int(p.ShotBlockingPotential)),
			util.GetPotentialGrade(int(p.GoalkeepingPotential)), util.GetPotentialGrade(int(p.GoalieVisionPotential)),
		}

		err := writer.Write(playerRow)
		if err != nil {
			log.Fatal("Cannot write player row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func WriteCollegePlayersExport(w http.ResponseWriter, players []structs.CollegePlayer) {
	ts := GetTimestamp()
	w.Header().Set("Content-Disposition", "attachment;filename="+strconv.Itoa(int(ts.Season))+"_chl_players.csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	writer := csv.NewWriter(w)

	writer.Write(getHeaderRow())

	for _, p := range players {
		idStr := strconv.Itoa(int(p.ID))

		playerRow := []string{
			idStr, p.Team, p.FirstName, p.LastName, p.Position, p.Archetype, strconv.Itoa(int(p.Height)), strconv.Itoa(int(p.Weight)), p.City, p.State, p.Country,
			strconv.Itoa(int(p.Stars)), strconv.Itoa(int(p.Age)), util.GetLetterGrade(int(p.Overall), p.Year), util.GetLetterGrade(int(p.Agility), p.Year), util.GetLetterGrade(int(p.Faceoffs), p.Year), util.GetLetterGrade(int(p.LongShotAccuracy), p.Year),
			util.GetLetterGrade(int(p.LongShotPower), p.Year), util.GetLetterGrade(int(p.CloseShotAccuracy), p.Year), util.GetLetterGrade(int(p.CloseShotPower), p.Year), util.GetLetterGrade(int(p.Passing), p.Year), util.GetLetterGrade(int(p.PuckHandling), p.Year), util.GetLetterGrade(int(p.Strength), p.Year),
			util.GetLetterGrade(int(p.BodyChecking), p.Year), util.GetLetterGrade(int(p.StickChecking), p.Year), util.GetLetterGrade(int(p.ShotBlocking), p.Year), util.GetLetterGrade(int(p.Goalkeeping), p.Year), util.GetLetterGrade(int(p.GoalieVision), p.Year), util.GetPotentialGrade(int(p.Stamina)),
			util.GetPotentialGrade(int(p.InjuryRating)), "?", "?", "?", "?",
			"?", "?", "?", "?",
			"?", "?", "?", "?",
			"?", "?",
		}

		err := writer.Write(playerRow)
		if err != nil {
			log.Fatal("Cannot write player row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}

func getHeaderRow() []string {
	return []string{"ID", "Team", "First Name", "Last Name", "Position", "Archetype",
		"Height", "Weight", "City", "Region", "Country", "Stars", "Age", "Overall",
		util.Agility, util.Faceoffs, util.LongShotAccuracy, util.LongShotPower, util.CloseShotAccuracy,
		util.CloseShotPower, util.Passing, util.PuckHandling, util.Strength, util.BodyChecking, util.StickChecking,
		util.ShotBlocking, util.Goalkeeping, util.GoalieVision, "Stamina", "Injury Rating", "Agility Pot.", "Faceoffs Pot.", "Long Shot Accuracy Pot.", "Long Shot Power Pot.",
		"Close Shot Accuracy Pot.", "Close Shot Power Pot.", "Passing Pot.", "Puck Handling Pot.",
		"Strength Pot.", "Body Checking Pot.", "Stick Checking Pot.", "Shot Blocking Pot.", "Goalkeeping Pot.", "Goalie Vision Pot."}
}
