import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:raseed_wallet/providers/wallet_provider.dart';
import 'package:raseed_wallet/utils/constants.dart';

class PassesScreen extends StatelessWidget {
  const PassesScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('Wallet Passes'),
      ),
      body: Consumer<WalletProvider>(
        builder: (context, walletProvider, child) {
          if (walletProvider.isLoading) {
            return const Center(child: CircularProgressIndicator());
          }

          if (walletProvider.walletPasses.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.account_balance_wallet, size: 64, color: AppColors.textTertiary),
                  const SizedBox(height: AppSizes.paddingM),
                  Text('No wallet passes yet', style: Theme.of(context).textTheme.headlineSmall),
                ],
              ),
            );
          }

          return ListView.builder(
            padding: const EdgeInsets.all(AppSizes.paddingM),
            itemCount: walletProvider.walletPasses.length,
            itemBuilder: (context, index) {
              final pass = walletProvider.walletPasses[index];
              return Card(
                margin: const EdgeInsets.only(bottom: AppSizes.paddingM),
                child: ListTile(
                  leading: CircleAvatar(
                    backgroundColor: AppColors.primary.withOpacity(0.1),
                    child: Icon(Icons.qr_code, color: AppColors.primary),
                  ),
                  title: Text(pass.title, style: const TextStyle(fontWeight: FontWeight.w600)),
                  subtitle: Text(pass.description),
                  trailing: Icon(Icons.chevron_right),
                  onTap: () {
                    // TODO: Show pass details
                  },
                ),
              );
            },
          );
        },
      ),
    );
  }
} 