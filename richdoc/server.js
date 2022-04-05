const {Delta,TextOperator}=require("richdoc")
const { pack, unpack,TableOperator,CellOperator,Operator} =require('richdoc')

const doc = new Delta([
    new TextOperator({ action: 'insert', data: 'Hello ' }),
    new TextOperator({ action: 'insert', data: 'World', attributes: { bold: true } })
  ]);

  console.log(doc)

  const change = new Delta([
    new TextOperator({ action: 'retain', data: 6 }),
    new TextOperator({ action: 'insert', data: 'Tom' }),
    new TextOperator({ action: 'remove', data: 5 })
  ]);
  
  const updatedDoc = doc.compose(change);
  
  console.log(updatedDoc)

  //A.compose(B') === B.compose(A')，这个变换过程就是通过 transform 实现的，即 A.compose(A.transform(B)) === B.compose(B.transform(A))。

  const A = new Delta([
    new TextOperator({ action: 'retain', data: 5 }),
    new TextOperator({ action: 'insert', data: ',' })
  ]);

  const B = new Delta([
    new TextOperator({ action: 'retain', data: 6 }),
    new TextOperator({ action: 'insert', data: 'Tom' }),
    new TextOperator({ action: 'remove', data: 5 })
  ]);

  console.log(A.transform(B))

  const AB = new Delta([
    new TextOperator({ action: 'retain', data: 7 }),
    new TextOperator({ action: 'insert', data: 'Tom' }),
    new TextOperator({ action: 'remove', data: 5 })
  ]);

  console.log(B.transform(A))

  const BA = new Delta([
    new TextOperator({ action: 'retain', data: 5 }),
    new TextOperator({ action: 'insert', data: ',' })
  ]);

  console.log(A.compose(AB))
  console.log(B.compose(BA))
  



const packed = pack(doc);
const unpacked = unpack(packed);

console.log(packed,unpacked)


//这里定义两个伪函数：T(A, B) = A.transform(B), A + B = A.compose(B)。可以得出：A + T(A, B) = B + T(B, A)。

const table = new Delta([
    new TableOperator({
      action: 'insert',
      data: {
        rows: new Delta([
          new Operator({ action: 'insert', data: 3 })
        ]),
        cols: new Delta([
          new Operator({ action: 'insert', data: 1 }),
          new Operator({ action: 'insert', data: 1, attributes: { width: 50 } }),
          new Operator({ action: 'insert', data: 1 })
        ]),
        cells: {
          A1: new CellOperator({
            action: 'insert',
            data: new Delta([
              new TextOperator({ action: 'insert', data: 'Hello ' }),
              new TextOperator({ action: 'insert', data: 'World', attributes: { bold: true } })
            ])
          })
        }
      }
    })
  ]);

  console.log(table)