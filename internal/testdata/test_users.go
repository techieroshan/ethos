package testdata

import (
	"time"
)

// TestUserType defines different user roles for test data
type TestUserType string

const (
	CEO                TestUserType = "CEO"
	VP_ENGINEERING     TestUserType = "VP_ENGINEERING"
	MANAGER            TestUserType = "MANAGER"
	SENIOR_DEV         TestUserType = "SENIOR_DEV"
	MID_LEVEL_DEV      TestUserType = "MID_LEVEL_DEV"
	JUNIOR_DEV         TestUserType = "JUNIOR_DEV"
	UX_DESIGNER        TestUserType = "UX_DESIGNER"
	HR_DIRECTOR        TestUserType = "HR_DIRECTOR"
	SALES_LEAD         TestUserType = "SALES_LEAD"
	SUPPORT_SPECIALIST TestUserType = "SUPPORT_SPECIALIST"
	FEEDBACK_AUTHOR_1  TestUserType = "FEEDBACK_AUTHOR_1"
	FEEDBACK_AUTHOR_2  TestUserType = "FEEDBACK_AUTHOR_2"
)

// TestUserProfile defines a test user with contextual data volumes
type TestUserProfile struct {
	ID            string
	FirstName     string
	LastName      string
	FullName      string
	Email         string
	Password      string
	Role          string
	Bio           string
	AvatarIndex   int
	FeedbackCount int
	RatingsCount  int
	ReviewsCount  int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TestFeedbackData represents feedback created by a test user
type TestFeedbackData struct {
	ID          string
	AuthorID    string
	RecipientID string
	Content     string
	Type        string
	Visibility  string
	Tags        []string
	Ratings     map[string]int // e.g., {"Integrity": 8, "Leadership": 7}
	Replies     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TestRatingData represents a rating given by a test user
type TestRatingData struct {
	ID          string
	UserID      string
	RatedUserID string
	Category    string
	Score       int
	Feedback    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TestReviewData represents a review/testimonial from a test user
type TestReviewData struct {
	ID        string
	AuthorID  string
	SubjectID string
	Title     string
	Content   string
	Rating    int
	Helpful   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetTestUserProfiles returns all 12 test user profiles
// Matches frontend test-avatars.ts exactly
func GetTestUserProfiles() map[TestUserType]*TestUserProfile {
	now := time.Now()

	return map[TestUserType]*TestUserProfile{
		CEO: {
			ID:            "test-user-ceo",
			FirstName:     "Jane",
			LastName:      "Smith",
			FullName:      "Jane Smith",
			Email:         "test-jane@ethos.local",
			Password:      "TestPass123!",
			Role:          "CEO",
			Bio:           "Chief Executive Officer - Leading innovation and strategic direction",
			AvatarIndex:   1,
			FeedbackCount: 30,
			RatingsCount:  15,
			ReviewsCount:  8,
			CreatedAt:     now.Add(-365 * 24 * time.Hour),
			UpdatedAt:     now,
		},
		VP_ENGINEERING: {
			ID:            "test-user-vp-eng",
			FirstName:     "John",
			LastName:      "Davis",
			FullName:      "John Davis",
			Email:         "test-john@ethos.local",
			Password:      "TestPass123!",
			Role:          "VP Engineering",
			Bio:           "Vice President of Engineering - Driving technical excellence",
			AvatarIndex:   2,
			FeedbackCount: 25,
			RatingsCount:  12,
			ReviewsCount:  6,
			CreatedAt:     now.Add(-350 * 24 * time.Hour),
			UpdatedAt:     now,
		},
		MANAGER: {
			ID:            "test-user-manager",
			FirstName:     "Sarah",
			LastName:      "Chen",
			FullName:      "Sarah Chen",
			Email:         "test-sarah@ethos.local",
			Password:      "TestPass123!",
			Role:          "Manager",
			Bio:           "Engineering Manager - Building high-performing teams",
			AvatarIndex:   3,
			FeedbackCount: 20,
			RatingsCount:  10,
			ReviewsCount:  5,
			CreatedAt:     now.Add(-300 * 24 * time.Hour),
			UpdatedAt:     now,
		},
		SENIOR_DEV: {
			ID:            "test-user-senior-dev",
			FirstName:     "Taylor",
			LastName:      "Jenkins",
			FullName:      "Taylor Jenkins",
			Email:         "test-taylor@ethos.local",
			Password:      "TestPass123!",
			Role:          "Senior Developer",
			Bio:           "Senior Software Engineer - Mentoring and architecting solutions",
			AvatarIndex:   4,
			FeedbackCount: 15,
			RatingsCount:  8,
			ReviewsCount:  4,
			CreatedAt:     now.Add(-250 * 24 * time.Hour),
			UpdatedAt:     now,
		},
		MID_LEVEL_DEV: {
			ID:            "test-user-mid-dev",
			FirstName:     "Alex",
			LastName:      "Rodriguez",
			FullName:      "Alex Rodriguez",
			Email:         "test-alex@ethos.local",
			Password:      "TestPass123!",
			Role:          "Mid-level Developer",
			Bio:           "Software Engineer - Contributing to core platform features",
			AvatarIndex:   5,
			FeedbackCount: 10,
			RatingsCount:  5,
			ReviewsCount:  3,
			CreatedAt:     now.Add(-200 * 24 * time.Hour),
			UpdatedAt:     now,
		},
		JUNIOR_DEV: {
			ID:            "test-user-junior-dev",
			FirstName:     "Emma",
			LastName:      "Wilson",
			FullName:      "Emma Wilson",
			Email:         "test-emma@ethos.local",
			Password:      "TestPass123!",
			Role:          "Junior Developer",
			Bio:           "Junior Software Engineer - Learning and growing with the team",
			AvatarIndex:   6,
			FeedbackCount: 8,
			RatingsCount:  4,
			ReviewsCount:  2,
			CreatedAt:     now.Add(-150 * 24 * time.Hour),
			UpdatedAt:     now,
		},
		UX_DESIGNER: {
			ID:            "test-user-designer",
			FirstName:     "Michael",
			LastName:      "Brown",
			FullName:      "Michael Brown",
			Email:         "test-michael@ethos.local",
			Password:      "TestPass123!",
			Role:          "UX Designer",
			Bio:           "Product Designer - Creating delightful user experiences",
			AvatarIndex:   7,
			FeedbackCount: 12,
			RatingsCount:  6,
			ReviewsCount:  3,
			CreatedAt:     now.Add(-180 * 24 * time.Hour),
			UpdatedAt:     now,
		},
		HR_DIRECTOR: {
			ID:            "test-user-hr",
			FirstName:     "Lisa",
			LastName:      "Anderson",
			FullName:      "Lisa Anderson",
			Email:         "test-lisa@ethos.local",
			Password:      "TestPass123!",
			Role:          "HR Director",
			Bio:           "Director of Human Resources - Building great teams and culture",
			AvatarIndex:   8,
			FeedbackCount: 18,
			RatingsCount:  9,
			ReviewsCount:  5,
			CreatedAt:     now.Add(-220 * 24 * time.Hour),
			UpdatedAt:     now,
		},
		SALES_LEAD: {
			ID:            "test-user-sales",
			FirstName:     "David",
			LastName:      "Martinez",
			FullName:      "David Martinez",
			Email:         "test-david@ethos.local",
			Password:      "TestPass123!",
			Role:          "Sales Lead",
			Bio:           "Sales Director - Driving business growth and partnerships",
			AvatarIndex:   9,
			FeedbackCount: 22,
			RatingsCount:  11,
			ReviewsCount:  6,
			CreatedAt:     now.Add(-280 * 24 * time.Hour),
			UpdatedAt:     now,
		},
		SUPPORT_SPECIALIST: {
			ID:            "test-user-support",
			FirstName:     "Rachel",
			LastName:      "Green",
			FullName:      "Rachel Green",
			Email:         "test-rachel@ethos.local",
			Password:      "TestPass123!",
			Role:          "Support Specialist",
			Bio:           "Customer Support Specialist - Ensuring customer success",
			AvatarIndex:   10,
			FeedbackCount: 14,
			RatingsCount:  7,
			ReviewsCount:  4,
			CreatedAt:     now.Add(-160 * 24 * time.Hour),
			UpdatedAt:     now,
		},
		// Feedback Authors - dedicated test personas for feedback cards
		FEEDBACK_AUTHOR_1: {
			ID:            "test-user-harold",
			FirstName:     "Harold",
			LastName:      "Stamm",
			FullName:      "Harold Stamm",
			Email:         "test-harold@ethos.local",
			Password:      "TestPass123!",
			Role:          "Tech Lead",
			Bio:           "Technical Lead at Tech Corp - Passionate about clean code and team collaboration",
			AvatarIndex:   11,
			FeedbackCount: 30,
			RatingsCount:  15,
			ReviewsCount:  8,
			CreatedAt:     now.Add(-400 * 24 * time.Hour),
			UpdatedAt:     now,
		},
		FEEDBACK_AUTHOR_2: {
			ID:            "test-user-cameron",
			FirstName:     "Cameron",
			LastName:      "Medhurst",
			FullName:      "Cameron Medhurst",
			Email:         "test-cameron@ethos.local",
			Password:      "TestPass123!",
			Role:          "Senior Engineer",
			Bio:           "Senior Software Engineer at Dev Studio - Building scalable solutions",
			AvatarIndex:   12,
			FeedbackCount: 25,
			RatingsCount:  12,
			ReviewsCount:  6,
			CreatedAt:     now.Add(-380 * 24 * time.Hour),
			UpdatedAt:     now,
		},
	}
}

// GetTestUserProfileByEmail returns a test user profile by email
func GetTestUserProfileByEmail(email string) *TestUserProfile {
	profiles := GetTestUserProfiles()
	for _, profile := range profiles {
		if profile.Email == email {
			return profile
		}
	}
	return nil
}

// GetTestUserProfileByID returns a test user profile by user ID
func GetTestUserProfileByID(userID string) *TestUserProfile {
	profiles := GetTestUserProfiles()
	for _, profile := range profiles {
		if profile.ID == userID {
			return profile
		}
	}
	return nil
}

// GetAllTestUserIDs returns all test user IDs
func GetAllTestUserIDs() []string {
	profiles := GetTestUserProfiles()
	ids := make([]string, 0, len(profiles))
	for _, profile := range profiles {
		ids = append(ids, profile.ID)
	}
	return ids
}

// IsTestUser checks if a user ID belongs to a test user
func IsTestUser(userID string) bool {
	testIDs := GetAllTestUserIDs()
	for _, id := range testIDs {
		if id == userID {
			return true
		}
	}
	return false
}

// GetSidebarPersonas returns the 5 personas used in the dashboard sidebar
// These are VP Engineering, Manager, Senior Dev, Mid Dev, Junior Dev
func GetSidebarPersonas() []*TestUserProfile {
	profiles := GetTestUserProfiles()
	return []*TestUserProfile{
		profiles[VP_ENGINEERING],
		profiles[MANAGER],
		profiles[SENIOR_DEV],
		profiles[MID_LEVEL_DEV],
		profiles[JUNIOR_DEV],
	}
}

// GetFeedbackAuthors returns the dedicated feedback author personas
func GetFeedbackAuthors() []*TestUserProfile {
	profiles := GetTestUserProfiles()
	return []*TestUserProfile{
		profiles[FEEDBACK_AUTHOR_1],
		profiles[FEEDBACK_AUTHOR_2],
	}
}

// GetTestUserProfilesAsStringMap returns all test user profiles keyed by user type string
// This is needed for compatibility with service.go
func GetTestUserProfilesAsStringMap() map[string]*TestUserProfile {
	profiles := GetTestUserProfiles()
	result := make(map[string]*TestUserProfile, len(profiles))
	for userType, profile := range profiles {
		result[string(userType)] = profile
	}
	return result
}
