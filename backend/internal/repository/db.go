package repository

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
	"github.com/stockapp/backend/internal/config"
)

var DB *sql.DB

func InitDB(cfg *config.Config) error {
	var err error
	DB, err = sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("Connected to SQLite:", cfg.DBPath)
	return runMigrations()
}

func runMigrations() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS usuarios (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			rol TEXT DEFAULT 'usuario' CHECK(rol IN ('admin', 'usuario')),
			nombre TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS categorias (
			id TEXT PRIMARY KEY,
			nombre TEXT UNIQUE NOT NULL,
			descripcion TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS productos (
			id TEXT PRIMARY KEY,
			codigo TEXT UNIQUE NOT NULL,
			codigo_barras TEXT,
			nombre TEXT NOT NULL,
			descripcion TEXT,
			categoria_id TEXT REFERENCES categorias(id) ON DELETE SET NULL,
			precio REAL NOT NULL DEFAULT 0,
			stock_actual INTEGER NOT NULL DEFAULT 0,
			stock_minimo INTEGER NOT NULL DEFAULT 0,
			imagen TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS movimientos (
			id TEXT PRIMARY KEY,
			producto_id TEXT NOT NULL REFERENCES productos(id) ON DELETE CASCADE,
			tipo TEXT NOT NULL CHECK(tipo IN ('entrada', 'salida')),
			cantidad INTEGER NOT NULL,
			motivo TEXT,
			usuario_id TEXT REFERENCES usuarios(id),
			fecha DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS devices (
			id TEXT PRIMARY KEY,
			user_id TEXT REFERENCES usuarios(id) ON DELETE CASCADE,
			token TEXT NOT NULL,
			platform TEXT CHECK(platform IN ('android', 'ios', 'web')),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_productos_codigo ON productos(codigo)`,
		`CREATE INDEX IF NOT EXISTS idx_productos_codigo_barras ON productos(codigo_barras)`,
		`CREATE INDEX IF NOT EXISTS idx_productos_categoria ON productos(categoria_id)`,
		`CREATE INDEX IF NOT EXISTS idx_movimientos_producto ON movimientos(producto_id)`,
		`CREATE INDEX IF NOT EXISTS idx_movimientos_fecha ON movimientos(fecha)`,
		`CREATE INDEX IF NOT EXISTS idx_devices_token ON devices(token)`,
	}

	for _, m := range migrations {
		if _, err := DB.Exec(m); err != nil {
			return err
		}
	}

	log.Println("Migrations completed")
	return nil
}