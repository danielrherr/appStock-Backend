import 'package:flutter/foundation.dart';
import '../../domain/entities/producto.dart';
import '../../domain/repositories/producto_repository.dart';
import '../../data/repositories/producto_repository_impl.dart';

class ProductoProvider extends ChangeNotifier {
  final ProductoRepository _productoRepository = ProductoRepositoryImpl();

  List<Producto> _productos = [];
  List<Producto> _stockBajo = [];
  bool _isLoading = false;
  String? _error;
  int _currentPage = 1;
  bool _hasMore = true;

  List<Producto> get productos => _productos;
  List<Producto> get stockBajo => _stockBajo;
  bool get isLoading => _isLoading;
  String? get error => _error;
  bool get hasMore => _hasMore;

  Future<void> fetchProductos({
    String? search,
    String? categoriaId,
    bool refresh = false,
  }) async {
    if (refresh) {
      _currentPage = 1;
      _hasMore = true;
      _productos = [];
    }

    if (!_hasMore || _isLoading) return;

    _isLoading = true;
    notifyListeners();

    try {
      final nuevos = await _productoRepository.getProductos(
        page: _currentPage,
        search: search,
        categoriaId: categoriaId,
      );

      if (nuevos.length < 20) {
        _hasMore = false;
      }

      _productos = [..._productos, ...nuevos];
      _currentPage++;
    } catch (e) {
      _error = e.toString();
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> fetchStockBajo() async {
    try {
      _stockBajo = await _productoRepository.getProductos(stockBajo: true);
      notifyListeners();
    } catch (e) {
      _error = e.toString();
    }
  }

  Future<bool> createProducto(Producto producto) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final nuevo = await _productoRepository.createProducto(producto);
      _productos.insert(0, nuevo);
      _isLoading = false;
      notifyListeners();
      return true;
    } catch (e) {
      _error = e.toString();
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }

  Future<bool> updateProducto(Producto producto) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final actualizado = await _productoRepository.updateProducto(producto);
      final index = _productos.indexWhere((p) => p.id == producto.id);
      if (index != -1) {
        _productos[index] = actualizado;
      }
      _isLoading = false;
      notifyListeners();
      return true;
    } catch (e) {
      _error = e.toString();
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }

  Future<bool> deleteProducto(String id) async {
    try {
      await _productoRepository.deleteProducto(id);
      _productos.removeWhere((p) => p.id == id);
      notifyListeners();
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<Producto?> searchByBarcode(String barcode) async {
    return await _productoRepository.searchByBarcode(barcode);
  }

  Future<String?> uploadImagen(String productoId, String filePath) async {
    try {
      return await _productoRepository.uploadImagen(productoId, filePath);
    } catch (e) {
      _error = e.toString();
      return null;
    }
  }

  void clearError() {
    _error = null;
    notifyListeners();
  }
}