@ECHO OFF
setlocal
:choice
echo This script will install the SSHrimp-Agent to run automatically for all users
echo This script MUST run with administrator privileges
set /P c=Do you want to proceed?  [Y/N]?
if /I "%c%" EQU "Y" goto :runcode
if /I "%c%" EQU "N" goto :EOF
goto :choice


:runcode

if not exist C:\Users\Public\SSHrimp MKDIR C:\Users\Public\SSHrimp
CALL XCOPY /y .\sshrimp-agent-windows.exe C:\Users\Public\SSHrimp\
CALL XCOPY /y .\sshrimp-windows.toml C:\Users\Public\SSHrimp\
CALL XCOPY /y .\sshrimp-run-in-background.vbs "C:\ProgramData\Microsoft\Windows\Start Menu\Programs\StartUp"
CALL "C:\ProgramData\Microsoft\Windows\Start Menu\Programs\StartUp\sshrimp-run-in-background.vbs"

echo "Deployment complete."
pause


:EOF
endlocal
@ECHO ON