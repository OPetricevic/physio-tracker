package anamneses

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"os"
	"path/filepath"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	out "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/outbound/anamneses"
	re "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/repoerrors"
	se "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/serviceerrors"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jung-kurt/gofpdf"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, doctorUUID string, req *pb.CreateAnamnesisRequest) (*pb.Anamnesis, error)
	Update(ctx context.Context, doctorUUID string, req *pb.UpdateAnamnesisRequest) (*pb.Anamnesis, error)
	List(ctx context.Context, doctorUUID, patientUUID, query string, pageSize, currentPage int) ([]*pb.Anamnesis, error)
	Delete(ctx context.Context, doctorUUID, uuid string) error
	Get(ctx context.Context, doctorUUID, uuid string) (*pb.Anamnesis, error)
	GeneratePDF(ctx context.Context, doctorUUID, patientUUID, anamnesisUUID string, include []string, onlyCurrent bool) ([]byte, error)
}

type service struct {
	repo        out.Repository
	patientRepo patientsRepo
	profileRepo profileRepo
}

type patientsRepo interface {
	Get(ctx context.Context, uuid string) (*pb.Patient, error)
}

type profileRepo interface {
	GetByDoctor(ctx context.Context, doctorUUID string) (*pb.DoctorProfile, error)
}

func NewService(repo out.Repository, pRepo patientsRepo, profRepo profileRepo) Service {
	return &service{repo: repo, patientRepo: pRepo, profileRepo: profRepo}
}

func (s *service) Create(ctx context.Context, doctorUUID string, req *pb.CreateAnamnesisRequest) (*pb.Anamnesis, error) {
	if strings.TrimSpace(req.GetPatientUuid()) == "" {
		return nil, fmt.Errorf("create anamnesis: %w", se.ErrInvalidRequest)
	}
	if strings.TrimSpace(req.GetAnamnesis()) == "" || strings.TrimSpace(req.GetDiagnosis()) == "" || strings.TrimSpace(req.GetTherapy()) == "" {
		return nil, fmt.Errorf("create anamnesis: %w", se.ErrInvalidRequest)
	}
	include := req.IncludeVisitUuids
	if include == nil {
		include = []string{}
	}
	now := time.Now().UTC()
	a := &pb.Anamnesis{
		Uuid:              uuid.NewString(),
		PatientUuid:       strings.TrimSpace(req.GetPatientUuid()),
		Anamnesis:         strings.TrimSpace(req.GetAnamnesis()),
		Diagnosis:         strings.TrimSpace(req.GetDiagnosis()),
		Therapy:           strings.TrimSpace(req.GetTherapy()),
		OtherInfo:         strings.TrimSpace(req.GetOtherInfo()),
		IncludeVisitUuids: include,
		CreatedAt:         timestamppb.New(now),
		UpdatedAt:         nil,
	}
	created, err := s.repo.Create(ctx, a)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("create anamnesis: %w", se.ErrConflict)
		}
		return nil, fmt.Errorf("create anamnesis: %w", err)
	}
	return created, nil
}

func (s *service) Update(ctx context.Context, doctorUUID string, req *pb.UpdateAnamnesisRequest) (*pb.Anamnesis, error) {
	if strings.TrimSpace(req.GetUuid()) == "" || strings.TrimSpace(req.GetPatientUuid()) == "" {
		return nil, fmt.Errorf("update anamnesis: %w", se.ErrInvalidRequest)
	}
	existing, err := s.repo.Get(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, re.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("update anamnesis: %w", se.ErrNotFound)
		}
		return nil, fmt.Errorf("load anamnesis for update: %w", err)
	}
	if strings.TrimSpace(existing.GetPatientUuid()) != strings.TrimSpace(req.GetPatientUuid()) {
		return nil, fmt.Errorf("update anamnesis: %w", se.ErrInvalidRequest)
	}
	if strings.TrimSpace(req.GetAnamnesis()) != "" {
		existing.Anamnesis = strings.TrimSpace(req.GetAnamnesis())
	}
	if strings.TrimSpace(req.GetDiagnosis()) != "" {
		existing.Diagnosis = strings.TrimSpace(req.GetDiagnosis())
	}
	if strings.TrimSpace(req.GetTherapy()) != "" {
		existing.Therapy = strings.TrimSpace(req.GetTherapy())
	}
	if strings.TrimSpace(req.GetOtherInfo()) != "" {
		existing.OtherInfo = strings.TrimSpace(req.GetOtherInfo())
	}
	if req.IncludeVisitUuids != nil {
		existing.IncludeVisitUuids = req.IncludeVisitUuids
	}
	existing.UpdatedAt = timestamppb.New(time.Now().UTC())

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		if errors.Is(err, re.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("update anamnesis: %w", se.ErrNotFound)
		}
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("update anamnesis: %w", se.ErrConflict)
		}
		return nil, fmt.Errorf("update anamnesis: %w", err)
	}
	return updated, nil
}

