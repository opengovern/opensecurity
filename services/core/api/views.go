package api

import (
	"time"
)

type GetViewsCheckpointResponse struct {
	Checkpoint time.Time `json:"checkpoint"`
}

type View struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Query        Query    `json:"query"`
	Dependencies []string `json:"dependencies"`
}

type GetViewsResponse struct {
	TotalCount int    `json:"total_count"`
	Views      []View `json:"views"`
}

type Query struct {
	ID             string       `json:"id"`
	QueryToExecute string       `json:"query_to_execute"`
	PrimaryTable   *string      `json:"primary_table"`
	ListOfTables   []string     `json:"list_of_tables"`
	Engine         string       `json:"engine"`
	Parameters     []QueryParameter `json:"parameters"`
	Global         bool         `json:"global"`
	CreatedAt      time.Time        `json:"createdAt" example:"2023-06-07T14:00:15.677558Z"`
	UpdatedAt      time.Time        `json:"updatedAt" example:"2023-06-16T14:58:08.759554Z"`
}

