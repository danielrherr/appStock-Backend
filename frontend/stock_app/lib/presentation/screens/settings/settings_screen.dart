import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../core/constants/app_constants.dart';
import '../../../core/theme/app_colors.dart';
import '../../providers/auth_provider.dart';
import '../auth/login_screen.dart';

class SettingsScreen extends StatelessWidget {
  const SettingsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Configuración'),
      ),
      body: ListView(
        children: [
          // Perfil
          Consumer<AuthProvider>(
            builder: (context, auth, _) {
              final user = auth.currentUser;
              return UserAccountsDrawerHeader(
                accountName: Text(user?.nombre ?? 'Usuario'),
                accountEmail: Text(user?.email ?? ''),
                currentAccountPicture: CircleAvatar(
                  backgroundColor: AppColors.primary,
                  child: Text(
                    (user?.nombre ?? 'U')[0].toUpperCase(),
                    style: const TextStyle(color: Colors.white, fontSize: 24),
                  ),
                ),
              );
            },
          ),
          const Divider(),

          // General
          const _SectionHeader(title: 'General'),
          ListTile(
            leading: const Icon(Icons.language),
            title: const Text('Idioma'),
            subtitle: const Text('Español'),
            onTap: () {},
          ),
          ListTile(
            leading: const Icon(Icons.attach_money),
            title: const Text('Moneda'),
            subtitle: const Text('ARS (\$)'),
            onTap: () {},
          ),
          const Divider(),

          // Notifications
          const _SectionHeader(title: 'Notificaciones'),
          SwitchListTile(
            secondary: const Icon(Icons.notifications),
            title: const Text('Notificaciones push'),
            subtitle: const Text('Alertas de stock bajo'),
            value: true,
            onChanged: (value) {},
          ),
          const Divider(),

          // Data
          const _SectionHeader(title: 'Datos'),
          ListTile(
            leading: const Icon(Icons.cloud_upload),
            title: const Text('Exportar datos'),
            subtitle: const Text('Descargar base de datos'),
            onTap: () {},
          ),
          ListTile(
            leading: const Icon(Icons.cloud_download),
            title: const Text('Importar datos'),
            subtitle: const Text('Restaurar respaldo'),
            onTap: () {},
          ),
          ListTile(
            leading: const Icon(Icons.delete_forever, color: AppColors.error),
            title: const Text(
              'Eliminar todos los datos',
              style: TextStyle(color: AppColors.error),
            ),
            subtitle: const Text('Borra todo el inventario'),
            onTap: () => _showDeleteDialog(context),
          ),
          const Divider(),

          // Acerca de
          const _SectionHeader(title: 'Acerca de'),
          ListTile(
            leading: const Icon(Icons.info),
            title: const Text('Versión'),
            subtitle: Text(AppConstants.appVersion),
          ),
          ListTile(
            leading: const Icon(Icons.description),
            title: const Text('Términos y condiciones'),
            onTap: () {},
          ),
          ListTile(
            leading: const Icon(Icons.privacy_tip),
            title: const Text('Política de privacidad'),
            onTap: () {},
          ),
          const Divider(),

          // Cerrar sesión
          ListTile(
            leading: const Icon(Icons.logout, color: AppColors.error),
            title: const Text(
              'Cerrar sesión',
              style: TextStyle(color: AppColors.error),
            ),
            onTap: () => _logout(context),
          ),
          const SizedBox(height: 24),
        ],
      ),
    );
  }

  void _logout(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Cerrar sesión'),
        content: const Text('¿Estás seguro de que quieres cerrar sesión?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancelar'),
          ),
          TextButton(
            onPressed: () {
              context.read<AuthProvider>().logout();
              Navigator.of(context).pushAndRemoveUntil(
                MaterialPageRoute(builder: (_) => const LoginScreen()),
                (route) => false,
              );
            },
            style: TextButton.styleFrom(foregroundColor: AppColors.error),
            child: const Text('Cerrar'),
          ),
        ],
      ),
    );
  }

  void _showDeleteDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Eliminar todos los datos'),
        content: const Text(
          'Esta acción no se puede deshacer. Se eliminarán todos los productos, categorías y movimientos.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancelar'),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(context);
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Text('Función no implementada aún'),
                ),
              );
            },
            style: TextButton.styleFrom(foregroundColor: AppColors.error),
            child: const Text('Eliminar'),
          ),
        ],
      ),
    );
  }
}

class _SectionHeader extends StatelessWidget {
  final String title;

  const _SectionHeader({required this.title});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
      child: Text(
        title.toUpperCase(),
        style: const TextStyle(
          color: AppColors.textSecondary,
          fontSize: 12,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }
}