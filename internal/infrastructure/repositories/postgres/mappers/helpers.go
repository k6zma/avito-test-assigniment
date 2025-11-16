package mappers

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDFromPg(x pgtype.UUID) uuid.UUID {
	var id uuid.UUID

	_ = id.UnmarshalBinary(x.Bytes[:])

	return id
}

func UUIDToPg(id uuid.UUID) pgtype.UUID {
	var out pgtype.UUID

	_ = out.Scan(id)

	return out
}

func StringPtr(s string) *string {
	return &s
}

func SafeString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
