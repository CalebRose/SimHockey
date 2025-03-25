package structs

import "gorm.io/gorm"

type FaceData struct {
	gorm.Model
	PlayerID        uint
	Accessories     uint8   // Range - 1 through 5
	Body            uint8   // Range - 1 through 3
	BodySize        float32 // Range - 0.8 through 1.2
	Ear             uint8   // Range - 1 through 3
	EarSize         float32 // 0.5 to 1.5
	Eye             uint8   // Range - 1 through 19
	EyeAngle        int8    // Range -10 to 15
	EyeLine         uint8   // Range - 1 through 7
	Eyebrow         uint8   // Range - 1 through 20
	EyeBrowAngle    int8    // Range -15 to 20
	FaceSize        float32 // 0 to 1
	FacialHair      uint8   // Range - 1 through 83
	FacialHairShave uint8   // 1 through 5 (make sure to define like rgba(0,0,0,0.x))
	Glasses         uint8   // Range - 1 through 7
	Hair            uint8   // Range - 1 through 49
	HairColor       uint8   //
	HairBG          uint8   // Range - 1 through 8 (feminine only?)
	HairFlip        bool
	Head            uint8 // Range - 1 through 22
	Jersey          uint8 // Range - 1 through 5 (specialized based on sport)
	MiscLine        uint8 // Range - 1 through 11
	Mouth           uint8 // Range - 1 through 18
	MouthFlip       bool
	Nose            uint8 // Range - 1 through 16
	NoseFlip        bool
	NoseSize        float32 // Range 0.5 to 1.25
	SmileLine       uint8   // Range - 1 through 5
	SmileLineSize   float32 // 0.25 to 2.25
	SkinTone        string
	SkinColor       uint8
}

// Will need to come up with a range of allowable characteristics for this view.

type FaceDataResponse struct {
	PlayerID        uint
	Accessories     string
	Body            string
	Ear             string
	Eye             string
	EyeLine         string
	Eyebrow         string
	FacialHair      string
	Glasses         string
	Hair            string
	HairBG          string
	HairFlip        bool
	Head            string
	Jersey          string
	MiscLine        string
	Mouth           string
	MouthFlip       bool
	Nose            string
	NoseFlip        bool
	SmileLine       string
	BodySize        float32 // Range - 0.8 through 1.2
	EarSize         float32 // 0.5 to 1.5
	EyeAngle        int8    // Range -10 to 15
	EyeBrowAngle    int8    // Range -15 to 20
	FaceSize        float32 // 0 to 1
	FacialHairShave string  // 1 through 5 (make sure to define like rgba(0,0,0,0.x))
	NoseSize        float32 // Range 0.5 to 1.25
	SmileLineSize   float32 // 0.25 to 2.25
	SkinColor       string
	HairColor       string
}
