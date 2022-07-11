# CRUD-Pessoa-Fisica

CRUD feito em Go com TDD sem a utilização de um banco de dados, permitindo criar, listar, atualizar e remover cadastros.

curl -XPOST localhost:8080/pessoa/ -d '{"id": 1, "nome": "jonas' (Faltando fechar as aspas duplas ou fechar as chaves no nome, dá erro na decodificação do JSON)

curl -XPOST localhost:8080/pessoa/ -d '{"id": 1, "nome": "jonas"}' (Tudo ok)

Execuntando os comandos no terminal.
