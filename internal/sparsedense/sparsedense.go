package sparsedense

// BetterOffDense returns true if something should be dense.
func BetterOffDense(count, capacity int) bool {
	return count > 100 && count*3 > capacity
}
