package formatter

import (
    "bytes"
    "encoding/json"
    "github.com/whencome/toml2x/util"
    "strings"
)

// FmtString 通用字符串处理，对字符串使用引号包围
func FmtString(str string) string {
    str = strings.TrimSpace(str)
    if str == "" {
        return "\"\""
    }
    chars := []rune(str)
    charsSize := len(chars)
    buffer := bytes.Buffer{}
    buffer.WriteRune('"')
    for i := 0; i < charsSize; i++ {
        if chars[i] == '"' && ((i == 0) || (i > 0 && chars[i-1] != '\\')) {
            buffer.WriteRune('\\')
            buffer.WriteRune('"')
            continue
        }
        if chars[i] == '\\' {
            if i+1 < charsSize {
                // 引号
                if chars[i+1] == '"' {
                    buffer.WriteRune('\\')
                    buffer.WriteRune('\\')
                    buffer.WriteRune('"')
                    i += 1
                    continue
                }
            }
        }
        if chars[i] == '\n' {
            buffer.WriteRune('\n')
            continue
        }
        if (i == 0 && chars[i] == '\\' && (i+1 == charsSize || chars[i+1] == '\n')) ||
            (i > 0 && chars[i] == '\\' && chars[i-1] != '\\' && (i+1 == charsSize || chars[i+1] == '\n')) {
            continue
        }
        buffer.WriteRune(chars[i])
    }
    buffer.WriteRune('"')
    return buffer.String()
}

// FmtPhpString 格式化为PHP字符串形式
func FmtPhpString(str string) string {
    str = strings.TrimSpace(str)
    if str == "" {
        return "''"
    }
    chars := []rune(str)
    charsSize := len(chars)
    buffer := bytes.Buffer{}
    buffer.WriteRune('\'')
    for i := 0; i < charsSize; i++ {
        if chars[i] == '\'' && ((i == 0) || (i > 0 && chars[i-1] != '\\')) {
            buffer.WriteRune('\\')
            buffer.WriteRune('\'')
            continue
        }
        // 去除转义
        if chars[i] == '\\' {
            if i+1 < charsSize {
                // 引号
                if chars[i+1] == '\'' {
                    buffer.WriteRune('\\')
                    buffer.WriteRune('\'')
                    i += 1
                    continue
                }
            }
        }
        if (i == 0 && chars[i] == '\\' && (i+1 == charsSize || chars[i+1] == '\n')) ||
            (i > 0 && chars[i] == '\\' && chars[i-1] != '\\' && (i+1 == charsSize || chars[i+1] == '\n')) {
            continue
        }
        buffer.WriteRune(chars[i])
    }
    buffer.WriteRune('\'')
    return buffer.String()
}

func FmtPhpKey(k string) string {
    if util.IsPositiveIntNumeric(k) {
        return k
    }
    return FmtPhpString(k)
}

func FmtJsonString(v interface{}) string {
    fmtRs, err := json.Marshal(v)
    if err != nil {
        return "null"
    }
    return string(fmtRs)
}
