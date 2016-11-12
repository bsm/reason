package hoeffding

import "github.com/bsm/reason/classifiers"

// Config configures behaviour
type Config struct {
	// The number of training instances a leaf node should observe
	// between split attempts.
	// Default: 200
	GracePeriod int

	// The number of training instances the tree should observe
	// between pruning attempts.
	// Default: 1000
	PrunePeriod int

	// The target heap size consumed by the tree. By default, trees are
	// allowed to grow twice the target heap size before they are pruned.
	// Please note that this is a rough estimate. Overall memory usage is
	// likely to be substantially higher.
	// Default: 128*1024*1024 (128MB)
	HeapTarget int

	// The split criterion to use for evaluating splits
	// Default: InformationGainSplitCriterion or VarReductionSplitCriterion
	SplitCriterion classifiers.SplitCriterion

	// The allowable error in a split decision - values closer
	// to zero will take longer to decide.
	// Default: 0.0000001
	SplitConfidence float64

	// Threshold below which a split will be forced to break ties
	// Default: 0.05
	TieThreshold float64

	// By enabling this option, tracing notification events will be
	// emitted via the Traces channel after each training cycle. This
	// is for debug purposes only. When enabled, you must consume
	// the Traces channel to avoid locked threads.
	// Default: false
	EnableTracing bool
}

func (c *Config) norm(isRegression bool) {
	if c.GracePeriod <= 0 {
		c.GracePeriod = 200
	}
	if c.PrunePeriod <= 0 {
		c.PrunePeriod = 1000
	}
	if c.HeapTarget <= 0 {
		c.HeapTarget = 128 * 1024 * 1024
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
