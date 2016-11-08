package hoeffding

type Trace struct {
	SplitPredictor string
	MeritGain      float64
	HoeffdingBound float64
	PossibleSplits []TracePossibleSplit
}

type TracePossibleSplit struct {
	Predictor string
	Merit     float64
}
