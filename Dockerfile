# --- Estágio 1: Build ---
# Usamos uma imagem oficial do Go como nossa imagem de "builder".
# 'alpine' é uma versão leve do Linux, boa para manter o tamanho baixo.
FROM golang:1.24-alpine AS builder

# Define o diretório de trabalho dentro do contêiner.
WORKDIR /app

# Copia os arquivos de gerenciamento de dependências primeiro.
# O Docker armazena em cache as camadas, então isso acelera builds futuros se as dependências não mudarem.
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o resto do código-fonte do nosso backend.
COPY . .

# Compila a aplicação.
# CGO_ENABLED=0 desabilita o CGO, criando um binário estaticamente vinculado.
# -o ./out/api aponta para o caminho do executável compilado.
# ./cmd/api é o caminho para o nosso pacote main.
RUN CGO_ENABLED=0 GOOS=linux go build -o ./out/api ./cmd/server

# --- Estágio 2: Final ---
# Começamos com uma imagem "do zero" (scratch), que é a menor possível, ou alpine.
# Alpine é uma ótima escolha pois é muito pequena mas ainda nos dá um shell, o que pode ajudar na depuração.
FROM alpine:latest

# Define o diretório de trabalho.
WORKDIR /app

# Copia APENAS o executável compilado do estágio 'builder'.
COPY --from=builder /app/out/api .

# Expõe a porta 8080, informando ao Docker que a aplicação dentro do contêiner escuta nesta porta.
EXPOSE 8080

# O comando que será executado quando o contêiner iniciar.
CMD ["./api"]
