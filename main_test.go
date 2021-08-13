package main

import (
	"fmt"
	"formulaParser/formula"
	"testing"
	"time"
)

func TestExec(t *testing.T) {
	start := time.Now()
	exp := "5/(2-1),"
	formula.Exec(exp)
	cost1 := time.Since(start)
	fmt.Println("time: " + cost1.String())
}

func TestExec2(t *testing.T) {
	start := time.Now()
	exp := "5+4+3"
	formula.Exec(exp)
	cost1 := time.Since(start)
	fmt.Println("time: " + cost1.String())
}

func TestExec3(t *testing.T) {
	start := time.Now()
	exp := "(5-(1*3)/4)"
	formula.Exec(exp)
	cost1 := time.Since(start)
	fmt.Println("time: " + cost1.String())
}

func TestExecS4(t *testing.T) {
	start := time.Now()
	effects := make(map[string]interface{})
	effects["levela"] = float64(10)
	effects["levelb"] = float64(40)
	effects["crita"] = float64(30)
	effects["caritdefb"] = float64(10)
	effects["buffa"] =  map[string]bool{"1": true}
	formula.SetContext(effects)
	cost := time.Since(start)
	fmt.Println("time: " + cost.String())
	start1 := time.Now()
	exp := "max(0,crita-caritdefb)*0.8/(max(0,crita-caritdefb)+levelb*15+4500*ahasbuff(1))"
	formula.Exec(exp)
	cost1 := time.Since(start1)
	fmt.Println("time: " + cost1.String())

	/*ret := 20*0.8/(20+40*15+4500 * 1 )
	fmt.Println(ret)*/
}

func TestExecS5(t *testing.T) {
	start := time.Now()
	effects := make(map[string]interface{})
	effects["damage"] = float64(1000)
	effects["skilllevela"] = float64(10)
	effects["hpa"] = float64(5000)
	effects["attacka"] = float64(30)
	effects["caritdefb"] = float64(10)
	effects["buffb"] =  map[string]bool{"1": true}
	formula.SetContext(effects)
	cost := time.Since(start)
	fmt.Println("time: " + cost.String())
	start1 := time.Now()
	exp := "damage*(15+4500*ahasbuff(1))"
	formula.Exec(exp)
	cost1 := time.Since(start1)
	fmt.Println("time: " + cost1.String())
}