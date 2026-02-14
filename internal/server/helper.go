package server

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"tenangantri/internal/model"
	"time"
)

func BuildFuncMap() template.FuncMap {
	return template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05")
		},
		"formatDuration": func(seconds int) string {
			if seconds < 60 {
				return "< 1 min"
			}
			minutes := seconds / 60
			if minutes < 60 {
				return "{{ . }} min"
			}
			// hours := minutes / 60
			// mins := minutes % 60
			return "{{ . }}h {{ . }}m"
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"mod": func(a, b int) int {
			return a % b
		},
		"sum": func(items interface{}, field string) int {
			total := 0
			switch v := items.(type) {
			case []model.Category:
				for _, item := range v {
					if field == "Priority" {
						total += item.Priority
					}
				}
			}
			return total
		},
		"countActive": func(items []model.Category) int {
			count := 0
			for _, item := range items {
				if item.IsActive {
					count++
				}
			}
			return count
		},
		"countInactive": func(items []model.Category) int {
			count := 0
			for _, item := range items {
				if !item.IsActive {
					count++
				}
			}
			return count
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call: even number of arguments required")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"js": func(s string) template.JS {
			return template.JS(s)
		},
		"uppercase": func(s string) string {
			return strings.ToUpper(s)
		},
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
		"now": func() time.Time {
			return time.Now()
		},
		"gt": func(a, b int) bool {
			return a > b
		},
		"buildPaginationURL": func(page int, filters map[string]interface{}, sortBy, sortOrder string) string {
			if page < 1 {
				page = 1
			}

			params := []string{}
			params = append(params, fmt.Sprintf("page=%d", page))

			// Add filters
			if dateFrom, ok := filters["date_from"].(string); ok && dateFrom != "" {
				params = append(params, fmt.Sprintf("date_from=%s", dateFrom))
			}
			if dateTo, ok := filters["date_to"].(string); ok && dateTo != "" {
				params = append(params, fmt.Sprintf("date_to=%s", dateTo))
			}
			if status, ok := filters["status"].(string); ok && status != "" {
				params = append(params, fmt.Sprintf("status=%s", status))
			}

			// Add sorting
			if sortBy != "" {
				params = append(params, fmt.Sprintf("sort_by=%s", sortBy))
			}
			if sortOrder != "" {
				params = append(params, fmt.Sprintf("sort_order=%s", sortOrder))
			}

			return "/staff/tickets?" + strings.Join(params, "&")
		},
	}
}

func LoadTemplate(funcMap template.FuncMap) (tmpl *template.Template, err error) {
	// Create a new template set
	tmpl = template.New("").Funcs(funcMap)
	err = filepath.Walk("web/templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			// Get relative path as template name (e.g., "customer/index.html")
			name, err := filepath.Rel("web/templates", path)
			if err != nil {
				return err
			}
			// Normalise path separators to /
			name = filepath.ToSlash(name)

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			_, err = tmpl.New(name).Parse(string(content))
			if err != nil {
				return err
			}
		}
		return nil
	})

	return tmpl, err
}
