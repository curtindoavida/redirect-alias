package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type redirectRule struct {
	redirectTo string
	typeCode   int // 1 = replace-in-place, else = send to https://redirectTo
}

var defaultRedirect string
var rules map[string]redirectRule

// Estruturas para carregar o JSON de configuração
type fileRule struct {
	RedirectTo string `json:"redirectTo"`
	Type       int    `json:"type"`
}

type configFile struct {
	DefaultRedirect string              `json:"defaultRedirect"`
	Rules           map[string]fileRule `json:"rules"`
}

func loadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var cfg configFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	defaultRedirect = cfg.DefaultRedirect
	rules = make(map[string]redirectRule, len(cfg.Rules))
	for host, fr := range cfg.Rules {
		rules[host] = redirectRule{redirectTo: fr.RedirectTo, typeCode: fr.Type}
	}
	return nil
}

func getPage(path string) string {
	segments := strings.Split(path, "/")
	if len(segments) > 0 {
		last := segments[len(segments)-1]
		if last == "" {
			return "none"
		}
		return last
	}
	return "none"
}

func robotsHandler(w http.ResponseWriter, r *http.Request) bool {
	if getPage(r.URL.Path) == "robots.txt" {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "User-agent: *")
		fmt.Fprintln(w, "Disallow: /")
		return true
	}
	return false
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	// robots.txt
	if robotsHandler(w, r) {
		return
	}

	requestedHost := r.Host
	requestURI := r.URL.RequestURI() // includes path and query

	for oldHost, rule := range rules {
		if strings.Contains(requestedHost, oldHost) {
			if rule.typeCode == 1 {
				// Replace old host inside full URL string construction
				scheme := "http"
				if r.TLS != nil {
					scheme = "https"
				}
				// Build full URL then replace host occurrence, emulating Java replaceAll
				full := fmt.Sprintf("%s://%s%s", scheme, requestedHost, requestURI)
				replaced := strings.ReplaceAll(full, oldHost, rule.redirectTo)
				// Permanent redirect
				w.Header().Set("Location", replaced)
				w.WriteHeader(http.StatusMovedPermanently)
				return
			}
			// tipo 2: redirect straight to https://redirectTo (already includes params when present)
			http.Redirect(w, r, "https://"+rule.redirectTo, http.StatusFound)
			return
		}
	}

	// default redirect
	http.Redirect(w, r, "https://"+defaultRedirect, http.StatusFound)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9080"
	}
	localhost := "localhost"

	// Carrega regras do arquivo externo
	rulesPath := os.Getenv("RULES_PATH")
	if rulesPath == "" {
		rulesPath = "redirects.json"
	}
	if err := loadConfig(rulesPath); err != nil {
		log.Fatalf("Falha ao carregar config: %v", err)
	}

	http.HandleFunc("/", mainHandler)
	log.Printf("Servidor iniciado em :%s", port)
	log.Fatal(http.ListenAndServe(localhost+":"+port, nil))
}
