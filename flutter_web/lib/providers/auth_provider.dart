import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';

class AuthProvider extends ChangeNotifier {
  bool _isAuthenticated = false;
  String? _userId;
  String? _userName;
  String? _userEmail;

  bool get isAuthenticated => _isAuthenticated;
  String? get userId => _userId;
  String? get userName => _userName;
  String? get userEmail => _userEmail;

  AuthProvider() {
    _loadAuthState();
  }

  Future<void> _loadAuthState() async {
    final prefs = await SharedPreferences.getInstance();
    _isAuthenticated = prefs.getBool('isAuthenticated') ?? false;
    _userId = prefs.getString('userId');
    _userName = prefs.getString('userName');
    _userEmail = prefs.getString('userEmail');
    notifyListeners();
  }

  Future<void> _saveAuthState() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setBool('isAuthenticated', _isAuthenticated);
    await prefs.setString('userId', _userId ?? '');
    await prefs.setString('userName', _userName ?? '');
    await prefs.setString('userEmail', _userEmail ?? '');
  }

  Future<void> login(String email, String password) async {
    // Mock authentication - in real app, this would call your backend
    if (email.isNotEmpty && password.isNotEmpty) {
      _isAuthenticated = true;
      _userId = 'test-user-123';
      _userName = 'Test User';
      _userEmail = email;
      
      await _saveAuthState();
      notifyListeners();
    }
  }

  Future<void> logout() async {
    _isAuthenticated = false;
    _userId = null;
    _userName = null;
    _userEmail = null;
    
    await _saveAuthState();
    notifyListeners();
  }

  Future<void> signUp(String name, String email, String password) async {
    // Mock signup - in real app, this would call your backend
    if (name.isNotEmpty && email.isNotEmpty && password.isNotEmpty) {
      _isAuthenticated = true;
      _userId = 'test-user-123';
      _userName = name;
      _userEmail = email;
      
      await _saveAuthState();
      notifyListeners();
    }
  }
} 