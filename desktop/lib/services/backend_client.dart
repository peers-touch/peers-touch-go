import 'dart:convert';

import 'package:http/http.dart' as http;

class BackendConfig {
  // TODO: make this configurable via settings or env.
  static const String defaultBaseUrl = 'http://127.0.0.1:8082';
}

class BackendClient {
  final String baseUrl;
  final http.Client _client;

  BackendClient({String? baseUrl, http.Client? client})
      : baseUrl = baseUrl ?? BackendConfig.defaultBaseUrl,
        _client = client ?? http.Client();

  Uri _uri(String path, [Map<String, String>? query]) {
    return Uri.parse(baseUrl).replace(
      path: path,
      queryParameters: query ?? {},
    );
  }

  Future<Map<String, dynamic>> getJson(String path,
      {Map<String, String>? query}) async {
    final resp = await _client.get(_uri(path, query));
    if (resp.statusCode >= 200 && resp.statusCode < 300) {
      try {
        return jsonDecode(resp.body) as Map<String, dynamic>;
      } catch (_) {
        return {'raw': resp.body};
      }
    }
    throw Exception('GET $path failed: ${resp.statusCode}');
  }

  Future<String> getRaw(String path, {Map<String, String>? query}) async {
    final resp = await _client.get(_uri(path, query));
    if (resp.statusCode >= 200 && resp.statusCode < 300) {
      return resp.body;
    }
    throw Exception('GET $path failed: ${resp.statusCode}');
  }

  BackendClient withBaseUrl(String newBaseUrl) {
    return BackendClient(baseUrl: newBaseUrl, client: _client);
  }

  // Station APIs
  Future<Map<String, dynamic>> listPhotos({String? album}) async {
    return getJson('/station/photo/list', query: {
      if (album != null && album.isNotEmpty) 'album': album,
    });
  }

  Future<String> getPhoto({required String album, required String filename}) {
    return getRaw('/station/photo/get', query: {
      'album': album,
      'filename': filename,
    });
  }

  // Manage APIs
  Future<Map<String, dynamic>> ping() async {
    return getJson('/ping');
  }

  Future<String> health() async {
    return getRaw('/health');
  }

  // Bootstrap APIs
  Future<Map<String, dynamic>> listBootstrapPeers({int no = 1, int size = 20}) async {
    return getJson('/sub-bootstrap/list-peers', query: {
      'no': '$no',
      'size': '$size',
    });
  }
}