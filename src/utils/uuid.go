package utils

import (
	"w3st/errors"

	"github.com/google/uuid"
)

func StringToUUID(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, errors.NewDomainErrorWithMessage(errors.InvalidParameter, "invalid uuid")
	}
	return id, nil
}

func UuidToString(id uuid.UUID) string {
	return id.String()
}

func UuidToUint(id uuid.UUID) (uint, error) {
	var uintID uint

	for i := 0; i < len(id); i++ {
		uintID += uint(id[i])
	}

	return uintID, nil
}
