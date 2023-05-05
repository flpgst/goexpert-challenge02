package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ViaCep struct {
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

type ApiCEP struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

func main() {
	if len(os.Args) == 1 {
		println("CEP não informado")
		return
	}

	c1 := make(chan ViaCep)
	c2 := make(chan ApiCEP)
	cep := os.Args[1]

	go func() {
		res := getCEP("https://viacep.com.br/ws/" + cep + "/json/")
		var viacep ViaCep
		err := json.Unmarshal(res, &viacep)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %x\n", err)
		}
		c1 <- viacep

	}()
	go func() {
		res := getCEP("https://cdn.apicep.com/file/apicep/" + cep + ".json")
		var apicep ApiCEP
		err := json.Unmarshal(res, &apicep)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %x\n", err)
		}
		// Avoiding Bad Request - Blocked by flood cdn
		if apicep.Status == 200 {
			c2 <- apicep
		}
	}()

	select {
	case cep := <-c1:
		fmt.Println("ViaCEP:", cep)
	case cep := <-c2:
		fmt.Println("ApiCEP: ", cep)
	case <-time.After(time.Second):
		println("Timeout!")
	}

}

func getCEP(url string) []byte {
	req, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %x\n", err)
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler resposta: %x\n", err)
	}
	return res

}
