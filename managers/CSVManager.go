package managers

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/structs"
)

func WritePlayByPlayCSVFile(playByPlays []structs.PlayByPlay, filename string, collegePlayerMap map[uint]structs.CollegePlayer, teamMap map[uint]structs.CollegeTeam) error {
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
		possessingTeam := teamMap[uint(play.TeamID)]
		zone := getZoneLabel(play.ZoneID)
		abbr := possessingTeam.Abbreviation
		penalty := getPenaltyByID(uint(play.PenaltyID))
		severity := getSeverityByID(play.Severity)
		isFight := "No"
		if play.IsFight {
			isFight = "Yes"
		}

		result := generateResultsString(play, event, outcome, collegePlayerMap, possessingTeam)
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

func FormatTimeToClock(timeInSeconds uint16) string {
	minutes := timeInSeconds / 60
	seconds := timeInSeconds % 60

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func generateResultsString(play structs.PlayByPlay, event, outcome string, playerMap map[uint]structs.CollegePlayer, possessingTeam structs.CollegeTeam) string {
	puckCarrier := playerMap[play.PuckCarrierID]
	puckCarrierLabel := getPlayerLabel(puckCarrier)
	receivingPlayer := playerMap[play.PassedPlayerID]
	receivingPlayerLabel := getPlayerLabel(receivingPlayer)
	assistingPlayer := playerMap[play.AssistingPlayerID]
	assistingPlayerLabel := getPlayerLabel(assistingPlayer)
	defendingPlayer := playerMap[play.DefenderID]
	defendingPlayerLabel := getPlayerLabel(defendingPlayer)
	goalie := playerMap[play.GoalieID]
	goalieLabel := getPlayerLabel(goalie)
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

func getPlayerLabel(player structs.CollegePlayer) string {
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
