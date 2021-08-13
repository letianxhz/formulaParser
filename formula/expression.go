package formula

import "fmt"

var precedence = map[string]int{"+": 2, "-": 2, "*": 4, "/": 4}

type Expression interface {
	Eval() float64
	ToString() string
}

type BinaryExp struct {
	Op string
	Lhs Expression
	Rhs Expression
}

func (c *BinaryExp) Eval() float64 {
	r := c.Rhs.Eval()
	l := c.Lhs.Eval()
	switch c.Op {
	case "+":
		return l + r
	case "-":
		return l - r
	case "*":
		return l * r
	case "/":
		return l / r
	default:
		return 0
	}
}

func (c *BinaryExp) ToString() string  {
	return fmt.Sprintf(
		"BinaryExprAST: (OP:%s Lhs:%s Rhs:%s)",
		c.Op,
		c.Lhs.ToString(),
		c.Rhs.ToString(),
	)
}


type ConstantExp struct {
	Val float64
}

func (c *ConstantExp) ToString() string  {
	return fmt.Sprintf(
		"ConstantExprAST:%f",
		c.Val,
	)
}

func (c *ConstantExp) Eval() float64  {
	return c.Val
}

type ParamExp struct {
	Val float64
	Name string
}

func (p *ParamExp) Eval() float64  {
	if _, ok := ParamContext[p.Name]; !ok {
		return 0
	}
	return ParamContext[p.Name].(float64)
}

func (p *ParamExp) ToString() string {
	return fmt.Sprintf(
		"ParamExpAST:%f name:%s",
		p.Val,
		p.Name,
	)
}

type FuncExp struct {
	Name string
	Arg  []Expression
}

func (f *FuncExp) Eval() float64  {
	if _, ok := funcMap[f.Name]; !ok {
		return 0
	}
	def := funcMap[f.Name]
	return def.function(f.Arg...)
}

func (f *FuncExp) ToString() string  {
	var s string
	for _, a := range f.Arg{
		s += a.ToString() +" "
	}
	return fmt.Sprintf(
		"FuncExpAST:%s name:%+v",
		f.Name,
		s,
	)
}