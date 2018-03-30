
echo off
cls
FOR /L %%A IN (1,1,2000) DO (
rem go run mainprog.go dataobj.go 
rem go run mygame.go
rem -race 
go run client.go DESKTOP-K53L4N7 alice
pause
)
