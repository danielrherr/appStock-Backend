import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../core/theme/app_colors.dart';
import '../../providers/producto_provider.dart';
import 'producto_form_screen.dart';
import 'producto_detail_screen.dart';

class ProductoListScreen extends StatefulWidget {
  const ProductoListScreen({super.key});

  @override
  State<ProductoListScreen> createState() => _ProductoListScreenState();
}

class _ProductoListScreenState extends State<ProductoListScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<ProductoProvider>().fetchProductos(refresh: true);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Productos'),
        actions: [
          IconButton(
            icon: const Icon(Icons.search),
            onPressed: () {
              // Implementar búsqueda
            },
          ),
          IconButton(
            icon: const Icon(Icons.qr_code_scanner),
            onPressed: () {
              // Implementar scanner
            },
          ),
        ],
      ),
      body: Consumer<ProductoProvider>(
        builder: (context, provider, _) {
          if (provider.isLoading && provider.productos.isEmpty) {
            return const Center(child: CircularProgressIndicator());
          }

          if (provider.productos.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(
                    Icons.inventory_2_outlined,
                    size: 64,
                    color: AppColors.textSecondary,
                  ),
                  const SizedBox(height: 16),
                  const Text(
                    'No hay productos',
                    style: TextStyle(
                      fontSize: 18,
                      color: AppColors.textSecondary,
                    ),
                  ),
                  const SizedBox(height: 24),
                  ElevatedButton.icon(
                    onPressed: () {
                      Navigator.of(context).push(
                        MaterialPageRoute(
                          builder: (_) => const ProductoFormScreen(),
                        ),
                      );
                    },
                    icon: const Icon(Icons.add),
                    label: const Text('Agregar Producto'),
                  ),
                ],
              ),
            );
          }

          return RefreshIndicator(
            onRefresh: () => provider.fetchProductos(refresh: true),
            child: ListView.builder(
              padding: const EdgeInsets.all(8),
              itemCount: provider.productos.length,
              itemBuilder: (context, index) {
                final producto = provider.productos[index];
                return Card(
                  child: ListTile(
                    leading: producto.imagen != null
                        ? Image.network(
                            producto.imagen!,
                            width: 48,
                            height: 48,
                            fit: BoxFit.cover,
                            errorBuilder: (_, __, ___) => const Icon(
                              Icons.inventory_2,
                              color: AppColors.primary,
                            ),
                          )
                        : const Icon(
                            Icons.inventory_2,
                            color: AppColors.primary,
                          ),
                    title: Text(producto.nombre),
                    subtitle: Text('Stock: ${producto.stockActual}'),
                    trailing: producto.tieneStockBajo
                        ? const Icon(Icons.warning, color: AppColors.error)
                        : null,
                    onTap: () {
                      Navigator.of(context).push(
                        MaterialPageRoute(
                          builder: (_) => ProductoDetailScreen(producto: producto),
                        ),
                      );
                    },
                  ),
                );
              },
            ),
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          Navigator.of(context).push(
            MaterialPageRoute(
              builder: (_) => const ProductoFormScreen(),
            ),
          );
        },
        child: const Icon(Icons.add),
      ),
    );
  }
}