package formula

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

//公式二元表达式树缓存
var FormulaCache = make(map[string]Expression)

//上下文信息用于传公式依赖的参数
var EvaluationContext = make(map[string]interface{})

//定义的函数实现
type funcInfo struct {
	argsCount int
	function  func(expr ...Expression) float64
}

var funcMap = map[string]funcInfo{
	"max": {2, func(expr ...Expression) float64 {
		return math.Max(expr[0].Eval(), expr[1].Eval())
	}},

	"min": {2, func(expr ...Expression) float64 {
		return math.Min(expr[0].Eval(), expr[1].Eval())
	}},

	"ahasbuff": {1, func(expr ...Expression) float64 {
		buffId := strconv.FormatFloat(expr[0].Eval(), 'f', 0, 64)
		if _, ok := EvaluationContext["buffa"]; !ok {
			return 0
		}
		buffs := EvaluationContext["buffa"].(map[string]bool)
		if _, ok := buffs[buffId]; ok {
			return 1
		}
		return 0
	}},

	"dhasbuff": {1, func(expr ...Expression) float64 {
		buffId := strconv.FormatFloat(expr[0].Eval(), 'f', 0, 64)
		if _, ok := EvaluationContext["buffb"]; !ok {
			return 0
		}
		buffs := EvaluationContext["buffb"].(map[string]bool)
		if _, ok := buffs[buffId]; ok {
			return 1
		}
		return 0
	}},
}

func SetContext(ctx map[string]interface{})  {
	EvaluationContext = ctx
}

func AddAstExpCache(exp string, e Expression)  {
	if _, ok := FormulaCache[exp]; !ok {
		FormulaCache[exp] = e
	}
}

func Exec(exp string) {
	tokens, err := Parse(exp)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}
	ast := NewAST(tokens, exp)
	if ast.Err != nil {
		fmt.Println("ERROR: " + ast.Err.Error())
		return
	}
	ae := ast.ParseExpression()
	if ast.Err != nil {
		fmt.Println("ERROR: " + ast.Err.Error())
		return
	}
	fmt.Printf("ExpAST: %+v\n", ae.ToString())
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("ERROR: ", e)
		}
	}()
	start := time.Now()
	AddAstExpCache(exp, ae)
	r := ae.Eval()
	cost1 := time.Since(start)
	fmt.Println("time for Eval: " + cost1.String())
	fmt.Println("result ...\t", r)
	fmt.Printf("%s = %v\n", exp, r)

	if _, ok := FormulaCache[exp]; ok {
		start := time.Now()
		fmt.Println("cache  result ...\t", FormulaCache[exp].Eval())
		cost := time.Since(start)
		fmt.Println("time for cache Eval: " + cost.String())
	}

}
