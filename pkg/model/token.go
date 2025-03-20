package model

type Token struct {
	store TokenStore
}

type TokenStore interface {
	SaveToken(email, token string) error
	GetToken(email string) (string, error)
	DeleteToken(email string) error
	ValidateSessionToken(token string) (string, error)
}

func NewToken(store TokenStore) *Token {
	return &Token{
		store: store,
	}
}

func (t *Token) SaveToken(email, token string) error {
	return t.store.SaveToken(email, token)
}

func (t *Token) GetToken(email string) (string, error) {
	return t.store.GetToken(email)
}

func (t *Token) DeleteToken(email string) error {
	return t.store.DeleteToken(email)
}

func (t *Token) ValidateSessionToken(token string) (string, error) {
	return t.store.ValidateSessionToken(token)
}
