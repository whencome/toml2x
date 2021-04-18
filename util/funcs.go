package util

import (
    "bytes"
    "regexp"
)

// RuneInArray 判断给定的rune是否在数组中
func RuneInArray(r rune, arr []rune) bool {
    size := len(arr)
    if size == 0 {
        return false
    }
    for i := 0; i < size; i++ {
        if r == arr[i] {
            return true
        }
    }
    return false
}

// RunesContains 判断rune列表中是否包含某个值
func RunesContains(arr []rune, r rune) bool {
    for _, v := range arr {
        if v == r {
            return true
        }
    }
    return false
}

// IsNumeric 判断给定的字符串是否是数字
func IsNumeric(str string) bool {
    matched, err := regexp.MatchString(`^(\+|\-)?(0|[1-9]\d*)((\.\d+)?((e|E)(\+|\-)?[1-9]\d*)?)?$`, str)
    if err != nil {
        return false
    }
    return matched
}

// IsPositiveIntNumeric 判断给定的数字是否是正整数
func IsPositiveIntNumeric(str string) bool {
    matched, err := regexp.MatchString(`^(0|[1-9]\d*)$`, str)
    if err != nil {
        return false
    }
    return matched
}

// IsNonNegativeInt 判断给定的数字是否是非负整数
func IsNonNegativeInt(str string) bool {
    matched, err := regexp.MatchString(`^(0|[1-9]\d*)$`, str)
    if err != nil {
        return false
    }
    return matched
}

// ParseTomlTableName Parses TOML table names and returns the hierarchy array of table names.
func ParseTomlTableName(chars []rune) []string {
    buffer := bytes.Buffer{}
    strOpen := false
    names := make([]string, 0)

    charsSize := len(chars)
    for i := 0; i < charsSize; i++ {
        if chars[i] == '"' {
            if !strOpen || (strOpen && chars[i-1] != '\\') {
                strOpen = !strOpen
            }
        } else if chars[i] == '.' && !strOpen {
            names = append(names, buffer.String())
            buffer.Reset()
            continue
        }
        buffer.WriteRune(chars[i])
    }
    if buffer.Len() > 0 {
        names = append(names, buffer.String())
    }
    return names
}
