@echo off
chcp 65001
echo.
echo Regenerating file
echo.
go run -v .\cmd\mysql_cmd\main.go  -addr %1 -user %2 -pass %3 -name %4 -tables %5
if %errorlevel% == 1 (
echo.
echo failed!!!
exit 1
)

echo.
echo Done.
