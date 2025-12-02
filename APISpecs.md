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

## Core Functionalities

* **User Registration:** Create new user accounts using email and password; initiate email verification.

* **Authentication (Login):** Authenticate credentials; issue JWT access and refresh tokens.

* **Email Verification:** Confirm user email ownership through secure tokens.

* **Password Management:** Change and reset passwords via email-based workflows.

* **Token Refresh:** Seamlessly refresh access tokens using valid refresh tokens.

* **Profile Retrieval:** Authenticated access to user profile details and personal information.

* **Profile Update:** Editable profile data and user preference management.

* **Account Security:** Two-factor authentication (2FA), security events log, and user data export.

*Example Scenarios:*

* A new user registers with their email, receives a verification link, and can then log in and manage their profile.

* An existing user requests a password reset, receives an email, and sets a new password securely.

* A mobile app fetches profile information for authenticated users to display personalized content.

## Architecture Overview

The Ethos API is architected as a stateless RESTful service, versioned at `/api/v1`, exposing endpoints for authentication and user profiles across environments (development, staging, production). Core architectural components include:

* **Authentication Service:** Manages login, registration, and token issuance/validation.

* **Profile Service:** Handles user profile retrieval and updates.

* **Token Management Layer:** Issues, validates, and refreshes JWTs for session management.

* **Security Management:** Endpoints for account-level security operations.

**Environment Base URLs (Examples):**

* Production: `https://api.ethos.example.com/api/v1`

* Staging: `https://staging.ethos.example.com/api/v1`

* Development: `http://localhost:8000/api/v1`

All requests and responses use JSON. The API is designed for horizontal scalability and high performance, supporting stateless interactions and optimized for both browser and mobile clients.

---

## Core Data Models

This section describes the canonical data types and structures used across Ethos API endpoints. Each model is shown as a TypeScript interface or type alias, with an explanation and references to relevant endpoints that return or accept this shape. All models are actively used by the corresponding endpoints.

### Primitive Type Aliases

type UserId = string; // Unique user identifier (e.g., "user-1234")  

type FeedbackId = string; // Unique feedback item identifier (e.g., "f-001")  

type CommentId = string; // Unique comment identifier (e.g., "c-1001")  

type NotificationId = string; // Unique notification identifier (e.g., "n-01")  

type Timestamp = string; // ISO 8601 formatted string, e.g. "2024-05-04T14:30:00Z"

Commonly referenced throughout all API models to ensure type safety and clarity.

---

### UserSummary

A minimal reference to a user, typically embedded within other models such as feedback or comments.

interface UserSummary { id: UserId; name: string; }

* **Used in:** `FeedbackItem`, `FeedbackComment`, and related endpoints for participant attribution.

---

### UserProfile

Extends `UserSummary` and surfaces additional profile details for the authenticated or searched user.

interface UserProfile extends UserSummary { email: string; email_verified: boolean; created_at: Timestamp; public_bio?: string; updated_at?: Timestamp; }

* **Returned by:** `GET /profile/me`, `GET /profile/:user_id`

* **Describes:** The complete profile for the requester or queried user.

---

### FeedbackVisibility and FeedbackType

Captures constraints and categorization for feedback.

type FeedbackVisibility = "public" | "private" | "team";  

type FeedbackType = "appreciation" | "suggestion" | "issue" | "other";

* **Used in:** `FeedbackItem` to control access and reporting, and allow for filtering by type.

---

### FeedbackDimensionScore

Optional dimension-level scoring for structured feedback.

interface FeedbackDimensionScore { dimension: string; // e.g., "clarity", "impact" score: number; // e.g., 1-5 }

* **Used in:** `FeedbackItem.dimensions` for advanced feedback analytics.

---

### FeedbackReactionSummary

Summarizes the reaction counts for a feedback item.

interface FeedbackReactionSummary { \[reaction: string\]: number; // e.g., "like": 3, "helpful": 2 }

* **Used in:** `FeedbackItem.reactions`

* **Returned by:** `GET /feedback/feed`, `GET /feedback/:feedback_id`

---

### FeedbackItem

The core model for all feedback-related data. Represents an individual feedback post.

interface FeedbackItem { feedback_id: FeedbackId; author: UserSummary; content: string; type?: FeedbackType; visibility?: FeedbackVisibility; reactions: FeedbackReactionSummary; dimensions?: FeedbackDimensionScore\[\]; comments_count: number; created_at: Timestamp; }

* **Returned by:**

