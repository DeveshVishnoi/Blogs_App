package main

import (
	"fmt"
	"social/internal/db"
	"social/internal/environment"
	"social/internal/mailer"
	"social/internal/store"
	"time"

	"go.uber.org/zap"
)

const version = "0.0.1"

var logger *zap.SugaredLogger

func init() {
	logger = zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	logger.Infow("Logger Initialize", "function", "init")
}

// Calling main functin for start the main thread.
func main() {

	logger.Infow("Inside the Main thread", "function", "main")
	// Set the environmant data for the running the application.
	config := config{
		addr:        environment.GetString("ADDR", ":9000"),
		frontendURL: environment.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{

			// In my Server postgreSQL running on the 5433 port number.
			// we can also use this
			// "postgres://postgres:admin@localhost:5433/social?sslmode=disable"
			addr:         environment.GetString("DB_ADDR", "host=localhost port=5433 user=postgres password=admin dbname=social sslmode=disable"),
			maxOpenConns: environment.GetIntegerValue("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: environment.GetIntegerValue("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  environment.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: environment.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days
			fromEmail: environment.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: environment.GetString("SENDGRID_API_KEY", ""),
			},
			mailTrap: mailTrapConfig{
				apiKey: environment.GetString("MAILTRAP_API_KEY", ""),
			},
		},
	}

	// Initilaize the database.
	db, err := db.New(config.db.addr, config.db.maxOpenConns, config.db.maxIdleConns, config.db.maxIdleTime)
	if err != nil {
		logger.Panic(err)
	}

	defer db.Close()

	fmt.Println("Database connection pool established ")

	mailtrap, err := mailer.NewMailTrapClient(config.mail.mailTrap.apiKey, config.mail.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	// Make Databases instanses or add as a layer of the database.
	storage := store.NewStorage(db)

	// Make the application instances.
	app := application{
		config: config,
		store:  storage,
		logger: logger,
		mailer: mailtrap,
	}

	app.logger.Infow("server has start", "addr", app.config.addr, "env", app.config.env)

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
