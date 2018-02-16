@echo off
set DNOTEPATH=%PROGRAMFILES%\Dnote CLI
set DNOTEDL=dnote-windows-amd64.exe
set DNOTETARGET=%DNOTEPATH%\dnote.exe
echo Checking for directory...
if not exist "%DNOTEPATH%\" (
  echo Creating directory...
  mkdir "%DNOTEPATH%"
)
echo Moving program to target directory...
move /Y %DNOTEDL% "%DNOTETARGET%"
echo "Adding directory to user PATH..."

REM retrieve only the user's PATH from registry,
REM to avoid attempting (and probably failing) to overwrite the
REM system path

set Key=HKCU\Environment
FOR /F "tokens=2* skip=1" %%G IN ('REG QUERY %Key% /v PATH') DO (
  echo %%H > user_path_backup.txt
  set t=%%H
  set "NEWPATH="
  :loop
  for /f "delims=; tokens=1*" %%a in ("%t%") do (
    set t=%%b
    if not "%%a" == "%DNOTEPATH%" (
      if defined NEWPATH (
        set NEWPATH=%NEWPATH%;%%a
      ) else (
        set NEWPATH=%%a
      )
    )
  )
  if defined t goto :loop
)
set NEWPATH=%NEWPATH%;%DNOTEPATH%
setx PATH "%NEWPATH%"
