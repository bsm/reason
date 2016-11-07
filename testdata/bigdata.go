package testdata

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"github.com/bsm/reason/core"
)

func BigClassificationModel() *core.Model {
	return core.NewModel(
		&core.Attribute{Name: "tv", Kind: core.AttributeKindNominal, Values: core.NewAttributeValues("c1", "c2")},
		&core.Attribute{Name: "c1", Kind: core.AttributeKindNominal, Values: core.NewAttributeValues("v1", "v2", "v3", "v4", "v5")},
		&core.Attribute{Name: "c2", Kind: core.AttributeKindNominal, Values: core.NewAttributeValues("v1", "v2", "v3", "v4", "v5")},
		&core.Attribute{Name: "c3", Kind: core.AttributeKindNominal, Values: core.NewAttributeValues("v1", "v2", "v3", "v4", "v5")},
		&core.Attribute{Name: "c4", Kind: core.AttributeKindNominal, Values: core.NewAttributeValues("v1", "v2", "v3", "v4", "v5")},
		&core.Attribute{Name: "c5", Kind: core.AttributeKindNominal, Values: core.NewAttributeValues("v1", "v2", "v3", "v4", "v5")},
		&core.Attribute{Name: "n1", Kind: core.AttributeKindNumeric},
		&core.Attribute{Name: "n2", Kind: core.AttributeKindNumeric},
		&core.Attribute{Name: "n3", Kind: core.AttributeKindNumeric},
		&core.Attribute{Name: "n4", Kind: core.AttributeKindNumeric},
		&core.Attribute{Name: "n5", Kind: core.AttributeKindNumeric},
	)
}

func BigRegressionModel() *core.Model {
	return core.NewModel(
		&core.Attribute{Name: "tv", Kind: core.AttributeKindNumeric},
		&core.Attribute{Name: "c1", Kind: core.AttributeKindNominal},
		&core.Attribute{Name: "c2", Kind: core.AttributeKindNominal},
		&core.Attribute{Name: "c3", Kind: core.AttributeKindNominal},
		&core.Attribute{Name: "c4", Kind: core.AttributeKindNominal},
		&core.Attribute{Name: "n1", Kind: core.AttributeKindNumeric},
	)
}

type BigDataStream struct {
	file  *os.File
	recs  *csv.Reader
	model *core.Model

	inst core.MapInstance
	err  error
}

func Open(fname string, model *core.Model) (*BigDataStream, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	recs := csv.NewReader(f)
	recs.FieldsPerRecord = model.NumPredictors() + 1

	return &BigDataStream{
		file:  f,
		recs:  recs,
		model: model,
		inst:  make(core.MapInstance, recs.FieldsPerRecord),
	}, nil
}

func (s *BigDataStream) Next() bool {
	if s.err != nil {
		return false
	}

	fields, err := s.recs.Read()
	if err != nil {
		s.err = err
		return false
	}

	// init instance
	s.inst = make(core.MapInstance, len(s.inst))

	// read predictors
	for i, attr := range s.model.Predictors() {
		str := fields[i]
		if attr.IsNominal() {
			s.inst[attr.Name] = str
		} else if s.inst[attr.Name], err = strconv.ParseFloat(str, 64); err != nil {
			s.err = err
			return false
		}
	}

	// read target
	str := fields[len(fields)-1]
	if target := s.model.Target(); target.IsNominal() {
		s.inst[target.Name] = str
	} else if s.inst[target.Name], err = strconv.ParseFloat(str, 64); err != nil {
		s.err = err
		return false
	}

	return true
}

func (s *BigDataStream) ReadN(n int) ([]core.Instance, error) {
	res := make([]core.Instance, 0, n)
	for s.Next() {
		res = append(res, s.Instance())
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

func (s *BigDataStream) Instance() core.Instance { return s.inst }
func (s *BigDataStream) Close() error            { return s.file.Close() }
