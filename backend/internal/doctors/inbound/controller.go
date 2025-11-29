package inbound

import (
	"context"

	"github.com/OPetricevic/physio-tracker/backend/internal/doctors"
)

type Controller struct {
	svc doctors.Service
}

func NewController(svc doctors.Service) *Controller {
	return &Controller{svc: svc}
}

func (c *Controller) Register(ctx context.Context, input doctors.RegisterInput) (doctors.Doctor, error) {
	return c.svc.Register(ctx, input)
}

func (c *Controller) Authenticate(ctx context.Context, email, password string) (doctors.Doctor, error) {
	return c.svc.Authenticate(ctx, email, password)
}

func (c *Controller) GetByUUID(ctx context.Context, uuid string) (doctors.Doctor, error) {
	return c.svc.GetByUUID(ctx, uuid)
}
