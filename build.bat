@echo off
REM Script de build do projeto Lince
echo Building Lince...
set OUTPUT=build\lince.exe
if not exist build mkdir build
go build -o %OUTPUT% .
echo Build concluido: %OUTPUT%
