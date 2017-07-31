package types

import (
	"time"
)

// BananasMonSKU holds SKU information needed for just-in-time calculations.
type BananasMonSKU struct {
	Sold    int
	LastUTC time.Time
	Pending bool
}

// BananasCfg holds the program configuration for hit-the-bananas.
type BananasCfg struct {
	LastUTC       time.Time
	PODays        []time.Weekday
	OrdXDaysWorth int
}
