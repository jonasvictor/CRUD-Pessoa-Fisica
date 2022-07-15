package pessoa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

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

		} else {
			return Servico{}, err
		}
	}

	// Se existir, leio o arquivo e atualizo a variável 'pessoas' do serviço com as pessoas do arquivo
	jsonArquivo, err := os.Open(dbFilePath)
	if err != nil {
		return Servico{}, fmt.Errorf("Erro ao tentar abrir arquivo que contém todas as pessoas: %s\n", err.Error())
	}

	jsonArquivoByte, err := ioutil.ReadAll(jsonArquivo)
	if err != nil {
		return Servico{}, fmt.Errorf("Erro ao tentar ler o arquivo: %s\n", err.Error())
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
		return fmt.Errorf("Erro ao tentar codificar pessoas como JSON: %s\n", err.Error())
	}

	err = ioutil.WriteFile(dbFilePath, pessoasJSON, 0755)
	if err != nil {
		return fmt.Errorf("Erro ao tentar gravar no arquivo. Erro: %s\n", err.Error())
	}

	return nil

}

func (s Servico) campoVazio(pessoa dominio.Pessoa) bool {
	// Verifica se todos os campos foram preechidos
	if pessoa.NomeCompleto == "" || pessoa.Endereco == "" || pessoa.DataNascimento == "" || pessoa.Cpf == "" || pessoa.Telefone == 0 {
		return true
	}
	return false
}

func (s Servico) quantDigitos(pessoa dominio.Pessoa) bool {
	// Verifica se telefone os 11 dígitos do telefone foram preechidos
	digitosTelefone := strconv.Itoa(pessoa.Telefone)
	if pessoa.Telefone != 0 {
		if len(digitosTelefone) != 11 {
			return true
		}
	}
	return false
}

func (s *Servico) Create(pessoa dominio.Pessoa) error {
	// Verifica se a pessoa já existe, se já existir, então retorna um erro
	if s.existe(pessoa) {
		return fmt.Errorf("Já existe uma pessoa com este ID cadastrado")
	}
	if s.campoVazio(pessoa) {
		return fmt.Errorf("Erro ao tentar criar pessoa, dados insuficientes")
	}

	// Verifica quantidade de números inseridos em telefone
	if s.quantDigitos(pessoa) {
		return fmt.Errorf("Erro ao tentar criar pessoa. Número de telefone incompleto/incorreto")
	}

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
	// Salva os dados no cadastro
	todasPessoasJSON, err := json.Marshal(s.pessoas)
	if err != nil {
		return fmt.Errorf("Erro ao tentar codificar pessoas como json: %s", err.Error())
	}
	return ioutil.WriteFile(s.dbFilePath, todasPessoasJSON, 0755)
}

func (s Servico) List() dominio.Pessoas {
	return s.pessoas
}

func (s Servico) GetByID(pessoaID int) (dominio.Pessoa, error) {
	// Busca cadastro de pessoa pelo ID informado
	for _, pessoaAtual := range s.pessoas.Pessoas {
		if pessoaAtual.ID == pessoaID {
			return pessoaAtual, nil
		}
	}
	return dominio.Pessoa{}, fmt.Errorf("Pessoa não encontrada")
}

func (s *Servico) Update(pessoa dominio.Pessoa) error {
	// Atualiza o cadastro da pessoa pelo ID informado
	var indexUpdate int = -1
	for index, pessoaAtual := range s.pessoas.Pessoas {
		if pessoaAtual.ID == pessoa.ID {
			indexUpdate = index
			break
		}
	}
	if indexUpdate < 0 {
		return fmt.Errorf("Não há pessoa com o ID fornecido")
	}

	s.pessoas.Pessoas[indexUpdate] = pessoa
	return s.salvaArquivo()

}

func (s *Servico) DeleteByID(pessoaID int) error {
	// Deleta cadastro da pessoa pelo ID informado
	var indexDelete int = -1
	for index, pessoaAtual := range s.pessoas.Pessoas {
		if pessoaAtual.ID == pessoaID {
			indexDelete = index
			break
		}
	}
	if indexDelete < 0 {
		return fmt.Errorf("Não há pessoa com o ID fornecido")
	}

	s.pessoas.Pessoas = append(
		s.pessoas.Pessoas[:indexDelete],
		s.pessoas.Pessoas[indexDelete+1:]...,
	)

	return s.salvaArquivo()
}
