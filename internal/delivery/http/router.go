package http

import (
	"net/http"

	"money_tracker/internal/delivery/http/handler"
	"money_tracker/internal/delivery/http/middleware"
	"money_tracker/internal/pkg/logger"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "money_tracker/docs" // swagger docs
)

type Router struct {
	userHandler            *handler.UserHandler
	authHandler            *handler.AuthHandler
	transactionHandler     *handler.TransactionHandler
	categoryHandler        *handler.CategoryHandler
	transactionRuleHandler *handler.TransactionRuleHandler
	jwtValidator           middleware.JWTValidator
	logger                 logger.Logger
}

func NewRouter(
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	transactionHandler *handler.TransactionHandler,
	categoryHandler *handler.CategoryHandler,
	transactionRuleHandler *handler.TransactionRuleHandler,
	jwtValidator middleware.JWTValidator,
	logger logger.Logger,
) *Router {
	return &Router{
		userHandler:            userHandler,
		authHandler:            authHandler,
		transactionHandler:     transactionHandler,
		categoryHandler:        categoryHandler,
		transactionRuleHandler: transactionRuleHandler,
		jwtValidator:           jwtValidator,
		logger:                 logger,
	}
}

func (rt *Router) Setup() http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.LoggerMiddleware(rt.logger))
	r.Use(middleware.RecoveryMiddleware(rt.logger))

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Post("/users", rt.userHandler.CreateUser)
	r.Post("/login", rt.authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(rt.jwtValidator))

		r.Post("/users/{userId}/transactions", rt.transactionHandler.CreateTransaction)
		r.Get("/users/{userId}/transactions", rt.transactionHandler.ListTransactions)
		r.Get("/users/{userId}/transactions/{id}", rt.transactionHandler.GetTransaction)
		r.Put("/users/{userId}/transactions/{id}", rt.transactionHandler.UpdateTransaction)

		r.Post("/users/{userId}/categories", rt.categoryHandler.CreateCategory)
		r.Get("/users/{userId}/categories", rt.categoryHandler.ListCategories)
		r.Get("/users/{userId}/categories/{id}", rt.categoryHandler.GetCategory)
		r.Put("/users/{userId}/categories/{id}", rt.categoryHandler.UpdateCategory)

		r.Post("/users/{userId}/transaction-rules", rt.transactionRuleHandler.CreateTransactionRule)
		r.Get("/users/{userId}/transaction-rules", rt.transactionRuleHandler.ListTransactionRules)
		r.Get("/users/{userId}/transaction-rules/{id}", rt.transactionRuleHandler.GetTransactionRule)
		r.Put("/users/{userId}/transaction-rules/{id}", rt.transactionRuleHandler.UpdateTransactionRule)
		r.Delete("/users/{userId}/transaction-rules/{id}", rt.transactionRuleHandler.DeleteTransactionRule)
	})

	return r
}
