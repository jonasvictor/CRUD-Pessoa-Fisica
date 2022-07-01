package pessoa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jonasvictor/CRUD-Pessoa-Fisica/dominio"
)

type Servico struct {
	dbFilePath string
	pessoas    dominio.Pessoas
}

func NovoServico(dbFilePath string) (Servico, error) {
	// Verifica se o arquivo existe
	_, err := os.Stat(dbFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Se não exixtir, crio arquivo vazio
			err := criaArquivoVazio(dbFilePath)
			if err != nil {
				return Servico{}, err
			}
			return Servico{
				dbFilePath: dbFilePath,
				pessoas:    dominio.Pessoas{},
			}, nil

		}
	}

	// Se existir, leio o arquivo e atualizo a variável 'pessoas' do serviço com as pessoas do arquivo
	jsonArquivo, err := os.Open(dbFilePath)
	if err != nil {
		return Servico{}, fmt.Errorf("Erro ao tentar abrir arquivo que contém todas as pessoas: %s", err.Error())
	}

	jsonArquivoByte, err := ioutil.ReadAll(jsonArquivo)
	if err != nil {
		return Servico{}, fmt.Errorf("Erro ao tentar ler o arquivo: %s", err.Error())
	}

	var todasPessoas dominio.Pessoas
	json.Unmarshal(jsonArquivoByte, &todasPessoas)

	return Servico{
		dbFilePath: dbFilePath,
		pessoas:    todasPessoas,
	}, nil

}

func criaArquivoVazio(dbFilePath string) error {
	var pessoas dominio.Pessoas = dominio.Pessoas{
		Pessoas: []dominio.Pessoa{},
	}
	pessoasJSON, err := json.Marshal(pessoas)
	if err != nil {
		return fmt.Errorf("Erro ao tentar codificar pessoas como JSON?: %s", err.Error())
	}

	err = ioutil.WriteFile(dbFilePath, pessoasJSON, 0755)
	if err != nil {
		return fmt.Errorf("Erro ao tentar gravar no arquivo. Erro: %s", err.Error())
	}

	return nil

}

func (s *Servico) Create(pessoa dominio.Pessoa) error {
	// Verificar se a pessoa já existe, se já existe então retorna um erro
	if s.existe(pessoa) {
		return fmt.Errorf("Erro ao tentar criar pessoa. Já existe uma pessoa com este ID cadastrada")
	}

	// Adiciona a pessoa na slice de pessoas
	s.pessoas.Pessoas = append(s.pessoas.Pessoas, pessoa)

	// Salvo o arquivo
	err := s.salvaArquivo()
	if err != nil {
		return fmt.Errorf("Erro ao tentar salvar arquivo no método 'Create'. Erro: %s", err.Error())
	}

	return nil
}

// Verifica se a pessoa já existe pelo ID
func (s Servico) existe(pessoa dominio.Pessoa) bool {
	for _, pessoaAtual := range s.pessoas.Pessoas {
		if pessoaAtual.ID == pessoa.ID {
			return true
		}
	}
	return false
}

func (s Servico) salvaArquivo() error {
	todasPessoasJSON, err := json.Marshal(s.pessoas)
	if err != nil {
		return fmt.Errorf("Error trying to encode people as json: %s", err.Error())
	}
	return ioutil.WriteFile(s.dbFilePath, todasPessoasJSON, 0755)
}

func (s Servico) List() dominio.Pessoas {
	return s.pessoas
}

func (s Servico) GetByID(pessoaID int) (dominio.Pessoa, error) {
	for _, pessoaAtual := range s.pessoas.Pessoas {
		if pessoaAtual.ID == pessoaID {
			return pessoaAtual, nil
		}
	}
	return dominio.Pessoa{}, fmt.Errorf("Pessoa não encontrada")
}

func (s *Servico) Update(pessoa dominio.Pessoa) error {
	var indexUpdate int = -1
	for index, pessoaAtual := range s.pessoas.Pessoas {
		if pessoaAtual.ID == pessoa.ID {
			indexUpdate = index
			break
		}
	}
	if indexUpdate < 0 {
		return fmt.Errorf("Não há nenhuma pessoa com o ID fornecido em nosso banco de dados")
	}

	s.pessoas.Pessoas[indexUpdate] = pessoa
	return s.salvaArquivo()

}

func (s *Servico) DeleteByID(pessoaID int) error {
	var indexDelete int = -1
	for index, pessoaAtual := range s.pessoas.Pessoas {
		if pessoaAtual.ID == pessoaID {
			indexDelete = index
			break
		}
	}
	if indexDelete < 0 {
		return fmt.Errorf("Não há nenhuma pessoa com o ID fornecido em nosso banco de dados")
	}

	s.pessoas.Pessoas = append(s.pessoas.Pessoas[:indexDelete], s.pessoas.Pessoas[indexDelete+1:]...)

	return s.salvaArquivo()
}
