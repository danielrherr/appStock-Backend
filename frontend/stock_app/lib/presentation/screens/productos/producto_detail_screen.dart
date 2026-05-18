import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../core/theme/app_colors.dart';
import '../../../domain/entities/producto.dart';
import '../../providers/producto_provider.dart';
import '../movimientos/movimiento_form_screen.dart';
import 'producto_form_screen.dart';

class ProductoDetailScreen extends StatelessWidget {
  final Producto producto;

  const ProductoDetailScreen({super.key, required this.producto});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(producto.nombre),
        actions: [
          IconButton(
            icon: const Icon(Icons.edit),
            onPressed: () => _edit(context),
          ),
          IconButton(
            icon: const Icon(Icons.delete),
            onPressed: () => _delete(context),
          ),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Imagen
            if (producto.imagen != null)
              Container(
                width: double.infinity,
                height: 200,
                decoration: BoxDecoration(
                  color: AppColors.background,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Image.network(
                  producto.imagen!,
                  fit: BoxFit.cover,
                  errorBuilder: (_, __, ___) => const Icon(
                    Icons.inventory_2,
                    size: 64,
                    color: AppColors.textSecondary,
                  ),
                ),
              )
            else
              Container(
                width: double.infinity,
                height: 200,
                decoration: BoxDecoration(
                  color: AppColors.background,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: const Icon(
                  Icons.inventory_2,
                  size: 64,
                  color: AppColors.textSecondary,
                ),
              ),
            const SizedBox(height: 24),

            // Stock card
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceAround,
                  children: [
                    _buildStockInfo(
                      'Stock Actual',
                      producto.stockActual.toString(),
                      _getStockColor(producto.stockActual, producto.stockMinimo),
                    ),
                    Container(
                      width: 1,
                      height: 40,
                      color: AppColors.border,
                    ),
                    _buildStockInfo(
                      'Stock Mínimo',
                      producto.stockMinimo.toString(),
                      AppColors.textSecondary,
                    ),
                    Container(
                      width: 1,
                      height: 40,
                      color: AppColors.border,
                    ),
                    _buildStockInfo(
                      'Precio',
                      '\$${producto.precio.toStringAsFixed(2)}',
                      AppColors.primary,
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),

            // Código
            _buildInfoRow('Código', producto.codigo),
            if (producto.codigoBarras != null)
              _buildInfoRow('Código de barras', producto.codigoBarras!),
            if (producto.categoriaNombre != null)
              _buildInfoRow('Categoría', producto.categoriaNombre!),
            if (producto.descripcion != null && producto.descripcion!.isNotEmpty)
              _buildInfoRow('Descripción', producto.descripcion!),
            _buildInfoRow(
              'Creado',
              producto.createdAt?.toString().split('.').first ?? '-',
            ),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => _registrarMovimiento(context),
        icon: const Icon(Icons.swap_horiz),
        label: const Text('Movimiento'),
      ),
    );
  }

  Widget _buildStockInfo(String label, String value, Color color) {
    return Column(
      children: [
        Text(
          value,
          style: TextStyle(
            fontSize: 20,
            fontWeight: FontWeight.bold,
            color: color,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: const TextStyle(
            fontSize: 12,
            color: AppColors.textSecondary,
          ),
        ),
      ],
    );
  }

  Widget _buildInfoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 120,
            child: Text(
              label,
              style: const TextStyle(
                color: AppColors.textSecondary,
                fontSize: 14,
              ),
            ),
          ),
          Expanded(
            child: Text(
              value,
              style: const TextStyle(
                fontSize: 14,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Color _getStockColor(int stock, int min) {
    if (stock < min) return AppColors.error;
    if (stock < min * 1.5) return AppColors.warning;
    return AppColors.success;
  }

  void _edit(BuildContext context) async {
    final result = await Navigator.of(context).push<bool>(
      MaterialPageRoute(
        builder: (_) => ProductoFormScreen(producto: producto),
      ),
    );
    if (result == true && context.mounted) {
      Navigator.of(context).pop();
    }
  }

  void _delete(BuildContext context) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Confirmar eliminación'),
        content: Text('¿Estás seguro de eliminar "${producto.nombre}"?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Cancelar'),
          ),
          TextButton(
            onPressed: () => Navigator.pop(context, true),
            style: TextButton.styleFrom(foregroundColor: AppColors.error),
            child: const Text('Eliminar'),
          ),
        ],
      ),
    );

    if (confirmed == true && context.mounted) {
      final success = await context.read<ProductoProvider>().deleteProducto(producto.id);
      if (success && context.mounted) {
        Navigator.of(context).pop();
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Producto eliminado'),
            backgroundColor: AppColors.success,
          ),
        );
      }
    }
  }

  void _registrarMovimiento(BuildContext context) {
    Navigator.of(context).push(
      MaterialPageRoute(
        builder: (_) => MovimientoFormScreen(producto: producto),
      ),
    );
  }
}