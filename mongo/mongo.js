//cnpm install mongodb -g
var MongoClient = require('mongodb').MongoClient;
var url = "mongodb://localhost:27017/quill-sharedb-cursors"
//"mongodb://localhost:27017/runoob";

MongoClient.connect(url, function (err, db) {
    if (err) throw err;
    console.log('数据库已创建');
    const myAsyncFunction = async () => {
        var dbase = db.db("quill-sharedb-cursors");
        if (!dbase.collection('site')) {
            await dbase.createCollection('site', function (err, res) {
                if (err) throw err;
                console.log("创建集合!");
            });
        }

        console.log('集合已创建');
        var dbo = db.db("quill-sharedb-cursors");
        var myobj = { name: "菜鸟教程", url: "www.runoob" };
        await dbo.collection("site").insertOne(myobj, function (err, res) {
            if (err) throw err;
            console.log("文档插入成功");
        });


        console.log('文档已创建');
        await dbo.collection("site").find({}).toArray(function (err, result) { // 返回集合中所有数据
            if (err) throw err;
            console.log(result);

        });
        console.log('文档已读取');
    }
    myAsyncFunction().then(() => {
        db.close();
    }
    );

});

//https://www.runoob.com/nodejs/nodejs-mongodb.html