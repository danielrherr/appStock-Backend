import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../core/constants/api_constants.dart';
import '../../../core/theme/app_colors.dart';
import '../../../domain/entities/producto.dart';
import '../../../domain/entities/movimiento.dart';
import '../../../data/api/api_client.dart';
import '../scanner/barcode_scanner_screen.dart';

class MovimientoFormScreen extends StatefulWidget {
  final Producto? producto;

  const MovimientoFormScreen({super.key, this.producto});

  @override
  State<MovimientoFormScreen> createState() => _MovimientoFormScreenState();
}

class _MovimientoFormScreenState extends State<MovimientoFormScreen> {
  final _formKey = GlobalKey<FormState>();
  final _cantidadController = TextEditingController();
  final _motivoController = TextEditingController();

  TipoMovimiento _tipo = TipoMovimiento.entrada;
  Producto? _selectedProducto;
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _selectedProducto = widget.producto;
  }

  @override
  void dispose() {
    _cantidadController.dispose();
    _motivoController.dispose();
    super.dispose();
  }

  Future<void> _scanProduct() async {
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
      setState(() => _selectedProducto = result);
    }
  }

  Future<void> _save() async {
    if (!_formKey.currentState!.validate()) return;
    if (_selectedProducto == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Selecciona un producto'),
          backgroundColor: AppColors.error,
        ),
      );
      return;
    }

    setState(() => _isLoading = true);

    try {
      await ApiClient.instance.post(
        ApiConstants.movimientos,
        data: {
          'producto_id': _selectedProducto!.id,
          'tipo': _tipo == TipoMovimiento.entrada ? 'entrada' : 'salida',
          'cantidad': int.parse(_cantidadController.text),
          'motivo': _motivoController.text.trim().isEmpty 
              ? null 
              : _motivoController.text.trim(),
        },
      );

      if (mounted) {
        Navigator.of(context).pop(true);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              _tipo == TipoMovimiento.entrada 
                  ? 'Entrada registrada' 
                  : 'Salida registrada',
            ),
            backgroundColor: AppColors.success,
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Error: ${e.toString()}'),
            backgroundColor: AppColors.error,
          ),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Registrar Movimiento'),
      ),
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              // Tipo de movimiento
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(8),
                  child: Row(
                    children: [
                      Expanded(
                        child: _buildTipoOption(
                          'Entrada',
                          Icons.arrow_downward,
                          TipoMovimiento.entrada,
                          AppColors.success,
                        ),
                      ),
                      const SizedBox(width: 8),
                      Expanded(
                        child: _buildTipoOption(
                          'Salida',
                          Icons.arrow_upward,
                          TipoMovimiento.salida,
                          AppColors.error,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 24),

              // Producto
              if (widget.producto == null) ...[
                Row(
                  children: [
                    Expanded(
                      child: TextFormField(
                        decoration: const InputDecoration(
                          labelText: 'Producto *',
                          prefixIcon: Icon(Icons.inventory_2),
                          hintText: 'Selecciona un producto',
                        ),
                        readOnly: true,
                        onTap: () {
                          // Por ahora mostrar mensaje
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(
                              content: Text('Usa el botón de scanner para buscar'),
                            ),
                          );
                        },
                        controller: TextEditingController(
                          text: _selectedProducto?.nombre ?? '',
                        ),
                      ),
                    ),
                    const SizedBox(width: 8),
                    IconButton.filled(
                      onPressed: _scanProduct,
                      icon: const Icon(Icons.qr_code_scanner),
                      tooltip: 'Escanear código de barras',
                    ),
                  ],
                ),
              ] else ...[
                Card(
                  child: ListTile(
                    leading: const Icon(Icons.inventory_2),
                    title: Text(_selectedProducto!.nombre),
                    subtitle: Text('Stock: ${_selectedProducto!.stockActual}'),
                  ),
                ),
              ],
              const SizedBox(height: 16),

              // Cantidad
              TextFormField(
                controller: _cantidadController,
                decoration: const InputDecoration(
                  labelText: 'Cantidad *',
                  prefixIcon: Icon(Icons.numbers),
                ),
                keyboardType: TextInputType.number,
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Ingresa la cantidad';
                  }
                  final cantidad = int.tryParse(value);
                  if (cantidad == null || cantidad <= 0) {
                    return 'Cantidad inválida';
                  }
                  if (_tipo == TipoMovimiento.salida && _selectedProducto != null) {
                    if (cantidad > _selectedProducto!.stockActual) {
                      return 'Stock insuficiente';
                    }
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),

              // Motivo
              TextFormField(
                controller: _motivoController,
                decoration: const InputDecoration(
                  labelText: 'Motivo (opcional)',
                  prefixIcon: Icon(Icons.note),
                  hintText: 'Ej: Compra proveedor, Venta, etc.',
                ),
                maxLines: 2,
              ),
              const SizedBox(height: 32),

              // Preview
              if (_selectedProducto != null && _cantidadController.text.isNotEmpty)
                Card(
                  color: AppColors.background,
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Text(
                          'Preview',
                          style: TextStyle(
                            fontWeight: FontWeight.bold,
                            color: AppColors.textSecondary,
                          ),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          _tipo == TipoMovimiento.entrada
                              ? 'Stock actual: ${_selectedProducto!.stockActual} + ${_cantidadController.text} = ${_selectedProducto!.stockActual + (int.tryParse(_cantidadController.text) ?? 0)}'
                              : 'Stock actual: ${_selectedProducto!.stockActual} - ${_cantidadController.text} = ${_selectedProducto!.stockActual - (int.tryParse(_cantidadController.text) ?? 0)}',
                        ),
                      ],
                    ),
                  ),
                ),
              const SizedBox(height: 24),

              // Botón
              ElevatedButton.icon(
                onPressed: _isLoading ? null : _save,
                icon: _isLoading
                    ? const SizedBox(
                        height: 20,
                        width: 20,
                        child: CircularProgressIndicator(
                          strokeWidth: 2,
                          color: Colors.white,
                        ),
                      )
                    : const Icon(Icons.check),
                label: Text(_tipo == TipoMovimiento.entrada 
                    ? 'REGISTRAR ENTRADA' 
                    : 'REGISTRAR SALIDA'),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildTipoOption(String label, IconData icon, TipoMovimiento tipo, Color color) {
    final isSelected = _tipo == tipo;
    return InkWell(
      onTap: () => setState(() => _tipo = tipo),
      borderRadius: BorderRadius.circular(8),
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 12),
        decoration: BoxDecoration(
          color: isSelected ? color.withValues(alpha: 0.1) : null,
          borderRadius: BorderRadius.circular(8),
          border: Border.all(
            color: isSelected ? color : AppColors.border,
            width: isSelected ? 2 : 1,
          ),
        ),
        child: Column(
          children: [
            Icon(icon, color: color, size: 28),
            const SizedBox(height: 4),
            Text(
              label,
              style: TextStyle(
                color: isSelected ? color : AppColors.textSecondary,
                fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
              ),
            ),
          ],
        ),
      ),
    );
  }
}