package repository

import "errors"

var (
	ErrStatsQueryFailed     = errors.New("failed to execute stats query")
	ErrStatsScanFailed      = errors.New("failed to scan stats row")
	ErrBannerNameExists     = errors.New("banner name already exists")
	ErrCreateBannerFailed   = errors.New("failed to create banner")
	ErrSaveStatisticsFailed = errors.New("failed to save statistics")
)
