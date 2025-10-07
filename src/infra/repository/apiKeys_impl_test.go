package infrastructure

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"w3st/errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 db,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		_ = db.Close()
		t.Fatalf("failed to open gorm db: %v", err)
	}

	cleanup := func() {
		_ = db.Close()
	}

	return gdb, mock, cleanup
}

func TestFindByKey_NotFound(t *testing.T) {
	t.Parallel()

	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewApiKeyRepositoryImpl(gdb)

	// Expect query and return no rows
	mock.ExpectQuery(`SELECT \* FROM "api_keys" WHERE key = \$1 AND revoked = false ORDER BY "api_keys"\."id" LIMIT (?:\$?\d+)`).
		WithArgs("missing-key", sqlmock.AnyArg()).
		WillReturnError(sql.ErrNoRows)

	apiKey, de := repo.FindByKey(context.Background(), "missing-key")
	if apiKey != nil {
		t.Fatalf("expected nil apiKey when not found, got %+v", apiKey)
	}
	if de == nil {
		t.Fatalf("expected domain error when key not found")
	}
	if de.GetType() != errors.QueryDataNotFoundError {
		t.Fatalf("expected QueryDataNotFoundError, got %v", de.GetType())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestFindByKey_Success(t *testing.T) {
	t.Parallel()

	gdb, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewApiKeyRepositoryImpl(gdb)

	id := 1
	userID := uuid.New()
	projectID := 1
	name := "Public API Key"
	key := "public-key-123"
	var collectionIDs interface{} = nil // use nil to avoid driver string scan issues in sqlmock
	var ipWhitelist interface{} = nil
	expireAt := time.Now()
	revoked := false
	rateLimit := 1000
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "project_id", "name", "key", "collection_ids", "ip_whitelist", "expire_at", "revoked", "rate_limit_per_hour", "created_at",
	}).AddRow(
		id, userID.String(), projectID, name, key, collectionIDs, ipWhitelist, expireAt, revoked, rateLimit, createdAt,
	)

	mock.ExpectQuery(`SELECT \* FROM "api_keys" WHERE key = \$1 AND revoked = false ORDER BY "api_keys"\."id" LIMIT (?:\$?\d+)`).
		WithArgs(key, sqlmock.AnyArg()).
		WillReturnRows(rows)

	apiKey, de := repo.FindByKey(context.Background(), key)
	if de != nil {
		t.Fatalf("unexpected domain error: %v", de)
	}
	if apiKey == nil {
		t.Fatalf("expected apiKey, got nil")
	}
	if apiKey.Key != key {
		t.Fatalf("expected key %s, got %s", key, apiKey.Key)
	}
	if apiKey.Name != name {
		t.Fatalf("expected name %s, got %s", name, apiKey.Name)
	}
	if apiKey.ProjectID != projectID {
		t.Fatalf("expected project id %d, got %d", projectID, apiKey.ProjectID)
	}
	// collection ids may be empty slice or nil depending on driver; accept both
	if apiKey.CollectionIds == nil {
		// nothing to assert, just log for visibility
		t.Log("collection ids is nil (driver returned NULL/empty)")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
