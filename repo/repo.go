package repo

import (
	"errors"
)

var (
	ErrUserNotFound error = errors.New("user not found")
)

type memoryDB struct {
	users map[int64]User
}

// AddUser implements db.
func (m *memoryDB) AddUser(user User) error {
	m.users[user.TelegramID] = User{
		ID:         int64(len(m.users) + 1),
		TelegramID: user.TelegramID,
		FullName:   user.FullName,
	}
	return nil
}

// Decryption implements db.
func (m *memoryDB) Decryption(telegramID int64, text string, decryptionType Type) (string, error) {

	user, ok := m.users[telegramID]

	if !ok {
		return "", ErrUserNotFound
	}
	user.Decryption = true

	return text, nil
}

// Encryption implements db.
func (m *memoryDB) Encryption(telegramID int64, text string, encryptionType Type) (string, error) {
	panic("unimplemented")
}

// GetUserByTelegramID implements db.
func (m *memoryDB) GetUserByTelegramID(telegramID int64) (User, error) {
	user, ok := m.users[telegramID]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

// UpdateUser implements db.
func (m *memoryDB) UpdateUser(user User) error {
	m.users[user.TelegramID] = user
	return nil
}

type Db interface {
	AddUser(user User) error
	GetUserByTelegramID(telegramID int64) (User, error)
	UpdateUser(user User) error
	Encryption(telegramID int64, text string, encryptionType Type) (string, error)
	Decryption(telegramID int64, text string, decryptionType Type) (string, error)
}

func NewDB() Db {
	return &memoryDB{
		users: make(map[int64]User),
	}
}
