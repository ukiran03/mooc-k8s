package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	model  *PingModel
	logger *slog.Logger
}

func main() {
	pingPort := os.Getenv("PING_PORT")
	if pingPort == "" {
		fmt.Println("env PING_PORT was unset\nUsing Port 3001 as pingPort")
		pingPort = "3001"
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dbDSN, err := getDbDSN()
	if dbDSN == "" || err != nil {
		log.Fatal("postgres DB DSN is wrong or missing")
	}

	db, err := openDB(dbDSN)
	if err != nil {
		logger.Error(
			"Failed to connect to database at startup (will retry via health check)",
			"error",
			err,
		)
	} else {
		defer db.Close()
		if err := initPingsTable(db); err != nil {
			logger.Error("Failed to initialize pings table", "error", err)
		}
	}

	app := &application{
		model:  &PingModel{DB: db, TTL: 3 * time.Second},
		logger: logger,
	}

	fmt.Printf("Starting Ping server on port %s...\n", pingPort)
	if err := http.ListenAndServe(":"+pingPort, app.routes()); err != nil {
		log.Fatalf("Server failed: %s\n", err)
	}
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", app.homeHandler)
	mux.HandleFunc("/healthz", app.healthcheckHandler)
	mux.HandleFunc("/pings", app.getPings)

	// Accept both /pingpong and /pingpong/ for POST requests
	mux.HandleFunc("POST /pingpong", app.pingHandler)
	mux.HandleFunc("POST /pingpong/", app.pingHandler)

	return mux
}

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "visit /pingpong")
}

// ReadinessProbe for the Ping-pong application. It should be ready when it has
// a connection to the database.
func (app *application) healthcheckHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if err := app.model.PgxPing(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR: DB not ready"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (app *application) getPings(w http.ResponseWriter, r *http.Request) {
	count, err := app.model.Get(r.Context())
	if err != nil {
		app.logger.Error("failed to get ping count", "error", err)
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "ping %d", count)
}

func (app *application) pingHandler(w http.ResponseWriter, r *http.Request) {
	count, err := app.model.Increment(r.Context())
	if err != nil {
		app.logger.Error("failed to increment ping", "error", err)
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "ping %d", count)
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 2
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}
	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

type PingModel struct {
	DB  *pgxpool.Pool
	TTL time.Duration
}

func (m *PingModel) PgxPing(ctx context.Context) error {
	if m.DB == nil {
		return fmt.Errorf("DB pool is not initialized")
	}
	if err := m.DB.Ping(ctx); err != nil {
		return fmt.Errorf("DB is not ready: %w", err)
	}
	return nil
}

func (m *PingModel) Increment(ctx context.Context) (int, error) {
	if m.DB == nil {
		return 0, fmt.Errorf("database is not available")
	}
	stmt := `UPDATE pings SET val = val + 1 RETURNING val`

	ctx, cancel := context.WithTimeout(ctx, m.TTL)
	defer cancel()

	var count int
	err := m.DB.QueryRow(ctx, stmt).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *PingModel) Get(ctx context.Context) (int, error) {
	if m.DB == nil {
		return 0, fmt.Errorf("database is not available")
	}
	stmt := `SELECT val FROM pings LIMIT 1`

	ctx, cancel := context.WithTimeout(ctx, m.TTL)
	defer cancel()

	var count int
	err := m.DB.QueryRow(ctx, stmt).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ---

func getDbDSN() (string, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	passwd := os.Getenv("DB_PASSWD")
	dbName := os.Getenv("DB_NAME")

	if host == "" {
		host = "postgres-svc"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" || passwd == "" || dbName == "" {
		return "", fmt.Errorf(
			`missing required environment variables (DB_USER, DB_PASSWD, or DB_NAME)`,
		)
	}
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, passwd, dbName,
	)
	return dsn, nil
}

func initPingsTable(db *pgxpool.Pool) error {
	ctx := context.Background()

	stmt := `CREATE TABLE IF NOT EXISTS pings (val INTEGER)`
	if _, err := db.Exec(ctx, stmt); err != nil {
		return err
	}

	var count int
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM pings").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = db.Exec(ctx, "INSERT INTO pings (val) VALUES (0)")
		if err != nil {
			return err
		}
		log.Println("Seeded database with initial row.")
	}
	return nil
}
