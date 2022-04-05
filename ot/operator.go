package ot

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Operation are essentially lists of ops. There are three types of ops:
// * Retain ops: Advance the cursor position by a given number of characters.  Retain 移动游标到指定位置  n是正数
//   Represented by positive ints.
// * Insert ops: Insert a given string at the current cursor position.   Insert 在游标位置插入一个字符串
//   Represented by strings.
// * Delete ops: Delete the next n characters. Represented by negative ints.  Delete 删除紧接着的n个字符，n是负数
type OPType int

const (
	Retain OPType = iota
	Insert
	Delete
)

type Operation struct {
	// actions (skip/delete/insert) are stored as an array
	//用户提交的操作，包含了一系列的子操作
	Ops []*Operation
	// 执行operation操作之前的长度
	BaseLength int64
	// 执行operation操作之后的长度
	TargetLength int64
	Type         OPType
	Val          int64  //如果是retain 这里是个正数，如果是删除这里是个负数
	StrVal       string //如果是 插入这里是个字符串
}

func (o *Operation) equal(other *Operation) bool {
	if o == nil && other == nil {
		return true
	}
	if o != nil && other == nil {
		return false
	}
	if o == nil && other != nil {
		return false
	}
	if o.BaseLength != other.BaseLength {
		return false
	}
	if o.TargetLength != other.TargetLength {
		return false
	}
	if len(o.Ops) != len(other.Ops) {
		return false
	}
	for i := 0; i < len(o.Ops); i++ {
		if o.Ops[i].equal(other.Ops[i]) {
			return false
		}
	}
	return true
}

func (o *Operation) IsRetain() bool {
	return o.Type == Retain
}
func (o *Operation) IsInsert() bool {
	return o.Type == Insert
}
func (o *Operation) IsDelete() bool {
	return o.Type == Delete
}

//跳过n个字符
func (o *Operation) retain(n string) (*Operation, error) {
	if o == nil {
		return nil, fmt.Errorf("current operation is nil")
	}
	val, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		return nil, err
	}
	if val == 0 {
		return o, nil
	}
	o.BaseLength += val
	o.TargetLength += val

	if len(o.Ops) > 0 && o.Ops[len(o.Ops)-1] != nil && o.Ops[len(o.Ops)-1].IsRetain() {
		//最后一个操作是保留，那么可以把输入的保留合并进来,上一次retain 5，本次是3 ，结果就是 retain 8
		o.Ops[len(o.Ops)-1].Val += val
	} else {
		newOp := &Operation{
			Type:         Retain,
			BaseLength:   val,
			TargetLength: val,
		}
		o.Ops = append(o.Ops, newOp) //创建一个新的操作
	}
	return o, nil
}

//插入字符 str
func (o *Operation) insert(str string) (*Operation, error) {
	if o == nil {
		return nil, fmt.Errorf("current operation is nil")
	}
	if str == "" {
		return o, nil
	}
	o.TargetLength += int64(len(str))
	if len(o.Ops) > 0 && o.Ops[len(o.Ops)-1].IsInsert() {
		// Merge insert op.
		//如果上一个操作是插入，本当前操作也合并进来
		o.Ops[len(o.Ops)-1].StrVal += str
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
			o.Ops[len(o.Ops)-2].StrVal += str
		} else {
			//上上一个操作是retain或者删除，用当前字符串替换
			o.Ops = append(o.Ops, &Operation{
				StrVal:       o.Ops[len(o.Ops)-1].StrVal,
				Type:         o.Ops[len(o.Ops)-1].Type,
				BaseLength:   o.Ops[len(o.Ops)-1].BaseLength,
				TargetLength: o.Ops[len(o.Ops)-1].TargetLength,
			})
			//ops[ops.length] = ops[ops.length-1]
			//ops[ops.length-2] = str
			//这里是个bug吧应该是 ops[ops.length-1] = str;
			o.Ops[len(o.Ops)-1] = &Operation{
				StrVal:       str,
				Type:         Insert,
				BaseLength:   int64(len(str)),
				TargetLength: int64(len(str)),
			}
		}
	} else {
		//上一个操作是retain，把当前插入放在后面
		o.Ops = append(o.Ops, &Operation{
			StrVal:       str,
			Type:         Insert,
			BaseLength:   int64(len(str)),
			TargetLength: int64(len(str)),
		})
	}
	return o, nil
}

