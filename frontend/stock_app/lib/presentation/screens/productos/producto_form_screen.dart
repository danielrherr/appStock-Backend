import 'dart:io';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../core/theme/app_colors.dart';
import '../../../domain/entities/producto.dart';
import '../../../domain/entities/categoria.dart';
import '../../providers/producto_provider.dart';
import '../../providers/categoria_provider.dart';
import '../../widgets/common/app_button.dart';
import '../../widgets/producto/producto_image_picker.dart';

class ProductoFormScreen extends StatefulWidget {
  final Producto? producto;

  const ProductoFormScreen({super.key, this.producto});

  @override
  State<ProductoFormScreen> createState() => _ProductoFormScreenState();
}

class _ProductoFormScreenState extends State<ProductoFormScreen> {
  final _formKey = GlobalKey<FormState>();
  final _codigoController = TextEditingController();
  final _codigoBarrasController = TextEditingController();
  final _nombreController = TextEditingController();
  final _descripcionController = TextEditingController();
  final _precioController = TextEditingController();
  final _stockController = TextEditingController();
  final _stockMinimoController = TextEditingController();

  String? _selectedCategoriaId;
  File? _selectedImage;
  bool _isLoading = false;

  bool get _isEditing => widget.producto != null;

  @override
  void initState() {
    super.initState();
    if (_isEditing) {
      _codigoController.text = widget.producto!.codigo;
      _codigoBarrasController.text = widget.producto!.codigoBarras ?? '';
      _nombreController.text = widget.producto!.nombre;
      _descripcionController.text = widget.producto!.descripcion ?? '';
      _precioController.text = widget.producto!.precio.toString();
      _stockController.text = widget.producto!.stockActual.toString();
      _stockMinimoController.text = widget.producto!.stockMinimo.toString();
      _selectedCategoriaId = widget.producto!.categoriaId;
    } else {
      _stockController.text = '0';
      _stockMinimoController.text = '5';
    }

    // Cargar categorías
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<CategoriaProvider>().fetchCategorias();
    });
  }

  @override
  void dispose() {
    _codigoController.dispose();
    _codigoBarrasController.dispose();
    _nombreController.dispose();
    _descripcionController.dispose();
    _precioController.dispose();
    _stockController.dispose();
    _stockMinimoController.dispose();
    super.dispose();
  }

  Future<void> _save() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isLoading = true);

    final producto = Producto(
      id: widget.producto?.id ?? '',
      codigo: _codigoController.text.trim(),
      codigoBarras: _codigoBarrasController.text.trim().isEmpty 
          ? null 
          : _codigoBarrasController.text.trim(),
      nombre: _nombreController.text.trim(),
      descripcion: _descripcionController.text.trim().isEmpty 
          ? null 
          : _descripcionController.text.trim(),
      categoriaId: _selectedCategoriaId,
      precio: double.tryParse(_precioController.text) ?? 0,
      stockActual: int.tryParse(_stockController.text) ?? 0,
      stockMinimo: int.tryParse(_stockMinimoController.text) ?? 5,
    );

    final provider = context.read<ProductoProvider>();
    bool success;

    if (_isEditing) {
      success = await provider.updateProducto(producto);
    } else {
      success = await provider.createProducto(producto);
    }

    if (mounted) {
      setState(() => _isLoading = false);
      
      if (success) {
        Navigator.of(context).pop(true);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(_isEditing ? 'Producto actualizado' : 'Producto creado'),
            backgroundColor: AppColors.success,
          ),
        );
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(provider.error ?? 'Error al guardar'),
            backgroundColor: AppColors.error,
          ),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(_isEditing ? 'Editar Producto' : 'Nuevo Producto'),
      ),
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              // Imagen
              ProductoImagePicker(
                imageUrl: widget.producto?.imagen,
                imageFile: _selectedImage,
                onImageSelected: (file) {
                  setState(() => _selectedImage = file);
                },
              ),
              const SizedBox(height: 24),

              // Código
              TextFormField(
                controller: _codigoController,
                decoration: const InputDecoration(
                  labelText: 'Código *',
                  prefixIcon: Icon(Icons.tag),
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Ingresa el código';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),

              // Código de barras
              TextFormField(
                controller: _codigoBarrasController,
                decoration: const InputDecoration(
                  labelText: 'Código de barras',
                  prefixIcon: Icon(Icons.qr_code),
                ),
              ),
              const SizedBox(height: 16),

              // Nombre
              TextFormField(
                controller: _nombreController,
                decoration: const InputDecoration(
                  labelText: 'Nombre *',
                  prefixIcon: Icon(Icons.inventory_2),
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Ingresa el nombre';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),

              // Descripción
              TextFormField(
                controller: _descripcionController,
                decoration: const InputDecoration(
                  labelText: 'Descripción',
                  prefixIcon: Icon(Icons.description),
                ),
                maxLines: 3,
              ),
              const SizedBox(height: 16),

              // Categoría
              Consumer<CategoriaProvider>(
                builder: (context, provider, _) {
                  return DropdownButtonFormField<String>(
                    value: _selectedCategoriaId,
                    decoration: const InputDecoration(
                      labelText: 'Categoría',
                      prefixIcon: Icon(Icons.category),
                    ),
                    items: [
                      const DropdownMenuItem(
                        value: null,
                        child: Text('Sin categoría'),
                      ),
                      ...provider.categorias.map((cat) => DropdownMenuItem(
                        value: cat.id,
                        child: Text(cat.nombre),
                      )),
                    ],
                    onChanged: (value) {
                      setState(() => _selectedCategoriaId = value);
                    },
                  );
                },
              ),
              const SizedBox(height: 16),

              // Precio
              TextFormField(
                controller: _precioController,
                decoration: const InputDecoration(
                  labelText: 'Precio *',
                  prefixIcon: Icon(Icons.attach_money),
                ),
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Ingresa el precio';
                  }
                  if (double.tryParse(value) == null) {
                    return 'Precio inválido';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),

              // Stock
              Row(
                children: [
                  Expanded(
                    child: TextFormField(
                      controller: _stockController,
                      decoration: const InputDecoration(
                        labelText: 'Stock',
                        prefixIcon: Icon(Icons.inventory),
                      ),
                      keyboardType: TextInputType.number,
                    ),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: TextFormField(
                      controller: _stockMinimoController,
                      decoration: const InputDecoration(
                        labelText: 'Stock mínimo',
                        prefixIcon: Icon(Icons.warning),
                      ),
                      keyboardType: TextInputType.number,
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'Requerido';
                        }
                        return null;
                      },
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 32),

              // Botón guardar
              AppButton(
                text: _isEditing ? 'ACTUALIZAR PRODUCTO' : 'GUARDAR PRODUCTO',
                onPressed: _save,
                isLoading: _isLoading,
                icon: Icons.save,
              ),
            ],
          ),
        ),
      ),
    );
  }
}