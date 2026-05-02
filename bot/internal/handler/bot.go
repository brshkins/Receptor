package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"receptor/bot/internal/client"
	"receptor/bot/internal/dto"
	"receptor/bot/internal/service"
)

type UserState string

const (
	UserStateIdle                     UserState = "idle"
	UserStateAwaitingEmail            UserState = "awaiting_email"
	UserStateAwaitingPass             UserState = "awaiting_password"
	UserStateAwaitingRegisterName     UserState = "awaiting_register_name"
	UserStateAwaitingRegisterEmail    UserState = "awaiting_register_email"
	UserStateAwaitingRegisterPassword UserState = "awaiting_register_password"
	UserStateAwaitingIngr             UserState = "awaiting_ingredients"
	UserStateAwaitingRecipeSearch     UserState = "awaiting_recipe_search"
)

type UserSession struct {
	State UserState
	Email string
	Token string

	AuthFlow string
	Name     string

	ProfileID    int64
	ProfileEmail string
	ProfileName  string

	ActiveNav       string 
	RecipesLastPage int
	FavsLastPage    int
	BackNav         string 
	BackPage        int

	MatchResults  []dto.MatchResponse `json:"match_results,omitempty"`
	MatchLastPage int                 `json:"match_last_page,omitempty"`

	// Каталог «Рецепты»: как на сайте — поиск по названию (RU/EN), категория, сортировка.
	RecipesSearchQuery    string `json:"recipes_search,omitempty"`
	RecipesFilterCategory string `json:"recipes_cat,omitempty"`
	RecipesSortBy         string `json:"recipes_sort_by,omitempty"` // name | time | category
	RecipesSortDesc       bool   `json:"recipes_sort_desc,omitempty"`

	PanelChatID    int64 `json:"panel_chat_id,omitempty"`
	PanelMessageID int   `json:"panel_message_id,omitempty"`
}

type TelegramHandler struct {
	token       string
	httpClient  *http.Client
	botService  service.BotService
	apiBaseURL  string
	fileBaseURL string

	store *SessionStore
}

func NewTelegramHandler(token string, httpClient *http.Client, botService service.BotService, store *SessionStore) (*TelegramHandler, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, fmt.Errorf("telegram token is empty")
	}
	if botService == nil {
		return nil, fmt.Errorf("bot service is nil")
	}
	if store == nil {
		store = &SessionStore{sessions: make(map[int64]*UserSession)}
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	return &TelegramHandler{
		token:       token,
		httpClient:  httpClient,
		botService:  botService,
		apiBaseURL:  "https://api.telegram.org/bot" + token,
		fileBaseURL: "https://api.telegram.org/file/bot" + token,
		store:       store,
	}, nil
}

func (h *TelegramHandler) mainMenuKeyboard() *tgInlineKeyboardMarkup {
	return &tgInlineKeyboardMarkup{
		InlineKeyboard: [][]tgInlineKeyboardButton{
			{{Text: "📖 Рецепты", CallbackData: "nav:recipes"}},
			{{Text: "❤️ Избранное", CallbackData: "nav:favorites"}},
			{{Text: "🔍 Подбор по продуктам", CallbackData: "nav:match"}},
			{{Text: "👤 Профиль", CallbackData: "nav:profile"}},
		},
	}
}

func (h *TelegramHandler) loginKeyboard() *tgInlineKeyboardMarkup {
	return &tgInlineKeyboardMarkup{
		InlineKeyboard: [][]tgInlineKeyboardButton{
			{{Text: "🔐 Войти", CallbackData: "nav:login"}},
			{{Text: "📝 Регистрация", CallbackData: "nav:register"}},
		},
	}
}

func (h *TelegramHandler) compactMenuRow() []tgInlineKeyboardButton {
	return []tgInlineKeyboardButton{{Text: "🏠 В меню", CallbackData: "nav:menu"}}
}

// Клавиатура экрана профиля: выход и переход в меню (для авторизованного пользователя).
func (h *TelegramHandler) profileScreenKeyboard() *tgInlineKeyboardMarkup {
	return &tgInlineKeyboardMarkup{
		InlineKeyboard: [][]tgInlineKeyboardButton{
			{{Text: "🚪 Выйти", CallbackData: "nav:logout"}},
			h.compactMenuRow(),
		},
	}
}

func (h *TelegramHandler) withBack(kb *tgInlineKeyboardMarkup) *tgInlineKeyboardMarkup {
	backRow := []tgInlineKeyboardButton{{Text: "⬅️ Назад", CallbackData: "nav:back"}}
	if kb == nil {
		return &tgInlineKeyboardMarkup{InlineKeyboard: [][]tgInlineKeyboardButton{backRow}}
	}
	rows := make([][]tgInlineKeyboardButton, 0, len(kb.InlineKeyboard)+1)
	rows = append(rows, kb.InlineKeyboard...)
	rows = append(rows, backRow)
	return &tgInlineKeyboardMarkup{InlineKeyboard: rows}
}

func localizeCategory(cat string) string {
	cat = strings.TrimSpace(strings.ToLower(cat))
	if cat == "" {
		return "—"
	}
	switch cat {
	case "breakfast":
		return "завтрак"
	case "lunch":
		return "обед"
	case "dinner":
		return "ужин"
	case "dessert":
		return "десерт"
	case "snack":
		return "перекус"
	case "pasta":
		return "паста"
	case "salad":
		return "салат"
	case "soup":
		return "суп"
	case "meat":
		return "мясное"
	case "vegetarian":
		return "вегетарианское"
	default:
		return recipeTitleFallback(cat)
	}
}

func localizeRecipeTitle(title string) string {
	key := strings.TrimSpace(strings.ToLower(title))
	if key == "" {
		return "—"
	}
	if ru, ok := recipeTitleENToRU[key]; ok && strings.TrimSpace(ru) != "" {
		return ru
	}
	return recipeTitleFallback(title)
}

// Частые английские слова в названиях рецептов → нейтральные русские подписи.
var recipeTitleWordENtoRU = map[string]string{
	"chicken": "курица", "beef": "говядина", "pork": "свинина", "turkey": "индейка",
	"soup": "суп", "pasta": "паста", "salad": "салат", "toast": "тост", "bread": "хлеб",
	"tomato": "помидор", "potato": "картофель", "rice": "рис", "egg": "яйцо", "eggs": "яйца",
	"cheese": "сыр", "milk": "молоко", "butter": "масло", "garlic": "чеснок", "onion": "лук",
	"fish": "рыба", "tuna": "тунец", "salmon": "лосось", "shrimp": "креветки",
	"vegetable": "овощи", "vegetables": "овощи", "fried": "жареное", "baked": "запечённое",
	"breakfast": "завтрак", "style": "", "with": "с", "and": "и", "bowl": "боул",
	"berry": "ягоды", "berries": "ягоды", "cream": "сливки", "honey": "мёд", "apple": "яблоко",
	"banana": "банан", "noodle": "лапша", "noodles": "лапша", "chili": "чили", "bean": "фасоль",
	"beans": "фасоль", "meat": "мясо", "burger": "бургер", "patty": "котлета", "patties": "котлеты",
	"spaghetti": "спагетти", "penne": "пенне", "farfalle": "фарфалле", "basil": "базилик", "grilled": "гриль",
	"roasted": "запечённое", "omelette": "омлет", "waffle": "вафли", "pancake": "блин",
	"pancakes": "блины", "cake": "торт", "cupcake": "кекс", "cupcakes": "кексы",
}

