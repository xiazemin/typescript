'use strict'
async function func1(){
    console.log("func1")
}

func1()

func1().then(()=>{
    console.log("func1 then")
})


async function func2(){
    return new Promise((resolve)=>{
        console.log("func2")
    })
}

func2()
func2().then(()=>{
    console.log("func2 then")
})

async function func3(){
    return new Promise((resolve)=>{
        console.log("func3")
        resolve()
    }).then(()=>{
        return 1
    })
}

func3()

func3().then(()=>{
    console.log("func3 then")
})


async function func4(){
    console.log("func4")
    let num=await func3()
    console.log(num)
}


func4()

func4().then(()=>{
    console.log("func4 then")
})

