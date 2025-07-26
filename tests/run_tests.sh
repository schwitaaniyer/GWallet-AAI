#!/bin/bash

# Test Runner Script for GWallet-AAI
# This script runs different types of tests for CLI and web interface components

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TEST_DIR="$PROJECT_ROOT/tests"
CLI_DIR="$PROJECT_ROOT/cli"
FLUTTER_WEB_DIR="$PROJECT_ROOT/flutter_web"
REPORTS_DIR="$PROJECT_ROOT/test_reports"

# Test types
TEST_TYPES=("unit" "integration" "e2e" "performance" "all")

# Functions
print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# Create reports directory
create_reports_dir() {
    if [ ! -d "$REPORTS_DIR" ]; then
        mkdir -p "$REPORTS_DIR"
        print_info "Created reports directory: $REPORTS_DIR"
    fi
}

# Check dependencies
check_dependencies() {
    print_header "Checking Dependencies"
    
    # Check Go
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed"
        exit 1
    fi
    print_success "Go is installed: $(go version)"
    
    # Check Flutter
    if ! command -v flutter &> /dev/null; then
        print_error "Flutter is not installed"
        exit 1
    fi
    print_success "Flutter is installed: $(flutter --version | head -n 1)"
    
    # Check Node.js (for web testing tools)
    if ! command -v node &> /dev/null; then
        print_warning "Node.js is not installed (optional for web testing)"
    else
        print_success "Node.js is installed: $(node --version)"
    fi
}

# Run CLI unit tests
run_cli_unit_tests() {
    print_header "Running CLI Unit Tests"
    
    cd "$CLI_DIR"
    
    # Run Go tests
    if go test -v -coverprofile=coverage.out ./...; then
        print_success "CLI unit tests passed"
        
        # Generate coverage report
        go tool cover -html=coverage.out -o coverage.html
        mv coverage.html "$REPORTS_DIR/cli_unit_coverage.html"
        mv coverage.out "$REPORTS_DIR/cli_unit_coverage.out"
        print_info "Coverage report saved to $REPORTS_DIR/cli_unit_coverage.html"
    else
        print_error "CLI unit tests failed"
        return 1
    fi
}

# Run CLI integration tests
run_cli_integration_tests() {
    print_header "Running CLI Integration Tests"
    
    cd "$TEST_DIR"
    
    # Run integration tests
    if go test -v -tags=integration ./...; then
        print_success "CLI integration tests passed"
    else
        print_error "CLI integration tests failed"
        return 1
    fi
}

# Run Flutter web unit tests
run_flutter_unit_tests() {
    print_header "Running Flutter Web Unit Tests"
    
    cd "$FLUTTER_WEB_DIR"
    
    # Run Flutter tests
    if flutter test --coverage; then
        print_success "Flutter unit tests passed"
        
        # Generate coverage report
        genhtml coverage/lcov.info -o coverage/html
        cp -r coverage/html "$REPORTS_DIR/flutter_unit_coverage"
        print_info "Coverage report saved to $REPORTS_DIR/flutter_unit_coverage"
    else
        print_error "Flutter unit tests failed"
        return 1
    fi
}

# Run Flutter web widget tests
run_flutter_widget_tests() {
    print_header "Running Flutter Web Widget Tests"
    
    cd "$FLUTTER_WEB_DIR"
    
    # Run widget tests
    if flutter test test/widget_test.dart; then
        print_success "Flutter widget tests passed"
    else
        print_error "Flutter widget tests failed"
        return 1
    fi
}

# Run E2E tests
run_e2e_tests() {
    print_header "Running End-to-End Tests"
    
    cd "$TEST_DIR"
    
    # Run E2E tests
    if go test -v -tags=e2e ./...; then
        print_success "E2E tests passed"
    else
        print_error "E2E tests failed"
        return 1
    fi
}

# Run performance tests
run_performance_tests() {
    print_header "Running Performance Tests"
    
    cd "$TEST_DIR"
    
    # Run performance benchmarks
    if go test -v -bench=. -benchmem ./...; then
        print_success "Performance tests completed"
    else
        print_error "Performance tests failed"
        return 1
    fi
}

# Run web interface tests with browser automation
run_web_interface_tests() {
    print_header "Running Web Interface Tests"
    
    cd "$FLUTTER_WEB_DIR"
    
    # Check if webdriver is available
    if command -v chromedriver &> /dev/null; then
        print_info "Running web interface tests with Chrome"
        
        # Start Chrome driver
        chromedriver --port=4444 &
        CHROME_PID=$!
        
        # Wait for Chrome driver to start
        sleep 2
        
        # Run web tests
        if flutter drive --driver=test_driver/integration_test.dart --target=integration_test/app_test.dart; then
            print_success "Web interface tests passed"
        else
            print_error "Web interface tests failed"
            kill $CHROME_PID 2>/dev/null || true
            return 1
        fi
        
        # Stop Chrome driver
        kill $CHROME_PID 2>/dev/null || true
    else
        print_warning "ChromeDriver not found, skipping web interface tests"
    fi
}

