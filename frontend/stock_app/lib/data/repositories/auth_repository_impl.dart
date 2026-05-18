import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import '../../core/constants/api_constants.dart';
import '../../core/constants/app_constants.dart';
import '../../domain/entities/user.dart';
import '../../domain/repositories/auth_repository.dart';
import '../api/api_client.dart';

class AuthRepositoryImpl implements AuthRepository {
  final ApiClient _apiClient = ApiClient.instance;

  @override
  Future<User> login(String email, String password) async {
    final response = await _apiClient.post(
      ApiConstants.login,
      data: {'email': email, 'password': password},
    );

    final data = response.data;
    final token = data['token'];
    final userData = data['user'];

    // Save token
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(AppConstants.tokenKey, token);
    await prefs.setString(AppConstants.userKey, jsonEncode(userData));

    return _mapUser(userData);
  }

  @override
  Future<User> register(String email, String password, String? nombre) async {
    final response = await _apiClient.post(
      ApiConstants.register,
      data: {
        'email': email,
        'password': password,
        if (nombre != null) 'nombre': nombre,
      },
    );

    final data = response.data;
    final token = data['token'];
    final userData = data['user'];

    // Save token
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(AppConstants.tokenKey, token);
    await prefs.setString(AppConstants.userKey, jsonEncode(userData));

    return _mapUser(userData);
  }

  @override
  Future<void> logout() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(AppConstants.tokenKey);
    await prefs.remove(AppConstants.userKey);
  }

  @override
  Future<User?> getCurrentUser() async {
    final prefs = await SharedPreferences.getInstance();
    final userJson = prefs.getString(AppConstants.userKey);
    if (userJson != null) {
      return _mapUser(jsonDecode(userJson));
    }
    return null;
  }

  @override
  Future<bool> isLoggedIn() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.containsKey(AppConstants.tokenKey);
  }

  User _mapUser(Map<String, dynamic> json) {
    return User(
      id: json['id'] ?? '',
      email: json['email'] ?? '',
      nombre: json['nombre'],
      rol: json['rol'] == 'admin' ? UserRole.admin : UserRole.usuario,
      createdAt: json['created_at'] != null 
          ? DateTime.parse(json['created_at']) 
          : null,
    );
  }
}