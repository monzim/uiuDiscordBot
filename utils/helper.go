package utils

const MAX_DESCRIPTION_LENGTH = 2000

func ConstructDescription(mentionRoleID, summary string) string {
	maxSummaryLength := MAX_DESCRIPTION_LENGTH - len(SUPPORT_MESSAGE) - len(mentionRoleID) - 2
	if len(summary) > maxSummaryLength {
		summary = summary[:maxSummaryLength-3] + "..."
	}

	description := ""
	if mentionRoleID != "" {
		description = mentionRoleID + " "
	}
	description += summary + "\n" + SUPPORT_MESSAGE

	if len(description) > MAX_DESCRIPTION_LENGTH {
		description = description[:MAX_DESCRIPTION_LENGTH]
	}

	return description
}
