package usecase

import "github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/repository"

import "github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/domain"

type CategoryUsecase interface {
	GetAll() ([]domain.Category, error)
	GetByID(id uint) (*domain.Category, error)
	Create(name string) (*domain.Category, error)
	Update(id uint, name string) (*domain.Category, error)
	Delete(id uint) error
}

type categoryUsecase struct {
	repo repository.CategoryRepository
}

func NewCategoryUsecase(repo repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{repo: repo}
}

func (uc *categoryUsecase) GetAll() ([]domain.Category, error) {
	return uc.repo.GetAll()
}

func (uc *categoryUsecase) GetByID(id uint) (*domain.Category, error) {
	return uc.repo.GetByID(id)
}

func (uc *categoryUsecase) Create(name string) (*domain.Category, error) {
	category := &domain.Category{Nama: name}
	if err := uc.repo.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (uc *categoryUsecase) Update(id uint, name string) (*domain.Category, error) {
	category, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, nil
	}

	category.Nama = name
	if err := uc.repo.Update(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (uc *categoryUsecase) Delete(id uint) error {
	return uc.repo.Delete(id)
}
