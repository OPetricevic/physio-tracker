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
	doctorprofilesoutboundport "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/outbound/doctorprofiles"
	outdoctors "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/outbound/doctors"
	outboundportpatients "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/outbound/patients"
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
	patientRepo outboundportpatients.Repository
	profileRepo doctorprofilesoutboundport.Repository
	doctorRepo  outdoctors.Repository
}

func NewService(
	repo out.Repository,
	pRepo outboundportpatients.Repository,
	profRepo doctorprofilesoutboundport.Repository,
	dRepo outdoctors.Repository) Service {
	return &service{
		repo:        repo,
		patientRepo: pRepo,
		profileRepo: profRepo,
		doctorRepo:  dRepo}
}

func (s *service) Create(ctx context.Context, doctorUUID string, req *pb.CreateAnamnesisRequest) (*pb.Anamnesis, error) {
	if strings.TrimSpace(req.GetPatientUuid()) == "" {
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
		Status:            strings.TrimSpace(req.GetStatus()),
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
	if req.Anamnesis != nil {
		existing.Anamnesis = strings.TrimSpace(req.Anamnesis.GetValue())
	}
	if req.Status != nil {
		existing.Status = strings.TrimSpace(req.Status.GetValue())
	}
	if req.Diagnosis != nil {
		existing.Diagnosis = strings.TrimSpace(req.Diagnosis.GetValue())
	}
	if req.Therapy != nil {
		existing.Therapy = strings.TrimSpace(req.Therapy.GetValue())
	}
	if req.OtherInfo != nil {
		existing.OtherInfo = strings.TrimSpace(req.OtherInfo.GetValue())
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
	doctor, _ := s.doctorRepo.Get(ctx, doctorUUID)           // optional

	return buildPDF(profile, doctor, patient, target, prior)
}

func buildPDF(profile *pb.DoctorProfile, doctor *pb.Doctor, patient *pb.Patient, current *pb.Anamnesis, prior []*pb.Anamnesis) ([]byte, error) {
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
	headerTop := 15.0
	logoSize := 24.0
	textX := 45.0
	if profile != nil && profile.GetLogoPath() != "" {
		if local := localFromStatic(profile.GetLogoPath()); local != "" {
			pdf.ImageOptions(local, 15, headerTop, logoSize, 0, false, gofpdf.ImageOptions{ImageType: "", ReadDpi: true}, 0, "")
		}
	}
	// printed date top-right
	pdf.SetFont("DejaVu", "", 10)
	pdf.SetXY(-70, headerTop)
	pdf.Cell(55, 5, tr("Datum ispisa: "+printedOn))

	pdf.SetFont("DejaVu", "B", 12)
	if profile != nil {
		pdf.SetXY(textX, headerTop)
		pdf.Cell(0, 6, tr(profile.GetPracticeName()))
		pdf.Ln(6)
		pdf.SetFont("DejaVu", "", 10)
		pdf.SetX(textX)
		roleDept := strings.TrimSpace(strings.TrimSpace(profile.GetRoleTitle()) + " " + strings.TrimSpace(profile.GetDepartment()))
		if roleDept != "" {
			pdf.Cell(0, 5, tr(roleDept))
			pdf.Ln(5)
			pdf.SetX(textX)
		}
		pdf.SetX(textX)
		pdf.Cell(0, 5, tr(profile.GetAddress()))
		pdf.Ln(5)
		pdf.SetFont("DejaVu", "", 10)
		pdf.SetX(textX)
		if strings.TrimSpace(profile.GetPhone()) != "" {
			pdf.Cell(0, 5, tr(profile.GetPhone()))
			pdf.Ln(5)
			pdf.SetX(textX)
		}
		if strings.TrimSpace(profile.GetEmail()) != "" {
			pdf.Cell(0, 5, tr(profile.GetEmail()))
			pdf.Ln(5)
			pdf.SetX(textX)
		}
		if strings.TrimSpace(profile.GetWebsite()) != "" {
			pdf.SetX(textX)
			pdf.Cell(0, 5, tr(strings.TrimSpace(profile.GetWebsite())))
			pdf.Ln(5)
		}
		pdf.Ln(6)
	}

	//centered title
	headerBottom := pdf.GetY()
	pdf.SetX(15)
	pdf.Line(15, headerBottom, 195, headerBottom)
	pdf.Ln(6)
	pdf.SetFont("DejaVu", "B", 12)
	pdf.CellFormat(0, 6, tr("MIŠLJENJE FIZIOTERAPEUTA"), "", 1, "C", false, 0, "")
	pdf.Ln(6)

	pdf.SetFont("DejaVu", "B", 10)
	pdf.Cell(40, 5, tr("Pacijent:"))
	pdf.SetFont("DejaVu", "", 10)
	pdf.Cell(0, 5, tr(strings.TrimSpace(patient.GetFirstName()+" "+patient.GetLastName())))
	pdf.Ln(5)
	if dob := patient.GetDateOfBirth(); dob != nil && strings.TrimSpace(dob.GetValue()) != "" {
		pdf.SetFont("DejaVu", "B", 10)
		pdf.Cell(40, 5, tr("Datum rođenja:"))
		pdf.SetFont("DejaVu", "", 10)
		pdf.Cell(0, 5, tr(formatPlainDate(dob.GetValue())))
		pdf.Ln(5)
	}
	if patient.GetPhone() != nil {
		pdf.SetFont("DejaVu", "B", 10)
		pdf.Cell(40, 5, tr("Telefon:"))
		pdf.SetFont("DejaVu", "", 10)
		pdf.Cell(0, 5, tr(patient.GetPhone().GetValue()))
		pdf.Ln(5)
	}
	if patient.GetAddress() != nil {
		pdf.SetFont("DejaVu", "B", 10)
		pdf.Cell(40, 5, tr("Adresa:"))
		pdf.SetFont("DejaVu", "", 10)
		pdf.Cell(0, 5, tr(patient.GetAddress().GetValue()))
		pdf.Ln(8)
	}

	// set line to separate header from content (match header thickness)
	pdf.SetLineWidth(0.2)
	pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
	pdf.Ln(6)

	writeField := func(label, value string) {
		if strings.TrimSpace(value) == "" {
			return
		}
		pdf.SetFont("DejaVu", "B", 10)
		pdf.MultiCell(0, 5, tr(label), "", "", false)
		pdf.SetFont("DejaVu", "", 10)
		pdf.MultiCell(0, 5, tr(value), "", "", false)
		pdf.Ln(2)
	}

	visits := make([]*pb.Anamnesis, 0, len(prior)+1)
	visits = append(visits, current)
	visits = append(visits, prior...)
	sort.SliceStable(visits, func(i, j int) bool {
		ti := time.Time{}
		tj := time.Time{}
		if visits[i].GetCreatedAt() != nil {
			ti = visits[i].GetCreatedAt().AsTime()
		}
		if visits[j].GetCreatedAt() != nil {
			tj = visits[j].GetCreatedAt().AsTime()
		}
		return ti.Before(tj)
	})

	for i, v := range visits {
		pdf.SetFont("DejaVu", "B", 11)
		label := fmt.Sprintf("%d. posjet - %s", i+1, formatDate(v.GetCreatedAt()))
		pdf.MultiCell(0, 5, tr(label), "", "", false)
		pdf.SetFont("DejaVu", "", 10)
		writeField("Anamneza", v.GetAnamnesis())
		writeField("Status", v.GetStatus())
		writeField("Dijagnoza", v.GetDiagnosis())
		writeField("Terapija", v.GetTherapy())
		writeField("Ostalo", v.GetOtherInfo())
		pdf.Ln(2)
	}

	// Footer signature (doctor)
	if doctor != nil {
		name := strings.TrimSpace(doctor.GetFirstName() + " " + doctor.GetLastName())
		if name != "" {
			pdf.SetY(-35)
			pdf.SetX(15)
			pdf.SetFont("DejaVu", "B", 9)
			pdf.Cell(0, 5, tr("Fizioterapeut:"))
			pdf.Ln(5)
			pdf.SetX(15)
			pdf.SetFont("DejaVu", "", 9)
			pdf.Cell(0, 5, tr("bacc.physioth "+name))
		}
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
