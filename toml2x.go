package toml2x

import (
    "github.com/whencome/toml2x/formatter"
    "github.com/whencome/toml2x/parser"
    "github.com/whencome/toml2x/xtype"
)

// parse 解析toml配置内容
// toml toml格式的配置内容
func parse(dataType string, toml string) (*xtype.Object, error) {
    toml, err := formatter.Normalize(toml)
    if err != nil {
        return nil, err
    }
    obj, err := parser.Parse(dataType, toml)
    if err != nil {
        return nil, err
    }
    return obj, nil
}

// Json 转换为json
// dataType 配置的数据类型，single，table
// toml toml配置内容
func Json(dataType string, toml string) (string, error) {
    obj, err := parse(dataType, toml)
    if err != nil {
        return "", err
    }
    return obj.Json(true), nil
}

// xml 转换为xml格式
// dataType 配置的数据类型，single，table
// toml toml配置内容
func Xml(dataType string, toml string) (string, error) {
    obj, err := parse(dataType, toml)
    if err != nil {
        return "", err
    }
    return obj.Xml(), nil
}

// xml 转换为php格式
// dataType 配置的数据类型，single，table
// toml toml配置内容
func Php(dataType string, toml string) (string, error) {
    obj, err := parse(dataType, toml)
    if err != nil {
        return "", err
    }
    return obj.Php(), nil
}
