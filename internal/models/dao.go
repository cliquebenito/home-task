package models

import "time"

type StatsResponse struct {
	TS time.Time `json:"ts"`
	V  int       `json:"v"`
}

type StatsQueryParams struct {
	BannerID int
	From     time.Time
	To       time.Time
}
