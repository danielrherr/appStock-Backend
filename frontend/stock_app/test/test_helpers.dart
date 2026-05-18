import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';

// Mock classes
class MockDio extends Mock implements Dio {}

class MockResponse extends Mock implements Response {
  @override
  final dynamic data;
  final int statusCode;

  MockResponse({required this.data, this.statusCode = 200});
}

// Test helper para crear respuestas mockeadas
Response<T> createMockResponse<T>(dynamic data, {int statusCode = 200}) {
  return Response<T>(
    data: data as T,
    statusCode: statusCode,
    requestOptions: RequestOptions(path: ''),
  );
}

// Matcher para verificar excepciones de Dio
Matcher throwsDioException([String? message]) {
  return throwsA(
    predicate((e) => e is DioException && (message == null || e.message == message)),
  );
}

// Matcher para respuestas exitosas
 Matcher isSuccessfulResponse([int statusCode = 200]) {
  return predicate<Response>((resp) => resp.statusCode == statusCode);
}

// Matcher para respuestas de error
Matcher isErrorResponse([int? statusCode]) {
  return predicate<Response>((resp) {
    if (statusCode != null) {
      return resp.statusCode == statusCode;
    }
    return resp.statusCode! >= 400;
  });
}