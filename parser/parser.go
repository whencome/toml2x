package parser

import (
    "bytes"
    "errors"
    "strconv"
    "strings"

    "github.com/whencome/toml2x/util"
    "github.com/whencome/toml2x/xtype"
)

// Parse 解析toml内容
func Parse(contentType string, toml string) (*xtype.Object, error) {
    if contentType == "single" {
        return ParseSingle(toml)
    }
    return ParseTable(toml)
}

// ParseTable 解析复杂数据
func ParseTable(toml string) (*xtype.Object, error) {
    arr := &xtype.Map{}

    // split lines
    arrToml := strings.Split(toml, "\n")
    arrSize := len(arrToml)

    var recurseKeys []string
    for ln := 0; ln < arrSize; ln++ {
        line := []rune(strings.TrimSpace(arrToml[ln]))
        lineSize := len(line)

        // Skip commented and empty lines
        if lineSize == 0 || line[0] == '#' {
            continue
        }

        // Array of Tables
        if string(line[0:2]) == "[[" && string(line[lineSize-2:]) == "]]" {
            tableName := line[2 : lineSize-2]
            aTables := parseTomlTableName(tableName)
            if len(aTables) <= 0 {
                continue
            }
            recurseKeys = arr.GetRecursiveIndexedKeys(aTables)
        } else if string(line[0:1]) == "[" && string(line[lineSize-1:]) == "]" {
            tableName := line[1 : lineSize-1]
            aTables := parseTomlTableName(tableName)
            if len(aTables) <= 0 {
                continue
            }
            recurseKeys = make([]string, len(aTables))
            copy(recurseKeys, aTables)
        } else if util.RunesContains(line, '=') {
            rawLine := string(line)
            pos := strings.Index(rawLine, "=")
            field := strings.TrimSpace(rawLine[0:pos])
            val := strings.TrimSpace(rawLine[pos+1:])
            valSize := len(val)
            if valSize >= 3 && val[0:3] == `"""` {
                if valSize == 3 || (valSize > 3 && val[valSize-3:] != `"""`) {
                    for {
                        ln++
                        nextLine := strings.TrimSpace(arrToml[ln])
                        val += "\n"
                        val += arrToml[ln]
                        if nextLine == `"""` || (len(nextLine) > 3 && nextLine[len(nextLine)-3:] == `"""`) {
                            break
                        }
                    }
                }
            }
            if valSize >= 3 && val[0:3] == `'''` {
                if valSize == 3 || (valSize > 3 && val[valSize-3:] != `'''`) {
                    for {
                        ln++
                        nextLine := strings.TrimSpace(arrToml[ln])
                        val += "\n"
                        val += arrToml[ln]
                        if nextLine == `'''` || (len(nextLine) > 3 && nextLine[len(nextLine)-3:] == `'''`) {
                            break
                        }
                    }
                }
            }
            // 支持单引号，对单引号的内容进行一次转义
            if valSize > 2 && val[0] == '\'' && val[1] != '\'' && val[valSize-1] == '\'' {
                buf := bytes.Buffer{}
                valChars := []rune(val)
                valCharsSize := len(valChars)
                for j := 0; j < valCharsSize; j++ {
                    // 单引号里面的内容，不对双引号再进行转义
                    if util.RuneInArray(valChars[j], []rune{'\\'}) {
                        buf.WriteRune('\\')
                        buf.WriteRune(valChars[j])
                    } else {
                        buf.WriteRune(valChars[j])
                    }
                }
                val = buf.String()
            }
            pathKeys := make([]string, len(recurseKeys))
            copy(pathKeys, recurseKeys)
            fieldKeys := parseTomlTableName([]rune(field))
            pathKeys = append(pathKeys, fieldKeys...)
            err := parseKeyValue(arr, pathKeys, val)
            if err != nil {
                return nil, err
            }
        } else if string(line[0:1]) == "[" && string(line[lineSize-1:]) != "]" {
            return nil, errors.New("Key groups have to be on a line by themselves: " + string(line))
        } else {
            return nil, errors.New("Syntax error on: " + string(line))
        }
    }

    return xtype.NewMapObject(arr), nil
}

// ParseSingle 解析单个值
func ParseSingle(val string) (*xtype.Object, error) {
    val = strings.TrimSpace(val)
    // 布尔值
    if val == "true" || val == "false" {
        return xtype.NewBoolObject(val), nil
    }
    // 字符串
    if util.IsNumeric(val) {
        return xtype.NewNumberObject(val), nil
    }
    // 字符串解析，可能是复杂对象
    chars := []rune(val)
    charsSize := len(chars)
    parsedVal := make([]rune, 0)
    // 多行字符串
    if string(chars[0:3]) == `'''` && string(chars[charsSize-3:charsSize]) == `'''` {
        parsedVal = chars[3 : charsSize-3]
        if parsedVal[0] == '\n' {
            parsedVal = parsedVal[1:]
        }
        return xtype.NewStringObject(string(parsedVal)), nil
    }
    if string(chars[0:3]) == `"""` && string(chars[charsSize-3:charsSize]) == `"""` {
        parsedVal = chars[3 : charsSize-3]
        if parsedVal[0] == '\n' {
            parsedVal = parsedVal[1:]
        }
        return xtype.NewStringObject(string(parsedVal)), nil
    }
    // 单行字符串
    if chars[0] == '\'' && chars[charsSize-1] == '\'' {
        if util.RunesContains(chars, '\n') {
            return nil, errors.New("new lines not allowed on single line string literals")
        }
        return xtype.NewStringObject(string(chars[1 : charsSize-1])), nil
    }
    if chars[0] == '"' && chars[charsSize-1] == '"' {
        return xtype.NewStringObject(string(chars[1 : charsSize-1])), nil
    }
    // Single line array (normalized)
    if chars[0] == '[' && chars[charsSize-1] == ']' {
        arr, err := parseArray(chars)
        if err != nil {
            return nil, err
        }
        return xtype.NewMapObject(arr), nil
    }
    // Inline table (normalized)
    if chars[0] == '{' && chars[charsSize-1] == '}' {
        arr, err := parseInlineTable(chars)
        if err != nil {
            return nil, err
        }
        return xtype.NewMapObject(arr), nil
    }
    return nil, errors.New("Unknown value type: " + val)
}

