package parser

import (
    "io/ioutil"
    "os"
    "testing"

    "github.com/whencome/toml2x/formatter"
)

func TestNormalize(t *testing.T) {
    tomlFile := "example.toml"
    file, err := os.Open(tomlFile)
    if err != nil {
        t.Logf("open file %s failed \n", tomlFile)
        t.Fail()
    }
    defer file.Close()

    tomlBytes, err := ioutil.ReadAll(file)
    if err != nil {
        t.Log("read file content failed\n")
        t.Fail()
    }

    toml, err := formatter.Normalize(string(tomlBytes))
    if err != nil {
        t.Logf("normalize toml failed: %s\n", err)
        t.Fail()
    }

    t.Log(toml)
}

func TestParseArray(t *testing.T) {
    tomlArr := `[ 'literal,', 'strings', 'quo"ted' ]`
    parsed, err := ParseSingle(tomlArr)
    if err != nil {
        t.Logf("TestParseArray failed: %s \n", err)
        t.Fail()
    }
    t.Logf("parsed: %+v\n", parsed)
}

func TestParseInlineTable(t *testing.T) {
    tomlInlineTable := `PR17 = [
  {title = "Home", url = "/", childs = []},
  {title = "Games", url = "/games", childs = [{title = "Game A", url = "/games/game-a", childs = []}, {title = "Game B", url = "/games/game-b", childs = []}]},
  {title = "About us", url = "/about", childs = []}
]`
    parsed, err := parseInlineTableFieldValue(tomlInlineTable)
    if err != nil {
        t.Logf("parseInlineTableFieldValue failed: %s \n", err)
        t.Fail()
        return
    }
    t.Logf("parsed: %+v\n", parsed)
}

func TestParseMultiLine(t *testing.T) {
    toml := `key3 = """
One
Two"""`
    rs, err := ParseSingle(toml)
    if err != nil {
        t.Logf("parse %s failed: %s\n", toml, err)
        t.Fail()
        return
    }
    t.Logf("parse %s success: %+v\n", toml, rs)
}

func TestParseSingle(t *testing.T) {
    var tomls = []string{
        // number
        `123554`,
        `23.4056`,
        `0.2399`,
        // boolean
        `true`,
        `false`,
        // string
        `"good"`,
        `"hello,world"`,
    }
    for _, toml := range tomls {
        rs, err := ParseSingle(toml)
        if err != nil {
            t.Logf("parse %s failed: %s\n", toml, err)
            t.Fail()
            continue
        }
        t.Logf("parse %s success: %+v\n\n", toml, rs)
    }
}

func TestParseTable(t *testing.T) {
    tomlFile := "example.toml"
    file, err := os.Open(tomlFile)
    if err != nil {
        t.Logf("open file %s failed \n", tomlFile)
        t.Fail()
    }
    defer file.Close()

    tomlBytes, err := ioutil.ReadAll(file)
    if err != nil {
        t.Log("read file content failed\n")
        t.Fail()
    }
    
    toml := string(tomlBytes)
    toml, err = formatter.Normalize(toml)
    if err != nil {
        t.Log("Normalize content failed\n")
        t.Fail()
    }

    rs, err := ParseTable(toml)
    if err != nil {
        t.Logf("parse table failed: %s\n", err)
        t.Fail()
    }

    t.Logf("\n=======================\n")

    t.Logf("%+v\n", rs)
}
