class Producto {
  final String id;
  final String codigo;
  final String? codigoBarras;
  final String nombre;
  final String? descripcion;
  final String? categoriaId;
  final String? categoriaNombre;
  final double precio;
  final int stockActual;
  final int stockMinimo;
  final String? imagen;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  Producto({
    required this.id,
    required this.codigo,
    this.codigoBarras,
    required this.nombre,
    this.descripcion,
    this.categoriaId,
    this.categoriaNombre,
    required this.precio,
    required this.stockActual,
    required this.stockMinimo,
    this.imagen,
    this.createdAt,
    this.updatedAt,
  });

  bool get tieneStockBajo => stockActual < stockMinimo;
  
  Producto copyWith({
    String? id,
    String? codigo,
    String? codigoBarras,
    String? nombre,
    String? descripcion,
    String? categoriaId,
    String? categoriaNombre,
    double? precio,
    int? stockActual,
    int? stockMinimo,
    String? imagen,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Producto(
      id: id ?? this.id,
      codigo: codigo ?? this.codigo,
      codigoBarras: codigoBarras ?? this.codigoBarras,
      nombre: nombre ?? this.nombre,
      descripcion: descripcion ?? this.descripcion,
      categoriaId: categoriaId ?? this.categoriaId,
      categoriaNombre: categoriaNombre ?? this.categoriaNombre,
      precio: precio ?? this.precio,
      stockActual: stockActual ?? this.stockActual,
      stockMinimo: stockMinimo ?? this.stockMinimo,
      imagen: imagen ?? this.imagen,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}