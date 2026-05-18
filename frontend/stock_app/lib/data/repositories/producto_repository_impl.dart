import '../../core/constants/api_constants.dart';
import '../../domain/entities/producto.dart';
import '../../domain/repositories/producto_repository.dart';
import '../api/api_client.dart';
import '../models/producto_model.dart';

class ProductoRepositoryImpl implements ProductoRepository {
  final ApiClient _apiClient = ApiClient.instance;

  @override
  Future<List<Producto>> getProductos({
    int page = 1, 
    String? search, 
    String? categoriaId, 
    bool? stockBajo,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'limit': 20,
      if (search != null && search.isNotEmpty) 'search': search,
      if (categoriaId != null) 'categoria_id': categoriaId,
      if (stockBajo == true) 'stock_bajo': 'true',
    };

    final response = await _apiClient.get(
      ApiConstants.productos,
      queryParams: queryParams,
    );

    final List<dynamic> data = response.data['data'] ?? [];
    return data.map((json) => ProductoModel.fromJson(json)).toList();
  }

  @override
  Future<Producto> getProducto(String id) async {
    final response = await _apiClient.get('${ApiConstants.productos}/$id');
    return ProductoModel.fromJson(response.data);
  }

  @override
  Future<Producto> createProducto(Producto producto) async {
    final model = ProductoModel(
      id: '',
      codigo: producto.codigo,
      codigoBarras: producto.codigoBarras,
      nombre: producto.nombre,
      descripcion: producto.descripcion,
      categoriaId: producto.categoriaId,
      precio: producto.precio,
      stockActual: producto.stockActual,
      stockMinimo: producto.stockMinimo,
    );

    final response = await _apiClient.post(
      ApiConstants.productos,
      data: model.toJson(),
    );
    return ProductoModel.fromJson(response.data);
  }

  @override
  Future<Producto> updateProducto(Producto producto) async {
    final response = await _apiClient.put(
      '${ApiConstants.productos}/${producto.id}',
      data: {
        'nombre': producto.nombre,
        'descripcion': producto.descripcion,
        'categoria_id': producto.categoriaId,
        'precio': producto.precio,
        'stock_minimo': producto.stockMinimo,
        if (producto.codigoBarras != null) 'codigo_barras': producto.codigoBarras,
      },
    );
    return ProductoModel.fromJson(response.data);
  }

  @override
  Future<void> deleteProducto(String id) async {
    await _apiClient.delete('${ApiConstants.productos}/$id');
  }

  @override
  Future<Producto?> searchByBarcode(String barcode) async {
    try {
      final response = await _apiClient.get(
        '${ApiConstants.productos}/barcode/$barcode',
      );
      return ProductoModel.fromJson(response.data);
    } catch (e) {
      return null;
    }
  }

  @override
  Future<String> uploadImagen(String productoId, String filePath) async {
    final response = await _apiClient.uploadFile(
      '${ApiConstants.productos}/$productoId/imagen',
      filePath,
      'imagen',
    );
    return response.data['imagen'] ?? '';
  }
}