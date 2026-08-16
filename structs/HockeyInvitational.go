package structs

import "gorm.io/gorm"

/*
 Hockey Invitationals List
 - Beanpot (Always Boston College, Boston U, Harvard, and Northeastern. Always at Madison Square Gardens)
 - Great Lakes Invitational
 - Lake Placid Invitational (Always Quinnipiac, Sacred Heart, Yale, and UConn)
 - Staten Island Invitational
 - Great North Invitational (Always include an Alaska school)
*/

// Static Data for Invitationals in the Simulation
type HockeyInvitational struct {
	gorm.Model
	Name                string
	ArenaID             string
	Arena               string
	TotalTeams          uint8 // Usually 4
	Week                uint8 // Reserve for week 3. Slots B and C.
	IsTournament        bool  // For Beanpot & Lake Placid only. All teams must face each other due to intraconference.
	InterConferenceOnly bool  // Restrict to interconference matchups only. Only one team per conference.
}

// Request Data for Teams Participating in Invitationals
type HockeyInvitationalRequest struct {
	gorm.Model
	TeamID         uint
	InvitationalID uint
	SeasonID       uint
	ConferenceID   uint
}
