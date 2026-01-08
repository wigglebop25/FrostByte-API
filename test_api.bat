@echo off
setlocal

echo Starting Docker services...
docker-compose up -d --build

echo Waiting for API to start (15 seconds)...
timeout /t 15

echo Renaming 'orders' folder to 'z_orders' to ensure correct test execution order...
if exist "bruno\orders" rename "bruno\orders" "z_orders"

echo Running Bruno tests...
call npx @usebruno/cli run bruno --output json --insecure

echo Restoring 'orders' folder name...
if exist "bruno\z_orders" rename "bruno\z_orders" "orders"

echo.
echo Tests completed.
pause
