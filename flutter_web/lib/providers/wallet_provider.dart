import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:raseed_wallet/utils/constants.dart';
import 'package:raseed_wallet/models/wallet_pass.dart';
import 'package:raseed_wallet/models/query.dart';

class WalletProvider extends ChangeNotifier {
  List<WalletPass> _walletPasses = [];
  List<Query> _queries = [];
  bool _isLoading = false;
  String? _error;

  List<WalletPass> get walletPasses => _walletPasses;
  List<Query> get queries => _queries;
  bool get isLoading => _isLoading;
  String? get error => _error;

  Future<void> loadWalletPasses(String userId) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final response = await http.get(
        Uri.parse('${ApiEndpoints.baseUrl}${ApiEndpoints.walletPasses}?user_id=$userId'),
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        _walletPasses = (data['wallet_passes'] as List)
            .map((json) => WalletPass.fromJson(json))
            .toList();
      } else {
        _error = 'Failed to load wallet passes';
      }
    } catch (e) {
      _error = 'Network error: $e';
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> loadQueries(String userId) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final response = await http.get(
        Uri.parse('${ApiEndpoints.baseUrl}${ApiEndpoints.queries}?user_id=$userId'),
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        _queries = (data['queries'] as List)
            .map((json) => Query.fromJson(json))
            .toList();
      } else {
        _error = 'Failed to load queries';
      }
    } catch (e) {
      _error = 'Network error: $e';
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<bool> submitQuery(String userId, String question, String language) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final response = await http.post(
        Uri.parse('${ApiEndpoints.baseUrl}${ApiEndpoints.queries}'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode({
          'user_id': userId,
          'query': question,
          'language': language,
        }),
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        final newQuery = Query.fromJson(data['query']);
        _queries.insert(0, newQuery);
        _isLoading = false;
        notifyListeners();
        return true;
      } else {
        _error = 'Failed to submit query';
        _isLoading = false;
        notifyListeners();
        return false;
      }
    } catch (e) {
      _error = 'Network error: $e';
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }

  Future<Map<String, dynamic>?> getSpendingAnalysis(String userId) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final response = await http.get(
        Uri.parse('${ApiEndpoints.baseUrl}${ApiEndpoints.analysis}?user_id=$userId'),
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        _isLoading = false;
        notifyListeners();
        return data;
      } else {
        _error = 'Failed to load spending analysis';
        _isLoading = false;
        notifyListeners();
        return null;
      }
    } catch (e) {
      _error = 'Network error: $e';
      _isLoading = false;
      notifyListeners();
      return null;
    }
  }

  void clearError() {
    _error = null;
    notifyListeners();
  }
} 