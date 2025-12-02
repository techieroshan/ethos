package converter

import (
	"time"

	peoplepb "ethos/api/proto/people"
	authModel "ethos/internal/auth/model"
)

// ProtoToUserProfile converts proto UserProfile to domain model
func ProtoToUserProfile(pb *peoplepb.UserProfile) *authModel.UserProfile {
	if pb == nil {
		return nil
	}

	profile := &authModel.UserProfile{
		ID:            pb.Id,
		Email:         pb.Email,
		Name:          pb.Name,
		EmailVerified: pb.EmailVerified,
		PublicBio:     pb.PublicBio,
	}

	// Convert created_at timestamp
	if pb.CreatedAt != nil {
		profile.CreatedAt = pb.CreatedAt.AsTime()
	} else {
		profile.CreatedAt = time.Now()
	}

	// Convert updated_at timestamp (optional)
	if pb.UpdatedAt != nil {
		updatedAt := pb.UpdatedAt.AsTime()
		profile.UpdatedAt = &updatedAt
	}

	return profile
}

