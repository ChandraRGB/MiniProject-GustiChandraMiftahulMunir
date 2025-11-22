package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Provinsi represents a province from EMSIFA API.
type Provinsi struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Kota represents a city/regency from EMSIFA API.
type Kota struct {
	ID         string `json:"id"`
	ProvinceID string `json:"province_id"`
	Name       string `json:"name"`
}

// ProvinceCityUsecase defines methods for fetching provinces and cities.
type ProvinceCityUsecase interface {
	GetProvinces() ([]Provinsi, error)
	GetCitiesByProvince(provID string) ([]Kota, error)
}

type provinceCityUsecase struct{}

// NewProvinceCityUsecase creates a new ProvinceCityUsecase.
func NewProvinceCityUsecase() ProvinceCityUsecase {
	return &provinceCityUsecase{}
}

const baseEMSIFAURL = "https://emsifa.github.io/api-wilayah-indonesia/api"

func (uc *provinceCityUsecase) GetProvinces() ([]Provinsi, error) {
	url := fmt.Sprintf("%s/provinces.json", baseEMSIFAURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var list []Provinsi
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}

	return list, nil
}

func (uc *provinceCityUsecase) GetCitiesByProvince(provID string) ([]Kota, error) {
	url := fmt.Sprintf("%s/regencies/%s.json", baseEMSIFAURL, provID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var list []Kota
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}

	return list, nil
}
