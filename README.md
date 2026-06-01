# Auction Go Expert

Serviço de leilões em Go usando Gin e MongoDB.

## Descrição

Projeto simples que permite criar leilões, consultar leilões, criar lances e fechar leilões automaticamente após um tempo configurado.

## Requisitos

- Go 1.20
- Docker
- Docker Compose

## Como executar os testes

No diretório do projeto:

```bash
go test ./...
```

Isso executa todos os testes do projeto.

## Como executar por Docker Compose

No diretório do projeto:

```bash
docker-compose up --build
```

Isso irá:

- subir a aplicação em `http://localhost:8080`
- subir o MongoDB em `localhost:27017`

O arquivo de ambiente usado é `cmd/auction/.env`.

## Teste manual para verificar o fechamento do leilão

### 1. Criar um leilão

Use o endpoint `POST /auction` com JSON:

```bash
curl -X POST http://localhost:8080/auction \
  -H "Content-Type: application/json" \
  -d '{
    "product_name": "Test Product",
    "category": "Tests",
    "description": "Test description",
    "condition": 1
  }'
```

### 2. Obter o ID do leilão

A resposta deve retornar o ID do leilão criado. Anote esse `auctionId`.

### 3. Aguardar o fechamento automático

O tempo de fechamento está configurado em `cmd/auction/.env` via a variável:

```env
AUCTION_DURATION=1m
```

Nesse exemplo, o leilão deve ser encerrado automaticamente após 1 minuto.

### 4. Consultar o leilão pelo ID

Use o endpoint `GET /auction/{auctionId}`:

```bash
curl http://localhost:8080/auction/<auctionId>
```

Substitua `<auctionId>` pelo ID retornado na criação.

### 5. Verificar o status

Quando o leilão estiver fechado, o campo `status` deve indicar que ele foi completado. ( statu = 1)

### 6. Verificar diretamente no MongoDB (opcional)

Conecte ao shell do MongoDB e use:

```bash
mongosh "mongodb://admin:admin@localhost:27017/auctions?authSource=admin"
use auctions
show collections
```

Em seguida, busque o leilão criado:

```bash
db.auctions.find({ _id: "<auctionId>" }).pretty()
```

Substitua `<auctionId>` pelo ID real do leilão.

## Observações

A aplicação usa a coleção `auctions` no MongoDB e o endpoint principal para criar leilões é `POST /auction`.