func recipeTitleFallback(title string) string {
	title = strings.TrimSpace(title)
	if title == "" {
		return "—"
	}
	// Уже кириллица — не трогаем. Символы вроде é в «sauté» не должны отключать перевод.
	for _, r := range title {
		if unicode.Is(unicode.Cyrillic, r) {
			return title
		}
	}
	words := strings.Fields(title)
	out := make([]string, 0, len(words))
	for _, raw := range words {
		w := strings.TrimFunc(raw, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})
		if w == "" {
			continue
		}
		key := strings.ToLower(w)
		if ru, ok := recipeTitleWordENtoRU[key]; ok {
			if strings.TrimSpace(ru) != "" {
				out = append(out, ru)
			}
			continue
		}
		allLatin := len(w) > 0
		for _, r := range w {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
				allLatin = false
				break
			}
		}
		if allLatin && len(w) > 0 {
			continue
		}
		runes := []rune(strings.ToLower(w))
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
			out = append(out, string(runes))
		}
	}
	if len(out) == 0 {
		return "Блюдо"
	}
	return strings.Join(out, " ")
}

var recipeTitleENToRU = make(map[string]string)

func parseIngredients(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	parts := strings.FieldsFunc(s, func(r rune) bool {
		switch r {
		case ' ', '\n', '\t', ',', ';', ':':
			return true
		default:
			return false
		}
	})

	out := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.Trim(p, ".,;:!?-–—")
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

const matchResultsPerPage = 5

func (h *TelegramHandler) sendMatches(ctx context.Context, chatID int64, userID int64, matches []dto.MatchResponse, panel *tgCallbackQueryMsg) error {
	if len(matches) == 0 {
		sess := h.store.GetOrCreate(userID)
		sess.MatchResults = nil
		sess.MatchLastPage = 0
		_ = h.store.Save()
		return h.present(ctx, chatID, userID, "😔 Пока нет рецептов под ваш набор продуктов.\n\nПопробуйте другие ингредиенты или фото получше 📷", h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: [][]tgInlineKeyboardButton{h.compactMenuRow()}}), panel)
	}

	sess := h.store.GetOrCreate(userID)
	sess.MatchResults = matches
	sess.MatchLastPage = 0
	sess.ActiveNav = "match"
	_ = h.store.Save()
	return h.renderMatchesPage(ctx, chatID, userID, 0, panel)
}

func (h *TelegramHandler) renderMatchesPage(ctx context.Context, chatID int64, userID int64, page int, panel *tgCallbackQueryMsg) error {
	sess := h.store.GetOrCreate(userID)
	all := sess.MatchResults
	if len(all) == 0 {
		return h.present(ctx, chatID, userID, "Список подбора устарел — снова откройте «🔍 Подбор по продуктам» или отправьте продукты текстом.", h.withBack(h.mainMenuKeyboard()), panel)
	}

	perPage := matchResultsPerPage
	totalPages := (len(all) + perPage - 1) / perPage
	if page < 0 {
		page = 0
	}
	if page >= totalPages {
		page = totalPages - 1
	}
	start := page * perPage
	end := start + perPage
	if end > len(all) {
		end = len(all)
	}
	shown := all[start:end]

	sess.MatchLastPage = page
	sess.ActiveNav = "match"
	_ = h.store.Save()

	const (
		maxMissingToShow  = 4
		maxIngredientsRow = 14
	)

	var b strings.Builder
	b.WriteString("🔍 Подбор по продуктам")
	if totalPages > 1 {
		b.WriteString(fmt.Sprintf(" · стр. %d из %d", page+1, totalPages))
	}
	b.WriteString("\n\n")

	rows := make([][]tgInlineKeyboardButton, 0, len(shown)+6)
	for i, m := range shown {
		if i > 0 {
			b.WriteString("\n\n")
		}

		b.WriteString("🍽 ")
		b.WriteString(localizeRecipeTitle(m.Title))
		b.WriteString("\n")
		if m.CookingTime > 0 {
			b.WriteString("⏱ ")
			b.WriteString(strconv.Itoa(m.CookingTime))
			b.WriteString(" мин\n")
		}

		pct := int(m.MatchPercent * 100)
		b.WriteString("📊 ")
		b.WriteString(strconv.Itoa(pct))
		b.WriteString("% совпадение по продуктам\n")

		if len(m.Ingredients) > 0 {
			b.WriteString("🥕 Ингредиенты: ")
			showIng := m.Ingredients
			if len(showIng) > maxIngredientsRow {
				showIng = showIng[:maxIngredientsRow]
			}
			ruIng := make([]string, 0, len(showIng))
			for _, x := range showIng {
				ruIng = append(ruIng, localizeIngredientDisplayName(x))
			}
			b.WriteString(strings.Join(ruIng, ", "))
			if remain := len(m.Ingredients) - len(showIng); remain > 0 {
				b.WriteString(fmt.Sprintf(" и ещё %d", remain))
			}
			b.WriteString("\n")
		}

		if len(m.MissingIngredients) == 0 {
			b.WriteString("✨ Можно приготовить прямо сейчас 🎉")
		} else {
			b.WriteString("❌ Не хватает: ")
			miss := m.MissingIngredients
			if len(miss) > maxMissingToShow {
				miss = miss[:maxMissingToShow]
			}
			ruMiss := make([]string, 0, len(miss))
			for _, x := range miss {
				ruMiss = append(ruMiss, localizeIngredientDisplayName(x))
			}
			b.WriteString(strings.Join(ruMiss, ", "))
			if remain := len(m.MissingIngredients) - len(miss); remain > 0 {
				b.WriteString(fmt.Sprintf(" и ещё %d", remain))
			}
		}

		matchRow := []tgInlineKeyboardButton{
			{Text: "Подробнее", CallbackData: fmt.Sprintf("recipe:%d", m.RecipeID)},
			{Text: "❤️", CallbackData: fmt.Sprintf("fav:%d", m.RecipeID)},
		}
		rows = append(rows, matchRow)
	}

	if totalPages > 1 {
		navRow := make([]tgInlineKeyboardButton, 0, 2)
		if page > 0 {
			navRow = append(navRow, tgInlineKeyboardButton{Text: "⬅️ Назад", CallbackData: "nav:match_prev"})
		}
		if page+1 < totalPages {
			navRow = append(navRow, tgInlineKeyboardButton{Text: "➡️ Далее", CallbackData: "nav:match_next"})
		}
		if len(navRow) > 0 {
			rows = append(rows, navRow)
		}
	}

	rows = append(rows, h.compactMenuRow())
	kb := &tgInlineKeyboardMarkup{InlineKeyboard: rows}
	return h.present(ctx, chatID, userID, b.String(), h.withBack(kb), panel)
}

