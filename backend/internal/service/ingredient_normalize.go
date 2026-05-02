package service

import "strings"

// ingredientInputToDBName maps user/ML input (RU/EN) to canonical ingredient.name values
// produced by seeds/seed.go (normalizeIngredientName + ingredientClassMap).
var ingredientInputToDBName = map[string]string{
	// EN meat / eggs / fish (align with seed class map)
	"chicken": "meat",
	"beef":    "meat",
	"pork":    "meat",
	"turkey":  "meat",
	"egg":     "eggs",
	"eggs":    "eggs",
	"tuna":    "fish",
	"salmon":  "fish",
	// RU — proteins
	"курица":     "meat",
	"куриное":    "meat",
	"куриный":    "meat",
	"говядина":   "meat",
	"свинина":    "meat",
	"индейка":    "meat",
	"мясо":       "meat",
	"рыба":       "fish",
	"лосось":     "fish",
	"тунец":      "fish",
	"яйцо":       "eggs",
	"яйца":       "eggs",
	"яиц":        "eggs",
	"креветки":   "shrimp",
	"креветка":   "shrimp",
	"shrimp":     "shrimp",
	// Vegetables & staples
	"помидор":        "tomato",
	"помидоры":       "tomato",
	"томат":          "tomato",
	"томаты":         "tomato",
	"tomato":         "tomato",
	"tomatoes":       "tomato",
	"сыр":            "cheese",
	"cheese":         "cheese",
	"молоко":         "milk",
	"milk":           "milk",
	"мука":           "flour",
	"flour":          "flour",
	"сливочное":       "butter",
	"сливочное масло": "butter",
	"butter":          "butter",
	"масло":           "olive oil",
	"оливковое масло": "olive oil",
	"оливковое":      "olive oil",
	"olive oil":      "olive oil",
	"oil":            "olive oil",
	"паста":          "pasta",
	"pasta":          "pasta",
	"спагетти":       "pasta",
	"spaghetti":      "pasta",
	"чеснок":         "garlic",
	"garlic":         "garlic",
	"лук":            "onion",
	"onion":          "onion",
	"onions":         "onion",
	"картофель":      "potato",
	"картошка":       "potato",
	"potato":         "potato",
	"potatoes":       "potato",
	"рис":            "rice",
	"rice":           "rice",
	"бекон":          "meat",
	"bacon":          "meat",
	"соль":           "salt",
	"salt":           "salt",
	"перец":          "pepper",
	"pepper":         "pepper",
	"peppers":        "pepper",
	"сахар":          "sugar",
	"sugar":          "sugar",
	"лимон":          "lemon",
	"lemon":          "lemon",
	"базилик":        "basil",
	"basil":          "basil",
	"петрушка":       "parsley",
	"parsley":        "parsley",
	"шпинат":         "spinach",
	"spinach":        "spinach",
	"морковь":        "carrot",
	"carrot":         "carrot",
	"carrots":        "carrot",
	"огурец":         "cucumber",
	"огурцы":         "cucumber",
	"cucumber":       "cucumber",
	"cucumbers":      "cucumber",
	"гриб":           "mushroom",
	"грибы":          "mushroom",
	"mushroom":       "mushroom",
	"mushrooms":      "mushroom",
	"орегано":        "oregano",
	"oregano":        "oregano",
	"чили":           "chili",
	"перец чили":     "chili",
	"chili":          "chili",
	"йогурт":         "yogurt",
	"yogurt":         "yogurt",
	"сметана":        "cream",
	"cream":          "cream",
	"банан":          "banana",
	"banana":         "banana",
	"яблоко":         "apple",
	"яблоки":         "apple",
	"apple":          "apple",
	"apples":         "apple",
	"ягоды":          "berry",
	"berry":          "berry",
	"berries":        "berry",
	"фасоль":         "beans",
	"beans":          "beans",
	"вода":           "water",
	"water":          "water",
}

func normalizeMatchIngredientNames(in []string) []string {
	seen := make(map[string]struct{})
	var out []string
	for _, raw := range in {
		s := strings.ToLower(strings.TrimSpace(raw))
		s = strings.ReplaceAll(s, "ё", "е")
		if s == "" {
			continue
		}
		if canon, ok := ingredientInputToDBName[s]; ok && canon != "" {
			s = canon
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
