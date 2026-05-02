package service

import (
	"strings"
	"unicode"
)

// Ingredient normalization & localization lives in the bot only.
// Canonical names must match backend DB (see seeds/seed.go ingredientClassMap).

var ingredientSynonymsToEN = map[string]string{
	// tomato
	"помидор":      "tomato",
	"помидоры":     "tomato",
	"томат":        "tomato",
	"томаты":       "tomato",
	"tomato":       "tomato",
	"tomatoes":     "tomato",
	// cheese
	"сыр":          "cheese",
	"cheese":       "cheese",
	// milk
	"молоко":       "milk",
	"milk":         "milk",
	// egg → DB name "eggs"
	"яйцо":         "eggs",
	"яйца":         "eggs",
	"яиц":          "eggs",
	"egg":          "eggs",
	"eggs":         "eggs",
	// meat (chicken/beef/pork collapse to "meat" in seed)
	"курица":       "meat",
	"куриное":      "meat",
	"куриный":      "meat",
	"chicken":      "meat",
	"говядина":     "meat",
	"beef":         "meat",
	"свинина":      "meat",
	"pork":         "meat",
	"индейка":      "meat",
	"turkey":       "meat",
	"бекон":        "meat",
	"bacon":        "meat",
	"мясо":         "meat",
	"meat":         "meat",
	// fish
	"тунец":        "fish",
	"tuna":         "fish",
	"лосось":       "fish",
	"salmon":       "fish",
	"рыба":         "fish",
	"fish":         "fish",
	// flour
	"мука":         "flour",
	"flour":        "flour",
	// cooking oil → olive oil (matches seeded recipes)
	"масло":           "olive oil",
	"растительное":    "olive oil",
	"oil":             "olive oil",
	"оливковое":       "olive oil",
	"оливковое масло": "olive oil",
	"olive oil":       "olive oil",
	// pasta
	"паста":        "pasta",
	"pasta":        "pasta",
	"спагетти":     "pasta",
	"spaghetti":    "pasta",
	// garlic
	"чеснок":       "garlic",
	"garlic":       "garlic",
	// onion
	"лук":          "onion",
	"onion":        "onion",
	// potato
	"картофель":    "potato",
	"картошка":     "potato",
	"potato":       "potato",
	"potatoes":     "potato",
	// rice
	"рис":          "rice",
	"rice":         "rice",
	// butter
	"сливочное":       "butter",
	"сливочное масло": "butter",
	"butter":          "butter",
	// salt / pepper
	"соль":         "salt",
	"salt":         "salt",
	"перец":        "pepper",
	"pepper":       "pepper",
	// sugar
	"сахар":        "sugar",
	"sugar":        "sugar",
	// lemon
	"лимон":        "lemon",
	"lemon":        "lemon",
	// basil / parsley / spinach
	"базилик":      "basil",
	"basil":        "basil",
	"петрушка":     "parsley",
	"parsley":      "parsley",
	"шпинат":       "spinach",
	"spinach":      "spinach",
	// carrot / cucumber
	"морковь":      "carrot",
	"carrot":       "carrot",
	"огурец":       "cucumber",
	"огурцы":       "cucumber",
	"cucumber":     "cucumber",
	// mushroom
	"гриб":         "mushroom",
	"грибы":        "mushroom",
	"mushroom":     "mushroom",
	"mushrooms":    "mushroom",
	// oregano / chili
	"орегано":      "oregano",
	"oregano":      "oregano",
	"чили":         "chili",
	"перец чили":   "chili",
	"chili":        "chili",
	// yogurt / cream
	"йогурт":       "yogurt",
	"yogurt":       "yogurt",
	"сметана":      "cream",
	"cream":        "cream",
	// banana / apple / berry
	"банан":        "banana",
	"banana":       "banana",
	"яблоко":       "apple",
	"яблоки":       "apple",
	"apple":        "apple",
	"ягоды":        "berry",
	"berry":        "berry",
	"berries":      "berry",
	"креветки":     "shrimp",
	"креветка":     "shrimp",
	"shrimp":       "shrimp",
	// beans
	"фасоль":       "beans",
	"beans":        "beans",
	// water
	"вода":         "water",
	"water":        "water",
}

var ingredientENToRU = map[string]string{
	"tomato":     "помидор",
	"cheese":     "сыр",
	"milk":       "молоко",
	"eggs":       "яйца",
	"egg":        "яйца",
	"meat":       "мясо",
	"chicken":    "курица",
	"beef":       "говядина",
	"pork":       "свинина",
	"bacon":      "бекон",
	"fish":       "рыба",
	"tuna":       "тунец",
	"salmon":     "лосось",
	"shrimp":     "креветки",
	"flour":      "мука",
	"oil":        "оливковое масло",
	"olive oil":  "оливковое масло",
	"pasta":      "паста",
	"garlic":     "чеснок",
	"onion":      "лук",
	"potato":     "картофель",
	"rice":       "рис",
	"butter":     "сливочное масло",
	"salt":       "соль",
	"pepper":     "перец",
	"sugar":      "сахар",
	"lemon":      "лимон",
	"basil":      "базилик",
	"parsley":    "петрушка",
	"spinach":    "шпинат",
	"carrot":     "морковь",
	"cucumber":   "огурец",
	"mushroom":   "грибы",
	"oregano":    "орегано",
	"chili":      "чили",
	"yogurt":     "йогурт",
	"cream":      "сметана",
	"banana":     "банан",
	"apple":      "яблоко",
	"berry":      "ягоды",
	"beans":      "фасоль",
	"water":      "вода",
}

func normalizeIngredientsForBackend(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, 0, len(in))
	seen := make(map[string]struct{}, len(in))
	for _, raw := range in {
		w := strings.ToLower(strings.TrimSpace(raw))
		// basic normalization: trim punctuation and normalize 'ё' => 'е'
		w = strings.ReplaceAll(w, "ё", "е")
		w = strings.Trim(w, " \t\r\n,.;:")
		if w == "" {
			continue
		}
		if norm, ok := ingredientSynonymsToEN[w]; ok && strings.TrimSpace(norm) != "" {
			w = norm
		}
		if _, ok := seen[w]; ok {
			continue
		}
		seen[w] = struct{}{}
		out = append(out, w)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func localizeIngredientENToRU(s string) string {
	key := strings.ToLower(strings.TrimSpace(s))
	if key == "" {
		return ""
	}
	if ru, ok := ingredientENToRU[key]; ok {
		return ru
	}
	return ingredientFallbackRU(s)
}

func ingredientFallbackRU(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	for _, r := range s {
		if r > 127 {
			return s
		}
	}
	words := strings.Fields(s)
	for i := range words {
		w := []rune(strings.ToLower(words[i]))
		if len(w) == 0 {
			continue
		}
		w[0] = unicode.ToUpper(w[0])
		words[i] = string(w)
	}
	return strings.Join(words, " ")
}

