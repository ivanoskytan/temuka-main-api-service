package dto

type CreateReportRequest struct {
	CommunityID     int    `json:"community_id"`
	ReportedUserID  int    `json:"reported_user_id"`
	CommunityRuleID *int   `json:"community_rule_id,omitempty"`
	TargetType      string `json:"target_type"`
	TargetID        int    `json:"target_id"`
	Reason          string `json:"reason"`
}

type UpdateReportRequest struct {
	CommunityID     int    `json:"community_id"`
	ReportedUserID  int    `json:"reported_user_id"`
	CommunityRuleID *int   `json:"community_rule_id,omitempty"`
	TargetType      string `json:"target_type"`
	TargetID        int    `json:"target_id"`
	Reason          string `json:"reason"`
}
