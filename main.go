package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type CepResponseDTO struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	SearchCEP("01153000")
}

func SearchCEP(cep string) {
	var cepData *CepResponseDTO
	wg := sync.WaitGroup{}
	wg.Add(1)

	ctx := context.Background()

	go func() {
		cepData, _ = GetCepDataFromBrasilAPI(cep, ctx)
		fmt.Printf("%+v\n", cepData)
		wg.Done()
	}()

	go func() {
		cepData, _ = GetCepDataFromViaCepAPI(cep, ctx)
		fmt.Printf("%+v\n", cepData)
		wg.Done()
	}()

	wg.Wait()
}

func GetCepDataFromBrasilAPI(cep string, ctx context.Context) (*CepResponseDTO, error) {
	log.Println("getting CEP information from BrasilAPI")

	getContext, cancel := context.WithTimeout(ctx, 1*time.Second)

	defer cancel()

	req, err := http.NewRequestWithContext(getContext, "GET", "https://brasilapi.com.br/api/cep/v1/"+cep, nil)

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("getting CEP information, error: ", err.Error())

		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var cepData CepResponseDTO

	err = json.Unmarshal(body, &cepData)

	if err != nil {
		return nil, err
	}

	cepData.Service = "BrasilAPI"

	log.Println("CEP information retrieved successfully")
	return &cepData, nil
}

func GetCepDataFromViaCepAPI(cep string, ctx context.Context) (*CepResponseDTO, error) {
	log.Println("getting CEP information from ViaCepAPI")

	getContext, cancel := context.WithTimeout(ctx, 1*time.Second)

	defer cancel()

	req, err := http.NewRequestWithContext(getContext, "GET", "https://viacep.com.br/ws/"+cep+"/json", nil)

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("getting CEP information, error: ", err.Error())

		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var cepData ViaCepResponse

	err = json.Unmarshal(body, &cepData)

	if err != nil {
		return nil, err
	}

	log.Println("CEP information retrieved successfully")

	return &CepResponseDTO{
		Cep:          cepData.Cep,
		State:        cepData.Uf,
		City:         cepData.Localidade,
		Neighborhood: cepData.Bairro,
		Street:       cepData.Logradouro,
		Service:      "ViaCep",
	}, nil
}
