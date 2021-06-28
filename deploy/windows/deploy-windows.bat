@ECHO OFF
setlocal
:choice
echo This script will install the SSHrimp-Agent to run automatically for the current user
echo This is done by placing a script in the "%APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup" folder
echo To uninstall, delete the script from this folder, as well as the sshrimp folder from %USERPROFILE%\SSHrimp\
echo This script does NOT need to be run with administrator privileges
set /P c=Do you want to proceed?  [Y/N]?
if /I "%c%" EQU "Y" goto :runcode
if /I "%c%" EQU "N" goto :EOF
goto :choice


:runcode

if not exist %USERPROFILE%\SSHrimp MKDIR %USERPROFILE%\SSHrimp
CALL XCOPY /y .\sshrimp-agent-windows.exe %USERPROFILE%\SSHrimp\
CALL XCOPY /y .\sshrimp-windows.toml %USERPROFILE%\SSHrimp\

@REM CALL XCOPY /y .\sshrimp-run-in-background.vbs "C:\ProgramData\Microsoft\Windows\Start Menu\Programs\StartUp"
CALL XCOPY /y .\sshrimp-run-in-background.vbs "%APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup"
start "" "%APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup\sshrimp-run-in-background.vbs"

echo "Deployment complete."
pause


:EOF
endlocal
@ECHO ON