# ⏳ TimeSheet PRO

> Sistema moderno de gestão de ponto e controle de horas para empresas (Multi-tenant).

**TimeSheet PRO** é uma aplicação *fullstack* desenvolvida em Go que substitui planilhas manuais por uma plataforma digital centralizada. O projeto utiliza **Server-Side Rendering (SSR)** com **Templ** e **HTMX** para oferecer uma experiência de usuário ágil e dinâmica, sem a complexidade de SPAs pesadas.

[Link](https://timesheet-pro.onrender.com) para a aplicação em produção 

## 🚀 Tecnologias Utilizadas

* **Backend:** [Go](https://go.dev/) (Golang)
* **Framework Web:** [Gin Gonic](https://github.com/gin-gonic/gin)
* **Template Engine:** [Templ](https://templ.guide/) (Type-safe HTML para Go)
* **Interatividade:** [HTMX](https://htmx.org/) (AJAX, CSS Transitions, WebSockets via HTML)
* **Banco de Dados:** [PostgreSQL](https://www.postgresql.org/) (Driver: `pgx`)
* **Estilização:** [TailwindCSS](https://tailwindcss.com/)
* **Integrações:** API ViaCEP (Autocompletar endereços)

---

## ✨ Funcionalidades Principais

* **Autenticação:** Cadastro e Login de usuários (JWT).
* **Multi-tenancy:** Criação e gestão de múltiplas Organizações.
* **Gestão de Membros:** Convite e remoção de membros, com papéis (Admin/Member).
* **Endereçamento Inteligente:** Preenchimento automático de endereço da empresa via CEP.
* **Controle de Ponto:** Registro de entradas e saídas (Daily Timesheets).
* **Relatórios:** Painéis administrativos para gestão de horas.

---

## 🛠️ Pré-requisitos

Antes de começar, certifique-se de ter instalado em sua máquina:

* [Go](https://go.dev/dl/) (Versão 1.23 ou superior)
* [PostgreSQL](https://www.postgresql.org/download/) (Ou rodando via Docker)
* [Make](https://www.gnu.org/software/make/) (Para rodar os comandos do Makefile)

---

## ⚙️ Configuração e Instalação

### 1. Clone o repositório
```bash
git clone [https://github.com/seu-usuario/timesheet-pro.git](https://github.com/seu-usuario/timesheet-pro.git)
cd timesheet-pro
````

### 2\. Instale as ferramentas de desenvolvimento

O projeto possui um comando `make` configurado para baixar o **Templ**, **Goose** e **Swag** automaticamente.

```bash
make setup
```

### 3\. Configure as Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto. Você pode usar o `.env.example` como base (se houver) ou configurar as seguintes variáveis:

```env
# Configuração do Servidor
PORT=8080
GIN_MODE=debug

# Configuração do Banco de Dados
DATABASE_URL=postgres://usuario:senha@localhost:5432/timesheet_db?sslmode=disable

# Configuração de Migrations (Goose)
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://usuario:senha@localhost:5432/timesheet_db?sslmode=disable

# Segurança
JWT_SECRET=sua_chave_secreta_aqui
```

-----

## 🗄️ Banco de Dados

Certifique-se de que seu PostgreSQL está rodando e que o banco de dados (ex: `timesheet_db`) foi criado.

### Rodar Migrações

Para criar as tabelas necessárias no banco de dados, utilize o comando configurado no Makefile:

```bash
make migrations/up
```

> **Nota:** Para criar uma nova migração no futuro, use:
> `make migrations/new name=nome_da_migracao`

-----

## ▶️ Como Rodar Localmente

### Usando Make (Padrão)

Este comando irá gerar os arquivos do Templ (`templ generate`) e iniciar o servidor Go:

```bash
make run
```

Acesse: `http://localhost:8080`

-----

## 📂 Estrutura do Projeto

```
.
├── cmd/
│   └── api/           # Ponto de entrada (main.go)
├── internal/
│   ├── domain/        # Modelos e Regras de Negócio (Structs)
│   ├── server/        # Configuração do Gin e Rotas
│   ├── service/       # Lógica de Aplicação
│   ├── repository/    # Acesso ao Banco de Dados (Queries SQL/PGX)
│   │   └── migrations # Arquivos .sql do Goose
│   └── templates/     # Componentes de UI (Arquivos .templ)
│       ├── components
│       ├── layouts
│       └── pages
├── Makefile           # Automação de tarefas
└── README.md
```

## 🤝 Contribuição

1.  Faça um Fork do projeto
2.  Crie uma Branch para sua Feature (`git checkout -b feature/MinhaFeature`)
3.  Faça o Commit (`git commit -m 'Adicionando funcionalidade X'`)
4.  Faça o Push (`git push origin feature/MinhaFeature`)
5.  Abra um Pull Request
