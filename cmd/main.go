package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Endereco struct {
	Cep                 string `json:"cep"`
	Logradouro          string `json:"logradouro,omitempty"`
	Bairro              string `json:"bairro,omitempty"`
	Cidade              string `json:"localidade,omitempty"`
	Uf                  string `json:"uf,omitempty"`
	LogradouroBrasilAPI string `json:"street,omitempty"`
	BairroBrasilAPI     string `json:"neighborhood,omitempty"`
	CidadeBrasilAPI     string `json:"city,omitempty"`
	UfBrasilAPI         string `json:"state,omitempty"`
}

const BrasilAPIURL = "https://brasilapi.com.br/api/cep/v1/%s"
const ViaCepURL = "https://viacep.com.br/ws/%s/json/"

func main() {

	c1 := make(chan Endereco)
	c2 := make(chan Endereco)

	go getCep("01153000", BrasilAPIURL, c1)
	go getCep("01153000", ViaCepURL, c2)

	select {
	case endereco := <-c1:
		fmt.Printf(
			"Brasil API: Cep: %s, Logradouro: %s, Bairro: %s, Cidade: %s, UF: %s\n",
			endereco.Cep,
			endereco.LogradouroBrasilAPI,
			endereco.BairroBrasilAPI,
			endereco.CidadeBrasilAPI,
			endereco.UfBrasilAPI,
		)
	case endereco := <-c2:
		fmt.Printf(
			"Via CEP: Cep: %s, Logradouro: %s, Bairro: %s, Cidade: %s, UF: %s\n",
			endereco.Cep,
			endereco.Logradouro,
			endereco.Bairro,
			endereco.Cidade,
			endereco.Uf,
		)
	case <-time.After(time.Second * 1):
		println("timeout")
	}
}

func getCep(cep, url string, c chan Endereco) {
	client := http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf(url, cep), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("User-Agent", "Go")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	endereco := deserializeBrasilAPICep(body)
	c <- endereco

}

func deserializeBrasilAPICep(jsonEndereco []byte) Endereco {
	var endereco Endereco
	err := json.Unmarshal(jsonEndereco, &endereco)
	if err != nil {
		panic(err)
	}
	return endereco
}
