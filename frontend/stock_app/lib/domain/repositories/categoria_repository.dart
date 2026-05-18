import '../entities/categoria.dart';

abstract class CategoriaRepository {
  Future<List<Categoria>> getCategorias();
  Future<Categoria> getCategoria(String id);
  Future<Categoria> createCategoria(Categoria categoria);
  Future<Categoria> updateCategoria(Categoria categoria);
  Future<void> deleteCategoria(String id);
}