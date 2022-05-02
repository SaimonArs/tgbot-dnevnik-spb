package main

import (
	"encoding/json"
	"log"
	"os"

	"main.go/clients/dnevnik2"
	tgClient "main.go/clients/telegram"
	event_consumer "main.go/consumer/event-consumer"
	"main.go/database/sqlite"
	"main.go/events/telegram"
)

func main() {
    cfg := mustConf()
    botTgClient := tgClient.New(cfg.TelegramHostClient, cfg.TelegramBotToken)
    dnevnikClient := dnevnik2.New(cfg.DnevnikHostClient, cfg.DnevnikBasePath)
    sqliteClient, err := sqlite.New(cfg.StoragePath)
    if err != nil {
        log.Fatalln("Ошибка при назначении базы данных")
    }

    eventProcessor := telegram.New(botTgClient, sqliteClient, *dnevnikClient)
    log.Println("service started")

    consumer := event_consumer.New(eventProcessor, eventProcessor, 100)

    if err := consumer.Start(); err != nil {
        log.Fatal("service is stoped", err)
    }

}
type Config struct {
    TelegramBotToken string
    TelegramHostClient string 
    DnevnikHostClient string
    DnevnikBasePath string
    StoragePath string
}

func mustConf() Config{
    file, _ := os.Open("config.json")
    decoder := json.NewDecoder(file)
    configuration := Config{}
    err := decoder.Decode(&configuration)
    if err != nil || configuration.TelegramBotToken == "" {
        log.Fatal("configuration", err) 
    }
    return configuration
}
