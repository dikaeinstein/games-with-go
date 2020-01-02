package apt

import "math"

// Node describes a node of the Abstract Picture Tree(APT)
type Node interface {
	Eval(x, y float32) float32
	String() string
}

// Leaf is a leaf node
type Leaf struct{}

// Single is a node with a one child node
type Single struct {
	Child Node
}

// Double is a node with two children nodes
type Double struct {
	LeftChild  Node
	RightChild Node
}

// OpPlus is the plus operator node
type OpPlus struct {
	Double
}

// Eval evaluates the plus operation on the operands
func (op *OpPlus) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) + op.RightChild.Eval(x, y)
}

func (op *OpPlus) String() string {
	return "( + " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

// OpX is the X constant node
type OpX Leaf

// Eval evaluates the value of x
func (OpX) Eval(x, y float32) float32 {
	return x
}

func (OpX) String() string {
	return "X"
}

// OpY is the Y constant node
type OpY Leaf

// Eval evaluates the value of y
func (OpY) Eval(x, y float32) float32 {
	return y
}

func (OpY) String() string {
	return "Y"
}

// OpSin is the Sine operation node
type OpSin struct {
	Single
}

// Eval evaluates the sin(x)
func (ops *OpSin) Eval(x, y float32) float32 {
	return float32(math.Sin(float64(ops.Child.Eval(x, y))))
}

func (ops *OpSin) String() string {
	return "( Sin " + ops.Child.String() + " )"
}
