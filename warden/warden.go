package warden

import (
	"time"
)

type Team struct {
	Name                 string    `json:"name"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	OwnerID              int       `json:"owner_id"`
	ParentTeamIdentifier string    `json:"parent_team_identifier"`
	Identifier           string    `json:"identifier"`
	ProductGroupName     string    `json:"product_group_name"`
	ProductGroupID       string    `json:"product_group_id"`
	Labels               any       `json:"labels"`
	ShortCode            string    `json:"short_code"`
}

type TeamResponse struct {
	Success bool `json:"success"`
	Data    Team `json:"data"`
}

type TeamListResponse struct {
	Success bool      `json:"success"`
	Data    TeamsData `json:"data"`
}

type TeamsData struct {
	Teams []Team `json:"teams"`
}

type TeamListRequest struct {
	Email string
}

type TeamByUUIDRequest struct {
	TeamUUID string
}
