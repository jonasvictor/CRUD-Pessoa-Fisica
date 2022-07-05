package dominio

type Pessoa struct {
	ID             int    `json:"id"`
	NomeCompleto   string `json:"nome"`
	Endereco       string `json:"endereco"`
	DataNascimento string `json:"data-de-nascimento"` //time
	Cpf            string `json:"cpf"`                //int
	Telefone       int    `json:"telefone"`           //int
}

type Pessoas struct {
	Pessoas []Pessoa `json:"pessoas"`
}
