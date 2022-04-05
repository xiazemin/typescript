package ot

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Delta interface {
	Compose(operation2 Delta) (Delta, error)
	Transform(operation1In, operation2In Delta) (Delta, Delta, error)
	Apply(str string) (string, error)
	retain(n string) (Delta, error)
	insert(str string) (Delta, error)
	delete(n string) (Delta, error)
}

type QuillDelta struct {
	Ops          []Op `json:"ops"`
	TargetLength int
	BaseLength   int
}

func (o *QuillDelta) Compose(operationIn Delta) (Delta, error) {
	operation2, ok := operationIn.(*QuillDelta)
	if !ok {
		return nil, fmt.Errorf("operationIn not QuillDelta")
	}
	operation1 := o
	if operation1.TargetLength != operation2.BaseLength {
		//后面的操作必须是紧紧挨者前面的操作的
		return nil, fmt.Errorf("The base length of the second operation has to be the target length of the first operation")
	}

	operation := &QuillDelta{} // the combined operation
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
		if op1 == nil && op2 == nil {
			break
		}

		//把ops1的所有删除操作执行完毕
		if op1 != nil && op1.IsDelete() {
			operation.delete(op1.GetStringVal())
			op1, i1 = next(ops1, op1, i1)
			continue
		}

		//把ops2的所有插入操作执行完毕
		if op2 != nil && op2.IsInsert() {
			operation.insert(op2.GetStringVal())
			op2, i2 = next(ops2, op2, i2)
			continue
		}

		//剩下的情况如下(因为相同的项在构建的时候已经合并了)
		//  op1: retain insert ;retain delete; insert retain;insert delete
		//  op2:	retain delete;retain insert; delete insert；delete retaim;

		// op1,op2 的组合有 retain retain； retain delete;insert retain；insert delete

		if op1 == nil || i1 == len(ops1) {
			fmt.Println("Cannot compose operations: first operation is too short.")
			for op2 != nil {
				if op2.IsRetain() {
					operation.retain(op2.GetStringVal())
				} else if op2.IsInsert() {
					operation.insert(op2.GetStringVal())
				} else if op2.IsDelete() {
					operation.delete(op2.GetStringVal())
				}
				op2, i2 = next(ops2, op2, i2)
			}
			break
		}
		if op2 == nil || i2 == len(ops2) {
			fmt.Println("Cannot compose operations: first operation is too long.")
			for op1 != nil {
				if op1.IsRetain() {
					operation.retain(op1.GetStringVal())
				} else if op1.IsInsert() {
					operation.insert(op1.GetStringVal())
				} else if op1.IsDelete() {
					operation.delete(op1.GetStringVal())
				}
				op1, i1 = next(ops1, op1, i1)
			}
			break
		}

		if op1.IsRetain() && op2.IsRetain() {
			if op1.GetVal() > op2.GetVal() {
				//都是重定位操作，重定位到更近的位置
				operation.retain(op2.GetStringVal())
				//把op1操作的重定位位置减小，保证，两个的和等于op1操作的位置
				op1.SetVal(op1.GetVal() - op2.GetVal())
				op2, i2 = next(ops2, op2, i2)
			} else if op1.GetVal() == op2.GetVal() {
				operation.retain(op1.GetStringVal())
				op1, i1 = next(ops1, op1, i1)
				op2, i2 = next(ops2, op2, i2)
			} else {
				operation.retain(op1.GetStringVal())
				op2.SetVal(op2.GetVal() - op1.GetVal())
				op1, i1 = next(ops1, op1, i1)
			}
		} else if op1.IsInsert() && op2.IsDelete() {
			if op1.GetVal() > -op2.GetVal() {
				//直接在op1上删除 ，丢弃删除操作 op2
				strVal := op1.GetStringVal()
				op1.SetVal(strVal[:op1.GetVal()+op2.GetVal()])
				op2, i2 = next(ops2, op2, i2)
			} else if op1.GetVal() == -op2.GetVal() {
				op1, i1 = next(ops1, op1, i1)
				op2, i2 = next(ops2, op2, i2)
			} else {
				//直接减少，删除的长度，丢弃插入操作
				op2.SetVal(op2.GetVal() + op1.GetVal())
				op1, i1 = next(ops1, op1, i1)
			}
		} else if op1.IsInsert() && op2.IsRetain() {
			if op1.GetVal() > op2.GetVal() {
				//把第一个插入拆成两部分操作
				operation.insert(op1.GetStringVal()[0:op2.GetVal()])
				op1.SetVal(op1.GetStringVal()[op2.GetVal():])
				op2, i2 = next(ops2, op2, i2)
			} else if op1.GetVal() == op2.GetVal() {
				operation.insert(op1.GetStringVal())
				op1, i1 = next(ops1, op1, i1)
				op2, i2 = next(ops2, op2, i2)
			} else {
				operation.insert(op1.GetStringVal())
				op2.SetVal(op2.GetVal() - op1.GetVal())
				op1, i1 = next(ops1, op1, i1)
			}
		} else if op1.IsRetain() && op2.IsDelete() {
			if op1.GetVal() > -op2.GetVal() {
				operation.delete(op2.GetStringVal())
				op1.SetVal(op1.GetVal() + op2.GetVal())
				op2, i2 = next(ops2, op2, i2)
			} else if op1.GetVal() == -op2.GetVal() {
				operation.delete(op2.GetStringVal())
				op1, i1 = next(ops1, op1, i1)
				op2, i2 = next(ops2, op2, i2)
			} else {
				operation.delete(op1.GetStringVal())
				op2.SetVal(op2.GetVal() + op1.GetVal())
				op1, i1 = next(ops1, op1, i1)
			}
		} else {
			d1, _ := json.Marshal(op1)
			d2, _ := json.Marshal(op2)
			return nil, fmt.Errorf("This shouldn't happen: op1: " + string(d1) + ", op2: " + string(d2))
		}
	}
	return operation, nil
}

