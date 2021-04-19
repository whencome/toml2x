/**
 * define parse target object structure.
 */
package xtype

import (
    "sort"

    "github.com/whencome/toml2x/util"
)

// define common data type
const (
    TypeNumber = iota
    TypeBoolean
    TypeString
    TypeMap    // key-value 值
    TypeArray  // array
)

// Object define a scalar object which save only a single value
type Object struct {
    Value interface{}
    Type  int // indicate the value type, can be int,float,string,bool,array or map
}

// Define the key of a map
type Key struct {
    Value     string
    IsNumeric bool // 键值是否是数字
}

func (k *Key) String() string {
    if k.IsNumeric {
        return k.Value
    }
    return util.FmtString(k.Value)
}

// Define a map struct
type Map struct {
    Keys []*Key
    Data map[*Key]*Object
}

// Define an array
type Array struct {
    Keys []int
    Data map[int]*Object
}

type KeyValuePair struct {
    Key   string
    Value *Object
}

// NewBoolObject create a boolean object
func NewBoolObject(val string) *Object {
    return &Object{
        Value: val,
        Type:  TypeBoolean,
    }
}

// NewNumberObject create a number object
func NewNumberObject(val string) *Object {
    return &Object{
        Value: val,
        Type:  TypeNumber,
    }
}

// NewStringObject create a string object
func NewStringObject(val string) *Object {
    return &Object{
        Value: val,
        Type:  TypeString,
    }
}

func NewMapObject(val *Map) *Object {
    return &Object{
        Value: val,
        Type:  TypeMap,
    }
}

func NewArrayObject(val *Array) *Object {
    return &Object{
        Value: val,
        Type:  TypeArray,
    }
}

func NewNumberKey(v string) *Key {
    return &Key{
        Value:     v,
        IsNumeric: true,
    }
}

func NewStringKey(v string) *Key {
    return &Key{
        Value:     v,
        IsNumeric: false,
    }
}

func NewMap() *Map {
    return &Map{
        Keys: make([]*Key, 0),
        Data: make(map[*Key]*Object, 0),
    }
}

// Create a new empty array
func NewArray() *Array {
    return &Array{
        Keys: make([]int, 0),
        Data: make(map[int]*Object, 0),
    }
}

func (arr *Array) Add(i int, obj *Object) {
    if _, ok := arr.Data[i]; !ok {
        arr.Keys = append(arr.Keys, i)
    }
    arr.Data[i] = obj
}

// 判断给定的key是否存在
func (m *Map) GetKey(k string) *Key {
    if len(m.Keys) == 0 {
        return nil
    }
    for _, ek := range m.Keys {
        if ek.Value == k {
            return ek
        }
    }
    return nil
}

func (m *Map) Add(k *Key, obj *Object) {
    if _, ok := m.Data[k]; !ok {
        m.Keys = append(m.Keys, k)
    }
    m.Data[k] = obj
}

func (m *Map) DeepAdd(fields []string, obj *Object) {
    fSize := len(fields)
    if fSize <= 0 {
        return
    }
    dst := m
    for i, field := range fields {
        k := dst.GetKey(field)
        if k == nil {
            k = NewStringKey(field)
            dst.Keys = append(dst.Keys, k)
        }
        if dst.Data == nil {
            dst.Data = make(map[*Key]*Object, 0)
        }
        if i == fSize-1 {
            dst.Data[k] = obj
            return
        }
        if kv, ok := dst.Data[k]; ok {
            if kv.Type == TypeMap {
                dst = dst.Data[k].Value.(*Map)
            }
        } else {
            innerMap := NewMap()
            dst.Data[k] = NewMapObject(innerMap)
            dst = innerMap
        }
    }
}

// Merge 合并对象
func (m *Map) Merge(m1 *Map) {
    if m1 == nil {
        return
    }
    dst := m
    for _, k := range m1.Keys {
        if _, ok := dst.Data[k]; !ok {
            dst.Keys = append(dst.Keys, k)
            dst.Data[k] = NewMapObject(NewMap())
        }
        if m1.Data[k].Type == TypeMap {
            dst.Data[k].Value.(*Map).Merge(m1.Data[k].Value.(*Map))
        } else {
            dst.Data[k] = m1.Data[k]
        }
    }
}

// 判断map对象是数组还是对象
func (m *Map) IsArray() bool {
    if len(m.Keys) == 0 {
        return false
    }
    numberKeys := make([]int, 0)
    for _, k := range m.Keys {
        // 数组的key必须是数字
        if !k.IsNumeric {
            return false
        }
        // 必须是非负整数
        if !util.IsNonNegativeInt(k.Value) {
            return false
        }
        // 转换成整数
        ik := util.Int(k.Value)
        numberKeys = append(numberKeys, ik)
    }
    // 对下标进行排序
    sort.Ints(numberKeys[:])
    // 检查下标是否是从0开始逐个递增
    for i, v := range numberKeys {
        if i != v {
            return false
        }
    }
    return true
}
