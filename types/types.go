package types

import (
	"time"
)

// MonitorSKU holds SKU information needed for just-in-time calculations.
type MonitorSKU struct {
	Sold    int
	Days    int
	Then    time.Time
	Pending bool
}
