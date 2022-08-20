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

var fileSortValues = map[string]FileSort{
	"CREATED": FileSortCreated,
	"UPDATED": FileSortUpdated,
}

func (e FileSort) IsValid() bool {
	switch e {
	case FileSortCreated, FileSortUpdated:
		return true
	}
	return false
}

func (e *FileSort) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = fileSortValues[str]
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FileSort enum", str)
	}
	return nil
}

func (e FileSort) MarshalGQL(w io.Writer) {
	switch e {
	case FileSortCreated:
		fmt.Fprint(w, strconv.Quote("CREATED"))
	case FileSortUpdated:
		fmt.Fprint(w, strconv.Quote("UPDATED"))
	}
}

var tagSortValues = map[string]TagSort{
	"CREATED": TagSortCreated,
	"UPDATED": TagSortUpdated,
}

func (e TagSort) IsValid() bool {
	switch e {
	case TagSortCreated, TagSortUpdated:
		return true
	}
	return false
}

func (e *TagSort) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = tagSortValues[str]
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TagSort enum", str)
	}
	return nil
}

func (e TagSort) MarshalGQL(w io.Writer) {
	switch e {
	case TagSortCreated:
		fmt.Fprint(w, strconv.Quote("CREATED"))
	case TagSortUpdated:
		fmt.Fprint(w, strconv.Quote("UPDATED"))
	}
}
