enum UserRole { admin, usuario }

class User {
  final String id;
  final String email;
  final String? nombre;
  final UserRole rol;
  final DateTime? createdAt;

  User({
    required this.id,
    required this.email,
    this.nombre,
    required this.rol,
    this.createdAt,
  });

  bool get isAdmin => rol == UserRole.admin;

  User copyWith({
    String? id,
    String? email,
    String? nombre,
    UserRole? rol,
    DateTime? createdAt,
  }) {
    return User(
      id: id ?? this.id,
      email: email ?? this.email,
      nombre: nombre ?? this.nombre,
      rol: rol ?? this.rol,
      createdAt: createdAt ?? this.createdAt,
    );
  }
}