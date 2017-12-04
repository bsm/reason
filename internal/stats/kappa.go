package stats

// Kappa represents Cohen's kappa
type Kappa struct {
	m     [][]float64
	ncols int
}

// NewKappa inits a new Kappa
func NewKappa() *Kappa { return new(Kappa) }

// Record records in instance of expected vs actual with a given weight
func (k *Kappa) Record(expected, actual int, weight float64) {
	if expected > actual {
		k.grow(expected + 1)
	} else {
		k.grow(actual + 1)
	}

	k.m[actual][expected] += weight
}

// Value returns the kappa value
func (k *Kappa) Value() float64 {
	pws := make([]float64, len(k.m))
	tws := make([]float64, k.ncols)
	sum := 0.0

	for i, v := range k.m {
		for j, w := range v {
			pws[i] += w
			tws[j] += w
			sum += w
		}
	}
	if sum == 0.0 {
		return 0.0
	}

	var obs, exp float64
	for i, w := range pws {
		obs += k.m[i][i]
		exp += w * tws[i] / sum
	}
	return (obs - exp) / (sum - exp)
}

func (k *Kappa) grow(n int) {
	if d := n - len(k.m); d > 0 {
		k.m = append(k.m, make([][]float64, d)...)
	}
	if k.ncols < n {
		for i, r := range k.m {
			if d := n - len(r); d > 0 {
				k.m[i] = append(r, make([]float64, d)...)
			}
		}
		k.ncols = n
	}
}
