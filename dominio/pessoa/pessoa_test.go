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
	Status  string `json:"status"`
	Results []struct {
		NomeCompleto   string `json:"nome"`
		Endereco       string `json:"endereco"`
		DataNascimento string `json:"data-de-nascimento"`
		Cpf            string `json:"cpf"`
		ID             int    `json:"id"`
		Telefone       int    `json:"telefone"`
	} `json:"results"`
}

// Retorna o cadastro da pessoa como JSON, se não houver ID já existente
func TestCreate(t *testing.T) {

	resp, err := http.Post(
		"http://localhost:8080/pessoa/",
		"application/json",
		bytes.NewReader([]byte(
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

	if resp.StatusCode == 201 {
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
			t.Errorf("Não foi possível buscar: %v", string(body))
		}
	} else {
		t.Errorf("Não foi possível criar: %v", string(body))
	}
}

// Não deve cadastrar com o ID zero
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
		t.Errorf("Não foi possível fazer a requisição: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		t.Errorf("Obrigatório preencher todos os campos: %v", err)
	}

	if resp.StatusCode == 400 {
		t.Error(string(body))
	}
}

// Não deve cadastrar com o ID negativo
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
		t.Errorf("Não foi possível fazer a requisição: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		t.Errorf("Obrigatório preencher todos os campos: %v", err)
	}

	if resp.StatusCode == 400 {
		t.Error(string(body))
	}
}

// Não deve realizar o cadastro da pessoa com algum dos campos vazios
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
		t.Errorf("Não foi possível fazer a requisição: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		t.Errorf("Obrigatório preencher todos os campos: %v", err)
	}

	if resp.StatusCode == 500 {
		t.Errorf("Não foi possível criar: %v", string(body))
	}
}

// Não deve cadastrar pessoa com o número de telefone abaixo de 11 dígitos
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
		t.Errorf("Não foi possível fazer a requisição: %v", err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		t.Errorf("Obrigatório preencher todos os campos: %v", err)
	}

	if resp.StatusCode == 500 {
		t.Errorf("Não foi possível criar: %v", string(body))
	}

}

// Retorna uma lista com todos os cadastros
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
		t.Errorf("Não foi possível listar: %v", string(body))
	}
}

// Retorna uma busca dos dados pelo ID do cadastro informado
func TestGetByID(t *testing.T) {
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
	}
	if resp.StatusCode != 200 {
		t.Errorf("Não foi possível listar: %v", string(body))
	}
}

// Não retorna dados cadastrados se o ID inofrmado não estiver cadastrado
func TestGetByIDErro(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/pessoa/50")
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
	}
	if resp.StatusCode != 404 {
		t.Errorf("Não foi possível listar: %v", string(body))
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

	if err != nil {
		t.Errorf("Obrigatório preencher todos os campos: %v", err)
	}

	if resp.StatusCode == 200 {
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
		}
		if resp.StatusCode != 200 {
			t.Errorf("Não foi possível buscar: %v", string(body))
		}
	} else {
		t.Errorf("Não foi possível buscar: %v", string(body))
	}
}

// Não atualiza os dados da pessoa se o ID não estiver cadastrado
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

	if resp.StatusCode == 500 {
		t.Errorf("Não foi possível editar os dados: %v", string(body))
	}
}

// Remove o cadastro da pessoa pelo ID informado
func TestDelete(t *testing.T) {

	req, err := http.NewRequest("DELETE", "http://localhost:8080/pessoa/50", nil)
	if err != nil {
		t.Error("*******************************")
		t.Error(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error("------------------------------")
		t.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		t.Error(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		resp, err := http.Get("http://localhost:8080/pessoa/50")
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
		}
		if resp.StatusCode != 200 {
			t.Errorf("Não foi possível buscar: %v", string(body))
		}
	} else {
		t.Errorf("Não foi possível buscar: %v", string(body))
	}
}

// Retorna um erro por não identificar o ID informado
func TestDeleteErro(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://localhost:8080/pessoa/80", nil)
	if err != nil {
		t.Error("*******************************")
		t.Error(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error("------------------------------")
		t.Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		t.Error(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 500 {
		t.Errorf("Não foi possível remover: %v", string(body))
	}
}
