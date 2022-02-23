cnpm install electron -g
electron

electron .

remote 错误汇总
https://cloud.tencent.com/developer/article/1884363

// Deprecated in Electron 12:
const { BrowserWindow } = require('electron').remote
// Replace with:
const { BrowserWindow } = require('@electron/remote')

https://www.electronjs.org/docs/latest/breaking-changes

https://www.psvmc.cn/article/2021-09-29-electron-api-change.html

cnpm i electron/remote -g

npm i electron/remote -g

 npm install --save @electron/remote -g
 https://github.com/electron/remote
 cnpm install --save @electron/remote -g


  Cannot find module '../dist/src/renderer'
  https://github.com/parcel-bundler/parcel/issues/1843
  

cnpm install electron-builder -g
https://segmentfault.com/a/1190000013924153


npm run dist
https://segmentfault.com/a/1190000013924153

https://www.cnblogs.com/mmykdbc/p/11468908.html

sudo npm -g install electron --unsafe-perm=true --allow-root

cnpm install electron-packager -g

npm run-script package

electron-packager 打包太慢解决方法
https://juejin.cn/post/6844903951792341005

export ELECTRON_MIRROR="https://npm.taobao.org/mirrors/electron/"
export ELECTRON_CUSTOM_DIR=17.0.1


cnpm run-script package

> elec@1.0.0 package /Users/xiazemin/source/typescript/electron/exp13
> electron-packager . '应用名称' --platform=darwin --arch=x64 --icon=Pikachu.ico --out=./dist --asar --app-version=1.0.0 --ignore="(dist|src|docs|.gitignore|LICENSE|README.md|webpack.config*|node_modules)"



https://www.sqlpy.com/blogs/894258980
cnpm install electron-installer-dmg -g

% npx electron-installer-dmg ./dist/应用名称-darwin-x64/应用名称.app 应用名称

git rm -r -f --cached **/node_modules/

git rm -r --cached .

 git diff origin/master..master  --name-only
 
 https://www.jianshu.com/p/c70766b05408

export NODE_PATH='/usr/local' 
% npm prefix -g             
/usr/local

export NODE_PATH='/usr/local/lib/node_modules/'

export NODE_HOME='/usr/local/bin'
export PATH=$NODE_HOME:$PATH

npm config set cache /usr/local/lib/node_modules
npm config set prefix /usr/local/lib/node_modules


vi ~/.zshrc 
export NODE_PATH="/usr/local/lib/node_modules/"
source ~/.zshrc
https://www.cnblogs.com/miaodi/p/6607812.html


npm install --save
https://stackoverflow.com/questions/35682131/electron-packager-cannot-find-module

cnpm install --save

https://segmentfault.com/a/1190000014209821


cnpm i electron --save
cnpm i @electron/remote --save

原因 不能 --ignore dist/src/node_modules