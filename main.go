package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jonasvictor/CRUD-Pessoa-Fisica/dominio"
	"github.com/jonasvictor/CRUD-Pessoa-Fisica/dominio/pessoa"
)

func main() {
	sercicoPessoa, err := pessoa.NovoServico("pessoa.json")
	if err != nil {
		fmt.Printf("Erro ao tentar cadastrar pessoa: %s\n", err.Error())
		return
	}

	http.HandleFunc("/pessoa/", func(w http.ResponseWriter, r *http.Request) {
		// Adiciona o cadastro da pessoa
		if r.Method == "POST" {
			var pessoa dominio.Pessoa
			err := json.NewDecoder(r.Body).Decode(&pessoa)
			if err != nil {
				fmt.Printf("Erro ao tentar decodificar. O corpo deve ser um json. Erro: %s\n", err.Error())
				http.Error(w, "Erro ao tentar cadastrar pessoa", http.StatusBadRequest)
				return
			}

			if pessoa.ID <= 0 {
				http.Error(w, "O ID da pessoa deve ser um número inteiro positivo", http.StatusBadRequest)
				return
			}

			// Criar pessoa
			err = sercicoPessoa.Create(pessoa)
			if err != nil {
				fmt.Printf("Erro ao tentar criar pessoa: %s\n", err.Error())
				http.Error(w, "Erro ao tentar criar pessoa", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			return
		}

		if r.Method == "PUT" {
			// Atualiza os dados de cadastro da pessoa
			var pessoa dominio.Pessoa
			err := json.NewDecoder(r.Body).Decode(&pessoa)
			if err != nil {
				fmt.Printf("O corpo deve ser um json. Erro: %s\n", err.Error())
				http.Error(w, "Erro ao tentar criar pessoa", http.StatusBadRequest)
				return
			}

			if pessoa.ID <= 0 {
				http.Error(w, "O ID da pessoa deve ser um número inteiro positivo", http.StatusBadRequest)
				return
			}

			err = sercicoPessoa.Update(pessoa)
			if err != nil {
				fmt.Printf("Erro ao tentar atualizar a pessoa: %s\n", err.Error())
				http.Error(w, "Erro ao tentar atualizar pessoa", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == "GET" {
			// Lista todas as pessoas cadastradas
			path := strings.TrimPrefix(r.URL.Path, "/pessoa/")
			if path == "" {
				w.Header().Set("Content-Type", "application/json") // Antes de codar
				w.WriteHeader(http.StatusOK)

				err := json.NewEncoder(w).Encode(sercicoPessoa.List())
				if err != nil {
					http.Error(w, "Erro ao tentar listar pessoas", http.StatusInternalServerError)
					return
				}

			} else {
				// Lista as pessoas cadastradas por ID
				pessoaID, err := strconv.Atoi(path)
				if err != nil {
					http.Error(w, "ID fornecido é inválido. O ID da pessoa deve ser um número inteiro", http.StatusBadRequest)
					return
				}
				pessoa, err := sercicoPessoa.GetByID(pessoaID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				err = json.NewEncoder(w).Encode(pessoa)
				if err != nil {
					http.Error(w, "Erro ao tentar encontrar pessoa", http.StatusInternalServerError)
					return
				}
			}
			return

		}

		if r.Method == "DELETE" {
			// Remove cadastro da pessoa
			path := strings.TrimPrefix(r.URL.Path, "/pessoa/")
			if path == "" {
				http.Error(w, "O ID deve ser fornecido no URL", http.StatusBadRequest)
				return
			} else {
				pessoaID, err := strconv.Atoi(path)
				if err != nil {
					http.Error(w, "ID fornecido é inválido. O ID da pessoa deve ser um número inteiro", http.StatusBadRequest)
					return
				}
				err = sercicoPessoa.DeleteByID(pessoaID)
				if err != nil {
					fmt.Printf("Erro ao tentar deletar o cadastro da pessoa: %s\n", err.Error())
					http.Error(w, "Erro ao tentar deletar a pessoa", http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)

			}
			return

		}

	})

	http.ListenAndServe(":8080", nil)
}
