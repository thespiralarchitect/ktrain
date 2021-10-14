package handler

import (
	"ktrain/cmd/api/user-api/dto"
	"ktrain/cmd/api/user-api/mapper"
	middleware2 "ktrain/cmd/api/user-api/middleware"
	"ktrain/cmd/repository"
	"ktrain/pkg/httputil"
	"ktrain/proto/pb"
	"net/http"
)

type activityLogHandler struct {
	//activityLogRepository repository.ActivityLogRepository
	activityLogClient pb.ActivityLogDMSServiceClient
}

func NewActivityLogHandler(activityLogClient pb.ActivityLogDMSServiceClient, activityLogRepository repository.ActivityLogRepository) *activityLogHandler {
	return &activityLogHandler{
		//activityLogRepository: activityLogRepository,
		activityLogClient: activityLogClient,
	}
}
func (h *activityLogHandler) GetActivity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	preq := &pb.GetLogActionRequest{
		Id: ctx.Value(middleware2.ContextKey("userID")).(int64),
	}
	action, err := h.activityLogClient.GetAllLogAction(r.Context(), preq)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting list action")
		return
	}
	resp := []*dto.ActionRequest{}
	for _, v := range action.UserActivityLogMessage {
		resp = append(resp, &dto.ActionRequest{
			ID:     v.Id,
			Action: v.Log,
		})
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToActionResponse(resp))
}
