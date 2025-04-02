-- ユーザーテーブル
CREATE TABLE users (
  id VARCHAR(36) PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  name VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  role ENUM('user', 'admin') DEFAULT 'user',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- マップテーブル
CREATE TABLE maps (
  id VARCHAR(36) PRIMARY KEY,
  map_id VARCHAR(255) UNIQUE NOT NULL,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  user_id VARCHAR(36) NOT NULL,
  is_publicly_editable BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- フロア（エリア）テーブル
CREATE TABLE floors (
  id VARCHAR(36) PRIMARY KEY,
  map_id VARCHAR(36) NOT NULL,
  floor_number INT NOT NULL,
  name VARCHAR(255) NOT NULL,
  image_url TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (map_id) REFERENCES maps(id) ON DELETE CASCADE
);

-- ピンテーブル
CREATE TABLE pins (
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
);

-- 公開編集者テーブル
CREATE TABLE public_editors (
  id VARCHAR(36) PRIMARY KEY,
  map_id VARCHAR(36) NOT NULL,
  nickname VARCHAR(255) NOT NULL,
  editor_token VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (map_id) REFERENCES maps(id) ON DELETE CASCADE
);

-- 管理者ユーザーの作成（パスワードは適切に変更してください）
INSERT INTO users (id, email, name, password, role, created_at)
VALUES (
  UUID(),
  'admin@example.com',
  'Admin User',
  '$2a$10$JqKCLdu1fJdqMZJMxZqO7O0RSl7GQv8MR4q4PrFGqU/xB6e7/y9tO', -- 'password123'
  'admin',
  NOW()
);