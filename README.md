# TreeHouse API (Gestio)

Backend service responsável pelo funcionamento do sistema **TreeHouse**, parte da plataforma **Gestio**, voltada para gestão educacional.

Esta API fornece a lógica de negócio e os endpoints necessários para gerenciar usuários, autenticação, alunos, turmas, aulas e contratos dentro da plataforma.

O projeto foi desenvolvido em **Go (Golang)** com foco em organização modular, clareza arquitetural e facilidade de manutenção.

---

# Tecnologias Utilizadas

* **Golang**
* **Gorilla Mux** (roteamento HTTP)
* **MySQL / MariaDB**
* **SQL puro (sem ORM)**
* **JWT para autenticação**
* Arquitetura baseada em camadas

---

# Estrutura do Projeto

A aplicação segue uma separação de responsabilidades entre camadas da aplicação.

```text
.
├── migrations.sql          # Script de criação e estrutura do banco de dados
├── main.go                 # Ponto de entrada da aplicação
├── go.mod                  # Dependências do projeto
├── .env                    # Variáveis de ambiente
├── .gitignore
│
└── src
    ├── authentication      # Geração e validação de tokens
    │   └── token.go
    │
    ├── config              # Configurações da aplicação
    │   └── config.go
    │
    ├── controllers         # Camada HTTP (handlers das rotas)
    │   ├── login.go
    │   └── users.go
    │
    ├── middlewares         # Middlewares HTTP (autenticação, validações etc)
    │   └── middlewares.go
    │
    ├── models              # Estruturas de dados da aplicação
    │   └── users.go
    │
    ├── persistency         # Conexão com banco de dados
    │   └── database.go
    │
    ├── pkg                 # Pacotes reutilizáveis
    │
    ├── repository          # Camada de acesso ao banco (queries SQL)
    │   └── users.go
    │
    ├── responses           # Padronização das respostas HTTP
    │   └── responses.go
    │
    ├── router              # Definição das rotas da aplicação
    │   ├── router.go
    │   └── routes
    │       ├── login.go
    │       ├── routes.go
    │       └── user.go
    │
    ├── security            # Hash de senha e funções de segurança
    │   └── security.go
    │
    └── utils               # Funções utilitárias
```

---

# Arquitetura da Aplicação

A API segue um fluxo de camadas bem definido:

```
HTTP Request
      │
      ▼
Router
      │
      ▼
Controllers
      │
      ▼
Repository
      │
      ▼
Database
```

Cada camada possui uma responsabilidade específica:

| Camada      | Responsabilidade            |
| ----------- | --------------------------- |
| Router      | Define os endpoints da API  |
| Controllers | Processa requisições HTTP   |
| Repository  | Executa queries no banco    |
| Models      | Define estruturas de dados  |
| Middlewares | Intercepta requisições      |
| Security    | Autenticação e criptografia |
| Persistency | Gerencia conexão com banco  |

---

# Banco de Dados

O banco **treehousedb** foi projetado para suportar a gestão completa da escola.

Principais entidades:

* **usuarios** – professores e administradores do sistema
* **permissoes** – controle de acesso
* **clientes** – responsáveis financeiros
* **alunos** – estudantes
* **turmas** – grupos de alunos
* **aulas** – registros de aulas
* **contratos** – contratos educacionais
* **enderecos** – normalização de endereços
* **logs_auditoria** – auditoria de alterações

Além disso existem tabelas de relacionamento:

* usuarios_permissoes
* enderecos_clientes
* clientes_alunos
* alunos_turmas
* alunos_aulas

O banco também possui **triggers automáticas para integridade e auditoria**.

---

# Como Rodar o Projeto

### 1. Clonar o repositório

```bash
git clone https://github.com/brunoob35/TreeHouse-API.git
cd TreeHouse-API
```

---

### 2. Configurar variáveis de ambiente

Crie um arquivo `.env`:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=treehousedb
```

---

### 3. Criar o banco de dados

Execute o script:

```
migrations.sql
```

---

### 4. Rodar a aplicação

```bash
go run main.go
```

---

# Endpoints

### Autenticação

| Método | Endpoint | Descrição               |
| ------ | -------- | ----------------------- |
| POST   | /login   | Autenticação de usuário |

### Usuários

| Método | Endpoint       | Descrição         |
| ------ | -------------- | ----------------- |
| POST   | /usuarios      | Criar usuário     |
| GET    | /usuarios      | Listar usuários   |
| GET    | /usuarios/{id} | Buscar usuário    |
| PUT    | /usuarios/{id} | Atualizar usuário |
| DELETE | /usuarios/{id} | Remover usuário   |

---

# Objetivo do Projeto

O **TreeHouse API** foi desenvolvido para servir como backend da plataforma **Gestio**, com foco em:

* gestão educacional
* controle de alunos e turmas
* registro de aulas
* gestão de contratos
* controle de usuários e permissões

A arquitetura foi pensada para permitir crescimento e evolução do sistema de forma organizada.

---

# Autores
- Bruno Schmaiske Quoos
