package handler

import (
	"fmt"
	"ktrain/cmd/api/user-api/dto"
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
func (h *mongoHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := dto.UserResquest{
		Fullname: "Hieu",
		Username: "",
		Gender:   "",
		Birthday: "",
	}
	fmt.Println("ok1")
	resp, err := h.mongoRepository.CreateUser(&user)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when create user")
		return
	}
	user1 := mapper.ToUserModel(resp)
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToUserResponse(user1))
}
