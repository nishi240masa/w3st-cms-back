package utils

import "github.com/google/uuid"


func StringToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func UuidToString(id uuid.UUID) (string) {
	return id.String()
}

func UuidToUint(id uuid.UUID) (uint, error) {
    
    var uintID uint

	for i := 0; i < len(id); i++ {
		uintID += uint(id[i])
	}

    return uintID, nil
}