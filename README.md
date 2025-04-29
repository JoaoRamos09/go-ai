# Go-AI

Uma API em Go que integra serviços de IA com funcionalidades de autenticação e autorização. O projeto utiliza OpenAI para processamento de linguagem natural e Pinecone para armazenamento e busca de vetores.

## Características Principais

- Integração com OpenAI para processamento de IA
- Sistema de vetorização e busca com Pinecone
- Autenticação JWT
- Sistema de roles (Admin/User)
- API RESTful usando Chi router
- Banco de dados PostgreSQL

## Requisitos

- Go 1.22.3+
- PostgreSQL
- Docker
- Chaves de API:
  - OpenAI
  - Pinecone

## Configuração

1. Clone o repositório
2. Copie o arquivo `.env.example` para `.env`
3. Configure as variáveis de ambiente no arquivo `.env`:
   - Configurações do banco de dados
   - Chaves de API (OpenAI e Pinecone)
   - Configurações JWT
   - Configurações do servidor

### Autenticação
- Utiliza JWT (JSON Web Tokens)
- Token deve ser enviado no header `Authorization` no formato `Bearer <token>`
- Tokens expiram em 1 hora

### Autorização
O sistema possui dois níveis de acesso:
- **USER**: Acesso básico às funcionalidades de IA
- **ADMIN**: Acesso completo, incluindo gerenciamento de documentos

## Desenvolvimento

1. Inicie o banco de dados:
```bash
docker-compose up -d
```

2. Execute o servidor:
```bash
make run
```

## Estrutura do Projeto

- `/cmd/api`: Ponto de entrada da aplicação
- `/config`: Configurações da aplicação
- `/internal`: Código interno da aplicação
  - `/database`: Conexões com bancos de dados
  - `/http`: Handlers HTTP e middlewares
  - `/user`: Lógica de usuários e autenticação
- `/ia`: Integrações com IA
- `/docs`: Documentação adicional

## Segurança

- Senhas são hasheadas usando bcrypt
- Tokens JWT são assinados com chave secreta
- Todas as rotas de IA requerem autenticação
- Sistema de roles para controle de acesso
- Validação de entrada em todas as requisições