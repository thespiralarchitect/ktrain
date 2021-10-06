package handler

import (
	"encoding/json"
	"io/ioutil"
	"ktrain/cmd/api/user-api/mapper"
	"ktrain/cmd/model"
	"ktrain/cmd/repository"
	"ktrain/pkg/httpreq"
	"ktrain/pkg/httputil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type userHandler struct {
	userRepository repository.IUserRepository
}

func NewUserHandler(userRepository repository.IUserRepository) *userHandler {
	return &userHandler{
		userRepository: userRepository,
	}
}
func readBodyRequest (w http.ResponseWriter,r *http.Request, u *httpreq.UserRequest )  {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
    if err != nil {
        httputil.RespondError(w, http.StatusInternalServerError, "Error read body")
        return
    }
	err2 := json.Unmarshal(b,&u)
	if err2 != nil {
        httputil.RespondError(w, http.StatusInternalServerError, "Error unmarshal ")
        return
    }
}
func (h *userHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := h.userRepository.GetUserByID(ctx.Value("userID").(int64))
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting user profile")
		return
	}

	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToUserResponse(user))
}

func (h *userHandler) GetListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepository.GetListUser()
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting users list")
		return
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToListUsersResponse(users))
}

func (h *userHandler) GetInformationUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		httputil.RespondError(w, http.StatusForbidden, "Error happend with user id ")
		return
	}
	user, err := h.userRepository.GetUserByID(int64(userID))
	if err != nil {
		if  err.Error() == "record not found"{
			httputil.RespondError(w, http.StatusNotFound, "user id not exist")
			return
		}else{
			httputil.RespondError(w, http.StatusInternalServerError, "Error when getting user information")
			return
		}
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToInformationUserResponse(user))
}

func (h *userHandler) PostNewUser(w http.ResponseWriter, r *http.Request){
	var u httpreq.UserRequest
	readBodyRequest(w,r,&u)
	// b, err := ioutil.ReadAll(r.Body)
	// defer r.Body.Close()
    // if err != nil {
    //     httputil.RespondError(w, http.StatusInternalServerError, "Error read body")
    //     return
    // }
	// err2 := json.Unmarshal(b,&u)
	// if err2 != nil {
    //     httputil.RespondError(w, http.StatusInternalServerError, "Error unmarshal ")
    //     return
    // }
	User := &model.User{
		Fullname:   u.Fullname,
		Username:   u.Username,
		Gender:     u.Gender,
		Birthday:   u.Birthday,
	}
	newUser, err := h.userRepository.CreateUser(User)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting user create")
		return
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToUserResponse(newUser))
}