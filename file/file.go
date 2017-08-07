package file

import "github.com/WedgeNix/awsapi/types"

// Any represents any package-level type of AWS file.
type Any interface {
	__()
}

// BananasMon maps SKUs to their respective just-in-time data.
type BananasMon struct {
	AvgWait float64
	SKUs    types.SKUs
}

func (_ BananasMon) __() {}

// BananasCfg maps SKUs to their respective just-in-time data.
type BananasCfg types.BananasCfg

// BananasCfgName holds the file name.
const BananasCfgName = "hit-the-bananas/config.json"

func (_ BananasCfg) __() {}
