package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	sercicoPessoa, err := pessoa.NovoServico("pessoa.json")
	if err != nil {
		fmt.Println("Erro ao tentar cadastrar pessoa")
		return
	}

	http.HandleFunc("/pessoa/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var pessoa dominio.Pessoa
			err := json.NewDecoder(r.Body).Decode(&pessoa)
			if err != nil {
				fmt.Printf("Erro ao tentar decodificar. O corpo deve ser um json. Erro: %s", err.Error())
				http.Error(w, "Erro ao tentar cadastrar pessoa", http.StatusBadRequest)
				return
			}

			if pessoa.ID <= 0 {
				http.Error(w, "Erro ao tentar cadastrar a pessoa, ID deve ser um número inteiro positivo", http.StatusBadRequest)
				return
			}

			// Criar pessoa
			err = sercicoPessoa.Create(pessoa)
			if err != nil {
				fmt.Printf("Erro ao tentar criar pessoa:%s", err.Error())
				http.Error(w, "Erro ao tentar criar pessoa", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			return
		}
		if r.Method == "GET" {
			// Lista todas as pessoas cadastradas
			path := strings.TrimPrefix(r.URL.Path, "/pessoa/")
			if path == "" {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-type", "application/json") // Antes de codar
				pessoas := sercicoPessoa.List()
				err := json.NewEncoder(w).Encode(pessoas)
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
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-type", "application/json")
				err = json.NewEncoder(w).Encode(pessoa)
				if err != nil {
					http.Error(w, "Erro ao tentar codificar pessoa como json", http.StatusInternalServerError)
					return
				}
			}

		}

		if r.Method == "PUT" {
			var pessoa dominio.Pessoa
			err := json.NewDecoder(r.Body).Decode(&pessoa)
			if err != nil {
				fmt.Printf("Erro ao tentar decodificar. O corpo deve ser um json. Erro: %s", err.Error())
				http.Error(w, "Erro ao tentar criar pessoa", http.StatusBadRequest)
				return
			}

			if pessoa.ID <= 0 {
				http.Error(w, "Erro ao tentar cadastrar a pessoa, ID deve ser um número inteiro positivo", http.StatusBadRequest)
				return
			}

			// Atualizar/Editar cadastro das pessoas
			err = sercicoPessoa.Update(pessoa)
			if err != nil {
				fmt.Printf("Erro ao tentar atualizar pessoa: %s", err.Error())
				http.Error(w, "Erro ao tentar atualizar pessoa", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			//return
		}

		if r.Method == "DELETE" {
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
				// Delete pessoa
				err = sercicoPessoa.DeleteByID(pessoaID)
				if err != nil {
					fmt.Printf("Erro ao tentar deletar o cadastro da pessoa: %s", err.Error())
					http.Error(w, "Erro ao tentar deletar o cadastro da pessoa", http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)

			}

		}

		//http.Error(w, "Não implementado", http.StatusInternalServerError)
	})

	http.ListenAndServe(":8080", nil)
}
