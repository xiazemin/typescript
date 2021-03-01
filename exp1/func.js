function buildName(firstName, lastName) {
    if (lastName)
        return firstName + " " + lastName;
    else
        return firstName;
}
var result1 = buildName("Bob"); // 正确
//let result2 = buildName("Bob", "Adams", "Sr.");  // 错误，参数太多了
var result3 = buildName("Bob", "Adams"); // 正确
