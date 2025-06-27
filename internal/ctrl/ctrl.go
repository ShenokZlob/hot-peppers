package ctrl

import (
	"encoding/json"
	"hotpepper/internal/domain"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type PeppersRepository interface {
	Get(offset, limit int) ([]domain.Pepper, error)
	GetByName(name string) (domain.Pepper, error)
}

type Peppers struct {
	r PeppersRepository

	tableTmpl *template.Template
	once      sync.Once
	tmplErr   error
}

func NewPeppersHandler(mux *http.ServeMux, r PeppersRepository) {
	p := &Peppers{r: r}

	mux.HandleFunc("/table", p.GetTable)
	mux.HandleFunc("/scroll", p.Scroll)
	mux.HandleFunc("/pepper/", p.Get)
}

// /table
func (p *Peppers) GetTable(w http.ResponseWriter, r *http.Request) {
	log.Println("Controllers - GetTable")

	p.once.Do(p.loadTableTemplate)
	if p.tmplErr != nil {
		log.Printf("Failed to parse template: %v\n", p.tmplErr)
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := p.tableTmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Failed to execute template: %v\n", err)
		http.Error(w, "Ошибка вывода шаблона", http.StatusInternalServerError)
		return
	}
}

func (p *Peppers) loadTableTemplate() {
	tmplPath := filepath.Join("data", "templates", "table.html")
	p.tableTmpl, p.tmplErr = template.ParseFiles(tmplPath)
}

// /scroll
func (p *Peppers) Scroll(w http.ResponseWriter, r *http.Request) {
	log.Println("Controlles - Scroll")

	query := r.URL.Query()

	offsetStr := query.Get("offset")
	limitStr := query.Get("limit")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		log.Printf("Invalid offset: %v\n", err)
		http.Error(w, "Invalid offset", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20 // значение по умолчанию
	}

	peppers, err := p.r.Get(offset, limit)
	if err != nil {
		log.Printf("Failed to get peppers: %v\n", err)
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peppers)
}

// get
func (p *Peppers) Get(w http.ResponseWriter, r *http.Request) {
	log.Println("Controllers - Get")

	pepperName := strings.TrimPrefix(r.URL.Path, "/pepper/")

	if pepperName == "" {
		log.Println("Empty pepper name")
		http.Error(w, "Нет такого перца", http.StatusNotFound)
		return
	}

	pepper, err := p.r.GetByName(pepperName)
	if err != nil {
		log.Printf("Failed to find pepper: %v\n", err)
		http.Error(w, "Failed to find pepper", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("data/templates/pepper.html")
	if err != nil {
		log.Printf("Failed to parse file: %v\n", err)
		http.Error(w, "Failed to parse pepper", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, pepper)
}
