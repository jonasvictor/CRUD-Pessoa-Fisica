package pessoa_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

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

// Test do Post/ Criar cadastro da pessoa
func TestCreate(t *testing.T) {

	resp, err := http.Post(
		"http://localhost:8080/pessoa/",
		"application/json",
		bytes.NewBuffer([]byte(
			`{"id": 10, 
				"nome": "jonas Victor Alves da Silva",
				"endereco": "Av.Principal, 100, RT-PB",
				"data-de-nascimento": "09/03/1997",
				"cpf": "987.654.321-02",
				"telefone": 83987654321}`)))

	if err != nil {
		t.Errorf("Não foi possível fazer a requisição: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("Obrigatório preencher todos os campos: %v", err)
	}

	if resp.StatusCode == http.StatusCreated {
		resp, err := http.Get("http://localhost:8080/pessoa/10")
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
			log.Println(err)
			t.Error(err)
		}
		if resp.StatusCode != 200 {
			t.Errorf("Erro: %v", string(body))
		}
	} else {
		t.Errorf("Erro: %v", string(body))
	}
}

func TestIDZero(t *testing.T) {
	resp, err := http.Post(
		"http://localhost:8080/pessoa/",
		"application/json",
		bytes.NewReader([]byte(
			`{"id": 0,
				"nome": "Jonas Victor",
				"endereco": "Rua Principal JP-PB",
				"data-de-nascimento": "09/03/1997",
				"cpf": "989.878.767-10",
				"telefone": 83987654321}`)))

	if err != nil {
		t.Errorf("Não foi possível fazer a requisição: %v", err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		t.Errorf("Obrigatório preencher todos os campos: %v", err.Error())
	}

	if resp.StatusCode == 400 {
		t.Error(string(body))
	}
}

/*
func TestCreateIDNegativo(t *testing.T) {
	resp, err := http.Post(
		"http://localhost:8080/pessoa/",
		"application/json",
		bytes.NewReader([]byte(
			`{"id": -20,
				"nome": "Jonas",
				"endereco": "Rua Principal",
				"data-de-nascimento": "09/03/1997",
				"cpf": "989.878.767-10",
				"telefone": 83987654321}`)))

	if err != nil {
		t.Errorf("Error ao fazer requisição: %v", err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		t.Errorf("Erro ao preencher os campos: %v", err.Error())
	}

	if resp.StatusCode != 400 {
		t.Error(string(body))
	}

	if string(body) == "Pessoa não encontrada" {
		t.Error(string(body))
	}
}

func TestCampoVazio(t *testing.T) {
	resp, err := http.Post(
		"http://localhost:8080/pessoa/",
		"application/json",
		bytes.NewReader([]byte(
			`{"id": 60,
				"nome": "Jonas",
				"endereco": "",
				"data-de-nascimento": "",
				"cpf": "",
				"telefone": 83987654321}`)))

	if err != nil {
		t.Errorf("Error ao fazer requisição: %v", err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		t.Errorf("Erro ao preencher os campos: %v", err.Error())
	}

	if resp.StatusCode != 400 {
		t.Error(string(body))
	}

	if string(body) == "Pessoa não encontrada" {
		t.Error(string(body))
	}
}

//
func TestQuantDigitos(t *testing.T) {
	resp, err := http.Post(
		"http://localhost:8080/pessoa/",
		"application/json",
		bytes.NewReader([]byte(
			`{"id": 40,
				"nome": "Jonas",
				"endereco": "Rua Principal",
				"data-de-nascimento": "09/03/1997",
				"cpf": "989.878.767-10",
				"telefone": 839876543}`)))

	if err != nil {
		t.Errorf("Error ao fazer requisição: %v", err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		t.Errorf("Erro ao preencher os campos: %v", err.Error())
	}

	if resp.StatusCode != 400 {
		t.Error(string(body))
	}

	if string(body) == "Pessoa não encontrada" {
		t.Error(string(body))
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
		t.Error(err)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Sem sucesso: %v", string(body))
	}

	if string(body) == "Pessoa não encontrada" {
		t.Error(string(body))
	}
}

//
func TestGetPessoasErro(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/pessoa/10")
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	//log.Println(string(body))
	pess := Pessoa{}
	err = json.Unmarshal([]byte(string(body)), &pess)
	if err != nil {

		t.Error(err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("Sem sucesso: %v", string(body))
	}

	if string(body) == "Pessoa não encontrada" {
		t.Error(string(body))
	}

}

// Busca a pessoa cadastrada pelo ID
func TestGetByID(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/pessoa/20")
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
		log.Println(err)
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Sem sucesso: %v", string(body))
	}

	if string(body) == "Pessoa não encontrada" {
		t.Error(string(body))
	}
}

//
func TestGetByIDErro(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/pessoa/100")
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
		log.Println(err)
		t.Error(err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("Sem sucesso: %v", string(body))
	}

	if string(body) == "Pessoa não encontrada" {
		t.Error(string(body))
	}
}

// Atualiza os dados da pessoa cadastrada pelo ID
func TestUpdate(t *testing.T) {

	req, err := http.NewRequest(
		"PUT",
		"http://localhost:8080/pessoa/",
		bytes.NewBuffer([]byte(
			`{"id": 20,
				"nome": "Laynne",
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		t.Error(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Sem sucesso: %v", string(body))
	}
	if string(body) == "Pessoa não encontrada" {
		t.Error(string(body))
	}
}

//
func TestUpdateErro(t *testing.T) {
	req, err := http.NewRequest(
		"PUT",
		"http://localhost:8080/pessoa/",
		bytes.NewBuffer([]byte(
			`{"id": 200,
			"nome": "Lay",
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		t.Error(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Errorf("Sem sucesso: %v", string(body))
	}
	if string(body) == "Pessoa não encontrada" {
		t.Error(string(body))
	}
}

// Deleta a pessoa cadastrada através do ID
func TestDelete(t *testing.T) {

	req, err := http.NewRequest("DELETE", "http://localhost:8080/pessoa/50", nil)
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		t.Error(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Sem sucesso: %v", string(body))
	}
	if string(body) == "Pessoa não encontrada" {
		t.Error(string(body))
	}
}

func TestDeleteErro(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://localhost:8080/pessoa/10", nil)
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		t.Error(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Errorf("Sem sucesso: %v", string(body))
	}
	if string(body) == "Pessoa não encontrada" {
		t.Error(string(body))
	}
}
*/
