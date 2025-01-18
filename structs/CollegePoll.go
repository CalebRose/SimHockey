package structs

import "gorm.io/gorm"

type CollegePollSubmission struct {
	gorm.Model
	Username string
	SeasonID uint
	WeekID   uint
	Week     uint
	Rank1    string
	Rank1ID  uint
	Rank2    string
	Rank2ID  uint
	Rank3    string
	Rank3ID  uint
	Rank4    string
	Rank4ID  uint
	Rank5    string
	Rank5ID  uint
	Rank6    string
	Rank6ID  uint
	Rank7    string
	Rank7ID  uint
	Rank8    string
	Rank8ID  uint
	Rank9    string
	Rank9ID  uint
	Rank10   string
	Rank10ID uint
	Rank11   string
	Rank11ID uint
	Rank12   string
	Rank12ID uint
	Rank13   string
	Rank13ID uint
	Rank14   string
	Rank14ID uint
	Rank15   string
	Rank15ID uint
	Rank16   string
	Rank16ID uint
	Rank17   string
	Rank17ID uint
	Rank18   string
	Rank18ID uint
	Rank19   string
	Rank19ID uint
	Rank20   string
	Rank20ID uint
}

type CollegePollOfficial struct {
	gorm.Model
	SeasonID       uint
	WeekID         uint
	Week           uint
	Rank1          string
	Rank1ID        uint
	Rank1Votes     uint
	Rank1No1Votes  uint
	Rank2          string
	Rank2ID        uint
	Rank2Votes     uint
	Rank2No1Votes  uint
	Rank3          string
	Rank3ID        uint
	Rank3Votes     uint
	Rank3No1Votes  uint
	Rank4          string
	Rank4ID        uint
	Rank4Votes     uint
	Rank4No1Votes  uint
	Rank5          string
	Rank5ID        uint
	Rank5Votes     uint
	Rank5No1Votes  uint
	Rank6          string
	Rank6ID        uint
	Rank6Votes     uint
	Rank6No1Votes  uint
	Rank7          string
	Rank7ID        uint
	Rank7Votes     uint
	Rank7No1Votes  uint
	Rank8          string
	Rank8ID        uint
	Rank8Votes     uint
	Rank8No1Votes  uint
	Rank9          string
	Rank9ID        uint
	Rank9Votes     uint
	Rank9No1Votes  uint
	Rank10         string
	Rank10ID       uint
	Rank10Votes    uint
	Rank10No1Votes uint
	Rank11         string
	Rank11ID       uint
	Rank11Votes    uint
	Rank11No1Votes uint
	Rank12         string
	Rank12ID       uint
	Rank12Votes    uint
	Rank12No1Votes uint
	Rank13         string
	Rank13ID       uint
	Rank13Votes    uint
	Rank13No1Votes uint
	Rank14         string
	Rank14ID       uint
	Rank14Votes    uint
	Rank14No1Votes uint
	Rank15         string
	Rank15ID       uint
	Rank15Votes    uint
	Rank15No1Votes uint
	Rank16         string
	Rank16ID       uint
	Rank16Votes    uint
	Rank16No1Votes uint
	Rank17         string
	Rank17ID       uint
	Rank17Votes    uint
	Rank17No1Votes uint
	Rank18         string
	Rank18ID       uint
	Rank18Votes    uint
	Rank18No1Votes uint
	Rank19         string
	Rank19ID       uint
	Rank19Votes    uint
	Rank19No1Votes uint
	Rank20         string
	Rank20ID       uint
	Rank20Votes    uint
	Rank20No1Votes uint
}

func (s *CollegePollSubmission) AssignID(id uint) {
	s.ID = id
}

func (s *CollegePollSubmission) MoveSubmissionToNextWeek(weekID, week uint) {
	s.WeekID = weekID
	s.Week = week
}

