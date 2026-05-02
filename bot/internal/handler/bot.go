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
)

type UserSession struct {
	State UserState
	Email string
	Token string

	// AuthFlow is "login" or "register" when in auth FSM states.
	AuthFlow string
	Name     string

	// Persisted user profile info from GET /auth/me.
	ProfileID    int64
	ProfileEmail string
	ProfileName  string

	// Navigation context for nav:back.
	ActiveNav       string // recipes|favorites|match|profile|menu
	RecipesLastPage int
	FavsLastPage    int
	BackNav         string // recipes|favorites|match|profile|menu
	BackPage        int

	// Одно «окно» диалога: последнее сообщение бота с inline-клавиатурой (редактируем его).
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

// Одна кнопка «в меню» вместо дублирования четырёх пунктов под длинными списками.
func (h *TelegramHandler) compactMenuRow() []tgInlineKeyboardButton {
	return []tgInlineKeyboardButton{{Text: "🏠 В меню", CallbackData: "nav:menu"}}
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

// Частые английские слова в названиях рецептов → нейтральные русские подписи (демо).
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
	for _, r := range title {
		if r > 127 {
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
			// Не показываем сырой английский в демо.
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

func recipeDescriptionLooksEnglish(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	hasLetter := false
	for _, r := range s {
		if unicode.IsLetter(r) {
			hasLetter = true
		}
		if unicode.Is(unicode.Cyrillic, r) {
			return false
		}
	}
	return hasLetter
}

// Filled from seeds + extras in recipe_title_ru_generated.go (init).
var recipeTitleENToRU = make(map[string]string)

func parseIngredients(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	parts := strings.FieldsFunc(s, func(r rune) bool {
		switch r {
		case ' ', '\n', '\t', ',', ';':
			return true
		default:
			return false
		}
	})

	out := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		p = strings.ToLower(p)
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

func (h *TelegramHandler) sendMatches(ctx context.Context, chatID int64, userID int64, matches []dto.MatchResponse, panel *tgCallbackQueryMsg) error {
	if len(matches) == 0 {
		return h.present(ctx, chatID, userID, "😔 Пока нет рецептов под ваш набор продуктов.\n\nПопробуйте другие ингредиенты или фото получше 📷", h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: [][]tgInlineKeyboardButton{h.compactMenuRow()}}), panel)
	}

	const (
		maxShownMatches  = 5
		maxMissingToShow = 4
	)

	shown := matches
	if len(shown) > maxShownMatches {
		shown = shown[:maxShownMatches]
	}

	var b strings.Builder
	rows := make([][]tgInlineKeyboardButton, 0, len(shown)+4)
	for i, m := range shown {
		if i > 0 {
			b.WriteString("\n\n")
		}

		b.WriteString("🍽 ")
		b.WriteString(localizeRecipeTitle(m.Title))
		b.WriteString("\n")

		pct := int(m.MatchPercent * 100)
		b.WriteString("📊 ")
		b.WriteString(strconv.Itoa(pct))
		b.WriteString("% совпадение по продуктам\n")

		if len(m.MissingIngredients) == 0 {
			b.WriteString("✨ Можно приготовить прямо сейчас 🎉")
		} else {
			b.WriteString("❌ Не хватает: ")
			miss := m.MissingIngredients
			if len(miss) > maxMissingToShow {
				miss = miss[:maxMissingToShow]
			}
			b.WriteString(strings.Join(miss, ", "))
			if remain := len(m.MissingIngredients) - len(miss); remain > 0 {
				b.WriteString(fmt.Sprintf(" и ещё %d", remain))
			}
		}

		rows = append(rows, []tgInlineKeyboardButton{
			{Text: "Подробнее", CallbackData: fmt.Sprintf("recipe:%d", m.RecipeID)},
			{Text: "❤️", CallbackData: fmt.Sprintf("fav:%d", m.RecipeID)},
		})
	}

	if remain := len(matches) - len(shown); remain > 0 {
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("… и ещё вариантов: %d", remain))
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

// RunLongPolling starts polling updates via getUpdates and handles:
// - /start
// - photos
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
			// transient network errors shouldn't kill the bot loop
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
		// /start must not wipe JWT/session; reset only FSM state.
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

	// In idle, we still must respond to the user with a helpful navigation.
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
			return h.present(ctx, chatID, userID, "Нажмите кнопку ниже или напишите «меню».\n\n✨ Главное меню 👇", h.mainMenuKeyboard(), nil)
		}
	}

	switch st.State {
	case UserStateAwaitingRegisterName:
		name := strings.TrimSpace(text)
		if name == "" {
			return h.present(ctx, chatID, userID, "Введите имя", h.withBack(nil), nil)
		}
		sess := h.store.GetOrCreate(userID)
		sess.Name = name
		sess.AuthFlow = "register"
		_ = h.store.Save()
		_ = h.store.Set(userID, UserStateAwaitingRegisterEmail, "", st.Token)
		return h.present(ctx, chatID, userID, "Введите email", h.withBack(nil), nil)

	case UserStateAwaitingRegisterEmail:
		email := strings.TrimSpace(text)
		if email == "" {
			return h.present(ctx, chatID, userID, "Введите email", h.withBack(nil), nil)
		}
		_ = h.store.Set(userID, UserStateAwaitingRegisterPassword, email, st.Token)
		return h.present(ctx, chatID, userID, "Введите пароль", h.withBack(nil), nil)

	case UserStateAwaitingRegisterPassword:
		password := strings.TrimSpace(text)
		if password == "" {
			return h.present(ctx, chatID, userID, "Введите пароль", h.withBack(nil), nil)
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
		return h.present(ctx, chatID, userID, formatProfileCard(me.Email, me.ID), h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: [][]tgInlineKeyboardButton{h.compactMenuRow()}}), nil)

	case UserStateAwaitingEmail:
		email := strings.TrimSpace(text)
		if email == "" {
			return h.present(ctx, chatID, userID, "Введите email", h.withBack(nil), nil)
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

		// 1) Save token in session immediately after successful auth (login/register).
		_ = h.store.Set(userID, UserStateIdle, email, token)

		// 3) Immediately call GET /auth/me.
		me, err := h.botService.Me(ctx, token)
		if err != nil {
			_ = h.store.Set(userID, UserStateIdle, "", token)
			sess.AuthFlow = ""
			sess.Name = ""
			_ = h.store.Save()
			return h.present(ctx, chatID, userID, err.Error(), h.withBack(nil), nil)
		}

		// 4) Save user in session.
		sess.ProfileID = me.ID
		sess.ProfileEmail = me.Email
		sess.ProfileName = me.Name

		_ = h.store.Set(userID, UserStateIdle, me.Email, token)
		sess.AuthFlow = ""
		sess.Name = ""
		_ = h.store.Save()
		return h.present(ctx, chatID, userID, formatProfileCard(me.Email, me.ID), h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: [][]tgInlineKeyboardButton{h.compactMenuRow()}}), nil)

	case UserStateAwaitingIngr:
		ingredients := parseIngredients(text)
		matches, err := h.botService.Match(ctx, ingredients)
		_ = h.store.Set(userID, UserStateIdle, "", st.Token)
		if err != nil {
			return h.present(ctx, chatID, userID, err.Error(), h.withBack(h.mainMenuKeyboard()), nil)
		}
		return h.sendMatches(ctx, chatID, userID, matches, nil)

	default:
		// Старые сохранённые сессии: одношаговая форма регистрации больше не используется.
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

func formatProfileCard(email string, id int64) string {
	return fmt.Sprintf("👤 Профиль\n\n📧 Email: %s\n🆔 ID: %d", email, id)
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
	return h.present(ctx, chatID, userID, formatProfileCard(me.Email, me.ID), h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: [][]tgInlineKeyboardButton{h.compactMenuRow()}}), panel)
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

	// REQUIRED: strict routing by callback prefix.
	switch {
	case strings.HasPrefix(data, "nav:"):
		// Any nav action resets FSM to idle first (without losing token/session).
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
			_ = h.store.Save()
			sess.RecipesLastPage = 0
			_ = h.store.Save()
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
			_ = h.store.Save()
			// Set awaiting ingredients after nav reset.
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

		case "login":
			// Start LOGIN: email -> password.
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
			sess.BackPage = 0
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
			b.WriteString("\n\n")
			if recipeDescriptionLooksEnglish(desc) {
				b.WriteString("📝 Суть блюда: сбалансированный вкус и простая подача — идеально для домашнего стола ✨")
			} else {
				b.WriteString("📝 ")
				b.WriteString(desc)
			}
		}

		if len(recipe.Steps) > 0 {
			b.WriteString("\n\n👩‍🍳 Как готовить:\n\n")
			stepNo := 0
			for _, s := range recipe.Steps {
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
		return h.present(ctx, chatID, userID, b.String(), h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: [][]tgInlineKeyboardButton{h.compactMenuRow()}}), panel)

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

		// Toggle делаем через conflict: AddFavorite -> conflict => RemoveFavorite.
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

