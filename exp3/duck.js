function addPoints(p1, p2) {
    var x = p1.x + p2.x;
    var y = p1.y + p2.y;
    return { x: x, y: y };
}
// 正确
var newPoint = addPoints({ x: 3, y: 4 }, { x: 5, y: 1 });
// 错误 
//var newPoint2 = addPoints({x:1},{x:4,y:3})
