import 'package:flutter/material.dart';
import '../../../core/theme/app_colors.dart';
import '../../../domain/entities/producto.dart';
import '../producto/producto_card.dart';

class StockBajoList extends StatelessWidget {
  final List<Producto> productos;
  final VoidCallback? onViewAll;
  final Function(Producto)? onProductoTap;

  const StockBajoList({
    super.key,
    required this.productos,
    this.onViewAll,
    this.onProductoTap,
  });

  @override
  Widget build(BuildContext context) {
    if (productos.isEmpty) {
      return const SizedBox.shrink();
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Row(
              children: [
                const Icon(Icons.warning, color: AppColors.error, size: 20),
                const SizedBox(width: 8),
                const Text(
                  'Stock Bajo',
                  style: TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 16,
                  ),
                ),
              ],
            ),
            if (onViewAll != null)
              TextButton(
                onPressed: onViewAll,
                child: const Text('Ver todos'),
              ),
          ],
        ),
        const SizedBox(height: 8),
        ...productos.take(3).map((producto) => Padding(
          padding: const EdgeInsets.only(bottom: 8),
          child: ProductoCard(
            producto: producto,
            onTap: onProductoTap != null 
                ? () => onProductoTap!(producto)
                : null,
          ),
        )),
        if (productos.length > 3)
          Padding(
            padding: const EdgeInsets.only(top: 8),
            child: Text(
              '+${productos.length - 3} más',
              style: const TextStyle(
                color: AppColors.textSecondary,
                fontSize: 12,
              ),
            ),
          ),
      ],
    );
  }
}