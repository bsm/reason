package internal

// NodeStats is an abstract interface to distinguish between classification
// and regression stats.
type NodeStats interface {
	isNode_Stats
	WeightSum() float64
}

// NewNodeStats_Classification inits a new node stats instance.
func NewNodeStats_Classification() *Node_Classification {
	return &Node_Classification{Classification: &Node_ClassificationStats{}}
}

// WeightSum implementes NodeStats.
func (s *Node_Classification) WeightSum() float64 {
	return s.Classification.WeightSum()
}

// NewNodeStats_Regression inits a new node stats instance.
func NewNodeStats_Regression() *Node_Regression {
	return &Node_Regression{Regression: &Node_RegressionStats{}}
}

// WeightSum implementes NodeStats.
func (s *Node_Regression) WeightSum() float64 {
	return s.Regression.Weight
}
