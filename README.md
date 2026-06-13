# Auction Go Expert

API de leilões desenvolvida em Go, com Gin e MongoDB.

O projeto permite criar e consultar leilões, registrar lances em lote e fechar leilões automaticamente após o tempo configurado.

## Requisitos

- Go 1.20
- Docker
- Docker Compose

## Como executar com Docker Compose

No diretório raiz do projeto, execute:

```bash
docker-compose up --build
```

A aplicação ficará disponível em:

```text
http://localhost:8080
```

O MongoDB ficará exposto em:

```text
localhost:27017
```

As variáveis de ambiente usadas pela aplicação estão em `cmd/auction/.env`.

## Variáveis principais

```env
BATCH_INSERT_INTERVAL=20s
MAX_BATCH_SIZE=4
AUCTION_INTERVAL=20s
AUCTION_DURATION=3m
MONGODB_URL=mongodb://admin:admin@mongodb:27017/auctions?authSource=admin
MONGODB_DB=auctions
```

- `AUCTION_DURATION`: tempo até o fechamento automático de um leilão criado.
- `BATCH_INSERT_INTERVAL`: intervalo máximo para persistir lances acumulados.
- `MAX_BATCH_SIZE`: quantidade máxima de lances antes de gravar o lote no MongoDB.

## Como executar os testes

```bash
go test ./...
```

O teste principal valida se um leilão é alterado para o status `Completed` após a duração configurada.

## Endpoints

### Criar leilão

```http
POST /auction
```

Exemplo:

```bash
curl -X POST http://localhost:8080/auction \
  -H "Content-Type: application/json" \
  -d '{
    "product_name": "Notebook",
    "category": "Eletronicos",
    "description": "Notebook usado em bom estado",
    "condition": 2
  }'
```

No payload atual, `condition` é um valor numérico. A validação da API aceita:

- `0`
- `1`
- `2`

Resposta esperada: `201 Created`.

### Listar leilões

```http
GET /auction?status=0
```

Também é possível filtrar por `category` e `productName`:

```bash
curl "http://localhost:8080/auction?status=0&category=Eletronicos&productName=Notebook"
```

Status:

- `0`: ativo
- `1`: completado

Observação: a implementação atual exige o parâmetro `status`; usando `status=0`, a listagem retorna os leilões sem aplicar filtro de status.

### Buscar leilão por ID

```http
GET /auction/{auctionId}
```

Exemplo:

```bash
curl http://localhost:8080/auction/<auctionId>
```

### Criar lance

```http
POST /bid
```

Exemplo:

```bash
curl -X POST http://localhost:8080/bid \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "<userId>",
    "auction_id": "<auctionId>",
    "amount": 150.50
  }'
```

Resposta esperada: `201 Created`.

### Listar lances de um leilão

```http
GET /bid/{auctionId}
```

### Buscar lance vencedor

```http
GET /auction/winner/{auctionId}
```

### Buscar usuário por ID

```http
GET /user/{userId}
```

## Teste manual do fechamento automático

1. Suba a aplicação:

```bash
docker-compose up --build
```

2. Crie um leilão com `POST /auction`.

3. Liste os leilões e copie o campo `id` do leilão criado:

```bash
curl "http://localhost:8080/auction?status=0"
```

4. Aguarde o tempo definido em `AUCTION_DURATION`.

5. Consulte o leilão pelo ID:

```bash
curl http://localhost:8080/auction/<auctionId>
```

6. Verifique o campo `status`. O valor esperado após o fechamento automático é `1`.

## Verificacao no MongoDB

Com o Docker Compose em execução, conecte no MongoDB:

```bash
mongosh "mongodb://admin:admin@localhost:27017/auctions?authSource=admin"
```

Consulte os leilões:

```javascript
use auctions
db.auctions.find().pretty()
```

Para buscar um leilão específico:

```javascript
db.auctions.find({ _id: "<auctionId>" }).pretty()
```
