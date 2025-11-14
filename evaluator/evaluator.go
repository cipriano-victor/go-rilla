package evaluator

import (
	"fmt"
	"go-rilla/ast"
	"go-rilla/object"
	"math"
)

var (
	NULL           = &object.Null{}
	TRUE           = &object.Boolean{Value: true}
	FALSE          = &object.Boolean{Value: false}
	breakSignal    = &object.Break{}
	continueSignal = &object.Continue{}
)

var loopDepth int

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := evalLetValue(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.BreakStatement:
		if loopDepth == 0 {
			return newError("break statement outside of loop")
		}
		return breakSignal
	case *ast.ContinueStatement:
		if loopDepth == 0 {
			return newError("continue statement outside of loop")
		}
		return continueSignal
	// Expressions
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.PrefixExpression:
		if node.Operator == "++" || node.Operator == "--" {
			return evalPrefixUpdateExpression(node, env)
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		if node.Operator == "=" {
			return evalAssignmentExpression(node, env)
		}
		if isCompoundAssignment(node.TokenLiteral()) {
			return evalCompoundAssignment(node, env)
		}
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right, node.Operator)
	case *ast.PostfixExpression:
		return evalPostfixExpression(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.WhileExpression:
		return evalWhileExpression(node, env)
	}

	return nil
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s %s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ && right.Type() != object.FLOAT_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	if right.Type() == object.INTEGER_OBJ {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	} else {
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}
}

func evalInfixExpression(operator string, left, right object.Object, display string) object.Object {
	if display == "" {
		display = operator
	}
	switch {
	case (left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ) || (left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ) || (left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ) || (left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ):
		return evalNumberInfixExpression(operator, left, right, display)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right, display)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		leftVal := left.(*object.Boolean).Value
		rightVal := right.(*object.Boolean).Value
		switch operator {
		case "&&":
			return nativeBoolToBooleanObject(leftVal && rightVal)
		case "||":
			return nativeBoolToBooleanObject(leftVal || rightVal)
		case "==":
			return nativeBoolToBooleanObject(leftVal == rightVal)
		case "!=":
			return nativeBoolToBooleanObject(leftVal != rightVal)
		}
		return newError("unknown operator: %s %s %s",
			left.Type(), display, right.Type())
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), display, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), display, right.Type())
	}
}

func isCompoundAssignment(tokenLiteral string) bool {
	return tokenLiteral == "+=" || tokenLiteral == "-="
}

func evalCompoundAssignment(node *ast.InfixExpression, env *object.Environment) object.Object {
	current := Eval(node.Left, env)
	if isError(current) {
		return current
	}

	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}

	result := evalInfixExpression(node.Operator, current, right, node.TokenLiteral())
	if isError(result) {
		return result
	}

	identifier, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newError("invalid assignment target: %s", node.Left.TokenLiteral())
	}

	env.Set(identifier.Value, result)
	return result
}

func evalNumberInfixExpression(operator string, left, right object.Object, display string) object.Object {
	var leftVal float64
	var rightVal float64

	if left.Type() == object.INTEGER_OBJ {
		leftVal = float64(left.(*object.Integer).Value)
	} else {
		leftVal = left.(*object.Float).Value
	}

	if right.Type() == object.INTEGER_OBJ {
		rightVal = float64(right.(*object.Integer).Value)
	} else {
		rightVal = right.(*object.Float).Value
	}

	switch operator {
	case "+":
		if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
			return &object.Integer{Value: int64(leftVal + rightVal)}
		}
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
			return &object.Integer{Value: int64(leftVal - rightVal)}
		}
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
			return &object.Integer{Value: int64(leftVal * rightVal)}
		}
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
			return &object.Integer{Value: int64(leftVal / rightVal)}
		}
		return &object.Float{Value: leftVal / rightVal}
	case "**":
		result := math.Pow(leftVal, rightVal)
		if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
			return &object.Integer{Value: int64(result)}
		}
		return &object.Float{Value: result}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), display, right.Type())
	}
}

func evalPrefixUpdateExpression(node *ast.PrefixExpression, env *object.Environment) object.Object {
	ident, ok := node.Right.(*ast.Identifier)
	if !ok {
		return newError("invalid prefix target: %s", node.Right.TokenLiteral())
	}

	current, ok := env.Get(ident.Value)
	if !ok {
		return newError("identifier not found: %s", ident.Value)
	}

	var result object.Object

	switch node.Operator {
	case "++":
		switch current.Type() {
		case object.INTEGER_OBJ:
			value := current.(*object.Integer).Value + 1
			result = &object.Integer{Value: value}
		case object.FLOAT_OBJ:
			value := current.(*object.Float).Value + 1
			result = &object.Float{Value: value}
		default:
			return newError("unknown operator: %s%s", current.Type(), node.Operator)
		}
	case "--":
		switch current.Type() {
		case object.INTEGER_OBJ:
			value := current.(*object.Integer).Value - 1
			result = &object.Integer{Value: value}
		case object.FLOAT_OBJ:
			value := current.(*object.Float).Value - 1
			result = &object.Float{Value: value}
		default:
			return newError("unknown operator: %s%s", current.Type(), node.Operator)
		}
	default:
		return newError("unknown operator: %s%s", current.Type(), node.Operator)
	}

	env.Set(ident.Value, result)
	return result
}

