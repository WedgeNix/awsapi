package file

import "github.com/WedgeNix/awsapi/types"

// Any represents any package-level type of AWS file.
type Any interface {
	__()
}

// BananasMon maps SKUs to their respective just-in-time data.
type BananasMon struct {
	AvgWait float64
	SKUs    map[string]types.BananasMonSKU
}

func (_ BananasMon) __() {}

// BananasCfg maps SKUs to their respective just-in-time data.
type BananasCfg types.BananasCfg

func (_ BananasCfg) __() {}
