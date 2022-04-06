package main

import (
	"encoding/json"
	"fmt"
	"ot-go/ot"
)

func main() {
	doc := &ot.QuillDelta{
		Ops: []ot.Op{
			&ot.InsertOp{
				Insert: "Hello ",
				Attributes: map[string]interface{}{
					"alt": "Lab Octocat",
				},
			},
			&ot.InsertOp{
				Insert: "World",
				Attributes: map[string]interface{}{
					"bold": true,
				},
			},
		},
	}

	change := &ot.QuillDelta{
		Ops: []ot.Op{
			&ot.RetainOp{
				Retain: 6,
			},
			&ot.DeleteOp{
				Delete: -3,
			},
			&ot.InsertOp{
				Insert: "xiazemin",
			},
			&ot.DeleteOp{
				Delete: -2,
			},
		},
	}
	docData, err := json.Marshal(doc)
	if err != nil {
		fmt.Println(err)
	}
	changeData, err := json.Marshal(change)
	if err != nil {
		fmt.Println(err)
	}
	updatedDoc, err := doc.Compose(change)
	if err != nil {
		fmt.Println(err)
	}
	updatedData, err := json.Marshal(updatedDoc)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("doc:", string(docData), "\nchange:", string(changeData), "\nupdated", string(updatedData))

	// 而 transform 则用于操作变基。如对一篇文档，A 先做了修改并提交到服务器，而 B 也在同一时刻对文档做了修改并提交到服务器。此时服务器先收到 A，后收到 B，且 A 和 B 都是对同一版本做的修改。为了合并这两个操作，需要变换 B 为 B' 使得 A.compose(B') === B.compose(A')，这个变换过程就是通过 transform 实现的，即 A.compose(A.transform(B)) === B.compose(B.transform(A))。
	A := &ot.QuillDelta{
		Ops: []ot.Op{
			&ot.RetainOp{
				Retain: 5,
			},
			&ot.InsertOp{
				Insert: ",",
			},
		},
	}
	B := &ot.QuillDelta{
		Ops: []ot.Op{
			&ot.RetainOp{
				Retain: 6,
			},
			&ot.InsertOp{
				Insert: "xiazemin",
			},
			&ot.DeleteOp{
				Delete: -5,
			},
		},
	}

	BA, AB, err := doc.Transform(A, B)
	if err != nil {
		fmt.Println(err)
	}
	ABData, _ := json.Marshal(AB)
	BAData, _ := json.Marshal(BA)
	AData, _ := json.Marshal(A)
	BData, _ := json.Marshal(B)
	fmt.Println("A", string(AData), "\nAB:", string(ABData), "\nB", string(BData), "\nBA:", string(BAData))

	// A.compose(AB) 和 B.compose(BA) 都为：
	// new Delta([
	// 	new TextOperator({ action: 'retain', data: 5 }),
	// 	new TextOperator({ action: 'insert', data: ',' }),
	// 	new TextOperator({ action: 'retain', data: 1 }),
	// 	new TextOperator({ action: 'insert', data: 'Tom' }),
	// 	new TextOperator({ action: 'remove', data: 5 })
	//   ]);
	// AB.(*ot.QuillDelta).TargetLength = 0
	AB.(*ot.QuillDelta).BaseLength = 0
	// BA.(*ot.QuillDelta).TargetLength = 0
	BA.(*ot.QuillDelta).BaseLength = 0

	AAB, err := A.Compose(AB)
	if err != nil {
		fmt.Println(err)
	}

	BBA, err := B.Compose(BA)
	if err != nil {
		fmt.Println(err)
	}

	AABData, _ := json.Marshal(AAB)
	BBAData, _ := json.Marshal(BBA)
	fmt.Println("AAB:", string(AABData), "\nBBA:", string(BBAData))

	//https://npmmirror.com/package/richdoc/v/0.3.7
	check(doc, A, B, BA.(*ot.QuillDelta), AB.(*ot.QuillDelta), AAB.(*ot.QuillDelta), BBA.(*ot.QuillDelta))
}

func check(doc, A, B, BA, AB, AAB, BBA *ot.QuillDelta) {
	fmt.Println("======================")
	// apply(apply(S, A), B) = apply(S, compose(A, B)) must hold.
	str := "hello world"
	fmt.Println(A.BaseLength)
	A.BaseLength = len(str)
	sa, err := A.Apply(str)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("sa:", sa, len(sa))
	B.BaseLength = len(sa)
	saa, err := B.Apply(sa)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("saa:", saa)

	AAB.BaseLength = len(str)
	sb, err := AAB.Apply(str)
	if err != nil {
		fmt.Println(err)
	}

	BBA.BaseLength = len(str)
	sb1, err := BBA.Apply(str)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(" apply(apply(S, A), B) :", saa, "\n apply(S, compose(A, B)):", sb, "\n apply(S, compose(A, B)):", sb1)

}