func (s *service) List(ctx context.Context, doctorUUID, patientUUID, query string, pageSize, currentPage int) ([]*pb.Anamnesis, error) {
	if strings.TrimSpace(patientUUID) == "" {
		return nil, fmt.Errorf("list anamneses: %w", se.ErrInvalidRequest)
	}
	if pageSize <= 0 {
		pageSize = 5
	}
	if currentPage <= 0 {
		currentPage = 1
	}
	offset := (currentPage - 1) * pageSize
	list, err := s.repo.List(ctx, patientUUID, doctorUUID, query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("list anamneses: %w", err)
	}
	return list, nil
}

func (s *service) Delete(ctx context.Context, doctorUUID, uuid string) error {
	if strings.TrimSpace(uuid) == "" {
		return fmt.Errorf("delete anamnesis: %w", se.ErrInvalidRequest)
	}
	if err := s.repo.Delete(ctx, uuid); err != nil {
		if errors.Is(err, re.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("delete anamnesis: %w", se.ErrNotFound)
		}
		return fmt.Errorf("delete anamnesis: %w", err)
	}
	return nil
}

func (s *service) Get(ctx context.Context, doctorUUID, uuid string) (*pb.Anamnesis, error) {
	if strings.TrimSpace(uuid) == "" {
		return nil, fmt.Errorf("get anamnesis: %w", se.ErrInvalidRequest)
	}
	anm, err := s.repo.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, re.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get anamnesis: %w", se.ErrNotFound)
		}
		return nil, fmt.Errorf("get anamnesis: %w", err)
	}
	return anm, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func (s *service) GeneratePDF(ctx context.Context, doctorUUID, patientUUID, anamnesisUUID string, include []string, onlyCurrent bool) ([]byte, error) {
	if strings.TrimSpace(doctorUUID) == "" || strings.TrimSpace(patientUUID) == "" || strings.TrimSpace(anamnesisUUID) == "" {
		return nil, fmt.Errorf("generate pdf: %w", se.ErrInvalidRequest)
	}

	patient, err := s.patientRepo.Get(ctx, patientUUID)
	if err != nil {
		if errors.Is(err, re.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("generate pdf: %w", se.ErrNotFound)
		}
		return nil, fmt.Errorf("generate pdf: load patient: %w", err)
	}
	if strings.TrimSpace(patient.GetDoctorUuid()) != strings.TrimSpace(doctorUUID) {
		return nil, fmt.Errorf("generate pdf: %w", se.ErrInvalidRequest)
	}

	target, err := s.repo.Get(ctx, anamnesisUUID)
	if err != nil {
		if errors.Is(err, re.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("generate pdf: %w", se.ErrNotFound)
		}
		return nil, fmt.Errorf("generate pdf: load anamnesis: %w", err)
	}
	if strings.TrimSpace(target.GetPatientUuid()) != strings.TrimSpace(patientUUID) {
		return nil, fmt.Errorf("generate pdf: %w", se.ErrInvalidRequest)
	}

	includeList := include
	if includeList == nil {
		includeList = []string{}
	}
	if onlyCurrent {
		includeList = []string{}
	} else if len(includeList) == 0 && len(target.IncludeVisitUuids) > 0 {
		includeList = target.IncludeVisitUuids
	}
	var prior []*pb.Anamnesis
	if len(includeList) > 0 {
		list, err := s.repo.ListByUUIDs(ctx, includeList)
		if err != nil {
			return nil, fmt.Errorf("generate pdf: load included visits: %w", err)
		}
		for _, v := range list {
			if strings.TrimSpace(v.GetPatientUuid()) == strings.TrimSpace(patientUUID) {
				prior = append(prior, v)
			}
		}
		sort.Slice(prior, func(i, j int) bool {
			return prior[i].GetCreatedAt().AsTime().After(prior[j].GetCreatedAt().AsTime())
		})
	}

	profile, _ := s.profileRepo.GetByDoctor(ctx, doctorUUID) // optional

	return buildPDF(profile, patient, target, prior)
}

