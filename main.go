package main

import (
	"context"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"

	core "github.com/tommzn/hdb-core"
)

func main() {

	ctx := context.Background()
	minion, bootstrapError := bootstrap(nil, ctx)
	exitOnError(bootstrapError)

	executionError := minion.Run(ctx)
	exitOnError(executionError)
}

// bootstrap loads config and creates a processor for SQS events.
func bootstrap(conf config.Config, ctx context.Context) (*core.Minion, error) {

	secretsManager := newSecretsManager()
	conf, err := loadConfig()
	if err != nil {
		return nil, err
	}
	logger := newLogger(conf, secretsManager, ctx)
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
func newLogger(conf config.Config, secretsMenager secrets.SecretsManager, ctx context.Context) log.Logger {
	logger := log.NewLoggerFromConfig(conf, secretsMenager)
	logger = log.WithNameSpace(logger, "hdb-message-agent")
	return log.WithK8sContext(logger)
}

func exitOnError(err error) {
	if err != nil {
		panic(err)
	}
}
