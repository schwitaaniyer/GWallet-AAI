import 'package:raseed_wallet/models/item.dart';
import 'package:raseed_wallet/models/location.dart';

class Receipt {
  final String id;
  final String userId;
  final String storeName;
  final double total;
  final double tax;
  final double subtotal;
  final DateTime date;
  final String imageUrl;
  final List<Item> items;
  final Location? location;
  final String status;
  final DateTime createdAt;
  final DateTime updatedAt;

  Receipt({
    required this.id,
    required this.userId,
    required this.storeName,
    required this.total,
    required this.tax,
    required this.subtotal,
    required this.date,
    required this.imageUrl,
    required this.items,
    this.location,
    required this.status,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Receipt.fromJson(Map<String, dynamic> json) {
    return Receipt(
      id: json['id'] ?? '',
      userId: json['user_id'] ?? '',
      storeName: json['store_name'] ?? '',
      total: (json['total'] ?? 0.0).toDouble(),
      tax: (json['tax'] ?? 0.0).toDouble(),
      subtotal: (json['subtotal'] ?? 0.0).toDouble(),
      date: DateTime.parse(json['date'] ?? DateTime.now().toIso8601String()),
      imageUrl: json['image_url'] ?? '',
      items: (json['items'] as List<dynamic>?)
              ?.map((item) => Item.fromJson(item))
              .toList() ??
          [],
      location: json['location'] != null
          ? Location.fromJson(json['location'])
          : null,
      status: json['status'] ?? 'pending',
      createdAt: DateTime.parse(json['created_at'] ?? DateTime.now().toIso8601String()),
      updatedAt: DateTime.parse(json['updated_at'] ?? DateTime.now().toIso8601String()),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'user_id': userId,
      'store_name': storeName,
      'total': total,
      'tax': tax,
      'subtotal': subtotal,
      'date': date.toIso8601String(),
      'image_url': imageUrl,
      'items': items.map((item) => item.toJson()).toList(),
      'location': location?.toJson(),
      'status': status,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }

  Receipt copyWith({
    String? id,
    String? userId,
    String? storeName,
    double? total,
    double? tax,
    double? subtotal,
    DateTime? date,
    String? imageUrl,
    List<Item>? items,
    Location? location,
    String? status,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Receipt(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      storeName: storeName ?? this.storeName,
      total: total ?? this.total,
      tax: tax ?? this.tax,
      subtotal: subtotal ?? this.subtotal,
      date: date ?? this.date,
      imageUrl: imageUrl ?? this.imageUrl,
      items: items ?? this.items,
      location: location ?? this.location,
      status: status ?? this.status,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
} 