package repository

import (
	"fmt"
	"strings"

	"github.com/stockapp/backend/internal/model"
	"github.com/stockapp/backend/internal/utils"
)

func CreateProducto(req *model.CreateProductoRequest) (*model.Producto, error) {
	id := utils.NewUUID()
	
	_, err := DB.Exec(
		`INSERT INTO productos (id, codigo, codigo_barras, nombre, descripcion, categoria_id, precio, stock_actual, stock_minimo) 
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, req.Codigo, req.CodigoBarras, req.Nombre, req.Descripcion, req.CategoriaID, req.Precio, req.StockActual, req.StockMinimo,
	)
	if err != nil {
		return nil, err
	}

	return GetProductoByID(id)
}

func GetProductoByID(id string) (*model.Producto, error) {
	var p model.Producto
	err := DB.QueryRow(
		`SELECT p.id, p.codigo, p.codigo_barras, p.nombre, p.descripcion, p.categoria_id, c.nombre, p.precio, p.stock_actual, p.stock_minimo, p.imagen, p.created_at, p.updated_at 
		 FROM productos p 
		 LEFT JOIN categorias c ON p.categoria_id = c.id 
		 WHERE p.id = ?`, id,
	).Scan(&p.ID, &p.Codigo, &p.CodigoBarras, &p.Nombre, &p.Descripcion, &p.CategoriaID, &p.CategoriaNombre, &p.Precio, &p.StockActual, &p.StockMinimo, &p.Imagen, &p.CreatedAt, &p.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetProductoByCodigo(codigo string) (*model.Producto, error) {
	var p model.Producto
	err := DB.QueryRow(
		`SELECT p.id, p.codigo, p.codigo_barras, p.nombre, p.descripcion, p.categoria_id, c.nombre, p.precio, p.stock_actual, p.stock_minimo, p.imagen, p.created_at, p.updated_at 
		 FROM productos p 
		 LEFT JOIN categorias c ON p.categoria_id = c.id 
		 WHERE p.codigo = ? OR p.codigo_barras = ?`, codigo, codigo,
	).Scan(&p.ID, &p.Codigo, &p.CodigoBarras, &p.Nombre, &p.Descripcion, &p.CategoriaID, &p.CategoriaNombre, &p.Precio, &p.StockActual, &p.StockMinimo, &p.Imagen, &p.CreatedAt, &p.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetProductos(page, limit int, search, categoriaID string, stockBajo bool) ([]model.Producto, int, error) {
	offset := (page - 1) * limit
	
	// Build query
	baseQuery := `FROM productos p LEFT JOIN categorias c ON p.categoria_id = c.id WHERE 1=1`
	var args []interface{}
	
	if search != "" {
		baseQuery += ` AND (p.nombre LIKE ? OR p.codigo LIKE ? OR p.codigo_barras LIKE ?)`
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}
	
	if categoriaID != "" {
		baseQuery += ` AND p.categoria_id = ?`
		args = append(args, categoriaID)
	}
	
	if stockBajo {
		baseQuery += ` AND p.stock_actual < p.stock_minimo`
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get data
	query := fmt.Sprintf(`SELECT p.id, p.codigo, p.codigo_barras, p.nombre, p.descripcion, p.categoria_id, c.nombre, p.precio, p.stock_actual, p.stock_minimo, p.imagen, p.created_at, p.updated_at %s LIMIT ? OFFSET ?`, baseQuery)
	args = append(args, limit, offset)
	
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var productos []model.Producto
	for rows.Next() {
		var p model.Producto
		if err := rows.Scan(&p.ID, &p.Codigo, &p.CodigoBarras, &p.Nombre, &p.Descripcion, &p.CategoriaID, &p.CategoriaNombre, &p.Precio, &p.StockActual, &p.StockMinimo, &p.Imagen, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		productos = append(productos, p)
	}
	return productos, total, nil
}

func GetProductosStockBajo() ([]model.Producto, error) {
	rows, err := DB.Query(
		`SELECT p.id, p.codigo, p.codigo_barras, p.nombre, p.descripcion, p.categoria_id, c.nombre, p.precio, p.stock_actual, p.stock_minimo, p.imagen, p.created_at, p.updated_at 
		 FROM productos p 
		 LEFT JOIN categorias c ON p.categoria_id = c.id 
		 WHERE p.stock_actual < p.stock_minimo 
		 ORDER BY p.stock_actual ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productos []model.Producto
	for rows.Next() {
		var p model.Producto
		if err := rows.Scan(&p.ID, &p.Codigo, &p.CodigoBarras, &p.Nombre, &p.Descripcion, &p.CategoriaID, &p.CategoriaNombre, &p.Precio, &p.StockActual, &p.StockMinimo, &p.Imagen, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		productos = append(productos, p)
	}
	return productos, nil
}

func UpdateProducto(id string, req *model.UpdateProductoRequest) (*model.Producto, error) {
	// Build dynamic update
	var sets []string
	var args []interface{}
	
	if req.Codigo != nil {
		sets = append(sets, "codigo = ?")
		args = append(args, *req.Codigo)
	}
	if req.CodigoBarras != nil {
		sets = append(sets, "codigo_barras = ?")
		args = append(args, *req.CodigoBarras)
	}
	if req.Nombre != nil {
		sets = append(sets, "nombre = ?")
		args = append(args, *req.Nombre)
	}
	if req.Descripcion != nil {
		sets = append(sets, "descripcion = ?")
		args = append(args, *req.Descripcion)
	}
	if req.CategoriaID != nil {
		sets = append(sets, "categoria_id = ?")
		args = append(args, *req.CategoriaID)
	}
	if req.Precio != nil {
		sets = append(sets, "precio = ?")
		args = append(args, *req.Precio)
	}
	if req.StockMinimo != nil {
		sets = append(sets, "stock_minimo = ?")
		args = append(args, *req.StockMinimo)
	}
	
	if len(sets) == 0 {
		return GetProductoByID(id)
	}
	
	sets = append(sets, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)
	
	query := "UPDATE productos SET " + strings.Join(sets, ", ") + " WHERE id = ?"
	_, err := DB.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	
	return GetProductoByID(id)
}

func UpdateStock(productoID string, cantidad int, esEntrada bool) error {
	operador := "+"
	if !esEntrada {
		operador = "-"
	}
	
	_, err := DB.Exec(
		fmt.Sprintf(`UPDATE productos SET stock_actual = stock_actual %s ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, operador),
		cantidad, productoID,
	)
	return err
}

func UpdateImagen(productoID, imagen string) (*model.Producto, error) {
	_, err := DB.Exec(
		`UPDATE productos SET imagen = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		imagen, productoID,
	)
	if err != nil {
		return nil, err
	}
	return GetProductoByID(productoID)
}

func DeleteProducto(id string) error {
	_, err := DB.Exec(`DELETE FROM productos WHERE id = ?`, id)
	return err
}

func CodigoExists(codigo string) (bool, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM productos WHERE codigo = ?`, codigo).Scan(&count)
	return count > 0, err
}