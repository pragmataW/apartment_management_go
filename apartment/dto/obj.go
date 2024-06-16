package dto

import "sync"

var (
	Mutx      sync.Mutex
	DuesPrice = 40.0
	PayDay    = 15
)