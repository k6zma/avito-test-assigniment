package mappers

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDFromPg(x pgtype.UUID) uuid.UUID {
	var id uuid.UUID

	if err := id.UnmarshalBinary(x.Bytes[:]); err != nil {
		return uuid.Nil
	}

	return id
}

func UUIDToPg(id uuid.UUID) pgtype.UUID {
	var out pgtype.UUID
	if err := out.Scan(id.String()); err != nil {
		return pgtype.UUID{}
	}

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
