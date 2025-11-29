package patients

import "context"

// Domain models (using proto-aligned field names where practical).
type Patient struct {
	UUID      string
	DoctorUUID string
	FirstName string
	LastName  string
	Phone     *string
	CreatedAt int64
	UpdatedAt int64
}

type Anamnesis struct {
	UUID        string
	PatientUUID string
	Note        string
	CreatedAt   int64
}

// Inbound port: exposed to controllers/handlers.
type Service interface {
	CreatePatient(ctx context.Context, input CreatePatientInput) (Patient, error)
	UpdatePatient(ctx context.Context, input UpdatePatientInput) (Patient, error)
	DeletePatient(ctx context.Context, patientUUID string) error
	ListPatients(ctx context.Context, filter ListPatientsFilter) ([]Patient, error)

	AddAnamnesis(ctx context.Context, input AddAnamnesisInput) (Anamnesis, error)
	ListAnamneses(ctx context.Context, patientUUID string, page, pageSize int) ([]Anamnesis, error)
	GenerateAnamnesisPDF(ctx context.Context, patientUUID, anamnesisUUID string) ([]byte, error)
	BackupPatient(ctx context.Context, patientUUID string) ([]byte, error)
}

// Outbound ports (adapters implement these).
type PatientRepository interface {
	Insert(ctx context.Context, p Patient) error
	Update(ctx context.Context, p Patient) error
	Delete(ctx context.Context, patientUUID string) error
	List(ctx context.Context, filter ListPatientsFilter) ([]Patient, error)
	Get(ctx context.Context, patientUUID string) (Patient, error)
}

type AnamnesisRepository interface {
	Insert(ctx context.Context, a Anamnesis) error
	ListByPatient(ctx context.Context, patientUUID string, page, pageSize int) ([]Anamnesis, error)
	Get(ctx context.Context, patientUUID, anamnesisUUID string) (Anamnesis, error)
	DeleteByPatient(ctx context.Context, patientUUID string) error
}

type PDFGenerator interface {
	Generate(ctx context.Context, patient Patient, anamnesis Anamnesis) ([]byte, error)
}

type BackupStore interface {
	BackupPatient(ctx context.Context, patient Patient, anamneses []Anamnesis) ([]byte, error)
}

type CreatePatientInput struct {
	DoctorUUID string
	FirstName  string
	LastName   string
	Phone      *string
}

type UpdatePatientInput struct {
	UUID       string
	DoctorUUID string
	FirstName  string
	LastName   string
	Phone      *string
}

type ListPatientsFilter struct {
	DoctorUUID string
	Query      string
}

type AddAnamnesisInput struct {
	PatientUUID string
	Note        string
}
