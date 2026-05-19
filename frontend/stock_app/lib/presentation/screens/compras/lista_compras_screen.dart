import 'package:flutter/material.dart';
import '../../../core/theme/app_colors.dart';
import '../../../domain/entities/producto.dart';
import '../../../domain/entities/lista_compra.dart';
import '../scanner/barcode_scanner_screen.dart';

class ListaComprasScreen extends StatefulWidget {
  const ListaComprasScreen({super.key});

  @override
  State<ListaComprasScreen> createState() => _ListaComprasScreenState();
}

class _ListaComprasScreenState extends State<ListaComprasScreen> {
  final List<ItemListaCompra> _items = [];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Lista de Compras'),
        actions: [
          if (_items.isNotEmpty)
            IconButton(
              icon: const Icon(Icons.delete_sweep),
              onPressed: _limpiarLista,
              tooltip: 'Limpiar lista',
            ),
        ],
      ),
      body: _items.isEmpty ? _buildEmpty() : _buildList(),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: _agregarProducto,
        icon: const Icon(Icons.add),
        label: const Text('Agregar'),
      ),
    );
  }

  Widget _buildEmpty() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.shopping_cart_outlined,
            size: 80,
            color: AppColors.textSecondary.withValues(alpha: 0.5),
          ),
          const SizedBox(height: 16),
          const Text(
            'Tu lista de compras está vacía',
            style: TextStyle(
              fontSize: 18,
              color: AppColors.textSecondary,
            ),
          ),
          const SizedBox(height: 8),
          const Text(
            'Escanea productos para agregarlos',
            style: TextStyle(
              fontSize: 14,
              color: AppColors.textSecondary,
            ),
          ),
          const SizedBox(height: 24),
          ElevatedButton.icon(
            onPressed: _agregarProducto,
            icon: const Icon(Icons.qr_code_scanner),
            label: const Text('Escanear Producto'),
          ),
        ],
      ),
    );
  }

  Widget _buildList() {
    final pendientes = _items.where((i) => !i.comprado).toList();
    final comprados = _items.where((i) => i.comprado).toList();

    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        // Resumen
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                _buildStat('Total', _items.length.toString(), AppColors.primary),
                _buildStat('Pendientes', pendientes.length.toString(), AppColors.warning),
                _buildStat('Comprados', comprados.length.toString(), AppColors.success),
              ],
            ),
          ),
        ),
        const SizedBox(height: 16),

        // Pendientes
        if (pendientes.isNotEmpty) ...[
          const Text(
            'Por comprar',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.bold,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 8),
          ...pendientes.map((item) => _buildItem(item)),
          const SizedBox(height: 16),
        ],

        // Comprados
        if (comprados.isNotEmpty) ...[
          const Text(
            'Comprados',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.bold,
              color: AppColors.textSecondary,
            ),
          ),
          const SizedBox(height: 8),
          ...comprados.map((item) => _buildItem(item)),
        ],
      ],
    );
  }

  Widget _buildStat(String label, String value, Color color) {
    return Column(
      children: [
        Text(
          value,
          style: TextStyle(
            fontSize: 24,
            fontWeight: FontWeight.bold,
            color: color,
          ),
        ),
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

  Widget _buildItem(ItemListaCompra item) {
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: ListTile(
        leading: Checkbox(
          value: item.comprado,
          onChanged: (value) => _toggleComprado(item),
        ),
        title: Text(
          item.producto.nombre,
          style: TextStyle(
            decoration: item.comprado ? TextDecoration.lineThrough : null,
          ),
        ),
        subtitle: Text(
          'Cantidad: ${item.cantidad} - \$${item.producto.precio.toStringAsFixed(2)}',
        ),
        trailing: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              '\$${(item.producto.precio * item.cantidad).toStringAsFixed(2)}',
              style: const TextStyle(fontWeight: FontWeight.bold),
            ),
            IconButton(
              icon: const Icon(Icons.delete_outline, color: AppColors.error),
              onPressed: () => _eliminarItem(item),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _agregarProducto() async {
    // Usar el scanner para agregar producto
    final result = await Navigator.of(context).push<Producto>(
      MaterialPageRoute(
        builder: (_) => BarcodeScannerScreen(
          onProductFound: (producto) {
            Navigator.of(context).pop(producto);
          },
          autoPop: true,
        ),
      ),
    );

    if (result != null) {
      // Mostrar diálogo para cantidad
      _mostrarDialogoCantidad(result);
    }
  }

  void _mostrarDialogoCantidad(Producto producto) {
    int cantidad = 1;
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(producto.nombre),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Precio: \$${producto.precio.toStringAsFixed(2)}'),
            Text('Stock actual: ${producto.stockActual}'),
            const SizedBox(height: 16),
            Row(
              children: [
                const Text('Cantidad: '),
                IconButton(
                  icon: const Icon(Icons.remove),
                  onPressed: () {
                    if (cantidad > 1) {
                      cantidad--;
                    }
                  },
                ),
                Text('$cantidad', style: const TextStyle(fontSize: 18)),
                IconButton(
                  icon: const Icon(Icons.add),
                  onPressed: () {
                    cantidad++;
                  },
                ),
              ],
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancelar'),
          ),
          ElevatedButton(
            onPressed: () {
              _agregarItem(producto, cantidad);
              Navigator.pop(context);
            },
            child: const Text('Agregar'),
          ),
        ],
      ),
    );
  }

  void _agregarItem(Producto producto, int cantidad) {
    final item = ItemListaCompra(
      id: DateTime.now().millisecondsSinceEpoch.toString(),
      producto: producto,
      cantidad: cantidad,
    );
    setState(() => _items.add(item));
  }

  void _toggleComprado(ItemListaCompra item) {
    setState(() {
      final index = _items.indexOf(item);
      if (index != -1) {
        _items[index] = item.copyWith(comprado: !item.comprado);
      }
    });
  }

  void _eliminarItem(ItemListaCompra item) {
    setState(() => _items.remove(item));
  }

  void _limpiarLista() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('¿Limpiar lista?'),
        content: const Text('Se eliminarán todos los productos de la lista.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancelar'),
          ),
          ElevatedButton(
            style: ElevatedButton.styleFrom(backgroundColor: AppColors.error),
            onPressed: () {
              setState(() => _items.clear());
              Navigator.pop(context);
            },
            child: const Text('Limpiar'),
          ),
        ],
      ),
    );
  }
}