# Run all tests
run_all_tests() {
    print_header "Running All Tests"
    
    local exit_code=0
    
    # Run CLI tests
    if run_cli_unit_tests; then
        print_success "CLI unit tests completed"
    else
        exit_code=1
    fi
    
    if run_cli_integration_tests; then
        print_success "CLI integration tests completed"
    else
        exit_code=1
    fi
    
    # Run Flutter tests
    if run_flutter_unit_tests; then
        print_success "Flutter unit tests completed"
    else
        exit_code=1
    fi
    
    if run_flutter_widget_tests; then
        print_success "Flutter widget tests completed"
    else
        exit_code=1
    fi
    
    # Run E2E tests
    if run_e2e_tests; then
        print_success "E2E tests completed"
    else
        exit_code=1
    fi
    
    # Run performance tests
    if run_performance_tests; then
        print_success "Performance tests completed"
    else
        exit_code=1
    fi
    
    # Run web interface tests
    if run_web_interface_tests; then
        print_success "Web interface tests completed"
    else
        exit_code=1
    fi
    
    return $exit_code
}

# Generate test report
generate_test_report() {
    print_header "Generating Test Report"
    
    local report_file="$REPORTS_DIR/test_report_$(date +%Y%m%d_%H%M%S).html"
    
    cat > "$report_file" << EOF
<!DOCTYPE html>
<html>
<head>
    <title>GWallet-AAI Test Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background-color: #f0f0f0; padding: 20px; border-radius: 5px; }
        .section { margin: 20px 0; }
        .success { color: green; }
        .error { color: red; }
        .warning { color: orange; }
        .coverage { background-color: #f9f9f9; padding: 10px; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>GWallet-AAI Test Report</h1>
        <p>Generated on: $(date)</p>
    </div>
    
    <div class="section">
        <h2>Test Summary</h2>
        <p>This report contains the results of automated testing for the GWallet-AAI project.</p>
    </div>
    
    <div class="section">
        <h2>Coverage Reports</h2>
        <div class="coverage">
            <p><a href="cli_unit_coverage.html">CLI Unit Test Coverage</a></p>
            <p><a href="flutter_unit_coverage/index.html">Flutter Unit Test Coverage</a></p>
        </div>
    </div>
    
    <div class="section">
        <h2>Test Results</h2>
        <p>Check the console output for detailed test results.</p>
    </div>
</body>
</html>
EOF
    
    print_success "Test report generated: $report_file"
}

# Clean up test artifacts
cleanup() {
    print_header "Cleaning Up"
    
    # Remove temporary files
    find "$PROJECT_ROOT" -name "*.tmp" -delete 2>/dev/null || true
    find "$PROJECT_ROOT" -name "test_*.jpg" -delete 2>/dev/null || true
    
    print_success "Cleanup completed"
}

# Show usage
show_usage() {
    echo "Usage: $0 [OPTIONS] [TEST_TYPE]"
    echo ""
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  -c, --clean    Clean up test artifacts before running"
    echo "  -r, --report   Generate HTML test report"
    echo ""
    echo "Test Types:"
    echo "  unit          Run unit tests only"
    echo "  integration   Run integration tests only"
    echo "  e2e           Run end-to-end tests only"
    echo "  performance   Run performance tests only"
    echo "  all           Run all tests (default)"
    echo ""
    echo "Examples:"
    echo "  $0                    # Run all tests"
    echo "  $0 unit              # Run unit tests only"
    echo "  $0 -c e2e            # Clean and run E2E tests"
    echo "  $0 -r all            # Run all tests and generate report"
}

# Main script
main() {
    local test_type="all"
    local clean_flag=false
    local report_flag=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_usage
                exit 0
                ;;
            -c|--clean)
                clean_flag=true
                shift
                ;;
            -r|--report)
                report_flag=true
                shift
                ;;
            unit|integration|e2e|performance|all)
                test_type="$1"
                shift
                ;;
            *)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    # Print header
    print_header "GWallet-AAI Test Runner"
    print_info "Test Type: $test_type"
    print_info "Project Root: $PROJECT_ROOT"
    
    # Check dependencies
    check_dependencies
    
    # Create reports directory
    create_reports_dir
    
    # Clean up if requested
    if [ "$clean_flag" = true ]; then
        cleanup
    fi
    
    # Run tests based on type
    case $test_type in
        unit)
            run_cli_unit_tests
            run_flutter_unit_tests
            ;;
        integration)
            run_cli_integration_tests
            run_flutter_widget_tests
            ;;
        e2e)
            run_e2e_tests
            run_web_interface_tests
            ;;
        performance)
            run_performance_tests
            ;;
        all)
            run_all_tests
            ;;
    esac
    
    # Generate report if requested
    if [ "$report_flag" = true ]; then
        generate_test_report
    fi
    
    # Final cleanup
    cleanup
    
    print_header "Test Execution Completed"
    print_success "All tests completed successfully!"
}

# Run main function
main "$@" 