  * `GET /feedback/feed` (lists)

  * `GET /profile/:user_id/feedback` (lists)

  * `GET /feedback/:feedback_id` (detail)

---

### FeedbackComment

Represents a comment or reply on a feedback post.

interface FeedbackComment { comment_id: CommentId; author: UserSummary; content: string; created_at: Timestamp; parent_comment_id?: CommentId; // for threading/replies }

* **Returned by:** `GET /feedback/:feedback_id/comments`

---

### DashboardSnapshot

Aggregated dashboard view for an authenticated user, including activity, statistics, and suggested actions.

interface DashboardSnapshot { recent_feedback: FeedbackItem\[\]; stats: { feedback_given: number; comments: number; \[key: string\]: number; }; suggested_actions: string\[\]; }

* **Returned by:** `GET /dashboard`

* **Purpose:** Personalized dashboard experience aggregating key user data.

---

### NotificationType

Defines the set of notification event categories.

type NotificationType = | "feedback_reply" | "feedback_received" | "new_comment" | "system_alert" | "reminder" | "other";

* **Used in:** `Notification.type`

---

### Notification

Represents a notification delivered to a user.

interface Notification { notification_id: NotificationId; type: NotificationType; message: string; read: boolean; created_at: Timestamp; }

* **Returned by:** `GET /notifications`

* **Purpose:** Track and display user notifications, including unread status.

---

### NotificationPreferences

Tracks delivery channel choices for user notifications.

interface NotificationPreferences { email: boolean; push: boolean; in_app: boolean; }

* **Returned by:** `GET /notifications/preferences`

* **Updated by:** `PUT /notifications/preferences`

* **Purpose:** Controls which notification channels a user has opted into.

---

## API Authentication

**JWT-based authentication** is enforced:

* On successful login, clients receive an **access token** (short-lived) and a **refresh token** (longer-lived).

* Protected endpoints require an `Authorization: Bearer <access_token>` header.

* Refresh tokens are used to obtain new access tokens without re-entering credentials.

* No API keys or session cookies are used.

**Security Best Practices:**

* Store tokens securely (not in localStorage if possible; prefer secure HTTP-only cookies on web).

* Always use HTTPS.

* Logout or invalidate refresh tokens on logout.

* Employ short expiry durations for access tokens and regular rotation of refresh tokens.

## Authentication Methods

Supported authentication methods:

1. **Email + Password Login**

  * **Step 1:** POST to `/auth/login` with email and password.

  * **Step 2:** Receive access and refresh tokens on success.

  **Request Example:**

  **POST** `/api/v1/auth/login`

  ```
  {
      "email": "user@example.com",
      "password": "UserPassword123!"
  }
```

  **Response Example:**

  ```
  {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..."
  }
```

2. **Token-based Access**

  * **Step 1:** Attach access token in `Authorization` header for protected requests.

  * **Example Header:**

    Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6...

## Token Management

* **Token Issuance:** After successful authentication, both access and refresh tokens are issued.

* **Access Token:** Short lifespan (e.g., 15 minutes), used for API calls.

* **Refresh Token:** Longer lifespan (e.g., 14 days), used at `/auth/refresh` endpoint.

* **Refreshing Tokens:**

  * POST your refresh token to `/auth/refresh` to receive a new access token.

* **Invalidation:** Log out or password change revokes refresh tokens.

**Example Refresh Request:**

POST /api/v1/auth/refresh { "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..." }

**Best Practices:**

* Store refresh tokens securely (not in client-side JavaScript if possible).

* Discard expired/invalid tokens promptly.

* Never expose tokens on URLs, logs, or browser storage that isn’t secure.

## Common Issues and Troubleshooting

**Frequent Causes:**

* Using expired or invalid access tokens

* Failing to include or misformatting the `Authorization` header

* Attempting to use a revoked or expired refresh token

* Accessing protected endpoints without login

**Quick Resolutions:**

* Always login again if tokens are expired

* Copy tokens exactly from login/refresh responses

* Ensure HTTP headers are correctly set

* Check server error messages for actionable hints

## Error Messages

The API uses standardized JSON error responses, typically in the format:

```
{
    "error": "Invalid credentials",
    "code": "AUTH_INVALID_CREDENTIALS"
}
```

Main error types include:

* **Authentication errors:** Invalid token, expired token, bad credentials, account not verified

* **Authorization errors:** Missing or malformed token, insufficient permissions

* **Resource errors:** User not found, email not verified, validation failures

* **Server errors:** Unexpected processing issues or outages

Each error contains a machine-readable code and human-friendly message to aid debugging.

## Error Code List

## Troubleshooting Guide

1. **Invalid Credentials on Login**

