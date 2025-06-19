package useful

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/k0kubun/pp"
	tele "gopkg.in/telebot.v3"
)

const (
	cmdUseful       = "/useful"
	prefixCategory  = "useful|"
	itemPrefix      = "useful_item|"
	itemsBackPrefix = "items_back|"
	srcFile         = "./data/sources.json"
)

// Register wires all useful-related handlers.
func Register(b *tele.Bot) {
	loadSources(srcFile)

	b.Handle(cmdUseful, onUseful)
	b.Handle("\fuseful_cat", onCategory)
	b.Handle("\fuseful_back", onUsefulBack)
	b.Handle("\f"+"useful_item", onItem)
	b.Handle("\fitems_back", onItemsBack)
}

// onUseful sends the list of categories.
func onUseful(c tele.Context) error {
	if len(categories) == 0 {
		return c.Send("manbalar mavjud emas")
	}

	pp.Println(makeCategoryKeyboard())
	return c.Send("kategoriyani tanlang", makeCategoryKeyboard())
}

// onCategory shows items inside a category.
func onCategory(c tele.Context) error {
	slog.Info("onCategory", "data", c.Data())
	data := c.Data()
	if !strings.HasPrefix(data, prefixCategory) {
		return nil
	}
	cat := strings.TrimPrefix(data, prefixCategory)
	kb, err := makeItemsKeyboard(cat)
	if err != nil {
		return c.Respond(&tele.CallbackResponse{Text: err.Error()})
	}
	return c.Edit("manbalar: "+cat, kb)
}

// onItem shows description and link for a selected item.
func onItem(c tele.Context) error {
	payload := strings.TrimPrefix(c.Data(), itemPrefix) // <cat>|<idx>
	parts := strings.SplitN(payload, "|", 2)
	if len(parts) != 2 {
		return c.Respond(&tele.CallbackResponse{Text: "xato ma'lumot"})
	}
	cat, idxStr := parts[0], parts[1]
	idx, _ := strconv.Atoi(idxStr)

	items, ok := categories[cat]
	if !ok || idx < 0 || idx >= len(items) {
		return c.Respond(&tele.CallbackResponse{Text: "kategoriya topilmadi"})
	}

	it := items[idx]
	kb := makeLinkKeyboard(cat, idx)

	desc := it.Desc
	if desc == "" {
		desc = it.Name
	}
	return c.Edit(desc, kb)
}

// onUsefulBack returns from category view to top-level categories.
func onUsefulBack(c tele.Context) error {
	return c.Edit("kategoriyani tanlang", makeCategoryKeyboard())
}

// onItemsBack returns from item view to list of items within its category.
func onItemsBack(c tele.Context) error {
	cat := strings.TrimPrefix(c.Data(), itemsBackPrefix)
	kb, err := makeItemsKeyboard(cat)
	if err != nil {
		return c.Respond(&tele.CallbackResponse{Text: err.Error()})
	}
	return c.Edit("Manbalar: "+cat, kb)
}
