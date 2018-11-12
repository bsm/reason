package reason_test

import (
	"math"

	"github.com/bsm/reason"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Feature", func() {

	DescribeTable("Number",
		func(f *reason.Feature, x reason.Example, exp float64) {
			if math.IsNaN(exp) {
				Expect(reason.IsNum(f.Number(x))).To(BeFalse())
			} else {
				Expect(f.Number(x)).To(Equal(exp))
			}
		},

		Entry("categorical",
			reason.NewCategoricalFeature("cat", nil), reason.MapExample{}, math.NaN()),
		Entry("numerical (no value)",
			reason.NewNumericalFeature("num"), reason.MapExample{}, math.NaN()),
		Entry("numerical (from float)",
			reason.NewNumericalFeature("num"), reason.MapExample{"num": 1.23}, 1.23),
		Entry("numerical (from int)",
			reason.NewNumericalFeature("num"), reason.MapExample{"num": 101}, 101.0),
		Entry("numerical (from string)",
			reason.NewNumericalFeature("num"), reason.MapExample{"num": "12.3"}, 12.3),
	)

	DescribeTable("Category",
		func(f *reason.Feature, x reason.Example, exp reason.Category) {
			Expect(f.Category(x)).To(Equal(exp))
		},

		Entry("numerical",
			reason.NewNumericalFeature("num"), reason.MapExample{}, reason.NoCategory),

		Entry("categorical (no value)",
			reason.NewCategoricalFeature("cat", []string{"x"}), reason.MapExample{}, reason.NoCategory),
		Entry("categorical (unknown value)",
			reason.NewCategoricalFeature("cat", []string{"x"}), reason.MapExample{"cat": "y"}, reason.NoCategory),
		Entry("categorical (string value)",
			reason.NewCategoricalFeature("cat", []string{"x"}), reason.MapExample{"cat": "x"}, reason.Category(0)),
		Entry("categorical (byte value)",
			reason.NewCategoricalFeature("cat", []string{"x"}), reason.MapExample{"cat": []byte("x")}, reason.Category(0)),
		Entry("categorical (int value)",
			reason.NewCategoricalFeature("cat", []string{"2"}), reason.MapExample{"cat": 2}, reason.Category(0)),

		Entry("categorical, no vocabulary, no hash buckets",
			reason.NewCategoricalFeature("cat", nil),
			reason.MapExample{"cat": "x"}, reason.NoCategory),
		Entry("categorical, vocabulary, no hash buckets (known value)",
			reason.NewCategoricalFeature("cat", []string{"a", "b"}),
			reason.MapExample{"cat": "b"}, reason.Category(1)),
		Entry("categorical, vocabulary, no hash buckets (unknown value)",
			reason.NewCategoricalFeature("cat", []string{"a"}),
			reason.MapExample{"cat": "x"}, reason.NoCategory),

		Entry("categorical, hash buckets",
			reason.NewCategoricalFeatureHashBuckets("cat", 10),
			reason.MapExample{"cat": "z"}, reason.Category(6)),

		Entry("categorical, vocabulary, hash buckets (known value)",
			reason.NewCategoricalFeatureVocabularyHashBuckets("cat", []string{"a", "b"}, 10),
			reason.MapExample{"cat": "a"}, reason.Category(0)),
		Entry("categorical, vocabulary, hash buckets (unknown value)",
			reason.NewCategoricalFeatureVocabularyHashBuckets("cat", []string{"a", "b"}, 10),
			reason.MapExample{"cat": "z"}, reason.Category(8)),

		Entry("categorical, vocabulary, expandable (known value)",
			reason.NewCategoricalFeatureExpandable("cat", []string{"a", "b"}),
			reason.MapExample{"cat": "b"}, reason.Category(1)),
		Entry("categorical, vocabulary, expandable (unknown value)",
			reason.NewCategoricalFeatureExpandable("cat", []string{"a", "b"}),
			reason.MapExample{"cat": "x"}, reason.Category(2)),

		Entry("categorical, identity (no value)",
			reason.NewCategoricalFeatureIdentity("cat"), reason.MapExample{}, reason.NoCategory),
		Entry("categorical, identity (bad value)",
			reason.NewCategoricalFeatureIdentity("cat"), reason.MapExample{"cat": "x"}, reason.NoCategory),
		Entry("categorical, identity (string value)",
			reason.NewCategoricalFeatureIdentity("cat"), reason.MapExample{"cat": "2"}, reason.Category(2)),
		Entry("categorical, identity (byte value)",
			reason.NewCategoricalFeatureIdentity("cat"), reason.MapExample{"cat": []byte("1")}, reason.Category(1)),
		Entry("categorical, identity (int value)",
			reason.NewCategoricalFeatureIdentity("cat"), reason.MapExample{"cat": 5}, reason.Category(5)),
		Entry("categorical, identity (bool value)",
			reason.NewCategoricalFeatureIdentity("cat"), reason.MapExample{"cat": true}, reason.Category(1)),
		Entry("categorical, identity (numeric value)",
			reason.NewCategoricalFeatureIdentity("cat"), reason.MapExample{"cat": customNumeric(8)}, reason.Category(8)),
		Entry("categorical, identity (stringish value)",
			reason.NewCategoricalFeatureIdentity("cat"), reason.MapExample{"cat": customString("7")}, reason.Category(7)),
		Entry("categorical, identity (int ptr)",
			reason.NewCategoricalFeatureIdentity("cat"), reason.MapExample{"cat": intPtr(6)}, reason.Category(6)),
	)

})
