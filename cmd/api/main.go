// Filename: cmd/api/main.go

package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/kelseyaban/National-Inservice-Training-Database/internal/data"
	"github.com/kelseyaban/National-Inservice-Training-Database/internal/mailer"
	_ "github.com/lib/pq"
)

// configuration holds all the runtime configuration settings for the app.
// Private (non-exportable) to this package (lowercase "configuration")
type configuration struct {
	port    int
	env     string // Application environment
	version string // Version number of the API
	db      struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	cors struct {
		trustedOrigins []string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

// Hold dependencies shared across handlers,
// such as config and logger.
type application struct {
	config configuration
	logger *slog.Logger
	// quoteModel      data.QuoteModel
	userModel              data.UserModel
	courseModel            data.CourseModel
	mailer                 mailer.Mailer
	wg                     sync.WaitGroup
	tokenModel             data.TokenModel
	permissionModel        data.PermissionModel
	roleModel              data.RoleModel
	facilitatorRatingModel data.FacilitatorRatingModel
	sessionModel           data.SessionModel
}

// loadConfig reads configuration from command line flags
func loadConfig() configuration {
	var cfg configuration

	// Register CLI flags and bind them to cfg fields.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")
	flag.StringVar(&cfg.version, "version", "1.0.0", "Application version")

	// Read in the dsn
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://quotes:whyme@localhost/quotes", "PostgreSQL DSN")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate Limiter maximum requests per second")

	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 5, "Rate Limiter maximum burst")

	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	// Flags for SMTP
	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	// We have port 25, 465, 587, 2525. If 25 doesn't work choose another
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "213718fe792166", "SMTP username")
	// Use your Password value provided by Mailtrap
	flag.StringVar(&cfg.smtp.password, "smtp-password", "54177abe81857f", "SMTP password")

	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Training Datatbase <no-reply@trainingdatabase.nationalinservice.net>", "SMTP sender")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	// Allow us to access space-seperted origins.
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space seperated)",
		func(val string) error {
			cfg.cors.trustedOrigins = strings.Fields(val)
			return nil
		})

	flag.Parse()

	return cfg
}

// setupLogger configures the application logger based on environment
func setupLogger() *slog.Logger {
	var logger *slog.Logger

	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	return logger
}

func openDB(settings configuration) (*sql.DB, error) {
	// open a connection pool
	db, err := sql.Open("postgres", settings.db.dsn)
	if err != nil {
		return nil, err
	}

	// add the database config options  from the command line flags
	db.SetMaxOpenConns(settings.db.maxOpenConns)
	db.SetMaxIdleConns(settings.db.maxIdleConns)
	db.SetConnMaxIdleTime(settings.db.maxIdleTime)

	// set a context to ensure DB operations don't take too long
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test if the connection pool was created
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	// return the connection pool (sql.DB)
	return db, nil

}

func main() {

	// Initialize configuration
	cfg := loadConfig()
	// Initialize logger
	logger := setupLogger()

	// Call to openDB() sets up our connection pool
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// release the database resources before exiting
	defer db.Close()

	logger.Info("database connection pool established")
	expvar.NewString("version").Set(cfg.version)

	// the number of active goroutines
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	// the database connection pool metrics
	expvar.Publish("database", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	// the current Unix timestamp
	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	// Initialize application struc with dependencies
	app := &application{
		config: cfg,
		logger: logger,
		// quoteModel: data.QuoteModel{DB: db},
		userModel:   data.UserModel{DB: db},
		courseModel: data.CourseModel{DB: db},
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port,
			cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
		tokenModel:             data.TokenModel{DB: db},
		permissionModel:        data.PermissionModel{DB: db},
		roleModel:              data.RoleModel{DB: db},
		facilitatorRatingModel: data.FacilitatorRatingModel{DB: db},
		sessionModel:           data.SessionModel{DB: db},
	}

	// Run the application
	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}
