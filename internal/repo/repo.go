package repo

import (
	"encoding/json"
	"fmt"
	"hotpepper/internal/domain"
	"os"
	"path/filepath"
)

type Peppers struct {
	all     []domain.Pepper
	peppers map[string]domain.Pepper
}

func NewPeppers() *Peppers {
	dirPath := filepath.Join("data", "peppers")
	files, err := os.ReadDir(dirPath)
	if err != nil {
		panic("failed to read peppers dir: " + err.Error())
	}

	var peppers []domain.Pepper
	m := make(map[string]domain.Pepper)
	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dirPath, f.Name()))
		if err != nil {
			continue
		}
		var pepper domain.Pepper
		if err := json.Unmarshal(data, &pepper); err != nil {
			continue
		}
		peppers = append(peppers, pepper)
		m[pepper.Name] = pepper
	}

	return &Peppers{all: peppers, peppers: m}
}

func (p *Peppers) Get(offset, limit int) ([]domain.Pepper, error) {
	if offset > len(p.all) {
		return []domain.Pepper{}, nil
	}
	end := offset + limit
	if end > len(p.all) {
		end = len(p.all)
	}
	return p.all[offset:end], nil
}

func (p *Peppers) GetByName(name string) (domain.Pepper, error) {
	v, ok := p.peppers[name]
	if !ok {
		return domain.Pepper{}, fmt.Errorf("not found pepper")
	}
	return v, nil
}
