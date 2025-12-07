package auth

import (
	"context"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
)

type Repository interface {
	Create(ctx context.Context, t *pb.AuthToken) (*pb.AuthToken, error)
	Get(ctx context.Context, token string) (*pb.AuthToken, error)
	Delete(ctx context.Context, token string) error
	DeleteByDoctor(ctx context.Context, doctorUUID string) error
}
