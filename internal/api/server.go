package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/golang-migrate/migrate/v4"
	postgresMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	"github.com/lastbyte32/gofemart/internal/api/handler"
	"github.com/lastbyte32/gofemart/internal/config"
	"github.com/lastbyte32/gofemart/internal/jwt"
	"github.com/lastbyte32/gofemart/internal/service"
	"github.com/lastbyte32/gofemart/internal/storage/postgres"
)

type Configurator interface {
	GetApiHost() string
	GetDSN() string
	GetSigningKey() string
}

type Server struct {
	logger *zerolog.Logger
	http   *http.Server
	db     *sqlx.DB
	cfg    Configurator
}

func New() (*Server, error) {

	c, err := config.New()
	if err != nil {
		return nil, err
	}

	zero := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		//Caller().
		//Int("pid", os.Getpid()).
		//Str("go_version", "dddfdfff").
		Logger()

	server := &http.Server{
		Addr:              c.GetApiHost(),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	return &Server{
		logger: &zero,
		http:   server,
		cfg:    c,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	db, err := s.configureDataBase(ctx, s.cfg.GetDSN())
	if err != nil {
		return err
	}
	s.db = db
	if err := s.Migrate(); err != nil {
		return err
	}
	s.logger.Info().Msg("database migrate complete")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		s.logger.Info().Msg("start shutdown watcher")
		<-ctx.Done()
		s.logger.Info().Msg("received signal, stopping application")
		s.stop()
		s.logger.Info().Msg("app terminated")
		wg.Done()
	}()

	router, errRouter := s.configureRoutes()
	if errRouter != nil {
		return errRouter
	}
	s.http.Handler = router

	s.logger.Info().Msgf("starting the api-server on %s", s.http.Addr)
	if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	wg.Wait()
	return nil
}

func (s *Server) stop() {
	s.logger.Info().Msg("shutdown initiated")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.http.Shutdown(ctx); err != nil {
		s.logger.Error().Msgf("shutdown api-server failed: %v", err)
	}
	s.logger.Info().Msg("api-server stopped successfully")

	err := s.db.Close()
	if err != nil {
		s.logger.Error().Msgf("shutdown database failed: %v", err)
	}
	s.logger.Info().Msg("database stopped successfully")
	s.logger.Info().Msg("shutdown completed")
}

func (s *Server) configureDataBase(ctx context.Context, dsn string) (*sqlx.DB, error) {
	s.logger.Info().Str("DSN", dsn).Msg("configure database")
	if _, err := pgx.ParseConfig(dsn); err != nil {
		return nil, fmt.Errorf("error dsn config: %w", err)
	}

	db, err := sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	s.logger.Info().Msg("database connection successfully")
	return db, nil
}

func (s *Server) Migrate() error {
	dbInstance, err := postgresMigrate.WithInstance(s.db.DB, &postgresMigrate.Config{})
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://migrations", "pgx", dbInstance)
	if err != nil {
		return err
	}
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (s *Server) configureRoutes() (chi.Router, error) {

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Heartbeat("/health"))

	tokenManager, _ := jwt.NewManager("1234")
	userStorage := postgres.NewUserStore(s.db)

	userService := service.NewUserService(userStorage, tokenManager)

	userHandler := handler.NewUserHandler(userService)
	userHandler.Routes(router)

	s.logger.Info().Msg("route list")
	chi.Walk(router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		var middlewareNames []string
		for _, mid := range middlewares {
			parts := strings.Split((runtime.FuncForPC(reflect.ValueOf(mid).Pointer()).Name()), "/")
			middlewareNames = append(middlewareNames, parts[len(parts)-1])
		}
		middlewareNamesStr := strings.Join(middlewareNames, ", ")
		s.logger.Info().Msgf("[%s] %s (%s)", method, route, middlewareNamesStr)
		return nil
	})

	return router, nil
}
