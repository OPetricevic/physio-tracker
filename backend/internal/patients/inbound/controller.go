package inbound

import (
	"context"

	"github.com/OPetricevic/physio-tracker/backend/internal/patients"
)

// Controller mediates between HTTP/transport and the service, handling validation/auth checks.
type Controller struct {
	svc patients.Service
}

func NewController(svc patients.Service) *Controller {
	return &Controller{svc: svc}
}

func (c *Controller) CreatePatient(ctx context.Context, input patients.CreatePatientInput) (patients.Patient, error) {
	return c.svc.CreatePatient(ctx, input)
}

func (c *Controller) UpdatePatient(ctx context.Context, input patients.UpdatePatientInput) (patients.Patient, error) {
	return c.svc.UpdatePatient(ctx, input)
}

func (c *Controller) DeletePatient(ctx context.Context, uuid string) error {
	return c.svc.DeletePatient(ctx, uuid)
}

func (c *Controller) ListPatients(ctx context.Context, filter patients.ListPatientsFilter) ([]patients.Patient, error) {
	return c.svc.ListPatients(ctx, filter)
}

func (c *Controller) AddAnamnesis(ctx context.Context, input patients.AddAnamnesisInput) (patients.Anamnesis, error) {
	return c.svc.AddAnamnesis(ctx, input)
}

func (c *Controller) ListAnamneses(ctx context.Context, patientUUID string, page, pageSize int) ([]patients.Anamnesis, error) {
	return c.svc.ListAnamneses(ctx, patientUUID, page, pageSize)
}

func (c *Controller) GenerateAnamnesisPDF(ctx context.Context, patientUUID, anamnesisUUID string) ([]byte, error) {
	return c.svc.GenerateAnamnesisPDF(ctx, patientUUID, anamnesisUUID)
}

func (c *Controller) BackupPatient(ctx context.Context, patientUUID string) ([]byte, error) {
	return c.svc.BackupPatient(ctx, patientUUID)
}
