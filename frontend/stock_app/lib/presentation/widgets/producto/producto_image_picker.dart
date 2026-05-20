import 'dart:io';
import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';
import '../../../core/theme/app_colors.dart';

class ProductoImagePicker extends StatelessWidget {
  final String? imageUrl;
  final File? imageFile;
  final Function(File) onImageSelected;

  const ProductoImagePicker({
    super.key,
    this.imageUrl,
    this.imageFile,
    required this.onImageSelected,
  });

  String _resolveImageUrl(String? url) {
    if (url == null || url.isEmpty) return '';
    if (url.startsWith('http')) return url;
    // Asumimos que es una ruta relativa del backend
    return 'https://appstock-backend1.onrender.com$url';
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      height: 200,
      decoration: BoxDecoration(
        color: AppColors.background,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.border),
      ),
      child: _buildContent(context),
    );
  }

  Widget _buildContent(BuildContext context) {
    // Priority: File > URL > Placeholder
    if (imageFile != null) {
      return Stack(
        fit: StackFit.expand,
        children: [
          ClipRRect(
            borderRadius: BorderRadius.circular(12),
            child: Image.file(
              imageFile!,
              fit: BoxFit.cover,
            ),
          ),
          _buildActions(context),
        ],
      );
    }

    if (imageUrl != null && imageUrl!.isNotEmpty) {
      final resolvedUrl = _resolveImageUrl(imageUrl);
      return Stack(
        fit: StackFit.expand,
        children: [
          ClipRRect(
            borderRadius: BorderRadius.circular(12),
            child: Image.network(
              resolvedUrl,
              fit: BoxFit.cover,
              errorBuilder: (_, __, ___) => _buildPlaceholder(context),
            ),
          ),
          _buildActions(context),
        ],
      );
    }

    return _buildPlaceholder(context);
  }

  Widget _buildPlaceholder(BuildContext context) {
    return InkWell(
      onTap: () => _pickImage(context),
      borderRadius: BorderRadius.circular(12),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(
            Icons.add_a_photo,
            size: 48,
            color: AppColors.textSecondary,
          ),
          const SizedBox(height: 8),
          const Text(
            'Agregar imagen',
            style: TextStyle(
              color: AppColors.textSecondary,
              fontSize: 14,
            ),
          ),
          const SizedBox(height: 4),
          const Text(
            '(opcional)',
            style: TextStyle(
              color: AppColors.textSecondary,
              fontSize: 12,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildActions(BuildContext context) {
    return Positioned(
      top: 8,
      right: 8,
      child: Row(
        children: [
          CircleAvatar(
            radius: 18,
            backgroundColor: Colors.black54,
            child: IconButton(
              icon: const Icon(Icons.camera_alt, size: 18),
              color: Colors.white,
              onPressed: () => _pickImage(context),
            ),
          ),
        ],
      ),
    );
  }

  Future<void> _pickImage(BuildContext context) async {
    final picker = ImagePicker();
    
    final source = await showModalBottomSheet<ImageSource>(
      context: context,
      builder: (context) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              leading: const Icon(Icons.camera_alt),
              title: const Text('Cámara'),
              onTap: () => Navigator.pop(context, ImageSource.camera),
            ),
            ListTile(
              leading: const Icon(Icons.photo_library),
              title: const Text('Galería'),
              onTap: () => Navigator.pop(context, ImageSource.gallery),
            ),
          ],
        ),
      ),
    );

    if (source != null) {
      final pickedFile = await picker.pickImage(
        source: source,
        maxWidth: 800,
        maxHeight: 800,
        imageQuality: 80,
      );

      if (pickedFile != null) {
        onImageSelected(File(pickedFile.path));
      }
    }
  }
}