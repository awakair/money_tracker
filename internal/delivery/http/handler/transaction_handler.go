package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"money_tracker/internal/delivery/http/dto"
	httpinternal "money_tracker/internal/delivery/http/internal"

	"github.com/go-chi/chi/v5"
)

type TransactionHandler struct {
	transactionUseCase TransactionUseCase
}

func NewTransactionHandler(transactionUseCase TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		transactionUseCase: transactionUseCase,
	}
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Create a new transaction for the authenticated user
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param request body dto.CreateTransactionRequest true "Transaction data"
// @Success 201 {object} dto.TransactionResponse "Transaction created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request body"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 404 {object} dto.ErrorResponse "Category not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /users/{userId}/transactions [post]
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")

	var req dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpinternal.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	transaction, err := h.transactionUseCase.CreateTransaction(r.Context(), userID, req.ToUsecaseDTO())
	if err != nil {
		processError(w, err)
		return
	}

	resp := dto.NewTransactionResponseFromUsecaseDTO(transaction)
	httpinternal.RespondWithJSON(w, http.StatusCreated, resp)
}

// ListTransactions godoc
// @Summary List user transactions
// @Description Get a paginated list of transactions for the authenticated user
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param limit query int false "Maximum number of transactions to return"
// @Param offset query int false "Number of transactions to skip"
// @Success 200 {object} dto.TransactionListResponse "List of transactions"
// @Failure 400 {object} dto.ErrorResponse "Invalid query parameters"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /users/{userId}/transactions [get]
func (h *TransactionHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")

	var limit, offset *int

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit < 0 {
			httpinternal.RespondWithError(w, http.StatusBadRequest, "invalid limit parameter")
			return
		}
		limit = &parsedLimit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			httpinternal.RespondWithError(w, http.StatusBadRequest, "invalid offset parameter")
			return
		}
		offset = &parsedOffset
	}

	transactions, total, err := h.transactionUseCase.ListTransactions(r.Context(), userID, limit, offset)
	if err != nil {
		processError(w, err)
		return
	}

	resp := dto.TransactionListResponse{
		Transactions: make([]dto.TransactionResponse, len(transactions)),
		Total:        total,
		Limit:        limit,
		Offset:       offset,
	}

	for i, tx := range transactions {
		resp.Transactions[i] = dto.NewTransactionResponseFromUsecaseDTO(tx)
	}

	httpinternal.RespondWithJSON(w, http.StatusOK, resp)
}

// GetTransaction godoc
// @Summary Get a transaction by ID
// @Description Get detailed information about a specific transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param id path string true "Transaction ID"
// @Success 200 {object} dto.TransactionResponse "Transaction details"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 404 {object} dto.ErrorResponse "Transaction not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /users/{userId}/transactions/{id} [get]
func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	transactionID := chi.URLParam(r, "id")

	transaction, err := h.transactionUseCase.GetTransaction(r.Context(), userID, transactionID)
	if err != nil {
		processError(w, err)
		return
	}

	resp := dto.NewTransactionResponseFromUsecaseDTO(transaction)
	httpinternal.RespondWithJSON(w, http.StatusOK, resp)
}

// UpdateTransaction godoc
// @Summary Update a transaction
// @Description Update an existing transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param id path string true "Transaction ID"
// @Param request body dto.UpdateTransactionRequest true "Updated transaction data"
// @Success 200 {object} dto.TransactionResponse "Transaction updated successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request body"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 404 {object} dto.ErrorResponse "Transaction not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /users/{userId}/transactions/{id} [put]
func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	transactionID := chi.URLParam(r, "id")

	var req dto.UpdateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpinternal.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	transaction, err := h.transactionUseCase.UpdateTransaction(r.Context(), userID, transactionID, req.ToUsecaseDTO())
	if err != nil {
		processError(w, err)
		return
	}

	resp := dto.NewTransactionResponseFromUsecaseDTO(transaction)
	httpinternal.RespondWithJSON(w, http.StatusOK, resp)
}
