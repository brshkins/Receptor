package handler

import (
	"regexp"
	"strings"
	"unicode"
)

var ingredientNameEnToRu = map[string]string{
	"apple":       "Яблоко",
	"apples":      "Яблоки",
	"banana":      "Банан",
	"basil":       "Базилик",
	"beef":        "Говядина",
	"butter":      "Масло сливочное",
	"carrot":      "Морковь",
	"cheese":      "Сыр",
	"chicken":     "Курица",
	"chili":       "Чили",
	"cinnamon":    "Корица",
	"cream":       "Сливки",
	"cucumber":    "Огурец",
	"egg":         "Яйцо",
	"flour":       "Мука",
	"garlic":      "Чеснок",
	"honey":       "Мёд",
	"lemon":       "Лимон",
	"milk":        "Молоко",
	"mushroom":    "Грибы",
	"olive oil":   "Оливковое масло",
	"onion":       "Лук",
	"oregano":     "Орегано",
	"parsley":     "Петрушка",
	"pasta":       "Паста",
	"pepper":      "Перец",
	"pork":        "Свинина",
	"potato":      "Картофель",
	"rice":        "Рис",
	"salt":        "Соль",
	"shrimp":      "Креветки",
	"sugar":       "Сахар",
	"tomato":      "Помидор",
	"tuna":        "Тунец",
	"walnut":      "Грецкие орехи",
	"walnuts":     "Грецкие орехи",
	"water":       "Вода",
	"meat":        "Мясо",
	"fish":        "Рыба",
	"eggs":        "Яйца",
	"tomatoes":    "Помидоры",
	"potatoes":    "Картофель",
	"onions":      "Лук",
	"carrots":     "Морковь",
	"mushrooms":   "Грибы",
	"beans":       "Фасоль",
	"berry":       "Ягоды",
	"berries":     "Ягоды",
	"spinach":     "Шпинат",
	"broccoli":    "Брокколи",
	"cabbage":     "Капуста",
	"cheddar":     "Чеддер",
	"mozzarella":  "Моцарелла",
	"parmesan":    "Пармезан",
	"yogurt":      "Йогурт",
	"instant_noodle": "Лапша быстрого приготовления",
}

// Фразы для описаний и шагов (как на фронте; длинные раньше).
var recipeBodyPhraseRu = []struct{ en, ru string }{
	{"warm apple bake with cinnamon and buttery crumble.", "Тёплая яблочная запеканка с корицей и масляным крамблом."},
	{"until al dente", "до состояния аль денте"},
	{"until cooked through", "пока не будет готово"},
	{"until golden and bubbling", "пока не зарумянится и не забулькает"},
	{"until golden", "пока не зарумянится"},
	{"until crisp", "пока не станет хрустящим"},
	{"until bubbling", "пока не забулькает"},
	{"until tender", "пока не станет мягким"},
	{"preheat oven to", "Разогрейте духовку до"},
	{"bring to a boil", "Доведите до кипения"},
	{"bring to boil", "Доведите до кипения"},
	{"simmer for", "Томите на слабом огне"},
	{"simmer ", "Томите на слабом огне "},
	{"boil in salted water", "варите в подсоленной воде"},
	{"boil ", "Варите "},
	{"bake until", "Выпекайте, пока"},
	{"bake for", "Выпекайте"},
	{"bake ", "Выпекайте "},
	{"serve with", "Подавайте с"},
	{"serve hot", "Подавайте горячим"},
	{"season with salt and pepper", "Посолите и поперчите"},
	{"season with", "Приправьте"},
	{"heat olive oil", "Разогрейте оливковое масло"},
	{"heat oil", "Разогрейте масло"},
	{"heat ", "Нагрейте "},
	{"add salt", "Добавьте соль"},
	{"add ", "Добавьте "},
	{"stir in", "Вмешайте"},
	{"stir ", "Помешайте "},
	{"mix well", "Хорошо перемешайте"},
	{"mix ", "Смешайте "},
	{"toss with", "Смешайте с"},
	{"slice ", "Нарежьте "},
	{"chop ", "Мелко нарежьте "},
	{"drain ", "Слейте "},
	{"top with", "Сверху выложите"},
	{"rub butter into flour", "Вотрите масло в муку"},
	{"fold in", "Аккуратно вмешайте"},
	{"cool before serving", "Остудите перед подачей"},
	{"classic spaghetti with tomato, garlic, and basil", "Классические спагетти с томатом, чесноком и базиликом"},
}

func localizeIngredientDisplayName(name string) string {
	key := strings.ToLower(strings.TrimSpace(name))
	if key == "" {
		return name
	}
	if ru, ok := ingredientNameEnToRu[key]; ok && strings.TrimSpace(ru) != "" {
		return ru
	}
	return name
}

func localizedRecipeDescription(title, en string) string {
	key := strings.ToLower(strings.TrimSpace(title))
	if e, ok := recipeBodyRuByTitleLower[key]; ok && strings.TrimSpace(e.Description) != "" {
		return e.Description
	}
	return translateRecipeBodyText(en)
}

func localizedRecipeSteps(title string, en []string) []string {
	key := strings.ToLower(strings.TrimSpace(title))
	if e, ok := recipeBodyRuByTitleLower[key]; ok && len(e.Steps) > 0 {
		return e.Steps
	}
	out := make([]string, len(en))
	for i, s := range en {
		out[i] = translateRecipeBodyText(s)
	}
	return out
}

func translateRecipeBodyText(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	hasCyr := false
	hasLongLatin := false
	latinRun := 0
	for _, r := range s {
		if unicode.Is(unicode.Cyrillic, r) {
			hasCyr = true
		}
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			latinRun++
			if latinRun >= 5 {
				hasLongLatin = true
			}
		} else {
			latinRun = 0
		}
	}
	if hasCyr && !hasLongLatin {
		return s
	}

	s = strings.ReplaceAll(s, "\u00b0C", "°C")
	s = regexp.MustCompile(`(?i)190В°C`).ReplaceAllString(s, "190 °C")
	s = regexp.MustCompile(`(\d+)В°C`).ReplaceAllString(s, "$1 °C")

	for _, p := range recipeBodyPhraseRu {
		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(p.en))
		s = re.ReplaceAllString(s, p.ru)
	}
	return s
}
