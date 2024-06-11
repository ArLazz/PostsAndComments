package main

import (
	"flag"
	"log"
	"postsandcomments/configs"
	"postsandcomments/internal/db"
	"postsandcomments/internal/server"
	"github.com/spf13/viper"
)

const (
	defaultPort            = "8080"
	InMemoryStorage string = "memory"
	PostgreStorage  string = "postgres"
)

func main(){
	if err := configs.InitConfig(); err != nil {
		log.Fatalf("error to open config: %v", err)
	}

	port := viper.GetString("port")
	if port == "" {
		port = defaultPort
	}

	var dataBase db.Database
	dbType := flag.String("storage-type", "", "Type of storage (memory or postgres)")
	flag.Parse()

	switch *dbType {
	case InMemoryStorage:
		dataBase = db.NewInMemoryDB()
	case PostgreStorage:
		db, err := db.NewPostgresDB(
			viper.GetString("postgres_host"),
			viper.GetInt("postgres_port"), 
			viper.GetString("postgres_user"), 
			viper.GetString("postgres_password"),
		)
		if err != nil {
			log.Fatalf("error to open postgresql: %v", err)
		}
		dataBase = db
	default:
		log.Fatalf("invalid storage type. Use --storage-type either memory or postgres.")
	}

	server.StartServer(port, dataBase)
}