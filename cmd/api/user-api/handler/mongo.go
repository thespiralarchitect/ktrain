package handler

import (
	"ktrain/cmd/api/user-api/mapper"
	"ktrain/cmd/repository"
	"ktrain/pkg/httputil"
	"net/http"
)

type mongoHandler struct {
	mongoRepository repository.MongoRepository
}

func NewMongoHandler(mongoRepository repository.MongoRepository) *mongoHandler {
	return &mongoHandler{
		mongoRepository: mongoRepository,
	}
}
func (h *mongoHandler) GetAction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	action, err := h.mongoRepository.GetAllLogAction(r.Context(), ctx.Value("userID").(int64))
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting list action")
		return
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToActionResponse(action))
}
