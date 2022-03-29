list=[{a:1,b:2},{a:3,b:4}]

for (let i = 0; i < list.length; i++) {
    const item = list[i]
    list[i]={
     a:item.a
    }
}
console.log(list)
