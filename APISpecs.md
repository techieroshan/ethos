# Ethos API Documentation

## Overview

The Ethos API is a secure RESTful interface designed to handle user authentication and profile management in modern applications. It provides robust, scalable endpoints supporting operations such as user registration, login, email verification, password management, and user profile retrieval, enabling seamless integration for web and mobile platforms. Built on token-based security leveraging JWTs, the API enforces industry best practices for data protection and authentication workflows.

## API Purpose

The Ethos API is designed to simplify and standardize authentication and user profile management for client applications. It solves the challenges of:

* Securely verifying user identities

* Streamlining the onboarding process (registration and email verification)

* Enabling secure password management and recovery

* Allowing controlled profile data access and user information updates

**Primary use cases include:**

* User sign-up and login flows

* Self-serve password reset

* Verifying email addresses during registration

* Fetching authenticated user profile data for personalization or account management

---

## Core Data Models

This section describes the canonical data types and structures used across Ethos API endpoints. Below are model extensions and new models referenced in the enhanced endpoints.

### Updated: FeedbackItem

**Now includes:**

* Reaction analytics

* Helpfulness metric

* Follow-up tracking

* Anonymous flag

```typescript
interface FeedbackItem {
  feedback_id: FeedbackId;
  author: UserSummary;
  content: string;
  type?: FeedbackType;
  visibility: FeedbackVisibility;
  reactions: FeedbackReactionSummary;
  helpfulness?: number; // New: percentage or score of helpful feedback
  reactions_analytics?: {
    \[reaction: string\]: {
      count: number;
      user_ids: UserId\[\];
    };
  };
  follow_ups?: FeedbackFollowUp\[\]; // Tracked follow-up discussions
  is_anonymous?: boolean; // New: Identifies anonymous submissions
  dimensions?: FeedbackDimensionScore\[\];
  comments_count: number;
  created_at: Timestamp;
}

```

### New: FeedbackTemplate

```typescript
interface FeedbackTemplate {
  template_id: string;
  name: string;
  description: string;
  context_tags: string\[\];
  template_fields: any; // JSON schema describing template content fields
}

```

### New: FeedbackFollowUp

```typescript
interface FeedbackFollowUp {
  follow_up_id: string;
  content: string;
  author: UserSummary;
  created_at: Timestamp;
}

```

### Updated: ModerationState (Richer States)

```typescript
type ModerationState =
  | "pending"
  | "warned"
  | "actioned"
  | "escalated"
  | "appealed";

```

---

## API Endpoints and Operations

### Feedback

Feedback Templates & Suggestions

GET `/api/feedback/templates`

**Description:** Retrieve a list of available feedback templates, with contextual details for suggestion.

* **Query Parameters:**

  * `context` (optional): Provide context (e.g., "performance review", "peer feedback")

  * `tags` (optional): Comma-separated list of tags to filter templates

**Request Example:**

GET `/api/feedback/templates?context=performance_review`

**Response Example:**

```json
{
  "results": \[
    {
      "template_id": "t-001",
      "name": "Appreciation Template",
      "description": "A short message to acknowledge great work.",
      "context_tags": \["appreciation", "general"\],
      "template_fields": {
        "fields": \[
          { "name": "what_went_well", "type": "text", "label": "What went well?" }
        \]
      }
    }
  \]
}

```

---

POST `/api/feedback/template_suggestions`

**Description:** Suggest a new template or signal an in-context template need.

* **Body:**

  * `suggested_by` (UserId): (Optional if user is authenticated)

  * `usage_context` (string): Context for which the template is needed (e.g., "1:1 meetings", "onboarding")

  * `details` (string): Free-text description of the needed template

  * `desired_fields` (optional): Array of desired field names/types

**Request Example:**

```json
{
  "usage_context": "peer review",
  "details": "Need a short, positive feedback template for quick peer recognition.",
  "desired_fields": \[
    { "name": "positive_note", "type": "text" }
  \]
}

```

