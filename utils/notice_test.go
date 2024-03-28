package utils_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/monzim/uiuBot/utils"
)

func TestConstructDescriptionV2(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	// Test case 1: Summary within the limit
	mentionRoleID := "@role123"
	summary := generateRandomString(100)
	expectedDescription := mentionRoleID + " " + summary + "\n" + utils.SUPPORT_MESSAGE
	description := utils.ConstructDescription(mentionRoleID, summary)
	if description != expectedDescription {
		t.Errorf("Description mismatch. Summary within the limit. \nExpected: %s, \nGot: %s", expectedDescription, description)
	}

	// Test case 2: Summary exceeding the limit
	mentionRoleID = "@role456"
	summary = generateRandomString(3000)
	description = utils.ConstructDescription(mentionRoleID, summary)
	if len(description) > utils.MAX_DESCRIPTION_LENGTH {
		t.Errorf("Description length exceeds the limit. \nExpected: %d, \nGot: %d", utils.MAX_DESCRIPTION_LENGTH, len(description))
	}

	// Test case 3: Summary equal to the limit
	mentionRoleID = "@role789"
	summary = generateRandomString(1993)
	description = utils.ConstructDescription(mentionRoleID, summary)
	if len(description) != utils.MAX_DESCRIPTION_LENGTH {
		t.Errorf("Description length doesn't match the limit. \nExpected: %d, \nGot: %d", utils.MAX_DESCRIPTION_LENGTH, len(description))
	}

}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
