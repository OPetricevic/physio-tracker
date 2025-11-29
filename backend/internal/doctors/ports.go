package doctors

import "context"

type Doctor struct {
	UUID         string
	Email        string
	FirstName    string
	LastName     string
	PasswordHash string
	CreatedAt    int64
}

// Inbound port.
type Service interface {
	Register(ctx context.Context, input RegisterInput) (Doctor, error)
	Authenticate(ctx context.Context, email, password string) (Doctor, error)
	GetByUUID(ctx context.Context, uuid string) (Doctor, error)
}

// Outbound ports.
type Repository interface {
	Insert(ctx context.Context, d Doctor) error
	GetByUUID(ctx context.Context, uuid string) (Doctor, error)
	GetByEmail(ctx context.Context, email string) (Doctor, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hashed string) error
}

type RegisterInput struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
}
