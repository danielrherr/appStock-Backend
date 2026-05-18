import 'package:flutter/foundation.dart';
import '../../data/api/api_client.dart';
import '../../core/constants/api_constants.dart';

class ReportProvider extends ChangeNotifier {
  bool _isLoading = false;
  String? _error;

  // Reportes
  List<dynamic>? _stockPorCategoria;
  List<dynamic>? _movimientosPorFecha;
  Map<String, dynamic>? _dashboardAvanzado;

  // Getters
  bool get isLoading => _isLoading;
  String? get error => _error;
  List<dynamic>? get stockPorCategoria => _stockPorCategoria;
  List<dynamic>? get movimientosPorFecha => _movimientosPorFecha;
  Map<String, dynamic>? get dashboardAvanzado => _dashboardAvanzado;

  Future<void> loadStockPorCategoria() async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _stockPorCategoria = await ApiClient.get(ApiConstants.stockPorCategoria);
    } catch (e) {
      _error = e.toString();
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> loadMovimientosPorFecha({String? desde, String? hasta}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final queryParams = <String, String>{};
      if (desde != null) queryParams['desde'] = desde;
      if (hasta != null) queryParams['hasta'] = hasta;

      _movimientosPorFecha = await ApiClient.get(
        ApiConstants.movimientosPorFecha,
        queryParams: queryParams,
      );
    } catch (e) {
      _error = e.toString();
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> loadDashboardAvanzado() async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _dashboardAvanzado = await ApiClient.get(ApiConstants.dashboardAvanzado);
    } catch (e) {
      _error = e.toString();
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  // Generar CSV de productos
  String generateCSV(List<dynamic> productos) {
    final buffer = StringBuffer();
    buffer.writeln('ID,Nombre,Categoría,Stock,Stock Mínimo,Precio,Activo');

    for (final p in productos) {
      buffer.writeln(
        '${p['id']},'
        '"${p['nombre']}",'
        '"${p['categoria_nombre'] ?? ''}",'
        '${p['stock']},'
        '${p['stock_minimo']},'
        '${p['precio']},'
        '${p['activo']}'
      );
    }

    return buffer.toString();
  }

  // Generar CSV de movimientos
  String generateMovimientosCSV(List<dynamic> movimientos) {
    final buffer = StringBuffer();
    buffer.writeln('Fecha,Producto,Tipo,Cantidad,Notas');

    for (final m in movimientos) {
      buffer.writeln(
        '${m['fecha']},'
        '"${m['producto_nombre'] ?? ''}",'
        '${m['tipo']},'
        '${m['cantidad']},'
        '"${m['notas'] ?? ''}"'
      );
    }

    return buffer.toString();
  }
}