func (h *TelegramHandler) answerCallbackQuery(ctx context.Context, callbackQueryID string, text string) error {
	endpoint := h.apiBaseURL + "/answerCallbackQuery"
	payload := struct {
		CallbackQueryID string `json:"callback_query_id"`
		Text            string `json:"text,omitempty"`
	}{
		CallbackQueryID: callbackQueryID,
		Text:            strings.TrimSpace(text),
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal answerCallbackQuery payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("create answerCallbackQuery request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("answerCallbackQuery request: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func (h *TelegramHandler) RunLongPolling(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("ctx is nil")
	}

	offset := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		updates, err := h.getUpdates(ctx, offset, 30)
		if err != nil {
			select {
			case <-time.After(1 * time.Second):
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		for _, upd := range updates {
			if upd.UpdateID >= offset {
				offset = upd.UpdateID + 1
			}

			if upd.Message != nil && upd.Message.Chat != nil {
				if err := h.handleMessage(ctx, upd.Message); err != nil {
					uid := int64(0)
					if upd.Message.From != nil {
						uid = upd.Message.From.ID
					}
					_ = h.present(ctx, upd.Message.Chat.ID, uid, "Что-то пошло не так 😔 Напишите /start или «меню»", h.mainMenuKeyboard(), nil)
				}
			}

			if upd.CallbackQuery != nil && upd.CallbackQuery.Message != nil && upd.CallbackQuery.Message.Chat != nil {
				if err := h.handleCallbackQuery(ctx, upd.CallbackQuery); err != nil {
					_ = h.answerCallbackQuery(ctx, upd.CallbackQuery.ID, "")
				}
			}
		}
	}
}

func (h *TelegramHandler) handleMessage(ctx context.Context, msg *tgMessage) error {
	if msg == nil || msg.Chat == nil || msg.From == nil {
		return nil
	}

	userID := msg.From.ID
	chatID := msg.Chat.ID

	text := strings.TrimSpace(msg.Text)
	if text != "" && strings.HasPrefix(text, "/start") {
		_ = h.store.SetState(userID, UserStateIdle)
		return h.present(ctx, chatID, userID, "👋 Добро пожаловать!\n\nВыберите раздел ниже — подскажем, что приготовить ✨", h.mainMenuKeyboard(), nil)
	}

	if text != "" {
		return h.handleIncomingText(ctx, chatID, userID, text)
	}

	if len(msg.Photo) == 0 {
		return nil
	}

	if err := h.handleIncomingPhoto(ctx, chatID, userID, msg.Photo[len(msg.Photo)-1].FileID); err != nil {
		return err
	}
	return nil
}

func (h *TelegramHandler) handleIncomingText(ctx context.Context, chatID int64, userID int64, text string) error {
	st := h.store.GetCopy(userID)

	if st.State == UserStateAwaitingRecipeSearch {
		sess := h.store.GetOrCreate(userID)
		sess.RecipesSearchQuery = strings.TrimSpace(text)
		sess.RecipesLastPage = 0
		_ = h.store.Set(userID, UserStateIdle, st.Email, st.Token)
		return h.renderRecipesPage(ctx, userID, chatID, 0, nil)
	}

	if st.State == UserStateIdle {
		t := strings.ToLower(strings.TrimSpace(text))
		switch t {
		case "меню", "menu", "start", "старт":
			return h.present(ctx, chatID, userID, "✨ Главное меню\n\nВыберите, что сделать дальше 👇", h.mainMenuKeyboard(), nil)
		case "рецепты":
			return h.present(ctx, chatID, userID, "Откройте «📖 Рецепты» в меню.", h.mainMenuKeyboard(), nil)
		case "подбор", "подобрать", "поиск", "match":
			return h.present(ctx, chatID, userID, "Откройте «🔍 Подбор по продуктам» в меню.", h.mainMenuKeyboard(), nil)
		case "профиль", "profile":
			return h.present(ctx, chatID, userID, "Откройте «👤 Профиль» в меню.", h.mainMenuKeyboard(), nil)
		case "войти", "логин", "login":
			return h.present(ctx, chatID, userID, "Выберите «🔐 Войти» или «📝 Регистрация».", h.loginKeyboard(), nil)
		default:
			// Сообщение вида «помидор» или «помидор: сыр» вне режима подбора — всё равно делаем матч (не email/ключ).
			if strings.Contains(text, "@") {
				return h.present(ctx, chatID, userID, "Нажмите кнопку ниже или напишите «меню».\n\n✨ Главное меню 👇", h.mainMenuKeyboard(), nil)
			}
			parsed := parseIngredients(text)
			if !service.HasIngredientKeyword(parsed) {
				return h.present(ctx, chatID, userID, "Нажмите кнопку ниже или напишите «меню».\n\nЧтобы подобрать рецепт по продуктам — «🔍 Подбор по продуктам», затем список или фото.\n\n✨ Главное меню 👇", h.mainMenuKeyboard(), nil)
			}
			ingredients := service.NormalizeIngredientsForBackend(parsed)
			if len(ingredients) == 0 {
				return h.present(ctx, chatID, userID, "Не удалось разобрать продукты. Пример: помидор, сыр, яйцо", h.mainMenuKeyboard(), nil)
			}
			matches, err := h.botService.Match(ctx, ingredients)
			if err != nil {
				return h.present(ctx, chatID, userID, err.Error(), h.mainMenuKeyboard(), nil)
			}
			return h.sendMatches(ctx, chatID, userID, matches, nil)
		}
	}

	switch st.State {
	case UserStateAwaitingRegisterName:
		name := strings.TrimSpace(text)
		if err := service.ValidateRegisterName(name); err != nil {
			return h.present(ctx, chatID, userID, err.Error(), h.withBack(nil), nil)
		}
		sess := h.store.GetOrCreate(userID)
		sess.Name = name
		sess.AuthFlow = "register"
		_ = h.store.Save()
		_ = h.store.Set(userID, UserStateAwaitingRegisterEmail, "", st.Token)
		return h.present(ctx, chatID, userID, "Введите email", h.withBack(nil), nil)

	case UserStateAwaitingRegisterEmail:
		email := strings.TrimSpace(strings.ToLower(text))
		if err := service.ValidateEmail(email); err != nil {
			return h.present(ctx, chatID, userID, err.Error(), h.withBack(nil), nil)
		}
		_ = h.store.Set(userID, UserStateAwaitingRegisterPassword, email, st.Token)
		return h.present(ctx, chatID, userID, "Введите пароль", h.withBack(nil), nil)

	case UserStateAwaitingRegisterPassword:
		password := strings.TrimSpace(text)
		if err := service.ValidateRegisterPassword(password); err != nil {
			return h.present(ctx, chatID, userID, err.Error(), h.withBack(nil), nil)
		}
		sess := h.store.GetOrCreate(userID)
		email := strings.TrimSpace(st.Email)
		name := strings.TrimSpace(sess.Name)
		token, err := h.botService.Register(ctx, email, password, name)
		if err != nil {
			_ = h.store.Set(userID, UserStateIdle, "", "")
			sess.AuthFlow = ""
			sess.Name = ""
			_ = h.store.Save()
			return h.present(ctx, chatID, userID, err.Error(), h.withBack(nil), nil)
		}
		if strings.TrimSpace(token) == "" {
			_ = h.store.Set(userID, UserStateIdle, "", "")
			sess.AuthFlow = ""
			sess.Name = ""
			_ = h.store.Save()
			return h.present(ctx, chatID, userID, "Ошибка авторизации: пустой токен", h.withBack(nil), nil)
		}
		_ = h.store.Set(userID, UserStateIdle, email, token)
		me, err := h.botService.Me(ctx, token)
		if err != nil {
			_ = h.store.Set(userID, UserStateIdle, "", token)
			sess.AuthFlow = ""
			sess.Name = ""
			_ = h.store.Save()
			return h.present(ctx, chatID, userID, err.Error(), h.withBack(nil), nil)
		}
		sess.ProfileID = me.ID
		sess.ProfileEmail = me.Email
		sess.ProfileName = me.Name
		_ = h.store.Set(userID, UserStateIdle, me.Email, token)
		sess.AuthFlow = ""
		sess.Name = ""
		_ = h.store.Save()
		return h.present(ctx, chatID, userID, formatProfileCard(me.Name, me.Email), h.withBack(h.profileScreenKeyboard()), nil)

	case UserStateAwaitingEmail:
		email := strings.TrimSpace(strings.ToLower(text))
		if err := service.ValidateEmail(email); err != nil {
			return h.present(ctx, chatID, userID, err.Error(), h.withBack(nil), nil)
		}
		_ = h.store.Set(userID, UserStateAwaitingPass, email, st.Token)
		return h.present(ctx, chatID, userID, "Введите пароль", h.withBack(nil), nil)

	case UserStateAwaitingPass:
		email := st.Email
		password := strings.TrimSpace(text)
		if password == "" {
			return h.present(ctx, chatID, userID, "Введите пароль", h.withBack(nil), nil)
		}

		sess := h.store.GetOrCreate(userID)
		token, err := h.botService.Login(ctx, email, password)
		if err != nil {
			_ = h.store.Set(userID, UserStateIdle, "", "")
			sess.AuthFlow = ""
			sess.Name = ""
			_ = h.store.Save()
			return h.present(ctx, chatID, userID, err.Error(), h.withBack(nil), nil)
		}

		if strings.TrimSpace(token) == "" {
			log.Printf("ERROR: empty token after login userID=%d", userID)
			_ = h.store.Set(userID, UserStateIdle, "", "")
			sess.AuthFlow = ""
			sess.Name = ""
			_ = h.store.Save()
			return h.present(ctx, chatID, userID, "Ошибка авторизации: пустой токен", h.withBack(nil), nil)
		}

		_ = h.store.Set(userID, UserStateIdle, email, token)

		me, err := h.botService.Me(ctx, token)
		if err != nil {
			_ = h.store.Set(userID, UserStateIdle, "", token)
			sess.AuthFlow = ""
			sess.Name = ""
			_ = h.store.Save()
			return h.present(ctx, chatID, userID, err.Error(), h.withBack(nil), nil)
		}

		sess.ProfileID = me.ID
		sess.ProfileEmail = me.Email
		sess.ProfileName = me.Name

		_ = h.store.Set(userID, UserStateIdle, me.Email, token)
		sess.AuthFlow = ""
		sess.Name = ""
		_ = h.store.Save()
		return h.present(ctx, chatID, userID, formatProfileCard(me.Name, me.Email), h.withBack(h.profileScreenKeyboard()), nil)

	case UserStateAwaitingIngr:
		parsed := parseIngredients(text)
		ingredients := service.NormalizeIngredientsForBackend(parsed)
		if len(ingredients) == 0 {
			_ = h.store.Set(userID, UserStateAwaitingIngr, "", st.Token)
			return h.present(ctx, chatID, userID, "Не удалось разобрать список. Напишите продукты через запятую, например: помидор, сыр, яйцо", h.withBack(h.mainMenuKeyboard()), nil)
		}
		matches, err := h.botService.Match(ctx, ingredients)
		_ = h.store.Set(userID, UserStateIdle, "", st.Token)
		if err != nil {
			return h.present(ctx, chatID, userID, err.Error(), h.withBack(h.mainMenuKeyboard()), nil)
		}
		return h.sendMatches(ctx, chatID, userID, matches, nil)

	default:
		if st.State == UserState("awaiting_register_form") {
			sess := h.store.GetOrCreate(userID)
			sess.AuthFlow = "register"
			sess.Name = ""
			_ = h.store.Save()
			_ = h.store.Set(userID, UserStateAwaitingRegisterName, "", st.Token)
			return h.present(ctx, chatID, userID, "Введите имя", h.withBack(nil), nil)
		}
		return nil
	}
}

func (h *TelegramHandler) handleIncomingPhoto(ctx context.Context, chatID int64, userID int64, fileID string) error {
	st := h.store.GetCopy(userID)
	if st.State != UserStateAwaitingIngr {
		return h.present(ctx, chatID, userID, "Сначала нажмите «🔍 Подбор по продуктам», затем пришлите фото 📷", h.withBack(h.mainMenuKeyboard()), nil)
	}

	filePath, err := h.getFilePath(ctx, fileID)
	if err != nil {
		return err
	}

	img, err := h.downloadFile(ctx, filePath)
	if err != nil {
		return err
	}

	matches, err := h.botService.UploadImage(ctx, img)
	_ = h.store.Set(userID, UserStateIdle, "", st.Token)
	if err != nil {
		msg := "Не удалось распознать продукты на фото 😔 Попробуйте другое изображение или отправьте список продуктов текстом."
		return h.present(ctx, chatID, userID, msg, h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: [][]tgInlineKeyboardButton{h.compactMenuRow()}}), nil)
	}
	return h.sendMatches(ctx, chatID, userID, matches, nil)
}

func formatProfileCard(name, email string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "—"
	}
	email = strings.TrimSpace(email)
	if email == "" {
		email = "—"
	}
	return fmt.Sprintf("👤 Профиль\n\n🪪 Имя: %s\n📧 Email: %s", name, email)
}

