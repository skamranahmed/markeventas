package utils

type ReplyType string

const (
	UserAccountDoesNotExist         ReplyType = "USER_ACCOUNT_DOES_NOT_EXIST"
	UserGoogleCalendarConsentAbsent ReplyType = "USER_GOOGLE_CALENDAR_CONSENT_ABSENT"
	UserGoogleCalendarEventCreated  ReplyType = "USER_GOOGLE_CALENDAR_EVENT_CREATED"
)

func (r ReplyType) Body() string {
	switch r {
	case UserAccountDoesNotExist:
		return `@%s Looks like you do not have not an account with us yet. You can create one by visiting here: `

	case UserGoogleCalendarConsentAbsent:
		return `@%s Looks like you have not provided the access to your Google Calendar yet. You can do it by visiting here: `

	case UserGoogleCalendarEventCreated:
		return `@%s The event has been created in your calendar. Here's the link: %s`
	}
	return ""
}
