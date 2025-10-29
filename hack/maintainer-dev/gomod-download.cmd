@echo off
rem initiate the retry number
set retryNumber=0
set maxRetries=3

:RESTORE
go mod download

rem problem?
IF NOT ERRORLEVEL 1 GOTO :EOF
@echo Oops, go mod download exited with code %ERRORLEVEL% - let us try again!
set /a retryNumber=%retryNumber%+1
IF %reTryNumber% LSS %maxRetries% (GOTO :RESTORE)
@echo Sorry, we tried downloading go modules %maxRetries% times and all attempts were unsuccessful!
EXIT /B 1
