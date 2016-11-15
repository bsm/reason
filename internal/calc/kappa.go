package calc

type KappaStat struct {
	m     [][]float64
	ncols int
}

func NewKappaStat() *KappaStat { return new(KappaStat) }

func (k *KappaStat) Record(expIndex, actIndex int, weight float64) {
	if expIndex > actIndex {
		k.grow(expIndex + 1)
	} else {
		k.grow(actIndex + 1)
	}

	k.m[actIndex][expIndex] += weight
}

func (k *KappaStat) Kappa() float64 {
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

func (k *KappaStat) grow(n int) {
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
