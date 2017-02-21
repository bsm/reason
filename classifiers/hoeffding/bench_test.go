package hoeffding

import (
	"sync/atomic"
	"testing"

	"github.com/bsm/reason/core"
	"github.com/bsm/reason/testdata"
)

func BenchmarkTree_Train_c(b *testing.B) {
	benchmarkC(b, benchmarkTrain, func(t *Tree, s []core.Instance) {
		max := len(s)
		for i := 0; i < b.N; i++ {
			t.Train(s[i%max])
		}
	})
}

func BenchmarkTree_Train_r(b *testing.B) {
	benchmarkR(b, benchmarkTrain, func(t *Tree, s []core.Instance) {
		max := len(s)
		for i := 0; i < b.N; i++ {
			t.Train(s[i%max])
		}
	})
}

func BenchmarkTree_Predict_c(b *testing.B) {
	benchmarkC(b, benchmarkPredict, func(t *Tree, s []core.Instance) {
		max := len(s)
		for i := 0; i < b.N; i++ {
			p := t.Predict(s[i%max])
			p.Release()
		}
	})
}

func BenchmarkTree_Predict_r(b *testing.B) {
	benchmarkR(b, benchmarkPredict, func(t *Tree, s []core.Instance) {
		max := len(s)
		for i := 0; i < b.N; i++ {
			p := t.Predict(s[i%max])
			p.Release()
		}
	})
}

func BenchmarkTree_Train_cpb(b *testing.B) {
	benchmarkC(b, benchmarkTrain, func(t *Tree, s []core.Instance) {
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

func BenchmarkTree_Train_rpb(b *testing.B) {
	benchmarkR(b, benchmarkTrain, func(t *Tree, s []core.Instance) {
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

func benchmarkC(b *testing.B, bm benchmarkFunc, fn func(*Tree, []core.Instance)) {
	model := testdata.BigClassificationModel()
	stream, err := testdata.Open("../../testdata/bigcls.csv", model)
	if err != nil {
		b.Fatal(err)
	}
	defer stream.Close()
	bm(b, model, stream, fn)
}

func benchmarkR(b *testing.B, bm benchmarkFunc, fn func(*Tree, []core.Instance)) {
	model := testdata.BigRegressionModel()
	stream, err := testdata.Open("../../testdata/bigreg.csv", model)
	if err != nil {
		b.Fatal(err)
	}
	defer stream.Close()
	bm(b, model, stream, fn)
}

type benchmarkFunc func(*testing.B, *core.Model, *testdata.BigDataStream, func(*Tree, []core.Instance))

func benchmarkTrain(b *testing.B, model *core.Model, stream *testdata.BigDataStream, fn func(*Tree, []core.Instance)) {
	sample, err := stream.ReadN(1000)
	if err != nil {
		b.Fatal(err)
	}

	tree := New(model, nil)
	b.ResetTimer()
	fn(tree, sample)
}

func benchmarkPredict(b *testing.B, model *core.Model, stream *testdata.BigDataStream, fn func(*Tree, []core.Instance)) {
	sample, err := stream.ReadN(50000)
	if err != nil {
		b.Fatal(err)
	}

	tree := New(model, nil)
	for _, inst := range sample {
		tree.Train(inst)
	}

	b.ResetTimer()
	fn(tree, sample)
}
