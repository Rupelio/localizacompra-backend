package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() *pgxpool.Pool {
	// Pega a URL do banco de dados a partir das variáveis de ambiente que definimos no docker
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("DATABASE_URL não foi definida.")
	}

	dbpool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("Não foi possível conectar ao banco de dados: %v\n", err)
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Ping para o banco de dados falhou: %v\n", err)
	}

	log.Println("Conexão com o banco de dados estabelecida com sucesso.")
	return dbpool
}
