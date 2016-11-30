package eval

import (
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Regression", func() {
	var subject *Regression
	model := testdata.RegressionModel()

	BeforeEach(func() {
		subject = NewRegression(model)
		subject.Record(core.MapInstance{"hours": 26}, core.Prediction{{AttributeValue: 25, Votes: 1}})
		subject.Record(core.MapInstance{"hours": 20}, core.Prediction{{AttributeValue: 25, Votes: 1}})
		subject.Record(core.MapInstance{"hours": 24}, core.Prediction{{AttributeValue: 22, Votes: 1}})
		subject.Record(core.MapInstance{"hours": 21}, core.Prediction{{AttributeValue: 23, Votes: 1}})
		subject.Record(core.MapInstance{"hours": 23}, core.Prediction{{AttributeValue: 24, Votes: 1}})
		subject.Record(core.MapInstance{"hours": 25}, core.Prediction{{AttributeValue: 29, Votes: 1}})
		subject.Record(core.MapInstance{"hours": 27}, core.Prediction{{AttributeValue: 28, Votes: 1}})
		subject.Record(core.MapInstance{"hours": 28, "@weight": 2.0}, core.Prediction{{AttributeValue: 26, Votes: 1}})
		subject.Record(core.MapInstance{"hours": 29}, core.Prediction{{AttributeValue: 30, Votes: 1}})
		subject.Record(core.MapInstance{"hours": 22}, core.Prediction{{AttributeValue: 18, Votes: 1}})
	})

	It("should calculate stats", func() {
		Expect(subject.TotalWeight()).To(Equal(11.0))
		Expect(subject.Mean()).To(BeNumerically("~", 24.8, 0.1))
		Expect(subject.MAE()).To(BeNumerically("~", 2.27, 0.01))
		Expect(subject.RMSE()).To(BeNumerically("~", 2.64, 0.01))
		Expect(subject.R2()).To(BeNumerically("~", 0.39, 0.01))

		subject.Record(core.MapInstance{"hours": 28, "@weight": 2.0}, core.Prediction{{AttributeValue: 28, Votes: 1}})
		Expect(subject.R2()).To(BeNumerically("~", 0.47, 0.01))
	})

})
