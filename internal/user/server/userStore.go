package userserver

import (
	"encoding/json"
	"hash/fnv"
	"net/http"
)

type UserStore struct {
	Users  map[uint32]User
	NextId uint32
}

// Make sure we conform to ServerInterface
var _ ServerInterface = (*UserStore)(nil)

func NewUserStore() *UserStore {
	return &UserStore{
		Users:  make(map[uint32]User),
		NextId: 1000,
	}
}

// This function wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendUserStoreError(w http.ResponseWriter, code int, message string) {
	userErr := Error{
		Code:    int(code),
		Message: message,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(userErr)
}

func genHash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (u *UserStore) GetUsers(w http.ResponseWriter, r *http.Request) {
	var result []User
	for _, user := range u.Users {
		result = append(result, user)
		// TODO: Add pagination
		if len(result) > 10 {
			break
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (u *UserStore) PostUser(w http.ResponseWriter, r *http.Request) {
	var newUser NewUser
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		sendUserStoreError(w, http.StatusBadRequest, "Invalid format for NewUser")
		return
	}

	id := genHash(string(newUser.Email))
	_, ok := u.Users[id]
	if ok {
		sendUserStoreError(w, http.StatusBadRequest, "User Already Exist")
	}

	var user User
	user.Name = newUser.Name
	user.Email = newUser.Email
	user.Role = newUser.Role
	id_int := int(id)
	user.Id = &id_int
	u.Users[id] = user

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (u *UserStore) DeleteUsersUserId(w http.ResponseWriter, r *http.Request, userId int) {}

func (u *UserStore) GetUsersUserId(w http.ResponseWriter, r *http.Request, userId int) {}

func (u *UserStore) PatchUsersUserId(w http.ResponseWriter, r *http.Request, userId int) {}
