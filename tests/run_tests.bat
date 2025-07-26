@echo off
REM Test Runner Script for GWallet-AAI (Windows Batch Version)
REM This script runs different types of tests for CLI and web interface components

setlocal enabledelayedexpansion

REM Configuration
set "PROJECT_ROOT=%~dp0.."
set "TEST_DIR=%PROJECT_ROOT%\tests"
set "CLI_DIR=%PROJECT_ROOT%\cli"
set "FLUTTER_WEB_DIR=%PROJECT_ROOT%\flutter_web"
set "REPORTS_DIR=%PROJECT_ROOT%\test_reports"

REM Check for help argument first
if "%~1"=="-h" goto :show_usage
if "%~1"=="--help" goto :show_usage
if "%~1"=="/?" goto :show_usage

REM Initialize variables
set "test_type=all"
set "clean_flag=false"
set "report_flag=false"

REM Parse arguments
:parse_args
if "%~1"=="" goto :main
if "%~1"=="-c" (
    set "clean_flag=true"
    shift
    goto :parse_args
)
if "%~1"=="--clean" (
    set "clean_flag=true"
    shift
    goto :parse_args
)
if "%~1"=="-r" (
    set "report_flag=true"
    shift
    goto :parse_args
)
if "%~1"=="--report" (
    set "report_flag=true"
    shift
    goto :parse_args
)
if "%~1"=="unit" (
    set "test_type=unit"
    shift
    goto :parse_args
)
if "%~1"=="integration" (
    set "test_type=integration"
    shift
    goto :parse_args
)
if "%~1"=="e2e" (
    set "test_type=e2e"
    shift
    goto :parse_args
)
if "%~1"=="performance" (
    set "test_type=performance"
    shift
    goto :parse_args
)
if "%~1"=="all" (
    set "test_type=all"
    shift
    goto :parse_args
)
echo Unknown option: %~1
goto :show_usage

:show_usage
echo.
echo Usage: run_tests.bat [OPTIONS] [TEST_TYPE]
echo.
echo Options:
echo   -h, --help     Show this help message
echo   -c, --clean    Clean up test artifacts before running
echo   -r, --report   Generate HTML test report
echo.
echo Test Types:
echo   unit          Run unit tests only
echo   integration   Run integration tests only
echo   e2e           Run end-to-end tests only
echo   performance   Run performance tests only
echo   all           Run all tests (default)
echo.
echo Examples:
echo   run_tests.bat                    # Run all tests
echo   run_tests.bat unit              # Run unit tests only
echo   run_tests.bat -c e2e            # Clean and run E2E tests
echo   run_tests.bat -r all            # Run all tests and generate report
echo.
exit /b 0

:print_header
echo.
echo ================================
echo %~1
echo ================================
echo.
goto :eof

:print_success
echo [SUCCESS] %~1
goto :eof

:print_error
echo [ERROR] %~1
goto :eof

:print_warning
echo [WARNING] %~1
goto :eof

:check_dependencies
call :print_header "Checking Dependencies"

REM Check if Go is installed
go version >nul 2>&1
if errorlevel 1 (
    call :print_error "Go is not installed or not in PATH"
    exit /b 1
)
call :print_success "Go is installed"

REM Check if Flutter is installed
flutter --version >nul 2>&1
if errorlevel 1 (
    call :print_warning "Flutter is not installed or not in PATH"
    call :print_warning "Web interface tests will be skipped"
    set "SKIP_FLUTTER=true"
) else (
    call :print_success "Flutter is installed"
    set "SKIP_FLUTTER=false"
)
goto :eof

:run_cli_unit_tests
call :print_header "Running CLI Unit Tests"

cd /d "%CLI_DIR%"

go test -v -coverprofile=coverage.out ./...
if errorlevel 1 (
    call :print_error "CLI unit tests failed"
    exit /b 1
)

call :print_success "CLI unit tests passed"

REM Generate coverage report
if exist coverage.out (
    go tool cover -html=coverage.out -o coverage.html
    call :print_success "Coverage report generated: coverage.html"
)
goto :eof

:run_integration_tests
call :print_header "Running Integration Tests"

cd /d "%TEST_DIR%"

