package main

import (
	"context"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"

	core "github.com/tommzn/hdb-core"
)

func main() {

	minion, bootstrapError := bootstrap(nil)
	exitOnError(bootstrapError)

	executionError := minion.Run(context.Background())
	exitOnError(executionError)
}

// bootstrap loads config and creates a processor for SQS events.
func bootstrap(conf config.Config) (*core.Minion, error) {

	conf, err := loadConfig()
	if err != nil {
		return nil, err
	}

	secretsManager := newSecretsManager()
	logger := newLogger(conf, secretsManager)
	agent, err := newAgent(conf, logger)
	return core.NewMinion(agent), err
}

// loadConfig from config file.
func loadConfig() (config.Config, error) {

	configSource, err := config.NewS3ConfigSourceFromEnv()
	if err != nil {
		return nil, err
	}

	conf, err := configSource.Load()
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// newSecretsManager retruns a new secrets manager from passed config.
func newSecretsManager() secrets.SecretsManager {
	secretsManager := secrets.NewDockerecretsManager("/run/secrets/token")
	secrets.ExportToEnvironment([]string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"}, secretsManager)
	return secretsManager
}

// newLogger creates a new logger from  passed config.
func newLogger(conf config.Config, secretsMenager secrets.SecretsManager) log.Logger {
	return log.NewLoggerFromConfig(conf, secretsMenager)
}

func exitOnError(err error) {
	if err != nil {
		panic(err)
	}
}
