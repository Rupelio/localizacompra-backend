package main

import (
	"context"
	"fmt"
	"localiza-compra/backend/internal/api/user"
	"localiza-compra/backend/internal/database"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 4 {
		fmt.Println("Uso: go run ./cmd/cli/main.go promote <email> <role>")
		return
	}

	command := args[1]
	email := args[2]
	role := args[3]

	if command != "promote" {
		fmt.Println("Comando inválido. Use 'promote'.")
		return
	}

	fmt.Printf("A promover o utilizador %s para o cargo %s...\n", email, role)

	db := database.Connect()
	defer db.Close()

	userRepo := user.NewRepository(db)

	err := userRepo.UpdateRole(context.Background(), email, role)
	if err != nil {
		fmt.Printf("Erro ao promover o utilizador: %v\n", err)
		return
	}

	fmt.Println("Utilizador promovido com sucesso!")
}
