package doctors

import (
	"context"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
)

type Repository interface {
	Create(ctx context.Context, d *pb.Doctor) (*pb.Doctor, error)
	Update(ctx context.Context, d *pb.Doctor) (*pb.Doctor, error)
	List(ctx context.Context, query string, limit, offset int) ([]*pb.Doctor, error)
	GetByIdentifier(ctx context.Context, identifier string) (*pb.Doctor, error) // email or username
	Get(ctx context.Context, uuid string) (*pb.Doctor, error)
	Delete(ctx context.Context, uuid string) error
}
