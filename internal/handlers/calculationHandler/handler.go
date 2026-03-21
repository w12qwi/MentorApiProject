package calculationHandler

import (
	"MentorApiProject/internal/adapter"
	"MentorApiProject/internal/domain/models"
	"MentorApiProject/internal/handlers/calculationHandler/dto"
	"MentorApiProject/internal/infrastructure/grpc"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"net/http"
)

type CalculationService interface {
	Calculate(ctx context.Context, calculation models.Calculation) (float64, error)
	GetCalculation(ctx context.Context, id string) (*models.Calculation, error)
	GetAllCalculations(ctx context.Context) ([]*models.Calculation, error)
	GetCalculationsWithFilters(ctx context.Context, filters models.CalculationsFilters) ([]*models.Calculation, error)
}

type Handler struct {
	calculationService CalculationService
}

func NewHandler(calculationService CalculationService) *Handler {
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
		slog.Error("Unable to encode response:", err)
	}
	return
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {

	req := dto.GetCalculationByIdRequest{
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

	id := uuid.MustParse(req.Id).String()

	resp, err := h.calculationService.GetCalculation(r.Context(), id)
	if err != nil {
		if errors.Is(err, grpc.CalcultionDoesNotExist) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			err = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		if err != nil {
			slog.Error("Unable to encode error response:", err)
		}
		return
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
	return
}

func (h *Handler) GetCalculations(w http.ResponseWriter, r *http.Request) {

	var req dto.GetCalculationstWithFiltersRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		if err == io.EOF {
			resp, err := h.calculationService.GetAllCalculations(r.Context())
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				err = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
				if err != nil {
					slog.Error("Unable to encode error response:", err)
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
			}

			response := adapter.DomainSliceToDto(resp)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				slog.Error("Unable to encode error response:", err)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: dto.InvalidRequestBodyError.Error()})
	}

	filters, err := adapter.DtoFiltersToDomain(req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
	}

	resp, err := h.calculationService.GetCalculationsWithFilters(r.Context(), filters)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
	}
	response := adapter.DomainSliceToDto(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /calculate", h.Calculate)
	mux.HandleFunc("GET /calculations/{id}", h.GetById)
	mux.HandleFunc("GET /calculations", h.GetCalculations)
}
