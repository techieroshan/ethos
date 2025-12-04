package testdata

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"ethos/internal/auth/model"
	"ethos/internal/feedback/model"
)

// TestDataService handles creation and cleanup of test data
type TestDataService struct {
	// Database or repository dependencies would go here
}

// NewTestDataService creates a new test data service
func NewTestDataService() *TestDataService {
	return &TestDataService{}
}

// CreateTestData creates all test users and their contextual data
func (s *TestDataService) CreateTestData(ctx context.Context) (map[string]interface{}, error) {
	profiles := GetTestUserProfiles()
	createdUsers := make([]map[string]interface{}, 0)
	createdFeedback := make([]map[string]interface{}, 0)
	createdRatings := make([]map[string]interface{}, 0)
	createdReviews := make([]map[string]interface{}, 0)

	// For each test user, create their profile, feedback, ratings, and reviews
	for userType, profile := range profiles {
		// Create user profile
		userPayload := map[string]interface{}{
			"id":        profile.ID,
			"email":     profile.Email,
			"firstName": profile.FirstName,
			"lastName":  profile.LastName,
			"role":      profile.Role,
			"bio":       profile.Bio,
			"avatarUrl": fmt.Sprintf("https://i.pravatar.cc/150?img=%d", profile.AvatarIndex),
			"createdAt": profile.CreatedAt,
		}
		createdUsers = append(createdUsers, userPayload)

		// Create contextual feedback items for this user
		feedbackItems := s.generateContextualFeedback(profile, profiles)
		createdFeedback = append(createdFeedback, feedbackItems...)

		// Create ratings for this user
		ratings := s.generateContextualRatings(profile, profiles)
		createdRatings = append(createdRatings, ratings...)

		// Create reviews for this user
		reviews := s.generateContextualReviews(profile, profiles)
		createdReviews = append(createdReviews, reviews...)

		fmt.Printf("[TESTDATA] Created test user: %s (%s) with %d feedback, %d ratings, %d reviews\n",
			profile.FullName, userType, profile.FeedbackCount, profile.RatingsCount, profile.ReviewsCount)
	}

	result := map[string]interface{}{
		"status":          "created",
		"users_count":     len(createdUsers),
		"feedback_count":  len(createdFeedback),
		"ratings_count":   len(createdRatings),
		"reviews_count":   len(createdReviews),
		"users":           createdUsers,
		"feedbackItems":   createdFeedback,
		"ratings":         createdRatings,
		"reviews":         createdReviews,
		"createdAt":       time.Now(),
		"testDataMarkers": s.generateTestDataMarkers(),
	}

	return result, nil
}

// generateContextualFeedback generates feedback items for a user
func (s *TestDataService) generateContextualFeedback(author *TestUserProfile, allProfiles map[string]*TestUserProfile) []map[string]interface{} {
	feedback := make([]map[string]interface{}, 0)
	rand.Seed(time.Now().UnixNano())

	// Get a random subset of other users as feedback recipients
	recipients := s.getRandomRecipients(allProfiles, author.ID, author.FeedbackCount)

	for i := 0; i < author.FeedbackCount; i++ {
		if i >= len(recipients) {
			break
		}

		recipient := recipients[i]
		ratings := map[string]int{
			"Integrity":      rand.Intn(10) + 1,
			"Leadership":     rand.Intn(10) + 1,
			"Communication":  rand.Intn(10) + 1,
			"Collaboration":  rand.Intn(10) + 1,
			"Technical":      rand.Intn(10) + 1,
		}

		feedbackItem := map[string]interface{}{
			"id":        fmt.Sprintf("feedback-%s-%d", author.ID, i),
			"authorId":  author.ID,
			"author":    author.FullName,
			"recipientId": recipient.ID,
			"recipient": recipient.FullName,
			"content":   s.generateFeedbackContent(author, recipient),
			"type":      "peer-feedback",
			"visibility": "public",
			"ratings":   ratings,
			"tags":      s.generateFeedbackTags(),
			"replies":   rand.Intn(5),
			"createdAt": time.Now().Add(-time.Duration(rand.Intn(90)) * 24 * time.Hour),
		}
		feedback = append(feedback, feedbackItem)
	}

	return feedback
}

