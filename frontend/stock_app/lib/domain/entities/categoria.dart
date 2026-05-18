class Categoria {
  final String id;
  final String nombre;
  final String? descripcion;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  Categoria({
    required this.id,
    required this.nombre,
    this.descripcion,
    this.createdAt,
    this.updatedAt,
  });

  Categoria copyWith({
    String? id,
    String? nombre,
    String? descripcion,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Categoria(
      id: id ?? this.id,
      nombre: nombre ?? this.nombre,
      descripcion: descripcion ?? this.descripcion,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}