//跳过n个字符
func (o *QuillDelta) retain(n string) (Delta, error) {
	if o == nil {
		return nil, fmt.Errorf("current operation is nil")
	}
	valInt64, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		return nil, err
	}
	val := int(valInt64)
	if val == 0 {
		return o, nil
	}
	o.BaseLength += val
	o.TargetLength += val

	if len(o.Ops) > 0 && o.Ops[len(o.Ops)-1] != nil && o.Ops[len(o.Ops)-1].IsRetain() {
		//最后一个操作是保留，那么可以把输入的保留合并进来,上一次retain 5，本次是3 ，结果就是 retain 8
		o.Ops[len(o.Ops)-1].SetVal(o.Ops[len(o.Ops)-1].GetVal() + val)
	} else {
		o.Ops = append(o.Ops, &RetainOp{
			Retain: val,
		}) //创建一个新的操作
	}
	return o, nil
}

//插入字符 str
func (o *QuillDelta) insert(str string) (Delta, error) {
	if o == nil {
		return nil, fmt.Errorf("current operation is nil")
	}
	if str == "" {
		return o, nil
	}
	o.TargetLength += len(str)
	if len(o.Ops) > 0 && o.Ops[len(o.Ops)-1].IsInsert() {
		// Merge insert op.
		//如果上一个操作是插入，本当前操作也合并进来
		o.Ops[len(o.Ops)-1].SetVal(o.Ops[len(o.Ops)-1].GetStringVal() + str)
	} else if len(o.Ops) > 0 && o.Ops[len(o.Ops)-1].IsDelete() {
		// It doesn't matter when an operation is applied whether the operation
		// is delete(3), insert("something") or insert("something"), delete(3).
		// Here we enforce that in this case, the insert op always comes first.
		// This makes all operations that have the same effect when applied to
		// a document of the right length equal in respect to the `equals` method.
		//我们忽略删除和插入的顺序，因为对所有的文档采用同样的操作，结果是一样的
		//我们定义一个策略，永远把插入放在删除前面
		if len(o.Ops) > 1 && o.Ops[len(o.Ops)-2].IsInsert() {
			//如果上一个操作是删除，上上一个是插入，把当前插入，放进去
			o.Ops[len(o.Ops)-2].SetVal(o.Ops[len(o.Ops)-2].GetStringVal() + str)
		} else {
			//上上一个操作是retain或者删除，用当前字符串替换
			o.Ops = append(o.Ops, &DeleteOp{
				Delete: o.Ops[len(o.Ops)-1].GetVal(),
			})
			//ops[ops.length] = ops[ops.length-1]
			//ops[ops.length-2] = str
			//这里是个bug吧应该是 ops[ops.length-1] = str;
			o.Ops[len(o.Ops)-1] = &InsertOp{
				Insert: str,
			}
		}
	} else {
		//上一个操作是retain，把当前插入放在后面
		o.Ops = append(o.Ops, &InsertOp{
			Insert: str,
		})
	}
	return o, nil
}