func (h *TelegramHandler) renderProfileCard(ctx context.Context, chatID int64, userID int64, token string, panel *tgCallbackQueryMsg) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return h.present(ctx, chatID, userID, "Вы не авторизованы", h.withBack(h.loginKeyboard()), panel)
	}
	me, err := h.botService.Me(ctx, token)
	if err != nil {
		if err.Error() == "unauthorized" {
			_ = h.store.Set(userID, UserStateIdle, "", "")
			sess := h.store.GetOrCreate(userID)
			sess.ProfileID = 0
			sess.ProfileEmail = ""
			sess.ProfileName = ""
			_ = h.store.Save()
			return h.present(ctx, chatID, userID, "Вы не авторизованы", h.withBack(h.loginKeyboard()), panel)
		}
		return h.present(ctx, chatID, userID, err.Error(), h.withBack(h.mainMenuKeyboard()), panel)
	}
	sess := h.store.GetOrCreate(userID)
	sess.ProfileID = me.ID
	sess.ProfileEmail = me.Email
	sess.ProfileName = me.Name
	_ = h.store.Save()
	return h.present(ctx, chatID, userID, formatProfileCard(me.Name, me.Email), h.withBack(h.profileScreenKeyboard()), panel)
}

func (h *TelegramHandler) handleCallbackQuery(ctx context.Context, cq *tgCallbackQuery) error {
	if cq == nil || cq.From == nil || cq.Message == nil || cq.Message.Chat == nil {
		return nil
	}
	defer func() { _ = h.answerCallbackQuery(ctx, cq.ID, "") }()

	userID := cq.From.ID
	chatID := cq.Message.Chat.ID
	panel := cq.Message
	data := strings.TrimSpace(cq.Data)

	st := h.store.GetCopy(userID)
	token := st.Token

	switch {
	case strings.HasPrefix(data, "rcp:"):
		_ = h.store.SetState(userID, UserStateIdle)
		sess := h.store.GetOrCreate(userID)
		rest := strings.TrimPrefix(data, "rcp:")
		switch {
		case rest == "clearsearch":
			sess.RecipesSearchQuery = ""
		case rest == "reset":
			sess.RecipesSearchQuery = ""
			sess.RecipesFilterCategory = ""
			sess.RecipesSortBy = ""
			sess.RecipesSortDesc = false
		case strings.HasPrefix(rest, "category:"):
			v := strings.TrimPrefix(rest, "category:")
			if v == "*" {
				sess.RecipesFilterCategory = ""
			} else {
				sess.RecipesFilterCategory = v
			}
		case strings.HasPrefix(rest, "sort:"):
			switch strings.TrimPrefix(rest, "sort:") {
			case "name", "time", "category":
				sess.RecipesSortBy = strings.TrimPrefix(rest, "sort:")
			default:
				return nil
			}
		case strings.HasPrefix(rest, "order:"):
			sess.RecipesSortDesc = strings.TrimPrefix(rest, "order:") == "desc"
		default:
			return nil
		}
		sess.RecipesLastPage = 0
		_ = h.store.Save()
		return h.renderRecipesPage(ctx, userID, chatID, 0, panel)

	case strings.HasPrefix(data, "nav:"):
		_ = h.store.SetState(userID, UserStateIdle)
		sess := h.store.GetOrCreate(userID)

		nav := strings.TrimPrefix(data, "nav:")
		switch nav {
		case "menu":
			sess.ActiveNav = "menu"
			sess.BackNav = ""
			sess.BackPage = 0
			sess.AuthFlow = ""
			sess.Name = ""
			_ = h.store.Save()
			return h.present(ctx, chatID, userID, "✨ Главное меню\n\nВыберите, что сделать дальше 👇", h.mainMenuKeyboard(), panel)

		case "recipes":
			sess.ActiveNav = "recipes"
			sess.BackNav = "menu"
			sess.BackPage = 0
			sess.AuthFlow = ""
			sess.Name = ""
			sess.RecipesSearchQuery = ""
			sess.RecipesFilterCategory = ""
			sess.RecipesSortBy = ""
			sess.RecipesSortDesc = false
			sess.RecipesLastPage = 0
			_ = h.store.Save()
			return h.renderRecipesPage(ctx, userID, chatID, 0, panel)

		case "recipes_search":
			_ = h.store.Set(userID, UserStateAwaitingRecipeSearch, st.Email, token)
			kb := [][]tgInlineKeyboardButton{
				{{Text: "⬅️ К списку", CallbackData: "nav:rcp_search_cancel"}},
				h.compactMenuRow(),
			}
			return h.present(ctx, chatID, userID, "🔍 Введите слово из названия (по-русски или по-английски).", &tgInlineKeyboardMarkup{InlineKeyboard: kb}, panel)

		case "rcp_search_cancel":
			_ = h.store.Set(userID, UserStateIdle, st.Email, token)
			return h.renderRecipesPage(ctx, userID, chatID, 0, panel)

		case "favorites":
			if strings.TrimSpace(token) == "" {
				log.Printf("ERROR: empty token userID=%d action=nav:favorites", userID)
				return h.present(ctx, chatID, userID, "Сначала войдите", h.withBack(h.loginKeyboard()), panel)
			}
			sess.ActiveNav = "favorites"
			sess.BackNav = "menu"
			sess.BackPage = 0
			sess.AuthFlow = ""
			sess.Name = ""
			_ = h.store.Save()
			return h.renderFavoritesPage(ctx, userID, chatID, token, panel)

		case "match":
			sess.ActiveNav = "match"
			sess.BackNav = "menu"
			sess.BackPage = 0
			sess.AuthFlow = ""
			sess.Name = ""
			sess.MatchResults = nil
			sess.MatchLastPage = 0
			_ = h.store.Save()
			_ = h.store.Set(userID, UserStateAwaitingIngr, "", token)
			return h.present(ctx, chatID, userID, "📷 Пришлите фото продуктов\nили\n✏️ напишите список через запятую (например: помидор, сыр, яйцо)", h.withBack(h.mainMenuKeyboard()), panel)

		case "profile":
			sess.ActiveNav = "profile"
			sess.BackNav = "menu"
			sess.BackPage = 0
			sess.AuthFlow = ""
			sess.Name = ""
			_ = h.store.Save()

			if strings.TrimSpace(token) == "" {
				log.Printf("ERROR: empty token userID=%d action=nav:profile", userID)
				return h.present(ctx, chatID, userID, "Вы не авторизованы", h.withBack(h.loginKeyboard()), panel)
			}
			return h.renderProfileCard(ctx, chatID, userID, token, panel)

		case "logout":
			_ = h.store.Set(userID, UserStateIdle, "", "")
			sess.ProfileID = 0
			sess.ProfileEmail = ""
			sess.ProfileName = ""
			sess.ActiveNav = "menu"
			sess.BackNav = ""
			sess.BackPage = 0
			sess.AuthFlow = ""
			sess.Name = ""
			_ = h.store.Save()
			return h.present(ctx, chatID, userID, "Вы вышли из аккаунта. При необходимости войдите снова 👋", h.mainMenuKeyboard(), panel)

		case "login":
			sess.AuthFlow = "login"
			sess.Name = ""
			_ = h.store.Save()
			_ = h.store.Set(userID, UserStateAwaitingEmail, "", token)
			return h.present(ctx, chatID, userID, "Введите email", h.withBack(nil), panel)

		case "register":
			sess.AuthFlow = "register"
			sess.Name = ""
			_ = h.store.Save()
			_ = h.store.Set(userID, UserStateAwaitingRegisterName, "", token)
			return h.present(ctx, chatID, userID, "Введите имя", h.withBack(nil), panel)

		case "back":
			target := strings.TrimSpace(sess.BackNav)
			switch target {
			case "recipes":
				return h.renderRecipesPage(ctx, userID, chatID, sess.RecipesLastPage, panel)
			case "favorites":
				if strings.TrimSpace(token) == "" {
					return h.present(ctx, chatID, userID, "Сначала войдите", h.withBack(h.loginKeyboard()), panel)
				}
				return h.renderFavoritesPage(ctx, userID, chatID, token, panel)
			case "match":
				if len(sess.MatchResults) > 0 {
					return h.renderMatchesPage(ctx, chatID, userID, sess.BackPage, panel)
				}
				_ = h.store.Set(userID, UserStateAwaitingIngr, "", token)
				return h.present(ctx, chatID, userID, "📷 Пришлите фото продуктов\nили\n✏️ напишите список через запятую (например: помидор, сыр, яйцо)", h.withBack(h.mainMenuKeyboard()), panel)
			case "profile":
				return h.renderProfileCard(ctx, chatID, userID, token, panel)
			default:
				sess.ActiveNav = "menu"
				sess.BackNav = ""
				sess.BackPage = 0
				_ = h.store.Save()
				return h.present(ctx, chatID, userID, "✨ Главное меню\n\nВыберите раздел 👇", h.mainMenuKeyboard(), panel)
			}
		case "recipes_next":
			next := sess.RecipesLastPage + 1
			sess.RecipesLastPage = next
			sess.BackNav = "menu"
			sess.BackPage = 0
			_ = h.store.Save()
			return h.renderRecipesPage(ctx, userID, chatID, next, panel)
		case "recipes_prev":
			prev := sess.RecipesLastPage - 1
			if prev < 0 {
				prev = 0
			}
			sess.RecipesLastPage = prev
			sess.BackNav = "menu"
			sess.BackPage = 0
			_ = h.store.Save()
			return h.renderRecipesPage(ctx, userID, chatID, prev, panel)
		case "match_next":
			next := sess.MatchLastPage + 1
			perPage := matchResultsPerPage
			totalPages := (len(sess.MatchResults) + perPage - 1) / perPage
			if len(sess.MatchResults) == 0 || next >= totalPages {
				return nil
			}
			sess.MatchLastPage = next
			_ = h.store.Save()
			return h.renderMatchesPage(ctx, chatID, userID, next, panel)
		case "match_prev":
			prev := sess.MatchLastPage - 1
			if prev < 0 || len(sess.MatchResults) == 0 {
				return nil
			}
			sess.MatchLastPage = prev
			_ = h.store.Save()
			return h.renderMatchesPage(ctx, chatID, userID, prev, panel)
		default:
			return nil
		}

	case strings.HasPrefix(data, "recipe:"):
		idStr := strings.TrimPrefix(data, "recipe:")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			return h.present(ctx, chatID, userID, "Не удалось открыть рецепт. Вернитесь в список и выберите снова 🙏", h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: [][]tgInlineKeyboardButton{h.compactMenuRow()}}), panel)
		}

		sess := h.store.GetOrCreate(userID)
		switch sess.ActiveNav {
		case "recipes":
			sess.BackNav = "recipes"
			sess.BackPage = sess.RecipesLastPage
		case "favorites":
			sess.BackNav = "favorites"
			sess.BackPage = sess.FavsLastPage
		case "match":
			sess.BackNav = "match"
			sess.BackPage = sess.MatchLastPage
		default:
			sess.BackNav = "menu"
			sess.BackPage = 0
		}
		_ = h.store.Save()

		recipe, err := h.botService.GetRecipeDetails(ctx, id)
		if err != nil {
			return h.present(ctx, chatID, userID, "Не получилось открыть рецепт. Попробуйте ещё раз чуть позже 🙏", h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: [][]tgInlineKeyboardButton{h.compactMenuRow()}}), panel)
		}

		var b strings.Builder
		b.WriteString("🍽 ")
		b.WriteString(localizeRecipeTitle(recipe.Title))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("⏱ %d мин · 📂 %s", recipe.CookingTime, localizeCategory(recipe.Category)))

		if desc := strings.TrimSpace(recipe.Description); desc != "" {
			b.WriteString("\n\n📝 ")
			b.WriteString(localizedRecipeDescription(recipe.Title, desc))
		}

		if len(recipe.Ingredients) > 0 {
			b.WriteString("\n\n🥕 Ингредиенты:\n")
			for _, ing := range recipe.Ingredients {
				name := strings.TrimSpace(ing.Name)
				if name == "" {
					continue
				}
				b.WriteString("· ")
				b.WriteString(localizeIngredientDisplayName(name))
				if a := formatRecipeAmountRu(ing.Amount); a != "" {
					b.WriteString(" — ")
					b.WriteString(a)
				}
				b.WriteString("\n")
			}
		}

		stepsRu := localizedRecipeSteps(recipe.Title, recipe.Steps)
		if len(stepsRu) > 0 {
			b.WriteString("\n👩‍🍳 Как готовить:\n\n")
			stepNo := 0
			for _, s := range stepsRu {
				s = strings.TrimSpace(s)
				if s == "" {
					continue
				}
				stepNo++
				b.WriteString(fmt.Sprintf("%d) %s\n", stepNo, s))
			}
		} else {
			b.WriteString("\n\n🥘 Способ приготовления: уточняется")
		}

		favLabel := "❤️ В избранное"
		if strings.TrimSpace(token) != "" {
			if favs, err := h.botService.GetFavorites(ctx, token); err == nil {
				for _, fr := range favs {
					if fr.ID == id {
						favLabel = "💔 Убрать из избранного"
						break
					}
				}
			}
		}

		detailKB := [][]tgInlineKeyboardButton{
			{{Text: favLabel, CallbackData: fmt.Sprintf("fav:%d", id)}},
		}
		detailKB = append(detailKB, h.compactMenuRow())
		return h.present(ctx, chatID, userID, b.String(), h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: detailKB}), panel, recipe.Image)

	case strings.HasPrefix(data, "fav:"):
		idStr := strings.TrimPrefix(data, "fav:")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			return h.answerCallbackQuery(ctx, cq.ID, "Некорректный id")
		}

		if strings.TrimSpace(token) == "" {
			log.Printf("ERROR: empty token userID=%d action=fav:toggle", userID)
			return h.answerCallbackQuery(ctx, cq.ID, "Нужно авторизоваться")
		}

		if err := h.botService.AddFavorite(ctx, token, id); err != nil {
			var be *client.BackendError
			if errors.As(err, &be) && be.StatusCode == http.StatusConflict {
				if err2 := h.botService.RemoveFavorite(ctx, token, id); err2 != nil {
					return h.answerCallbackQuery(ctx, cq.ID, err2.Error())
				}
				return h.answerCallbackQuery(ctx, cq.ID, "Удалено из избранного ✅")
			}
			return h.answerCallbackQuery(ctx, cq.ID, err.Error())
		}

		return h.answerCallbackQuery(ctx, cq.ID, "Добавлено в избранное ✅")
	}

	return nil
}