// generateContextualRatings generates ratings from a user to others
func (s *TestDataService) generateContextualRatings(rater *TestUserProfile, allProfiles map[string]*TestUserProfile) []map[string]interface{} {
	ratings := make([]map[string]interface{}, 0)
	rand.Seed(time.Now().UnixNano())

	recipients := s.getRandomRecipients(allProfiles, rater.ID, rater.RatingsCount)

	categories := []string{"Integrity", "Leadership", "Communication", "Collaboration", "Technical", "Innovation"}

	for i := 0; i < rater.RatingsCount; i++ {
		if i >= len(recipients) {
			break
		}

		recipient := recipients[i]
		rating := map[string]interface{}{
			"id":        fmt.Sprintf("rating-%s-%d", rater.ID, i),
			"userId":    rater.ID,
			"ratedUserId": recipient.ID,
			"category":  categories[rand.Intn(len(categories))],
			"score":     rand.Intn(10) + 1,
			"feedback":  s.generateRatingFeedback(),
			"createdAt": time.Now().Add(-time.Duration(rand.Intn(90)) * 24 * time.Hour),
		}
		ratings = append(ratings, rating)
	}

	return ratings
}

// generateContextualReviews generates review/testimonials from a user
func (s *TestDataService) generateContextualReviews(reviewer *TestUserProfile, allProfiles map[string]*TestUserProfile) []map[string]interface{} {
	reviews := make([]map[string]interface{}, 0)
	rand.Seed(time.Now().UnixNano())

	recipients := s.getRandomRecipients(allProfiles, reviewer.ID, reviewer.ReviewsCount)

	for i := 0; i < reviewer.ReviewsCount; i++ {
		if i >= len(recipients) {
			break
		}

		recipient := recipients[i]
		review := map[string]interface{}{
			"id":        fmt.Sprintf("review-%s-%d", reviewer.ID, i),
			"authorId":  reviewer.ID,
			"subjectId": recipient.ID,
			"title":     s.generateReviewTitle(reviewer, recipient),
			"content":   s.generateReviewContent(reviewer, recipient),
			"rating":    rand.Intn(5) + 1,
			"helpful":   rand.Intn(10),
			"createdAt": time.Now().Add(-time.Duration(rand.Intn(90)) * 24 * time.Hour),
		}
		reviews = append(reviews, review)
	}

	return reviews
}

// Helper functions to generate realistic content

func (s *TestDataService) generateFeedbackContent(author, recipient *TestUserProfile) string {
	templates := []string{
		fmt.Sprintf("%s demonstrated excellent %s skills on the recent project. %s's attention to detail and %s mindset were key to our success.",
			recipient.FirstName, s.getRandomSkill(), recipient.FirstName, s.getRandomTrait()),
		fmt.Sprintf("Working with %s was a great experience. %s brought strong %s and %s to the team.",
			recipient.FirstName, recipient.FirstName, s.getRandomSkill(), s.getRandomTrait()),
		fmt.Sprintf("%s shows great promise in %s. Their %s approach and %s are impressive.",
			recipient.FirstName, s.getRandomArea(), s.getRandomTrait(), s.getRandomQuality()),
		fmt.Sprintf("I had the pleasure of collaborating with %s on several initiatives. %s's %s and %s made a significant impact.",
			recipient.FirstName, recipient.FirstName, s.getRandomTrait(), s.getRandomQuality()),
	}
	return templates[rand.Intn(len(templates))]
}

func (s *TestDataService) generateFeedbackTags() []string {
	allTags := []string{"Teamwork", "Innovation", "Leadership", "Communication", "Reliability", "Creative", "Problem-solver", "Quick-learner"}
	tagCount := rand.Intn(3) + 1
	tags := make([]string, 0)
	used := make(map[int]bool)

	for len(tags) < tagCount {
		idx := rand.Intn(len(allTags))
		if !used[idx] {
			tags = append(tags, allTags[idx])
			used[idx] = true
		}
	}
	return tags
}

