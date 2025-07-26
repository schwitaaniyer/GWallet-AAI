# GWallet-AAI Testing Strategy

This document outlines the comprehensive testing strategy for the GWallet-AAI project, covering both CLI and web interface components.

## Overview

The testing strategy is designed to ensure reliability, performance, and user experience across all components of the GWallet-AAI system. It includes:

- **Unit Tests**: Testing individual components in isolation
- **Integration Tests**: Testing component interactions
- **End-to-End Tests**: Testing complete user workflows
- **Performance Tests**: Testing system performance under load
- **Web Interface Tests**: Testing Flutter web components

## Test Structure

```
tests/
├── basic_test.go           # Basic unit tests
├── integration_test.go     # Integration tests
├── e2e_test.go            # End-to-end tests
├── cli_test.go            # CLI-specific tests
├── test_config.yaml       # Test configuration
├── run_tests.sh           # Test runner script (Linux/macOS)
├── run_tests.bat          # Test runner script (Windows)
└── README.md              # This file
```

## Test Types

### 1. Unit Tests

Unit tests focus on testing individual functions and components in isolation.

#### CLI Unit Tests (`cli/cli_test.go`)
- **Purpose**: Test individual CLI commands and functions
- **Coverage**: All CLI commands (health, upload-receipt, query, etc.)
- **Mocking**: HTTP responses using `httptest.Server`
- **Execution**: `go test ./cli/`

**Example Test Cases:**
- Health command returns correct status
- Receipt upload with valid file
- Receipt upload with non-existent file
- Query submission and response parsing
- Error handling for network failures

#### Flutter Unit Tests (`flutter_web/test/widget_test.dart`)
- **Purpose**: Test individual Flutter widgets and providers
- **Coverage**: UI components, state management, data models
- **Mocking**: Provider state and HTTP responses
- **Execution**: `flutter test`

**Example Test Cases:**
- Authentication flow
- Receipt list display
- Wallet pass creation
- Analytics data visualization
- Error state handling

### 2. Integration Tests

Integration tests verify that components work together correctly.

#### Backend Integration Tests (`tests/integration_test.go`)
- **Purpose**: Test API endpoints and data flow
- **Coverage**: All REST API endpoints
- **Mocking**: Database and external services
- **Execution**: `go test -tags=integration`

**Test Scenarios:**
- Receipt upload and processing pipeline
- Query submission and AI response
- Wallet pass creation and retrieval
- Data consistency across endpoints
- Error handling and recovery

### 3. End-to-End Tests

E2E tests simulate real user workflows across the entire system.

#### Complete Workflow Tests (`tests/e2e_test.go`)
- **Purpose**: Test complete user journeys
- **Coverage**: CLI to web interface data flow
- **Mocking**: Full backend simulation
- **Execution**: `go test -tags=e2e`

**Test Scenarios:**
- User uploads receipt via CLI, views in web interface
- User submits query via CLI, sees response in web
- Data consistency between CLI and web components
- Cross-component authentication
- Performance under realistic load

### 4. Performance Tests

Performance tests measure system behavior under various load conditions.

#### Benchmark Tests
- **Purpose**: Measure response times and throughput
- **Coverage**: Critical operations (receipt upload, query processing)
- **Execution**: `go test -bench=. -benchmem`

**Benchmarks:**
- Receipt upload performance
- Query processing latency
- Concurrent user handling
- Memory usage patterns

## Running Tests

### Prerequisites

1. **Go**: Version 1.21 or higher
2. **Flutter**: Version 3.0 or higher
3. **Node.js**: Version 18 or higher (optional, for web testing tools)

### Using the Test Runner Script

The test runner scripts provide a convenient way to run different types of tests:

#### On Linux/macOS (Bash):
```bash
# Make the script executable (first time only)
chmod +x tests/run_tests.sh

# Run all tests
./tests/run_tests.sh

# Run specific test types
./tests/run_tests.sh unit
./tests/run_tests.sh integration
./tests/run_tests.sh e2e
./tests/run_tests.sh performance

# Run with options
./tests/run_tests.sh -c -r all  # Clean, run all tests, generate report
```

#### On Windows:
```cmd
# Using the Windows batch file
tests\run_tests.bat

# Run specific test types
tests\run_tests.bat unit
tests\run_tests.bat integration
tests\run_tests.bat e2e
tests\run_tests.bat performance

# Run with options
tests\run_tests.bat -c -r all  # Clean, run all tests, generate report
```

#### On Windows with Git Bash:
```bash
# Make the script executable (first time only)
chmod +x tests/run_tests.sh

# Run all tests
./tests/run_tests.sh

# Same options as Linux/macOS
```

### Manual Test Execution

