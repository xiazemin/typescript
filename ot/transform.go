package ot

import (
	"fmt"
	"strconv"
)

// Transform takes two operations A and B that happened concurrently and
// produces two operations A' and B' (in an array) such that
// `apply(apply(S, A), B') = apply(apply(S, B), A')`. This function is the
// heart of OT.
func (o *Operation) transform(operation1, operation2 *Operation) (*Operation, *Operation, error) {
	if operation1 == nil || operation2 == nil || operation1.BaseLength != operation2.BaseLength {
		return nil, nil, fmt.Errorf("Both operations have to have the same base length")
	}

	//A = R(3),I('c')
	// B = R(3), I('d')

	operation1prime := &Operation{}
	operation2prime := &Operation{}

	ops1 := operation1.Ops
	ops2 := operation2.Ops
	var i1, i2 int
	op1 := ops1[i1]
	i1++
	op2 := ops2[i2]
	i2++

	for true {
		//insert是用string 表示的，retain用正数表示，delete 用负数表示
		// At every iteration of the loop, the imaginary cursor that both
		// operation1 and operation2 have that operates on the input string must
		// have the same position in the input string.

		if op1 == nil && op2 == nil {
			// end condition: both ops1 and ops2 have been processed
			break
		}
		if i1 == len(ops1) && i2 == len(ops2) {
			break
		}

		// next two cases: one or both ops are insert ops
		// => insert the string in the corresponding prime operation, skip it in
		// the other one. If both op1 and op2 are insert ops, prefer op1.
		if op1.IsInsert() {
			//先把插入处理完毕
			//在op1‘ 插入的同时需要在op2’ 增加一个retain
			operation1prime.insert(op1.StrVal)
			operation2prime.retain(strconv.FormatInt(int64(len(op1.StrVal)), 10))
			op1 = ops1[i1]
			i1++
			//A = R(3),I('c')
			// B = R(4), I('d')

			continue
		}
		if op2.IsInsert() {
			operation1prime.retain(strconv.FormatInt(int64(len(op2.StrVal)), 10))
			operation2prime.insert(op2.StrVal)
			op2 = ops2[i2]
			i2++
			//A = R(4),I('c')
			// B = R(4), I('d')
			continue
		}

		if op1 == nil {
			return nil, nil, fmt.Errorf("Cannot compose operations: first operation is too short.")
		}
		if op2 == nil {
			return nil, nil, fmt.Errorf("Cannot compose operations: first operation is too long.")
		}

		var minl int
		if op1.IsRetain() && op2.IsRetain() { //处理都是保留的情况
			// Simple case: retain/retain
			if op1.Val > op2.Val { //7， 3
				minl = int(op2.Val)         //   3
				op1.Val = op1.Val - op2.Val // 4    op1'+op2=op1转换后两个保留的和为保留最大长度
				op2 = ops2[i2]              //指向下一个坐标
				i2++
			} else if op1.Val == op2.Val {
				minl = int(op2.Val)
				op1 = ops1[i1]
				i1++
				op2 = ops2[i2]
				i2++
				//同时指向下一个坐标
			} else {
				minl = int(op1.Val)
				op2.Val = op2.Val - op1.Val
				op1 = ops1[i1]
				i1++
			}
			operation1prime.retain(strconv.FormatInt(int64(minl), 10)) //转换后保留两者的最小位置
			operation2prime.retain(strconv.FormatInt(int64(minl), 10))
		} else if op1.IsDelete() && op2.IsDelete() {
			// Both operations delete the same string at the same position. We don't
			// need to produce any operations, we just skip over the delete ops and
			// handle the case that one operation deletes more than the other.
			if -op1.Val > -op2.Val {
				op1.Val = op1.Val - op2.Val //都是删除的话，保留删除后的最小位置
				op2 = ops2[i2]
				i2++
			} else if op1.Val == op2.Val {
				op1 = ops1[i1]
				i1++
				op2 = ops2[i2]
				i2++
			} else {
				op2.Val = op2.Val - op1.Val
				op1 = ops1[i1]
				i1++
			}
			// next two cases: delete/retain and retain/delete
		} else if op1.IsDelete() && op1.IsRetain() {
			//一个删除一个保留的情况
			if -op1.Val > op2.Val {
				minl = int(op2.Val)
				op1.Val = op1.Val + op2.Val
				op2 = ops2[i2]
				i2++
			} else if -op1.Val == op2.Val {
				minl = int(op2.Val)
				op1 = ops1[i1]
				i1++
				op2 = ops2[i2]
				i2++
			} else {
				minl = int(-op1.Val)
				op2.Val = op2.Val + op1.Val
				op1 = ops1[i1]
				i1++
			}
			operation1prime.delete(strconv.FormatInt(int64(minl), 10))
		} else if op1.IsRetain() && op2.IsDelete() {
			if op1.Val > -op2.Val { //永远保留最小位置的情况
				minl = int(-op2.Val)
				op1.Val = op1.Val + op2.Val
				op2 = ops2[i2]
				i2++
			} else if op1.Val == -op2.Val {
				minl = int(op1.Val)
				op1 = ops1[i1]
				i1++
				op2 = ops2[i2]
				i2++
			} else {
				minl = int(op1.Val)
				op2.Val = op2.Val + op1.Val
				op1 = ops1[i1]
				i1++
			}
			operation2prime.delete(strconv.FormatInt(int64(minl), 10))
		} else {
			return nil, nil, fmt.Errorf("The two operations aren't compatible")
		}
	}

	return operation1prime, operation2prime, nil
}
