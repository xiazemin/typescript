MongoDB以BSON格式的文档（Documents）形式存储。Databases中包含集合（Collections），集合（Collections）中存储文档（Documents）。

BSON是一个二进制形式的JSON文档，它比JSON包含更多的数据类型。对于BSON规格，可参见bsonspec.org，也可参考BSON类型。


ocument结构
MongoDB的文件是由field和value对的结构组成，例如下面这样的结构：

{
   field1: value1,
   field2: value2,
}

value值可以是任何BSON数据类型，包括：其他document，数字，和document数组。


Documents中的filed名有下列限制：

_id被保留用于主键；其值必须是集合中唯一的、不可变的、并且可以是数组以外的任何数据类型
不能以美元符号$开头
不能包含点字符.
不能包含空字符

Field Value限制
对于索引的collections，索引字段中的值有最大长度限制


圆点符号
MongoDB中使用圆点符号.访问数组中的元素，也可以访问嵌入式Documents的fields。

Arrays数组
通过圆点符号.来链接Arrays数组名字和从0开始的数字位置，来定位和访问一个元素数组：

"<array>.<index>"

要访问contribs数组中的第三个元素，可以这样访问：

"contribs.2"

嵌入式Documents
通过圆点符号.来链接嵌入式document的名字和field名，来定位和访问嵌入式document：

"<embedded document>.<field>"

要访问name中的last字段，可以这样使用：

"name.last"


Document Field顺序
MongoDB中field的顺序默认是按照写操作的顺序来保存的，除了下面几种情况：

_id总是document的第一个field
可能会导致文档中的字段的重新排序的更新，包括字段名重命名。


_id字段有以下行为和限制：

默认情况下，MongoDB会在创建collection时创建一个_id字段的唯一索引
_id字段总是documents中的第一个字段。如果服务器接收到一个docuement，它的第一个字段不是_id，那么服务器会将_id字段移在开头
_id字段可以是除了array数组之外的任何BSON数据格式
以下是存储_id值的常用选项：

使用ObjectId
最好使用自然的唯一标识符，可以节省空间并避免额外的索引
生成一个自动递增的数字。请参阅创建一个自动递增序列字段
在您的应用程序代码中生成UUID。为了更高效的在collection和_id索引中存储UUID值，可以用BSON的BinData类型存储UUID。

https://www.cnblogs.com/kkdn/p/9435257.html



