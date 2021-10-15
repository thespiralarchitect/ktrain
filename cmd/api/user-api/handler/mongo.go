package handler

import (
	"ktrain/cmd/api/user-api/dto"
	"ktrain/cmd/api/user-api/mapper"
	"ktrain/pkg/httputil"
	"ktrain/proto/pb"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type activityLogHandler struct {
	activityLogClient pb.ActivityLogDMSServiceClient
	logger            *zap.SugaredLogger
}

func NewActivityLogHandler(activityLogClient pb.ActivityLogDMSServiceClient, logger *zap.SugaredLogger) *activityLogHandler {
	return &activityLogHandler{
		activityLogClient: activityLogClient,
		logger:            logger,
	}
}
func (h *activityLogHandler) GetActivity(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	preq := &pb.GetLogActionRequest{
		Id: int64(id),
	}
	action, err := h.activityLogClient.GetAllLogAction(r.Context(), preq)
	if err != nil {
		h.logger.Errorw("Error when getting list action", "error", err)
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
