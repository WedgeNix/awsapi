package dir

import (
	"github.com/WedgeNix/awsapi/file"
)

// Any represents any package-level type of AWS directory.
type Any interface {
	__()
}

// BananasMon is a mapping of file names to monitor files.
type BananasMon map[string]file.BananasMon

func (_ BananasMon) __() {}
