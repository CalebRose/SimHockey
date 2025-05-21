package managers

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/engine"
	"github.com/CalebRose/SimHockey/repository"
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

func WriteProPlayersExport(w http.ResponseWriter, players []structs.ProfessionalPlayer, filename string) {
	ts := GetTimestamp()
	w.Header().Set("Content-Disposition", "attachment;filename="+strconv.Itoa(int(ts.Season))+filename)
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
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

func WriteCollegePlayersExport(w http.ResponseWriter, players []structs.CollegePlayer, filename string) {
	ts := GetTimestamp()
	w.Header().Set("Content-Disposition", "attachment;filename="+strconv.Itoa(int(ts.Season))+filename)
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
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

func WriteCollegeRecruitsExport(w http.ResponseWriter, players []structs.Recruit, filename string) {
	ts := GetTimestamp()
	w.Header().Set("Content-Disposition", "attachment;filename="+strconv.Itoa(int(ts.Season))+filename)
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	w.Header().Set("Transfer-Encoding", "chunked")
	writer := csv.NewWriter(w)

	writer.Write(getHeaderRow())

	for _, p := range players {
		idStr := strconv.Itoa(int(p.ID))

		year := 1
		playerRow := []string{
			idStr, p.Team, p.FirstName, p.LastName, p.Position, p.Archetype, strconv.Itoa(int(p.Height)), strconv.Itoa(int(p.Weight)), p.City, p.State, p.Country,
			strconv.Itoa(int(p.Stars)), strconv.Itoa(int(p.Age)), util.GetLetterGrade(int(p.Overall), year), util.GetLetterGrade(int(p.Agility), year), util.GetLetterGrade(int(p.Faceoffs), year), util.GetLetterGrade(int(p.LongShotAccuracy), year),
			util.GetLetterGrade(int(p.LongShotPower), year), util.GetLetterGrade(int(p.CloseShotAccuracy), year), util.GetLetterGrade(int(p.CloseShotPower), year), util.GetLetterGrade(int(p.Passing), year), util.GetLetterGrade(int(p.PuckHandling), year), util.GetLetterGrade(int(p.Strength), year),
			util.GetLetterGrade(int(p.BodyChecking), year), util.GetLetterGrade(int(p.StickChecking), year), util.GetLetterGrade(int(p.ShotBlocking), year), util.GetLetterGrade(int(p.Goalkeeping), year), util.GetLetterGrade(int(p.GoalieVision), year), util.GetPotentialGrade(int(p.Stamina)),
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

func ExportCollegeStats(seasonID, weekID, viewType, gameType string, w http.ResponseWriter) {
	stats := SearchCollegeStats(seasonID, weekID, viewType, gameType)
	seasonIDNum := util.ConvertStringToInt(seasonID)
	seasonIDNum += 2024
	seasonStr := strconv.Itoa(seasonIDNum)
	weekStr := ""
	if viewType != "SEASON" && (weekID != "" && weekID != "0") {
		weekNum := util.ConvertStringToInt(weekID)
		weekNum = weekNum - 2500
		weekStr = "_WEEK_" + strconv.Itoa(weekNum) + "_"
	}
	baseName := fmt.Sprintf("chl_stats_%s_%s", seasonStr, viewType)
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", baseName))
	w.Header().Set("Transfer-Encoding", "chunked")
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()
	// Initialize writer
	fileName := "chl_player_stats_" + seasonStr + weekStr + ".csv"
	teamFileName := "chl_team_stats_" + seasonStr + weekStr + ".csv"

	collegePlayerMap := GetCollegePlayersMap()
	chlTeamMap := GetCollegeTeamMap()

	writeCSVIntoZip(zipWriter, fileName, func(csvW *csv.Writer) error {
		header := util.GetCHLPlayerHeaderRows()
		if err := csvW.Write(header); err != nil {
			return err
		}

		if viewType == "WEEK" {
			for _, stat := range stats.CHLPlayerGameStats {
				p := collegePlayerMap[stat.PlayerID]
				if p.ID == 0 {
					continue
				}
				Year, RedshirtStatus := util.GetYearAndRedshirtStatus(p.Year, p.IsRedshirt)

				team := chlTeamMap[stat.TeamID]

				answer := "No."
				diceRoll := util.GenerateIntFromRange(1, 1000)
				if diceRoll == 1000 {
					answer = "Yes."
				}

				injury := "No"
				if stat.IsInjured {
					injury = "Yes"
				}

				timeOnIce := FormatTimeToClock(uint16(stat.TimeOnIce))

				pr := []string{strconv.Itoa(int(p.ID)), p.FirstName, p.LastName, p.Position,
					p.Archetype, Year, RedshirtStatus, p.Team, team.Conference, strconv.Itoa(int(p.Age)), strconv.Itoa(int(p.Stars)),
					strconv.Itoa(int(stat.Goals)), strconv.Itoa(int(stat.Assists)), strconv.Itoa(int(stat.Points)), strconv.Itoa(int(stat.PlusMinus)),
					strconv.Itoa(int(stat.PenaltyMinutes)), strconv.Itoa(int(stat.EvenStrengthGoals)), strconv.Itoa(int(stat.EvenStrengthPoints)), strconv.Itoa(int(stat.PowerPlayGoals)),
					strconv.Itoa(int(stat.PowerPlayPoints)), strconv.Itoa(int(stat.ShorthandedGoals)), strconv.Itoa(int(stat.ShorthandedPoints)), strconv.Itoa(int(stat.OvertimeGoals)),
					strconv.Itoa(int(stat.GameWinningGoals)), strconv.Itoa(int(stat.Shots)), strconv.Itoa(int(stat.ShootingPercentage)), timeOnIce,
					strconv.Itoa(int(stat.FaceOffWinPercentage)), strconv.Itoa(int(stat.FaceOffsWon)), strconv.Itoa(int(stat.FaceOffs)), strconv.Itoa(int(stat.GoalieWins)),
					strconv.Itoa(int(stat.GoalieLosses)), strconv.Itoa(int(stat.GoalieTies)), strconv.Itoa(int(stat.OvertimeLosses)), strconv.Itoa(int(stat.ShotsAgainst)),
					strconv.Itoa(int(stat.Saves)), strconv.Itoa(int(stat.GoalsAgainst)), strconv.Itoa(int(stat.SavePercentage)), strconv.Itoa(int(stat.Shutouts)),
					strconv.Itoa(int(stat.ShotsBlocked)), strconv.Itoa(int(stat.BodyChecks)), strconv.Itoa(int(stat.StickChecks)), injury,
					stat.InjuryName, stat.InjuryType, answer,
				}
				err := csvW.Write(pr)
				if err != nil {
					log.Fatal("Cannot write player row to CSV", err)
				}

				csvW.Flush()
				err = csvW.Error()
				if err != nil {
					log.Fatal("Error while writing to file ::", err)
				}
			}
		} else {
			for _, stat := range stats.CHLPlayerSeasonStats {
				p := collegePlayerMap[stat.PlayerID]
				if p.ID == 0 {
					continue
				}
				Year, RedshirtStatus := util.GetYearAndRedshirtStatus(p.Year, p.IsRedshirt)

				team := chlTeamMap[stat.TeamID]

				answer := "No."
				diceRoll := util.GenerateIntFromRange(1, 1000)
				if diceRoll == 1000 {
					answer = "Yes."
				}

				injury := "No"
				if stat.IsInjured {
					injury = "Yes"
				}

				timeOnIce := FormatTimeToClock(uint16(stat.TimeOnIce))

				pr := []string{strconv.Itoa(int(p.ID)), p.FirstName, p.LastName, p.Position,
					p.Archetype, Year, RedshirtStatus, p.Team, team.Conference, strconv.Itoa(int(p.Age)), strconv.Itoa(int(p.Stars)),
					strconv.Itoa(int(stat.Goals)), strconv.Itoa(int(stat.Assists)), strconv.Itoa(int(stat.Points)), strconv.Itoa(int(stat.PlusMinus)),
					strconv.Itoa(int(stat.PenaltyMinutes)), strconv.Itoa(int(stat.EvenStrengthGoals)), strconv.Itoa(int(stat.EvenStrengthPoints)), strconv.Itoa(int(stat.PowerPlayGoals)),
					strconv.Itoa(int(stat.PowerPlayPoints)), strconv.Itoa(int(stat.ShorthandedGoals)), strconv.Itoa(int(stat.ShorthandedPoints)), strconv.Itoa(int(stat.OvertimeGoals)),
					strconv.Itoa(int(stat.GameWinningGoals)), strconv.Itoa(int(stat.Shots)), strconv.Itoa(int(stat.ShootingPercentage)), timeOnIce,
					strconv.Itoa(int(stat.FaceOffWinPercentage)), strconv.Itoa(int(stat.FaceOffsWon)), strconv.Itoa(int(stat.FaceOffs)), strconv.Itoa(int(stat.GoalieWins)),
					strconv.Itoa(int(stat.GoalieLosses)), strconv.Itoa(int(stat.GoalieTies)), strconv.Itoa(int(stat.OvertimeLosses)), strconv.Itoa(int(stat.ShotsAgainst)),
					strconv.Itoa(int(stat.Saves)), strconv.Itoa(int(stat.GoalsAgainst)), strconv.Itoa(int(stat.SavePercentage)), strconv.Itoa(int(stat.Shutouts)),
					strconv.Itoa(int(stat.ShotsBlocked)), strconv.Itoa(int(stat.BodyChecks)), strconv.Itoa(int(stat.StickChecks)), injury,
					stat.InjuryName, stat.InjuryType, answer,
				}
				err := csvW.Write(pr)
				if err != nil {
					log.Fatal("Cannot write player row to CSV", err)
				}

				csvW.Flush()
				err = csvW.Error()
				if err != nil {
					log.Fatal("Error while writing to file ::", err)
				}
			}
		}

		return csvW.Error()
	})

	writeCSVIntoZip(zipWriter, teamFileName, func(csvW *csv.Writer) error {
		header := util.GetCHLTeamsHeaderRows()
		if err := csvW.Write(header); err != nil {
			return err
		}

		if viewType == "WEEK" {
			for _, stat := range stats.CHLTeamGameStats {
				team := chlTeamMap[stat.TeamID]

				answer := "No."
				diceRoll := util.GenerateIntFromRange(1, 1000)
				if diceRoll == 1000 {
					answer = "Yes."
				}

				pr := []string{strconv.Itoa(int(team.ID)), team.TeamName, team.Conference,
					strconv.Itoa(int(stat.GoalsFor)), strconv.Itoa(int(stat.GoalsAgainst)), strconv.Itoa(int(stat.Assists)), strconv.Itoa(int(stat.Points)),
					strconv.Itoa(int(stat.Period1Score)), strconv.Itoa(int(stat.Period2Score)), strconv.Itoa(int(stat.Period3Score)), strconv.Itoa(int(stat.OTScore)),
					strconv.Itoa(int(stat.PlusMinus)),
					strconv.Itoa(int(stat.PenaltyMinutes)), strconv.Itoa(int(stat.EvenStrengthGoals)), strconv.Itoa(int(stat.EvenStrengthPoints)), strconv.Itoa(int(stat.PowerPlayGoals)),
					strconv.Itoa(int(stat.PowerPlayPoints)), strconv.Itoa(int(stat.ShorthandedGoals)), strconv.Itoa(int(stat.ShorthandedPoints)), strconv.Itoa(int(stat.OvertimeGoals)),
					strconv.Itoa(int(stat.GameWinningGoals)), strconv.Itoa(int(stat.Shots)), strconv.Itoa(int(stat.ShootingPercentage)),
					strconv.Itoa(int(stat.FaceOffWinPercentage)), strconv.Itoa(int(stat.FaceOffsWon)), strconv.Itoa(int(stat.FaceOffs)), strconv.Itoa(int(stat.ShotsAgainst)),
					strconv.Itoa(int(stat.Saves)), strconv.Itoa(int(stat.GoalsAgainst)), strconv.Itoa(int(stat.SavePercentage)), strconv.Itoa(int(stat.Shutouts)),
					answer,
				}
				err := csvW.Write(pr)
				if err != nil {
					log.Fatal("Cannot write player row to CSV", err)
				}

				csvW.Flush()
				err = csvW.Error()
				if err != nil {
					log.Fatal("Error while writing to file ::", err)
				}
			}
		} else {
			for _, stat := range stats.CHLTeamSeasonStats {
				team := chlTeamMap[stat.TeamID]

				answer := "No."
				diceRoll := util.GenerateIntFromRange(1, 1000)
				if diceRoll == 1000 {
					answer = "Yes."
				}

				pr := []string{strconv.Itoa(int(team.ID)), team.TeamName, team.Conference,
					strconv.Itoa(int(stat.GoalsFor)), strconv.Itoa(int(stat.GoalsAgainst)), strconv.Itoa(int(stat.Assists)), strconv.Itoa(int(stat.Points)),
					strconv.Itoa(int(stat.Period1Score)), strconv.Itoa(int(stat.Period2Score)), strconv.Itoa(int(stat.Period3Score)), strconv.Itoa(int(stat.OTScore)),
					strconv.Itoa(int(stat.PlusMinus)),
					strconv.Itoa(int(stat.PenaltyMinutes)), strconv.Itoa(int(stat.EvenStrengthGoals)), strconv.Itoa(int(stat.EvenStrengthPoints)), strconv.Itoa(int(stat.PowerPlayGoals)),
					strconv.Itoa(int(stat.PowerPlayPoints)), strconv.Itoa(int(stat.ShorthandedGoals)), strconv.Itoa(int(stat.ShorthandedPoints)), strconv.Itoa(int(stat.OvertimeGoals)),
					strconv.Itoa(int(stat.GameWinningGoals)), strconv.Itoa(int(stat.Shots)), strconv.Itoa(int(stat.ShootingPercentage)),
					strconv.Itoa(int(stat.FaceOffWinPercentage)), strconv.Itoa(int(stat.FaceOffsWon)), strconv.Itoa(int(stat.FaceOffs)), strconv.Itoa(int(stat.ShotsAgainst)),
					strconv.Itoa(int(stat.Saves)), strconv.Itoa(int(stat.GoalsAgainst)), strconv.Itoa(int(stat.SavePercentage)), strconv.Itoa(int(stat.Shutouts)),
					answer,
				}
				err := csvW.Write(pr)
				if err != nil {
					log.Fatal("Cannot write player row to CSV", err)
				}

				csvW.Flush()
				err = csvW.Error()
				if err != nil {
					log.Fatal("Error while writing to file ::", err)
				}
			}
		}

		return csvW.Error()
	})
}

func ExportProStats(seasonID, weekID, viewType, gameType string, w http.ResponseWriter) {
	stats := SearchProStats(seasonID, weekID, viewType, gameType)
	seasonIDNum := util.ConvertStringToInt(seasonID)
	seasonIDNum += 2024
	seasonStr := strconv.Itoa(seasonIDNum)
	weekStr := ""
	if viewType != "SEASON" && (weekID != "" && weekID != "0") {
		weekNum := util.ConvertStringToInt(weekID)
		weekNum = weekNum - 2500
		weekStr = "WEEK_" + strconv.Itoa(weekNum) + "_"
	}
	baseName := fmt.Sprintf("phl_stats_%s_%s", seasonStr, viewType)
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", baseName))
	w.Header().Set("Transfer-Encoding", "chunked")
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()
	// Initialize writer
	fileName := "phl_player_stats_" + seasonStr + "_" + weekStr
	teamFileName := "phl_team_stats_" + seasonStr + "_" + weekStr

	proPlayerMap := GetProPlayersMap()
	phlTeamMap := GetProTeamMap()

	writeCSVIntoZip(zipWriter, fileName, func(csvW *csv.Writer) error {
		header := util.GetPHLPlayerHeaderRows()
		if err := csvW.Write(header); err != nil {
			return err
		}

		if viewType == "WEEK" {
			for _, stat := range stats.PHLPlayerGameStats {
				p := proPlayerMap[stat.PlayerID]
				if p.ID == 0 {
					continue
				}
				team := phlTeamMap[stat.TeamID]

				answer := "No."
				diceRoll := util.GenerateIntFromRange(1, 1000)
				if diceRoll == 1000 {
					answer = "Yes."
				}

				injury := "No"
				if stat.IsInjured {
					injury = "Yes"
				}

				timeOnIce := FormatTimeToClock(uint16(stat.TimeOnIce))

				pr := []string{strconv.Itoa(int(p.ID)), p.FirstName, p.LastName, p.Position,
					p.Archetype, strconv.Itoa(p.Year), p.Team, team.Division, strconv.Itoa(int(p.Age)), strconv.Itoa(int(p.Stars)),
					strconv.Itoa(int(stat.Goals)), strconv.Itoa(int(stat.Assists)), strconv.Itoa(int(stat.Points)), strconv.Itoa(int(stat.PlusMinus)),
					strconv.Itoa(int(stat.PenaltyMinutes)), strconv.Itoa(int(stat.EvenStrengthGoals)), strconv.Itoa(int(stat.EvenStrengthPoints)), strconv.Itoa(int(stat.PowerPlayGoals)),
					strconv.Itoa(int(stat.PowerPlayPoints)), strconv.Itoa(int(stat.ShorthandedGoals)), strconv.Itoa(int(stat.ShorthandedPoints)), strconv.Itoa(int(stat.OvertimeGoals)),
					strconv.Itoa(int(stat.GameWinningGoals)), strconv.Itoa(int(stat.Shots)), strconv.Itoa(int(stat.ShootingPercentage)), timeOnIce,
					strconv.Itoa(int(stat.FaceOffWinPercentage)), strconv.Itoa(int(stat.FaceOffsWon)), strconv.Itoa(int(stat.FaceOffs)), strconv.Itoa(int(stat.GoalieWins)),
					strconv.Itoa(int(stat.GoalieLosses)), strconv.Itoa(int(stat.GoalieTies)), strconv.Itoa(int(stat.OvertimeLosses)), strconv.Itoa(int(stat.ShotsAgainst)),
					strconv.Itoa(int(stat.Saves)), strconv.Itoa(int(stat.GoalsAgainst)), strconv.Itoa(int(stat.SavePercentage)), strconv.Itoa(int(stat.Shutouts)),
					strconv.Itoa(int(stat.ShotsBlocked)), strconv.Itoa(int(stat.BodyChecks)), strconv.Itoa(int(stat.StickChecks)), injury,
					stat.InjuryName, stat.InjuryType, answer,
				}
				err := csvW.Write(pr)
				if err != nil {
					log.Fatal("Cannot write player row to CSV", err)
				}

				csvW.Flush()
				err = csvW.Error()
				if err != nil {
					log.Fatal("Error while writing to file ::", err)
				}
			}
		} else {
			for _, stat := range stats.PHLPlayerSeasonStats {
				p := proPlayerMap[stat.PlayerID]
				if p.ID == 0 {
					continue
				}
				team := phlTeamMap[stat.TeamID]

				answer := "No."
				diceRoll := util.GenerateIntFromRange(1, 1000)
				if diceRoll == 1000 {
					answer = "Yes."
				}

				injury := "No"
				if stat.IsInjured {
					injury = "Yes"
				}

				timeOnIce := FormatTimeToClock(uint16(stat.TimeOnIce))

				pr := []string{strconv.Itoa(int(p.ID)), p.FirstName, p.LastName, p.Position,
					p.Archetype, strconv.Itoa(p.Year), p.Team, team.Conference, strconv.Itoa(int(p.Age)), strconv.Itoa(int(p.Stars)),
					strconv.Itoa(int(stat.Goals)), strconv.Itoa(int(stat.Assists)), strconv.Itoa(int(stat.Points)), strconv.Itoa(int(stat.PlusMinus)),
					strconv.Itoa(int(stat.PenaltyMinutes)), strconv.Itoa(int(stat.EvenStrengthGoals)), strconv.Itoa(int(stat.EvenStrengthPoints)), strconv.Itoa(int(stat.PowerPlayGoals)),
					strconv.Itoa(int(stat.PowerPlayPoints)), strconv.Itoa(int(stat.ShorthandedGoals)), strconv.Itoa(int(stat.ShorthandedPoints)), strconv.Itoa(int(stat.OvertimeGoals)),
					strconv.Itoa(int(stat.GameWinningGoals)), strconv.Itoa(int(stat.Shots)), strconv.Itoa(int(stat.ShootingPercentage)), timeOnIce,
					strconv.Itoa(int(stat.FaceOffWinPercentage)), strconv.Itoa(int(stat.FaceOffsWon)), strconv.Itoa(int(stat.FaceOffs)), strconv.Itoa(int(stat.GoalieWins)),
					strconv.Itoa(int(stat.GoalieLosses)), strconv.Itoa(int(stat.GoalieTies)), strconv.Itoa(int(stat.OvertimeLosses)), strconv.Itoa(int(stat.ShotsAgainst)),
					strconv.Itoa(int(stat.Saves)), strconv.Itoa(int(stat.GoalsAgainst)), strconv.Itoa(int(stat.SavePercentage)), strconv.Itoa(int(stat.Shutouts)),
					strconv.Itoa(int(stat.ShotsBlocked)), strconv.Itoa(int(stat.BodyChecks)), strconv.Itoa(int(stat.StickChecks)), injury,
					stat.InjuryName, stat.InjuryType, answer,
				}
				err := csvW.Write(pr)
				if err != nil {
					log.Fatal("Cannot write player row to CSV", err)
				}

				csvW.Flush()
				err = csvW.Error()
				if err != nil {
					log.Fatal("Error while writing to file ::", err)
				}
			}
		}

		return csvW.Error()
	})

	writeCSVIntoZip(zipWriter, teamFileName, func(csvW *csv.Writer) error {
		header := util.GetPHLTeamsHeaderRows()
		if err := csvW.Write(header); err != nil {
			return err
		}

		if viewType == "WEEK" {
			for _, stat := range stats.PHLTeamGameStats {
				team := phlTeamMap[stat.TeamID]

				answer := "No."
				diceRoll := util.GenerateIntFromRange(1, 1000)
				if diceRoll == 1000 {
					answer = "Yes."
				}

				pr := []string{strconv.Itoa(int(team.ID)), team.TeamName, team.Conference,
					strconv.Itoa(int(stat.GoalsFor)), strconv.Itoa(int(stat.GoalsAgainst)), strconv.Itoa(int(stat.Assists)), strconv.Itoa(int(stat.Points)),
					strconv.Itoa(int(stat.Period1Score)), strconv.Itoa(int(stat.Period2Score)), strconv.Itoa(int(stat.Period3Score)), strconv.Itoa(int(stat.OTScore)),
					strconv.Itoa(int(stat.PlusMinus)),
					strconv.Itoa(int(stat.PenaltyMinutes)), strconv.Itoa(int(stat.EvenStrengthGoals)), strconv.Itoa(int(stat.EvenStrengthPoints)), strconv.Itoa(int(stat.PowerPlayGoals)),
					strconv.Itoa(int(stat.PowerPlayPoints)), strconv.Itoa(int(stat.ShorthandedGoals)), strconv.Itoa(int(stat.ShorthandedPoints)), strconv.Itoa(int(stat.OvertimeGoals)),
					strconv.Itoa(int(stat.GameWinningGoals)), strconv.Itoa(int(stat.Shots)), strconv.Itoa(int(stat.ShootingPercentage)),
					strconv.Itoa(int(stat.FaceOffWinPercentage)), strconv.Itoa(int(stat.FaceOffsWon)), strconv.Itoa(int(stat.FaceOffs)), strconv.Itoa(int(stat.ShotsAgainst)),
					strconv.Itoa(int(stat.Saves)), strconv.Itoa(int(stat.GoalsAgainst)), strconv.Itoa(int(stat.SavePercentage)), strconv.Itoa(int(stat.Shutouts)),
					answer,
				}
				err := csvW.Write(pr)
				if err != nil {
					log.Fatal("Cannot write player row to CSV", err)
				}

				csvW.Flush()
				err = csvW.Error()
				if err != nil {
					log.Fatal("Error while writing to file ::", err)
				}
			}
		} else {
			for _, stat := range stats.PHLTeamSeasonStats {
				team := phlTeamMap[stat.TeamID]

				answer := "No."
				diceRoll := util.GenerateIntFromRange(1, 1000)
				if diceRoll == 1000 {
					answer = "Yes."
				}

				pr := []string{strconv.Itoa(int(team.ID)), team.TeamName, team.Conference,
					strconv.Itoa(int(stat.GoalsFor)), strconv.Itoa(int(stat.GoalsAgainst)), strconv.Itoa(int(stat.Assists)), strconv.Itoa(int(stat.Points)),
					strconv.Itoa(int(stat.Period1Score)), strconv.Itoa(int(stat.Period2Score)), strconv.Itoa(int(stat.Period3Score)), strconv.Itoa(int(stat.OTScore)),
					strconv.Itoa(int(stat.PlusMinus)),
					strconv.Itoa(int(stat.PenaltyMinutes)), strconv.Itoa(int(stat.EvenStrengthGoals)), strconv.Itoa(int(stat.EvenStrengthPoints)), strconv.Itoa(int(stat.PowerPlayGoals)),
					strconv.Itoa(int(stat.PowerPlayPoints)), strconv.Itoa(int(stat.ShorthandedGoals)), strconv.Itoa(int(stat.ShorthandedPoints)), strconv.Itoa(int(stat.OvertimeGoals)),
					strconv.Itoa(int(stat.GameWinningGoals)), strconv.Itoa(int(stat.Shots)), strconv.Itoa(int(stat.ShootingPercentage)),
					strconv.Itoa(int(stat.FaceOffWinPercentage)), strconv.Itoa(int(stat.FaceOffsWon)), strconv.Itoa(int(stat.FaceOffs)), strconv.Itoa(int(stat.ShotsAgainst)),
					strconv.Itoa(int(stat.Saves)), strconv.Itoa(int(stat.GoalsAgainst)), strconv.Itoa(int(stat.SavePercentage)), strconv.Itoa(int(stat.Shutouts)),
					answer,
				}
				err := csvW.Write(pr)
				if err != nil {
					log.Fatal("Cannot write player row to CSV", err)
				}

				csvW.Flush()
				err = csvW.Error()
				if err != nil {
					log.Fatal("Error while writing to file ::", err)
				}
			}
		}

		return csvW.Error()
	})
}

func writeCSVIntoZip(z *zip.Writer, filename string, writeRows func(*csv.Writer) error) {
	f, err := z.Create(filename)
	if err != nil {
		// handle error (log, panic, or return to client)
		panic("unable to create zip entry: " + err.Error())
	}
	csvW := csv.NewWriter(f)
	if err := writeRows(csvW); err != nil {
		panic("error writing CSV data: " + err.Error())
	}
}

func ExportHCKGameResults(w http.ResponseWriter, seasonID, weekID, timeslot string) {
	ts := GetTimestamp()
	baseName := fmt.Sprintf("hck_game_results_%s_%s_%s", seasonID, weekID, timeslot)
	chlFilename := fmt.Sprintf("chl_game_results_%s_%s_%s.csv", seasonID, weekID, timeslot)
	phlFilename := fmt.Sprintf("phl_game_results_%s_%s_%s.csv", seasonID, weekID, timeslot)

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", baseName))
	w.Header().Set("Transfer-Encoding", "chunked")
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()
	isExactWeek := weekID == strconv.Itoa(int(ts.WeekID)) && seasonID == strconv.Itoa(int(ts.SeasonID))
	gameNotRan := (timeslot == "A" && !ts.GamesARan) ||
		(timeslot == "B" && !ts.GamesBRan) ||
		(timeslot == "C" && !ts.GamesCRan) ||
		(timeslot == "D" && !ts.GamesDRan)
	// Get All needed data
	matchChn := make(chan []structs.CollegeGame)
	phlMatchChn := make(chan []structs.ProfessionalGame)

	go func() {
		matches := repository.FindCollegeGames(repository.GamesClauses{WeekID: weekID, Timeslot: timeslot, IsPreseason: ts.IsPreseason})
		matchChn <- matches
	}()

	go func() {
		proGames := repository.FindProfessionalGames(repository.GamesClauses{WeekID: weekID, Timeslot: timeslot, IsPreseason: ts.IsPreseason})
		phlMatchChn <- proGames
	}()

	collegePlayers := GetAllCollegePlayers()
	historicPlayers := GetAllHistoricCollegePlayers()
	chlTeamMap := GetCollegeTeamMap()
	proTeamMap := GetProTeamMap()

	for _, hp := range historicPlayers {
		player := structs.CollegePlayer{Model: hp.Model, BasePlayer: hp.BasePlayer}
		collegePlayers = append(collegePlayers, player)
	}

	collegePlayerMap := MakeCollegePlayerMap(collegePlayers)

	proPlayers := GetAllProPlayers()
	retiredPlayers := GetAllRetiredPlayers()
	for _, r := range retiredPlayers {
		player := structs.ProfessionalPlayer{Model: r.Model, BasePlayer: r.BasePlayer}
		proPlayers = append(proPlayers, player)
	}

	proPlayerMap := MakeProfessionalPlayerMap(proPlayers)

	collegeGames := <-matchChn
	close(matchChn)
	proGames := <-phlMatchChn
	close(phlMatchChn)

	HeaderRow := []string{
		"League", "Week", "Home Team", "Home Score",
		"Away Team", "Away Score", "Is Overtime", "IsShootout", "HT Shootout Score", "AT Shootout Score", "Home Coach", "Home Rank", "Away Coach", "Away Rank", "Game Title",
		"Neutral Site", "Conference", "Game Day", "Arena", "Attendance", "City", "State", "Country", "Third Star", "Second Star", "First Star",
	}

	writeCSVIntoZip(zipWriter, chlFilename, func(csvW *csv.Writer) error {
		err := csvW.Write(HeaderRow)
		if err != nil {
			log.Fatal("Cannot write header row", err)
		}
		for _, m := range collegeGames {
			if isExactWeek && gameNotRan {
				m.HideScore()
			}
			neutralStr := "N"
			if m.IsNeutralSite {
				neutralStr = "Y"
			}
			confStr := "N"
			if m.IsConference {
				confStr = "Y"
			}
			otStr := "N"
			if m.IsOvertime {
				otStr = "Y"
			}
			soStr := "N"
			if m.IsShootout {
				soStr = "Y"
			}

			thirdStarStr := ""
			if m.StarThree > 0 {
				player := collegePlayerMap[m.StarThree]
				thirdStarStr = strconv.Itoa(int(player.ID)) + " " + player.Position + " " + player.FirstName + " " + player.LastName
			}

			secondStarStr := ""
			if m.StarTwo > 0 {
				player := collegePlayerMap[m.StarTwo]
				secondStarStr = strconv.Itoa(int(player.ID)) + " " + player.Position + " " + player.FirstName + " " + player.LastName
			}

			firstStarStr := ""
			if m.StarOne > 0 {
				player := collegePlayerMap[m.StarOne]
				firstStarStr = strconv.Itoa(int(player.ID)) + " " + player.Position + " " + player.FirstName + " " + player.LastName
			}

			homeTeam := chlTeamMap[m.HomeTeamID]
			awayTeam := chlTeamMap[m.AwayTeamID]

			row := []string{
				"CHL", strconv.Itoa(int(m.Week)), homeTeam.Abbreviation, strconv.Itoa(int(m.HomeTeamScore)),
				awayTeam.Abbreviation, strconv.Itoa(int(m.AwayTeamScore)), otStr, soStr, strconv.Itoa(int(m.HomeTeamShootoutScore)), strconv.Itoa(int(m.AwayTeamShootoutScore)), m.HomeTeamCoach,
				strconv.Itoa(int(m.HomeTeamRank)), m.AwayTeamCoach, strconv.Itoa(int(m.AwayTeamRank)), m.GameTitle,
				neutralStr, confStr, m.GameDay, m.Arena, strconv.Itoa(int(m.AttendanceCount)), m.City, m.State, m.Country, thirdStarStr, secondStarStr, firstStarStr,
			}
			err = csvW.Write(row)
			if err != nil {
				log.Fatal("Cannot write row to CSV", err)
			}

			csvW.Flush()
			err = csvW.Error()
			if err != nil {
				log.Fatal("Error while writing to file ::", err)
			}
		}
		return csvW.Error()
	})
	writeCSVIntoZip(zipWriter, phlFilename, func(csvW *csv.Writer) error {
		err := csvW.Write(HeaderRow)
		if err != nil {
			log.Fatal("Cannot write header row", err)
		}
		for _, m := range proGames {
			neutralStr := "N"
			if m.IsNeutralSite {
				neutralStr = "Y"
			}
			confStr := "N"
			if m.IsConference {
				confStr = "Y"
			}
			otStr := "N"
			if m.IsOvertime {
				otStr = "Y"
			}
			soStr := "N"
			if m.IsShootout {
				soStr = "Y"
			}

			thirdStarStr := ""
			if m.StarThree > 0 {
				player := proPlayerMap[m.StarThree]
				thirdStarStr = strconv.Itoa(int(player.ID)) + " " + player.Position + " " + player.FirstName + " " + player.LastName
			}

			secondStarStr := ""
			if m.StarTwo > 0 {
				player := proPlayerMap[m.StarTwo]
				secondStarStr = strconv.Itoa(int(player.ID)) + " " + player.Position + " " + player.FirstName + " " + player.LastName
			}

			firstStarStr := ""
			if m.StarOne > 0 {
				player := proPlayerMap[m.StarOne]
				firstStarStr = strconv.Itoa(int(player.ID)) + " " + player.Position + " " + player.FirstName + " " + player.LastName
			}
			homeTeam := proTeamMap[m.HomeTeamID]
			awayTeam := proTeamMap[m.AwayTeamID]

			row := []string{
				"PHL", strconv.Itoa(int(m.Week)), homeTeam.Abbreviation, strconv.Itoa(int(m.HomeTeamScore)),
				awayTeam.Abbreviation, strconv.Itoa(int(m.AwayTeamScore)), otStr, soStr, strconv.Itoa(int(m.HomeTeamShootoutScore)), strconv.Itoa(int(m.AwayTeamShootoutScore)), m.HomeTeamCoach,
				strconv.Itoa(int(m.HomeTeamRank)), m.AwayTeamCoach, strconv.Itoa(int(m.AwayTeamRank)), m.GameTitle,
				neutralStr, confStr, m.GameDay, m.Arena, strconv.Itoa(int(m.AttendanceCount)), m.City, m.State, m.Country, thirdStarStr, secondStarStr, firstStarStr,
			}
			err = csvW.Write(row)
			if err != nil {
				log.Fatal("Cannot write row to CSV", err)
			}

			csvW.Flush()
			err = csvW.Error()
			if err != nil {
				log.Fatal("Error while writing to file ::", err)
			}
		}
		return csvW.Error()
	})

}