func (h *TelegramHandler) renderRecipesPage(ctx context.Context, userID int64, chatID int64, page int, panel *tgCallbackQueryMsg) error {
	const limit = 10

	recipes, supported, err := h.botService.GetRecipesPaged(ctx, page, limit)
	if err != nil {
		return h.present(ctx, chatID, userID, "Не удалось загрузить рецепты. Проверьте связь и нажмите «📖 Рецепты» ещё раз 🙏", h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: [][]tgInlineKeyboardButton{h.compactMenuRow()}}), panel)
	}

	// If backend doesn't support paging: show only first 10, no pagination buttons.
	if !supported {
		all, err := h.botService.GetRecipes(ctx)
		if err != nil {
			return h.present(ctx, chatID, userID, "Не удалось загрузить рецепты. Проверьте связь и попробуйте снова 🙏", h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: [][]tgInlineKeyboardButton{h.compactMenuRow()}}), panel)
		}
		if len(all) > limit {
			all = all[:limit]
		}
		sess := h.store.GetOrCreate(userID)
		sess.ActiveNav = "recipes"
		sess.RecipesLastPage = 0
		sess.BackNav = "menu"
		sess.BackPage = 0
		_ = h.store.Save()
		return h.renderRecipeListPage(ctx, chatID, "Рецепты", all, userID, panel)
	}

	// backend paging supported
	sess := h.store.GetOrCreate(userID)
	sess.ActiveNav = "recipes"
	sess.RecipesLastPage = page
	sess.BackNav = "menu"
	sess.BackPage = 0
	_ = h.store.Save()

	return h.renderRecipeListPageRecipesPaged(ctx, chatID, "Рецепты", recipes, page, len(recipes) == limit, userID, panel)
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

	// Pagination callback formats are intentionally removed to comply with strict callback routing formats.
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

		rows = append(rows, []tgInlineKeyboardButton{
			{Text: "Подробнее", CallbackData: fmt.Sprintf("recipe:%d", r.ID)},
			{Text: "❤️", CallbackData: fmt.Sprintf("fav:%d", r.ID)},
		})
	}

	rows = append(rows, h.compactMenuRow())
	kb := &tgInlineKeyboardMarkup{InlineKeyboard: rows}
	return h.present(ctx, chatID, userID, strings.TrimSpace(b.String()), h.withBack(kb), panel)
}

