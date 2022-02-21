var electron=require('electron')

var app=electron.app
var BrowserWindow=electron.BrowserWindow

var mainWindow=null//要打开的主窗口

app.on('ready',()=>{
    mainWindow=new BrowserWindow({
        with:800,
        height:800,
        webPreferences:{
            nodeIntegration:true,
            contextIsolation:false//https://blog.csdn.net/qq_35066582/article/details/114457490
        }
    })
    mainWindow.loadFile('index.html')
    mainWindow.on('closed',()=>{
        mainWindow=null
    })
})