#### CLI Tests
```bash
# Unit tests
cd cli
go test -v -coverprofile=coverage.out ./...

# Integration tests
cd tests
go test -v -tags=integration ./...

# E2E tests
go test -v -tags=e2e ./...

# Performance tests
go test -v -bench=. -benchmem ./...
```

#### Flutter Tests
```bash
# Unit tests
cd flutter_web
flutter test --coverage

# Widget tests
flutter test test/widget_test.dart

# Integration tests (requires ChromeDriver)
flutter drive --driver=test_driver/integration_test.dart --target=integration_test/app_test.dart
```

## Test Configuration

The `test_config.yaml` file contains:

- **Test Environments**: Local, staging, production configurations
- **Mock Data**: Predefined test data for consistent testing
- **Test Scenarios**: Step-by-step test workflows
- **Performance Parameters**: Load testing configurations
- **Error Scenarios**: Expected error handling behavior

## Test Data Management

### Mock Data
- **Users**: Test user accounts with different permission levels
- **Receipts**: Sample receipt images and data
- **Queries**: Predefined AI queries and expected responses
- **Wallet Passes**: Sample wallet pass configurations

### Test Files
- **Receipt Images**: Various sizes and formats for testing
- **Invalid Files**: Files that should trigger error handling
- **Large Files**: Performance testing with large uploads

## Error Testing

The testing strategy includes comprehensive error scenario testing:

### Network Errors
- Connection timeouts
- Server unavailability
- Network interruptions

### Validation Errors
- Invalid file formats
- Missing required fields
- Malformed data

### Authentication Errors
- Invalid credentials
- Expired tokens
- Permission denied

## Test Reporting

### Coverage Reports
- **CLI Coverage**: HTML coverage report for Go code
- **Flutter Coverage**: LCOV coverage report for Dart code
- **Coverage Thresholds**: Minimum coverage requirements

### Test Reports
- **HTML Reports**: Comprehensive test execution reports
- **Performance Metrics**: Response times and throughput data
- **Error Logs**: Detailed error information for debugging

### Continuous Integration
- **Automated Testing**: Tests run on every commit
- **Coverage Tracking**: Coverage trends over time
- **Performance Monitoring**: Performance regression detection

## Best Practices

### Test Organization
1. **Group Related Tests**: Use test groups for related functionality
2. **Descriptive Names**: Use clear, descriptive test names
3. **Setup and Teardown**: Properly clean up test resources
4. **Mock External Dependencies**: Avoid external service calls in tests

### Test Data
1. **Consistent Data**: Use predefined test data for consistency
2. **Isolation**: Each test should be independent
3. **Cleanup**: Remove test data after each test
4. **Realistic Data**: Use realistic but anonymized test data

### Performance Testing
1. **Baseline Measurements**: Establish performance baselines
2. **Load Testing**: Test under realistic load conditions
3. **Resource Monitoring**: Monitor CPU, memory, and network usage
4. **Regression Detection**: Alert on performance regressions

## Troubleshooting

### Common Issues

#### CLI Tests Fail
- Check Go version compatibility
- Verify network connectivity for integration tests
- Ensure test data files exist

#### Flutter Tests Fail
- Check Flutter version compatibility
- Verify web dependencies are installed
- Check for ChromeDriver availability

#### Performance Tests Fail
- Check system resources
- Verify test environment configuration
- Review performance thresholds

### Debugging Tips
1. **Verbose Output**: Use `-v` flag for detailed test output
2. **Single Test**: Run individual tests for focused debugging
3. **Log Analysis**: Review test logs for error details
4. **Environment Check**: Verify test environment setup

## Future Enhancements

### Planned Improvements
1. **Visual Regression Testing**: Automated UI comparison testing
2. **Accessibility Testing**: Automated accessibility compliance testing
3. **Security Testing**: Automated security vulnerability scanning
4. **Load Testing**: More sophisticated load testing scenarios

### Test Automation
1. **CI/CD Integration**: Automated test execution in pipelines
2. **Test Parallelization**: Parallel test execution for faster feedback
3. **Test Data Management**: Automated test data generation and cleanup
4. **Monitoring Integration**: Real-time test result monitoring

## Contributing

When adding new tests:

1. **Follow Naming Conventions**: Use descriptive test names
2. **Add Documentation**: Document new test scenarios
3. **Update Configuration**: Add new test data to config files
4. **Maintain Coverage**: Ensure adequate test coverage
5. **Performance Impact**: Consider test execution time

## Support

For questions or issues with testing:

1. **Check Documentation**: Review this README and inline comments
2. **Review Test Logs**: Check test output for error details
3. **Verify Environment**: Ensure all dependencies are installed
4. **Contact Team**: Reach out to the development team for assistance 