func (h *TelegramHandler) renderRecipeListPageRecipesPaged(ctx context.Context, chatID int64, title string, recipes []dto.RecipeResponse, page int, hasNext bool, userID int64, panel *tgCallbackQueryMsg) error {
	if len(recipes) == 0 {
		if page > 0 {
			msg := "😔 На этой странице пусто — вернитесь назад или в меню"
			rows := [][]tgInlineKeyboardButton{
				{{Text: "⬅️ Предыдущая", CallbackData: "nav:recipes_prev"}},
				h.compactMenuRow(),
			}
			return h.present(ctx, chatID, userID, msg, h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: rows}), panel)
		}
		msg := "😔 Рецепты не найдены\n\nПопробуйте позже или зайдите в подбор по продуктам 🔍"
		return h.present(ctx, chatID, userID, msg, h.withBack(&tgInlineKeyboardMarkup{InlineKeyboard: [][]tgInlineKeyboardButton{h.compactMenuRow()}}), panel)
	}

	var b strings.Builder
	b.WriteString(title)
	b.WriteString(fmt.Sprintf(" (стр. %d)\n\n", page+1))

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

		rows = append(rows, []tgInlineKeyboardButton{
			{Text: "Подробнее", CallbackData: fmt.Sprintf("recipe:%d", r.ID)},
			{Text: "❤️", CallbackData: fmt.Sprintf("fav:%d", r.ID)},
		})
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

