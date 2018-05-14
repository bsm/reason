package ftrl

import (
	"math"
	"strconv"

	"github.com/bsm/reason/core"
	"github.com/cespare/xxhash"
)

func featureKV(feat *core.Feature, x core.Example) (key string, val float64) {
	if feat.Kind.IsCategorical() {
		if cat := feat.Category(x); cat > -1 {
			key = feat.Name + "." + strconv.FormatUint(uint64(cat), 10)
			val = 1.0
		}
	} else if feat.Kind.IsNumerical() {
		if num := feat.Number(x); !math.IsNaN(num) {
			key = feat.Name
			val = num
		}
	}
	return
}

func featureBV(feat *core.Feature, x core.Example, numBuckets uint32) (int, float64) {
	key, val := featureKV(feat, x)
	if key == "" {
		return -1, 0
	}
	return int(xxhash.Sum64String(key) % uint64(numBuckets)), val
}
