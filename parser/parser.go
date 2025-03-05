package parser

import (
    "encoding/json"
	"fmt"
	"regexp"
	"strings"
    
	"github.com/PuerkitoBio/goquery"
)

// GetScripts извлекает содержимое всех тегов tagName из HTML-контента.
// По умолчанию ищет <script>.
func GetScripts(htmlContent string, tagName string) []string {
	if tagName == "" {
		tagName = "script"
	}

	var scripts []string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		// Если ошибка при парсинге, возвращаем nil
		return nil
	}

	doc.Find(tagName).Each(func(i int, s *goquery.Selection) {
		scripts = append(scripts, s.Text())
	})

	return scripts
}



func FindJSON(scriptText, pattern string) (*json.RawMessage, error) {
	// Экранируем специальные символы в паттерне
	escapedPattern := regexp.QuoteMeta(pattern)
	
	// Более широкий поиск с максимальной гибкостью
	re := regexp.MustCompile(fmt.Sprintf(`%s\s*=\s*(\{.*?\});`, escapedPattern))
	
	match := re.FindStringSubmatch(scriptText)
	
	if len(match) < 2 {
		return nil, fmt.Errorf("pattern not found: %s", pattern)
	}
	
	var jsonData json.RawMessage
	
	err := json.Unmarshal([]byte(match[1]), &jsonData)
	if err != nil {
		return nil, fmt.Errorf("JSON decode error: %v", err)
	}
	
	return &jsonData, nil
}


func GetYTVideoStats(htmlContent string)(*json.RawMessage, error) {
    scripts := GetScripts(htmlContent, "script")

    for _, script := range scripts {
    //for _, number := range foundIndices {
        j, err := FindJSON(script, "var ytInitialData")
        if err == nil {
            return j, nil
        }

        // Дальнейшая обработка j
    }
    var j *json.RawMessage
    return j, fmt.Errorf("Not found stats script")
}



