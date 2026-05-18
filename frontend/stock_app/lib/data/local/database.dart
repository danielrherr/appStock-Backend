import 'dart:io';
import 'package:drift/drift.dart';
import 'package:drift/native.dart';
import 'package:path_provider/path_provider.dart';
import 'package:path/path.dart' as p;

part 'database.g.dart';

// Tablas
class Productos extends Table {
  TextColumn get id => text()();
  TextColumn get nombre => text()();
  TextColumn get descripcion => text().nullable()();
  TextColumn get categoriaId => text()();
  TextColumn get categoriaNombre => text().nullable()();
  TextColumn get codigoBarra => text().nullable()();
  IntColumn get stock => integer().withDefault(const Constant(0))();
  IntColumn get stockMinimo => integer().withDefault(const Constant(5))();
  RealColumn get precio => real().withDefault(const Constant(0.0))();
  TextColumn get imagenUrl => text().nullable()();
  BoolColumn get activo => boolean().withDefault(const Constant(true))();
  TextColumn get usuarioId => text()();
  DateTimeColumn get createdAt => dateTime()();
  DateTimeColumn get updatedAt => dateTime()();

  @override
  Set<Column> get primaryKey => {id};
}

class Categorias extends Table {
  TextColumn get id => text()();
  TextColumn get nombre => text()();
  TextColumn get descripcion => text().nullable()();
  TextColumn get usuarioId => text()();
  DateTimeColumn get createdAt => dateTime()();
  DateTimeColumn get updatedAt => dateTime()();

  @override
  Set<Column> get primaryKey => {id};
}

class MovimientosOffline extends Table {
  IntColumn get id => integer().autoIncrement()();
  TextColumn get productoId => text()();
  TextColumn get productoNombre => text()();
  TextColumn get tipo => text()(); // entrada o salida
  IntColumn get cantidad => integer()();
  TextColumn get notas => text().nullable()();
  DateTimeColumn get createdAt => dateTime()();
  BoolColumn get synced => boolean().withDefault(const Constant(false))();
}

@DriftDatabase(tables: [Productos, Categorias, MovimientosOffline])
class AppDatabase extends _$AppDatabase {
  AppDatabase() : super(_openConnection());

  @override
  int get schemaVersion => 1;

  // Productos
  Future<List<Producto>> getAllProductos() => select(productos).get();

  Stream<List<Producto>> watchProductos() => select(productos).watch();

  Future<Producto?> getProductoById(String id) =>
      (select(productos)..where((t) => t.id.equals(id))).getSingleOrNull();

  Future<void> insertProducto(ProductosCompanion producto) =>
      into(productos).insert(producto);

  Future<void> insertProductos(List<ProductosCompanion> productosList) async {
    await batch((batch) {
      batch.insertAll(productos, productosList, mode: InsertMode.insertOrReplace);
    });
  }

  Future<void> updateProducto(ProductosCompanion producto) =>
      (update(productos)..where((t) => t.id.equals(producto.id.value)))
          .write(producto);

  Future<void> deleteProducto(String id) =>
      (delete(productos)..where((t) => t.id.equals(id))).go();

  Future<void> clearProductos() => delete(productos).go();

  // Categorías
  Future<List<Categoria>> getAllCategorias() => select(categorias).get();

  Stream<List<Categoria>> watchCategorias() => select(categorias).watch();

  Future<void> insertCategoria(CategoriasCompanion categoria) =>
      into(categorias).insert(categoria);

  Future<void> insertCategorias(List<CategoriasCompanion> categoriasList) async {
    await batch((batch) {
      batch.insertAll(categorias, categoriasList, mode: InsertMode.insertOrReplace);
    });
  }

  Future<void> clearCategorias() => delete(categorias).go();

  // Movimientos Offline (Queue)
  Future<List<MovimientosOfflineData>> getPendingMovimientos() =>
      (select(movimientosOffline)..where((t) => t.synced.equals(false))).get();

  Stream<List<MovimientosOfflineData>> watchPendingMovimientos() =>
      (select(movimientosOffline)..where((t) => t.synced.equals(false))).watch();

  Future<void> addMovimientoOffline(MovimientosOfflineCompanion movimiento) =>
      into(movimientosOffline).insert(movimiento);

  Future<void> markMovimientoSynced(int id) =>
      (update(movimientosOffline)..where((t) => t.id.equals(id)))
          .write(const MovimientosOfflineCompanion(synced: Value(true)));

  Future<void> clearSyncedMovimientos() =>
      (delete(movimientosOffline)..where((t) => t.synced.equals(true))).go();
}

LazyDatabase _openConnection() {
  return LazyDatabase(() async {
    final dbFolder = await getApplicationDocumentsDirectory();
    final file = File(p.join(dbFolder.path, 'stockapp.db'));
    return NativeDatabase.createInBackground(file);
  });
}