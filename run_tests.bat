@echo off
echo Starting Docker services...
docker-compose up -d --build

echo Waiting for API to start (15 seconds)...
timeout /t 15

echo Running Bruno tests...
call npx @usebruno/cli run bruno --output json --insecure

echo.
echo Tests completed.
pause
