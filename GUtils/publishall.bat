for /d %%d in (*.*) do if EXIST %%d\publish.bat (
	pushd %%d
	call publish.bat
	popd
)
pause