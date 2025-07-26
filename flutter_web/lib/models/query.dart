class Query {
  final String id;
  final String userId;
  final String query;
  final String language;
  final String? response;
  final String? intent;
  final Map<String, dynamic>? metadata;
  final DateTime createdAt;
  final DateTime? respondedAt;
  final String status;

  Query({
    required this.id,
    required this.userId,
    required this.query,
    required this.language,
    this.response,
    this.intent,
    this.metadata,
    required this.createdAt,
    this.respondedAt,
    required this.status,
  });

  factory Query.fromJson(Map<String, dynamic> json) {
    return Query(
      id: json['id'] ?? '',
      userId: json['user_id'] ?? '',
      query: json['query'] ?? '',
      language: json['language'] ?? 'en',
      response: json['response'],
      intent: json['intent'],
      metadata: json['metadata'],
      createdAt: DateTime.parse(json['created_at'] ?? DateTime.now().toIso8601String()),
      respondedAt: json['responded_at'] != null
          ? DateTime.parse(json['responded_at'])
          : null,
      status: json['status'] ?? 'pending',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'user_id': userId,
      'query': query,
      'language': language,
      'response': response,
      'intent': intent,
      'metadata': metadata,
      'created_at': createdAt.toIso8601String(),
      'responded_at': respondedAt?.toIso8601String(),
      'status': status,
    };
  }

  bool get isAnswered => response != null && response!.isNotEmpty;
  bool get isPending => status == 'pending';
  bool get isProcessing => status == 'processing';

  Query copyWith({
    String? id,
    String? userId,
    String? query,
    String? language,
    String? response,
    String? intent,
    Map<String, dynamic>? metadata,
    DateTime? createdAt,
    DateTime? respondedAt,
    String? status,
  }) {
    return Query(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      query: query ?? this.query,
      language: language ?? this.language,
      response: response ?? this.response,
      intent: intent ?? this.intent,
      metadata: metadata ?? this.metadata,
      createdAt: createdAt ?? this.createdAt,
      respondedAt: respondedAt ?? this.respondedAt,
      status: status ?? this.status,
    );
  }
} 