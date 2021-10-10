package handler

import (
	"ktrain/cmd/api/user-api/mapper"
	"ktrain/cmd/repository"
	"ktrain/pkg/httputil"
	"net/http"
)

type activityLogHandler struct {
	activityLogRepository repository.ActivityLogRepository
}

func NewActivityLogHandler(activityLogRepository repository.ActivityLogRepository) *activityLogHandler {
	return &activityLogHandler{
		activityLogRepository: activityLogRepository,
	}
}
func (h *activityLogHandler) GetActivity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	action, err := h.activityLogRepository.GetAllLogAction(r.Context(), ctx.Value("userID").(int64))
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting list action")
		return
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToActionResponse(action))
}
