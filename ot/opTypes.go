package ot

import "strconv"

type Op interface {
	IsRetain() bool
	IsInsert() bool
	IsDelete() bool
	GetVal() int
	SetVal(interface{})
	GetStringVal() string
}

// Keep and bold the next 5 characters
//    { retain: 5, attributes: { bold: true } }
// Keep and unbold the next 5 characters
// More specifically, remove the bold key in the attributes Object
// in the next 5 characters
//    { retain: 5, attributes: { bold: null } }
type RetainOp struct {
	Retain     int               `json:"retain"`
	Attributes map[string]string `json:"attributes"`
}

func (r *RetainOp) IsRetain() bool {
	return true
}
func (r *RetainOp) IsInsert() bool {
	return false
}
func (r *RetainOp) IsDelete() bool {
	return false
}
func (r *RetainOp) GetVal() int {
	return r.Retain
}
func (r *RetainOp) SetVal(v interface{}) {
	r.Retain = v.(int)
}

func (r *RetainOp) GetStringVal() string {
	return strconv.FormatInt(int64(r.Retain), 10)
}

// Insert a link
//    { insert: "Google", attributes: { href: 'https://www.google.com' } }

// Insert an embed
//    {
//    insert: { image: 'https://octodex.github.com/images/labtocat.png' },
//    attributes: { alt: "Lab Octocat" }
//    }
type InsertOp struct {
	Insert     interface{}            `json:"insert"`
	Attributes map[string]interface{} `json:"attributes"`
}

func (r *InsertOp) IsRetain() bool {
	return false
}
func (r *InsertOp) IsInsert() bool {
	return true
}
func (r *InsertOp) IsDelete() bool {
	return false
}
func (r *InsertOp) GetVal() int {
	val, ok := r.Insert.(string)
	if ok {
		return len(val)
	}
	return 0
}
func (r *InsertOp) SetVal(v interface{}) {
	r.Insert = v.(string)
}
func (r *InsertOp) GetStringVal() string {
	val, ok := r.Insert.(string)
	if ok {
		return val
	}
	return ""
}

// Delete the next 10 characters
//    { delete: 10 }
type DeleteOp struct {
	Delete int `json:"delete"`
}

func (r *DeleteOp) IsRetain() bool {
	return false
}
func (r *DeleteOp) IsInsert() bool {
	return false
}
func (r *DeleteOp) IsDelete() bool {
	return true
}
func (r *DeleteOp) GetVal() int {
	return r.Delete
}
func (r *DeleteOp) SetVal(v interface{}) {
	r.Delete = v.(int)
}
func (r *DeleteOp) GetStringVal() string {
	return strconv.FormatInt(int64(r.Delete), 10)
}
