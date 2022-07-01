package pessoa_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jonasvictor/CRUD-Pessoa-Fisica/dominio"
)

func TestCreate(t *testing.T) {
	pessoa := dominio.Pessoa{
		ID:             1,
		NomeCompleto:   "Jonas Victor Alves da Silva",
		Endereco:       "Rua Cônego Fautino, 100, Itapororoca-PB",
		DataNascimento: "09/03/997",
		Cpf:            "999.999.999-99",
		Telefone:       "(83)98888-8888",
	}

	bodyRequestJson := new(bytes.Buffer)
	encodeJson, erro := json.Marshal(pessoa)
	if erro != nil {
		return erro
	}
	bodyRequestJson.Write(encodeJson)

	request, erro := http.NewRequest("POST", "/pessoa/", bodyRequestJson)
	// Executando o request
	client := &http.Client{}
	resp, erro := client.Do(request)
	if erro != nil {
		return erro
		// Alertar o erro
		// Tratar o erro

	}

	// Testar as respostas dos usuarios

}
