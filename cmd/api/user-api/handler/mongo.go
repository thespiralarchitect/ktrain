package handler

import (
	"ktrain/cmd/api/user-api/dto"
	"ktrain/cmd/api/user-api/mapper"
	"ktrain/pkg/httputil"
	"ktrain/pkg/logger"
	"ktrain/proto/pb"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type activityLogHandler struct {
	activityLogClient pb.ActivityLogDMSServiceClient
}

func NewActivityLogHandler(activityLogClient pb.ActivityLogDMSServiceClient) *activityLogHandler {
	return &activityLogHandler{
		activityLogClient: activityLogClient,
	}
}
func (h *activityLogHandler) GetActivity(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	preq := &pb.GetLogActionRequest{
		Id: int64(id),
	}
	action, err := h.activityLogClient.GetAllLogAction(r.Context(), preq)
	if err != nil {
		logger.Log().Errorw("Error when getting list action", "error", err)
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting list action")
		return
	}
	resp := []*dto.ActionRequest{}
	for _, v := range action.UserActivityLog {
		resp = append(resp, &dto.ActionRequest{
			ID:     v.Id,
			Action: v.Log,
		})
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToActionResponse(resp))
}
