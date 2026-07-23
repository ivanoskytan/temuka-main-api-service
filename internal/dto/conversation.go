package dto

type AddConversationRequest struct {
	Title          string `json:"title"`
	ParticipantIDs []int  `json:"participant_ids"`
	UserID         int    `json:"user_id"`
}

type AddMessageRequest struct {
	ParticipantID  int    `json:"participant_id"`
	ConversationID int    `json:"conversation_id"`
	UserID         int    `json:"user_id"`
	Text           string `json:"text"`
}

type AddParticipantRequest struct {
	ConversationID int `json:"conversation_id"`
	UserID         int `json:"user_id"`
}
