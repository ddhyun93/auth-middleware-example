package api

import (
	"encoding/json"
	"fmt"
	"go-auth-with-chi/domain"
	"go-auth-with-chi/dto"
	"go-auth-with-chi/ioc"
	"go-auth-with-chi/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var data dto.Msg
	var input dto.CreateUserReq
	var res dto.TokenRes
	data.Result = "fail"

	// check method
	if r.Method != http.MethodPost {
		errMsg := fmt.Sprintf("http method type %s not supported", r.Method)
		data.Error = errMsg
		jsonMsg, _ := json.Marshal(data)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(jsonMsg)
		return
	}

	// decode input
	_ = json.NewDecoder(r.Body).Decode(&input)

	// validate input
	ok, msg := input.Validate()
	if !ok {
		data.Error = msg
		jsonMsg, _ := json.Marshal(data)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(jsonMsg)
		return
	}

	// check user
	userDao, err := ioc.Repo.Users.GetByEmail(input.Email)
	if err != nil {
		errMsg := fmt.Sprintf("can not find user %s", input.Email)
		data.Error = errMsg
		jsonMsg, _ := json.Marshal(data)
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonMsg)
		return
	}

	// check hashed password
	if !utils.CheckHashedPassword(userDao.Password, input.Password) {
		data.Error = "unmatched password"
		jsonMsg, _ := json.Marshal(data)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(jsonMsg)
		return
	}

	// gen token
	token, _ := utils.GenerateToken(userDao.ID.Hex(), 24)
	refreshToken, _ := utils.GenerateToken(userDao.ID.Hex(), 24*365)
	res.Token = token
	res.RefreshToken = refreshToken
	data.Data = res
	data.Result = "success"

	// update device info & refresh token
	userDao.UpdatedAt = time.Now()
	userDao.RefreshToken = refreshToken
	userDao, _ = ioc.Repo.Users.Upsert(userDao)

	// return
	jsonMsg, _ := json.Marshal(data)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonMsg)
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	var data dto.Msg
	var input dto.CreateUserReq
	var res dto.TokenRes
	var user domain.UserDAO
	data.Result = "fail"

	// check method
	if r.Method != http.MethodPost {
		errMsg := fmt.Sprintf("http method type %s not supported", r.Method)
		data.Error = errMsg
		jsonMsg, _ := json.Marshal(data)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(jsonMsg)
		return
	}

	// decode input body
	_ = json.NewDecoder(r.Body).Decode(&input)

	// check duplicated user
	_, err := ioc.Repo.Users.GetByEmail(input.Email)
	if err == nil {
		errMsg := fmt.Sprintf("already registered username %s", input.Email)
		data.Error = errMsg
		jsonMsg, _ := json.Marshal(data)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(jsonMsg)
		return
	}

	// validate input
	ok, msg := input.Validate()
	if !ok {
		data.Error = msg
		jsonMsg, _ := json.Marshal(data)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(jsonMsg)
		return
	}

	// hash password
	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		data.Error = err.Error()
		jsonMsg, _ := json.Marshal(data)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(jsonMsg)
		return
	}

	// persist user model
	user.ID = primitive.NewObjectID()
	user.Email = input.Email
	user.Password = hashed
	user.ActivationCode = utils.CreateRandCode()
	userDao, err := ioc.Repo.Users.Upsert(&user)
	if err != nil {
		data.Error = err.Error()
		jsonMsg, _ := json.Marshal(data)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonMsg)
		return
	}

	// generate token
	token, _ := utils.GenerateToken(userDao.ID.Hex(), 24)
	res.Token = token

	// marshaling json
	data.Data = res
	data.Result = "success"
	jsonMsg, _ := json.Marshal(data)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonMsg)
	go utils.SendAuthMail(userDao.Email, user.ActivationCode)
	return
}
