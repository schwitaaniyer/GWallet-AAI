import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:intl/intl.dart';
import 'package:raseed_wallet/providers/auth_provider.dart';
import 'package:raseed_wallet/providers/receipt_provider.dart';
import 'package:raseed_wallet/providers/wallet_provider.dart';
import 'package:raseed_wallet/utils/constants.dart';
import 'package:raseed_wallet/widgets/quick_action_card.dart';
import 'package:raseed_wallet/widgets/stats_card.dart';

class DashboardScreen extends StatefulWidget {
  const DashboardScreen({super.key});

  @override
  State<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends State<DashboardScreen> {
  final TextEditingController _queryController = TextEditingController();

  @override
  void dispose() {
    _queryController.dispose();
    super.dispose();
  }

  void _submitQuery() {
    if (_queryController.text.trim().isEmpty) return;

    final authProvider = context.read<AuthProvider>();
    final walletProvider = context.read<WalletProvider>();

    if (authProvider.userId != null) {
      walletProvider.submitQuery(
        authProvider.userId!,
        _queryController.text.trim(),
        'en',
      );
      _queryController.clear();
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Query submitted! Check your passes for the response.')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      body: CustomScrollView(
        slivers: [
          // App Bar
          SliverAppBar(
            expandedHeight: 120,
            floating: false,
            pinned: true,
            backgroundColor: AppColors.primary,
            flexibleSpace: FlexibleSpaceBar(
              title: Consumer<AuthProvider>(
                builder: (context, authProvider, child) {
                  return Text(
                    'Welcome, ${authProvider.userName ?? 'User'}!',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 18,
                      fontWeight: FontWeight.w600,
                    ),
                  );
                },
              ),
              background: Container(
                decoration: const BoxDecoration(
                  gradient: LinearGradient(
                    begin: Alignment.topLeft,
                    end: Alignment.bottomRight,
                    colors: [AppColors.primary, Color(0xFF0D47A1)],
                  ),
                ),
              ),
            ),
            actions: [
              IconButton(
                icon: const Icon(Icons.notifications, color: Colors.white),
                onPressed: () {
                  // TODO: Show notifications
                },
              ),
              IconButton(
                icon: const Icon(Icons.person, color: Colors.white),
                onPressed: () {
                  // TODO: Show profile
                },
              ),
            ],
          ),

          // Content
          SliverPadding(
            padding: const EdgeInsets.all(AppSizes.paddingM),
            sliver: SliverList(
              delegate: SliverChildListDelegate([
                // Stats Cards
                Consumer<ReceiptProvider>(
                  builder: (context, receiptProvider, child) {
                    final receipts = receiptProvider.receipts;
                    final totalSpent = receipts.fold<double>(
                      0,
                      (sum, receipt) => sum + receipt.total,
                    );
                    final thisMonth = receipts.where((r) =>
                        r.date.month == DateTime.now().month &&
                        r.date.year == DateTime.now().year).toList();
                    final thisMonthSpent = thisMonth.fold<double>(
                      0,
                      (sum, receipt) => sum + receipt.total,
                    );

                    return Column(
                      children: [
                        Row(
                          children: [
                            Expanded(
                              child: StatsCard(
                                title: 'Total Spent',
                                value: '\$${totalSpent.toStringAsFixed(2)}',
                                icon: Icons.account_balance_wallet,
                                color: AppColors.primary,
                              ),
                            ),
                            const SizedBox(width: AppSizes.paddingM),
                            Expanded(
                              child: StatsCard(
                                title: 'This Month',
                                value: '\$${thisMonthSpent.toStringAsFixed(2)}',
                                icon: Icons.calendar_today,
                                color: AppColors.secondary,
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: AppSizes.paddingM),
                        Row(
                          children: [
                            Expanded(
                              child: StatsCard(
                                title: 'Receipts',
                                value: receipts.length.toString(),
                                icon: Icons.receipt_long,
                                color: AppColors.accent,
                              ),
                            ),
                            const SizedBox(width: AppSizes.paddingM),
                            Expanded(
                              child: Consumer<WalletProvider>(
                                builder: (context, walletProvider, child) {
                                  return StatsCard(
                                    title: 'Passes',
                                    value: walletProvider.walletPasses.length.toString(),
                                    icon: Icons.qr_code,
                                    color: AppColors.warning,
                                  );
                                },
                              ),
                            ),
                          ],
                        ),
                      ],
                    );
                  },
                ),

                const SizedBox(height: AppSizes.paddingL),

                // Quick Actions
                Text(
                  'Quick Actions',
                  style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: AppColors.textPrimary,
                  ),
                ),
                const SizedBox(height: AppSizes.paddingM),
                GridView.count(
                  shrinkWrap: true,
                  physics: const NeverScrollableScrollPhysics(),
                  crossAxisCount: 2,
                  crossAxisSpacing: AppSizes.paddingM,
                  mainAxisSpacing: AppSizes.paddingM,
                  childAspectRatio: 1.2,
                  children: [
                    QuickActionCard(
                      title: 'Upload Receipt',
                      icon: Icons.camera_alt,
                      color: AppColors.primary,
                      onTap: () {
                        // TODO: Navigate to upload screen
                      },
                    ),
                    QuickActionCard(
                      title: 'Ask AI',
                      icon: Icons.chat,
                      color: AppColors.secondary,
                      onTap: () {
                        _showQueryDialog(context);
                      },
                    ),
                    QuickActionCard(
                      title: 'View Passes',
                      icon: Icons.account_balance_wallet,
                      color: AppColors.accent,
                      onTap: () {
                        // TODO: Navigate to passes screen
                      },
                    ),
                    QuickActionCard(
                      title: 'Analytics',
                      icon: Icons.analytics,
                      color: AppColors.warning,
                      onTap: () {
                        // TODO: Navigate to analytics screen
                      },
                    ),
                  ],
                ),

                const SizedBox(height: AppSizes.paddingL),

                // Recent Activity
                Text(
                  'Recent Activity',
                  style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: AppColors.textPrimary,
                  ),
                ),
                const SizedBox(height: AppSizes.paddingM),
                Consumer<ReceiptProvider>(
                  builder: (context, receiptProvider, child) {
                    final recentReceipts = receiptProvider.receipts
                        .take(3)
                        .toList();

                    if (recentReceipts.isEmpty) {
                      return Card(
                        child: Padding(
                          padding: const EdgeInsets.all(AppSizes.paddingL),
                          child: Column(
                            children: [
                              Icon(
                                Icons.receipt_long,
                                size: 48,
                                color: AppColors.textTertiary,
                              ),
                              const SizedBox(height: AppSizes.paddingM),
                              Text(
                                'No receipts yet',
                                style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                                  color: AppColors.textSecondary,
                                ),
                              ),
                              const SizedBox(height: AppSizes.paddingS),
                              Text(
                                'Upload your first receipt to get started',
                                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                                  color: AppColors.textTertiary,
                                ),
                                textAlign: TextAlign.center,
                              ),
                            ],
                          ),
                        ),
                      );
                    }

                    return Column(
                      children: recentReceipts.map((receipt) {
                        return Card(
                          margin: const EdgeInsets.only(bottom: AppSizes.paddingS),
                          child: ListTile(
                            leading: CircleAvatar(
                              backgroundColor: AppColors.primary.withOpacity(0.1),
                              child: Icon(
                                Icons.receipt_long,
                                color: AppColors.primary,
                              ),
                            ),
                            title: Text(
                              receipt.storeName,
                              style: const TextStyle(fontWeight: FontWeight.w600),
                            ),
                            subtitle: Text(
                              DateFormat('MMM dd, yyyy').format(receipt.date),
                              style: TextStyle(color: AppColors.textSecondary),
                            ),
                            trailing: Text(
                              '\$${receipt.total.toStringAsFixed(2)}',
                              style: const TextStyle(
                                fontWeight: FontWeight.bold,
                                fontSize: 16,
                              ),
                            ),
                            onTap: () {
                              // TODO: Navigate to receipt details
                            },
                          ),
                        );
                      }).toList(),
                    );
                  },
                ),
              ]),
            ),
          ),
        ],
      ),
    );
  }

  void _showQueryDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Ask AI Assistant'),
        content: TextField(
          controller: _queryController,
          decoration: const InputDecoration(
            hintText: 'Ask me anything about your receipts...',
            border: OutlineInputBorder(),
          ),
          maxLines: 3,
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              _submitQuery();
              Navigator.pop(context);
            },
            child: const Text('Ask'),
          ),
        ],
      ),
    );
  }
} 