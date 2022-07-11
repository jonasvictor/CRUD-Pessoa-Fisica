package dominio

type Pessoa struct {
	ID             int    `json:"id"`
	NomeCompleto   string `json:"nome"`
	Endereco       string `json:"endereco"`
	DataNascimento string `json:"data-de-nascimento"`
	Cpf            string `json:"cpf"`
	Telefone       int    `json:"telefone"`
}

type Pessoas struct {
	Pessoas []Pessoa `json:"pessoas"`
}
