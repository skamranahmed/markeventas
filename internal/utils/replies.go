package utils

var (
	userAccountDoesNotExistReplyBody         = `@%s Hey, looks like you do not have not an account with us yet. You can create one by visiting here: `
	userGoogleCalendarConsentAbsentReplyBody = `@%s Hey, looks like you have not provided the access to your Google Calendar yet. You can do it by visiting here: `
	userGoogleCalendarEventCreatedReplyBody  = `@%s The event has been created in your calendar. Here's the link: %s`
)

type ReplyType string

var (
	userAccountDoesNotExist         ReplyType
	userGoogleCalendarConsentAbsent ReplyType
	userGoogleCalendarEventCreated  ReplyType
)

type BotReplyType struct {
	Code int
	Name ReplyType
	Body string
}

var (

	// UserAccountDoesNotExistReply
	UserAccountDoesNotExistReply = BotReplyType{
		Code: 1,
		Name: userAccountDoesNotExist,
		Body: userAccountDoesNotExistReplyBody,
	}

	// UserGoogleCalendarConsentAbsentReply
	UserGoogleCalendarConsentAbsentReply = BotReplyType{
		Code: 2,
		Name: userGoogleCalendarConsentAbsent,
		Body: userGoogleCalendarConsentAbsentReplyBody,
	}

	// UserGoogleCalendarEventCreatedReply
	UserGoogleCalendarEventCreatedReply = BotReplyType{
		Code: 3,
		Name: userGoogleCalendarEventCreated,
		Body: userGoogleCalendarEventCreatedReplyBody,
	}
)