func (o *QuillDelta) delete(n string) (Delta, error) {
	if o == nil {
		return nil, fmt.Errorf("current operation is nil")
	}
	valInt64, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		return o, err
	}
	val := int(valInt64)
	if val == 0 {
		return o, nil
	}

	if val > 0 {
		fmt.Println("删除操作必须是负数", val)
		val = -val
	}
	o.BaseLength -= val //把baseLength加上删除长度，
	if len(o.Ops) > 0 && o.Ops[len(o.Ops)-1].IsDelete() {
		o.Ops[len(o.Ops)-1].SetVal(o.Ops[len(o.Ops)-1].GetVal() + val)
	} else {
		o.Ops = append(o.Ops, &DeleteOp{
			Delete: val,
		})
	}
	return o, nil
}

func next(ops []Op, op Op, i int) (Op, int) {
	if i < len(ops) {
		op = ops[i]
	} else {
		op = nil
	}
	i++
	return op, i
}

// Transform takes two operations A and B that happened concurrently and
// produces two operations A' and B' (in an array) such that
// `apply(apply(S, A), B') = apply(apply(S, B), A')`. This function is the
// heart of OT.
func (o *QuillDelta) Transform(operation1In, operation2In Delta) (Delta, Delta, error) {
	operation1, ok := operation1In.(*QuillDelta)
	if !ok {
		return nil, nil, fmt.Errorf("operation1In not QuillDelta")
	}

	operation2, ok := operation2In.(*QuillDelta)
	if !ok {
		return nil, nil, fmt.Errorf("operation2In not QuillDelta")
	}

	if operation1 == nil || operation2 == nil || operation1.BaseLength != operation2.BaseLength {
		return nil, nil, fmt.Errorf("Both operations have to have the same base length")
	}

	//A = R(3),I('c')
	// B = R(3), I('d')

	operation1prime := &QuillDelta{}
	operation2prime := &QuillDelta{}

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
		if op1 != nil && op1.IsInsert() {
			//先把插入处理完毕
			//在op1‘ 插入的同时需要在op2’ 增加一个retain
			operation1prime.insert(op1.GetStringVal())
			operation2prime.retain(strconv.FormatInt(int64(op1.GetVal()), 10))
			op1, i1 = next(ops1, op1, i1)
			//A = R(3),I('c')
			// B = R(4), I('d')

			continue
		}
		fmt.Println(i1, i2)
		if op2 != nil && op2.IsInsert() {
			operation1prime.retain(strconv.FormatInt(int64(op1.GetVal()), 10))
			operation2prime.insert(op2.GetStringVal())
			op2, i2 = next(ops2, op2, i2)
			//A = R(4),I('c')
			// B = R(4), I('d')
			continue
		}

		if op1 == nil {
			fmt.Println("Cannot compose operations: first operation is too short.")
			for i2 <= len(ops2) {
				fmt.Println("xxxxxxx", i2, len(ops2), op2.GetStringVal())
				if op2.IsRetain() {
					operation2prime.retain(op2.GetStringVal())
				} else if op2.IsInsert() {
					operation2prime.insert(op2.GetStringVal())
				} else if op2.IsDelete() {
					operation2prime.delete(op2.GetStringVal())
				}
				op2, i2 = next(ops2, op2, i2)
			}
			return operation1prime, operation2prime, nil
		}
		if op2 == nil {
			for i1 < len(ops1) {
				if op1.IsRetain() {
					operation1prime.retain(op1.GetStringVal())
				} else if op1.IsInsert() {
					operation1prime.insert(op1.GetStringVal())
				} else if op1.IsDelete() {
					operation1prime.delete(op1.GetStringVal())
				}
				op1, i1 = next(ops1, op1, i1)
			}
			fmt.Println("Cannot compose operations: first operation is too long.")
			return operation1prime, operation2prime, nil
		}

		var minl int
		if op1.IsRetain() && op2.IsRetain() { //处理都是保留的情况
			// Simple case: retain/retain
			if op1.GetVal() > op2.GetVal() { //7， 3
				minl = int(op2.GetVal())                //   3
				op1.SetVal(op1.GetVal() - op2.GetVal()) // 4    op1'+op2=op1转换后两个保留的和为保留最大长度
				op2, i2 = next(ops2, op2, i2)           //指向下一个坐标

			} else if op1.GetVal() == op2.GetVal() {
				minl = int(op2.GetVal())
				op1, i1 = next(ops1, op1, i1)
				op2, i2 = next(ops2, op2, i2)
				//同时指向下一个坐标
			} else {
				minl = int(op1.GetVal())
				op2.SetVal(op2.GetVal() - op1.GetVal())
				op1, i1 = next(ops1, op1, i1)
				fmt.Println("xiazemin 1:", minl, op1.GetVal(), op1.IsInsert(), op2.GetVal(), op2.IsRetain())
			}
			operation1prime.retain(strconv.FormatInt(int64(minl), 10)) //转换后保留两者的最小位置
			operation2prime.retain(strconv.FormatInt(int64(minl), 10))

			ABData, _ := json.Marshal(operation1prime)
			BAData, _ := json.Marshal(operation2prime)
			fmt.Println("xiazemin **** A':", string(ABData), "\nB':", string(BAData))
		} else if op1.IsDelete() && op2.IsDelete() {
			// Both operations delete the same string at the same position. We don't
			// need to produce any operations, we just skip over the delete ops and
			// handle the case that one operation deletes more than the other.
			if -op1.GetVal() > -op2.GetVal() {
				op1.SetVal(op1.GetVal() - op2.GetVal()) //都是删除的话，保留删除后的最小位置
				op2, i2 = next(ops2, op2, i2)
			} else if op1.GetVal() == op2.GetVal() {
				op1, i1 = next(ops1, op1, i1)
				op2, i2 = next(ops2, op2, i2)
			} else {
				op2.SetVal(op2.GetVal() - op1.GetVal())
				op1, i1 = next(ops1, op1, i1)
			}
			// next two cases: delete/retain and retain/delete
		} else if op1.IsDelete() && op2.IsRetain() {
			//一个删除一个保留的情况
			if -op1.GetVal() > op2.GetVal() {
				minl = int(op2.GetVal())
				op1.SetVal(op1.GetVal() + op2.GetVal())
				op2, i2 = next(ops2, op2, i2)
			} else if -op1.GetVal() == op2.GetVal() {
				minl = int(op2.GetVal())
				op1, i1 = next(ops1, op1, i1)
				op2, i2 = next(ops2, op2, i2)
			} else {
				minl = int(-op1.GetVal())
				op2.SetVal(op2.GetVal() + op1.GetVal())
				op1, i1 = next(ops1, op1, i1)
			}
			operation1prime.delete(strconv.FormatInt(int64(minl), 10))
		} else if op1.IsRetain() && op2.IsDelete() {
			if op1.GetVal() > -op2.GetVal() { //永远保留最小位置的情况
				minl = int(-op2.GetVal())
				op1.SetVal(op1.GetVal() + op2.GetVal())
				op2, i2 = next(ops2, op2, i2)
			} else if op1.GetVal() == -op2.GetVal() {
				minl = int(op1.GetVal())
				op1, i1 = next(ops1, op1, i1)
				op2, i2 = next(ops2, op2, i2)
			} else {
				minl = int(op1.GetVal())
				op2.SetVal(op2.GetVal() + op1.GetVal())
				op1, i1 = next(ops1, op1, i1)
			}
			operation2prime.delete(strconv.FormatInt(int64(minl), 10))
		} else {
			op1Data, _ := json.Marshal(op1)
			op2Data, _ := json.Marshal(op2)
			return nil, nil, fmt.Errorf("The two operations aren't compatible %s,%s", string(op1Data), string(op2Data))
		}
	}

	return operation1prime, operation2prime, nil
}

// Apply an operation to a string, returning a new string. Throws an error if
// there's a mismatch between the input string and the operation.
//把操作应用到str上产生一个新的操作
func (o *QuillDelta) Apply(str string) (string, error) {
	var operation = o
	if len(str) != operation.BaseLength {
		return "", fmt.Errorf("The operation's base length must be equal to the string's length.")
	}

	var newStr []string
	var j, strIndex int
	for i := 0; i < len(o.Ops); i++ {
		if o.Ops[i].IsRetain() {
			if strIndex+o.Ops[i].GetVal() > len(str) {
				return "", fmt.Errorf("QuillDelta can't retain more characters than are left in the string.")
			}
			// Copy skipped part of the old string.
			newStr[j] = str[strIndex : strIndex+int(o.Ops[i].GetVal())]
			j++
			strIndex += int(o.Ops[i].GetVal())
		} else if o.Ops[i].IsInsert() {
			// Insert string.
			newStr[j] = o.Ops[i].GetStringVal()
			j++
		} else { // delete op
			strIndex -= int(o.Ops[i].GetVal())
		}
	}
	if strIndex != len(str) {
		return "", fmt.Errorf("The operation didn't operate on the whole string.")
	}
	return strings.Join(newStr, ""), nil
}
