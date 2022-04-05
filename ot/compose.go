package ot

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Compose merges two consecutive operations into one operation, that
// preserves the changes of both. Or, in other words, for each input string S
// and a pair of consecutive operations A and B,
// apply(apply(S, A), B) = apply(S, compose(A, B)) must hold.
func (o *Operation) Compose(operation2 *Operation) (*Operation, error) {
	operation1 := o
	if operation1.TargetLength != operation2.BaseLength {
		//后面的操作必须是紧紧挨者前面的操作的
		return nil, fmt.Errorf("The base length of the second operation has to be the target length of the first operation")
	}

	operation := &Operation{} // the combined operation
	ops1 := operation1.Ops
	ops2 := operation2.Ops // for fast access
	var i1, i2 int         // current index into ops1 respectively ops2
	op1 := ops1[i1]
	op2 := ops2[i2] // current ops
	i1++
	i2++

	for true {
		// Dispatch on the type of op1 and op2
		if i1 == len(ops1) && i2 == len(ops2) {
			// end condition: both ops1 and ops2 have been processed
			break
		}

		//把ops1的所有删除操作执行完毕
		if op1.IsDelete() {
			operation.delete(strconv.FormatInt(op1.Val, 10))
			op1 = ops1[i1]
			i1++
			continue
		}

		//把ops2的所有插入操作执行完毕
		if op2.IsInsert() {
			operation.insert(op2.StrVal)
			op2 = ops2[i2]
			i2++
			continue
		}

		//剩下的情况如下(因为相同的项在构建的时候已经合并了)
		//  op1: retain insert ;retain delete; insert retain;insert delete
		//  op2:	retain delete;retain insert; delete insert；delete retaim;

		// op1,op2 的组合有 retain retain； retain delete;insert retain；insert delete

		if op1 == nil || i1 == len(ops1) {
			return nil, fmt.Errorf("Cannot compose operations: first operation is too short.")
		}
		if op2 == nil || i2 == len(ops2) {
			return nil, fmt.Errorf("Cannot compose operations: first operation is too long.")
		}

		if op1.IsRetain() && op2.IsRetain() {
			if op1.Val > op2.Val {
				//都是重定位操作，重定位到更近的位置
				operation.retain(strconv.FormatInt(op2.Val, 10))
				//把op1操作的重定位位置减小，保证，两个的和等于op1操作的位置
				op1.Val -= op2.Val
				op2 = ops2[i2]
				i2++
			} else if op1.Val == op2.Val {
				operation.retain(strconv.FormatInt(op1.Val, 10))
				op1 = ops1[i1]
				i1++
				op2 = ops2[i2]
				i2++
			} else {
				operation.retain(strconv.FormatInt(op1.Val, 10))
				op2.Val -= op1.Val
				op1 = ops1[i1]
				i1++
			}
		} else if op1.IsInsert() && op2.IsDelete() {
			if len(op1.StrVal) > int(-op2.Val) {
				//直接在op1上删除 ，丢弃删除操作 op2
				op1.StrVal = op1.StrVal[:len(op1.StrVal)+int(op2.Val)]
				op2 = ops2[i2]
				i2++
			} else if len(op1.StrVal) == int(-op2.Val) {
				op1 = ops1[i1]
				i1++
				op2 = ops2[i2]
				i2++
			} else {
				//直接减少，删除的长度，丢弃插入操作
				op2.Val = op2.Val + int64(len(op1.StrVal))
				op1 = ops1[i1]
				i1++
			}
		} else if op1.IsInsert() && op2.IsRetain() {
			if len(op1.StrVal) > int(op2.Val) {
				//把第一个插入拆成两部分操作
				operation.insert(op1.StrVal[0:op2.Val])
				op1.StrVal = op1.StrVal[op2.Val:]
				op2 = ops2[i2]
				i2++
			} else if len(op1.StrVal) == int(op2.Val) {
				operation.insert(op1.StrVal)
				op1 = ops1[i1]
				i1++
				op2 = ops2[i2]
				i2++
			} else {
				operation.insert(op1.StrVal)
				op2.Val -= int64(len(op1.StrVal))
				op1 = ops1[i1]
				i1++
			}
		} else if op1.IsRetain() && op2.IsDelete() {
			if op1.Val > -op2.Val {
				operation.delete(strconv.FormatInt(op2.Val, 10))
				op1.Val += op2.Val
				op2 = ops2[i2]
				i2++
			} else if op1.Val == -op2.Val {
				operation.delete(strconv.FormatInt(op2.Val, 10))
				op1 = ops1[i1]
				i1++
				op2 = ops2[i2]
				i2++
			} else {
				operation.delete(strconv.FormatInt(op1.Val, 10))
				op2.Val += op1.Val
				op1 = ops1[i1]
				i1++
			}
		} else {
			d1, _ := json.Marshal(op1)
			d2, _ := json.Marshal(op2)
			return nil, fmt.Errorf("This shouldn't happen: op1: " + string(d1) + ", op2: " + string(d2))
		}
	}
	return operation, nil
}

func getSimpleOp(operation *Operation) *Operation {
	if operation == nil {
		return operation
	}

	switch len(operation.Ops) {
	case 1:
		return operation.Ops[0]
	case 2:
		if operation.Ops[0].IsRetain() {
			return operation.Ops[1]
		}
		if operation.Ops[1].IsRetain() {
			return operation.Ops[0]
		}
		return nil
	case 3:
		if operation.Ops[0].IsRetain() && operation.Ops[2].IsRetain() {
			return operation.Ops[1]
		}
	}
	return nil
}

func getStartIndex(operation *Operation) int {
	if operation == nil {
		return 0
	}
	if operation.Ops[0].IsRetain() {
		return int(operation.Ops[0].Val)
	}
	return 0
}
