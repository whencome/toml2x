package toml2x

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/whencome/toml2x/formatter"
	"github.com/whencome/toml2x/parser"
)

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
		`'{"payment_delegate":{"merchant_id":"2","merchant_no":"123456","plan_id":"53206"}}'`,
		`"{\"payment_delegate\":{\"merchant_id\":\"2\",\"merchant_no\":\"123456\",\"plan_id\":\"53206\"}}"`,
	}
	for _, toml := range tomls {
		rs, err := parser.ParseSingle(toml)
		if err != nil {
			t.Logf("parse %s failed: %s\n", toml, err)
			t.Fail()
			continue
		}
		t.Logf("parse %s success: %+v\n\n", toml, rs.Json(true))
	}
}

func TestParseJsonSingle(t *testing.T) {
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
		`"https://www.baidu.com/"`,
		`'{"payment_delegate":{"merchant_id":"2","merchant_no":"123456","plan_id":"53206"}}'`,
		`"{\"payment_delegate\":{\"merchant_id\":\"2\",\"merchant_no\":\"123456\",\"plan_id\":\"53206\"}}"`,
	}
	for _, toml := range tomls {
		rs, err := Json("single", toml)
		if err != nil {
			t.Logf("parse %s failed: %s\n", toml, err)
			t.Fail()
			continue
		}
		t.Logf("parse %s success: %+v\n\n", toml, rs)
	}
}

func TestParseArray(t *testing.T) {
	tomlArr := `[ 'literal,', 'strings', 'quo"ted' ]`
	parsed, err := parser.ParseSingle(tomlArr)
	if err != nil {
		t.Logf("TestParseArray failed: %s \n", err)
		t.Fail()
	}
	t.Logf("parsed: %+v\n", parsed.Json(true))
}

func TestParseInlineTable(t *testing.T) {
	tomlInlineTable := `PR17 = [
  {title = "Home", url = "/", childs = []},
  {title = "Games", url = "/games", childs = [{title = "Game A", url = "/games/game-a", childs = []}, {title = "Game B", url = "/games/game-b", childs = []}]},
  {title = "About us", url = "/about", childs = []}
]`
	tomlInlineTable, _ = formatter.Normalize(tomlInlineTable)
	parsed, err := parser.ParseTable(tomlInlineTable)
	if err != nil {
		t.Logf("parseInlineTableFieldValue failed: %s \n", err)
		t.Fail()
		return
	}
	t.Logf("parsed: %+v\n", parsed.Json(true))
}

func TestParseSimple(t *testing.T) {
	tomlInlineTable := `key = "value"
bare_key = "value"
bare-key = "value"
1234 = "value"`
	tomlInlineTable, _ = formatter.Normalize(tomlInlineTable)
	parsed, err := parser.ParseTable(tomlInlineTable)
	if err != nil {
		t.Logf("TestParseSimple failed: %s \n", err)
		t.Fail()
		return
	}
	t.Logf("parsed: %+v\n", parsed.Json(true))
}

func TestParseSimpleTable(t *testing.T) {
	tomlTable := `# 模式，debug or release
mode = "debug"
# http端口
port = 8808
# 跨域访问设置
[site.cors]
# 是否允许跨域(为true时全局启用，为false时只启用IP白名单)
is_enabled = false
# IP白名单(is_enabled=false时启用)
ip_whitelist = ["1.2.3.4"]
# 测试json内容
ext_single_json = '{"payment_delegate":{"merchant_id":"2","merchant_no":"123456","plan_id":"53206"}}'
ext_double_json = "{\"payment_delegate\":{\"merchant_id\":\"2\",\"merchant_no\":\"123456\",\"plan_id\":\"53206\"}}"`
	tomlTable, _ = formatter.Normalize(tomlTable)
	parsed, err := parser.ParseTable(tomlTable)
	if err != nil {
		t.Logf("TestParseSimple failed: %s \n", err)
		t.Fail()
		return
	}
	t.Logf("parsed: %+v\n", parsed.Xml())
}

func TestParseTable(t *testing.T) {
	tomlFile := "example_2.toml"
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
	rs, err := Php("table", toml)
	if err != nil {
		t.Log("Normalize content failed\n")
		t.Fail()
	}

	t.Logf("%+v\n", rs)
}
