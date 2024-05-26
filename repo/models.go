package repo

type Type string

const (
	AES Type = "AES"
	DES Type = "DES"
)

type User struct {
	ID         int64
	TelegramID int64
	FullName   string
	Type       Type
	Encryption bool
	Decryption bool
}
