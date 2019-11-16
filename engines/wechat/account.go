package wechat

// Account ...
type Account struct {
	WechatID  string
	Email     string
	Favorites map[string]bool
	History   map[string]bool
}

// NewAccount ...
func NewAccount(userID string) *Account {
	return &Account{
		WechatID:  userID,
		Favorites: make(map[string]bool),
		History:   make(map[string]bool),
	}
}

// AddFavor 将加入收藏夹
func (a *Account) AddFavor(bookID string) {
	a.Favorites[bookID] = true
}

// GetFavor 判断书否已收藏
func (a *Account) GetFavor(bookID string) bool {
	if b, found := a.Favorites[bookID]; found {
		return b
	}

	return false
}

// AddHistory 增加到下载历史
func (a *Account) AddHistory(bookID string) {
	a.History[bookID] = true
}

// GetHistory ...
func (a *Account) GetHistory(bookID string) bool {
	if b, found := a.History[bookID]; found {
		return b
	}

	return false
}