// parseArray 解析数组
func parseArray(chars []rune) (*xtype.Map, error) {
    openBrackets := 0
    openString := false
    openCurlyBraces := 0
    openLString := false
    buffer := ""

    charsSize := len(chars)
    arr := xtype.NewMap()
    keyPos := 0
    for i := 0; i < charsSize; i++ {
        if chars[i] == '[' && !openString && !openLString {
            openBrackets++
            if openBrackets == 1 {
                continue
            }
        } else if chars[i] == ']' && !openString && !openLString {
            openBrackets--
            if openBrackets == 0 {
                if strings.TrimSpace(buffer) != "" {
                    obj, err := ParseSingle(strings.TrimSpace(buffer))
                    if err != nil {
                        return nil, err
                    }
                    arr.Add(xtype.NewNumberKey(strconv.Itoa(keyPos)), obj)
                    keyPos++
                }
                return arr, nil
            }
        } else if chars[i] == '"' && ((i == 0) || (i > 0 && chars[i-1] != '\\')) && !openLString {
            openString = !openString
        } else if chars[i] == '\'' && !openString {
            openLString = !openLString
        } else if chars[i] == '{' && !openString && !openLString {
            openCurlyBraces++
        } else if chars[i] == '}' && !openString && !openLString {
            openCurlyBraces--
        }

        if (chars[i] == ',' || chars[i] == '}') && !openString && !openLString && openBrackets == 1 && openCurlyBraces == 0 {
            if chars[i] == '}' {
                buffer += string(chars[i])
            }
            buffer = strings.TrimSpace(buffer)
            if buffer != "" {
                obj, err := ParseSingle(strings.TrimSpace(buffer))
                if err != nil {
                    return nil, err
                }
                arr.Add(xtype.NewNumberKey(strconv.Itoa(keyPos)), obj)
                keyPos++
            }
            buffer = ""
        } else {
            buffer += string(chars[i])
        }
    }
    return nil, errors.New("Wrong array definition:" + string(chars))
}

// parseInlineTable Parse inline tables into common table array
func parseInlineTable(chars []rune) (*xtype.Map, error) {
    charsSize := len(chars)
    if chars[0] == '{' && chars[charsSize-1] == '}' {
        chars = chars[1 : charsSize-1]
    } else {
        return nil, errors.New("invalid inline table definition: " + string(chars))
    }

    charsSize = len(chars)
    openString := false
    openLString := false
    openBrackets := 0
    buf := bytes.Buffer{}

    arr := xtype.NewMap()
    for i := 0; i < charsSize; i++ {
        if chars[i] == '"' && chars[i-1] != '\\' {
            openString = !openString
        } else if chars[i] == '\'' {
            openLString = !openLString
        } else if chars[i] == '[' && !openString && !openLString {
            openBrackets++
        } else if chars[i] == ']' && !openString && !openLString {
            openBrackets--
        }

        if chars[i] == ',' && !openString && !openLString && openBrackets == 0 {
            obj, err := parseInlineTableFieldValue(buf.String())
            if err != nil {
                return nil, err
            }
            arr.Merge(obj)
            // keyPos++
            buf.Reset()
        } else {
            buf.WriteRune(chars[i])
        }
    }

    // parse last buffer
    obj, err := parseInlineTableFieldValue(buf.String())
    if err != nil {
        return nil, err
    }
    arr.Merge(obj)
    // keyPos++
    return arr, nil
}

// parseInlineTableFieldValue 解析键值对内容
func parseInlineTableFieldValue(snippet string) (*xtype.Map, error) {
    pos := strings.Index(snippet, "=")
    if pos <= 0 {
        return nil, errors.New("[split] invalid inline toml table data: " + snippet)
    }
    field := strings.TrimSpace(snippet[0:pos])
    val := strings.TrimSpace(snippet[pos+1:])
    obj, err := ParseSingle(val)
    if err != nil {
        return nil, errors.New("[parse] invalid inline toml table data: " + val + " <= " + snippet)
    }
    arr := xtype.NewMap()
    arr.DeepAdd([]string{field}, obj)
    return arr, nil
}

// parseTomlTableName Parses TOML table names and returns the hierarchy array of table names.
func parseTomlTableName(chars []rune) []string {
    buffer := bytes.Buffer{}
    strOpen := false
    strOpenChar := '"'
    names := make([]string, 0)

    charsSize := len(chars)
    for i := 0; i < charsSize; i++ {
        if chars[i] == '"' {
            if !strOpen || (strOpen && chars[i-1] != '\\' && strOpenChar == '"') {
                strOpen = !strOpen
                if strOpen {
                    strOpenChar = '"'
                }
                continue
            }
        } else if chars[i] == '\'' {
            if !strOpen || (strOpen && chars[i-1] != '\\' && strOpenChar == '\'') {
                strOpen = !strOpen
                if strOpen {
                    strOpenChar = '\''
                }
                continue
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

// parseKeyValue 解析键值对
func parseKeyValue(arr *xtype.Map, keys []string, val string) error {
    val = strings.TrimSpace(val)
    obj, err := ParseSingle(val)
    if err != nil {
        return err
    }
    arr.DeepAdd(keys, obj)
    return nil
}