// present обновляет одно и то же сообщение бота: при callback — правка сообщения с кнопкой;
// при вводе текста/фото — правка по сохранённому message_id, иначе первая отправка.
func (h *TelegramHandler) present(ctx context.Context, chatID int64, userID int64, text string, kb *tgInlineKeyboardMarkup, panel *tgCallbackQueryMsg) error {
	text = strings.TrimSpace(text)
	if text == "" {
		text = "\u00a0"
	}

	tryEdit := func(cid int64, mid int) bool {
		if mid <= 0 || cid == 0 {
			return false
		}
		if err := h.editMessageText(ctx, cid, mid, text, kb); err != nil {
			log.Printf("present: edit failed chat=%d msg=%d: %v", cid, mid, err)
			return false
		}
		return true
	}

	if panel != nil && panel.Chat != nil && tryEdit(panel.Chat.ID, panel.MessageID) {
		return h.store.SetPanelMessage(userID, panel.Chat.ID, panel.MessageID)
	}

	st := h.store.GetOrCreate(userID)
	if st.PanelMessageID > 0 && st.PanelChatID == chatID && tryEdit(chatID, st.PanelMessageID) {
		return nil
	}

	mid, err := h.sendMessageReturnID(ctx, chatID, text, kb)
	if err != nil {
		return err
	}
	return h.store.SetPanelMessage(userID, chatID, mid)
}

func (h *TelegramHandler) editMessageText(ctx context.Context, chatID int64, messageID int, text string, kb *tgInlineKeyboardMarkup) error {
	endpoint := h.apiBaseURL + "/editMessageText"
	payload := struct {
		ChatID      int64                   `json:"chat_id"`
		MessageID   int                     `json:"message_id"`
		Text        string                  `json:"text"`
		ReplyMarkup *tgInlineKeyboardMarkup `json:"reply_markup,omitempty"`
	}{
		ChatID:      chatID,
		MessageID:   messageID,
		Text:        text,
		ReplyMarkup: kb,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal editMessageText: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("create editMessageText request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("editMessageText request: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
	low := strings.ToLower(string(bodyBytes))
	if strings.Contains(low, "message is not modified") {
		return nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("editMessageText failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}
	var out struct {
		OK bool `json:"ok"`
	}
	if err := json.Unmarshal(bodyBytes, &out); err != nil {
		return fmt.Errorf("decode editMessageText: %w", err)
	}
	if !out.OK && !strings.Contains(low, "message is not modified") {
		return fmt.Errorf("editMessageText: ok=false body=%s", strings.TrimSpace(string(bodyBytes)))
	}
	return nil
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
	CallbackData string `json:"callback_data"`
}
