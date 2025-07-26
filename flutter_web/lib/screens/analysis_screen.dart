import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:raseed_wallet/providers/wallet_provider.dart';
import 'package:raseed_wallet/utils/constants.dart';

class AnalysisScreen extends StatelessWidget {
  const AnalysisScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('Spending Analysis'),
      ),
      body: Consumer<WalletProvider>(
        builder: (context, walletProvider, child) {
          if (walletProvider.isLoading) {
            return const Center(child: CircularProgressIndicator());
          }

          return SingleChildScrollView(
            padding: const EdgeInsets.all(AppSizes.paddingM),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Card(
                  child: Padding(
                    padding: const EdgeInsets.all(AppSizes.paddingL),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Monthly Overview',
                          style: Theme.of(context).textTheme.titleLarge?.copyWith(
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: AppSizes.paddingM),
                        Row(
                          children: [
                            Expanded(
                              child: _buildStatItem('Total Spent', '\$1,234.56', AppColors.primary),
                            ),
                            const SizedBox(width: AppSizes.paddingM),
                            Expanded(
                              child: _buildStatItem('Avg/Day', '\$41.15', AppColors.secondary),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
                const SizedBox(height: AppSizes.paddingL),
                Text(
                  'Categories',
                  style: Theme.of(context).textTheme.titleLarge?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: AppSizes.paddingM),
                Card(
                  child: Padding(
                    padding: const EdgeInsets.all(AppSizes.paddingL),
                    child: Column(
                      children: [
                        _buildCategoryItem('Food & Dining', 450.00, 0.36, AppColors.primary),
                        _buildCategoryItem('Shopping', 320.00, 0.26, AppColors.secondary),
                        _buildCategoryItem('Transportation', 280.00, 0.23, AppColors.accent),
                        _buildCategoryItem('Entertainment', 184.56, 0.15, AppColors.warning),
                      ],
                    ),
                  ),
                ),
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _buildStatItem(String title, String value, Color color) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          title,
          style: TextStyle(color: AppColors.textSecondary),
        ),
        const SizedBox(height: AppSizes.paddingXS),
        Text(
          value,
          style: TextStyle(
            fontSize: 24,
            fontWeight: FontWeight.bold,
            color: color,
          ),
        ),
      ],
    );
  }

  Widget _buildCategoryItem(String category, double amount, double percentage, Color color) {
    return Padding(
      padding: const EdgeInsets.only(bottom: AppSizes.paddingM),
      child: Row(
        children: [
          Container(
            width: 12,
            height: 12,
            decoration: BoxDecoration(
              color: color,
              borderRadius: BorderRadius.circular(6),
            ),
          ),
          const SizedBox(width: AppSizes.paddingM),
          Expanded(
            child: Text(category, style: const TextStyle(fontWeight: FontWeight.w500)),
          ),
          Text(
            '\$${amount.toStringAsFixed(2)}',
            style: const TextStyle(fontWeight: FontWeight.w600),
          ),
          const SizedBox(width: AppSizes.paddingM),
          Text(
            '${(percentage * 100).toInt()}%',
            style: TextStyle(color: AppColors.textSecondary),
          ),
        ],
      ),
    );
  }
} 