enum TipoMovimiento { entrada, salida }

class Movimiento {
  final String id;
  final String productoId;
  final String? productoNombre;
  final TipoMovimiento tipo;
  final int cantidad;
  final String? motivo;
  final String? usuarioId;
  final DateTime fecha;

  Movimiento({
    required this.id,
    required this.productoId,
    this.productoNombre,
    required this.tipo,
    required this.cantidad,
    this.motivo,
    this.usuarioId,
    required this.fecha,
  });

  bool get esEntrada => tipo == TipoMovimiento.entrada;
  bool get esSalida => tipo == TipoMovimiento.salida;

  Movimiento copyWith({
    String? id,
    String? productoId,
    String? productoNombre,
    TipoMovimiento? tipo,
    int? cantidad,
    String? motivo,
    String? usuarioId,
    DateTime? fecha,
  }) {
    return Movimiento(
      id: id ?? this.id,
      productoId: productoId ?? this.productoId,
      productoNombre: productoNombre ?? this.productoNombre,
      tipo: tipo ?? this.tipo,
      cantidad: cantidad ?? this.cantidad,
      motivo: motivo ?? this.motivo,
      usuarioId: usuarioId ?? this.usuarioId,
      fecha: fecha ?? this.fecha,
    );
  }
}