func (h *TelegramHandler) filterRecipeCatalog(all []dto.RecipeResponse, sess *UserSession) []dto.RecipeResponse {
	out := append([]dto.RecipeResponse(nil), all...)
	q := strings.TrimSpace(strings.ToLower(sess.RecipesSearchQuery))
	if q != "" {
		next := out[:0]
		for _, r := range out {
			titleRu := strings.ToLower(localizeRecipeTitle(r.Title))
			titleEn := strings.ToLower(strings.TrimSpace(r.Title))
			if strings.Contains(titleRu, q) || strings.Contains(titleEn, q) {
				next = append(next, r)
			}
		}
		out = next
	}
	if cat := strings.TrimSpace(sess.RecipesFilterCategory); cat != "" {
		next := out[:0]
		for _, r := range out {
			if strings.EqualFold(strings.TrimSpace(r.Category), cat) {
				next = append(next, r)
			}
		}
		out = next
	}
	sortBy := strings.TrimSpace(strings.ToLower(sess.RecipesSortBy))
	if sortBy == "" {
		sortBy = "name"
	}
	desc := sess.RecipesSortDesc
	sort.SliceStable(out, func(i, j int) bool {
		switch sortBy {
		case "time":
			ti, tj := out[i].CookingTime, out[j].CookingTime
			if desc {
				return ti > tj
			}
			return ti < tj
		case "category":
			ci := localizeCategory(out[i].Category)
			cj := localizeCategory(out[j].Category)
			if desc {
				return ci > cj
			}
			return ci < cj
		default:
			ti := strings.TrimSpace(out[i].Title)
			tj := strings.TrimSpace(out[j].Title)
			if desc {
				return ti > tj
			}
			return ti < tj
		}
	})
	return out
}

