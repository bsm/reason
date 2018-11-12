package internal

import (
	"bufio"
	"bytes"
	fmt "fmt"
	"io"

	"github.com/bsm/reason"
	"github.com/bsm/reason/common/observer"
	"github.com/bsm/reason/internal/iocount"
	"github.com/bsm/reason/internal/protoio"
	"github.com/gogo/protobuf/proto"
)

// New inits a new NaiveBayes classifier.
func New(model *reason.Model, target string) *NaiveBayes {
	return &NaiveBayes{Model: model, Target: target}
}

// ReadFrom reads from a Reader.
func (n *NaiveBayes) ReadFrom(r io.Reader) (int64, error) {
	rc := &iocount.Reader{R: r}
	rp := &protoio.Reader{Reader: bufio.NewReader(rc)}

	for {
		tag, wire, err := rp.ReadField()
		if err == io.EOF {
			return rc.N, nil
		} else if err != nil {
			return rc.N, err
		}

		switch tag {
		case 1: // model
			if wire != proto.WireBytes {
				return rc.N, proto.ErrInternalBadWireType
			}

			model := new(reason.Model)
			if err := rp.ReadMessage(model); err != nil {
				return rc.N, err
			}
			n.Model = model
		case 2: // target
			if wire != proto.WireBytes {
				return rc.N, proto.ErrInternalBadWireType
			}

			str, err := rp.ReadString()
			if err != nil {
				return rc.N, err
			}
			n.Target = str
		case 3: // target stats
			if wire != proto.WireBytes {
				return rc.N, proto.ErrInternalBadWireType
			}

			if err := rp.ReadMessage(&n.TargetStats); err != nil {
				return rc.N, err
			}
		case 4: // feature stats
			if wire != proto.WireBytes {
				return rc.N, proto.ErrInternalBadWireType
			}

			// read message size
			_, err := rp.ReadVarint()
			if err != nil {
				return rc.N, err
			}

			if err := n.readFeatureStats(rp); err != nil {
				return rc.N, err
			}
		default:
			return rc.N, fmt.Errorf("bayes: unexpected field tag %d", tag)
		}
	}
}

func (n *NaiveBayes) readFeatureStats(rp *protoio.Reader) error {
	if n.FeatureStats == nil {
		n.FeatureStats = make(map[string]*NaiveBayes_Observer)
	}

	key := ""
	for i := 0; i < 2; i++ {
		tag, wire, err := rp.ReadField()
		if err != nil {
			return err
		}
		if wire != proto.WireBytes {
			return proto.ErrInternalBadWireType
		}

		switch tag {
		case 1: // key
			if key, err = rp.ReadString(); err != nil {
				return err
			}
		case 2: // value
			obs := new(NaiveBayes_Observer)
			if err := rp.ReadMessage(obs); err != nil {
				return err
			}
			n.FeatureStats[key] = obs
		default:
			return fmt.Errorf("bayes: unexpected field tag %d", tag)
		}
	}
	return nil
}

// WriteTo writes to a Writer.
func (n *NaiveBayes) WriteTo(w io.Writer) (int64, error) {
	wc := &iocount.Writer{W: w}
	wp := &protoio.Writer{Writer: bufio.NewWriter(wc)}

	if err := wp.WriteMessageField(1, n.Model); err != nil {
		return wc.N, err
	}
	if err := wp.WriteStringField(2, n.Target); err != nil {
		return wc.N, err
	}
	if err := wp.WriteMessageField(3, &n.TargetStats); err != nil {
		return wc.N, err
	}
	for key, obs := range n.FeatureStats {
		buf := new(bytes.Buffer)
		sub := &protoio.Writer{Writer: bufio.NewWriter(buf)}
		if err := sub.WriteStringField(1, key); err != nil {
			return wc.N, err
		}
		if err := sub.WriteMessageField(2, obs); err != nil {
			return wc.N, err
		}
		if err := sub.Flush(); err != nil {
			return wc.N, err
		}
		if err := wp.WriteBinaryField(4, buf.Bytes()); err != nil {
			return wc.N, err
		}
	}
	return wc.N, wp.Flush()
}

// ObserveWeight observes example and updates stats.
func (n *NaiveBayes) ObserveWeight(x reason.Example, tcat reason.Category, weight float64) {
	n.TargetStats.Incr(int(tcat), weight)

	for name, feat := range n.Model.Features {
		if name != n.Target {
			n.observeFeat(x, feat, tcat, weight)
		}
	}
}

func (n *NaiveBayes) observeFeat(x reason.Example, feat *reason.Feature, tcat reason.Category, weight float64) {
	if n.FeatureStats == nil {
		n.FeatureStats = make(map[string]*NaiveBayes_Observer)
	}

	wrp, ok := n.FeatureStats[feat.Name]
	if !ok {
		wrp = &NaiveBayes_Observer{}
		n.FeatureStats[feat.Name] = wrp
	}

	switch feat.Kind {
	case reason.Feature_CATEGORICAL:
		obs := wrp.GetCat()
		if obs == nil {
			obs = observer.NewClassificationCategorical()
			wrp.Kind = &NaiveBayes_Observer_Cat{Cat: obs}
		}
		obs.ObserveWeight(feat.Category(x), tcat, weight)
	case reason.Feature_NUMERICAL:
		obs := wrp.GetNum()
		if obs == nil {
			obs = observer.NewClassificationNumerical(0)
			wrp.Kind = &NaiveBayes_Observer_Num{Num: obs}
		}
		obs.ObserveWeight(feat.Number(x), tcat, weight)
	}
}
