# mini-redis-go

Um servidor Redis minimalista implementado em Go, inspirado no funcionamento do Redis real. Ideal para estudos, desafios técnicos e compreensão de protocolos e comandos básicos do Redis.

## Funcionalidades

- **Protocolo RESP2**: Aceita conexões TCP e interpreta comandos no formato RESP2.
- **Comandos suportados**:
  - `PING`
  - `ECHO <mensagem>`
  - `SET <chave> <valor> [PX <milissegundos>]`
  - `GET <chave>`
  - `KEYS *`
  - `CONFIG GET <DIR|DBFILENAME>`
  - `SAVE`
- **Expiração de chaves**: Suporte ao parâmetro `PX` no comando `SET`.
- **Persistência**: Comando `SAVE` para salvar o banco de dados em disco.
- **Configuração via argumentos**: Defina diretório e nome do arquivo de banco com `--dir` e `--dbfilename`.

## Como rodar

```bash
go run app/main.go --dir ./data --dbfilename dump.rdb
```

O servidor irá escutar na porta padrão `6379`.

## Exemplos de uso

```bash
# PING
$ redis-cli ping
PONG

# SET e GET
$ redis-cli set foo bar
OK
$ redis-cli get foo
bar

# ECHO
$ redis-cli echo "Olá, mundo!"
Olá, mundo!

# KEYS
$ redis-cli keys *
1) "foo"
```

## Estrutura do projeto

```
app/
  commands/         # Implementação dos comandos Redis
  protocol_parser/  # Parser do protocolo RESP2
  server_config/    # Configuração do servidor
  main.go           # Ponto de entrada do servidor
```

## Créditos

Desenvolvido como parte do desafio [Codecrafters Redis Go](https://codecrafters.io/challenges/redis). 