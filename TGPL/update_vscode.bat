for /d %%f in (*.*) do (
	pushd %%f
	call go mod init %%f
	if exist .vscode\ copy /Y C:\Development\GitHub\Go\TGPL\ch01_01_dup1\.vscode .vscode
	popd
)