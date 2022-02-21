require('@electron/remote/main').initialize()
const { BrowserView } = require('electron')
//https://www.cnblogs.com/jspdelphi/p/15606260.html
// electron版本>=14.0.0,每个单独的webContents想要使用remote module，必须使用新的enable API来一个个使能.默认remote module是不可用的
// electron版本<14.0.0  版本可以使用enableRemoteModule来控制
// webPreferences{enableRemoteModule：false}可以禁用remote module

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

    require('./main/menu')

    require('@electron/remote/main').enable(mainWindow.webContents);
    
    //打开开发工具页面
    mainWindow.webContents.openDevTools();
    mainWindow.loadFile('index.html')
    
    /*
    var BrowserView=electron.BrowserView
    var view =new BrowserView()
    mainWindow.setBrowserView(view)
    view.setBounds({x:0,y:120,width:1000,height:680})
    view.webContents.loadURL('https://ww.baidu.com')*/

   // window.open()

    mainWindow.on('closed',()=>{
        mainWindow=null
    })
})
