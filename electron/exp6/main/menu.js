const {Menu, BrowserWindow}=require('electron')

var template=[
    {
        label:"菜单1",
        accelerator:'ctrl+n',
        submenu:[
            {
                label:"子菜单1",
                accelerator:'ctrl+m',
                click:()=>{
                    var win =new BrowserWindow({
                        width:500,
                        height:500,
                        webPreferences:{nodeIntegration:true}
                    })
                    win.loadFile('yellow.html')
                    win.on('closed',()=>{
                        win=null
                    })
                }
            },
            {label:"子菜单2"}
        ],
    },
    {
        label:"菜单2",
        submenu:[
            {label:"2子菜单1"},
            {label:"2子菜单2"}
        ],
    }
]

var m=Menu.buildFromTemplate(template)

Menu.setApplicationMenu(m)