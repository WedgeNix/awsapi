package types

import (
	"time"
)

//D2sVendor maches the stucture of the JSON file that keeps track of vendor processed data.
type D2sVendor struct {
	Days      int
	Processed bool
}

// BananasMonSKU holds SKU information needed for just-in-time calculations.
type BananasMonSKU struct {
	Sold            int
	Days            int
	LastUTC         time.Time
	Pending         time.Time
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
