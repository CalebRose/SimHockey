package structs

import "gorm.io/gorm"

// College Request & Pro Requests
type TeamRequest struct {
	TeamID     uint
	Username   string
	IsApproved bool
	IsActive   bool
	Role       string
}

type CollegeTeamRequest struct {
	gorm.Model
	TeamRequest
}

type ProTeamRequest struct {
	gorm.Model
	TeamRequest
}

func (r *TeamRequest) ApproveTeamRequest() {
	r.IsApproved = true
}

func (r *TeamRequest) RejectTeamRequest() {
	r.IsApproved = false
	r.IsActive = false
}

func (r *TeamRequest) Reactivate() {
	r.IsApproved = false
	r.IsActive = true
}

type TeamRequestsResponse struct {
	CollegeRequests []CollegeTeamRequest
	ProRequest      []ProTeamRequest
}
