package hoeffding

import "github.com/bsm/reason/classifiers"

// Config configures behaviour
type Config struct {
	// The number of training instances a leaf node should observe
	// between split attempts.
	// Default: 200
	GracePeriod int

	// The split criterion to use for classifications.
	SplitCriterion classifiers.SplitCriterion

	// The allowable error in a split decision - values closer
	// to zero will take longer to decide.
	// Default: 0.0000001
	SplitConfidence float64

	// Threshold below which a split will be forced to break ties
	// Default: 0.05
	TieThreshold float64
}

func (c *Config) norm(isRegression bool) {
	if c.GracePeriod <= 0 {
		c.GracePeriod = 200
	}
	if c.SplitConfidence <= 0 {
		c.SplitConfidence = 1e-7
	}
	if c.TieThreshold <= 0 {
		c.TieThreshold = 0.05
	}
	if c.SplitCriterion == nil {
		c.SplitCriterion = classifiers.DefaultSplitCriterion(isRegression)
	}
}
