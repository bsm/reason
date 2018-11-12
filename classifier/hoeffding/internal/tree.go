package internal

import (
	"bufio"
	fmt "fmt"
	"io"

	"github.com/bsm/reason"
	"github.com/bsm/reason/internal/iocount"
	"github.com/bsm/reason/internal/protoio"
	"github.com/gogo/protobuf/proto"
)

// NewTree inits a brand-new tree
func NewTree(model *reason.Model, target string) *Tree {
	t := &Tree{
		Model:  model,
		Target: target,
	}
	t.Root = t.AddLeaf(nil, 0)
	return t
}

// AddLeaf adds a new node to registry and returns the reference.
func (t *Tree) AddLeaf(stats isNode_Stats, weightAtLastEval float64) int64 {
	leaf := &LeafNode{WeightAtLastEval: weightAtLastEval}
	node := &Node{Kind: &Node_Leaf{Leaf: leaf}, Stats: stats}
	t.Nodes = append(t.Nodes, node)
	return int64(len(t.Nodes))
}

// ReplaceNode replaces a node by reference.
func (t *Tree) ReplaceNode(nref int64, node *Node) {
	if pos := int(nref - 1); pos > -1 && pos < len(t.Nodes) {
		t.Nodes[pos] = node
	}
}

// NumNodes returns the number of nodes.
func (t *Tree) NumNodes() int {
	return len(t.Nodes)
}

// GetNode retrieves a node by ref from registry.
func (t *Tree) GetNode(nref int64) *Node {
	if pos := int(nref - 1); pos > -1 && pos < len(t.Nodes) {
		return t.Nodes[pos]
	}
	return nil
}

// GetLeaf retrieves a leaf node by ref from registry.
func (t *Tree) GetLeaf(nref int64) *LeafNode {
	if node := t.GetNode(nref); node != nil {
		return node.GetLeaf()
	}
	return nil
}

// GetSplit retrieves a split node by ref from registry.
func (t *Tree) GetSplit(nref int64) *SplitNode {
	if node := t.GetNode(nref); node != nil {
		return node.GetSplit()
	}
	return nil
}

// SplitNode splits the nref node by feature name.
// Returns true if successful.
func (t *Tree) SplitNode(nref int64, feature string, postSplit PostSplit, pivot float64) bool {
	node := t.GetNode(nref)
	if node == nil {
		return false
	}

	leaf := node.GetLeaf()
	if leaf == nil {
		return false
	}

	split := &SplitNode{
		Feature: feature,
		Pivot:   pivot,
	}
	postSplit.forEach(func(pos int, stats isNode_Stats, weight float64) {
		split.SetChild(pos, t.AddLeaf(stats, weight))
	})
	node.Kind = &Node_Split{Split: split}
	return true
}

// Traverse traverses the tree with an example x starting at a given nref.
func (t *Tree) Traverse(x reason.Example, nref int64, parent *Node, ppos int) (*Node, int64, *Node, int) {
	node := t.GetNode(nref)
	if node == nil {
		return node, nref, parent, ppos
	}

	if split := node.GetSplit(); split != nil {
		feature := t.Model.Feature(split.Feature)

		if npos := int(split.childPos(feature, x)); npos > -1 {
			if cnref := split.GetChild(npos); cnref > 0 {
				return t.Traverse(x, cnref, node, npos)
			}
			return nil, nref, node, npos
		}
	}
	return node, nref, parent, ppos
}

