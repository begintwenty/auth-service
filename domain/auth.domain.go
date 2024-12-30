package domain

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Auth struct {
	UserID      UserID      `bson:"user_id" json:"user_id"`
	Permissions Permissions `bson:"permissions" json:"permissions"`
	UserType    UserType    `bson:"user_type" json:"user_type"`
}

type UserType string
type UserID string

const (
	UserTypeAdmin UserType = "admin"
	UserTypeUser  UserType = "user"
)

type Permissions map[Permission]bool

type Permission string

const (
	PermissionUserAuth Permission = "user.auth"
)

func NewUserID() UserID {
	return UserID("user_" + strings.ReplaceAll(uuid.New().String(), "-", ""))
}

func (u *UserID) AsString() string {
	return string(*u)
}

func VerifyUserID(userId string) (*UserID, error) {
	if !strings.HasPrefix(userId, "user_") {
		return nil, fmt.Errorf("invalid user ID: %s", userId)
	}
	validatedUserId := UserID(userId)
	return &validatedUserId, nil
}

func (p Permissions) MarshalJSON() ([]byte, error) {
	m := make(map[string]bool, len(p))
	for k, v := range p {
		m[string(k)] = v
	}
	return json.Marshal(m)
}

func (p *Permissions) UnmarshalJSON(data []byte) error {
	m := make(map[string]bool)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	*p = make(Permissions, len(m))
	for k, v := range m {
		(*p)[Permission(k)] = v
	}
	return nil
}
