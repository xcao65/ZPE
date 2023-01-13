package userserver

import (
	"encoding/json"
	"hash/fnv"
	"net/http"
)

type UserStore struct {
	Users    map[int]User
	Priority map[string]int
}

// Make sure we conform to ServerInterface
var _ ServerInterface = (*UserStore)(nil)

func NewUserStore() *UserStore {
	return &UserStore{
		Users:    make(map[int]User),
		Priority: map[string]int{string(Admin): 2, string(Modifier): 1, string(Watcher): 0},
	}
}

// This function wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendUserStoreError(w http.ResponseWriter, code int, message string) {
	code_int := &code
	userErr := Error{
		Code:    code_int,
		Message: &message,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(userErr)
}

func genHash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
		sendUserStoreError(w, 409, "User Already Exist")
		return
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

func (u *UserStore) DeleteUsersUserId(w http.ResponseWriter, r *http.Request, userId int) {
	_, ok := u.Users[userId]
	if !ok {
		sendUserStoreError(w, 404, "User does not exist")
		return
	}
	delete(u.Users, userId)
}

func (u *UserStore) GetUsersUserId(w http.ResponseWriter, r *http.Request, userId int) {
	user, ok := u.Users[userId]
	if !ok {
		sendUserStoreError(w, 404, "User does not exist")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (u *UserStore) PatchUsersUserId(w http.ResponseWriter, r *http.Request, userId int) {
	var newUser NewUser
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		sendUserStoreError(w, http.StatusBadRequest, "Invalid format for NewUser")
		return
	}

	user, ok := u.Users[userId]
	if !ok {
		sendUserStoreError(w, 404, "User does not exist")
		return
	}

	oldRoles := make(map[string]int)
	var maxP int
	maxP = 0
	for _, role := range user.Role {
		rP := u.Priority[string(role)]
		maxP = max(maxP, rP)
		oldRoles[string(role)] = rP
	}
	for _, role := range newUser.Role {
		_, exist := oldRoles[string(role)]
		if !exist {
			if u.Priority[string(role)] > maxP {
				sendUserStoreError(w, 403, "You don't have previllige to add this role")
				return
			}
			user.Role = append(user.Role, role)
			oldRoles[string(role)] = u.Priority[string(role)]
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
