require('@electron/remote/main').initialize()
const { BrowserView } = require('electron')

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
            contextIsolation:false,//https://blog.csdn.net/qq_35066582/article/details/114457490
            enableRemoteModule:true,//允许远程模块 https://blog.csdn.net/liu19721018/article/details/108922594
            backgroundThrottling: false,
            webSecurity: false 
        }
    })

    require('@electron/remote/main').enable(mainWindow.webContents);
    
    //打开开发工具页面
    mainWindow.webContents.openDevTools();
    mainWindow.loadFile('index.html')

    mainWindow.on('closed',()=>{
        mainWindow=null
    })
})