package course

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllCatalogue() ([]OffreFrontend, error) {
	return s.repo.GetAllCatalogue()
}

func (s *Service) GetAllApprovedCatalogue() ([]OffreFrontend, error) {
	return s.repo.GetAllApprovedCatalogue()
}

func (s *Service) GetAll() ([]Course, error) {
	return s.repo.GetAll()
}

func (s *Service) GetAllForAdmin() ([]AdminCourseView, error) {
	return s.repo.GetAllForAdmin()
}

func (s *Service) GetById(id pgtype.Int8) (*Course, error) {
	return s.repo.GetById(id)
}

func (s *Service) GetUserCourses(userId pgtype.Int8) ([]UserCourse, error) {
	return s.repo.GetUserCourses(userId)
}

func (s *Service) GetCoursesByCreator(userId pgtype.Int8) ([]Course, error) {
	return s.repo.GetCoursesByCreator(userId)
}

func (s *Service) Create(c Course) (pgtype.Int8, error) {
	return s.repo.Create(c)
}

func (s *Service) Update(id pgtype.Int8, c Course) error {
	return s.repo.Update(id, c)
}

func (s *Service) Approve(id pgtype.Int8, approvedBy pgtype.Int8) error {
	return s.repo.Approve(id, approvedBy)
}

func (s *Service) Disapprove(id pgtype.Int8) error {
	return s.repo.Disapprove(id)
}

func (s *Service) Delete(id pgtype.Int8) error {
	return s.repo.Delete(id)
}

func (s *Service) CreateDocument(doc CourseDocument) (pgtype.Int8, error) {
	return s.repo.CreateDocument(doc)
}

func (s *Service) GetDocumentsByCourseId(courseId pgtype.Int8) ([]CourseDocument, error) {
	return s.repo.GetDocumentsByCourseId(courseId)
}

func (s *Service) GetDocumentById(id pgtype.Int8) (*CourseDocument, error) {
	return s.repo.GetDocumentById(id)
}

func (s *Service) DeleteDocument(id pgtype.Int8) error {
	return s.repo.DeleteDocument(id)
}

func (s *Service) GetSessionsByCourseId(courseId pgtype.Int8) ([]CourseSession, error) {
	return s.repo.GetSessionsByCourseId(courseId)
}

func (s *Service) IsUserEnrolled(courseId pgtype.Int8, userId pgtype.Int8) (bool, error) {
	return s.repo.IsUserEnrolled(courseId, userId)
}

func (s *Service) ReplaceSessions(courseId pgtype.Int8, sessions []CourseSession) error {
	if err := s.repo.DeleteSessionsByCourseId(courseId); err != nil {
		return err
	}
	for _, sess := range sessions {
		sess.CourseId = courseId
		if _, err := s.repo.CreateSession(sess); err != nil {
			return err
		}
	}
	return nil
}
