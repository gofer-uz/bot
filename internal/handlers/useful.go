package handlers

import (
	"encoding/json"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	tele "gopkg.in/telebot.v3"
)

type item struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
	Link string `json:"link"`
}

// categories maps category name -> list of items.
var (
	categories     map[string][]item
	categoryOrder  []string
	categoriesOnce sync.Once
)

// text and href helpers for button generation.
func (it item) text() string { return it.Name }

func (it item) href() string { return it.Link }

// loadSources reads and parses sources.json exactly once.
func loadSources(path string) {
	categoriesOnce.Do(func() {
		data, err := os.ReadFile(path)
		if err != nil {
			slog.Error("failed to read sources.json", "err", err)
			return
		}

		// sources.json format: {"category":[{name,desc,link},...]}
		if err := json.Unmarshal(data, &categories); err != nil {
			slog.Error("failed to parse sources.json", "err", err)
			return
		}

		// build ordered list of category names and sort them so menu order is stable
		for name := range categories {
			categoryOrder = append(categoryOrder, name)
		}
		sort.Strings(categoryOrder)
	})
}

// registerUseful sets up the /useful command and callback handling.
func registerUseful(b *tele.Bot) {
	const (
		cmd             = "/useful"
		prefix          = "useful|" // callback data prefix
		srcFile         = "./data/sources.json"
		itemPrefix      = "useful_item|" // item selection callback data
		itemsBackPrefix = "items_back|"  // back to items list
	)

	loadSources(srcFile)

	// helper function to send categories list
	sendCategories := func(c tele.Context) error {
		if len(categories) == 0 {
			return c.Send("manbalar mavjud emas")
		}
		markup := &tele.ReplyMarkup{}
		var buttons []tele.Btn
		for _, name := range categoryOrder {
			buttons = append(buttons, markup.Data(name, "useful_cat", prefix+name))
		}
		rows := markup.Split(2, buttons)
		markup.Inline(rows...)
		return c.Send("kategoriyani tanlang", markup)
	}

	// handler for /useful command – shows top-level categories
	b.Handle(cmd, func(c tele.Context) error {
		return sendCategories(c)
	})

	// callback handler – processes category selection by unique id
	b.Handle("\fuseful_cat", func(c tele.Context) error {
		data := c.Data()
		if !strings.HasPrefix(data, prefix) {
			// not our callback
			return nil
		}
		catName := strings.TrimPrefix(data, prefix)
		items, ok := categories[catName]
		if !ok {
			return c.Respond(&tele.CallbackResponse{Text: "nomalum kategoriya"})
		}

		markup := &tele.ReplyMarkup{}
		var rows []tele.Row
		for idx, it := range items {
			text := it.text()
			if text == "" {
				continue
			}
			data := itemPrefix + catName + "|" + strconv.Itoa(idx)
			rows = append(rows, markup.Row(markup.Data(text, "useful_item", data)))
		}
		// add back button
		backBtn := markup.Data("⬅️Orqaga", "useful_back", itemsBackPrefix+catName)
		rows = append(rows, markup.Row(backBtn))
		markup.Inline(rows...)

		return c.Edit("manbalar: "+catName, markup)
	})

	// back handler shows categories again
	b.Handle("\fuseful_back", func(c tele.Context) error {
		// simply edit message to categories list
		// reuse helper but on edit
		if len(categories) == 0 {
			return c.Respond(&tele.CallbackResponse{Text: "manbalar topilmadi"})
		}
		markup := &tele.ReplyMarkup{}
		var buttons []tele.Btn
		for _, name := range categoryOrder {
			buttons = append(buttons, markup.Data(name, "useful_cat", prefix+name))
		}
		rows := markup.Split(2, buttons)
		markup.Inline(rows...)
		return c.Edit("kategoriyani tanlang", markup)
	})

	// item selection handler – shows description and link
	b.Handle("\fuseful_item", func(c tele.Context) error {
		data := c.Data() // useful_item|<cat>|<idx>
		if !strings.HasPrefix(data, itemPrefix) {
			return nil
		}
		payload := strings.TrimPrefix(data, itemPrefix)
		parts := strings.SplitN(payload, "|", 2)
		if len(parts) != 2 {
			return c.Respond(&tele.CallbackResponse{Text: "xato ma'lumot"})
		}
		catName, idxStr := parts[0], parts[1]
		items, ok := categories[catName]
		if !ok {
			return c.Respond(&tele.CallbackResponse{Text: "kategoriya topilmadi"})
		}
		i, err := strconv.Atoi(idxStr)
		if err != nil || i < 0 || i >= len(items) {
			return c.Respond(&tele.CallbackResponse{Text: "band topilmadi"})
		}
		it := items[i]

		markup := &tele.ReplyMarkup{}
		openBtn := markup.URL("havolani ochish", it.href())
		backBtn := markup.Data("⬅️Orqaga", "items_back", itemsBackPrefix+catName)
		markup.Inline(markup.Row(openBtn), markup.Row(backBtn))

		desc := it.Desc
		if desc == "" {
			desc = it.text()
		}

		return c.Edit(desc, markup)
	})

	// back from item to items list
	b.Handle("\fitems_back", func(c tele.Context) error {
		data := c.Data() // items_back|<cat>
		if !strings.HasPrefix(data, itemsBackPrefix) {
			return nil
		}
		cat := strings.TrimPrefix(data, itemsBackPrefix)
		items, ok := categories[cat]
		if !ok {
			return c.Respond(&tele.CallbackResponse{Text: "kategoria topilmadi"})
		}

		markup := &tele.ReplyMarkup{}
		var rows []tele.Row
		for idx, it := range items {
			text := it.text()
			if text == "" {
				continue
			}
			data := itemPrefix + cat + "|" + strconv.Itoa(idx)
			rows = append(rows, markup.Row(markup.Data(text, "useful_item", data)))
		}
		backBtn := markup.Data("⬅️Orqaga", "useful_back", "useful_back")
		rows = append(rows, markup.Row(backBtn))
		markup.Inline(rows...)

		return c.Edit("Manbalar: "+cat, markup)
	})
}
