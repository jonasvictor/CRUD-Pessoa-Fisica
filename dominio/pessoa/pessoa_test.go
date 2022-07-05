package pessoa_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

// Define
type Pessoa struct {
	Results []struct {
		ID             int    `json:"id"`
		NomeCompleto   string `json:"nome"`
		Endereco       string `json:"endereco"`
		DataNascimento string `json:"data-de-nascimento"`
		Cpf            string `json:"cpf"`
		Telefone       int    `json:"telefone"`
	} `json:"results"`
	Status string `json:"status"`
}

// Test do Post
func TestCreate(t *testing.T) {
	resp, err := http.Post(
		"http://localhost:8080/pessoa/",
		"application/json",
		bytes.NewReader([]byte(`{"id": 1, "nome": "Jonas",
								 "endereco": "Rua Principal",
								 "data-de-nascimento": "09/03/1997",
								 "cpf": "999.999.999-99",
								 "telefone": 83987654321}`)))

	if err != nil {
		t.Errorf("Error ao fazer requisição: %v", err.Error())
	}

	defer resp.Body.Close()

	if err != nil {
		t.Errorf("Erro ao preencher os campos: %v", err.Error())
	}

}

// Faz a lista de todas as pessoas
func TestGetPessoas(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/pessoa/")
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	log.Println(string(body))
	pess := Pessoa{}
	err = json.Unmarshal([]byte(string(body)), &pess)
	if err != nil {
		log.Println(err)
		t.Error(err)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Sem sucesso: %d", resp.StatusCode)
	}
}

// Busca a pessoa cadastrada pelo ID
func TestGetByID(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/pessoa/1")
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		t.Error(err)
	}

	log.Println(string(body))
	pess := Pessoa{}
	err = json.Unmarshal([]byte(string(body)), &pess)

	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		fmt.Printf("Sem sucesso: %d", resp.StatusCode)
	}

}

// Atualiza os dados da pessoa cadastrada pelo ID
func TestUpdate(t *testing.T) {
	req, err := http.NewRequest(
		"PUT",
		"http://localhost:8080/pessoa/1",
		bytes.NewBuffer([]byte(`"id": 1, "nome": "Jonas",
								"endereco": "Rua Principal",
								"data-de-nascimento": "09/03/1997",
								"cpf": "999.999.999-99",
								"telefone": 83987654321}`)))

	if err != nil {
		t.Error(err)
	}

	defer req.Body.Close()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Sem sucesso pessoa não cadastrada: %d", resp.StatusCode)
	}
}

// Deleta a pessoa cadastrada através do ID
func TestDelete(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://localhost:8080/pessoa/2", nil)
	if err != nil {
		t.Error("*********************************************")
		t.Error(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error("---------------------------------------------")
		t.Error(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Sem sucesso pessoa não cadastrada: %d", resp.StatusCode)
	}
}
