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
	// データソース名を環境変数から取得
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// データベースに接続
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// すでにマイグレーション済みかチェック (users テーブルがあるか確認)
	var exists bool
	checkSQL := `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users');`
	if err := db.Raw(checkSQL).Scan(&exists).Error; err != nil {
		log.Fatalf("Failed to check existing tables: %v", err)
	}

	// すでにテーブルがある場合はマイグレーションをスキップ
	if exists {
		log.Println("Database schema already exists. Skipping migration.")
		return db
	}

	log.Println("Running database migrations...")

	// マイグレーション
	initSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(100) NOT NULL,
		role VARCHAR(50) DEFAULT 'user',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS apiSchema (
		id SERIAL PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		field_id VARCHAR(100) NOT NULL,
		view_name VARCHAR(100) NOT NULL,
		field_type VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS field_data (
		id SERIAL PRIMARY KEY,
		apiSchema_id INT NOT NULL REFERENCES apiSchema(id) ON DELETE CASCADE,
		field_type VARCHAR(50) NOT NULL,
		field_value JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS list_options (
		id SERIAL PRIMARY KEY,
		apiSchema_id INT NOT NULL REFERENCES apiSchema(id) ON DELETE CASCADE,
		value VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS api_kind_relation (
		id SERIAL PRIMARY KEY,
		apiSchema_id INT NOT NULL REFERENCES apiSchema(id) ON DELETE CASCADE,
		related_id INT NOT NULL REFERENCES apiSchema(id) ON DELETE CASCADE,
		relation_type VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- トリガー関数の作成 (IF NOT EXISTS はトリガー関数には使えないので、関数がすでに存在するかチェック)
	DO $$ BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'cascade_delete_apiSchema') THEN
			CREATE OR REPLACE FUNCTION cascade_delete_apiSchema()
			RETURNS TRIGGER AS $$
			BEGIN
				DELETE FROM field_data WHERE apiSchema_id = OLD.id;
				DELETE FROM list_options WHERE apiSchema_id = OLD.id;
				DELETE FROM api_kind_relation WHERE apiSchema_id = OLD.id OR related_id = OLD.id;
				RETURN OLD;
			END;
			$$ LANGUAGE plpgsql;
		END IF;
	END $$;

	DO $$ BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'validate_list_options') THEN
			CREATE OR REPLACE FUNCTION validate_list_options()
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
	END $$;

	DO $$ BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'check_cyclic_relation') THEN
			CREATE OR REPLACE FUNCTION check_cyclic_relation()
			RETURNS TRIGGER AS $$
			DECLARE
				is_cyclic BOOLEAN;
			BEGIN
				WITH RECURSIVE relation_path AS (
					SELECT related_id FROM api_kind_relation WHERE apiSchema_id = NEW.related_id
					UNION ALL
					SELECT r.related_id FROM api_kind_relation r INNER JOIN relation_path rp ON rp.related_id = r.apiSchema_id
				)
				SELECT EXISTS (SELECT 1 FROM relation_path WHERE related_id = NEW.apiSchema_id) INTO is_cyclic;
				IF is_cyclic THEN
					RAISE EXCEPTION 'Cyclic relation detected';
				END IF;
				RETURN NEW;
			END;
			$$ LANGUAGE plpgsql;
		END IF;
	END $$;

	-- 各テーブルのトリガー設定
	CREATE TRIGGER IF NOT EXISTS delete_related_data
	AFTER DELETE ON apiSchema
	FOR EACH ROW EXECUTE FUNCTION cascade_delete_apiSchema();

	CREATE TRIGGER IF NOT EXISTS validate_options
	BEFORE INSERT ON list_options
	FOR EACH ROW EXECUTE FUNCTION validate_list_options();

	CREATE TRIGGER IF NOT EXISTS prevent_cyclic_relation
	BEFORE INSERT ON api_kind_relation
	FOR EACH ROW EXECUTE FUNCTION check_cyclic_relation();
	`

	// SQL実行
	if err := db.Exec(initSQL).Error; err != nil {
		log.Fatalf("Error executing initSQL: %v", err)
	}

	log.Println("Database migration completed successfully.")
	return db
}