// ReadFrom reads a tree from a Reader.
func (t *Tree) ReadFrom(r io.Reader) (int64, error) {
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
			t.Model = model
		case 2: // target
			if wire != proto.WireBytes {
				return rc.N, proto.ErrInternalBadWireType
			}

			str, err := rp.ReadString()
			if err != nil {
				return rc.N, err
			}
			t.Target = str
		case 3: // root
			if wire != proto.WireVarint {
				return rc.N, proto.ErrInternalBadWireType
			}

			u, err := rp.ReadVarint()
			if err != nil {
				return rc.N, err
			}
			t.Root = int64(u)
		case 4: // nodes
			if wire != proto.WireBytes {
				return rc.N, proto.ErrInternalBadWireType
			}

			node := new(Node)
			if err := rp.ReadMessage(node); err != nil {
				return rc.N, err
			}
			t.Nodes = append(t.Nodes, node)
		default:
			return rc.N, fmt.Errorf("hoeffding: unexpected field tag %d", tag)
		}
	}
}

// WriteTo writes a tree to a Writer.
func (t *Tree) WriteTo(w io.Writer) (int64, error) {
	wc := &iocount.Writer{W: w}
	wp := &protoio.Writer{Writer: bufio.NewWriter(wc)}

	if err := wp.WriteMessageField(1, t.Model); err != nil {
		return wc.N, err
	}
	if err := wp.WriteStringField(2, t.Target); err != nil {
		return wc.N, err
	}
	if err := wp.WriteVarintField(3, uint64(t.Root)); err != nil {
		return wc.N, err
	}
	for _, node := range t.Nodes {
		if err := wp.WriteMessageField(4, node); err != nil {
			return wc.N, err
		}
	}
	return wc.N, wp.Flush()
}

// WriteText appends node information to a text document.
func (t *Tree) WriteText(w io.Writer, nref int64, indent, name string) (written int64, err error) {
	// Get the node
	node := t.GetNode(nref)
	if node == nil {
		return
	}

	var (
		nn int
		sz int64
	)

	// Print node stats
	nn, err = fmt.Fprintf(w, indent+name+" [weight:%.0f]\n", node.Weight())
	written += int64(nn)
	if err != nil {
		return
	}

	// Recurse if a split node
	if split := node.GetSplit(); split != nil {
		feat := t.Model.Feature(split.Feature)
		if feat == nil {
			return
		}

		subIndent := indent + "\t"
		for i, childRef := range split.Children {
			if childRef == 0 {
				continue
			}

			sz, err = t.WriteText(w, childRef, subIndent, nodeCondition(feat, i, split.Pivot))
			written += sz
			if err != nil {
				return
			}
		}
	}
	return
}

// WriteDOT appends node information to a DOT document.
func (t *Tree) WriteDOT(w io.Writer, nref int64, name, label string) (written int64, err error) {
	// Get the node
	node := t.GetNode(nref)
	if node == nil {
		return
	}

	var (
		nn int
		sz int64
	)

	// Write node data
	nn, err = fmt.Fprintf(w, `  %s [label="%sweight: %.0f"];`+"\n", name, label, node.Weight())
	written += int64(nn)
	if err != nil {
		return
	}

	// Recurse if a split node
	if split := node.GetSplit(); split != nil {
		feat := t.Model.Feature(split.Feature)
		if feat == nil {
			return
		}

		for i, childRef := range split.Children {
			if childRef == 0 {
				continue
			}

			subName := fmt.Sprintf("%s_%d", name, i)
			nn, err = fmt.Fprintf(w, "  %s -> %s;\n", name, subName)
			written += int64(nn)
			if err != nil {
				return
			}

			sz, err = t.WriteDOT(w, childRef, subName, nodeCondition(feat, i, split.Pivot)+`\n`)
			written += sz
			if err != nil {
				return
			}
		}
	}
	return
}

// nodeCondition
func nodeCondition(feat *reason.Feature, pos int, pivot float64) string {
	if feat.Kind.IsNumerical() {
		if pos == 0 {
			return fmt.Sprintf(`%s <= %.2f`, feat.Name, pivot)
		} else {
			return fmt.Sprintf(`%s > %.2f`, feat.Name, pivot)
		}
	}

	cat := reason.Category(pos)
	return fmt.Sprintf(`%s = %s`, feat.Name, feat.ValueOf(cat))
}