  * Double-check provided email and password.

  * Ensure the account exists and is verified.

2. **Access Token Expired**

  * Use the refresh token at `/auth/refresh` to get a new access token.

  * If refresh token is also expired, log in again.

3. **Invalid or Missing Authorization Header**

  * Verify the header is present and precisely formatted as `Bearer <token>`.

  * Do not include additional spaces or characters.

4. **Email Not Verified**

  * Ensure user has completed email verification.

  * Resend verification email if necessary.

5. **Profile Not Found**

  * Confirm successful login and correct user account.

  * New users may need to complete setup.

**General Best Practices:**

* Use recommended token storage methods.

* Always validate inputs before sending requests.

* Monitor for consistent error codes and handle gracefully.

## Support and Resources

* **Documentation Portal:** Comprehensive API docs and guides (Coming Soon)

* 

* **Community Forums:** forums.ethosapi.example.com (Coming Soon)

* **Knowledge Base:** FAQs and troubleshooting resources (Coming Soon)

*For integration or urgent production issues, please email the support address above.*

## API Endpoints and Operations

### Domain Grouping

## GET

### `/auth/me`

* **Status:** Current

* **Usage:** Retrieve information about the currently authenticated user, including email and verification status.

* **Authorization:** Required (Bearer access token)

* **Headers:**

  * `Authorization: Bearer <access_token>`

* **Request Example:**

  * GET `/api/v1/auth/me`

* **Sample Response:**

  ```
  {
      "id": "user-1234",
      "email": "user@example.com",
      "email_verified": true,
      "created_at": "2023-10-02T13:21:00Z"
  }
```

  *Response Body: UserProfile*

* **Errors:** AUTH_TOKEN_EXPIRED, AUTH_TOKEN_INVALID, AUTH_EMAIL_UNVERIFIED

### `/profile/me`

* **Status:** Current

* **Usage:** Retrieve authenticated user's profile.

* **Authorization:** Required (Bearer access token)

* **Headers:**

  * `Authorization: Bearer <access_token>`

* **Request Example:**

  * GET `/api/v1/profile/me`

* **Sample Response:**

  ```
  {
      "id": "user-1234",
      "email": "user@example.com",
      "name": "Jane Doe",
      "email_verified": true,
      "created_at": "2023-10-02T13:21:00Z"
  }
```

  *Response Body: UserProfile*

* **Errors:** PROFILE_NOT_FOUND, AUTH_TOKEN_EXPIRED

### `/profile/user-profile`

* **Status:** Current

* **Usage:** Search or list publicly accessible user profiles.

* **Authorization:** Optional, depending on profile visibility

* **Query Parameters:**

  * `q` (search term),

  * `limit` (default: 25),

  * `offset`

* **Request Example:**

  * GET `/api/v1/profile/user-profile?q=Jane`

* **Sample Response:**

  ```
  {
      "results": \[
          {
              "id": "user-5678",
              "name": "Jane Smith",
              "public_bio": "Educator, runner, mentor."
          }
      \],
      "count": 1
  }
```

  *Each user object conforms to UserProfile (with public fields only)*

* **Errors:** VALIDATION_FAILED, SERVER_ERROR

### `/account/security-events`

* **Status:** Current

* **Usage:** Retrieve a list of significant security-related events for the authenticated user’s account.

* **Authorization:** Required

* **Request Example:**

  * GET `/api/v1/account/security-events`

* **Sample Response:**

  ```
  {
      "events": \[
          {
              "event_id": "evt-01",
              "type": "login",
              "timestamp": "2024-01-02T15:09:10Z",
              "ip": "203.0.113.5",
              "location": "US"
          }
      \]
  }
```

* **Errors:** AUTH_TOKEN_EXPIRED, SERVER_ERROR

### `/account/export-data/:export_id/status`

* **Status:** Current

* **Usage:** Retrieve the status of a previously requested account data export.

* **Authorization:** Required

* **Request Example:**

