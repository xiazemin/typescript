const btn=this.document.querySelector("#btn")
//require('electron').hideInternalModules()
//https://www.electronjs.org/zh/blog/electron-api-changes#an-easier-way-to-use-the-remote-module

//https://wizardforcel.gitbooks.io/electron-doc/content/api/remote.html
// const { BrowserWindow } =require('remote');
const { BrowserWindow } = require('@electron/remote')
//const BrowserWindow=require('electron').remote.BrowserWindow;//通过 electron中的remote模块间接的调用主进程中的模块 https://cloud.tencent.com/developer/ask/sof/821881
//https://github.com/electron/electron/blob/v9.0.0-beta.17/docs/api/remote.md
//https://github.com/electron/electron/issues/32151
//https://www.cnblogs.com/cc11001100/p/14290450.html
https://github.com/electron-userland/electron-remote
//https://github.com/electron-userland/electron-remote
window.onload=function(){
    btn.onclick=()=>{
       newWin=new BrowserWindow(
        {
            width:500,
            height:500,
            webPreferences:{
                nodeIntegration:true,
                enableRemoteModule:true
              }
        })
        newWin.loadFile('yellow.html')
        newWin.on('close',()=>{
            newWin=null
        })
    }
}