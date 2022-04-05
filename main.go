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
	fmt.Println("A':", string(ABData), "\nB':", string(BAData))

	fmt.Println("======================")
	// A.compose(AB) 和 B.compose(BA) 都为：
	A1 := &ot.QuillDelta{
		Ops: []ot.Op{
			&ot.RetainOp{
				Retain: 5,
			},
			&ot.InsertOp{
				Insert: ",",
			},
		},
	}
	B1 := &ot.QuillDelta{
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
	db, _ := doc.Compose(B1)
	BBA, err := db.Compose(BA)
	if err != nil {
		fmt.Println(err)
	}

	da, _ := doc.Compose(A1)
	AAB, err := da.Compose(AB)
	if err != nil {
		fmt.Println(err)
	}

	AABData, _ := json.Marshal(AAB)
	BBAData, _ := json.Marshal(BBA)

	fmt.Println("AAB:", string(AABData), "\nBBA:", string(BBAData))

	fmt.Println("==================")
	// apply(apply(S, A), B) = apply(S, compose(A, B)) must hold.
	sa, err := A.Apply("hello world")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("sa:", sa)
	saa, err := B.Apply(sa)
	if err != nil {
		fmt.Println(err)
	}

	cab, err := A.Compose(B)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("cab:", cab)
	sb, err := cab.Apply("hello world")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(" apply(apply(S, A), B) :", saa, "\n apply(S, compose(A, B)):", sb)
}
