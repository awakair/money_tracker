package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"money_tracker/internal/delivery/http/dto"
	httpinternal "money_tracker/internal/delivery/http/internal"

	"github.com/go-chi/chi/v5"
)

type TransactionRuleHandler struct {
	transactionRuleUseCase TransactionRuleUseCase
}

func NewTransactionRuleHandler(transactionRuleUseCase TransactionRuleUseCase) *TransactionRuleHandler {
	return &TransactionRuleHandler{
		transactionRuleUseCase: transactionRuleUseCase,
	}
}

// CreateTransactionRule godoc
// @Summary Create a new transaction rule
// @Description Create a new recurring transaction rule for the authenticated user
// @Tags transaction-rules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param request body dto.CreateTransactionRuleRequest true "Transaction rule data"
// @Success 201 {object} dto.TransactionRuleResponse "Transaction rule created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request body or cron expression"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 404 {object} dto.ErrorResponse "Category not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /users/{userId}/transaction-rules [post]
func (h *TransactionRuleHandler) CreateTransactionRule(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")

	var req dto.CreateTransactionRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpinternal.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	rule, err := h.transactionRuleUseCase.CreateTransactionRule(r.Context(), userID, req.ToUsecaseDTO())
	if err != nil {
		processError(w, err)
		return
	}

	resp := dto.NewTransactionRuleResponseFromUsecaseDTO(rule)
	httpinternal.RespondWithJSON(w, http.StatusCreated, resp)
}

// ListTransactionRules godoc
// @Summary List user transaction rules
// @Description Get a paginated list of transaction rules for the authenticated user
// @Tags transaction-rules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param limit query int false "Maximum number of rules to return"
// @Param offset query int false "Number of rules to skip"
// @Success 200 {object} dto.TransactionRuleListResponse "List of transaction rules"
// @Failure 400 {object} dto.ErrorResponse "Invalid query parameters"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /users/{userId}/transaction-rules [get]
func (h *TransactionRuleHandler) ListTransactionRules(w http.ResponseWriter, r *http.Request) {
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

	rules, total, err := h.transactionRuleUseCase.ListTransactionRules(r.Context(), userID, limit, offset)
	if err != nil {
		processError(w, err)
		return
	}

	resp := dto.TransactionRuleListResponse{
		Rules:  make([]dto.TransactionRuleResponse, len(rules)),
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	for i, rule := range rules {
		resp.Rules[i] = dto.NewTransactionRuleResponseFromUsecaseDTO(rule)
	}

	httpinternal.RespondWithJSON(w, http.StatusOK, resp)
}

// GetTransactionRule godoc
// @Summary Get a transaction rule by ID
// @Description Get detailed information about a specific transaction rule
// @Tags transaction-rules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param id path string true "Transaction rule ID"
// @Success 200 {object} dto.TransactionRuleResponse "Transaction rule details"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 404 {object} dto.ErrorResponse "Transaction rule not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /users/{userId}/transaction-rules/{id} [get]
func (h *TransactionRuleHandler) GetTransactionRule(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	ruleID := chi.URLParam(r, "id")

	rule, err := h.transactionRuleUseCase.GetTransactionRule(r.Context(), userID, ruleID)
	if err != nil {
		processError(w, err)
		return
	}

	resp := dto.NewTransactionRuleResponseFromUsecaseDTO(rule)
	httpinternal.RespondWithJSON(w, http.StatusOK, resp)
}

// UpdateTransactionRule godoc
// @Summary Update a transaction rule
// @Description Update an existing transaction rule
// @Tags transaction-rules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param id path string true "Transaction rule ID"
// @Param request body dto.UpdateTransactionRuleRequest true "Updated transaction rule data"
// @Success 200 {object} dto.TransactionRuleResponse "Transaction rule updated successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request body or cron expression"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 404 {object} dto.ErrorResponse "Transaction rule not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /users/{userId}/transaction-rules/{id} [put]
func (h *TransactionRuleHandler) UpdateTransactionRule(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	ruleID := chi.URLParam(r, "id")

	var req dto.UpdateTransactionRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpinternal.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	rule, err := h.transactionRuleUseCase.UpdateTransactionRule(r.Context(), userID, ruleID, req.ToUsecaseDTO())
	if err != nil {
		processError(w, err)
		return
	}

	resp := dto.NewTransactionRuleResponseFromUsecaseDTO(rule)
	httpinternal.RespondWithJSON(w, http.StatusOK, resp)
}

// DeleteTransactionRule godoc
// @Summary Delete a transaction rule
// @Description Delete an existing transaction rule
// @Tags transaction-rules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param id path string true "Transaction rule ID"
// @Success 204 "Transaction rule deleted successfully"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 404 {object} dto.ErrorResponse "Transaction rule not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /users/{userId}/transaction-rules/{id} [delete]
func (h *TransactionRuleHandler) DeleteTransactionRule(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	ruleID := chi.URLParam(r, "id")

	err := h.transactionRuleUseCase.DeleteTransactionRule(r.Context(), userID, ruleID)
	if err != nil {
		processError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