func paginateRecipePage(recipes []dto.RecipeResponse, page, limit int) (slice []dto.RecipeResponse, hasNext bool) {
	if page < 0 {
		page = 0
	}
	if limit <= 0 {
		limit = 5
	}
	start := page * limit
	if start >= len(recipes) {
		return nil, false
	}
	end := start + limit
	if end > len(recipes) {
		end = len(recipes)
	}
	slice = append([]dto.RecipeResponse(nil), recipes[start:end]...)
	hasNext = end < len(recipes)
	return slice, hasNext
}

func (h *TelegramHandler) recipeCatalogControlRows(sess *UserSession) [][]tgInlineKeyboardButton {
	sortBy := strings.TrimSpace(strings.ToLower(sess.RecipesSortBy))
	if sortBy == "" {
		sortBy = "name"
	}
	markSort := func(key, label string) string {
		if sortBy == key {
			return "✓ " + label
		}
		return label
	}
	markCat := func(cat, label string) string {
		if cat == "" {
			if sess.RecipesFilterCategory == "" {
				return "✓ " + label
			}
			return label
		}
		if strings.EqualFold(sess.RecipesFilterCategory, cat) {
			return "✓ " + label
		}
		return label
	}
	ascLabel, descLabel := "↑ Возр.", "↓ Убыв."
	if !sess.RecipesSortDesc {
		ascLabel = "✓ " + ascLabel
	} else {
		descLabel = "✓ " + descLabel
	}

	rows := [][]tgInlineKeyboardButton{
		{
			{Text: "🔍 Название", CallbackData: "nav:recipes_search"},
		},
	}
	if strings.TrimSpace(sess.RecipesSearchQuery) != "" {
		rows = append(rows, []tgInlineKeyboardButton{
			{Text: "✖ Очистить поиск", CallbackData: "rcp:clearsearch"},
		})
	}
	rows = append(rows,
		[]tgInlineKeyboardButton{
			{Text: markCat("", "Все"), CallbackData: "rcp:category:*"},
			{Text: markCat("pasta", "Паста"), CallbackData: "rcp:category:pasta"},
			{Text: markCat("meat", "Мясо"), CallbackData: "rcp:category:meat"},
			{Text: markCat("vegetarian", "Вегет."), CallbackData: "rcp:category:vegetarian"},
		},
		[]tgInlineKeyboardButton{
			{Text: markCat("breakfast", "Завтрак"), CallbackData: "rcp:category:breakfast"},
			{Text: markCat("dessert", "Десерт"), CallbackData: "rcp:category:dessert"},
			{Text: markCat("soup", "Суп"), CallbackData: "rcp:category:soup"},
			{Text: markCat("salad", "Салат"), CallbackData: "rcp:category:salad"},
		},
		[]tgInlineKeyboardButton{
			{Text: markCat("dinner", "Ужин"), CallbackData: "rcp:category:dinner"},
			{Text: markCat("lunch", "Обед"), CallbackData: "rcp:category:lunch"},
		},
		[]tgInlineKeyboardButton{
			{Text: markSort("name", "Имя"), CallbackData: "rcp:sort:name"},
			{Text: markSort("time", "Время"), CallbackData: "rcp:sort:time"},
			{Text: markSort("category", "Раздел"), CallbackData: "rcp:sort:category"},
		},
		[]tgInlineKeyboardButton{
			{Text: ascLabel, CallbackData: "rcp:order:asc"},
			{Text: descLabel, CallbackData: "rcp:order:desc"},
			{Text: "↺ Сброс", CallbackData: "rcp:reset"},
		},
	)
	return rows
}

