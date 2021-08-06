const _ = require('lodash')

let oldV = [1, 2]
let newV = [12]
if (!_.isEqual(oldV, newV)) {
    console.log("case1 has diff", oldV, newV)
}

let oldV0 = {
    a: 'a',
    b: 'b',
    c: {
        a: 'a',
        b: 'b'
    }
}
var newV0 = {
    a: 'a'
}
newV0.b = 'b'
newV0["c"] = {
    a: 'a',
    b: 'b'
}
if (!_.isEqual(oldV0, newV0)) {
    console.log("case2 has diff", oldV0, newV0)
}

oldV = {
    a: 'a',
    b: 'b'
}
newV = {
    a: 'a',
    b: 'b'
}
if (!_.isEqual(oldV, newV)) {
    console.log("case3 has diff", oldV, newV)
}

oldV = {
    a: 'a',
    b: 'b'
}
newV = {
    a: 'a'
}
newV['b'] = 'b'
if (!_.isEqual(oldV, newV)) {
    console.log("case4 has diff", oldV, newV)
}

const newV1 = {
    b: {
        a: 'a',
        b: 'b'
    },
    a: {
        a: 'a',
        b: 'b'
    },
    d: null
}

var oldV1 = {
    a: {
        a: 'a',
        b: 'b'
    },
    b: {
        a: 'a',
        b: 'b'
    },
    c: newV1.c,
    d: null
}

newV1['a'] = {
    b: 'b',
    a: 'a'
}
newV1.b = {
    a: 'a',
    b: 'b'
}

if (!_.isEqual(oldV1, newV1)) {
    console.log("case3 has  diff", oldV1, newV1)
    console.log("case3 has  diff", JSON.stringify(oldV1), JSON.stringify(newV1)) //
}

/*
case3 has  diff { a: { a: 'a', b: 'b' },
  b: { a: 'a', b: 'b' },
  c: undefined,
  d: null } { b: { a: 'a', b: 'b' }, a: { b: 'b', a: 'a' }, d: null }
case3 has  diff {"a":{"a":"a","b":"b"},"b":{"a":"a","b":"b"},"d":null} {"b":{"a":"a","b":"b"},"a":{"b":"b","a":"a"},"d":null}
*/