**Response Example:**

```json
{
  "status": "suggestion_received"
}

```

---

Impact & Analytics

Model Enhancements

* `helpfulness` (number): Score or ratio representing community evaluation of this feedback.

* `reactions_analytics` (object): For each reaction, a count and optionally a list of user IDs who reacted.

* `follow_ups` (array): List of follow-up messages on the feedback.

* All new fields are returned in `GET /feedback/feed`, `GET /feedback/:feedback_id`, and relevant listing endpoints.

GET `/api/feedback/impact`

**Description:** Aggregate analytics on feedback collected, including reactions, helpfulness over time, and follow-up trends.

* **Query Parameters:**

  * `user_id` (optional): Filter impact by feedback received by a specific user

  * `from`, `to` (optional): Date range for aggregation (ISO dates)

**Request Example:**

GET `/api/feedback/impact?user_id=user-8822&from=2024-01-01&to=2024-07-01`

**Response Example:**

```json
{
  "feedback_count": 31,
  "average_helpfulness": 0.87,
  "reaction_totals": { "like": 120, "helpful": 53, "insightful": 12 },
  "follow_up_count": 7,
  "trends": \[
    { "date": "2024-06-01", "helpfulness": 0.91, "feedback_submitted": 4 }
  \]
}

```

---

Anonymous and Batch Feedback

Model Update

* `is_anonymous` (boolean): Indicates if feedback was submitted without surface author identity. Anonymous feedbacks mask the `author` field in certain responses based on permission.

POST `/api/feedback/batch`

**Description:** Batch create feedback items in a single operation, supporting anonymous and mixed-visibility submissions.

**Request Example:**

```json
{
  "items": \[
    {
      "content": "Great presentation in the meeting!",
      "type": "appreciation",
      "visibility": "public",
      "is_anonymous": false
    },
    {
      "content": "Consider slowing down during the Q&A.",
      "type": "suggestion",
      "visibility": "org",
      "is_anonymous": true
    }
  \]
}

```

**Response Example:**

```json
{
  "submitted": \[
    { "feedback_id": "f-741", "status": "created" },
    { "feedback_id": "f-742", "status": "created" }
  \]
}

```

---

### People / Search

Advanced Feed & People Filters

Enhanced Filters for Feedback and People

`GET /api/feedback/feed` and `GET /people/search` now accept:

* `reviewer_type`: `"public"` or `"org"` (to filter by type of feedback reviewer)

* `context`: Filter feedback/people by context (e.g., "project", "team", "initiative")

* `verification`: `"verified"`, `"unverified"` (filter results accordingly)

* `tags`: Comma-separated list to match feedback or user tags

Example: GET `/api/feedback/feed?reviewer_type=org&context=team-building&tags=leadership,initiative`

Reviewer Context in Models

**All applicable models now include:**

* `reviewer_context` (object): Rich JSON depicting the review scope, e.g.:

  * `{"type": "org", "tenant_id": "org-1234", "department": "Engineering"}`

  * This is returned with all feedback search and listing endpoints.

---

Bookmark and Export Endpoints

GET `/api/feedback/bookmarks`

* **Description:** Get paginated list of feedback items the user has bookmarked.

* **Sample Response:** 

  ```json
  {
    "results": \[
      {
        "feedback_id": "f-53091",
        "content": "Excellent mentoring this quarter!",
        "reviewer_context": { "type": "org", "tenant_id": "org-12" }
      }
    \]
  }
  
```

POST `/api/feedback/bookmarks/:feedback_id`

* **Description:** Bookmark or remove bookmark from a feedback item.

* **Body:** `{ "action": "add" | "remove" }`

GET `/api/feedback/export`

* **Description:** Export feedback matching filters (by type, tag, reviewer, date, etc.)

* **Query Parameters:** Accepts the same filter set as `/feedback/feed`

* **Response:** Download/export file (CSV, JSON), delivered as direct payload or via a downloadable export asset.

