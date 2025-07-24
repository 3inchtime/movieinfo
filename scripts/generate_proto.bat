@echo off
REM gRPC Code Generation Script - Windows Version
REM Generate Go code based on simplified gRPC protocol definitions

echo Starting gRPC code generation...

REM Set project root directory
set PROJECT_ROOT=%~dp0..
cd /d "%PROJECT_ROOT%"

REM Check if protoc is installed
protoc --version >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: protoc not installed or not in PATH
    echo Please install Protocol Buffers compiler
    pause
    exit /b 1
)

REM Check if protoc-gen-go is installed
protoc-gen-go --version >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: protoc-gen-go not installed
    echo Please run: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    pause
    exit /b 1
)

REM Check if protoc-gen-go-grpc is installed
protoc-gen-go-grpc --version >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: protoc-gen-go-grpc not installed
    echo Please run: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    pause
    exit /b 1
)

REM Create output directories
if not exist "proto\gen" mkdir "proto\gen"
if not exist "proto\gen\common" mkdir "proto\gen\common"
if not exist "proto\gen\user" mkdir "proto\gen\user"
if not exist "proto\gen\movie" mkdir "proto\gen\movie"
if not exist "proto\gen\rating" mkdir "proto\gen\rating"

echo Generating common module code...
protoc --proto_path=proto --go_out=proto/gen --go_opt=paths=source_relative --go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative common/common.proto common/error.proto

if %errorlevel% neq 0 (
    echo Error: Failed to generate common module code
    pause
    exit /b 1
)

echo Generating user service code...
protoc --proto_path=proto --go_out=proto/gen --go_opt=paths=source_relative --go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative user/user.proto user/user_service.proto

if %errorlevel% neq 0 (
    echo Error: Failed to generate user service code
    pause
    exit /b 1
)

echo Generating movie service code...
protoc --proto_path=proto --go_out=proto/gen --go_opt=paths=source_relative --go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative movie/movie.proto movie/movie_service.proto

if %errorlevel% neq 0 (
    echo Error: Failed to generate movie service code
    pause
    exit /b 1
)

echo Generating rating service code...
protoc --proto_path=proto --go_out=proto/gen --go_opt=paths=source_relative --go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative rating/rating.proto rating/rating_service.proto

if %errorlevel% neq 0 (
    echo Error: Failed to generate rating service code
    pause
    exit /b 1
)

echo.
echo gRPC code generation completed!
echo Generated files are located at: proto/gen/
echo.
echo File structure:
echo   proto/gen/common/    - Common data types and error definitions
echo   proto/gen/user/      - User service related code
echo   proto/gen/movie/     - Movie service related code
echo   proto/gen/rating/    - Rating service related code
echo.
pause