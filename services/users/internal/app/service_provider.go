package app

import (
	"context"
	"log"

	"screw-it-desk-backend/services/users/internal/client/db"
	"screw-it-desk-backend/services/users/internal/client/db/pg"
	"screw-it-desk-backend/services/users/internal/client/db/transaction"
	"screw-it-desk-backend/services/users/internal/closer"
	"screw-it-desk-backend/services/users/internal/config"
	"screw-it-desk-backend/services/users/internal/config/env"
	"screw-it-desk-backend/services/users/internal/service"

	"screw-it-desk-backend/services/users/internal/repository"

	userRepository "screw-it-desk-backend/services/users/internal/repository/users"

	noteService "screw-it-desk-backend/services/users/internal/service/user"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	httpConfig config.HTTPConfig

	dbClient       db.Client
	txManager      db.TxManager
	noteRepository repository.UserRepository

	noteService service.UserService
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) NoteRepository(ctx context.Context) repository.UserRepository {
	if s.noteRepository == nil {
		s.noteRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.noteRepository
}

func (s *serviceProvider) NoteService(ctx context.Context) service.UserService {
	if s.noteService == nil {
		s.noteService = noteService.NewService(
			s.NoteRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.noteService
}
