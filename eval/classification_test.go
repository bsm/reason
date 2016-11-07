package eval

import (
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Classification", func() {
	var subject *Classification
	model := testdata.ClassificationModel()

	BeforeEach(func() {
		yes := core.Prediction{{Value: 0, Votes: 1}}
		no := core.Prediction{{Value: 1, Votes: 1}}

		subject = NewClassification(model)
		subject.Record(core.MapInstance{"play": "yes"}, yes)
		subject.Record(core.MapInstance{"play": "yes"}, yes)
		subject.Record(core.MapInstance{"play": "no"}, yes)
		subject.Record(core.MapInstance{"play": "no"}, no)
		subject.Record(core.MapInstance{"play": "no"}, no)
		subject.Record(core.MapInstance{"play": "yes"}, no)
		subject.Record(core.MapInstance{"play": "yes"}, yes)
		subject.Record(core.MapInstance{"play": "yes"}, yes)
		subject.Record(core.MapInstance{"play": "yes", "@weight": 2.0}, yes)
		subject.Record(core.MapInstance{"play": "no"}, no)
		subject.Record(core.MapInstance{"play": "yes"}, no)
	})

	It("should calculate stats", func() {
		Expect(subject.Kappa()).To(BeNumerically("~", 0.471, 0.001))
		Expect(subject.Correct()).To(Equal(0.75))
		Expect(subject.TotalWeight()).To(Equal(12.0))
	})

})
