import 'package:flutter/material.dart';
import 'package:mobile_scanner/mobile_scanner.dart';
import '../../../core/theme/app_colors.dart';
import '../../../domain/entities/producto.dart';
import '../../../data/api/api_client.dart';
import '../../../core/constants/api_constants.dart';
import '../../../data/models/producto_model.dart';
import '../productos/producto_form_screen.dart';
import '../productos/producto_detail_screen.dart';

class BarcodeScannerScreen extends StatefulWidget {
  const BarcodeScannerScreen({super.key});

  @override
  State<BarcodeScannerScreen> createState() => _BarcodeScannerScreenState();
}

class _BarcodeScannerScreenState extends State<BarcodeScannerScreen> {
  final MobileScannerController _scannerController = MobileScannerController();
  bool _isScanning = true;
  String? _lastScanned;
  Producto? _foundProduct;

  @override
  void dispose() {
    _scannerController.dispose();
    super.dispose();
  }

  Future<void> _onDetect(BarcodeCapture capture) async {
    if (!_isScanning) return;

    final barcodes = capture.barcodes;
    if (barcodes.isEmpty) return;

    final code = barcodes.first.rawValue;
    if (code == null || code == _lastScanned) return;

    setState(() {
      _isScanning = false;
      _lastScanned = code;
    });

    // Buscar producto
    await _searchProduct(code);
  }

  Future<void> _searchProduct(String code) async {
    try {
      final response = await ApiClient.instance.get(
        '${ApiConstants.productos}/barcode/$code',
      );
      final producto = ProductoModel.fromJson(response.data);
      setState(() => _foundProduct = producto);
    } catch (e) {
      // Producto no encontrado - está OK
      setState(() => _foundProduct = null);
    }
  }

  void _reset() {
    setState(() {
      _isScanning = true;
      _lastScanned = null;
      _foundProduct = null;
    });
    _scannerController.start();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Escanear Código'),
        actions: [
          IconButton(
            icon: ValueListenableBuilder(
              valueListenable: _scannerController,
              builder: (context, state, _) {
                return Icon(
                  state.torchState == TorchState.on
                      ? Icons.flash_on
                      : Icons.flash_off,
                );
              },
            ),
            onPressed: () => _scannerController.toggleTorch(),
          ),
          IconButton(
            icon: const Icon(Icons.flip_camera_android),
            onPressed: () => _scannerController.switchCamera(),
          ),
        ],
      ),
      body: _buildBody(),
    );
  }

  Widget _buildBody() {
    if (_foundProduct != null || _lastScanned != null) {
      return _buildResult();
    }

    return Stack(
      children: [
        MobileScanner(
          controller: _scannerController,
          onDetect: _onDetect,
        ),
        // Overlay
        Center(
          child: Container(
            width: 250,
            height: 250,
            decoration: BoxDecoration(
              border: Border.all(color: AppColors.primary, width: 2),
              borderRadius: BorderRadius.circular(12),
            ),
          ),
        ),
        Positioned(
          bottom: 0,
          left: 0,
          right: 0,
          child: Container(
            padding: const EdgeInsets.all(24),
            decoration: BoxDecoration(
              gradient: LinearGradient(
                begin: Alignment.bottomCenter,
                end: Alignment.topCenter,
                colors: [
                  Colors.black.withValues(alpha: 0.8),
                  Colors.transparent,
                ],
              ),
            ),
            child: Column(
              children: [
                const Text(
                  'Apunta al código de barras',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                  ),
                ),
                const SizedBox(height: 8),
                const Text(
                  'El código se detectará automáticamente',
                  style: TextStyle(
                    color: Colors.white70,
                    fontSize: 12,
                  ),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildResult() {
    final code = _lastScanned!;
    final producto = _foundProduct;

    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          // Código escaneado
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                children: [
                  const Icon(Icons.qr_code, size: 48, color: AppColors.primary),
                  const SizedBox(height: 8),
                  const Text(
                    'Código escaneado',
                    style: TextStyle(color: AppColors.textSecondary),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    code,
                    style: const TextStyle(
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),

          if (producto != null) ...[
            // Producto encontrado
            Card(
              color: AppColors.success.withValues(alpha: 0.1),
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        const Icon(Icons.check_circle, color: AppColors.success),
                        const SizedBox(width: 8),
                        const Text(
                          'Producto encontrado',
                          style: TextStyle(
                            fontWeight: FontWeight.bold,
                            color: AppColors.success,
                          ),
                        ),
                      ],
                    ),
                    const Divider(),
                    Text(
                      producto.nombre,
                      style: const TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text('Código: ${producto.codigo}'),
                    Text('Stock: ${producto.stockActual}'),
                    Text('Precio: \$${producto.precio.toStringAsFixed(2)}'),
                    const SizedBox(height: 16),
                    ElevatedButton(
                      onPressed: () {
                        Navigator.of(context).push(
                          MaterialPageRoute(
                            builder: (_) => ProductoDetailScreen(producto: producto),
                          ),
                        );
                      },
                      child: const Text('VER DETALLE'),
                    ),
                  ],
                ),
              ),
            ),
          ] else ...[
            // Producto no encontrado
            Card(
              color: AppColors.warning.withValues(alpha: 0.1),
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  children: [
                    const Icon(Icons.warning, color: AppColors.warning, size: 48),
                    const SizedBox(height: 8),
                    const Text(
                      'Producto no encontrado',
                      style: TextStyle(
                        fontWeight: FontWeight.bold,
                        color: AppColors.warning,
                      ),
                    ),
                    const SizedBox(height: 8),
                    const Text(
                      '¿Deseas crear un nuevo producto con este código?',
                      textAlign: TextAlign.center,
                    ),
                    const SizedBox(height: 16),
                    ElevatedButton.icon(
                      onPressed: () {
                        Navigator.of(context).push(
                          MaterialPageRoute(
                            builder: (_) => ProductoFormScreen(
                              producto: Producto(
                                id: '',
                                codigo: code,
                                codigoBarras: code,
                                nombre: '',
                                precio: 0,
                                stockActual: 0,
                                stockMinimo: 5,
                              ),
                            ),
                          ),
                        );
                      },
                      icon: const Icon(Icons.add),
                      label: const Text('CREAR PRODUCTO'),
                    ),
                  ],
                ),
              ),
            ),
          ],
          const SizedBox(height: 16),

          // Botón escanear otro
          OutlinedButton.icon(
            onPressed: _reset,
            icon: const Icon(Icons.qr_code_scanner),
            label: const Text('ESCANEAR OTRO'),
          ),
        ],
      ),
    );
  }
}