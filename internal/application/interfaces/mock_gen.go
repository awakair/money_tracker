package interfaces

//go:generate mockgen -destination=mocks/mock_database.go -package=mocks money_tracker/internal/application/interfaces Database
//go:generate mockgen -destination=mocks/mock_id_generator.go -package=mocks money_tracker/internal/application/interfaces IDGenerator
//go:generate mockgen -destination=mocks/mock_password_hasher.go -package=mocks money_tracker/internal/application/interfaces PasswordHasher
//go:generate mockgen -destination=mocks/mock_transaction_repository.go -package=mocks money_tracker/internal/application/interfaces TransactionRepositoryRW
//go:generate mockgen -destination=mocks/mock_transaction_repository_r.go -package=mocks money_tracker/internal/application/interfaces TransactionRepositoryR
//go:generate mockgen -destination=mocks/mock_transaction_rule_repository.go -package=mocks money_tracker/internal/application/interfaces TransactionRuleRepositoryRW
//go:generate mockgen -destination=mocks/mock_transaction_rule_repository_r.go -package=mocks money_tracker/internal/application/interfaces TransactionRuleRepositoryR
//go:generate mockgen -destination=mocks/mock_category_repository.go -package=mocks money_tracker/internal/application/interfaces CategoryRepositoryRW
//go:generate mockgen -destination=mocks/mock_category_repository_r.go -package=mocks money_tracker/internal/application/interfaces CategoryRepositoryR
//go:generate mockgen -destination=mocks/mock_user_repository.go -package=mocks money_tracker/internal/application/interfaces UserRepositoryRW
//go:generate mockgen -destination=mocks/mock_user_repository_r.go -package=mocks money_tracker/internal/application/interfaces UserRepositoryR
