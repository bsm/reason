package util

type KappaStat NumMatrix

func (k KappaStat) Record(expected, actual int, weight float64) KappaStat {
	m := NumMatrix(k)
	row := NumVector(m.GetRow(actual))
	row = row.Incr(expected, weight)
	return KappaStat(m.SetRow(actual, row))
}

func (k KappaStat) Kappa() float64 {
	m := NumMatrix(k)
	rsums, tsum := m.SumRowsPlusTotal()
	csums := m.SumCols()

	var obs, exp float64

	for i, sum := range rsums {
		obs += m[i][i]
		exp += sum * csums[i] / tsum
	}
	return (obs - exp) / (tsum - exp)
}
