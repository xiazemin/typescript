当使用第三方库时，我们需要引用它的声明文件，才能获得对应的代码补全、接口提示等功能。

第三方申明文件
推荐的是使用 @types 统一管理第三方库的声明文件。
@types 的使用方式很简单，直接用 npm 安装对应的声明模块即可，以 jQuery 举例：

    npm install @types/jquery --save-dev

https://blog.csdn.net/qq_41914181/article/details/109062989

npm包的声明文件可能存于两个地方：

与该npm包绑定在一起。判断依据是package.json中有types字段，或者有一个index.d.ts声明文件。这种模式不需要额外安装其他包，是最为退件的，所以以后我们自己创建npm包的时候，最好也将声明文件与npm包绑定在一起。
发布到@types里。我们只需要尝试安装一下对应的@types包就知道是否存在改声明文件npm install @types/foo --save-dev。这种模式一般是由于 npm 包的维护者没有提供声明文件，所以只能由其他人将声明文件发布到 @types 里了。
假如上面两种方式都没有找到对应的声明文件，那么我们就需要自己为他写声明文件了。由于是通过import语句导入的模块，所以声明文件存放的文职也有所约束，一般有两种方案：

创建一个 node_modules/@types/foo/index.d.ts 文件，存放 foo 模块的声明文件。这种方式不需要额外的配置，但是 node_modules 目录不稳定，代码也没有被保存到仓库中，无法回溯版本，有不小心被删除的风险，故不太建议用这种方案，一般只用作临时测试。
创建一个 types 目录，专门用来管理自己写的声明文件，将 foo 的声明文件放到 types/foo/index.d.ts 中。这种方式需要配置下 tsconfig.json 中的 paths 和 baseUrl 字段。

https://tasaid.com/blog/20171102225101.html