package useful

import (
	"encoding/json"
	"log/slog"
	"os"
	"sort"
	"sync"
)

// item mirrors an entry in data/sources.json.
// only name and link are used for buttons; desc may be used later.
type item struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
	Link string `json:"link"`
}

var (
	categories     map[string][]item
	categoryOrder  []string
	categoriesOnce sync.Once
)

// loadSources reads and parses sources.json exactly once and populates caches.
func loadSources(path string) {
	categoriesOnce.Do(func() {
		data, err := os.ReadFile(path)
		if err != nil {
			slog.Error("failed to read sources.json", "err", err)
			return
		}
		if err := json.Unmarshal(data, &categories); err != nil {
			slog.Error("failed to parse sources.json", "err", err)
			return
		}
		for name := range categories {
			categoryOrder = append(categoryOrder, name)
		}
		sort.Strings(categoryOrder)
	})
}
