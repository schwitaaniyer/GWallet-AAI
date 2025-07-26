import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:raseed_wallet/providers/auth_provider.dart';
import 'package:raseed_wallet/providers/receipt_provider.dart';
import 'package:raseed_wallet/providers/wallet_provider.dart';
import 'package:raseed_wallet/screens/dashboard_screen.dart';
import 'package:raseed_wallet/screens/receipts_screen.dart';
import 'package:raseed_wallet/screens/passes_screen.dart';
import 'package:raseed_wallet/screens/analysis_screen.dart';
import 'package:raseed_wallet/screens/settings_screen.dart';
import 'package:raseed_wallet/utils/constants.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  int _currentIndex = 0;
  late PageController _pageController;

  final List<Widget> _screens = [
    const DashboardScreen(),
    const ReceiptsScreen(),
    const PassesScreen(),
    const AnalysisScreen(),
    const SettingsScreen(),
  ];

  @override
  void initState() {
    super.initState();
    _pageController = PageController();
    _loadData();
  }

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  void _loadData() {
    final authProvider = context.read<AuthProvider>();
    final receiptProvider = context.read<ReceiptProvider>();
    final walletProvider = context.read<WalletProvider>();

    if (authProvider.userId != null) {
      receiptProvider.loadReceipts(authProvider.userId!);
      walletProvider.loadWalletPasses(authProvider.userId!);
      walletProvider.loadQueries(authProvider.userId!);
    }
  }

  void _onTabTapped(int index) {
    setState(() {
      _currentIndex = index;
    });
    _pageController.animateToPage(
      index,
      duration: const Duration(milliseconds: 300),
      curve: Curves.easeInOut,
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: PageView(
        controller: _pageController,
        onPageChanged: (index) {
          setState(() {
            _currentIndex = index;
          });
        },
        children: _screens,
      ),
      bottomNavigationBar: Container(
        decoration: BoxDecoration(
          boxShadow: [
            BoxShadow(
              color: AppColors.shadow,
              blurRadius: 8,
              offset: const Offset(0, -2),
            ),
          ],
        ),
        child: BottomNavigationBar(
          type: BottomNavigationBarType.fixed,
          currentIndex: _currentIndex,
          onTap: _onTabTapped,
          selectedItemColor: AppColors.primary,
          unselectedItemColor: AppColors.textTertiary,
          backgroundColor: AppColors.surface,
          elevation: 0,
          items: const [
            BottomNavigationBarItem(
              icon: Icon(Icons.dashboard),
              label: AppStrings.home,
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.receipt_long),
              label: AppStrings.receipts,
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.account_balance_wallet),
              label: AppStrings.passes,
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.analytics),
              label: AppStrings.analysis,
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.settings),
              label: AppStrings.settings,
            ),
          ],
        ),
      ),
      floatingActionButton: _currentIndex == 1 ? FloatingActionButton(
        onPressed: () {
          _showUploadOptions(context);
        },
        backgroundColor: AppColors.primary,
        child: const Icon(Icons.add, color: Colors.white),
      ) : null,
    );
  }

  void _showUploadOptions(BuildContext context) {
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => Container(
        padding: const EdgeInsets.all(AppSizes.paddingL),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: AppColors.border,
                borderRadius: BorderRadius.circular(2),
              ),
            ),
            const SizedBox(height: AppSizes.paddingL),
            Text(
              'Upload Receipt',
              style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: AppSizes.paddingL),
            ListTile(
              leading: const Icon(Icons.camera_alt, color: AppColors.primary),
              title: const Text('Take Photo'),
              onTap: () {
                Navigator.pop(context);
                // TODO: Implement camera functionality
              },
            ),
            ListTile(
              leading: const Icon(Icons.photo_library, color: AppColors.primary),
              title: const Text('Choose from Gallery'),
              onTap: () {
                Navigator.pop(context);
                // TODO: Implement gallery functionality
              },
            ),
            const SizedBox(height: AppSizes.paddingM),
          ],
        ),
      ),
    );
  }
} 