import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:provider/provider.dart';
import 'package:raseed_wallet/main.dart';
import 'package:raseed_wallet/providers/auth_provider.dart';
import 'package:raseed_wallet/providers/receipt_provider.dart';
import 'package:raseed_wallet/providers/wallet_provider.dart';
import 'package:raseed_wallet/screens/home_screen.dart';
import 'package:raseed_wallet/screens/login_screen.dart';
import 'package:raseed_wallet/screens/receipts_screen.dart';
import 'package:raseed_wallet/screens/passes_screen.dart';
import 'package:raseed_wallet/screens/analysis_screen.dart';
import 'package:raseed_wallet/widgets/quick_action_card.dart';
import 'package:raseed_wallet/widgets/stats_card.dart';
import 'package:raseed_wallet/models/receipt.dart';
import 'package:raseed_wallet/models/wallet_pass.dart';

// Mock HTTP client for testing
class MockHttpClient {
  static Map<String, dynamic> mockResponses = {
    '/health': {'status': 'healthy'},
    '/receipts': [
      {
        'id': 'test_receipt_123',
        'user_id': 'test_user_123',
        'store_name': 'Test Store',
        'total_amount': 45.99,
        'tax_amount': 3.50,
        'items': [
          {
            'name': 'Milk',
            'price': 4.99,
            'quantity': 2,
            'category': 'dairy'
          }
        ],
        'date': '2024-01-15T10:30:00Z',
        'image_url': 'https://example.com/receipt.jpg',
        'created_at': '2024-01-15T10:30:00Z',
        'updated_at': '2024-01-15T10:30:00Z'
      }
    ],
    '/queries': [
      {
        'id': 'test_query_123',
        'user_id': 'test_user_123',
        'query': 'What can I cook with my recent purchases?',
        'language': 'en',
        'response': 'Based on your recent purchases, you can make: 1. Scrambled eggs with toast 2. Pasta with tomato sauce',
        'created_at': '2024-01-15T10:30:00Z'
      }
    ],
    '/wallet-passes': [
      {
        'id': 'test_pass_123',
        'user_id': 'test_user_123',
        'type': 'receipt',
        'title': 'Receipt - Test Store',
        'description': 'Total: \$45.99, Items: 2',
        'data': '{"receipt_id": "test_receipt_123", "store_name": "Test Store"}',
        'created_at': '2024-01-15T10:30:00Z'
      }
    ],
    '/analysis': {
      'total_spent': 245.67,
      'category_spending': {
        'groceries': 120.50,
        'restaurants': 85.25,
        'transportation': 40.00
      },
      'receipt_count': 15,
      'average_per_receipt': 16.38
    }
  };
}

