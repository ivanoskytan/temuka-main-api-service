package constant

const (
	SuggestionExchange = "temuka_suggestion_exchange"

	SuggestionSyncRoutingKey = "suggestion.sync"

	EventOperationCreate = "CREATE"
	EventOperationDelete = "DELETE"
	EventOperationUpdate = "UPDATE"
	EventOperationView   = "VIEW"
	EventOperationLike   = "LIKE"
	EventOperationJoin   = "JOIN"

	EventEntityTypePost       = "post"
	EventEntityTypeUser       = "user"
	EventEntityTypeCommunity  = "community"
	EventEntityTypeUniversity = "university"
	EventEntityTypeMajor      = "major"
)
