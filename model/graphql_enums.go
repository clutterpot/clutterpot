package model

import (
	"fmt"
	"io"
	"strconv"
)

var userKindValues = map[string]UserKind{
	"BANNED": UserKindBanned,
	"USER":   UserKindUser,
	"ADMIN":  UserKindAdmin,
}

func (e UserKind) IsValid() bool {
	switch e {
	case UserKindBanned, UserKindUser, UserKindAdmin:
		return true
	}
	return false
}

func (e *UserKind) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = userKindValues[str]
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UserKind enum", str)
	}
	return nil
}

func (e UserKind) MarshalGQL(w io.Writer) {
	switch e {
	case UserKindBanned:
		fmt.Fprint(w, strconv.Quote("BANNED"))
	case UserKindUser:
		fmt.Fprint(w, strconv.Quote("USER"))
	case UserKindAdmin:
		fmt.Fprint(w, strconv.Quote("ADMIN"))
	}
}

var userSortValues = map[string]UserSort{
	"CREATED": UserSortCreated,
	"UPDATED": UserSortUpdated,
}

func (e UserSort) IsValid() bool {
	switch e {
	case UserSortCreated, UserSortUpdated:
		return true
	}
	return false
}

func (e *UserSort) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = userSortValues[str]
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UserSort enum", str)
	}
	return nil
}

func (e UserSort) MarshalGQL(w io.Writer) {
	switch e {
	case UserSortCreated:
		fmt.Fprint(w, strconv.Quote("CREATED"))
	case UserSortUpdated:
		fmt.Fprint(w, strconv.Quote("UPDATED"))
	}
}