go test -v -run TestIntegration ./...
if errorlevel 1 (
    call :print_error "Integration tests failed"
    exit /b 1
)

call :print_success "Integration tests passed"
goto :eof

:run_e2e_tests
call :print_header "Running End-to-End Tests"

cd /d "%TEST_DIR%"

go test -v -run TestE2E ./...
if errorlevel 1 (
    call :print_error "E2E tests failed"
    exit /b 1
)

call :print_success "E2E tests passed"
goto :eof

:run_performance_tests
call :print_header "Running Performance Tests"

cd /d "%TEST_DIR%"

go test -v -bench=. -benchmem ./...
if errorlevel 1 (
    call :print_error "Performance tests failed"
    exit /b 1
)

call :print_success "Performance tests passed"
goto :eof

:run_flutter_tests
if "%SKIP_FLUTTER%"=="true" (
    call :print_warning "Skipping Flutter tests (Flutter not installed)"
    goto :eof
)

call :print_header "Running Flutter Web Tests"

cd /d "%FLUTTER_WEB_DIR%"

flutter test
if errorlevel 1 (
    call :print_error "Flutter tests failed"
    exit /b 1
)

call :print_success "Flutter tests passed"
goto :eof

:clean_artifacts
call :print_header "Cleaning Test Artifacts"

if exist "%CLI_DIR%\coverage.out" del "%CLI_DIR%\coverage.out"
if exist "%CLI_DIR%\coverage.html" del "%CLI_DIR%\coverage.html"
if exist "%REPORTS_DIR%" rmdir /s /q "%REPORTS_DIR%"

call :print_success "Test artifacts cleaned"
goto :eof

:generate_report
call :print_header "Generating Test Report"

if not exist "%REPORTS_DIR%" mkdir "%REPORTS_DIR%"

echo ^<!DOCTYPE html^> > "%REPORTS_DIR%\test_report.html"
echo ^<html^> >> "%REPORTS_DIR%\test_report.html"
echo ^<head^>^<title^>GWallet-AAI Test Report^</title^>^</head^> >> "%REPORTS_DIR%\test_report.html"
echo ^<body^> >> "%REPORTS_DIR%\test_report.html"
echo ^<h1^>GWallet-AAI Test Report^</h1^> >> "%REPORTS_DIR%\test_report.html"
echo ^<p^>Generated on: %date% %time%^</p^> >> "%REPORTS_DIR%\test_report.html"
echo ^<p^>All tests completed successfully!^</p^> >> "%REPORTS_DIR%\test_report.html"
echo ^</body^>^</html^> >> "%REPORTS_DIR%\test_report.html"

call :print_success "Test report generated: %REPORTS_DIR%\test_report.html"
goto :eof

:main
call :print_header "GWallet-AAI Test Runner"

REM Check dependencies
call :check_dependencies
if errorlevel 1 exit /b 1

REM Clean artifacts if requested
if "%clean_flag%"=="true" call :clean_artifacts

REM Run tests based on type
if "%test_type%"=="unit" (
    call :run_cli_unit_tests
    if errorlevel 1 exit /b 1
    call :run_flutter_tests
    if errorlevel 1 exit /b 1
) else if "%test_type%"=="integration" (
    call :run_integration_tests
    if errorlevel 1 exit /b 1
) else if "%test_type%"=="e2e" (
    call :run_e2e_tests
    if errorlevel 1 exit /b 1
) else if "%test_type%"=="performance" (
    call :run_performance_tests
    if errorlevel 1 exit /b 1
) else (
    REM Run all tests
    call :run_cli_unit_tests
    if errorlevel 1 exit /b 1
    call :run_integration_tests
    if errorlevel 1 exit /b 1
    call :run_e2e_tests
    if errorlevel 1 exit /b 1
    call :run_performance_tests
    if errorlevel 1 exit /b 1
    call :run_flutter_tests
    if errorlevel 1 exit /b 1
)

REM Generate report if requested
if "%report_flag%"=="true" call :generate_report

call :print_header "All Tests Completed Successfully!"
echo.
call :print_success "Test execution completed"
if "%report_flag%"=="true" (
    echo Test report available at: %REPORTS_DIR%\test_report.html
)
echo. 