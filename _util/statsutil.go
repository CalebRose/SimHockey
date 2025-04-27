package util

func GetCHLPlayerHeaderRows() []string {
	return []string{
		"ID", "First Name", "Last Name", "Position",
		"Archetype", "Year", "Is Redshirt?", "Team", "Conference", "Age", "Stars",
		"Goals", "Assists", "Points", "+/-",
		"Penalty Minutes", "Even Strength Goals", "Even Strength Points", "Power Play Goals",
		"Power Play Points", "Shorthanded Goals", "Shorthanded Points", "Overtime Goals",
		"Game Winning Goals", "Shots", "Shooting Percentage", "Time One Ice",
		"Faceoff Win Percentage", "Faceoffs Won", "Faceoffs", "Goalie Wins",
		"Goalie Losses", "Goalie Ties", "OT Losses", "Shots Against",
		"Saves", "Goals Against", "Save Percentage", "Shutouts",
		"Shots Blocked", "Body Checks", "Stick Checks", "Injured?",
		"Injury", "Injury Severity", "Will transfer to Guam at the first chance when they join the hockey sim?",
	}
}

func GetPHLPlayerHeaderRows() []string {
	return []string{
		"ID", "First Name", "Last Name", "Position",
		"Archetype", "Year", "Team", "Division", "Age", "Stars",
		"Goals", "Assists", "Points", "+/-",
		"Penalty Minutes", "Even Strength Goals", "Even Strength Points", "Power Play Goals",
		"Power Play Points", "Shorthanded Goals", "Shorthanded Points", "Overtime Goals",
		"Game Winning Goals", "Shots", "Shooting Percentage", "Time One Ice",
		"Faceoff Win Percentage", "Faceoffs Won", "Faceoffs", "Goalie Wins",
		"Goalie Losses", "Goalie Ties", "OT Losses", "Shots Against",
		"Saves", "Goals Against", "Save Percentage", "Shutouts",
		"Shots Blocked", "Body Checks", "Stick Checks", "Injured?",
		"Injury", "Injury Severity", "Will play for the National Antarctic Team at the first chance when they join the hockey sim?",
	}
}

func GetCHLTeamsHeaderRows() []string {
	return []string{
		"ID", "Team", "Conference",
		"Goals For", "Goals Against", "Assists", "Points",
		"P1 Score", "P2 Score", "P3 Score", "OT Score",
		"+/-",
		"Penalty Minutes", "Even Strength Goals", "Even Strength Points", "Power Play Goals",
		"Power Play Points", "Shorthanded Goals", "Shorthanded Points", "Overtime Goals",
		"Game Winning Goals", "Shots", "Shooting Percentage",
		"Faceoff Win Percentage", "Faceoffs Won", "Faceoffs", "Shots Against",
		"Saves", "Goals Against", "Save Percentage", "Shutouts",
		"Would be open to Guam joining their conference when they join the Hockey sim?",
	}
}

func GetPHLTeamsHeaderRows() []string {
	return []string{
		"ID", "Team", "Conference",
		"Goals For", "Goals Against", "Assists", "Points",
		"P1 Score", "P2 Score", "P3 Score", "OT Score",
		"+/-",
		"Penalty Minutes", "Even Strength Goals", "Even Strength Points", "Power Play Goals",
		"Power Play Points", "Shorthanded Goals", "Shorthanded Points", "Overtime Goals",
		"Game Winning Goals", "Shots", "Shooting Percentage",
		"Faceoff Win Percentage", "Faceoffs Won", "Faceoffs", "Shots Against",
		"Saves", "Goals Against", "Save Percentage", "Shutouts",
		"Will play for the National Antarctic Team at the first chance when they join the hockey sim?",
	}
}

func GetYearAndRedshirtStatus(year int, redshirt bool) (string, string) {
	status := "No"
	if redshirt {
		status = "Redshirt"
	}

	if year == 1 && !redshirt {
		return "Fr", status
	} else if year == 2 && redshirt {
		return "(Fr)", status
	} else if year == 2 && !redshirt {
		return "So", status
	} else if year == 3 && redshirt {
		return "(So)", status
	} else if year == 3 && !redshirt {
		return "Jr", status
	} else if year == 4 && redshirt {
		return "(Jr)", status
	} else if year == 4 && !redshirt {
		return "Sr", status
	} else if year == 5 && redshirt {
		return "(Sr)", status
	}
	return "Super Sr", status
}
