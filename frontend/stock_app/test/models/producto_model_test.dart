import 'package:flutter_test/flutter_test.dart';
import 'package:stock_app/data/models/producto_model.dart';
import 'package:stock_app/domain/entities/producto.dart';

void main() {
  group('ProductoModel', () {
    const testJson = {
      'id': 'test-uuid-123',
      'nombre': 'Test Producto',
      'descripcion': 'Descripción del producto',
      'categoria_id': 'cat-uuid-456',
      'categoria_nombre': 'Electrónica',
      'codigo_barra': '1234567890123',
      'stock': 50,
      'stock_minimo': 10,
      'precio': 99.99,
      'imagen_url': 'https://example.com/image.jpg',
      'activo': true,
      'usuario_id': 'user-uuid-789',
      'created_at': '2024-01-15T10:30:00Z',
      'updated_at': '2024-01-16T14:20:00Z',
    };

    test('should parse from JSON correctly', () {
      final producto = ProductoModel.fromJson(testJson);

      expect(producto.id, equals('test-uuid-123'));
      expect(producto.nombre, equals('Test Producto'));
      expect(producto.descripcion, equals('Descripción del producto'));
      expect(producto.categoriaId, equals('cat-uuid-456'));
      expect(producto.categoriaNombre, equals('Electrónica'));
      expect(producto.codigoBarra, equals('1234567890123'));
      expect(producto.stock, equals(50));
      expect(producto.stockMinimo, equals(10));
      expect(producto.precio, equals(99.99));
      expect(producto.imagenUrl, equals('https://example.com/image.jpg'));
      expect(producto.activo, isTrue);
      expect(producto.usuarioId, equals('user-uuid-789'));
    });

    test('should convert to JSON correctly', () {
      final producto = ProductoModel.fromJson(testJson);
      final json = producto.toJson();

      expect(json['id'], equals(testJson['id']));
      expect(json['nombre'], equals(testJson['nombre']));
      expect(json['stock'], equals(testJson['stock']));
    });

    test('should handle null optional fields', () {
      const minimalJson = {
        'id': 'test-id',
        'nombre': 'Minimal Product',
        'categoria_id': 'cat-id',
        'stock': 0,
      };

      final producto = ProductoModel.fromJson(minimalJson);

      expect(producto.descripcion, isNull);
      expect(producto.codigoBarra, isNull);
      expect(producto.imagenUrl, isNull);
    });

    test('should convert to entity correctly', () {
      final modelo = ProductoModel.fromJson(testJson);
      final entity = modelo.toEntity();

      expect(entity, isA<Producto>());
      expect(entity.id, equals(modelo.id));
      expect(entity.nombre, equals(modelo.nombre));
      expect(entity.stock, equals(modelo.stock));
    });
  });

  group('ProductoModel - Stock Alerts', () {
    test('should detect stock bajo when stock < stock_minimo', () {
      const lowStockJson = {
        'id': 'test-id',
        'nombre': 'Low Stock Product',
        'categoria_id': 'cat-id',
        'stock': 5,
        'stock_minimo': 10,
      };

      final producto = ProductoModel.fromJson(lowStockJson);
      expect(producto.tieneStockBajo, isTrue);
    });

    test('should NOT detect stock bajo when stock >= stock_minimo', () {
      const normalStockJson = {
        'id': 'test-id',
        'nombre': 'Normal Stock Product',
        'categoria_id': 'cat-id',
        'stock': 15,
        'stock_minimo': 10,
      };

      final producto = ProductoModel.fromJson(normalStockJson);
      expect(producto.tieneStockBajo, isFalse);
    });

    test('should handle zero stock_minimo', () {
      const zeroMinJson = {
        'id': 'test-id',
        'nombre': 'Zero Min Product',
        'categoria_id': 'cat-id',
        'stock': 0,
        'stock_minimo': 0,
      };

      final producto = ProductoModel.fromJson(zeroMinJson);
      expect(producto.tieneStockBajo, isFalse);
    });
  });
}