  * GET `/api/v1/account/export-data/exp-12345/status`

* **Sample Response:**

  ```
  {
      "export_id": "exp-12345",
      "status": "completed",
      "download_url": "https://download.ethos.example.com/exp-12345.zip",
      "expires_at": "2024-04-03T10:00:00Z"
  }
```

* **Errors:** VALIDATION_FAILED, AUTH_TOKEN_EXPIRED, SERVER_ERROR

---

### `/feedback/feed`

* **Status:** Current

* **Purpose:** Retrieve a paginated feed of community feedback posts for the main feedback wall.

* **Authorization:** Required (Bearer access token)

* **Query Parameters:**

  * `limit` (default: 20)

  * `offset`

* **Request Example:**

  * GET `/api/v1/feedback/feed?limit=10`

* **Sample Response:**

  ```
  {
      "results": \[
          {
              "feedback_id": "f-001",
              "author": { "id": "user-234", "name": "Lisa K." },
              "content": "Really enjoying the new feature!",
              "reactions": { "like": 5, "helpful": 2 },
              "comments_count": 3,
              "created_at": "2024-05-04T14:30:00Z"
          }
      \],
      "count": 1
  }
```

  *Each item in results conforms to FeedbackItem*

### `/feedback/:feedback_id`

* **Status:** Current

* **Purpose:** Retrieve the detail of a specific feedback post, including content, author, total reactions, and comments summary.

* **Authorization:** Required (Bearer access token)

* **Request Example:**

  * GET `/api/v1/feedback/f-001`

* **Sample Response:**

  ```
  {
      "feedback_id": "f-001",
      "author": { "id": "user-234", "name": "Lisa K." },
      "content": "Really enjoying the new feature!",
      "reactions": { "like": 5, "helpful": 2 },
      "created_at": "2024-05-04T14:30:00Z"
  }
```

  *Response Body: FeedbackItem*

### `/feedback/:feedback_id/comments`

* **Status:** Current

* **Purpose:** List all comments and replies for a given feedback post.

* **Authorization:** Required (Bearer access token)

* **Request Example:**

  * GET `/api/v1/feedback/f-001/comments`

* **Sample Response:**

  ```
  {
      "comments": \[
          {
              "comment_id": "c-1001",
              "author": { "id": "user-111", "name": "Joan P." },
              "content": "I agree!",
              "created_at": "2024-05-05T08:22:00Z"
          }
      \],
      "count": 1
  }
```

  *Each comment in comments is a FeedbackComment*

### `/people/search`

* **Status:** Current

* **Purpose:** Search for people matching the query string.

* **Authorization:** Required (Bearer access token)

* **Query Parameters:**

  * `q` (search term), `limit`, `offset`

* **Request Example:**

  * GET `/api/v1/people/search?q=alex`

* **Sample Response:**

  ```
  {
      "results": \[
          {
              "id": "user-8910",
              "name": "Alex Johnson",
              "role": "Designer"
          }
      \],
      "count": 1
  }
```

### `/people/recommendations`

* **Status:** Current

* **Purpose:** Return a personalized list of people the user may want to connect with.

* **Authorization:** Required (Bearer access token)

* **Request Example:**

  * GET `/api/v1/people/recommendations`

* **Sample Response:**

  ```
  {
      "recommendations": \[
          {
              "id": "user-1527",
              "name": "Morgan C.",
              "role": "Engineer"
          }
      \]
  }
```

### `/dashboard`

* **Status:** Current

* **Purpose:** Retrieve a personalized dashboard snapshot for the authenticated user, including key activity, stats, and recommendations.

* **Authorization:** Required (Bearer access token)

* **Request Example:**

  * GET `/api/v1/dashboard`

* **Sample Response:**

  ```
  {
      "recent_feedback": \[ ... \],
      "stats": { "feedback_given": 10, "comments": 22 },
      "suggested_actions": \[ "Complete your profile", "Give feedback" \]
  }
```

  *Response Body: DashboardSnapshot*

### `/notifications`

* **Status:** Current

* **Purpose:** List notifications for the authenticated user (paginated).

* **Authorization:** Required (Bearer access token)

* **Query Parameters:**

  * `limit`, `offset`

* **Request Example:**

  * GET `/api/v1/notifications?limit=20`

* **Sample Response:**

