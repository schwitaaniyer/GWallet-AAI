# Raseed Flutter Web Interface

A modern, responsive Flutter web application that provides a Google Wallet-like interface for Project Raseed AI Agent.

## Features

- **Modern UI/UX**: Material Design 3 with Google Wallet-inspired interface
- **Responsive Design**: Works seamlessly on desktop, tablet, and mobile
- **Real-time Data**: Live updates for receipts, passes, and analytics
- **AI Integration**: Direct access to AI query processing
- **Receipt Management**: Upload, view, and manage receipts
- **Wallet Passes**: View and manage Google Wallet passes
- **Spending Analytics**: Visual spending analysis and insights
- **Authentication**: Secure user authentication system

## Screenshots

- **Dashboard**: Overview with stats and quick actions
- **Receipts**: List and detailed view of all receipts
- **Passes**: Wallet passes management
- **Analytics**: Spending analysis and charts
- **Settings**: User preferences and configuration

## Technology Stack

- **Frontend**: Flutter Web
- **State Management**: Provider
- **HTTP Client**: http package
- **UI Components**: Material Design 3
- **Charts**: fl_chart
- **Icons**: Material Icons + Custom SVGs
- **Fonts**: Google Fonts (Inter)

## Prerequisites

- Flutter SDK 3.0.0 or later
- Dart SDK 2.17.0 or later
- Chrome browser (for web development)
- Access to deployed Raseed backend

## Installation

### 1. Clone the Repository

```bash
git clone <repository-url>
cd flutter_web
```

### 2. Install Dependencies

```bash
flutter pub get
```

### 3. Configure Backend URL

Update the backend URL in `lib/utils/constants.dart`:

```dart
class ApiEndpoints {
  static const String baseUrl = 'https://your-backend-url.com';
  // ... other endpoints
}
```

### 4. Run the Application

```bash
# For development
flutter run -d chrome

# For production build
flutter build web --release
```

## Project Structure

```
lib/
├── main.dart                 # App entry point
├── models/                   # Data models
│   ├── receipt.dart
│   ├── item.dart
│   ├── location.dart
│   ├── wallet_pass.dart
│   └── query.dart
├── providers/                # State management
│   ├── auth_provider.dart
│   ├── receipt_provider.dart
│   └── wallet_provider.dart
├── screens/                  # UI screens
│   ├── login_screen.dart
│   ├── home_screen.dart
│   ├── dashboard_screen.dart
│   ├── receipts_screen.dart
│   ├── passes_screen.dart
│   ├── analysis_screen.dart
│   └── settings_screen.dart
├── widgets/                  # Reusable widgets
│   ├── stats_card.dart
│   └── quick_action_card.dart
└── utils/                    # Utilities
    └── constants.dart
```

## Key Features Implementation

### 1. Authentication

The app uses a mock authentication system for development:

```dart
// Demo login credentials
Email: demo@raseed.com
Password: demo123
```

### 2. Receipt Upload

- Camera integration for photo capture
- Gallery selection for existing images
- Real-time upload progress
- AI processing status updates

### 3. AI Query Processing

- Natural language query submission
- Real-time response processing
- Query history and management
- Intent detection and categorization

### 4. Wallet Passes

- Google Wallet pass creation
- QR code generation
- Pass management and organization
- Expiry notifications

### 5. Analytics Dashboard

- Spending overview
- Category-wise breakdown
- Monthly trends
- Savings insights

## API Integration

The app communicates with the Raseed backend through these endpoints:

```dart
// Health check
GET /health

// Receipts
POST /receipts (upload)
GET /receipts (list)
DELETE /receipts/{id}

// Queries
POST /queries (submit)
GET /queries (list)

// Wallet passes
GET /wallet-passes (list)

// Analytics
GET /analysis (spending analysis)
```

## State Management

The app uses Provider for state management:

```dart
// AuthProvider - User authentication
// ReceiptProvider - Receipt data management
// WalletProvider - Wallet passes and queries
```

## Styling and Theming

### Color Scheme

```dart
class AppColors {
  static const Color primary = Color(0xFF1A73E8);    // Google Blue
  static const Color secondary = Color(0xFF34A853);  // Google Green
  static const Color accent = Color(0xFFEA4335);     // Google Red
  static const Color warning = Color(0xFFFBBC04);    // Google Yellow
}
```

### Typography

- **Primary Font**: Inter (Google Fonts)
- **Headings**: Bold weights for hierarchy
- **Body Text**: Regular weight for readability

## Development

### Adding New Screens

1. Create screen file in `lib/screens/`
2. Add route in `main.dart`
3. Update navigation in `home_screen.dart`

### Adding New Widgets

1. Create widget file in `lib/widgets/`
2. Export from `lib/widgets/widgets.dart`
3. Import and use in screens

### Testing

```bash
# Run unit tests
flutter test

# Run integration tests
flutter test integration_test/
```

## Deployment

### Local Build

```bash
flutter build web --release
```

### Docker Deployment

```bash
# Build Docker image
docker build -t raseed-web .

# Run locally
docker run -p 8080:8080 raseed-web

# Deploy to Cloud Run
gcloud builds submit --tag gcr.io/PROJECT_ID/raseed-web
gcloud run deploy raseed-web \
    --image gcr.io/PROJECT_ID/raseed-web \
    --platform managed \
    --region us-central1 \
    --allow-unauthenticated
```

## Performance Optimization

### Web Optimization

- **Tree Shaking**: Removes unused code
- **Code Splitting**: Lazy loading of routes
- **Asset Optimization**: Compressed images and fonts
- **Caching**: Browser and CDN caching

### Bundle Size

- **Initial Load**: ~2MB (compressed)
- **Runtime**: Minimal memory usage
- **Caching**: Aggressive caching strategy

## Browser Support

- **Chrome**: 90+
- **Firefox**: 88+
- **Safari**: 14+
- **Edge**: 90+

## Security

- **HTTPS Only**: Secure communication
- **CORS**: Proper cross-origin configuration
- **Input Validation**: Client-side validation
- **XSS Protection**: Content Security Policy

## Monitoring

### Error Tracking

- Console logging for development
- Error reporting for production
- Performance monitoring

### Analytics

- User interaction tracking
- Performance metrics
- Error rate monitoring

## Troubleshooting

### Common Issues

1. **Build Errors**: Check Flutter version compatibility
2. **API Errors**: Verify backend URL and connectivity
3. **Performance**: Check bundle size and caching
4. **Styling**: Verify Material Design theme setup

### Debug Mode

```bash
flutter run -d chrome --web-renderer html
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is part of Project Raseed and follows the same license terms.

## Support

For support and questions:
- Create an issue in the repository
- Check the documentation
- Review the API documentation 