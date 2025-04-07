package view

import (
	"log"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

var (
	templates map[string]*template.Template // Templates cache

	functions = template.FuncMap{ // Built-in custom template functions
		"formatDate": formatDate,
		"now":        now,
		"upper":      upper,
		"lower":      lower,
	}

	config = Config{
		ViewsDir:   "views",
		LayoutsDir: "layouts",
		Extension:  ".html",
		Cache:      false,
	}
)

// Config represents the view configuration
type Config struct {
	ViewsDir   string
	LayoutsDir string
	Extension  string
	Cache      bool
}

// Initialize initializes the view package with optional config overrides
func Initialize(cfg Config) {
	if cfg.ViewsDir != "" {
		config.ViewsDir = cfg.ViewsDir
	}
	if cfg.LayoutsDir != "" {
		config.LayoutsDir = cfg.LayoutsDir
	}
	if cfg.Extension != "" {
		config.Extension = cfg.Extension
	}
	config.Cache = cfg.Cache

	if config.Cache {
		templates = make(map[string]*template.Template)
	}
}

// Render renders a template using the default layout
func Render(w http.ResponseWriter, name string, data interface{}) error {
	return RenderWithLayout(w, name, "app", data)
}

// getTemplate loads or builds the template with all components
func getTemplate(name, layout string) (*template.Template, error) {
    name = ensureExtension(name)
    layout = ensureExtension(layout)
    cacheKey := name + ":" + layout

    if config.Cache {
        if templates == nil {
            templates = make(map[string]*template.Template)
        }
        if tmpl, ok := templates[cacheKey]; ok {
            return tmpl, nil
        }
    }

    // Build all paths
    layoutPath := filepath.Join(config.ViewsDir, config.LayoutsDir, layout)
    templatePath := filepath.Join(config.ViewsDir, name)
    partialsGlob := filepath.Join(config.ViewsDir, "partials", "*"+config.Extension)
    stylesGlob := filepath.Join(config.ViewsDir, "styles", "*"+config.Extension)

    // Use layout name without extension as root template name
    templateName := strings.TrimSuffix(layout, filepath.Ext(layout))
    tmpl := template.New(templateName).Funcs(functions)

    // 1. Parse layout first (parent template)
    tmpl, err := tmpl.ParseFiles(layoutPath)
    if err != nil {
        return nil, fmt.Errorf("error parsing layout: %v", err)
    }

    // 2. Parse main template (fills blocks)
    tmpl, err = tmpl.ParseFiles(templatePath)
    if err != nil {
        return nil, fmt.Errorf("error parsing page template: %v", err)
    }

    // 3. Parse partials (optional)
    if _, err := tmpl.ParseGlob(partialsGlob); err != nil {
        log.Printf("Note: partials not found at %s", partialsGlob)
    }

    // 4. Parse style templates (optional)
    if _, err := tmpl.ParseGlob(stylesGlob); err != nil {
        log.Printf("Note: styles not found at %s", stylesGlob)
    }

    // Debug: List all available templates
    if config.Cache {
        templates[cacheKey] = tmpl
        log.Printf("Cached template: %s", cacheKey)
    }
    
    log.Println("Available templates:")
    for _, t := range tmpl.Templates() {
        log.Println("-", t.Name())
    }

    return tmpl, nil
}

// RenderWithLayout remains the same as previous solution
func RenderWithLayout(w http.ResponseWriter, name, layout string, data interface{}) error {
    tmpl, err := getTemplate(name, layout)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return err
    }
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    
    templateName := strings.TrimSuffix(layout, filepath.Ext(layout))
    return tmpl.ExecuteTemplate(w, templateName, data)
}

// ensureExtension appends the configured extension if it's missing
func ensureExtension(name string) string {
	if filepath.Ext(name) == "" {
		return name + config.Extension
	}
	return name
}

// RegisterFunction adds a custom function to the template function map
func RegisterFunction(name string, fn interface{}) {
	functions[name] = fn
}

// ==========================
// Built-in template functions
// ==========================

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func now() time.Time {
	return time.Now()
}

func upper(s string) string {
	return strings.ToUpper(s)
}

func lower(s string) string {
	return strings.ToLower(s)
}
