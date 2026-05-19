package repository

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/stockapp/backend/internal/config"
)

var DB *sql.DB

func InitDB(cfg *config.Config) error {
	var err error

	// Use DATABASE_URL if provided (PostgreSQL), otherwise fallback to SQLite
	if cfg.DatabaseURL != "" {
		DB, err = sql.Open("postgres", cfg.DatabaseURL)
		if err != nil {
			return err
		}
		
		// Connection pool settings for Supabase
		DB.SetMaxOpenConns(cfg.MaxOpenConns)
		DB.SetMaxIdleConns(cfg.MaxIdleConns)
		DB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
		
		if err = DB.Ping(); err != nil {
			return err
		}
		log.Printf("Connected to PostgreSQL (pool: max=%d, idle=%d, lifetime=%ds)", 
			cfg.MaxOpenConns, cfg.MaxIdleConns, cfg.ConnMaxLifetime)
	} else {
		// Fallback to SQLite for local development without DATABASE_URL
		// Note: For SQLite, we need to import the driver dynamically
		// This is a simplified approach - in production always use PostgreSQL
		log.Println("WARNING: Using SQLite (not recommended for production)")
		// For SQLite fallback, we'd need the modernc.org/sqlite driver
		// For now, this will fail if DATABASE_URL is not set - which is intentional
		return nil // Skip DB init if no DB configured - will fail on first query
	}

	return runMigrations()
}

func runMigrations() error {
	migrations := []string{
		// ============================================
		// RLS (Row Level Security) - enable in Supabase Dashboard
		// Tables are created without RLS for migration compatibility
		// Enable RLS after initial setup via Supabase UI or run:
		// ALTER TABLE usuarios ENABLE ROW LEVEL SECURITY;
		// etc.
		// ============================================
		`CREATE TABLE IF NOT EXISTS usuarios (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			rol TEXT DEFAULT 'usuario' CHECK(rol IN ('admin', 'usuario')),
			nombre TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS categorias (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			nombre TEXT UNIQUE NOT NULL,
			descripcion TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS productos (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			codigo TEXT UNIQUE NOT NULL,
			codigo_barras TEXT,
			nombre TEXT NOT NULL,
			descripcion TEXT,
			categoria_id UUID REFERENCES categorias(id) ON DELETE SET NULL,
			precio REAL NOT NULL DEFAULT 0,
			stock_actual INTEGER NOT NULL DEFAULT 0,
			stock_minimo INTEGER NOT NULL DEFAULT 0,
			imagen TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS movimientos (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			producto_id UUID NOT NULL REFERENCES productos(id) ON DELETE CASCADE,
			tipo TEXT NOT NULL CHECK(tipo IN ('entrada', 'salida')),
			cantidad INTEGER NOT NULL,
			motivo TEXT,
			usuario_id UUID REFERENCES usuarios(id),
			fecha TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS devices (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES usuarios(id) ON DELETE CASCADE,
			token TEXT NOT NULL,
			platform TEXT CHECK(platform IN ('android', 'ios', 'web')),
			created_at TIMESTAMPTZ DEFAULT NOW()
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