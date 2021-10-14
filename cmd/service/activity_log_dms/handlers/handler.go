package handlers

import (
	"context"
	"ktrain/cmd/repository"
	"ktrain/proto/pb"
)

type ActivityLogHandler struct {
	pb.UnimplementedActivityLogDMSServiceServer
	activityLogRepository repository.ActivityLogRepository
}

func NewActivityLogHandler(activityLogRepository repository.ActivityLogRepository) (*ActivityLogHandler, error) {
	return &ActivityLogHandler{
		activityLogRepository: activityLogRepository,
	}, nil
}
func (h *ActivityLogHandler) CreateAction(ctx context.Context, in *pb.CreateActionRequest) (*pb.CreateActionResponse, error) {
	action, err := h.activityLogRepository.CreateAction(ctx, in.Id, in.Log)
	if err != nil {
		return nil, err
	}
	return &pb.CreateActionResponse{
		Id:  action.ID,
		Log: action.Log,
	}, nil
}
func (h *ActivityLogHandler) GetAllLogAction(ctx context.Context, in *pb.GetLogActionRequest) (*pb.GetLogActionResponse, error) {
	listaction, err := h.activityLogRepository.GetAllLogAction(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	resp := []*pb.UserActivityLog{}
	for _, v := range listaction {
		action := &pb.UserActivityLog{
			Id:  v.ID,
			Log: v.Action,
		}
		resp = append(resp, action)
	}
	return &pb.GetLogActionResponse{
		UserActivityLog: resp,
	}, nil
}
