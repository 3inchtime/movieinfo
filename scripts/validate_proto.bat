@echo off
REM Proto文件语法验证脚本
REM 验证proto文件是否有语法错误

echo Validating proto files syntax...
echo.

REM 设置项目根目录
set PROJECT_ROOT=%~dp0..
cd /d "%PROJECT_ROOT%"

REM 验证通用模块
echo Validating common proto files...
protoc --proto_path=proto --descriptor_set_out=NUL common/common.proto
if %errorlevel% neq 0 (
    echo ERROR: common/common.proto has syntax errors
    goto :error
)
echo   common/common.proto - OK

protoc --proto_path=proto --descriptor_set_out=NUL common/error.proto
if %errorlevel% neq 0 (
    echo ERROR: common/error.proto has syntax errors
    goto :error
)
echo   common/error.proto - OK

REM 验证用户服务
echo Validating user service proto files...
protoc --proto_path=proto --descriptor_set_out=NUL user/user.proto
if %errorlevel% neq 0 (
    echo ERROR: user/user.proto has syntax errors
    goto :error
)
echo   user/user.proto - OK

protoc --proto_path=proto --descriptor_set_out=NUL user/user_service.proto
if %errorlevel% neq 0 (
    echo ERROR: user/user_service.proto has syntax errors
    goto :error
)
echo   user/user_service.proto - OK

REM 验证电影服务
echo Validating movie service proto files...
protoc --proto_path=proto --descriptor_set_out=NUL movie/movie.proto
if %errorlevel% neq 0 (
    echo ERROR: movie/movie.proto has syntax errors
    goto :error
)
echo   movie/movie.proto - OK

protoc --proto_path=proto --descriptor_set_out=NUL movie/movie_service.proto
if %errorlevel% neq 0 (
    echo ERROR: movie/movie_service.proto has syntax errors
    goto :error
)
echo   movie/movie_service.proto - OK

REM 验证评分服务
echo Validating rating service proto files...
protoc --proto_path=proto --descriptor_set_out=NUL rating/rating.proto
if %errorlevel% neq 0 (
    echo ERROR: rating/rating.proto has syntax errors
    goto :error
)
echo   rating/rating.proto - OK

protoc --proto_path=proto --descriptor_set_out=NUL rating/rating_service.proto
if %errorlevel% neq 0 (
    echo ERROR: rating/rating_service.proto has syntax errors
    goto :error
)
echo   rating/rating_service.proto - OK

echo.
echo All proto files are valid!
echo Proto validation completed successfully.
goto :end

:error
echo.
echo Proto validation failed!
echo Please check the error messages above and fix the syntax errors.
pause
exit /b 1

:end
pause