func (h *TelegramHandler) renderRecipesPage(ctx context.Context, userID int64, chatID int64, page int, panel *tgCallbackQueryMsg) error {
	const limit = 5

	all, err := h.botService.GetRecipes(ctx)
	if err != nil {
		return h.present(ctx, chatID, userID, "Не удалось загрузить рецепты. Проверьте связь и нажмите «📖 Рецепты» ещё раз 🙏", h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: [][]tgInlineKeyboardButton{h.compactMenuRow()}}), panel)
	}

	sess := h.store.GetOrCreate(userID)
	filtered := h.filterRecipeCatalog(all, sess)
	slice, hasNext := paginateRecipePage(filtered, page, limit)

	sess.ActiveNav = "recipes"
	sess.RecipesLastPage = page
	sess.BackNav = "menu"
	sess.BackPage = 0
	_ = h.store.Save()

	return h.renderRecipeListPageRecipesPaged(ctx, chatID, "Рецепты", len(filtered), slice, page, hasNext, userID, panel)
}

func (h *TelegramHandler) renderFavoritesPage(ctx context.Context, userID int64, chatID int64, token string, panel *tgCallbackQueryMsg) error {
	recipes, err := h.botService.GetFavorites(ctx, token)
	if err != nil {
		return h.present(ctx, chatID, userID, "Избранное не загрузилось. Проверьте вход и попробуйте ещё раз 💛", h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: [][]tgInlineKeyboardButton{h.compactMenuRow()}}), panel)
	}
	sess := h.store.GetOrCreate(userID)
	sess.ActiveNav = "favorites"
	sess.FavsLastPage = 0
	sess.BackNav = "menu"
	sess.BackPage = 0
	_ = h.store.Save()
	return h.renderRecipeListPage(ctx, chatID, "Избранное", recipes, userID, panel)
}

func (h *TelegramHandler) renderRecipeListPage(ctx context.Context, chatID int64, title string, recipes []dto.RecipeResponse, userID int64, panel *tgCallbackQueryMsg) error {
	if len(recipes) == 0 {
		msg := "😔 Пока здесь пусто\n\nЗагляните в «🔍 Подобрать рецепт» — подберём блюдо по продуктам ✨"
		return h.present(ctx, chatID, userID, msg, h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: [][]tgInlineKeyboardButton{h.compactMenuRow()}}), panel)
	}

	const maxShown = 10
	slice := recipes
	if len(slice) > maxShown {
		slice = slice[:maxShown]
	}

	var b strings.Builder
	b.WriteString(title)
	b.WriteString("\n\n")

	rows := make([][]tgInlineKeyboardButton, 0, len(slice)+6)
	for _, r := range slice {
		b.WriteString("🍝 ")
		b.WriteString(localizeRecipeTitle(r.Title))
		b.WriteString("\n")
		b.WriteString("⏱ ")
		b.WriteString(strconv.Itoa(r.CookingTime))
		b.WriteString(" мин\n")
		b.WriteString("📂 ")
		b.WriteString(localizeCategory(r.Category))
		b.WriteString("\n\n")

		listRow := []tgInlineKeyboardButton{
			{Text: "Подробнее", CallbackData: fmt.Sprintf("recipe:%d", r.ID)},
			{Text: "❤️", CallbackData: fmt.Sprintf("fav:%d", r.ID)},
		}
		rows = append(rows, listRow)
	}

	rows = append(rows, h.compactMenuRow())
	kb := &tgInlineKeyboardMarkup{InlineKeyboard: rows}
	return h.present(ctx, chatID, userID, strings.TrimSpace(b.String()), h.withBack(kb), panel)
}

func (h *TelegramHandler) renderRecipeListPageRecipesPaged(ctx context.Context, chatID int64, title string, totalFiltered int, recipes []dto.RecipeResponse, page int, hasNext bool, userID int64, panel *tgCallbackQueryMsg) error {
	sess := h.store.GetCopy(userID)
	if len(recipes) == 0 {
		if page > 0 {
			msg := "😔 На этой странице пусто — вернитесь назад или в меню"
			rows := [][]tgInlineKeyboardButton{
				{{Text: "⬅️ Предыдущая", CallbackData: "nav:recipes_prev"}},
			}
			if title == "Рецепты" {
				rows = append(rows, h.recipeCatalogControlRows(&sess)...)
			}
			rows = append(rows, h.compactMenuRow())
			return h.present(ctx, chatID, userID, msg, h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: rows}), panel)
		}
		msg := "😔 Рецепты не найдены\n\nПопробуйте изменить поиск или фильтры — или зайдите в подбор по продуктам 🔍"
		rows := [][]tgInlineKeyboardButton{}
		if title == "Рецепты" {
			rows = append(rows, h.recipeCatalogControlRows(&sess)...)
		}
		rows = append(rows, h.compactMenuRow())
		return h.present(ctx, chatID, userID, msg, h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: rows}), panel)
	}

	var b strings.Builder
	b.WriteString(title)
	b.WriteString(fmt.Sprintf("\nНайдено: %d · стр. %d", totalFiltered, page+1))
	if title == "Рецепты" {
		if q := strings.TrimSpace(sess.RecipesSearchQuery); q != "" {
			b.WriteString("\n🔍 ")
			b.WriteString(q)
		}
		if c := strings.TrimSpace(sess.RecipesFilterCategory); c != "" {
			b.WriteString("\n📂 ")
			b.WriteString(localizeCategory(c))
		}
	}
	b.WriteString("\n\n")

	rows := make([][]tgInlineKeyboardButton, 0, len(recipes)+6)
	for _, r := range recipes {
		b.WriteString("🍝 ")
		b.WriteString(localizeRecipeTitle(r.Title))
		b.WriteString("\n")
		b.WriteString("⏱ ")
		b.WriteString(strconv.Itoa(r.CookingTime))
		b.WriteString(" мин\n")
		b.WriteString("📂 ")
		b.WriteString(localizeCategory(r.Category))
		b.WriteString("\n\n")

		pageRow := []tgInlineKeyboardButton{
			{Text: "Подробнее", CallbackData: fmt.Sprintf("recipe:%d", r.ID)},
			{Text: "❤️", CallbackData: fmt.Sprintf("fav:%d", r.ID)},
		}
		rows = append(rows, pageRow)
	}

	navRow := make([]tgInlineKeyboardButton, 0, 2)
	if page > 0 {
		navRow = append(navRow, tgInlineKeyboardButton{Text: "⬅️ Назад", CallbackData: "nav:recipes_prev"})
	}
	if hasNext {
		navRow = append(navRow, tgInlineKeyboardButton{Text: "➡️ Далее", CallbackData: "nav:recipes_next"})
	}
	if len(navRow) > 0 {
		rows = append(rows, navRow)
	}

	if title == "Рецепты" {
		rows = append(rows, h.recipeCatalogControlRows(&sess)...)
	}

	rows = append(rows, h.compactMenuRow())
	kb := &tgInlineKeyboardMarkup{InlineKeyboard: rows}
	return h.present(ctx, chatID, userID, strings.TrimSpace(b.String()), h.withBack(kb), panel)
}

func (h *TelegramHandler) getUpdates(ctx context.Context, offset int, timeoutSec int) ([]tgUpdate, error) {
	endpoint := h.apiBaseURL + "/getUpdates"
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse getUpdates url: %w", err)
	}

	q := u.Query()
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	if timeoutSec > 0 {
		q.Set("timeout", strconv.Itoa(timeoutSec))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create getUpdates request: %w", err)
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getUpdates request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return nil, fmt.Errorf("getUpdates failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	var out struct {
		OK     bool       `json:"ok"`
		Result []tgUpdate `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode getUpdates: %w", err)
	}
	if !out.OK {
		return nil, fmt.Errorf("getUpdates: ok=false")
	}

	return out.Result, nil
}

func (h *TelegramHandler) getFilePath(ctx context.Context, fileID string) (string, error) {
	if strings.TrimSpace(fileID) == "" {
		return "", fmt.Errorf("file_id is empty")
	}

	endpoint := h.apiBaseURL + "/getFile"
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("parse getFile url: %w", err)
	}
	q := u.Query()
	q.Set("file_id", fileID)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("create getFile request: %w", err)
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("getFile request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return "", fmt.Errorf("getFile failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	var out struct {
		OK     bool `json:"ok"`
		Result struct {
			FilePath string `json:"file_path"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode getFile: %w", err)
	}
	if !out.OK {
		return "", fmt.Errorf("getFile: ok=false")
	}
	if strings.TrimSpace(out.Result.FilePath) == "" {
		return "", fmt.Errorf("getFile: empty file_path")
	}

	return out.Result.FilePath, nil
}

