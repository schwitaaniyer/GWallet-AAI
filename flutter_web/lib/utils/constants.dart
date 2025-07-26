import 'package:flutter/material.dart';

class AppColors {
  static const Color primary = Color(0xFF1A73E8);
  static const Color secondary = Color(0xFF34A853);
  static const Color accent = Color(0xFFEA4335);
  static const Color warning = Color(0xFFFBBC04);
  static const Color background = Color(0xFFF8F9FA);
  static const Color surface = Color(0xFFFFFFFF);
  static const Color textPrimary = Color(0xFF202124);
  static const Color textSecondary = Color(0xFF5F6368);
  static const Color textTertiary = Color(0xFF9AA0A6);
  static const Color border = Color(0xFFDADCE0);
  static const Color shadow = Color(0x1F000000);
}

class AppSizes {
  static const double paddingXS = 4.0;
  static const double paddingS = 8.0;
  static const double paddingM = 16.0;
  static const double paddingL = 24.0;
  static const double paddingXL = 32.0;
  
  static const double radiusS = 8.0;
  static const double radiusM = 12.0;
  static const double radiusL = 16.0;
  static const double radiusXL = 24.0;
  
  static const double iconS = 16.0;
  static const double iconM = 24.0;
  static const double iconL = 32.0;
  static const double iconXL = 48.0;
}

class ApiEndpoints {
  static const String baseUrl = 'http://localhost:8080';
  static const String health = '/health';
  static const String receipts = '/receipts';
  static const String queries = '/queries';
  static const String walletPasses = '/wallet-passes';
  static const String analysis = '/analysis';
}

class AppStrings {
  static const String appName = 'Raseed Wallet';
  static const String appDescription = 'AI-powered personal assistant for receipt management';
  
  // Navigation
  static const String home = 'Home';
  static const String receipts = 'Receipts';
  static const String passes = 'Passes';
  static const String analysis = 'Analysis';
  static const String settings = 'Settings';
  
  // Actions
  static const String uploadReceipt = 'Upload Receipt';
  static const String askQuestion = 'Ask Question';
  static const String viewDetails = 'View Details';
  static const String delete = 'Delete';
  static const String edit = 'Edit';
  static const String save = 'Save';
  static const String cancel = 'Cancel';
  
  // Messages
  static const String uploadSuccess = 'Receipt uploaded successfully!';
  static const String uploadError = 'Failed to upload receipt';
  static const String processingReceipt = 'Processing receipt...';
  static const String noReceipts = 'No receipts found';
  static const String noPasses = 'No wallet passes found';
  
  // Placeholders
  static const String searchReceipts = 'Search receipts...';
  static const String askAnything = 'Ask me anything about your receipts...';
  static const String enterAmount = 'Enter amount';
  static const String enterDescription = 'Enter description';
}

class AppIcons {
  static const String home = 'assets/icons/home.svg';
  static const String receipt = 'assets/icons/receipt.svg';
  static const String pass = 'assets/icons/pass.svg';
  static const String analysis = 'assets/icons/analysis.svg';
  static const String settings = 'assets/icons/settings.svg';
  static const String upload = 'assets/icons/upload.svg';
  static const String search = 'assets/icons/search.svg';
  static const String camera = 'assets/icons/camera.svg';
  static const String gallery = 'assets/icons/gallery.svg';
  static const String qrCode = 'assets/icons/qr_code.svg';
  static const String notification = 'assets/icons/notification.svg';
  static const String profile = 'assets/icons/profile.svg';
} 