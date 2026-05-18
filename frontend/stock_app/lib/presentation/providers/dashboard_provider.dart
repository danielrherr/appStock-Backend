import 'package:flutter/foundation.dart';
import '../../core/constants/api_constants.dart';
import '../../data/api/api_client.dart';

class DashboardProvider extends ChangeNotifier {
  final ApiClient _apiClient = ApiClient.instance;

  int _totalProductos = 0;
  int _stockBajo = 0;
  int _movimientosHoy = 0;
  double _valorTotalStock = 0;
  bool _isLoading = false;
  String? _error;

  int get totalProductos => _totalProductos;
  int get stockBajo => _stockBajo;
  int get movimientosHoy => _movimientosHoy;
  double get valorTotalStock => _valorTotalStock;
  bool get isLoading => _isLoading;
  String? get error => _error;

  Future<void> fetchDashboard() async {
    _isLoading = true;
    notifyListeners();

    try {
      final response = await _apiClient.get(ApiConstants.dashboard);
      final data = response.data;
      
      _totalProductos = data['total_productos'] ?? 0;
      _stockBajo = data['productos_stock_bajo'] ?? 0;
      _movimientosHoy = data['movimientos_hoy'] ?? 0;
      _valorTotalStock = (data['valor_total_stock'] ?? 0).toDouble();
    } catch (e) {
      _error = e.toString();
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  void clearError() {
    _error = null;
    notifyListeners();
  }
}