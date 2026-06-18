package scraper

// MapLocationRank converts raw classification strings from sources to internal ranks.
func MapLocationRank(source, sourceType string) string {
	if sourceType == "Acadex" {
		switch source {
		case "A":
			return "Rank S"
		case "X":
			return "Rank A"
		case "Y":
			return "Rank B"
		case "Z":
			return "Rank C"
		default:
			return "Rank D"
		}
	}

	if sourceType == "IHappy" {
		switch source {
		case "signature":
			return "Rank S"
		case "prestige":
			return "Rank A"
		case "super premium":
			return "Rank B"
		case "premium":
			return "Rank C"
		case "promotion":
			return "Rank D"
		default:
			return "Rank D"
		}
	}

	return "Rank D"
}