void main() {
  group('Raseed Wallet App Tests', () {
    testWidgets('App should start with login screen when not authenticated', (WidgetTester tester) async {
      await tester.pumpWidget(
        MultiProvider(
          providers: [
            ChangeNotifierProvider(create: (_) => AuthProvider()),
            ChangeNotifierProvider(create: (_) => ReceiptProvider()),
            ChangeNotifierProvider(create: (_) => WalletProvider()),
          ],
          child: const MaterialApp(home: AuthWrapper()),
        ),
      );

      await tester.pumpAndSettle();

      // Verify login screen is shown
      expect(find.byType(LoginScreen), findsOneWidget);
      expect(find.byType(HomeScreen), findsNothing);
    });

    testWidgets('App should show home screen when authenticated', (WidgetTester tester) async {
      final authProvider = AuthProvider();
      authProvider.authenticate('test_user_123');

      await tester.pumpWidget(
        MultiProvider(
          providers: [
            ChangeNotifierProvider.value(value: authProvider),
            ChangeNotifierProvider(create: (_) => ReceiptProvider()),
            ChangeNotifierProvider(create: (_) => WalletProvider()),
          ],
          child: const MaterialApp(home: AuthWrapper()),
        ),
      );

      await tester.pumpAndSettle();

      // Verify home screen is shown
      expect(find.byType(HomeScreen), findsOneWidget);
      expect(find.byType(LoginScreen), findsNothing);
    });
  });

  group('Login Screen Tests', () {
    testWidgets('Login screen should have email and password fields', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(home: const LoginScreen()),
      );

      expect(find.byType(TextField), findsNWidgets(2));
      expect(find.byType(ElevatedButton), findsOneWidget);
    });

    testWidgets('Login button should trigger authentication', (WidgetTester tester) async {
      final authProvider = AuthProvider();
      
      await tester.pumpWidget(
        MultiProvider(
          providers: [
            ChangeNotifierProvider.value(value: authProvider),
          ],
          child: MaterialApp(home: const LoginScreen()),
        ),
      );

      // Enter credentials
      await tester.enterText(find.byType(TextField).first, 'test@example.com');
      await tester.enterText(find.byType(TextField).last, 'password123');
      
      // Tap login button
      await tester.tap(find.byType(ElevatedButton));
      await tester.pumpAndSettle();

      // Verify authentication was attempted
      expect(authProvider.isLoading, isTrue);
    });
  });

  group('Home Screen Tests', () {
    testWidgets('Home screen should display quick action cards', (WidgetTester tester) async {
      final authProvider = AuthProvider();
      authProvider.authenticate('test_user_123');

      await tester.pumpWidget(
        MultiProvider(
          providers: [
            ChangeNotifierProvider.value(value: authProvider),
            ChangeNotifierProvider(create: (_) => ReceiptProvider()),
            ChangeNotifierProvider(create: (_) => WalletProvider()),
          ],
          child: MaterialApp(home: const HomeScreen()),
        ),
      );

      await tester.pumpAndSettle();

      // Verify quick action cards are present
      expect(find.byType(QuickActionCard), findsNWidgets(4));
      expect(find.text('Upload Receipt'), findsOneWidget);
      expect(find.text('Ask AI'), findsOneWidget);
      expect(find.text('View Passes'), findsOneWidget);
      expect(find.text('Analytics'), findsOneWidget);
    });

    testWidgets('Quick action cards should navigate to correct screens', (WidgetTester tester) async {
      final authProvider = AuthProvider();
      authProvider.authenticate('test_user_123');

      await tester.pumpWidget(
        MultiProvider(
          providers: [
            ChangeNotifierProvider.value(value: authProvider),
            ChangeNotifierProvider(create: (_) => ReceiptProvider()),
            ChangeNotifierProvider(create: (_) => WalletProvider()),
          ],
          child: MaterialApp(
            home: const HomeScreen(),
            routes: {
              '/receipts': (context) => const ReceiptsScreen(),
              '/passes': (context) => const PassesScreen(),
              '/analysis': (context) => const AnalysisScreen(),
            },
          ),
        ),
      );

      await tester.pumpAndSettle();

      // Tap on "View Passes" card
      await tester.tap(find.text('View Passes'));
      await tester.pumpAndSettle();

      // Verify navigation to passes screen
      expect(find.byType(PassesScreen), findsOneWidget);
    });
  });

  group('Receipts Screen Tests', () {
    testWidgets('Receipts screen should display receipt list', (WidgetTester tester) async {
      final receiptProvider = ReceiptProvider();
      
      // Mock receipt data
      final receipt = Receipt(
        id: 'test_receipt_123',
        userId: 'test_user_123',
        storeName: 'Test Store',
        totalAmount: 45.99,
        taxAmount: 3.50,
        items: [
          ReceiptItem(
            name: 'Milk',
            price: 4.99,
            quantity: 2,
            category: 'dairy',
          ),
        ],
        date: DateTime.now(),
        imageUrl: 'https://example.com/receipt.jpg',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );
      
      receiptProvider.receipts = [receipt];

      await tester.pumpWidget(
        MultiProvider(
          providers: [
            ChangeNotifierProvider.value(value: receiptProvider),
          ],
          child: MaterialApp(home: const ReceiptsScreen()),
        ),
      );

      await tester.pumpAndSettle();

      // Verify receipt is displayed
      expect(find.text('Test Store'), findsOneWidget);
      expect(find.text('\$45.99'), findsOneWidget);
    });

    testWidgets('Receipts screen should have upload button', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(home: const ReceiptsScreen()),
      );

      expect(find.byType(FloatingActionButton), findsOneWidget);
      expect(find.byIcon(Icons.add), findsOneWidget);
    });
  });

  group('Passes Screen Tests', () {
    testWidgets('Passes screen should display wallet passes', (WidgetTester tester) async {
      final walletProvider = WalletProvider();
      
      // Mock wallet pass data
      final pass = WalletPass(
        id: 'test_pass_123',
        userId: 'test_user_123',
        type: 'receipt',
        title: 'Receipt - Test Store',
        description: 'Total: \$45.99, Items: 2',
        data: '{"receipt_id": "test_receipt_123", "store_name": "Test Store"}',
        createdAt: DateTime.now(),
      );
      
      walletProvider.passes = [pass];

      await tester.pumpWidget(
        MultiProvider(
          providers: [
            ChangeNotifierProvider.value(value: walletProvider),
          ],
          child: MaterialApp(home: const PassesScreen()),
        ),
      );

      await tester.pumpAndSettle();

      // Verify pass is displayed
      expect(find.text('Receipt - Test Store'), findsOneWidget);
      expect(find.text('Total: \$45.99, Items: 2'), findsOneWidget);
    });
  });

  group('Analysis Screen Tests', () {
    testWidgets('Analysis screen should display spending statistics', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(home: const AnalysisScreen()),
      );

      await tester.pumpAndSettle();

      // Verify stats cards are present
      expect(find.byType(StatsCard), findsNWidgets(4));
      expect(find.text('Total Spent'), findsOneWidget);
      expect(find.text('Receipt Count'), findsOneWidget);
      expect(find.text('Average per Receipt'), findsOneWidget);
    });
  });

  group('Widget Component Tests', () {
    testWidgets('QuickActionCard should display title and icon', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: QuickActionCard(
              title: 'Test Action',
              icon: Icons.star,
              onTap: () {},
            ),
          ),
        ),
      );

      expect(find.text('Test Action'), findsOneWidget);
      expect(find.byIcon(Icons.star), findsOneWidget);
    });

    testWidgets('StatsCard should display value and label', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: StatsCard(
              value: '\$245.67',
              label: 'Total Spent',
              icon: Icons.attach_money,
            ),
          ),
        ),
      );

      expect(find.text('\$245.67'), findsOneWidget);
      expect(find.text('Total Spent'), findsOneWidget);
      expect(find.byIcon(Icons.attach_money), findsOneWidget);
    });
  });

  group('Provider Tests', () {
    test('AuthProvider should authenticate user', () {
      final authProvider = AuthProvider();
      
      expect(authProvider.isAuthenticated, isFalse);
      expect(authProvider.userId, isNull);
      
      authProvider.authenticate('test_user_123');
      
      expect(authProvider.isAuthenticated, isTrue);
      expect(authProvider.userId, equals('test_user_123'));
    });

    test('AuthProvider should logout user', () {
      final authProvider = AuthProvider();
      authProvider.authenticate('test_user_123');
      
      expect(authProvider.isAuthenticated, isTrue);
      
      authProvider.logout();
      
      expect(authProvider.isAuthenticated, isFalse);
      expect(authProvider.userId, isNull);
    });

    test('ReceiptProvider should add receipt', () {
      final receiptProvider = ReceiptProvider();
      
      expect(receiptProvider.receipts, isEmpty);
      
      final receipt = Receipt(
        id: 'test_receipt_123',
        userId: 'test_user_123',
        storeName: 'Test Store',
        totalAmount: 45.99,
        taxAmount: 3.50,
        items: [],
        date: DateTime.now(),
        imageUrl: 'https://example.com/receipt.jpg',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );
      
      receiptProvider.addReceipt(receipt);
      
      expect(receiptProvider.receipts, hasLength(1));
      expect(receiptProvider.receipts.first.id, equals('test_receipt_123'));
    });

    test('WalletProvider should add pass', () {
      final walletProvider = WalletProvider();
      
      expect(walletProvider.passes, isEmpty);
      
      final pass = WalletPass(
        id: 'test_pass_123',
        userId: 'test_user_123',
        type: 'receipt',
        title: 'Test Receipt',
        description: 'Test Description',
        data: '{}',
        createdAt: DateTime.now(),
      );
      
      walletProvider.addPass(pass);
      
      expect(walletProvider.passes, hasLength(1));
      expect(walletProvider.passes.first.id, equals('test_pass_123'));
    });
  });

  group('Integration Tests', () {
    testWidgets('Complete user flow: login -> upload receipt -> view analysis', (WidgetTester tester) async {
      final authProvider = AuthProvider();
      final receiptProvider = ReceiptProvider();
      final walletProvider = WalletProvider();

      // Start with login
      await tester.pumpWidget(
        MultiProvider(
          providers: [
            ChangeNotifierProvider.value(value: authProvider),
            ChangeNotifierProvider.value(value: receiptProvider),
            ChangeNotifierProvider.value(value: walletProvider),
          ],
          child: MaterialApp(
            home: const AuthWrapper(),
            routes: {
              '/home': (context) => const HomeScreen(),
              '/receipts': (context) => const ReceiptsScreen(),
              '/analysis': (context) => const AnalysisScreen(),
            },
          ),
        ),
      );

      await tester.pumpAndSettle();

      // Verify login screen
      expect(find.byType(LoginScreen), findsOneWidget);

      // Authenticate user
      authProvider.authenticate('test_user_123');
      await tester.pumpAndSettle();

      // Verify home screen
      expect(find.byType(HomeScreen), findsOneWidget);

      // Navigate to receipts
      await tester.tap(find.text('Upload Receipt'));
      await tester.pumpAndSettle();

      // Verify receipts screen
      expect(find.byType(ReceiptsScreen), findsOneWidget);

      // Navigate to analysis
      await tester.tap(find.text('Analytics'));
      await tester.pumpAndSettle();

      // Verify analysis screen
      expect(find.byType(AnalysisScreen), findsOneWidget);
    });
  });
} 