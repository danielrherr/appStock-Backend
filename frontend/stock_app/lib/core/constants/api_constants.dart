import 'package:flutter/foundation.dart';

class ApiConstants {
  static String get baseUrl {
    const envBaseUrl = String.fromEnvironment('API_BASE_URL');
    if (envBaseUrl.isNotEmpty) {
      return envBaseUrl;
    }

    if (kIsWeb) {
      return _resolveWebBaseUrl();
    }

    switch (defaultTargetPlatform) {
      case TargetPlatform.android:
        return 'http://10.0.2.2:8080/api/v1';
      default:
        return 'http://localhost:8080/api/v1';
    }
  }

  static String _resolveWebBaseUrl() {
    final host = Uri.base.host;
    final path = Uri.base.path;
    final isLocalHost = host == 'localhost' || host == '127.0.0.1';
    final isServedByBackend = path.startsWith('/web') || Uri.base.port == 8080;

    if (isServedByBackend || !isLocalHost) {
      return '/api/v1';
    }

    return 'http://$host:8080/api/v1';
  }

  // Auth
  static const String login = '/auth/login';
  static const String register = '/auth/register';

  // Productos
  static const String productos = '/productos';
  static const String stockBajo = '/productos/stock-bajo';

  // Categorías
  static const String categorias = '/categorias';

  // Movimientos
  static const String movimientos = '/movimientos';

  // Dashboard
  static const String dashboard = '/dashboard';

  // Devices
  static const String devices = '/devices';

  // Reportes
  static const String stockPorCategoria = '/reportes/stock-categoria';
  static const String movimientosPorFecha = '/reportes/movimientos-fecha';
  static const String dashboardAvanzado = '/reportes/dashboard-avanzado';
}
