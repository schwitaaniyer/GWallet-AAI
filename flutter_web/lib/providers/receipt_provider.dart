import 'dart:convert';
import 'dart:io';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:raseed_wallet/utils/constants.dart';
import 'package:raseed_wallet/models/receipt.dart';

class ReceiptProvider extends ChangeNotifier {
  List<Receipt> _receipts = [];
  bool _isLoading = false;
  String? _error;

  List<Receipt> get receipts => _receipts;
  bool get isLoading => _isLoading;
  String? get error => _error;

  Future<void> loadReceipts(String userId) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final response = await http.get(
        Uri.parse('${ApiEndpoints.baseUrl}${ApiEndpoints.receipts}?user_id=$userId'),
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        _receipts = (data['receipts'] as List)
            .map((json) => Receipt.fromJson(json))
            .toList();
      } else {
        _error = 'Failed to load receipts';
      }
    } catch (e) {
      _error = 'Network error: $e';
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<bool> uploadReceipt(String userId, File imageFile) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      var request = http.MultipartRequest(
        'POST',
        Uri.parse('${ApiEndpoints.baseUrl}${ApiEndpoints.receipts}'),
      );

      request.fields['user_id'] = userId;
      request.files.add(
        await http.MultipartFile.fromPath(
          'receipt',
          imageFile.path,
        ),
      );

      final response = await request.send();
      final responseData = await response.stream.bytesToString();

      if (response.statusCode == 200) {
        final data = json.decode(responseData);
        final newReceipt = Receipt.fromJson(data['receipt']);
        _receipts.insert(0, newReceipt);
        _isLoading = false;
        notifyListeners();
        return true;
      } else {
        _error = 'Failed to upload receipt';
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

  Future<void> deleteReceipt(String receiptId) async {
    try {
      final response = await http.delete(
        Uri.parse('${ApiEndpoints.baseUrl}${ApiEndpoints.receipts}/$receiptId'),
      );

      if (response.statusCode == 200) {
        _receipts.removeWhere((receipt) => receipt.id == receiptId);
        notifyListeners();
      }
    } catch (e) {
      _error = 'Failed to delete receipt: $e';
      notifyListeners();
    }
  }

  void clearError() {
    _error = null;
    notifyListeners();
  }
} 