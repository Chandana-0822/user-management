package services

import (
	"fmt"
)

func GenerateUsernameSuggestions(firstName, lastName string) ([]string, error) {
	// Define potential username patterns
	potentialUsernames := []string{
		fmt.Sprintf("%s.%s", firstName, lastName),
		fmt.Sprintf("%s_%s", firstName, lastName),
		fmt.Sprintf("%s%s", firstName, lastName),
		fmt.Sprintf("%s.%s123", firstName, lastName),
		fmt.Sprintf("%s_%s123", firstName, lastName),
		fmt.Sprintf("%s%s123", firstName, lastName),
	}

	uniqueUsernames := []string{}
	for _, username := range potentialUsernames {
		exists, err := SearchUsername(username)
		if err != nil {
			return nil, fmt.Errorf("error checking username: %w", err)
		}

		// Add only unique usernames
		if !exists {
			uniqueUsernames = append(uniqueUsernames, username)
		}

		// Stop if we have three unique suggestions
		if len(uniqueUsernames) >= 3 {
			break
		}
	}

	return uniqueUsernames, nil
}
