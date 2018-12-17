package hoeffding_test

import (
	"testing"

	"github.com/bsm/reason/classifier/hoeffding"
	"github.com/bsm/reason/testdata"
)

func BenchmarkTree_Train_classification(b *testing.B) {
	const N = 1000

	stream, err := testdata.OpenBigData("classification", "../../testdata")
	if err != nil {
		b.Fatal(err)
	}
	defer stream.Close()

	examples, err := stream.ReadN(N)
	if err != nil {
		b.Fatal(err)
	}

	tree, err := hoeffding.New(stream.Model(), "target", nil)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		tree.Train(examples[i%N])
	}
}

func BenchmarkTree_Train_classification_parallel(b *testing.B) {
	const N = 1000

	stream, err := testdata.OpenBigData("classification", "../../testdata")
	if err != nil {
		b.Fatal(err)
	}
	defer stream.Close()

	examples, err := stream.ReadN(N)
	if err != nil {
		b.Fatal(err)
	}

	tree, err := hoeffding.New(stream.Model(), "target", nil)
	if err != nil {
		b.Fatal(err)
	}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			tree.Train(examples[i%N])
		}
	})
}

func BenchmarkTree_Train_regression(b *testing.B) {
	const N = 1000

	stream, err := testdata.OpenBigData("regression", "../../testdata")
	if err != nil {
		b.Fatal(err)
	}
	defer stream.Close()

	examples, err := stream.ReadN(N)
	if err != nil {
		b.Fatal(err)
	}

	tree, err := hoeffding.New(stream.Model(), "target", nil)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		tree.Train(examples[i%N])
	}
}

func BenchmarkTree_Train_regression_parallel(b *testing.B) {
	const N = 1000

	stream, err := testdata.OpenBigData("regression", "../../testdata")
	if err != nil {
		b.Fatal(err)
	}
	defer stream.Close()

	examples, err := stream.ReadN(N)
	if err != nil {
		b.Fatal(err)
	}

	tree, err := hoeffding.New(stream.Model(), "target", nil)
	if err != nil {
		b.Fatal(err)
	}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			tree.Train(examples[i%N])
		}
	})
}
