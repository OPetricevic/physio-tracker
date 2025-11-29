package patients

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type service struct {
	patientRepo   PatientRepository
	anamnesisRepo AnamnesisRepository
	pdfGen        PDFGenerator
	backupStore   BackupStore
}

func NewService(patientRepo PatientRepository, anamnesisRepo AnamnesisRepository, pdfGen PDFGenerator, backupStore BackupStore) Service {
	return &service{
		patientRepo:   patientRepo,
		anamnesisRepo: anamnesisRepo,
		pdfGen:        pdfGen,
		backupStore:   backupStore,
	}
}

func (s *service) CreatePatient(ctx context.Context, input CreatePatientInput) (Patient, error) {
	if strings.TrimSpace(input.DoctorUUID) == "" {
		return Patient{}, errors.New("doctor uuid je obavezan")
	}
	if _, err := uuid.Parse(input.DoctorUUID); err != nil {
		return Patient{}, errors.New("neispravan doctor uuid")
	}
	if strings.TrimSpace(input.FirstName) == "" || strings.TrimSpace(input.LastName) == "" {
		return Patient{}, errors.New("ime i prezime su obavezni")
	}
	now := time.Now().UTC().Unix()
	p := Patient{
		UUID:       uuid.NewString(),
		DoctorUUID: input.DoctorUUID,
		FirstName:  strings.TrimSpace(input.FirstName),
		LastName:   strings.TrimSpace(input.LastName),
		Phone:      input.Phone,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.patientRepo.Insert(ctx, p); err != nil {
		return Patient{}, err
	}
	return p, nil
}

func (s *service) UpdatePatient(ctx context.Context, input UpdatePatientInput) (Patient, error) {
	if _, err := uuid.Parse(strings.TrimSpace(input.UUID)); err != nil {
		return Patient{}, errors.New("neispravan uuid")
	}
	if strings.TrimSpace(input.DoctorUUID) == "" {
		return Patient{}, errors.New("doctor uuid je obavezan")
	}
	existing, err := s.patientRepo.Get(ctx, input.UUID)
	if err != nil {
		return Patient{}, err
	}
	// Optional: enforce doctor ownership match here if needed.
	existing.FirstName = strings.TrimSpace(input.FirstName)
	existing.LastName = strings.TrimSpace(input.LastName)
	existing.Phone = input.Phone
	existing.UpdatedAt = time.Now().UTC().Unix()
	if err := s.patientRepo.Update(ctx, existing); err != nil {
		return Patient{}, err
	}
	return existing, nil
}

func (s *service) DeletePatient(ctx context.Context, patientUUID string) error {
	if _, err := uuid.Parse(strings.TrimSpace(patientUUID)); err != nil {
		return errors.New("neispravan uuid")
	}
	if err := s.anamnesisRepo.DeleteByPatient(ctx, patientUUID); err != nil {
		return err
	}
	return s.patientRepo.Delete(ctx, patientUUID)
}

func (s *service) ListPatients(ctx context.Context, filter ListPatientsFilter) ([]Patient, error) {
	if filter.DoctorUUID != "" {
		if _, err := uuid.Parse(strings.TrimSpace(filter.DoctorUUID)); err != nil {
			return nil, errors.New("neispravan doctor uuid")
		}
	}
	return s.patientRepo.List(ctx, filter)
}

func (s *service) AddAnamnesis(ctx context.Context, input AddAnamnesisInput) (Anamnesis, error) {
	if _, err := uuid.Parse(strings.TrimSpace(input.PatientUUID)); err != nil {
		return Anamnesis{}, errors.New("neispravan uuid pacijenta")
	}
	if strings.TrimSpace(input.Note) == "" {
		return Anamnesis{}, errors.New("bilješka je obavezna")
	}
	now := time.Now().UTC().Unix()
	a := Anamnesis{
		UUID:        uuid.NewString(),
		PatientUUID: input.PatientUUID,
		Note:        input.Note,
		CreatedAt:   now,
	}
	if err := s.anamnesisRepo.Insert(ctx, a); err != nil {
		return Anamnesis{}, err
	}
	return a, nil
}

func (s *service) ListAnamneses(ctx context.Context, patientUUID string, page, pageSize int) ([]Anamnesis, error) {
	if _, err := uuid.Parse(strings.TrimSpace(patientUUID)); err != nil {
		return nil, errors.New("neispravan uuid pacijenta")
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 5
	}
	return s.anamnesisRepo.ListByPatient(ctx, patientUUID, page, pageSize)
}

func (s *service) GenerateAnamnesisPDF(ctx context.Context, patientUUID, anamnesisUUID string) ([]byte, error) {
	if s.pdfGen == nil {
		return nil, errors.New("pdf generator nije konfiguriran")
	}
	patient, err := s.patientRepo.Get(ctx, patientUUID)
	if err != nil {
		return nil, err
	}
	a, err := s.anamnesisRepo.Get(ctx, patientUUID, anamnesisUUID)
	if err != nil {
		return nil, err
	}
	return s.pdfGen.Generate(ctx, patient, a)
}

func (s *service) BackupPatient(ctx context.Context, patientUUID string) ([]byte, error) {
	if s.backupStore == nil {
		return nil, errors.New("backup nije konfiguriran")
	}
	patient, err := s.patientRepo.Get(ctx, patientUUID)
	if err != nil {
		return nil, err
	}
	notes, err := s.anamnesisRepo.ListByPatient(ctx, patientUUID, 1, 10_000)
	if err != nil {
		return nil, err
	}
	return s.backupStore.BackupPatient(ctx, patient, notes)
}
