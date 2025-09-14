# Cep Temperatuer - Google Cloud Run
## Desafio da Pós Go Expert 
Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin). Esse sistema deverá ser publicado no Google Cloud Run.

# Rodando o Projeto

## Pré-reqs

### API KEY
Para prosseguir, é necessário ter uma API KEY do https://www.weatherapi.com/

Com a API KEY em mãos, preencha o arquivo `.env.example` que se encontra na raíz do projeto, e então renomeie-o para `.env`

**[Importante]**: O arquivo `.env` já está no `.gitignore`, mas o arquivo `.env.example` não está, *NÃO* coloque dados sensíveis no arquivo `.env.example`.

### Rodando local sem Docker
Na raiz do projeto, execute o comando:
```
WEATHER_API_KEY=YOUR_API_KEY go run cmd/main.go
```

### Rodando local com Docker
Na raiz do projeto, execute os comandos:
```
docker build -t your-user/cep-temperature .
docker run --rm --env-file .env -p 8080:8080 your-user/cep-temperature:latest
```

### Rodando local com Docker-compose
Na raiz do projeto, execute o comando:
```
docker compose up -d
```

### Executando no Google Cloud Run
Para executar no Google Cloud Run, é só acessar a URL abaixo:
```
https://lab-google-cloud-run-514874811791.us-central1.run.app/cep-temperature/{cep}
```

### Utilitários
Para auxiliar nos testes, existe um arquivo chamado `cep-temperature.http` na pasta `test`.
A partir desse arquivo, é possível executar chamados tanto locais, quanto direto no Google Cloud Run.
Use a Extensão REST Client disponível no VS Code para usar este arquivo.

# Testes unitários
Os testes unitários podem ser rodados executando o comando abaixo na raíz do projeto:
```
go test ./...
```