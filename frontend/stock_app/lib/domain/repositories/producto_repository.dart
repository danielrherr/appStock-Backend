import '../entities/producto.dart';

abstract class ProductoRepository {
  Future<List<Producto>> getProductos({int page = 1, String? search, String? categoriaId, bool? stockBajo});
  Future<Producto> getProducto(String id);
  Future<Producto> createProducto(Producto producto);
  Future<Producto> updateProducto(Producto producto);
  Future<void> deleteProducto(String id);
  Future<Producto?> searchByBarcode(String barcode);
  Future<String> uploadImagen(String productoId, String filePath);
}