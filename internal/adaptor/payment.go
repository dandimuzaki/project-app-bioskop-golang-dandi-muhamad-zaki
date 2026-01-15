package adaptor

import (
	"encoding/json"
	"net/http"

	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/internal/usecase"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type PaymentHandler struct {
	Usecase usecase.Usecase
	Logger *zap.Logger
	Config utils.Configuration
}

func NewPaymentHandler(uc usecase.Usecase, log *zap.Logger, config utils.Configuration) PaymentHandler {
	return PaymentHandler{
		Usecase: uc,
		Logger: log,
		Config: config,
	}
}

func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.PaymentRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Error decode request body to dto payment request: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Execute create payment
	result, err := h.Usecase.PaymentUsecase.Create(req)
	if err != nil {
		h.Logger.Error("Error handling create payment: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "create payment failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "create payment success", result)
}

func (h *PaymentHandler) Callback(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdatePayment

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Error decode request body to dto payment request: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Execute update payment
	err = h.Usecase.PaymentUsecase.Update(req)
	if err != nil {
		h.Logger.Error("Error handling update payment: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "update payment failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "update payment success", nil)
}