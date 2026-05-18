import '../../domain/entities/movimiento.dart';

class MovimientoModel extends Movimiento {
  MovimientoModel({
    required super.id,
    required super.productoId,
    super.productoNombre,
    required super.tipo,
    required super.cantidad,
    super.motivo,
    super.usuarioId,
    required super.fecha,
  });

  factory MovimientoModel.fromJson(Map<String, dynamic> json) {
    return MovimientoModel(
      id: json['id'] ?? '',
      productoId: json['producto_id'] ?? '',
      productoNombre: json['producto_nombre'],
      tipo: json['tipo'] == 'entrada' 
          ? TipoMovimiento.entrada 
          : TipoMovimiento.salida,
      cantidad: json['cantidad'] ?? 0,
      motivo: json['motivo'],
      usuarioId: json['usuario_id'],
      fecha: json['fecha'] != null 
          ? DateTime.parse(json['fecha']) 
          : DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'producto_id': productoId,
      'tipo': tipo == TipoMovimiento.entrada ? 'entrada' : 'salida',
      'cantidad': cantidad,
      'motivo': motivo,
    };
  }
}