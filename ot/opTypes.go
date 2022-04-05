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
	if r == nil {
		return false
	}
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

func (i *InsertOp) IsRetain() bool {
	return false
}
func (i *InsertOp) IsInsert() bool {
	if i == nil {
		return false
	}
	return true
}
func (i *InsertOp) IsDelete() bool {
	return false
}
func (i *InsertOp) GetVal() int {
	val, ok := i.Insert.(string)
	if ok {
		return len(val)
	}
	return 0
}
func (i *InsertOp) SetVal(v interface{}) {
	i.Insert = v.(string)
}
func (i *InsertOp) GetStringVal() string {
	val, ok := i.Insert.(string)
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

func (d *DeleteOp) IsRetain() bool {
	return false
}
func (d *DeleteOp) IsInsert() bool {
	return false
}
func (d *DeleteOp) IsDelete() bool {
	if d == nil {
		return false
	}
	return true
}
func (d *DeleteOp) GetVal() int {
	return d.Delete
}
func (d *DeleteOp) SetVal(v interface{}) {
	d.Delete = v.(int)
}
func (d *DeleteOp) GetStringVal() string {
	return strconv.FormatInt(int64(d.Delete), 10)
}
