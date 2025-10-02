package infra

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupDB initializes the database connection and creates tables if they don't exist
func SetupDB() *gorm.DB {
	sslmode := "require"

	// local環境での接続設定
	if os.Getenv("DB_HOST") == "posttgresql-db" {
		//	sslmodeを無効にする
		sslmode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Tokyo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Check if 'users' table exists
	var exists bool
	checkSQL := `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users');`
	if err := db.Raw(checkSQL).Scan(&exists).Error; err != nil {
		log.Fatalf("Failed to check existing tables: %v", err)
	}

	if exists {
		log.Println("Database schema already exists. Skipping migration.")
		return db
	}

	log.Println("Running database migrations...")

	// Migrations
	createSQL := `
	CREATE EXTENSION IF NOT EXISTS "pgcrypto";

	CREATE TABLE users (
		id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(100) NOT NULL,
		role VARCHAR(50) DEFAULT 'user',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE apiSchema (
		id SERIAL PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		field_id VARCHAR(100) NOT NULL,
		view_name VARCHAR(100) NOT NULL,
		field_type VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE field_data (
		id SERIAL PRIMARY KEY,
		apiSchema_id INT NOT NULL REFERENCES apiSchema(id) ON DELETE CASCADE,
		field_type VARCHAR(50) NOT NULL,
		field_value JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE list_options (
		id SERIAL PRIMARY KEY,
		apiSchema_id INT NOT NULL REFERENCES apiSchema(id) ON DELETE CASCADE,
		value VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE api_kind_relation (
		id SERIAL PRIMARY KEY,
		apiSchema_id INT NOT NULL REFERENCES apiSchema(id) ON DELETE CASCADE,
		related_id INT NOT NULL REFERENCES apiSchema(id) ON DELETE CASCADE,
		relation_type VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if err := db.Exec(createSQL).Error; err != nil {
		log.Fatalf("Error executing table creation: %v", err)
	}

	// Create trigger functions and triggers
	triggerSQL := `
	DO $$
	BEGIN
		-- cascade_delete_apiSchema
		IF NOT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'cascade_delete_apischema') THEN
			CREATE FUNCTION cascade_delete_apischema()
			RETURNS TRIGGER AS $$
			BEGIN
				DELETE FROM field_data WHERE apiSchema_id = OLD.id;
				DELETE FROM list_options WHERE apiSchema_id = OLD.id;
				DELETE FROM api_kind_relation WHERE apiSchema_id = OLD.id OR related_id = OLD.id;
				RETURN OLD;
			END;
			$$ LANGUAGE plpgsql;
		END IF;

		-- validate_list_options
		IF NOT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'validate_list_options') THEN
			CREATE FUNCTION validate_list_options()
			RETURNS TRIGGER AS $$
			DECLARE
				schema_type VARCHAR(50);
			BEGIN
				SELECT field_type INTO schema_type FROM apiSchema WHERE id = NEW.apiSchema_id;
				IF schema_type NOT IN ('select', 'dropdown') THEN
					RAISE EXCEPTION 'list_options can only be added to select or dropdown fields';
				END IF;
				RETURN NEW;
			END;
			$$ LANGUAGE plpgsql;
		END IF;

		-- check_cyclic_relation
		IF NOT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'check_cyclic_relation') THEN
			CREATE FUNCTION check_cyclic_relation()
			RETURNS TRIGGER AS $$
			DECLARE
				is_cyclic BOOLEAN;
			BEGIN
				WITH RECURSIVE relation_path AS (
					SELECT related_id FROM api_kind_relation WHERE apiSchema_id = NEW.related_id
					UNION ALL
					SELECT r.related_id FROM api_kind_relation r INNER JOIN relation_path rp ON rp.related_id = r.apiSchema_id
				)
				SELECT EXISTS (
					SELECT 1 FROM relation_path WHERE related_id = NEW.apiSchema_id
				) INTO is_cyclic;

				IF is_cyclic THEN
					RAISE EXCEPTION 'Cyclic relation detected';
				END IF;
				RETURN NEW;
			END;
			$$ LANGUAGE plpgsql;
		END IF;
	END
	$$;

	-- Triggers
	DO $$
	BEGIN
		IF NOT EXISTS (
			SELECT 1 FROM pg_trigger WHERE tgname = 'delete_related_data'
		) THEN
			CREATE TRIGGER delete_related_data
			AFTER DELETE ON apiSchema
			FOR EACH ROW EXECUTE FUNCTION cascade_delete_apischema();
		END IF;

		IF NOT EXISTS (
			SELECT 1 FROM pg_trigger WHERE tgname = 'validate_options'
		) THEN
			CREATE TRIGGER validate_options
			BEFORE INSERT ON list_options
			FOR EACH ROW EXECUTE FUNCTION validate_list_options();
		END IF;

		IF NOT EXISTS (
			SELECT 1 FROM pg_trigger WHERE tgname = 'prevent_cyclic_relation'
		) THEN
			CREATE TRIGGER prevent_cyclic_relation
			BEFORE INSERT ON api_kind_relation
			FOR EACH ROW EXECUTE FUNCTION check_cyclic_relation();
		END IF;
	END
	$$;
	`

	if err := db.Exec(triggerSQL).Error; err != nil {
		log.Fatalf("Error executing trigger SQL: %v", err)
	}

	log.Println("Database migration completed successfully.")
	return db
}
