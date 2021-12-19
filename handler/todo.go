package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	tm, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{TODO: tm}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{TODOs: todos}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	tm, err := h.svc.UpdateTODO(ctx, int64(req.ID), req.Subject, req.Description)
	if err != nil {
		return nil, err
	}

	return &model.UpdateTODOResponse{TODO: tm}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	if err := h.svc.DeleteTODO(ctx, req.IDs); err != nil {
		return nil, err
	}
	return &model.DeleteTODOResponse{}, nil
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		size := r.URL.Query().Get("size")
		log.Println(size)
		var err error
		size64 := int64(5)
		if size != "" {
			size64, err = strconv.ParseInt(size, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		prevId := r.URL.Query().Get("prev_id")
		prevId64 := int64(0)
		if prevId != "" {
			prevId64, err = strconv.ParseInt(prevId, 10, 64)
			if err != nil {
				log.Println(err)
				return
			}
		}
		request := &model.ReadTODORequest{Size: size64, PrevID: prevId64}

		response, err := h.Read(r.Context(), request)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

		je := json.NewEncoder(w)

		if err := je.Encode(response); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}

	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var request model.CreateTODORequest
		err := decoder.Decode(&request)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if request.Subject == "" {
			log.Println("Subject not found")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response, err := h.Create(r.Context(), &request)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

		je := json.NewEncoder(w)

		if err := je.Encode(response); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}

	case http.MethodPut:
		decoder := json.NewDecoder(r.Body)
		var request model.UpdateTODORequest
		err := decoder.Decode(&request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}
		if request.ID == 0 {
			log.Println("ID not found")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if request.Subject == "" {
			log.Println("Subject not found")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response, err := h.Update(r.Context(), &request)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

		je := json.NewEncoder(w)

		if err := je.Encode(response); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		decoder := json.NewDecoder(r.Body)
		var request model.DeleteTODORequest
		err := decoder.Decode(&request)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response, err := h.Delete(r.Context(), &request)

		switch err {
		case nil:
			break
		case model.ErrNotFound{}:
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

		je := json.NewEncoder(w)

		if err := je.Encode(response); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
