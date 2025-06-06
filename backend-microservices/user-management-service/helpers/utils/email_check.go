package utils

import (
	"context"
	"log"
	"microservice/user/data-layer/repository"

	"github.com/google/uuid"
)

// CheckEmailExistsForOtherUser verifies if an email exists for any user except the one with the given ID
func CheckEmailExistsForOtherUser(ctx context.Context, repo repository.UserRepository, email string, excludeUserID uuid.UUID) (bool, error) {
	// First check if email exists at all
	exists, err := repo.ExistsByEmailOrUsername(ctx, nil, email, "")
	if err != nil {
		log.Printf("Error checking email existence: %v", err)
		return false, err
	}

	// If email doesn't exist, we're good
	if !exists {
		return false, nil
	}

	// If email exists, check if it belongs to the user being updated
	user, err := repo.FindUserByEmail(ctx, nil, email)
	if err != nil {
		// If we get an error here but the exists check was true,
		// it likely means the email exists for another user
		log.Printf("Error finding user by email: %v", err)
		return true, nil
	}

	// Check if the found user is the same as the one being updated
	if user.ID == excludeUserID {
		// The email belongs to the same user, so no conflict
		return false, nil
	}

	// Email exists for another user
	return true, nil
}
