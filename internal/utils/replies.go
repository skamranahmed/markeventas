package utils

var (
	userAccountDoesNotExistBody         = `@%s Hey, looks like you do not have not an account with us yet. You can create one by visiting here: `
	userGoogleCalendarConsentAbsentBody = `@%s Hey, looks like you have not provided the access to your Google Calendar yet. You can do it by visiting here: `
)

type ReplyType string

var (
	userAccountDoesNotExist         ReplyType
	userGoogleCalendarConsentAbsent ReplyType
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
		Body: userAccountDoesNotExistBody,
	}

	// UserGoogleCalendarConsentAbsentReply
	UserGoogleCalendarConsentAbsentReply = BotReplyType{
		Code: 2,
		Name: userGoogleCalendarConsentAbsent,
		Body: userGoogleCalendarConsentAbsentBody,
	}
)
