import 'package:shared_preferences/shared_preferences.dart';
import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import '../core/constants/api_constants.dart';

// Firebase solo para mobile - conditional import
import 'package:firebase_messaging/firebase_messaging.dart'
    show FirebaseMessaging
    if (dart.library.js) 'package:flutter/foundation.dart';

class NotificationService {
  static NotificationService? _instance;
  late final _NotificationServiceImpl _impl;
  
  factory NotificationService() {
    _instance ??= NotificationService._internal();
    return _instance!;
  }
  
  NotificationService._internal() {
    if (kIsWeb) {
      _impl = _NotificationServiceWeb();
    } else {
      _impl = _NotificationServiceMobile();
    }
  }

  Function(Map<String, dynamic>)? get onNotificationTap => _impl.onNotificationTap;
  set onNotificationTap(Function(Map<String, dynamic>)? value) => _impl.onNotificationTap = value;

  Future<void> initialize() => _impl.initialize();
  String get token => _impl.token;
  Future<void> refreshToken() => _impl.refreshToken();
  Future<void> subscribeToTopic(String topic) => _impl.subscribeToTopic(topic);
  Future<void> unsubscribeFromTopic(String topic) => _impl.unsubscribeFromTopic(topic);
  Future<void> deleteToken() => _impl.deleteToken();
}

abstract class _NotificationServiceImpl {
  Function(Map<String, dynamic>)? onNotificationTap;
  Future<void> initialize();
  String get token;
  Future<void> refreshToken();
  Future<void> subscribeToTopic(String topic);
  Future<void> unsubscribeFromTopic(String topic);
  Future<void> deleteToken();
}

// Mobile implementation
class _NotificationServiceMobile extends _NotificationServiceImpl {
  late final FirebaseMessaging _firebaseMessaging;
  String? _deviceToken;

  _NotificationServiceMobile() : _firebaseMessaging = FirebaseMessaging.instance;

  @override
  Function(Map<String, dynamic>)? onNotificationTap;

  @override
  Future<void> initialize() async {
    final settings = await _firebaseMessaging.requestPermission(
      alert: true, badge: true, sound: true,
    );
    print('Notification permissions: ${settings.authorizationStatus}');
    
    _deviceToken = await _firebaseMessaging.getToken();
    print('Device Token: $_deviceToken');
    
    if (_deviceToken != null) {
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('device_token', _deviceToken!);
      await _registerToken();
    }

    FirebaseMessaging.onMessage.listen((message) {
      print('Received message: ${message.notification?.title}');
      if (message.data.isNotEmpty) _processNotificationData(message.data);
    });

    FirebaseMessaging.onMessageOpenedApp.listen((message) {
      print('App opened from notification');
      if (onNotificationTap != null) onNotificationTap!(message.data);
    });
  }

  Future<void> _registerToken() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final token = prefs.getString('auth_token');
      if (token == null) return;

      final dio = Dio(BaseOptions(
        baseUrl: ApiConstants.baseUrl,
        headers: {'Authorization': 'Bearer $token', 'Content-Type': 'application/json'},
      ));
      await dio.post('/devices', data: {'token': _deviceToken, 'platform': 'mobile'});
    } catch (e) { print('Error registering token: $e'); }
  }

  void _processNotificationData(Map<String, dynamic> data) {
    final type = data['type'] as String?;
    switch (type) {
      case 'stock_alert': print('Stock alert: ${data['product_id']}'); break;
      case 'product_update': print('Product update: ${data['product_id']}'); break;
      default: print('Unknown: $type');
    }
  }

  @override String get token => _deviceToken ?? '';
  @override Future<void> refreshToken() async => _registerToken();
  @override Future<void> subscribeToTopic(String topic) => _firebaseMessaging.subscribeToTopic(topic);
  @override Future<void> unsubscribeFromTopic(String topic) => _firebaseMessaging.unsubscribeFromTopic(topic);
  @override Future<void> deleteToken() async {
    await _firebaseMessaging.deleteToken();
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('device_token');
    _deviceToken = null;
  }
}

// Web stub
class _NotificationServiceWeb extends _NotificationServiceImpl {
  @override
  Function(Map<String, dynamic>)? onNotificationTap;

  @override
  Future<void> initialize() async => print('Web: no push notifications');

  @override
  String get token => '';

  @override
  Future<void> refreshToken() async {}

  @override
  Future<void> subscribeToTopic(String topic) async => print('Web: no topics');

  @override
  Future<void> unsubscribeFromTopic(String topic) async {}

  @override
  Future<void> deleteToken() async {}
}