import 'dart:async';
import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:drift/drift.dart';
import 'package:dio/dio.dart';
import '../data/local/database.dart';
import '../core/constants/api_constants.dart';

class SyncService {
  static final SyncService _instance = SyncService._internal();
  factory SyncService() => _instance;
  SyncService._internal();

  final AppDatabase _db = AppDatabase();
  final Connectivity _connectivity = Connectivity();
  
  StreamSubscription<List<ConnectivityResult>>? _connectivitySubscription;
  bool _isOnline = true;
  bool _isSyncing = false;

  // Callbacks
  Function(bool)? onConnectivityChanged;
  Function()? onSyncComplete;

  // Getters
  bool get isOnline => _isOnline;
  bool get isSyncing => _isSyncing;

  Future<void> initialize() async {
    // Check initial connectivity
    final result = await _connectivity.checkConnectivity();
    _updateConnectivity(result);

    // Listen for changes
    _connectivitySubscription = _connectivity.onConnectivityChanged.listen(_updateConnectivity);
  }

  void _updateConnectivity(List<ConnectivityResult> result) {
    final wasOnline = _isOnline;
    _isOnline = result.isNotEmpty && !result.contains(ConnectivityResult.none);

    if (onConnectivityChanged != null) {
      onConnectivityChanged!(_isOnline);
    }

    // Trigger sync when coming online
    if (!wasOnline && _isOnline) {
      syncPendingData();
    }
  }

  Future<void> syncPendingData() async {
    if (_isSyncing || !_isOnline) return;

    _isSyncing = true;

    try {
      // Get pending movimientos
      final pendingMovimientos = await _db.getPendingMovimientos();

      if (pendingMovimientos.isEmpty) {
        _isSyncing = false;
        return;
      }

      // Get auth token
      // Note: You would get this from your auth provider/prefs
      final token = await _getToken();
      if (token == null) {
        _isSyncing = false;
        return;
      }

      final dio = Dio(BaseOptions(
        baseUrl: ApiConstants.baseUrl,
        headers: {
          'Authorization': 'Bearer $token',
          'Content-Type': 'application/json',
        },
      ));

      // Sync each movimiento
      for (final mov in pendingMovimientos) {
        try {
          await dio.post('/movimientos', data: {
            'producto_id': mov.productoId,
            'tipo': mov.tipo,
            'cantidad': mov.cantidad,
            'notas': mov.notas,
          });

          // Mark as synced
          await _db.markMovimientoSynced(mov.id);
        } catch (e) {
          print('Error syncing movimiento ${mov.id}: $e');
          // Continue with next
        }
      }

      // Clean up synced
      await _db.clearSyncedMovimientos();

      if (onSyncComplete != null) {
        onSyncComplete!();
      }
    } catch (e) {
      print('Sync error: $e');
    } finally {
      _isSyncing = false;
    }
  }

  Future<String?> _getToken() async {
    // Implement based on your auth implementation
    // Could use shared_preferences or auth provider
    return null;
  }

  Future<void> cacheProductos(List<Map<String, dynamic>> productos) async {
    final companions = productos.map((p) => ProductosCompanion(
      id: Value(p['id']),
      nombre: Value(p['nombre']),
      descripcion: Value(p['descripcion']),
      categoriaId: Value(p['categoria_id']),
      categoriaNombre: Value(p['categoria_nombre']),
      codigoBarra: Value(p['codigo_barra']),
      stock: Value(p['stock']),
      stockMinimo: Value(p['stock_minimo'] ?? 5),
      precio: Value((p['precio'] ?? 0).toDouble()),
      imagenUrl: Value(p['imagen_url']),
      activo: Value(p['activo'] ?? true),
      usuarioId: Value(p['usuario_id']),
      createdAt: Value(DateTime.parse(p['created_at'])),
      updatedAt: Value(DateTime.parse(p['updated_at'])),
    )).toList();

    await _db.insertProductos(companions);
  }

  Future<void> cacheCategorias(List<Map<String, dynamic>> categorias) async {
    final companions = categorias.map((c) => CategoriasCompanion(
      id: Value(c['id']),
      nombre: Value(c['nombre']),
      descripcion: Value(c['descripcion']),
      usuarioId: Value(c['usuario_id']),
      createdAt: Value(DateTime.parse(c['created_at'])),
      updatedAt: Value(DateTime.parse(c['updated_at'])),
    )).toList();

    await _db.insertCategorias(companions);
  }

  Future<void> addPendingMovimiento({
    required String productoId,
    required String productoNombre,
    required String tipo,
    required int cantidad,
    String? notas,
  }) async {
    await _db.addMovimientoOffline(MovimientosOfflineCompanion(
      productoId: Value(productoId),
      productoNombre: Value(productoNombre),
      tipo: Value(tipo),
      cantidad: Value(cantidad),
      notas: Value(notas),
      createdAt: Value(DateTime.now()),
      synced: const Value(false),
    ));

    // Try to sync immediately if online
    if (_isOnline) {
      syncPendingData();
    }
  }

  Future<int> getPendingCount() async {
    final pending = await _db.getPendingMovimientos();
    return pending.length;
  }

  void dispose() {
    _connectivitySubscription?.cancel();
  }
}