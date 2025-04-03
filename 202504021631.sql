-- データベースの作成
CREATE DATABASE IF NOT EXISTS mapapp CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE mapapp;

SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS public_editors;
DROP TABLE IF EXISTS pins;
DROP TABLE IF EXISTS floors;
DROP TABLE IF EXISTS maps;
DROP TABLE IF EXISTS users;

SET FOREIGN_KEY_CHECKS = 1;


-- ユーザーテーブル
CREATE TABLE IF NOT EXISTS users (
  id VARCHAR(36) PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  name VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  role ENUM('user', 'admin') DEFAULT 'user',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- マップテーブル
CREATE TABLE IF NOT EXISTS maps (
  id VARCHAR(36) PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  user_id VARCHAR(36) NOT NULL,
  is_publicly_editable BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- フロア（エリア）テーブル
CREATE TABLE IF NOT EXISTS floors (
  id VARCHAR(36) PRIMARY KEY,
  floor_number INT NOT NULL,
  name VARCHAR(255) NOT NULL,
  image_url TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (map_id) REFERENCES maps(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- ピンテーブル
CREATE TABLE IF NOT EXISTS pins (
  id VARCHAR(36) PRIMARY KEY,
  floor_id VARCHAR(36) NOT NULL,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  x_position DECIMAL(10,2) NOT NULL,
  y_position DECIMAL(10,2) NOT NULL,
  image_url TEXT,
  editor_id VARCHAR(255),
  editor_nickname VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (floor_id) REFERENCES floors(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- 公開編集者テーブル
CREATE TABLE IF NOT EXISTS public_editors (
  id VARCHAR(36) PRIMARY KEY,
  nickname VARCHAR(255) NOT NULL,
  editor_token VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (map_id) REFERENCES maps(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- インデックスの作成
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_maps_map_id ON maps(map_id);
CREATE INDEX idx_maps_user_id ON maps(user_id);
CREATE INDEX idx_floors_map_id ON floors(map_id);
CREATE INDEX idx_floors_floor_number ON floors(map_id, floor_number);
CREATE INDEX idx_pins_floor_id ON pins(floor_id);
CREATE INDEX idx_public_editors_map_id ON public_editors(map_id);
CREATE INDEX idx_public_editors_token ON public_editors(editor_token);

-- サンプル管理者ユーザー（開発環境用）
-- パスワードは 'password123' (bcrypt ハッシュ)
INSERT INTO users (id, email, name, password, role, created_at)
VALUES (
  UUID(),
  'admin@example.com',
  'Admin User',
  '$2a$10$JqKCLdu1fJdqMZJMxZqO7O0RSl7GQv8MR4q4PrFGqU/xB6e7/y9tO',
  'admin',
  NOW()
) ON DUPLICATE KEY UPDATE email = email;