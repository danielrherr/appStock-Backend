import 'producto.dart';

class ItemListaCompra {
  final String id;
  final Producto producto;
  final int cantidad;
  final bool comprado;
  final DateTime createdAt;

  ItemListaCompra({
    required this.id,
    required this.producto,
    required this.cantidad,
    this.comprado = false,
    DateTime? createdAt,
  }) : createdAt = createdAt ?? DateTime.now();

  ItemListaCompra copyWith({
    String? id,
    Producto? producto,
    int? cantidad,
    bool? comprado,
    DateTime? createdAt,
  }) {
    return ItemListaCompra(
      id: id ?? this.id,
      producto: producto ?? this.producto,
      cantidad: cantidad ?? this.cantidad,
      comprado: comprado ?? this.comprado,
      createdAt: createdAt ?? this.createdAt,
    );
  }
}