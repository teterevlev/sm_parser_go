package parser

import (
	"encoding/json"
	"strings"
	"testing"
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








func TestFindStartingWithNormalCase(t *testing.T) {
	scripts := []string{"hello world", "hello go", "hi there", "hello universe"}
	toBeFound := "hello"
	expected := []int{0, 1, 3}
	
	result := FindStartingWith(scripts, toBeFound)
	
	if len(result) != len(expected) {
		t.Fatalf("Разная длина: ожидалось %d, получено %d", len(expected), len(result))
	}
	
	for i := range expected {
		if result[i] != expected[i] {
			t.Fatalf("Несовпадение на индексе %d: ожидалось %d, получено %d", i, expected[i], result[i])
		}
	}
}

func TestFindStartingWithEmptyScriptsList(t *testing.T) {
	scripts := []string{}
	toBeFound := "test"
	expected := []int{}
	
	result := FindStartingWith(scripts, toBeFound)
	
	if len(result) != len(expected) {
		t.Fatalf("Для пустого списка некорректная длина: ожидалось %d, получено %d", len(expected), len(result))
	}
}

func TestFindStartingWithNonExistentPrefix(t *testing.T) {
	scripts := []string{"apple", "banana", "cherry"}
	toBeFound := "grape"
	expected := []int{}
	
	result := FindStartingWith(scripts, toBeFound)
	
	if len(result) != len(expected) {
		t.Fatalf("Для несуществующего префикса некорректная длина: ожидалось %d, получено %d", len(expected), len(result))
	}
}

func TestFindStartingWithCaseSensitivity(t *testing.T) {
	scripts := []string{"Hello", "hello", "HELLO"}
	toBeFound := "hello"
	expected := []int{1}
	
	result := FindStartingWith(scripts, toBeFound)
	
	if len(result) != len(expected) {
		t.Fatalf("Разная длина: ожидалось %d, получено %d", len(expected), len(result))
	}
	
	for i := range expected {
		if result[i] != expected[i] {
			t.Fatalf("Несовпадение на индексе %d: ожидалось %d, получено %d", i, expected[i], result[i])
		}
	}
}

func TestFindStartingWithEmptyPrefix(t *testing.T) {
	scripts := []string{"test", "example", "sample"}
	toBeFound := ""
	expected := []int{0, 1, 2}
	
	result := FindStartingWith(scripts, toBeFound)
	
	if len(result) != len(expected) {
		t.Fatalf("Разная длина: ожидалось %d, получено %d", len(expected), len(result))
	}
	
	for i := range expected {
		if result[i] != expected[i] {
			t.Fatalf("Несовпадение на индексе %d: ожидалось %d, получено %d", i, expected[i], result[i])
		}
	}
}

func BenchmarkFindStartingWith(b *testing.B) {
	scripts := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		scripts[i] = "script_" + strings.Repeat("x", i%100)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindStartingWith(scripts, "script_")
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