package types

import (
	"time"
)

// BananasMonSKU holds SKU information needed for just-in-time calculations.
type BananasMonSKU struct {
	Sold            int
	Days            int
	LastUTC         time.Time
	Pending         bool
	ProbationPeriod int
}

// BananasCfg holds the program configuration for hit-the-bananas.
type BananasCfg struct {
	LastLA        time.Time
	PODays        []time.Weekday
	OrdXDaysWorth int
}

// SKUs maps SKUs to their respective monitor data.
type SKUs map[string]BananasMonSKU
