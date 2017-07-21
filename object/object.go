package object

import (
	"time"
)

// Any represents any package-level type.
type Any interface {
	__()
}

// Monitor holds SKU information needed for just-in-time calculations.
type Monitor struct {
	Sold    int
	Days    int
	Then    time.Time
	Pending bool
}

func (_ Monitor) __() {}

// Monitors maps SKUs to their respective just-in-time data.
type Monitors map[string]Monitor

func (_ Monitors) __() {}
