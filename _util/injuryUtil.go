package util

func GetInjuryNameByID(injuryID uint8) string {
	switch injuryID {
	case 0:
		return "Concussion"
	case 1:
		return "Shoulder Separation"
	case 2:
		return "Broken Wrist"
	case 3:
		return "Broken Hand"
	case 4:
		return "Elbow Injury"
	case 5:
		return "Rib Injury"
	case 6:
		return "Back Strain"
	case 7:
		return "Cut/Laceration"
	case 8:
		return "Bruising"
	case 9:
		return "Groin Strain"
	case 10:
		return "Knee Sprain"
	case 11:
		return "Ankle Sprain"
	case 12:
		return "Hip Pointer"
	case 13:
		return "Hamstring Strain"
	case 14:
		return "General Soreness"
	default:
		return "Unknown Injury"
	}
}

func GetInjuryIDByName(injury string) uint8 {
	switch injury {
	case "Concussion":
		return 0
	case "Shoulder Separation":
		return 1
	case "Broken Wrist":
		return 2
	case "Broken Hand":
		return 3
	case "Elbow Injury":
		return 4
	case "Rib Injury":
		return 5
	case "Back Strain":
		return 6
	case "Cut/Laceration":
		return 7
	case "Bruising":
		return 8
	case "Groin Strain":
		return 9
	case "Knee Sprain":
		return 10
	case "Ankle Sprain":
		return 11
	case "Hip Pointer":
		return 12
	case "Hamstring Strain":
		return 13
	case "General Soreness":
		return 14
	default:
		return 255 // Unknown injury ID
	}
}

func GetInjurySeverityByID(severityID uint8) string {
	switch severityID {
	case 0:
		return "Minor"
	case 1:
		return "Moderate"
	case 2:
		return "Severe"
	default:
		return "Critical"
	}
}
