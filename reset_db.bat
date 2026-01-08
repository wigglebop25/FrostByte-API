@echo off
echo Stopping containers...
docker-compose down -v

echo.
echo Starting fresh containers (Database will be reset)...
docker-compose up -d --build

echo.
echo Waiting for API and Database to initialize (20 seconds)...
timeout /t 20

echo.
echo Database reset complete. You can now register a new admin user.
pause