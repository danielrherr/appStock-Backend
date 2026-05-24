package repository

import (
	"fmt"
	"strings"

	"github.com/stockapp/backend/internal/model"
)

func CreateProducto(req *model.CreateProductoRequest) (*model.Producto, error) {
	// PostgreSQL generates UUID automatically via gen_random_uuid()
	// We just need to pass NULL or omit the id column
	var id string
	err := DB.QueryRow(
		`INSERT INTO productos (codigo, codigo_barras, nombre, descripcion, categoria_id, precio, stock_actual, stock_minimo) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id`,
		req.Codigo, req.CodigoBarras, req.Nombre, req.Descripcion, req.CategoriaID, req.Precio, req.StockActual, req.StockMinimo,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return GetProductoByID(id)
}

func GetProductoByID(id string) (*model.Producto, error) {
	var p model.Producto
	err := DB.QueryRow(
		`SELECT p.id, p.codigo, p.codigo_barras, p.nombre, p.descripcion, p.categoria_id, COALESCE(c.nombre, ''), p.precio, p.stock_actual, p.stock_minimo, p.imagen, p.created_at, p.updated_at 
		 FROM productos p 
		 LEFT JOIN categorias c ON p.categoria_id = c.id 
		 WHERE p.id = $1`, id,
	).Scan(&p.ID, &p.Codigo, &p.CodigoBarras, &p.Nombre, &p.Descripcion, &p.CategoriaID, &p.CategoriaNombre, &p.Precio, &p.StockActual, &p.StockMinimo, &p.Imagen, &p.CreatedAt, &p.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetProductoByCodigo(codigo string) (*model.Producto, error) {
	var p model.Producto
	err := DB.QueryRow(
		`SELECT p.id, p.codigo, p.codigo_barras, p.nombre, p.descripcion, p.categoria_id, COALESCE(c.nombre, ''), p.precio, p.stock_actual, p.stock_minimo, p.imagen, p.created_at, p.updated_at 
		 FROM productos p 
		 LEFT JOIN categorias c ON p.categoria_id = c.id 
		 WHERE p.codigo = $1 OR p.codigo_barras = $1`, codigo,
	).Scan(&p.ID, &p.Codigo, &p.CodigoBarras, &p.Nombre, &p.Descripcion, &p.CategoriaID, &p.CategoriaNombre, &p.Precio, &p.StockActual, &p.StockMinimo, &p.Imagen, &p.CreatedAt, &p.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetProductos(page, limit int, search, categoriaID string, stockBajo bool) ([]model.Producto, int, error) {
	offset := (page - 1) * limit
	
	// Build query with dynamic placeholders
	baseQuery := `FROM productos p LEFT JOIN categorias c ON p.categoria_id = c.id WHERE 1=1`
	var args []interface{}
	argNum := 1
	
	if search != "" {
		baseQuery += fmt.Sprintf(` AND (p.nombre LIKE $%d OR p.codigo LIKE $%d OR p.codigo_barras LIKE $%d)`, argNum, argNum+1, argNum+2)
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
		argNum += 3
	}
	
	if categoriaID != "" {
		baseQuery += fmt.Sprintf(` AND p.categoria_id = $%d`, argNum)
		args = append(args, categoriaID)
		argNum++
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

	// Get data - append LIMIT/OFFSET at the end
	query := fmt.Sprintf(`SELECT p.id, p.codigo, p.codigo_barras, p.nombre, p.descripcion, p.categoria_id, COALESCE(c.nombre, ''), p.precio, p.stock_actual, p.stock_minimo, p.imagen, p.created_at, p.updated_at %s LIMIT $%d OFFSET $%d`, baseQuery, argNum, argNum+1)
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
		`SELECT p.id, p.codigo, p.codigo_barras, p.nombre, p.descripcion, p.categoria_id, COALESCE(c.nombre, ''), p.precio, p.stock_actual, p.stock_minimo, p.imagen, p.created_at, p.updated_at 
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
	argNum := 1
	
	if req.Codigo != nil {
		sets = append(sets, fmt.Sprintf("codigo = $%d", argNum))
		args = append(args, *req.Codigo)
		argNum++
	}
	if req.CodigoBarras != nil {
		sets = append(sets, fmt.Sprintf("codigo_barras = $%d", argNum))
		args = append(args, *req.CodigoBarras)
		argNum++
	}
	if req.Nombre != nil {
		sets = append(sets, fmt.Sprintf("nombre = $%d", argNum))
		args = append(args, *req.Nombre)
		argNum++
	}
	if req.Descripcion != nil {
		sets = append(sets, fmt.Sprintf("descripcion = $%d", argNum))
		args = append(args, *req.Descripcion)
		argNum++
	}
	if req.CategoriaID != nil {
		sets = append(sets, fmt.Sprintf("categoria_id = $%d", argNum))
		args = append(args, *req.CategoriaID)
		argNum++
	}
	if req.Precio != nil {
		sets = append(sets, fmt.Sprintf("precio = $%d", argNum))
		args = append(args, *req.Precio)
		argNum++
	}
	if req.StockActual != nil {
		sets = append(sets, fmt.Sprintf("stock_actual = $%d", argNum))
		args = append(args, *req.StockActual)
		argNum++
	}
	if req.StockMinimo != nil {
		sets = append(sets, fmt.Sprintf("stock_minimo = $%d", argNum))
		args = append(args, *req.StockMinimo)
		argNum++
	}
	
	if len(sets) == 0 {
		return GetProductoByID(id)
	}
	
	sets = append(sets, "updated_at = NOW()")
	args = append(args, id)
	
	query := fmt.Sprintf("UPDATE productos SET %s WHERE id = $%d", strings.Join(sets, ", "), argNum)
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
		fmt.Sprintf(`UPDATE productos SET stock_actual = stock_actual %s $1, updated_at = NOW() WHERE id = $2`, operador),
		cantidad, productoID,
	)
	return err
}

func UpdateImagen(productoID, imagen string) (*model.Producto, error) {
	_, err := DB.Exec(
		`UPDATE productos SET imagen = $1, updated_at = NOW() WHERE id = $2`,
		imagen, productoID,
	)
	if err != nil {
		return nil, err
	}
	return GetProductoByID(productoID)
}

func DeleteProducto(id string) error {
	_, err := DB.Exec(`DELETE FROM productos WHERE id = $1`, id)
	return err
}

func CodigoExists(codigo string) (bool, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM productos WHERE codigo = $1`, codigo).Scan(&count)
	return count > 0, err
}

func CodigoExistsForOtherProduct(codigo, id string) (bool, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM productos WHERE codigo = $1 AND id <> $2`, codigo, id).Scan(&count)
	return count > 0, err
}

func GetProductoCountByCategoria(categoriaID string) (int, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM productos WHERE categoria_id = $1`, categoriaID).Scan(&count)
	return count, err
}