func (s *TestDataService) generateRatingFeedback() string {
	feedbacks := []string{
		"Consistently delivers high-quality work",
		"Excellent communication and collaboration skills",
		"Strong technical expertise",
		"Great problem-solving ability",
		"Demonstrates leadership potential",
		"Reliable and dependable team member",
		"Proactive and takes initiative",
		"Adapts well to challenges",
	}
	return feedbacks[rand.Intn(len(feedbacks))]
}

func (s *TestDataService) generateReviewTitle(author, subject *TestUserProfile) string {
	templates := []string{
		fmt.Sprintf("Great working with %s on the team", subject.FirstName),
		fmt.Sprintf("%s is an excellent %s", subject.FirstName, subject.Role),
		fmt.Sprintf("Highly recommend %s for any project", subject.FirstName),
		fmt.Sprintf("%s brings amazing energy to the team", subject.FirstName),
	}
	return templates[rand.Intn(len(templates))]
}

func (s *TestDataService) generateReviewContent(author, subject *TestUserProfile) string {
	templates := []string{
		fmt.Sprintf("%s is a fantastic colleague and very easy to work with. Highly skilled and always ready to help the team.", subject.FirstName),
		fmt.Sprintf("It was a pleasure working with %s. %s brought excellent skills and a positive attitude to every interaction.", subject.FirstName, subject.FirstName),
		fmt.Sprintf("%s exceeded expectations in every way. A true professional with a passion for excellence.", subject.FirstName),
	}
	return templates[rand.Intn(len(templates))]
}

func (s *TestDataService) getRandomSkill() string {
	skills := []string{"technical", "leadership", "communication", "project management", "creative thinking", "analytical"}
	return skills[rand.Intn(len(skills))]
}

func (s *TestDataService) getRandomTrait() string {
	traits := []string{"collaborative spirit", "can-do attitude", "dedication", "attention to detail", "innovative mindset", "professional approach"}
	return traits[rand.Intn(len(traits))]
}

func (s *TestDataService) getRandomArea() string {
	areas := []string{"backend development", "frontend design", "system architecture", "product management", "team leadership", "technical mentoring"}
	return areas[rand.Intn(len(areas))]
}

func (s *TestDataService) getRandomQuality() string {
	qualities := []string{"attention to detail", "creative thinking", "strong communication", "problem-solving ability", "technical expertise", "leadership potential"}
	return qualities[rand.Intn(len(qualities))]
}

func (s *TestDataService) getRandomRecipients(allProfiles map[string]*TestUserProfile, exceptID string, count int) []*TestUserProfile {
	var recipients []*TestUserProfile
	for _, profile := range allProfiles {
		if profile.ID != exceptID {
			recipients = append(recipients, profile)
		}
	}

	// Shuffle and limit to count
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(recipients), func(i, j int) {
		recipients[i], recipients[j] = recipients[j], recipients[i]
	})

	if len(recipients) > count {
		recipients = recipients[:count]
	}
	return recipients
}

func (s *TestDataService) generateTestDataMarkers() map[string]interface{} {
	return map[string]interface{}{
		"isTestData":      true,
		"createdBy":       "test-data-seed-endpoint",
		"environment":     "dev-staging",
		"totalUsersCount": 10,
		"timestamp":       time.Now().Unix(),
	}
}

// CleanupTestData removes all test data
func (s *TestDataService) CleanupTestData(ctx context.Context) (map[string]interface{}, error) {
	testIDs := GetAllTestUserIDs()

	result := map[string]interface{}{
		"status":        "cleaned",
		"deletedUsers":  len(testIDs),
		"userIds":       testIDs,
		"timestamp":     time.Now().Unix(),
		"message":       "All test data has been removed successfully",
	}

	return result, nil
}
