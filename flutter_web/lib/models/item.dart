class Item {
  final String name;
  final double price;
  final int quantity;
  final String? category;
  final String? brand;

  Item({
    required this.name,
    required this.price,
    required this.quantity,
    this.category,
    this.brand,
  });

  factory Item.fromJson(Map<String, dynamic> json) {
    return Item(
      name: json['name'] ?? '',
      price: (json['price'] ?? 0.0).toDouble(),
      quantity: json['quantity'] ?? 1,
      category: json['category'],
      brand: json['brand'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'name': name,
      'price': price,
      'quantity': quantity,
      'category': category,
      'brand': brand,
    };
  }

  double get totalPrice => price * quantity;

  Item copyWith({
    String? name,
    double? price,
    int? quantity,
    String? category,
    String? brand,
  }) {
    return Item(
      name: name ?? this.name,
      price: price ?? this.price,
      quantity: quantity ?? this.quantity,
      category: category ?? this.category,
      brand: brand ?? this.brand,
    );
  }
} 