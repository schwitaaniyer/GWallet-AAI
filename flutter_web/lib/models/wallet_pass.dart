class WalletPass {
  final String id;
  final String userId;
  final String type;
  final String title;
  final String description;
  final String? imageUrl;
  final String? qrCode;
  final Map<String, dynamic> data;
  final DateTime createdAt;
  final DateTime? expiresAt;
  final bool isActive;

  WalletPass({
    required this.id,
    required this.userId,
    required this.type,
    required this.title,
    required this.description,
    this.imageUrl,
    this.qrCode,
    required this.data,
    required this.createdAt,
    this.expiresAt,
    required this.isActive,
  });

  factory WalletPass.fromJson(Map<String, dynamic> json) {
    return WalletPass(
      id: json['id'] ?? '',
      userId: json['user_id'] ?? '',
      type: json['type'] ?? '',
      title: json['title'] ?? '',
      description: json['description'] ?? '',
      imageUrl: json['image_url'],
      qrCode: json['qr_code'],
      data: json['data'] ?? {},
      createdAt: DateTime.parse(json['created_at'] ?? DateTime.now().toIso8601String()),
      expiresAt: json['expires_at'] != null
          ? DateTime.parse(json['expires_at'])
          : null,
      isActive: json['is_active'] ?? true,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'user_id': userId,
      'type': type,
      'title': title,
      'description': description,
      'image_url': imageUrl,
      'qr_code': qrCode,
      'data': data,
      'created_at': createdAt.toIso8601String(),
      'expires_at': expiresAt?.toIso8601String(),
      'is_active': isActive,
    };
  }

  bool get isExpired {
    if (expiresAt == null) return false;
    return DateTime.now().isAfter(expiresAt!);
  }

  WalletPass copyWith({
    String? id,
    String? userId,
    String? type,
    String? title,
    String? description,
    String? imageUrl,
    String? qrCode,
    Map<String, dynamic>? data,
    DateTime? createdAt,
    DateTime? expiresAt,
    bool? isActive,
  }) {
    return WalletPass(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      type: type ?? this.type,
      title: title ?? this.title,
      description: description ?? this.description,
      imageUrl: imageUrl ?? this.imageUrl,
      qrCode: qrCode ?? this.qrCode,
      data: data ?? this.data,
      createdAt: createdAt ?? this.createdAt,
      expiresAt: expiresAt ?? this.expiresAt,
      isActive: isActive ?? this.isActive,
    );
  }
} 