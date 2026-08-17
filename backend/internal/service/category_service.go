package service

import (
	"log/slog"

	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/repository"
	"github.com/flowershop/backend/internal/util"
)

// CategoryService 分类服务。
type CategoryService struct {
	repo     *repository.CategoryRepository
	auditSvc *AuditService
	logger   *slog.Logger
}

func NewCategoryService(repo *repository.CategoryRepository, auditSvc *AuditService, logger *slog.Logger) *CategoryService {
	return &CategoryService{repo: repo, auditSvc: auditSvc, logger: logger}
}

func (s *CategoryService) Create(req dto.CategoryCreateRequest) (*model.Category, error) {
	level := 1
	if req.ParentID > 0 {
		parent, err := s.repo.FindByID(req.ParentID)
		if err != nil {
			return nil, util.NewAppError(constants.CodeCategoryNotFound, "parent category not found, parent_id="+util.UintString(req.ParentID))
		}
		level = parent.Level + 1
	}
	c := &model.Category{ParentID: req.ParentID, Name: req.Name, Level: level, Sort: req.Sort}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogCategoryCreated, "id", c.ID, "name", c.Name, "level", c.Level, "parent_id", c.ParentID)
	return c, nil
}

func (s *CategoryService) Update(id uint, req dto.CategoryUpdateRequest) (*model.Category, error) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.NewAppError(constants.CodeCategoryNotFound, "category not found, id="+util.UintString(id))
	}
	if req.Name != "" {
		c.Name = req.Name
	}
	if req.ParentID != nil {
		c.ParentID = *req.ParentID
	}
	if req.Sort != nil {
		c.Sort = *req.Sort
	}
	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CategoryService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return util.NewAppError(constants.CodeCategoryNotFound, "delete category failed, id="+util.UintString(id))
	}
	return nil
}

// ListTree 返回分类树（复用：handler 列表与商品查询）。
func (s *CategoryService) ListTree() ([]dto.CategoryVO, error) {
	all, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}
	children := map[uint][]dto.CategoryVO{}
	for _, c := range all {
		children[c.ParentID] = append(children[c.ParentID], dto.CategoryVO{ID: c.ID, ParentID: c.ParentID, Name: c.Name, Level: c.Level, Sort: c.Sort})
	}
	var build func(pid uint) []dto.CategoryVO
	build = func(pid uint) []dto.CategoryVO {
		out := make([]dto.CategoryVO, 0, len(children[pid]))
		for _, c := range children[pid] {
			c.Children = build(c.ID)
			out = append(out, c)
		}
		return out
	}
	return build(0), nil
}
