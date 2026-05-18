import '../../core/constants/api_constants.dart';
import '../../domain/entities/categoria.dart';
import '../../domain/repositories/categoria_repository.dart';
import '../api/api_client.dart';
import '../models/categoria_model.dart';

class CategoriaRepositoryImpl implements CategoriaRepository {
  final ApiClient _apiClient = ApiClient.instance;

  @override
  Future<List<Categoria>> getCategorias() async {
    final response = await _apiClient.get(ApiConstants.categorias);
    final List<dynamic> data = response.data['data'] ?? [];
    return data.map((json) => CategoriaModel.fromJson(json)).toList();
  }

  @override
  Future<Categoria> getCategoria(String id) async {
    final response = await _apiClient.get('${ApiConstants.categorias}/$id');
    return CategoriaModel.fromJson(response.data);
  }

  @override
  Future<Categoria> createCategoria(Categoria categoria) async {
    final model = CategoriaModel(
      id: '',
      nombre: categoria.nombre,
      descripcion: categoria.descripcion,
    );
    final response = await _apiClient.post(
      ApiConstants.categorias,
      data: model.toJson(),
    );
    return CategoriaModel.fromJson(response.data);
  }

  @override
  Future<Categoria> updateCategoria(Categoria categoria) async {
    final response = await _apiClient.put(
      '${ApiConstants.categorias}/${categoria.id}',
      data: {
        'nombre': categoria.nombre,
        'descripcion': categoria.descripcion,
      },
    );
    return CategoriaModel.fromJson(response.data);
  }

  @override
  Future<void> deleteCategoria(String id) async {
    await _apiClient.delete('${ApiConstants.categorias}/$id');
  }
}