func buildPDF(profile *pb.DoctorProfile, patient *pb.Patient, current *pb.Anamnesis, prior []*pb.Anamnesis) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	// UTF-8 fonts to render regional characters correctly.
	fontDir := filepath.Join("assets", "fonts")
	pdf.SetFontLocation(fontDir)
	pdf.AddUTF8Font("DejaVu", "", "DejaVuSans.ttf")
	pdf.AddUTF8Font("DejaVu", "B", "DejaVuSans-Bold.ttf")
	if pdf.Err() {
		return nil, fmt.Errorf("generate pdf: load fonts from %s: %v", fontDir, pdf.Error())
	}
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()
	printedOn := time.Now().Format("02.01.2006.")
	tr := func(s string) string { return s }

	// Header
	if profile != nil && profile.GetLogoPath() != "" {
		if local := localFromStatic(profile.GetLogoPath()); local != "" {
			pdf.ImageOptions(local, 15, 15, 30, 0, false, gofpdf.ImageOptions{ImageType: "", ReadDpi: true}, 0, "")
		}
	}
	// printed date top-right
	pdf.SetFont("DejaVu", "", 10)
	pdf.SetXY(-70, 15)
	pdf.Cell(55, 5, tr("Datum ispisa: "+printedOn))

	pdf.SetFont("DejaVu", "B", 14)
	if profile != nil {
		pdf.SetXY(50, 20)
		pdf.Cell(0, 6, tr(profile.GetPracticeName()))
		pdf.Ln(6)
		pdf.SetFont("DejaVu", "", 10)
		pdf.SetX(50)
		pdf.Cell(0, 5, tr(strings.TrimSpace(profile.GetRoleTitle()+" "+profile.GetDepartment())))
		pdf.Ln(5)
		pdf.SetX(50)
		pdf.Cell(0, 5, tr(profile.GetAddress()))
		pdf.Ln(5)
		pdf.SetX(50)
		pdf.Cell(0, 5, tr(strings.TrimSpace(profile.GetPhone()+" "+profile.GetEmail())))
		pdf.Ln(5)
		if strings.TrimSpace(profile.GetWebsite()) != "" {
			pdf.SetX(50)
			pdf.Cell(0, 5, tr(strings.TrimSpace(profile.GetWebsite())))
			pdf.Ln(5)
		}
		pdf.Ln(4)
	}

	pdf.SetFont("DejaVu", "B", 12)
	pdf.Cell(0, 6, tr("Nalaz / Anamneza"))
	pdf.Ln(10)

	pdf.SetFont("DejaVu", "", 10)
	pdf.Cell(40, 5, tr("Pacijent:"))
	pdf.Cell(0, 5, tr(strings.TrimSpace(patient.GetFirstName()+" "+patient.GetLastName())))
	pdf.Ln(5)
	if dob := patient.GetDateOfBirth(); dob != nil && strings.TrimSpace(dob.GetValue()) != "" {
		pdf.Cell(40, 5, tr("Datum rođenja:"))
		pdf.Cell(0, 5, tr(formatPlainDate(dob.GetValue())))
		pdf.Ln(5)
	}
	if patient.GetPhone() != nil {
		pdf.Cell(40, 5, tr("Telefon:"))
		pdf.Cell(0, 5, tr(patient.GetPhone().GetValue()))
		pdf.Ln(5)
	}
	if patient.GetAddress() != nil {
		pdf.Cell(40, 5, tr("Adresa:"))
		pdf.Cell(0, 5, tr(patient.GetAddress().GetValue()))
		pdf.Ln(8)
	}

	pdf.SetFont("DejaVu", "B", 11)
	pdf.Cell(0, 6, tr(fmt.Sprintf("Posjet: %s", formatDate(current.GetCreatedAt()))))
	pdf.Ln(7)
	pdf.SetFont("DejaVu", "", 10)
	pdf.MultiCell(0, 5, tr("Dijagnoza: "+current.GetDiagnosis()), "", "", false)
	pdf.Ln(2)
	pdf.MultiCell(0, 5, tr("Terapija: "+current.GetTherapy()), "", "", false)
	pdf.Ln(2)
	if strings.TrimSpace(current.GetAnamnesis()) != "" {
		pdf.MultiCell(0, 5, tr("Anamneza: "+current.GetAnamnesis()), "", "", false)
		pdf.Ln(2)
	}
	if strings.TrimSpace(current.GetOtherInfo()) != "" {
		pdf.MultiCell(0, 5, tr("Ostalo: "+current.GetOtherInfo()), "", "", false)
		pdf.Ln(2)
	}
	pdf.Ln(5)

	if len(prior) > 0 {
		pdf.SetFont("DejaVu", "B", 11)
		pdf.Cell(0, 6, tr("Prethodni posjeti"))
		pdf.Ln(7)
		pdf.SetFont("DejaVu", "", 10)
		for _, v := range prior {
			pdf.SetFont("DejaVu", "B", 10)
			pdf.MultiCell(0, 5, tr(fmt.Sprintf("Datum: %s", formatDate(v.GetCreatedAt()))), "", "", false)
			pdf.SetFont("DejaVu", "", 10)
			pdf.MultiCell(0, 5, tr("Dijagnoza: "+v.GetDiagnosis()), "", "", false)
			if strings.TrimSpace(v.GetTherapy()) != "" {
				pdf.MultiCell(0, 5, tr("Terapija: "+v.GetTherapy()), "", "", false)
			}
			if strings.TrimSpace(v.GetAnamnesis()) != "" {
				pdf.MultiCell(0, 5, tr("Anamneza: "+v.GetAnamnesis()), "", "", false)
			}
			if strings.TrimSpace(v.GetOtherInfo()) != "" {
				pdf.MultiCell(0, 5, tr("Ostalo: "+v.GetOtherInfo()), "", "", false)
			}
			pdf.Ln(3)
		}
	}

	if profile != nil && strings.TrimSpace(profile.GetFooterNote()) != "" {
		pdf.SetY(-30)
		pdf.SetFont("DejaVu", "", 9)
		pdf.MultiCell(0, 5, tr(profile.GetFooterNote()), "", "C", false)
	}

	if pdf.Err() {
		return nil, fmt.Errorf("generate pdf: prepare content: %v", pdf.Error())
	}

	var buf strings.Builder
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("generate pdf: render: %w", err)
	}
	return []byte(buf.String()), nil
}

func formatDate(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	t := ts.AsTime()
	return t.Format("02.01.2006.")
}

func formatPlainDate(val string) string {
	s := strings.TrimSpace(val)
	if s == "" {
		return ""
	}
	layouts := []string{"2006-01-02", time.RFC3339}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t.Format("02.01.2006.")
		}
	}
	return s
}

func localFromStatic(path string) string {
	const prefix = "/static/"
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	rel := strings.TrimPrefix(path, prefix)
	local := filepath.Join("uploads", rel)
	if _, err := os.Stat(local); err == nil {
		return local
	}
	return ""
}
