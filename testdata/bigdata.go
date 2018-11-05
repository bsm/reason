package testdata

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bsm/reason/core"
)

var BigClassificationModel = core.NewModel(
	core.NewCategoricalFeature("c1", []string{"v1", "v2", "v3", "v4", "v5"}),
	core.NewCategoricalFeature("c2", []string{"v1", "v2", "v3", "v4", "v5"}),
	core.NewCategoricalFeature("c3", []string{"v1", "v2", "v3", "v4", "v5"}),
	core.NewCategoricalFeature("c4", []string{"v1", "v2", "v3", "v4", "v5"}),
	core.NewCategoricalFeature("c5", []string{"v1", "v2", "v3", "v4", "v5"}),
	core.NewNumericalFeature("n1"),
	core.NewNumericalFeature("n2"),
	core.NewNumericalFeature("n3"),
	core.NewNumericalFeature("n4"),
	core.NewNumericalFeature("n5"),
	core.NewCategoricalFeature("target", []string{"c1", "c2"}),
)

var BigRegressionModel = core.NewModel(
	core.NewCategoricalFeature("c1", []string{"v1", "v2", "v3", "v4", "v5"}),
	core.NewCategoricalFeature("c2", []string{"v1", "v2", "v3", "v4", "v5"}),
	core.NewCategoricalFeature("c3", []string{"v1", "v2", "v3", "v4", "v5"}),
	core.NewCategoricalFeature("c4", []string{"v1", "v2", "v3", "v4", "v5"}),
	core.NewCategoricalFeature("c5", []string{"v1", "v2", "v3", "v4", "v5"}),
	core.NewNumericalFeature("n1"),
	core.NewNumericalFeature("n2"),
	core.NewNumericalFeature("n3"),
	core.NewNumericalFeature("n4"),
	core.NewNumericalFeature("n5"),
	core.NewNumericalFeature("target"),
)

// BigDataStream is a stream of test events.
type BigDataStream struct {
	file    *os.File
	recs    *csv.Reader
	model   *core.Model
	mapping map[string]int

	x   core.MapExample
	err error
}

// OpenBigData opens a bigdata stream.
func OpenBigData(kind, root string) (*BigDataStream, *core.Model, error) {
	switch kind {
	case "classification":
		return bigDataClassification(root)
	case "regression":
		return bigDataRegression(root)
	}
	panic("no such kind: " + kind)
}

func bigDataClassification(root string) (*BigDataStream, *core.Model, error) {
	return open(root, BigClassificationModel, map[string]int{
		"c1": 0, "c2": 1, "c3": 2, "c4": 3, "c5": 4, "n1": 5, "n2": 6, "n3": 7, "n4": 8, "n5": 9, "target": 10,
	})
}

func bigDataRegression(root string) (*BigDataStream, *core.Model, error) {
	return open(root, BigRegressionModel, map[string]int{
		"c1": 0, "c2": 1, "c3": 2, "c4": 3, "c5": 4, "n1": 5, "n2": 6, "n3": 7, "n4": 8, "n5": 9, "target": 11,
	})
}

func open(root string, model *core.Model, mapping map[string]int) (*BigDataStream, *core.Model, error) {
	f, err := os.Open(filepath.Join(root, "bigdata.csv"))
	if err != nil {
		return nil, nil, err
	}

	recs := csv.NewReader(f)
	recs.FieldsPerRecord = len(model.Features) + 1

	return &BigDataStream{
		file:    f,
		recs:    recs,
		model:   model,
		mapping: mapping,
	}, model, nil
}

func (s *BigDataStream) Next() bool {
	if s.err != nil {
		return false
	}

	row, err := s.recs.Read()
	if err != nil {
		s.err = err
		return false
	}

	// read predictors
	s.x = make(core.MapExample, s.recs.FieldsPerRecord)
	for name, feat := range s.model.Features {
		if str := row[s.mapping[name]]; str == "" {
			continue
		} else if feat.Kind.IsCategorical() {
			s.x[name] = str
		} else if feat.Kind.IsNumerical() {
			if s.x[name], err = strconv.ParseFloat(str, 64); err != nil {
				s.err = err
				return false
			}
		}
	}
	return true
}

func (s *BigDataStream) ReadN(n int) ([]core.Example, error) {
	res := make([]core.Example, 0, n)
	for s.Next() {
		res = append(res, s.Example())
		if len(res) == n {
			break
		}
	}
	return res, s.Err()
}

func (s *BigDataStream) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}

func (s *BigDataStream) Example() core.Example { return s.x }
func (s *BigDataStream) Close() error          { return s.file.Close() }
