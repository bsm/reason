package testdata

import (
	"os"
	"path/filepath"

	"github.com/bsm/reason/ingest"

	"github.com/bsm/reason"
)

// BigClassificationModel definition
var BigClassificationModel = reason.NewModel(
	reason.NewCategoricalFeature("c1", []string{"v1", "v2", "v3", "v4", "v5"}),
	reason.NewCategoricalFeature("c2", []string{"v1", "v2", "v3", "v4", "v5"}),
	reason.NewCategoricalFeature("c3", []string{"v1", "v2", "v3", "v4", "v5"}),
	reason.NewCategoricalFeature("c4", []string{"v1", "v2", "v3", "v4", "v5"}),
	reason.NewCategoricalFeature("c5", []string{"v1", "v2", "v3", "v4", "v5"}),
	reason.NewNumericalFeature("n1"),
	reason.NewNumericalFeature("n2"),
	reason.NewNumericalFeature("n3"),
	reason.NewNumericalFeature("n4"),
	reason.NewNumericalFeature("n5"),
	reason.NewCategoricalFeature("target", []string{"c1", "c2"}),
)

// BigRegressionModel definition
var BigRegressionModel = reason.NewModel(
	reason.NewCategoricalFeature("c1", []string{"v1", "v2", "v3", "v4", "v5"}),
	reason.NewCategoricalFeature("c2", []string{"v1", "v2", "v3", "v4", "v5"}),
	reason.NewCategoricalFeature("c3", []string{"v1", "v2", "v3", "v4", "v5"}),
	reason.NewCategoricalFeature("c4", []string{"v1", "v2", "v3", "v4", "v5"}),
	reason.NewCategoricalFeature("c5", []string{"v1", "v2", "v3", "v4", "v5"}),
	reason.NewNumericalFeature("n1"),
	reason.NewNumericalFeature("n2"),
	reason.NewNumericalFeature("n3"),
	reason.NewNumericalFeature("n4"),
	reason.NewNumericalFeature("n5"),
	reason.NewNumericalFeature("target"),
)

// BigDataStream is a stream of test events.
type BigDataStream struct {
	*ingest.CSVStream

	file  *os.File
	model *reason.Model
}

// OpenBigData opens a bigdata stream.
func OpenBigData(kind, root string) (*BigDataStream, error) {
	switch kind {
	case "classification":
		return bigDataClassification(root)
	case "regression":
		return bigDataRegression(root)
	}
	panic("no such kind: " + kind)
}

func bigDataClassification(root string) (*BigDataStream, error) {
	return open(root, BigClassificationModel, []string{
		"c1", "c2", "c3", "c4", "c5", "n1", "n2", "n3", "n4", "n5", "target", "",
	})
}

func bigDataRegression(root string) (*BigDataStream, error) {
	return open(root, BigRegressionModel, []string{
		"c1", "c2", "c3", "c4", "c5", "n1", "n2", "n3", "n4", "n5", "", "target",
	})
}

func open(root string, model *reason.Model, headers []string) (*BigDataStream, error) {
	f, err := os.Open(filepath.Join(root, "bigdata.csv"))
	if err != nil {
		return nil, err
	}

	s, err := ingest.NewCSVStream(f, model, &ingest.CSVOptions{
		Headers:  headers,
		SkipRows: 1,
	})
	if err != nil {
		_ = f.Close()
		return nil, err
	}

	return &BigDataStream{
		CSVStream: s,
		file:      f,
	}, nil
}

// ReadN reads a batch of examples.
func (s *BigDataStream) ReadN(n int) ([]reason.Example, error) {
	batch := make([]reason.Example, 0, n)
	for i := 0; i < n; i++ {
		x, err := s.Next()
		if err != nil {
			return nil, err
		}
		batch = append(batch, x)
	}
	return batch, nil
}

// Close closes the stream.
func (s *BigDataStream) Close() error {
	_ = s.CSVStream.Close()
	return s.file.Close()
}
