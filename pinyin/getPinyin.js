const pinyin=require("pinyin")
const _ = require('lodash')

function getPinyin (hanzi) {
    hanzi = hanzi.replace(/[^\u4E00-\u9FFFa-zA-Z0-9]+$/g, ' ')
    //把结尾地方非汉字，字母、数字的标点干掉
    console.log("hanzi replace:"+hanzi)
    const hanziArray = hanzi.split(/\s+/g)
   //按照一个或者多个空白字符拆分
    console.log(hanziArray)
    const initials = [] //首字母集合
    const words = []   //拼音的第一个选项集合
    const shorten = function (str) {
      let startIdx = 0
      const retArr = []
      const maxLen = 10
      const splits = Math.ceil(str.length / 10)
      while (startIdx < splits) {
        retArr.push(str.substring(startIdx * maxLen, (1 + startIdx) * maxLen))
        startIdx++
      }
      return retArr
    }
    //10个字符拆分
    for(var i=0;i<hanziArray.length;i++){
        var shor=shorten(hanziArray[i])
       console.log(shor)
       console.log(_.flatten(shor))
       _.flatten(shor).forEach((it)=>{
        const wg=pinyin(it,{
            style: pinyin.STYLE_NORMAL
          })   
          console.log("------------")
        console.log(wg)
        console.log("------------")
        wg.forEach((i)=>{
            console.log(i[0][0])
        })
       })

    }
    _.flatten(hanziArray.map(shorten)).forEach(function (item) {
      const wordsGroups = pinyin(item, {
        style: pinyin.STYLE_NORMAL
      })
      return wordsGroups.forEach(function (word) {
        console.log("--xxxxxxxxxxxxxxx--")
        console.log(JSON.stringify(word))
        console.log("--xxxxxxxxxxxxxxx--")
        initials.push(word[0][0])
        words.push('\'' + word[0])
      })
    })
    //首字母集合+拼音的第一个选项集合，用｜ 分割
    return (initials.join('') + '|' + words.join('')).toLowerCase()
  }

  var a="测试一个很长很长。。。很长很长。。。很长很长。。。很长很长。。。很长很长。。。很长很长。。。很长很长。。。很长很长。。。的中文!!!!@ 我和我的#$%^&*祖国，一刻也不能分割,我和我的#$%^&*祖国，一刻也不能分割！！  @ xx #c !&^ I love china 哈哈#"

  console.log(getPinyin(a))

  //console.log(a.replace(/[^\u4E00-\u9FFFa-zA-Z0-9]+$/g," "))

  //console.log(a.replace(/[\u4E00-\u9FFFa-zA-Z0-9]+$/g," "))