func (o *Operation) delete(n string) (*Operation, error) {
	if o == nil {
		return nil, fmt.Errorf("current operation is nil")
	}
	val, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		return o, err
	}
	if val == 0 {
		return o, nil
	}

	if val > 0 {
		fmt.Println("删除操作必须是负数", val)
		val = -val
	}
	o.BaseLength -= val //把baseLength加上删除长度，
	if len(o.Ops) > 0 && o.Ops[len(o.Ops)-1].IsDelete() {
		o.Ops[len(o.Ops)-1].Val += val
	} else {
		o.Ops = append(o.Ops, &Operation{
			BaseLength:   -val,
			TargetLength: -val,
			Val:          val,
			Type:         Delete,
		})
	}
	return o, nil
}

func (o *Operation) isNoop() bool {
	return o == nil || len(o.Ops) == 0 || (len(o.Ops) == 1 && o.Ops[0].IsRetain())
}

func (o *Operation) toString() string {
	// map: build a new array by applying a function to every element in an old
	// array.
	toStr := func(op *Operation) string {
		if op == nil {
			fmt.Println("invalid op")
			return ""
		}

		if op.IsRetain() {
			return "retain " + strconv.FormatInt(op.Val, 10)
		} else if op.IsInsert() {
			return "insert \"" + op.StrVal + "\""
		} else if op.IsDelete() {
			return "delete " + strconv.FormatInt(-op.Val, 10)
		}
		fmt.Println("unknown op type")
		return ""
	}

	var newArr []string
	for i := 0; i < len(o.Ops); i++ {
		newArr[i] = toStr(o.Ops[i])
	}
	return strings.Join(newArr, ",")
}

func (o *Operation) toJSON() string {
	data, _ := json.Marshal(o.Ops)
	return string(data)
}

func (o *Operation) fromJSON(ops string) error {
	o1 := &Operation{}
	err := json.Unmarshal([]byte(ops), o1)
	if err != nil {
		return err
	}

	for i := 0; i < len(o1.Ops); i++ {
		var op = o1.Ops[i]
		if op.IsRetain() {
			o.retain(op.StrVal)
		} else if op.IsInsert() {
			o.insert(op.StrVal)
		} else if op.IsDelete() {
			o.delete(op.StrVal)
		} else {
			data, _ := json.Marshal(op)
			fmt.Println("unknown operation: " + string(data))
		}
	}
	return nil
}

// Apply an operation to a string, returning a new string. Throws an error if
// there's a mismatch between the input string and the operation.
//把操作应用到str上产生一个新的操作
func (o *Operation) apply(str string) (string, error) {
	var operation = o
	if int64(len(str)) != operation.BaseLength {
		return "", fmt.Errorf("The operation's base length must be equal to the string's length.")
	}

	var newStr []string
	var j, strIndex int
	for i := 0; i < len(o.Ops); i++ {
		if o.Ops[i].IsRetain() {
			if strIndex+int(o.Ops[i].Val) > len(str) {
				return "", fmt.Errorf("Operation can't retain more characters than are left in the string.")
			}
			// Copy skipped part of the old string.
			newStr[j] = str[strIndex : strIndex+int(o.Ops[i].Val)]
			j++
			strIndex += int(o.Ops[i].Val)
		} else if o.Ops[i].IsInsert() {
			// Insert string.
			newStr[j] = o.Ops[i].StrVal
			j++
		} else { // delete op
			strIndex -= int(o.Ops[i].Val)
		}
	}
	if strIndex != len(str) {
		return "", fmt.Errorf("The operation didn't operate on the whole string.")
	}
	return strings.Join(newStr, ""), nil
}

// Computes the inverse of an operation. The inverse of an operation is the
// operation that reverts the effects of the operation, e.g. when you have an
// operation 'insert("hello "); skip(6);' then the inverse is 'delete("hello ");
// skip(6);'. The inverse should be used for implementing undo.
//反向操作是用来做undo的
func (o *Operation) invert(str string) *Operation {
	var strIndex int
	inverse := &Operation{}

	for i := 0; i < len(o.Ops); i++ {
		op := o.Ops[i]
		if op.IsRetain() {
			inverse.retain(strconv.FormatInt(op.Val, 10))
			strIndex += int(op.Val)
		} else if op.IsInsert() {
			inverse.delete(strconv.FormatInt(int64(len(str)), 10))
		} else { // delete op
			inverse.insert(str[strIndex : strIndex-int(op.Val)])
			strIndex -= int(op.Val)
		}
	}
	return inverse
}