func (c *CollegePollOfficial) AssignRank(idx int, vote TeamVote) {
	if idx == 0 {
		c.Rank1 = vote.Team
		c.Rank1Votes = vote.TotalVotes
		c.Rank1ID = vote.TeamID
		c.Rank1No1Votes = vote.Number1Votes
	} else if idx == 1 {
		c.Rank2 = vote.Team
		c.Rank2Votes = vote.TotalVotes
		c.Rank2ID = vote.TeamID
		c.Rank2No1Votes = vote.Number1Votes
	} else if idx == 2 {
		c.Rank3 = vote.Team
		c.Rank3Votes = vote.TotalVotes
		c.Rank3ID = vote.TeamID
		c.Rank3No1Votes = vote.Number1Votes
	} else if idx == 3 {
		c.Rank4 = vote.Team
		c.Rank4Votes = vote.TotalVotes
		c.Rank4ID = vote.TeamID
		c.Rank4No1Votes = vote.Number1Votes
	} else if idx == 4 {
		c.Rank5 = vote.Team
		c.Rank5Votes = vote.TotalVotes
		c.Rank5ID = vote.TeamID
		c.Rank5No1Votes = vote.Number1Votes
	} else if idx == 5 {
		c.Rank6 = vote.Team
		c.Rank6Votes = vote.TotalVotes
		c.Rank6ID = vote.TeamID
		c.Rank6No1Votes = vote.Number1Votes
	} else if idx == 6 {
		c.Rank7 = vote.Team
		c.Rank7Votes = vote.TotalVotes
		c.Rank7ID = vote.TeamID
		c.Rank7No1Votes = vote.Number1Votes
	} else if idx == 7 {
		c.Rank8 = vote.Team
		c.Rank8Votes = vote.TotalVotes
		c.Rank8ID = vote.TeamID
		c.Rank8No1Votes = vote.Number1Votes
	} else if idx == 8 {
		c.Rank9 = vote.Team
		c.Rank9Votes = vote.TotalVotes
		c.Rank9ID = vote.TeamID
		c.Rank9No1Votes = vote.Number1Votes
	} else if idx == 9 {
		c.Rank10 = vote.Team
		c.Rank10Votes = vote.TotalVotes
		c.Rank10ID = vote.TeamID
		c.Rank10No1Votes = vote.Number1Votes
	} else if idx == 10 {
		c.Rank11 = vote.Team
		c.Rank11Votes = vote.TotalVotes
		c.Rank11ID = vote.TeamID
		c.Rank11No1Votes = vote.Number1Votes
	} else if idx == 11 {
		c.Rank12 = vote.Team
		c.Rank12Votes = vote.TotalVotes
		c.Rank12ID = vote.TeamID
		c.Rank12No1Votes = vote.Number1Votes
	} else if idx == 12 {
		c.Rank13 = vote.Team
		c.Rank13Votes = vote.TotalVotes
		c.Rank13ID = vote.TeamID
		c.Rank13No1Votes = vote.Number1Votes
	} else if idx == 13 {
		c.Rank14 = vote.Team
		c.Rank14Votes = vote.TotalVotes
		c.Rank14ID = vote.TeamID
		c.Rank14No1Votes = vote.Number1Votes
	} else if idx == 14 {
		c.Rank15 = vote.Team
		c.Rank15Votes = vote.TotalVotes
		c.Rank15ID = vote.TeamID
		c.Rank15No1Votes = vote.Number1Votes
	} else if idx == 15 {
		c.Rank16 = vote.Team
		c.Rank16Votes = vote.TotalVotes
		c.Rank16ID = vote.TeamID
		c.Rank16No1Votes = vote.Number1Votes
	} else if idx == 16 {
		c.Rank17 = vote.Team
		c.Rank17Votes = vote.TotalVotes
		c.Rank17ID = vote.TeamID
		c.Rank17No1Votes = vote.Number1Votes
	} else if idx == 17 {
		c.Rank18 = vote.Team
		c.Rank18Votes = vote.TotalVotes
		c.Rank18ID = vote.TeamID
		c.Rank18No1Votes = vote.Number1Votes
	} else if idx == 18 {
		c.Rank19 = vote.Team
		c.Rank19Votes = vote.TotalVotes
		c.Rank19ID = vote.TeamID
		c.Rank19No1Votes = vote.Number1Votes
	} else if idx == 19 {
		c.Rank20 = vote.Team
		c.Rank20Votes = vote.TotalVotes
		c.Rank20ID = vote.TeamID
		c.Rank20No1Votes = vote.Number1Votes
	}
}

type TeamVote struct {
	Team         string
	TeamID       uint
	TotalVotes   uint
	Number1Votes uint
}

func (t *TeamVote) AddVotes(num uint) {
	t.TotalVotes += (26 - num)
	if num == 1 {
		t.Number1Votes++
	}
}

type PollDataResponse struct {
	Poll          CollegePollSubmission
	Matches       []CollegeGame
	OfficialPolls []CollegePollOfficial
}
