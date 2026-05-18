import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../core/constants/api_constants.dart';
import '../../../core/theme/app_colors.dart';
import '../../../domain/entities/movimiento.dart';
import '../../../data/api/api_client.dart';
import '../../../data/models/movimiento_model.dart';
import 'movimiento_form_screen.dart';

class MovimientoListScreen extends StatefulWidget {
  const MovimientoListScreen({super.key});

  @override
  State<MovimientoListScreen> createState() => _MovimientoListScreenState();
}

class _MovimientoListScreenState extends State<MovimientoListScreen> {
  List<Movimiento> _movimientos = [];
  bool _isLoading = true;
  String? _error;

  // Filtros
  String? _filterTipo;
  DateTime? _filterFechaDesde;
  DateTime? _filterFechaHasta;

  @override
  void initState() {
    super.initState();
    _loadMovimientos();
  }

  Future<void> _loadMovimientos() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final queryParams = <String, dynamic>{
        'page': 1,
        'limit': 50,
      };
      if (_filterTipo != null) queryParams['tipo'] = _filterTipo;
      if (_filterFechaDesde != null) {
        queryParams['fecha_desde'] = _filterFechaDesde!.toString().split(' ')[0];
      }
      if (_filterFechaHasta != null) {
        queryParams['fecha_hasta'] = _filterFechaHasta!.toString().split(' ')[0];
      }

      final response = await ApiClient.instance.get(
        ApiConstants.movimientos,
        queryParams: queryParams,
      );

      final List<dynamic> data = response.data['data'] ?? [];
      setState(() {
        _movimientos = data.map((json) => MovimientoModel.fromJson(json)).toList();
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Movimientos'),
        actions: [
          IconButton(
            icon: const Icon(Icons.filter_list),
            onPressed: _showFilterDialog,
          ),
        ],
      ),
      body: _buildBody(),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          Navigator.of(context).push(
            MaterialPageRoute(
              builder: (_) => const MovimientoFormScreen(),
            ),
          );
        },
        child: const Icon(Icons.add),
      ),
    );
  }

  Widget _buildBody() {
    if (_isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_error != null) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.error, color: AppColors.error, size: 48),
            const SizedBox(height: 16),
            Text(_error!),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: _loadMovimientos,
              child: const Text('Reintentar'),
            ),
          ],
        ),
      );
    }

    if (_movimientos.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(
              Icons.swap_horiz,
              size: 64,
              color: AppColors.textSecondary,
            ),
            const SizedBox(height: 16),
            const Text(
              'No hay movimientos',
              style: TextStyle(
                fontSize: 18,
                color: AppColors.textSecondary,
              ),
            ),
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: _loadMovimientos,
      child: ListView.builder(
        padding: const EdgeInsets.all(8),
        itemCount: _movimientos.length,
        itemBuilder: (context, index) {
          final movimiento = _movimientos[index];
          return _buildMovimientoCard(movimiento);
        },
      ),
    );
  }

  Widget _buildMovimientoCard(Movimiento movimiento) {
    final isEntrada = movimiento.tipo == TipoMovimiento.entrada;
    final color = isEntrada ? AppColors.success : AppColors.error;
    final icon = isEntrada ? Icons.arrow_downward : Icons.arrow_upward;

    return Card(
      child: ListTile(
        leading: CircleAvatar(
          backgroundColor: color.withValues(alpha: 0.1),
          child: Icon(icon, color: color),
        ),
        title: Text(movimiento.productoNombre ?? 'Producto'),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              '${isEntrada ? '+' : '-'}${movimiento.cantidad} unidades',
              style: TextStyle(
                color: color,
                fontWeight: FontWeight.bold,
              ),
            ),
            if (movimiento.motivo != null)
              Text(
                movimiento.motivo!,
                style: const TextStyle(fontSize: 12),
              ),
            Text(
              _formatDate(movimiento.fecha),
              style: const TextStyle(
                fontSize: 12,
                color: AppColors.textSecondary,
              ),
            ),
          ],
        ),
        isThreeLine: movimiento.motivo != null,
      ),
    );
  }

  String _formatDate(DateTime date) {
    final now = DateTime.now();
    final diff = now.difference(date);

    if (diff.inDays == 0) {
      return 'Hoy ${date.hour.toString().padLeft(2, '0')}:${date.minute.toString().padLeft(2, '0')}';
    } else if (diff.inDays == 1) {
      return 'Ayer';
    } else {
      return '${date.day}/${date.month}/${date.year}';
    }
  }

  void _showFilterDialog() {
    showModalBottomSheet(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setModalState) => Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              const Text(
                'Filtrar movimientos',
                style: TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 16),
              // Tipo
              const Text('Tipo', style: TextStyle(fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              Row(
                children: [
                  ChoiceChip(
                    label: const Text('Todos'),
                    selected: _filterTipo == null,
                    onSelected: (_) {
                      setModalState(() => _filterTipo = null);
                    },
                  ),
                  const SizedBox(width: 8),
                  ChoiceChip(
                    label: const Text('Entrada'),
                    selected: _filterTipo == 'entrada',
                    onSelected: (_) {
                      setModalState(() => _filterTipo = 'entrada');
                    },
                  ),
                  const SizedBox(width: 8),
                  ChoiceChip(
                    label: const Text('Salida'),
                    selected: _filterTipo == 'salida',
                    onSelected: (_) {
                      setModalState(() => _filterTipo = 'salida');
                    },
                  ),
                ],
              ),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: () {
                  Navigator.pop(context);
                  _loadMovimientos();
                },
                child: const Text('Aplicar filtros'),
              ),
              const SizedBox(height: 8),
              TextButton(
                onPressed: () {
                  setModalState(() {
                    _filterTipo = null;
                    _filterFechaDesde = null;
                    _filterFechaHasta = null;
                  });
                },
                child: const Text('Limpiar filtros'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}