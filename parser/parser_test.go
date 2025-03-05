package parser

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тест на корректное извлечение содержимого <script>
func TestGetScripts(t *testing.T) {
	html := `
		<html>
			<head>
				<script>console.log("Hello, world!");</script>
			</head>
			<body>
				<script>alert("Test");</script>
			</body>
		</html>
	`

	scripts := GetScripts(html, "script")

	expected := []string{
		`console.log("Hello, world!");`,
		`alert("Test");`,
	}

	if len(scripts) != len(expected) {
		t.Errorf("Ожидалось %d скриптов, но получено %d", len(expected), len(scripts))
	}

	for i, script := range scripts {
		if script != expected[i] {
			t.Errorf("Ожидалось: %q, но получено: %q", expected[i], script)
		}
	}
}

// Тест на случай пустого HTML
func TestGetScripts_EmptyHTML(t *testing.T) {
	html := ``
	scripts := GetScripts(html, "script")

	if len(scripts) != 0 {
		t.Errorf("Ожидался пустой список, но получено %d элементов", len(scripts))
	}
}

// Тест на случай, если указанный тег отсутствует
func TestGetScripts_NoScripts(t *testing.T) {
	html := `<html><body><p>Просто текст</p></body></html>`
	scripts := GetScripts(html, "script")

	if len(scripts) != 0 {
		t.Errorf("Ожидался пустой список, но получено %d элементов", len(scripts))
	}
}

// Тест на использование тега по умолчанию (если передана пустая строка)
func TestGetScripts_DefaultTag(t *testing.T) {
	html := `
		<html>
			<head>
				<script>console.log("Default tag works");</script>
			</head>
		</html>
	`

	scripts := GetScripts(html, "")

	if len(scripts) != 1 || scripts[0] != `console.log("Default tag works");` {
		t.Errorf("Ошибка: ожидался один скрипт, но получено %v", scripts)
	}
}

// Тест на извлечение других тегов (например, <style>)
func TestGetScripts_ExtractStyleTag(t *testing.T) {
	html := `
		<html>
			<head>
				<style>body { background-color: red; }</style>
			</head>
		</html>
	`

	styles := GetScripts(html, "style")

	if len(styles) != 1 || styles[0] != `body { background-color: red; }` {
		t.Errorf("Ошибка: ожидался один style, но получено %v", styles)
	}
}










func TestFindJSONValidSimple(t *testing.T) {
	script := `var testData = {"key": "value"};`
	pattern := "var testData"
	
	result, err := FindJSON(script, pattern)
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	var data map[string]interface{}
	err = json.Unmarshal(*result, &data)
	
	if err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}
	
	if data["key"] != "value" {
		t.Errorf("Expected value 'value', got %v", data["key"])
	}
}


func TestFindJSONInvalidJSON(t *testing.T) {
	script := `var testData = {broken json};`
	pattern := "var testData"
	
	_, err := FindJSON(script, pattern)
	
	if err == nil {
		t.Fatal("Expected JSON decode error, got none")
	}
}

func TestFindJSONComplexNested(t *testing.T) {
	script := `var ytInitialData = {"users": [{"name": "John"}, {"name": "Jane"}]};`
	pattern := "var ytInitialData"
	
	result, err := FindJSON(script, pattern)
	
	if err != nil {
		t.Fatalf("Unexpected error: %v\nScript: %s\nPattern: %s", err, script, pattern)
	}
	
	var data map[string]interface{}
	err = json.Unmarshal(*result, &data)
	
	if err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}
	
	users, ok := data["users"].([]interface{})
	if !ok {
		t.Fatal("Expected 'users' to be a slice")
	}
	
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}

func BenchmarkFindJSON(b *testing.B) {
	script := `var testData = {"users": [{"name": "John"}, {"name": "Jane"}]};`
	pattern := "var testData"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindJSON(script, pattern)
	}
}





// Тестовые HTML-контенты с различными сценариями
const (
	htmlWithValidScript = `
		<html>
			<head>
				<script>
					var someOtherScript = "data";
					var ytInitialData = {"videoDetails": {"viewCount": "1000", "title": "Test Video"}};
				</script>
			</head>
		</html>
	`

	htmlWithoutYTScript = `
		<html>
			<head>
				<script>
					var someOtherScript = "data";
				</script>
			</head>
		</html>
	`

	htmlWithMultipleScripts = `
		<html>
			<head>
				<script>
					var firstScript = "data";
				</script>
				<script>
					var ytInitialData = {"videoDetails": {"viewCount": "2000", "title": "Another Test Video"}};
				</script>
				<script>
					var anotherScript = "more data";
				</script>
			</head>
		</html>
	`
)

func TestGetYTVideoStatsWithValidScript(t *testing.T) {
	result, err := GetYTVideoStats(htmlWithValidScript)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	var data map[string]interface{}
	err = json.Unmarshal(*result, &data)
	assert.NoError(t, err)

	videoDetails, ok := data["videoDetails"].(map[string]interface{})
	assert.True(t, ok, "videoDetails should be a map")

	viewCount, ok := videoDetails["viewCount"].(string)
	assert.True(t, ok, "viewCount should be a string")
	assert.Equal(t, "1000", viewCount)
}

func TestGetYTVideoStatsWithoutYTScript(t *testing.T) {
	result, err := GetYTVideoStats(htmlWithoutYTScript)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Not found stats script", err.Error())
}

func TestGetYTVideoStatsWithMultipleScripts(t *testing.T) {
	result, err := GetYTVideoStats(htmlWithMultipleScripts)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	var data map[string]interface{}
	err = json.Unmarshal(*result, &data)
	assert.NoError(t, err)

	videoDetails, ok := data["videoDetails"].(map[string]interface{})
	assert.True(t, ok, "videoDetails should be a map")

	viewCount, ok := videoDetails["viewCount"].(string)
	assert.True(t, ok, "viewCount should be a string")
	assert.Equal(t, "2000", viewCount)
}


func BenchmarkGetYTVideoStats(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetYTVideoStats(htmlWithValidScript)
	}
}