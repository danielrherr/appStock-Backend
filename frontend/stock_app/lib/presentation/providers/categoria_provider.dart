import 'package:flutter/foundation.dart';
import '../../domain/entities/categoria.dart';
import '../../domain/repositories/categoria_repository.dart';
import '../../data/repositories/categoria_repository_impl.dart';

class CategoriaProvider extends ChangeNotifier {
  final CategoriaRepository _categoriaRepository = CategoriaRepositoryImpl();

  List<Categoria> _categorias = [];
  bool _isLoading = false;
  String? _error;

  List<Categoria> get categorias => _categorias;
  bool get isLoading => _isLoading;
  String? get error => _error;

  Future<void> fetchCategorias() async {
    _isLoading = true;
    notifyListeners();

    try {
      _categorias = await _categoriaRepository.getCategorias();
    } catch (e) {
      _error = e.toString();
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<bool> createCategoria(Categoria categoria) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final nueva = await _categoriaRepository.createCategoria(categoria);
      _categorias.add(nueva);
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

  Future<bool> updateCategoria(Categoria categoria) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final actualizada = await _categoriaRepository.updateCategoria(categoria);
      final index = _categorias.indexWhere((c) => c.id == categoria.id);
      if (index != -1) {
        _categorias[index] = actualizada;
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

  Future<bool> deleteCategoria(String id) async {
    try {
      await _categoriaRepository.deleteCategoria(id);
      _categorias.removeWhere((c) => c.id == id);
      notifyListeners();
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  String getCategoriaNombre(String? id) {
    if (id == null) return 'Sin categoría';
    final cat = _categorias.firstWhere(
      (c) => c.id == id,
      orElse: () => Categoria(id: '', nombre: 'Sin categoría'),
    );
    return cat.nombre;
  }

  void clearError() {
    _error = null;
    notifyListeners();
  }
}