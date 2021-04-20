package xtype

import (
    "bytes"
    "strings"

    "github.com/whencome/toml2x/formatter"
    "github.com/whencome/toml2x/util"
)

// 将对象转换为字符串
func (o *Object) Json() string {
    if o == nil {
        return "null"
    }
    switch o.Type {
    case TypeBoolean:
        return util.String(o.Value)
    case TypeNumber:
        v := util.String(o.Value)
        if strings.HasPrefix(v, "+") {
            v = string([]rune(v)[1:])
        }
        return v
    case TypeString:
        return formatter.FmtJsonString(util.String(o.Value))
    case TypeMap:
        return o.Value.(*Map).Json()
    }
    return "\"\""
}

// 将对象转换为字符串
func (o *Object) Xml() string {
    if o == nil {
        return "<xml><single>null</single></xml>"
    }
    switch o.Type {
    case TypeBoolean:
        return "<xml><single>" + util.String(o.Value) + "</single></xml>"
    case TypeNumber:
        v := util.String(o.Value)
        if strings.HasPrefix(v, "+") {
            v = string([]rune(v)[1:])
        }
        return "<xml><single>" + v + "</single></xml>"
    case TypeString:
        return "<xml><single><![CDATA[" + util.String(o.Value)  + "]]></single></xml>"
    case TypeMap:
        return "<xml><table>" + o.Value.(*Map).Xml() + "</table></xml>"
    }
    return "<xml><single>null</single></xml>"
}

// 将对象转换为php
func (o *Object) Php() string {
    if o == nil {
        return "''"
    }
    switch o.Type {
    case TypeBoolean:
        return util.String(o.Value)
    case TypeNumber:
        v := util.String(o.Value)
        if strings.HasPrefix(v, "+") {
            v = string([]rune(v)[1:])
        }
        return v
    case TypeString:
        return formatter.FmtPhpString(util.String(o.Value))
    case TypeMap:
        return o.Value.(*Map).Php(0)
    }
    return "''"
}

// 将map转换为json
func (m *Map) Json() string {
    if m.IsArray() {
        return m.jsonArray()
    }
    return m.jsonObject()
}

func (m *Map) jsonArray() string {
    buf := bytes.Buffer{}
    buf.WriteString("[")
    for i, k := range m.Keys {
        if i > 0 {
            buf.WriteString(",")
        }
        v := m.Data[k]
        buf.WriteString(v.Json())
    }
    buf.WriteString("]")
    return buf.String()
}

func (m *Map) jsonObject() string {
    buf := bytes.Buffer{}
    buf.WriteString("{")
    for i, k := range m.Keys {
        if i > 0 {
            buf.WriteString(",")
        }
        v := m.Data[k]
        buf.WriteString(formatter.FmtJsonKey(k.Value))
        buf.WriteString(":")
        buf.WriteString(v.Json())
    }
    buf.WriteString("}")
    return buf.String()
}

// 将map转换为xml
func (m *Map) Xml() string {
    buf := bytes.Buffer{}
    isArr := m.IsArray()
    for _, k := range m.Keys {
        v := m.Data[k]
        if isArr {
            buf.WriteString("<item>")
        } else {
            buf.WriteString("<" + k.Value + ">")
        }
        switch v.Type {
        case TypeBoolean:
            fallthrough
        case TypeString:
            buf.WriteString("<![CDATA[")
            buf.WriteString(util.String(v.Value))
            buf.WriteString("]]>")
        case TypeNumber:
            n := util.String(v)
            if strings.HasPrefix(n, "+") {
                n = string([]rune(n)[1:])
            }
            buf.WriteString("<![CDATA[")
            buf.WriteString(util.String(v.Value))
            buf.WriteString("]]>")
        case TypeMap:
            buf.WriteString(v.Value.(*Map).Xml())
        }
        if isArr {
            buf.WriteString("</item>")
        } else {
            buf.WriteString("</" + k.Value + ">")
        }
    }
    return buf.String()
}

// 将map转换为php数组
func (m *Map) Php(depth int) string {
    if depth < 0 {
        depth = 0
    }
    indent := "    "
    buf := bytes.Buffer{}
    buf.WriteString("array(\n")
    isArr := m.IsArray()
    for _, k := range m.Keys {
        v := m.Data[k]
        buf.WriteString(strings.Repeat(indent, depth+1))
        if isArr {
            buf.WriteString(k.Value)
        } else {
            buf.WriteString(formatter.FmtPhpKey(k.Value))
        }
        buf.WriteString(" => ")
        switch v.Type {
        case TypeBoolean:
            fallthrough
        case TypeString:
            buf.WriteString(formatter.FmtPhpString(util.String(v.Value)))
            buf.WriteString(",\n")
        case TypeNumber:
            n := util.String(v.Value)
            if strings.HasPrefix(n, "+") {
                n = string([]rune(n)[1:])
            }
            buf.WriteString(n)
            buf.WriteString(",\n")
        case TypeMap:
            buf.WriteString(v.Value.(*Map).Php(depth+1))
        }
    }
    buf.WriteString(strings.Repeat(indent, depth))
    buf.WriteString(")")
    if depth > 0 {
        buf.WriteString(",\n")
    } else {
        buf.WriteString("\n")
    }
    return buf.String()
}