  ```
  {
      "notifications": \[
          {
              "notification_id": "n-01",
              "type": "feedback_reply",
              "message": "You received a reply to your feedback.",
              "read": false,
              "created_at": "2024-05-07T09:00:00Z"
          }
      \],
      "unread_count": 1
  }
```

  *Each notification object conforms to Notification*

### `/notifications/preferences`

* **Status:** Current

* **Purpose:** Get the authenticated user's notification delivery preferences.

* **Authorization:** Required (Bearer access token)

* **Request Example:**

  * GET `/api/v1/notifications/preferences`

* **Sample Response:**

  ```
  {
      "preferences": {
          "email": true,
          "push": true,
          "in_app": true
      }
  }
```

  *Response Body: NotificationPreferences*

### `/community/rules`

* **Status:** Current

* **Purpose:** Retrieve the current community rules and acceptable use policy.

* **Authorization:** Not required

* **Request Example:**

  * GET `/api/v1/community/rules`

* **Sample Response:**

  ```
  {
      "title": "Community Rules",
      "content": "Please respect others and do not post prohibited content."
  }
```

## POST

*(No changes required for core model notation in this section)*

## PUT

### `/profile/me`

* **Status:** Current

* **Purpose:** Update authenticated user’s profile details (such as name, display information, or public bio).

* **Authorization:** Required

* **Headers:**

  * `Authorization: Bearer <access_token>`

* **Input Example:**

  ```
  {
      "name": "Jane A. Doe",
      "public_bio": "Lead engineer, coffee enthusiast."
  }
```

* **Sample Response:**

  ```
  {
      "id": "user-1234",
      "email": "user@example.com",
      "name": "Jane A. Doe",
      "public_bio": "Lead engineer, coffee enthusiast.",
      "updated_at": "2024-03-01T12:00:00Z"
  }
```

  *Response Body: UserProfile*

* **Errors:** VALIDATION_FAILED, AUTH_TOKEN_EXPIRED

### `/notifications/preferences`

* **Status:** Current

* **Purpose:** Update user's preferences for notification delivery channels.

* **Authorization:** Required (Bearer access token)

* **Input Example:**

  ```
  {
      "email": false,
      "push": true,
      "in_app": true
  }
```

* **Sample Response:**

  ```
  {
      "preferences": {
          "email": false,
          "push": true,
          "in_app": true
      }
  }
```

  *Response Body: NotificationPreferences*

## PATCH

### `/profile/me/preferences`

* **Status:** Current

* **Purpose:** Update user-specific profile preferences (notifications, locale, etc.).

* **Authorization:** Required

* **Headers:**

  * `Authorization: Bearer <access_token>`

* **Input Example:**

  ```
  {
      "notify_on_login": false,
      "locale": "en-US"
  }
```

* **Sample Response:**

  ```
  {
      "preferences": {
          "notify_on_login": false,
          "locale": "en-US"
      }
  }
```

* **Errors:** VALIDATION_FAILED, AUTH_TOKEN_EXPIRED

## DELETE

### `/profile/me`

* **Status:** Current

* **Purpose:** Request permanent deletion of the authenticated user’s account.

* **Authorization:** Required

* **Headers:**

  * `Authorization: Bearer <access_token>`

* **Request Example:**

  * DELETE `/api/v1/profile/me`

* **Sample Response:**

  ```
  {
      "message": "Account scheduled for deletion. You will receive a confirmation email."
  }
```

* **Errors:** AUTH_TOKEN_EXPIRED, SERVER_ERROR

### `/auth/setup-2fa`

* **Status:** Current

* **Purpose:** Disable or remove two-factor authentication from the user account.

* **Authorization:** Required

* **Headers:**

  * `Authorization: Bearer <access_token>`

* **Request Example:**

  * DELETE `/api/v1/auth/setup-2fa`

* **Sample Response:**

  ```
  {
      "message": "Two-factor authentication disabled for your account."
  }
```

* **Errors:** AUTH_TOKEN_EXPIRED, SERVER_ERROR

---

### `/feedback/:feedback_id/react`

* **Status:** Current

* **Purpose:** Remove a reaction from a given feedback post for the current user.

* **Authorization:** Required (Bearer access token)

* **Request Example:**

  * DELETE `/api/v1/feedback/f-101/react`

* **Sample Response:**

  ```
  {
      "feedback_id": "f-101",
      "message": "Reaction removed"
  }
```

---

*End of Ethos API Documentation*