func (h *TelegramHandler) downloadFile(ctx context.Context, filePath string) ([]byte, error) {
	filePath = strings.TrimPrefix(filePath, "/")
	if strings.TrimSpace(filePath) == "" {
		return nil, fmt.Errorf("file_path is empty")
	}

	endpoint := h.fileBaseURL + "/" + filePath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create download request: %w", err)
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download file request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return nil, fmt.Errorf("download file failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	return io.ReadAll(resp.Body)
}

func (h *TelegramHandler) present(ctx context.Context, chatID int64, userID int64, text string, kb *tgInlineKeyboardMarkup, panel *tgCallbackQueryMsg, photoURL ...string) error {
	photo := ""
	if len(photoURL) > 0 {
		photo = strings.TrimSpace(photoURL[0])
	}
	text = strings.TrimSpace(text)
	if text == "" && photo == "" {
		text = "\u00a0"
	}

	var oldMid int
	if panel != nil && panel.Chat != nil && panel.MessageID > 0 {
		oldMid = panel.MessageID
	}
	if oldMid == 0 {
		st := h.store.GetOrCreate(userID)
		if st.PanelMessageID > 0 && st.PanelChatID == chatID {
			oldMid = st.PanelMessageID
		}
	}
	if oldMid > 0 {
		_ = h.deleteMessage(ctx, chatID, oldMid)
	}

	var newMid int
	var err error
	if photo != "" {
		cap := truncateTelegramCaption(text, 1024)
		newMid, err = h.sendPhotoReturnID(ctx, chatID, photo, cap, kb)
		if err != nil {
			log.Printf("present: sendPhoto failed, fallback to text: %v", err)
			if text == "" {
				text = "\u00a0"
			}
			newMid, err = h.sendMessageReturnID(ctx, chatID, text, kb)
		}
	} else {
		if text == "" {
			text = "\u00a0"
		}
		newMid, err = h.sendMessageReturnID(ctx, chatID, text, kb)
	}
	if err != nil {
		return err
	}
	return h.store.SetPanelMessage(userID, chatID, newMid)
}

func truncateTelegramCaption(s string, maxRunes int) string {
	s = strings.TrimSpace(s)
	r := []rune(s)
	if len(r) <= maxRunes {
		return s
	}
	cut := maxRunes - 1
	for cut > 0 && r[cut] != ' ' {
		cut--
	}
	if cut < maxRunes/2 {
		cut = maxRunes - 1
	}
	return strings.TrimSpace(string(r[:cut])) + "…"
}

func (h *TelegramHandler) deleteMessage(ctx context.Context, chatID int64, messageID int) error {
	if messageID <= 0 {
		return nil
	}
	endpoint := h.apiBaseURL + "/deleteMessage"
	payload := struct {
		ChatID    int64 `json:"chat_id"`
		MessageID int   `json:"message_id"`
	}{
		ChatID:    chatID,
		MessageID: messageID,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal deleteMessage: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("create deleteMessage request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("deleteMessage request: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("deleteMessage failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}
	var out struct {
		OK bool `json:"ok"`
	}
	if err := json.Unmarshal(bodyBytes, &out); err != nil {
		return fmt.Errorf("decode deleteMessage: %w", err)
	}
	if !out.OK {
		return fmt.Errorf("deleteMessage: ok=false body=%s", strings.TrimSpace(string(bodyBytes)))
	}
	return nil
}

func (h *TelegramHandler) sendPhotoReturnID(ctx context.Context, chatID int64, photo string, caption string, kb *tgInlineKeyboardMarkup) (int, error) {
	endpoint := h.apiBaseURL + "/sendPhoto"
	payload := struct {
		ChatID      int64                   `json:"chat_id"`
		Photo       string                  `json:"photo"`
		Caption     string                  `json:"caption,omitempty"`
		ReplyMarkup *tgInlineKeyboardMarkup `json:"reply_markup,omitempty"`
	}{
		ChatID:      chatID,
		Photo:       photo,
		Caption:     caption,
		ReplyMarkup: kb,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("marshal sendPhoto: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return 0, fmt.Errorf("create sendPhoto request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("sendPhoto request: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("sendPhoto failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}
	var out struct {
		OK     bool `json:"ok"`
		Result struct {
			MessageID int `json:"message_id"`
		} `json:"result"`
	}
	if err := json.Unmarshal(bodyBytes, &out); err != nil {
		return 0, fmt.Errorf("decode sendPhoto: %w", err)
	}
	if !out.OK {
		return 0, fmt.Errorf("sendPhoto: ok=false body=%s", strings.TrimSpace(string(bodyBytes)))
	}
	if out.Result.MessageID <= 0 {
		return 0, fmt.Errorf("sendPhoto: no message_id in response")
	}
	return out.Result.MessageID, nil
}

func (h *TelegramHandler) sendMessageReturnID(ctx context.Context, chatID int64, text string, kb *tgInlineKeyboardMarkup) (int, error) {
	endpoint := h.apiBaseURL + "/sendMessage"
	payload := struct {
		ChatID      int64                   `json:"chat_id"`
		Text        string                  `json:"text"`
		ReplyMarkup *tgInlineKeyboardMarkup `json:"reply_markup,omitempty"`
	}{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: kb,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("marshal sendMessage payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return 0, fmt.Errorf("create sendMessage request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("sendMessage request: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("sendMessage failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}
	var out struct {
		OK     bool `json:"ok"`
		Result struct {
			MessageID int `json:"message_id"`
		} `json:"result"`
	}
	if err := json.Unmarshal(bodyBytes, &out); err != nil {
		return 0, fmt.Errorf("decode sendMessage: %w", err)
	}
	if !out.OK {
		return 0, fmt.Errorf("sendMessage: ok=false body=%s", strings.TrimSpace(string(bodyBytes)))
	}
	if out.Result.MessageID <= 0 {
		return 0, fmt.Errorf("sendMessage: no message_id in response")
	}
	return out.Result.MessageID, nil
}

type tgUpdate struct {
	UpdateID      int              `json:"update_id"`
	Message       *tgMessage       `json:"message,omitempty"`
	CallbackQuery *tgCallbackQuery `json:"callback_query,omitempty"`
}

type tgCallbackQuery struct {
	ID      string              `json:"id"`
	From    *tgUser             `json:"from"`
	Message *tgCallbackQueryMsg `json:"message"`
	Data    string              `json:"data"`
}

type tgCallbackQueryMsg struct {
	MessageID int     `json:"message_id"`
	Chat      *tgChat `json:"chat"`
}

type tgMessage struct {
	MessageID int           `json:"message_id"`
	From      *tgUser       `json:"from"`
	Chat      *tgChat       `json:"chat"`
	Text      string        `json:"text"`
	Photo     []tgPhotoSize `json:"photo"`
}

type tgUser struct {
	ID int64 `json:"id"`
}

type tgChat struct {
	ID int64 `json:"id"`
}

type tgPhotoSize struct {
	FileID   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size"`
}

type tgInlineKeyboardMarkup struct {
	InlineKeyboard [][]tgInlineKeyboardButton `json:"inline_keyboard"`
}

type tgInlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
	URL          string `json:"url,omitempty"`
}
