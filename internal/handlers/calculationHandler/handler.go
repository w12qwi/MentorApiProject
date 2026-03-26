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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"io"
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
	tracer             trace.Tracer
}

func NewHandler(calculationService CalculationService, tracer trace.Tracer) *Handler {
	return &Handler{calculationService: calculationService, tracer: tracer}
}

func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {

	ctx, span := h.tracer.Start(r.Context(), "handler.calculate")
	defer span.End()
	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.route", r.URL.Path),
	)

	var req dto.CalculateRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: dto.InvalidRequestBodyError.Error()})
		span.SetAttributes(
			attribute.Int("http.status_code", http.StatusBadRequest),
			attribute.Bool("request.decode_failed", true))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	err = req.Validate()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		span.SetAttributes(
			attribute.Int("http.status_code", http.StatusUnprocessableEntity),
			attribute.Bool("validation.failed", true))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	calculation := models.Calculation{
		NumA: req.NumA,
		NumB: req.NumB,
		Sign: req.Sign,
	}

	resp, err := h.calculationService.Calculate(ctx, calculation)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		span.SetAttributes(attribute.Int("http.status_code", http.StatusInternalServerError))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.CalculateResponse{Result: resp})
	span.SetAttributes(attribute.Int("http.status_code", http.StatusOK))
	return
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {

	ctx, span := h.tracer.Start(r.Context(), "handler.getCalculationById")
	defer span.End()
	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.route", r.URL.Path),
	)

	req := dto.GetCalculationByIdRequest{
		Id: r.PathValue("id"),
	}

	err := req.Validate()
	if err != nil {
		span.SetAttributes(
			attribute.Bool("validation.failed", true),
			attribute.Int("http.status_code", http.StatusUnprocessableEntity))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	id := uuid.MustParse(req.Id).String()

	resp, err := h.calculationService.GetCalculation(ctx, id)
	if err != nil {
		if errors.Is(err, grpc.CalcultionDoesNotExist) {
			span.SetAttributes(
				attribute.Bool("calculation.not_found", true),
				attribute.Int("http.status_code", http.StatusNotFound))
			span.SetStatus(codes.Error, err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
			return
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	response := adapter.DomainToDto(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	return
}

func (h *Handler) GetCalculations(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "handler.getCalculations")
	defer span.End()
	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.route", r.URL.Path),
	)

	var req dto.GetCalculationstWithFiltersRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		if err == io.EOF {
			resp, err := h.calculationService.GetAllCalculations(ctx)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				span.SetAttributes(attribute.Int("http.status_code", http.StatusInternalServerError))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
				return
			}
			response := adapter.DomainSliceToDto(resp)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: dto.InvalidRequestBodyError.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		span.SetAttributes(
			attribute.Bool("validation.failed", true),
			attribute.Int("http.status_code", http.StatusUnprocessableEntity))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	filters, err := adapter.DtoFiltersToDomain(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.Int("http.status_code", http.StatusBadRequest))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	var filterAttrs []attribute.KeyValue
	switch {
	case filters.Sign != nil:
		filterAttrs = append(filterAttrs, attribute.String("Applied filter: sign", *filters.Sign))
	case filters.Date != nil:
		filterAttrs = append(filterAttrs, attribute.String("Applied filter: date", filters.Date.String()))
	case filters.DateFrom != nil:
		filterAttrs = append(filterAttrs, attribute.String("Applied filter: dateFrom", filters.DateFrom.String()))
	case filters.DateTo != nil:
		filterAttrs = append(filterAttrs, attribute.String("Applied filter: dateTo", filters.DateTo.String()))
	}
	span.SetAttributes(filterAttrs...)

	resp, err := h.calculationService.GetCalculationsWithFilters(ctx, filters)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.Int("http.status_code", http.StatusInternalServerError))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	response := adapter.DomainSliceToDto(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /calculate", h.Calculate)
	mux.HandleFunc("GET /calculations/{id}", h.GetById)
	mux.HandleFunc("GET /calculations", h.GetCalculations)
}
