package ftrl

import "github.com/bsm/reason"

// Mock regression with acceptable values ranging between 0 and 1.
// Values < 0.5 indicate that outcome 0 is most likely while values >= 0.5 suggest that it's not.
type classification float64

// Category returns the more likely category.
func (v classification) Category() reason.Category {
	if v < 0.5 {
		return 0
	}
	return 1
}

// Prob returns the probability of the given category.
// Only categories 0 and 1 may yield results > 0.
func (v classification) Prob(cat reason.Category) float64 {
	if cat == 0 {
		return 1 - float64(v)
	} else if reason.IsCat(cat) {
		return float64(v)
	}
	return 0.0
}
