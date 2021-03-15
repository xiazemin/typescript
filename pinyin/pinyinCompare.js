const pinyin = require('pinyin');

const data = '我要排序'.split('');
console.log(data)
console.log(pinyin(data[0]))
// 建议将汉字的拼音持久化存储起来。
const pinyinData = data.map(han => ({
  han: han,
  pinyin: pinyin(han)[0][0], // 可以自行选择不同的生成拼音方案和风格。
}));
const sortedData = pinyinData.sort((a, b) => {
  return a.pinyin.localeCompare(b.pinyin);
}).map(d => d.han);


console.log(sortedData)