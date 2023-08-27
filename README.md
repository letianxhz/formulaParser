# formulaParser
game skill formula parser

go实现的技能公式表达式解析器
支持变量，自定义函数

示例:
```
 effects := make(map[string]interface{})
	effects["levela"] = float64(10)
	effects["levelb"] = float64(40)
	effects["crita"] = float64(30)
	effects["caritdefb"] = float64(10)
	effects["buffa"] =  map[string]bool{"1": true}
  
  effects 为技能需要的参数。在实际游戏中，从施法者和受击者身上获取
	formula.SetContext(effects)
	exp := "max(0,crita-caritdefb)*0.8/(max(0,crita-caritdefb)+levelb*15+4500*ahasbuff(1))"
	formula.Exec(exp)
```
 
