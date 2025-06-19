package useful

import (
	"errors"
	"strconv"

	tele "gopkg.in/telebot.v3"
)

// makeCategoryKeyboard builds inline buttons for all categories.
func makeCategoryKeyboard() *tele.ReplyMarkup {
	mk := &tele.ReplyMarkup{}
	var buttons []tele.Btn
	for _, name := range categoryOrder {
		buttons = append(buttons, mk.Data(name, "useful_cat", prefixCategory+name))
	}
	rows := mk.Split(2, buttons)
	mk.Inline(rows...)
	return mk
}

// makeItemsKeyboard builds buttons for items within a category.
func makeItemsKeyboard(cat string) (*tele.ReplyMarkup, error) {
	items, ok := categories[cat]
	if !ok {
		return nil, errors.New("kategoriya topilmadi")
	}
	mk := &tele.ReplyMarkup{}
	var rows []tele.Row
	for idx, it := range items {
		if it.Name == "" {
			continue
		}
		data := itemPrefix + cat + "|" + strconv.Itoa(idx)
		rows = append(rows, mk.Row(mk.Data(it.Name, "useful_item", data)))
	}
	back := mk.Data("⬅️Orqaga", "useful_back", itemsBackPrefix+cat)
	rows = append(rows, mk.Row(back))
	mk.Inline(rows...)
	return mk, nil
}

// makeLinkKeyboard builds open link and back buttons for a specific item.
func makeLinkKeyboard(cat string, idx int) *tele.ReplyMarkup {
	items := categories[cat]
	it := items[idx]
	mk := &tele.ReplyMarkup{}
	open := mk.URL("havolani ochish", it.Link)
	back := mk.Data("⬅️Orqaga", "items_back", itemsBackPrefix+cat)
	mk.Inline(mk.Row(open), mk.Row(back))
	return mk
}
