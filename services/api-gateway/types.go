package main

import "ride-sharing/shared/types"

type previewTripRequest struct {
	UserID string           `json:"userID"`
	Pickup types.Coordinate `json:"pickup"`
	Dest   types.Coordinate `json:"destination"`
}
