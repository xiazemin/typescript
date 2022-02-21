var electron=require('electron')

var app=electron.app
var BrowserWindow=electron.BrowserWindow

var mainWindow=null//要打开的主窗口

app.on('ready',()=>{
    mainWindow=new BrowserWindow({
        with:800,
        height:800
    })
    mainWindow.loadFile('index.html')
    mainWindow.on('closed',()=>{
        mainWindow=null
    })
})