func evalLetValue(expr ast.Expression, env *object.Environment) object.Object {
	if prefix, ok := expr.(*ast.PrefixExpression); ok && (prefix.Operator == "++" || prefix.Operator == "--") {
		ident, ok := prefix.Right.(*ast.Identifier)
		if !ok {
			return newError("invalid prefix target: %s", prefix.Right.TokenLiteral())
		}

		current, ok := env.Get(ident.Value)
		if !ok {
			return newError("identifier not found: %s", ident.Value)
		}

		var updated object.Object

		switch prefix.Operator {
		case "++":
			switch current.Type() {
			case object.INTEGER_OBJ:
				value := current.(*object.Integer).Value + 1
				updated = &object.Integer{Value: value}
			case object.FLOAT_OBJ:
				value := current.(*object.Float).Value + 1
				updated = &object.Float{Value: value}
			default:
				return newError("unknown operator: %s%s", current.Type(), prefix.Operator)
			}
		case "--":
			switch current.Type() {
			case object.INTEGER_OBJ:
				value := current.(*object.Integer).Value - 1
				updated = &object.Integer{Value: value}
			case object.FLOAT_OBJ:
				value := current.(*object.Float).Value - 1
				updated = &object.Float{Value: value}
			default:
				return newError("unknown operator: %s%s", current.Type(), prefix.Operator)
			}
		default:
			return newError("unknown operator: %s%s", current.Type(), prefix.Operator)
		}

		env.Set(ident.Value, updated)
		return current
	}

	return Eval(expr, env)
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
	display string,
) object.Object {

	if operator != "+" {
		return newError("unknown operator: %s %s %s",
			left.Type(), display, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalWhileExpression(node *ast.WhileExpression, env *object.Environment) object.Object {
	if node.Init != nil {
		result := Eval(node.Init, env)
		if shouldHaltLoop(result) {
			return result
		}
	}

	var loopResult object.Object

	for {
		if node.Condition != nil {
			condition := Eval(node.Condition, env)
			if shouldHaltLoop(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		loopDepth++
		bodyResult := Eval(node.Body, env)
		loopDepth--

		if isBreak(bodyResult) {
			return loopResult
		}

		shouldContinue := false
		if isContinue(bodyResult) {
			shouldContinue = true
		} else if shouldHaltLoop(bodyResult) {
			return bodyResult
		} else if bodyResult != nil {
			loopResult = bodyResult
		}

		if node.Post != nil {
			postResult := Eval(node.Post, env)
			if shouldHaltLoop(postResult) {
				return postResult
			}
		}

		if shouldContinue {
			continue
		}
	}

	return loopResult
}

func shouldHaltLoop(obj object.Object) bool {
	if obj == nil {
		return false
	}
	t := obj.Type()
	return t == object.RETURN_VALUE_OBJ || t == object.ERROR_OBJ
}

func isBreak(obj object.Object) bool {
	if obj == nil {
		return false
	}
	return obj.Type() == object.BREAK_OBJ
}

func isContinue(obj object.Object) bool {
	if obj == nil {
		return false
	}
	return obj.Type() == object.CONTINUE_OBJ
}

func evalAssignmentExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	identifier, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newError("invalid assignment target: %s", node.Left.TokenLiteral())
	}

	value := Eval(node.Right, env)
	if isError(value) {
		return value
	}

	env.Set(identifier.Value, value)
	return value
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		case *object.Break:
			return newError("break statement outside of loop")
		case *object.Continue:
			return newError("continue statement outside of loop")
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ ||
				rt == object.ERROR_OBJ ||
				rt == object.BREAK_OBJ ||
				rt == object.CONTINUE_OBJ {
				return result
			}
		}
	}
	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", node.Value)
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		savedLoopDepth := loopDepth
		loopDepth = 0
		evaluated := Eval(fn.Body, extendedEnv)
		loopDepth = savedLoopDepth
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObject.Elements[idx]
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	identifier, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newError("invalid postfix target: %s", node.Left.TokenLiteral())
	}

	current := Eval(node.Left, env)
	if isError(current) {
		return current
	}

	var result object.Object

	switch node.Operator {
	case "++":
		switch current.Type() {
		case object.INTEGER_OBJ:
			value := current.(*object.Integer).Value + 1
			result = &object.Integer{Value: value}
		case object.FLOAT_OBJ:
			value := current.(*object.Float).Value + 1
			result = &object.Float{Value: value}
		default:
			return newError("unknown operator: %s%s", current.Type(), node.Operator)
		}
	case "--":
		switch current.Type() {
		case object.INTEGER_OBJ:
			value := current.(*object.Integer).Value - 1
			result = &object.Integer{Value: value}
		case object.FLOAT_OBJ:
			value := current.(*object.Float).Value - 1
			result = &object.Float{Value: value}
		default:
			return newError("unknown operator: %s%s", current.Type(), node.Operator)
		}
	default:
		return newError("unknown operator: %s%s", current.Type(), node.Operator)
	}

	env.Set(identifier.Value, result)
	return current
}
