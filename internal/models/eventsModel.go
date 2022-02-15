package models

type Event struct {
	Name string
	Caption string
}



var Events = []Event{
	{"NEW_USER", "User added"},
	{"USER_ALREADY_CREATE", "User is already registered"},
	{"USER_NOT_FOUND", "User is not registered"},
	{"INVALID_PASSWORD", "Incorrect password specified"},
	{"TRY_USER_BLOCKED", "Check if user is blocked"},
	{"USER_BLOCKED", "User is blocked"},
	{"USER_UNBLOCKED", "User is unblocked"},
	{"NEW_SESSION", "New session created"},
	{"INVALID_SESSION", "Malicious auth attempt"},
	{"SESSION_EXPIRED", "Refresh token is too old"},
	{"UPDATED_SESSION", "Session is updated"},
	{"DELETED_SESSION", "Session is deleted"},
}
