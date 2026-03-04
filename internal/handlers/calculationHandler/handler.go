package calculationHandler

import (
	"MentorApiProject/internal/domain/models"
	"MentorApiProject/internal/handlers/calculationHandler/dto"
	"MentorApiProject/internal/service/calculationService"
	"encoding/json"
	"log/slog"
	"net/http"
)

type Handler struct {
	calculationService calculationService.Service
}

func NewHandler(calculationService calculationService.Service) *Handler {
	return &Handler{calculationService: calculationService}
}

func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {

	var req dto.CalculateRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: dto.InvalidRequestBodyError.Error()})
		if err != nil {
			slog.Error("Unable to encode error response:", err)
			return
		}
		return
	}

	err = req.Validate()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		err = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		if err != nil {
			slog.Error("Unable to encode error response:", err)
		}
		return
	}

	calculation := models.Calculation{
		NumA: req.NumA,
		NumB: req.NumB,
		Sign: req.Sign,
	}

	resp, err := h.calculationService.Calculate(r.Context(), calculation)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		if err != nil {
			slog.Error("Unable to encode error response:", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dto.CalculateResponse{Result: resp})
	if err != nil {
		slog.Error("Unable to encode error response:", err)
	}

}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {

	req := dto.GetByIdRequest{
		Id: r.PathValue("id"),
	}

	err := req.Validate()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		err = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		if err != nil {
			slog.Error("Unable to encode error response:", err)
		}
		return
	}

	resp, err := h.calculationService.GetById(r.Context(), req.UUID())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		if err != nil {
			slog.Error("Unable to encode error response:", err)
		}
	}

	response := dto.CalculationResponse{
		Id:        resp.Id.String(),
		NumA:      resp.NumA,
		NumB:      resp.NumB,
		Sign:      resp.Sign,
		Result:    resp.Result,
		CreatedAt: resp.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Unable to encode error response:", err)
	}

}

func (h *Handler) GetAllCalculations(w http.ResponseWriter, r *http.Request) {

	resp, err := h.calculationService.GetAllCalculations(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		if err != nil {
			slog.Error("Unable to encode error response:", err)
		}
	}

	response := []dto.CalculationResponse{}
	for _, calculation := range resp {
		response = append(response, dto.CalculationResponse{
			Id:        calculation.Id.String(),
			NumA:      calculation.NumA,
			NumB:      calculation.NumB,
			Sign:      calculation.Sign,
			Result:    calculation.Result,
			CreatedAt: calculation.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Unable to encode error response:", err)
	}

}

func (h *Handler) GetCalculationsByDate(w http.ResponseWriter, r *http.Request) {

	var req dto.GetCalculationsByDateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: dto.InvalidRequestBodyError.Error()})
		if err != nil {
			slog.Error("Unable to encode error response:", err)
			return
		}

	}

	err = req.Validate()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		err = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		if err != nil {
			slog.Error("Unable to encode error response:", err)
		}
	}

	resp, err := h.calculationService.GetCalculationsByDate(r.Context(), req.UTCDate())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		if err != nil {
			slog.Error("Unable to encode error response:", err)
			return
		}

	}

	response := []dto.CalculationResponse{}
	for _, calculation := range resp {
		response = append(response, dto.CalculationResponse{
			Id:        calculation.Id.String(),
			NumA:      calculation.NumA,
			NumB:      calculation.NumB,
			Sign:      calculation.Sign,
			Result:    calculation.Result,
			CreatedAt: calculation.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Unable to encode error response:", err)
	}
}

func (h *Handler) GetCalculationsByDateRange(w http.ResponseWriter, r *http.Request) {

	req := dto.GetCalculationsByDateRangeRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		err = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: dto.InvalidRequestBodyError.Error()})
		if err != nil {
			slog.Error("Unable to encode error response:", err)
			return
		}
	}
	err = req.Validate()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		err = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: dto.InvalidDateFormatError.Error()})
		if err != nil {
			slog.Error("Unable to encode error response:", err)
		}
	}

	resp, err := h.calculationService.GetCalculationsByDateRange(r.Context(), req.UTCDateFrom(), req.UTCDateTo())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		if err != nil {
			slog.Error("Unable to encode error response:", err)
			return
		}
	}

	response := []dto.CalculationResponse{}
	for _, calculation := range resp {
		response = append(response, dto.CalculationResponse{
			Id:        calculation.Id.String(),
			NumA:      calculation.NumA,
			NumB:      calculation.NumB,
			Sign:      calculation.Sign,
			Result:    calculation.Result,
			CreatedAt: calculation.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Unable to encode error response:", err)
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /calculate", h.Calculate)
	mux.HandleFunc("GET /сalculations/{id}", h.GetById)
	mux.HandleFunc("GET /calculations", h.GetAllCalculations)
	mux.HandleFunc("GET /calculations/by-date", h.GetCalculationsByDate)
	mux.HandleFunc("GET /calculations/by-date-range", h.GetCalculationsByDateRange)

}
