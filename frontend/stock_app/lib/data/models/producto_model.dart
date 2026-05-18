import '../../domain/entities/producto.dart';

class ProductoModel extends Producto {
  ProductoModel({
    required super.id,
    required super.codigo,
    super.codigoBarras,
    required super.nombre,
    super.descripcion,
    super.categoriaId,
    super.categoriaNombre,
    required super.precio,
    required super.stockActual,
    required super.stockMinimo,
    super.imagen,
    super.createdAt,
    super.updatedAt,
  });

  factory ProductoModel.fromJson(Map<String, dynamic> json) {
    return ProductoModel(
      id: json['id'] ?? '',
      codigo: json['codigo'] ?? '',
      codigoBarras: json['codigo_barras'],
      nombre: json['nombre'] ?? '',
      descripcion: json['descripcion'],
      categoriaId: json['categoria_id'],
      categoriaNombre: json['categoria_nombre'],
      precio: (json['precio'] ?? 0).toDouble(),
      stockActual: json['stock_actual'] ?? 0,
      stockMinimo: json['stock_minimo'] ?? 0,
      imagen: json['imagen'],
      createdAt: json['created_at'] != null 
          ? DateTime.parse(json['created_at']) 
          : null,
      updatedAt: json['updated_at'] != null 
          ? DateTime.parse(json['updated_at']) 
          : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'codigo': codigo,
      'codigo_barras': codigoBarras,
      'nombre': nombre,
      'descripcion': descripcion,
      'categoria_id': categoriaId,
      'precio': precio,
      'stock_actual': stockActual,
      'stock_minimo': stockMinimo,
      'imagen': imagen,
    };
  }
}