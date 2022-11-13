package users

import (
	"encoding/json"
	"fmt"
	"helpers"
	"jwt"
	"mime/multipart"
	"models"
	"net/http"
	"strconv"
	"strings"
	"time"

	fc "controllers/files"
	db "db/users"
	mr "models/response"
)

/* Insert Permits to create a user in the DB */
func Insert(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(helpers.RequestUserKey{}).(models.User)

	if len(user.Password) < 6 {
		http.Error(w, "The password must have at least 6 characters", http.StatusBadRequest)

		return
	}

	_, isFound, _ := db.IsUser(user.Email)

	if isFound {
		http.Error(w, "The user already exists", http.StatusBadRequest)

		return
	}

	_, status, err := db.InsertUser(user)

	if err != nil {
		http.Error(w, "There was an error trying to regist the user"+err.Error(), http.StatusBadRequest)

		return
	}

	if !status {
		http.Error(w, "The user registry could not be inserted into the DB", http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

/* Login Does the login */
func Login(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(helpers.RequestUserKey{}).(models.User)
	userDb, isUser := db.TryLogin(user.Email, user.Password)

	if !isUser {
		http.Error(w, "User and/or password invalid", http.StatusBadRequest)

		return
	}

	jwtKey, err := jwt.GenerateJWT(userDb)

	if err != nil {
		http.Error(w, "Something went wrong"+err.Error(), http.StatusInternalServerError)

		return
	}

	response := mr.LoginResponse{
		Token: jwtKey,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)

	expirationTime := time.Now().Add(24 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   jwtKey,
		Expires: expirationTime,
	})
}

/* GetProfile Gets an user profile */
func GetProfile(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(helpers.RequestQueryIdKey{}).(string)
	profile, err := db.GetProfile(id)

	if err != nil {
		http.Error(w, "An error occurred when trying to find a registry in the DB: "+err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(profile)
}

/* Modify Allows to modify a registry */
func Modify(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Invalid data"+err.Error(), http.StatusBadRequest)

		return
	}

	status, err := db.ModifyRegistry(user, jwt.UserId)

	if err != nil {
		http.Error(w, "An error has occurred when trying to modify the registry: "+err.Error(), http.StatusBadRequest)

		return
	}

	if !status {
		http.Error(w, "The registry was not modified", http.StatusNotModified)

		return
	}

	w.WriteHeader(http.StatusOK)
}

/* Upload Uploads an user's avatar */
func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	file, header, err := fc.GetRequestFile("avatar", r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	isRemote, _ := strconv.ParseBool(helpers.GetEnvVariable("FILES_REMOTE"))
	filename := ""

	if isRemote {
		var profile models.User

		profile, err = db.GetProfile(jwt.UserId)

		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)

			return
		}

		filename, err = uploadRemote(file, profile.Avatar, "avatar")
	} else {
		filename, err = uploadLocal("uploads/avatars", header, file)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	user := models.User{
		Avatar: filename,
	}

	err = saveToDB(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

/* Upload Uploads an user's banner */
func UploadBanner(w http.ResponseWriter, r *http.Request) {
	file, header, err := fc.GetRequestFile("banner", r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	isRemote, _ := strconv.ParseBool(helpers.GetEnvVariable("FILES_REMOTE"))
	filename := ""

	if isRemote {
		var profile models.User

		profile, err = db.GetProfile(jwt.UserId)

		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)

			return
		}

		filename, err = uploadRemote(file, profile.Banner, "banner")
	} else {
		filename, err = uploadLocal("uploads/banners", header, file)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	user := models.User{
		Banner: filename,
	}

	err = saveToDB(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

/* GetAvatar Gets the user's file avatar */
func GetAvatar(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(helpers.RequestQueryIdKey{}).(string)
	profile, err := db.GetProfile(id)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)

		return
	}

	isRemote, _ := strconv.ParseBool(helpers.GetEnvVariable("FILES_REMOTE"))
	filepath := profile.Avatar

	if !isRemote {
		filepath = fmt.Sprintf("uploads/avatars/%s", profile.Avatar)
	}

	err = fc.SetFileToResponse(filepath, w, isRemote)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// This line is optional because the response defaults to 200 OK
	//w.WriteHeader(http.StatusOK)
}

/* GetBanner Gets the user's file banner */
func GetBanner(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(helpers.RequestQueryIdKey{}).(string)
	profile, err := db.GetProfile(id)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)

		return
	}

	isRemote, _ := strconv.ParseBool(helpers.GetEnvVariable("FILES_REMOTE"))
	filepath := profile.Banner

	if !isRemote {
		filepath = fmt.Sprintf("uploads/banners/%s", profile.Banner)
	}

	err = fc.SetFileToResponse(filepath, w, isRemote)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// This line is optional because the response defaults to 200 OK
	//w.WriteHeader(http.StatusOK)
}

func uploadLocal(filepath string, header *multipart.FileHeader, file multipart.File) (string, error) {
	extension := strings.Split(header.Filename, ".")[1]
	filename := fmt.Sprintf("%s.%s", jwt.UserId, extension)
	filepath = fmt.Sprintf("%s/%s", filepath, filename)

	err := helpers.UploadFileLocal(filepath, file)

	if err != nil {
		return "", err
	}

	return filename, nil
}

func uploadRemote(file multipart.File, fileUrl string, tag string) (string, error) {
	publicId := fmt.Sprintf("%s-%s", jwt.UserId, tag)

	if fileUrl != "" {
		err := helpers.DestroyRemote(publicId)

		if err != nil {
			return "", err
		}
	}

	fileUrl, err := helpers.UploadRemote(file, publicId)

	if err != nil {
		return "", err
	}

	return fileUrl, nil
}

func saveToDB(user models.User) error {
	status, err := db.ModifyRegistry(user, jwt.UserId)

	if err != nil || !status {
		return fmt.Errorf("Error when saving the file in the DB: " + err.Error())
	}

	return nil
}
