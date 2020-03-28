package users

import "encoding/json"

type PublicUser struct {
	Id          int64  `json:"id"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
}

type PrivateUser struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
}

func (users Users) Marshall(isPubic bool) interface{} {
	var results = make([]interface{}, 0)
	for _, user := range users {
		results = append(results, user.Marshall(isPubic))
	}
	return results
}

func (user *User) Marshall(isPublic bool) interface{} {
	if isPublic {
		return PublicUser{
			Id:          user.Id,
			DateCreated: user.DateCreated,
			Status:      user.Status,
		}
	}
	userJSON, _ := json.Marshal(user)
	var privateUser PrivateUser
	if err := json.Unmarshal(userJSON, &privateUser); err != nil {
		return nil
	}
	return privateUser
}
