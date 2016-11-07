package hoeffding

import (
	"sync/atomic"
	"testing"

	"github.com/bsm/reason/core"
	"github.com/bsm/reason/testdata"
)

func BenchmarkTreeTrain_C(b *testing.B) {
	benchmarkTrainC(b, func(t *Tree, s []core.Instance) {
		max := len(s)
		for i := 0; i < b.N; i++ {
			t.Train(s[i%max])
		}
	})
}

func BenchmarkTreeTrain_R(b *testing.B) {
	benchmarkTrainR(b, func(t *Tree, s []core.Instance) {
		max := len(s)
		for i := 0; i < b.N; i++ {
			t.Train(s[i%max])
		}
	})
}

func BenchmarkTreeTrain_CParallel(b *testing.B) {
	benchmarkTrainC(b, func(t *Tree, s []core.Instance) {
		max := len(s)
		cnt := int64(0)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				i := int(atomic.AddInt64(&cnt, 1))
				t.Train(s[i%max])
			}
		})
	})
}

func BenchmarkTreeTrain_RParallel(b *testing.B) {
	benchmarkTrainR(b, func(t *Tree, s []core.Instance) {
		max := len(s)
		cnt := int64(0)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				i := int(atomic.AddInt64(&cnt, 1))
				t.Train(s[i%max])
			}
		})
	})
}

func benchmarkTrainC(b *testing.B, cb func(*Tree, []core.Instance)) {
	model := testdata.BigClassificationModel()
	stream, err := testdata.Open("../../testdata/bigcls.csv", model)
	if err != nil {
		b.Fatal(err)
	}
	defer stream.Close()
	benchmarkTrain(b, model, stream, cb)
}

func benchmarkTrainR(b *testing.B, cb func(*Tree, []core.Instance)) {
	model := testdata.BigRegressionModel()
	stream, err := testdata.Open("../../testdata/bigreg.csv", model)
	if err != nil {
		b.Fatal(err)
	}
	defer stream.Close()
	benchmarkTrain(b, model, stream, cb)
}

func benchmarkTrain(b *testing.B, model *core.Model, stream *testdata.BigDataStream, cb func(*Tree, []core.Instance)) {
	sample, err := stream.ReadN(1000)
	if err != nil {
		b.Fatal(err)
	}

	tree := New(model, nil)
	b.ResetTimer()
	cb(tree, sample)
}
