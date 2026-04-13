package all

import (
	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"
	authmod "github.com/masudur-rahman/expense-tracker-bot/modules/auth"
	authrepo "github.com/masudur-rahman/expense-tracker-bot/repos/auth"
	"github.com/masudur-rahman/expense-tracker-bot/repos/budgets"
	"github.com/masudur-rahman/expense-tracker-bot/repos/event"
	"github.com/masudur-rahman/expense-tracker-bot/repos/transaction"
	"github.com/masudur-rahman/expense-tracker-bot/repos/user"
	"github.com/masudur-rahman/expense-tracker-bot/repos/wallets"
	"github.com/masudur-rahman/expense-tracker-bot/services"
	authsvc "github.com/masudur-rahman/expense-tracker-bot/services/auth"
	budgetsvc "github.com/masudur-rahman/expense-tracker-bot/services/budgets"
	eventsvc "github.com/masudur-rahman/expense-tracker-bot/services/event"
	txnsvc "github.com/masudur-rahman/expense-tracker-bot/services/transaction"
	usersvc "github.com/masudur-rahman/expense-tracker-bot/services/user"
	walletsvc "github.com/masudur-rahman/expense-tracker-bot/services/wallets"

	"github.com/masudur-rahman/styx"
)

type Services struct {
	User    services.ProfileService
	Wallet  services.WalletService
	Contact services.ContactService
	Txn     services.TransactionService
	Event   services.EventService
	Budget  services.BudgetService
	Auth    services.AuthService
}

var svc *Services

func GetServices() *Services {
	return svc
}

func InitiateSQLServices(uow styx.UnitOfWork, logger logr.Logger) {
	userRepo := user.NewSQLUserRepository(uow.SQL, logger)
	walletRepo := wallets.NewSQLWalletRepository(uow.SQL, logger)
	contactRepo := user.NewSQLContactRepository(uow.SQL, logger)
	txnRepo := transaction.NewSQLTransactionRepository(uow.SQL, logger)
	eventRepo := event.NewSQLEventRepository(uow.SQL, logger)
	budgetRepo := budgets.NewSQLBudgetRepository(uow.SQL, logger)

	userSvc := usersvc.NewProfileService(userRepo)
	walletSvc := walletsvc.NewWalletService(walletRepo)
	contactSvc := usersvc.NewContactService(contactRepo)
	txnSvc := txnsvc.NewTxnService(uow, walletRepo, contactRepo, txnRepo, eventRepo)
	eventSvc := eventsvc.NewEventService(eventRepo)
	budgetSvc := budgetsvc.NewBudgetService(budgetRepo, txnRepo)

	svc = &Services{
		User:    userSvc,
		Wallet:  walletSvc,
		Contact: contactSvc,
		Txn:     txnSvc,
		Event:   eventSvc,
		Budget:  budgetSvc,
	}
}

// InitiateWebServices wires the auth service when the web dashboard is enabled.
func InitiateWebServices(
	messenger authmod.Messenger,
	jwtSecret, refreshSecret, botUsername string,
	uow styx.UnitOfWork,
	logger logr.Logger,
) {
	userRepo := user.NewSQLUserRepository(uow.SQL, logger)
	ar := authrepo.NewSQLAuthRepository(uow.SQL, logger)
	svc.Auth = authsvc.NewAuthService(
		userRepo, ar, messenger,
		jwtSecret, refreshSecret, botUsername,
		logger,
	)
}