---

### Moderation / Admin

Moderation Appeal

POST `/api/moderation/appeals`

* **Description:** Submit an appeal for a moderation decision.

* **Body:**

  * `moderated_item_id` (string)

  * `item_type` ("feedback" | "comment" | "profile" | ...)

  * `reason` (string)

  * `details` (optional)

* **Request Example:** 

  ```json
  {
    "moderated_item_id": "f-2098",
    "item_type": "feedback",
    "reason": "Feedback was incorrectly removed",
    "details": "The post complies with all community guidelines."
  }
  
```

* **Response Example:** 

  ```json
  {
    "status": "appeal_submitted",
    "appeal_id": "a-10233"
  }
  
```

---

Rich Moderation State

Moderation status for any moderated item is now one of:

* `pending`: Awaiting review

* `warned`: User has been warned, item may still be visible/limited

* `actioned`: Item was removed or restricted due to violation

* `escalated`: Issue referred to higher-level admin

* `appealed`: Appeal is under review

This field is present in all moderation-related response objects.

---

GET `/api/moderation/context`

* **Description:** Get the moderation context and active rules for a given item.

* **Query Parameters:** `item_id`, `item_type`

**Request Example:**  

GET `/api/moderation/context?item_id=f-2098&item_type=feedback`

**Response Example:**

```json
{
  "item_id": "f-2098",
  "item_type": "feedback",
  "current_state": "warned",
  "rules_applied": \[
    { "rule_id": "r-offensive-language", "description": "Avoid offensive or discriminatory language.", "status": "applied" }
  \],
  "reviewer_notes": "Content was borderline but repeated offenses were observed."
}

```

---

Moderation Workflow

1. **Pending:** All new flagged items enter `pending` review.

2. **Warned:** Moderator can issue a warning, with item still visible but limited.

3. **Actioned:** Item restricted or removed after review.

4. **Escalated:** For serious/complex cases, item is escalated to a senior moderator or admin.

5. **Appealed:** If appealed by the user, status changes to `appealed` until re-review.

6. **Outcome:** Review decision is delivered to user, status is updated in feedback or item record.

---

### Privacy / Data Control

Data Control & Opt-Out Endpoints

POST `/api/profile/opt-out`

* **Description:** User requests to opt-out from certain features or data processing (e.g., public profile search).

* **Body:**

  * `from`: which data/feature to opt out from (e.g., `"public_search"`, `"analytics_use"`)

  * `reason` (optional)

* **Request Example:**

  ```json
  {
    "from": "public_search",
    "reason": "Prefer not to show up in company-wide searches."
  }
  
```

* **Response Example:**

  ```json
  { "status": "opted_out", "changed": true }
  
```

---

POST `/api/profile/anonymize`

* **Description:** User requests anonymization of their personal data to meet privacy regulations.

* **Response:** Success confirmation, eventual propagated effect across all references (e.g., their feedback appears as from "Anonymous User").

**Response Example:**

```json
{ "status": "in_progress", "expected_completion": "2024-07-15T12:00:00Z" }

```

---

POST `/api/profile/delete_request`

* **Description:** User requests complete deletion of profile and associated data. This triggers an irreversible deletion workflow.

* **Request Example:**

  ```json
  { "confirm": true, "reason": "Leaving the company" }
  
```

* **Response Example:**

  ```json
  {
    "status": "delete_requested",
    "expected_completion": "2024-07-20T18:00:00Z"
  }
  
```

---

### Scenario and Model Updates

All relevant user, feedback, and moderation data models reflect new optional boolean fields for anonymity and opt-out request tracking:

* `UserProfile.opt_outs: string[]`

* `FeedbackItem.is_anonymous: boolean`

* Moderation objects include detailed workflow states (`pending`, `warned`, etc.)

---

## Summary Table of Enhancements

---

*This documentation is ready for implementation by both product and engineering teams, with all new endpoints and model changes action-ready.*