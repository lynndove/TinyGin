# TinyGin

main 是使用

tinygin是框架

2.28

Tire树实现的动态路由具备以下两个功能。

- 参数匹配:。例如 /p/:lang/doc，可以匹配 /p/c/doc 和 /p/go/doc。

- 通配*。例如 /static/*filepath，可以匹配/static/fav.ico，也可以匹配/static/js/jQuery.js，这种模式常用于静态服务器，能够递归地匹配子路径。


Trie 树支持节点的插入与查询

插入

递归查找每一层的节点，如果没有匹配到当前part的节点，则新建一个.

例如：/p/:lang/doc只有在第三层节点，即doc节点，pattern才会设置为/p/:lang/doc。p和:lang节点的pattern
属性皆为空。因此，当匹配结束时，我们可以使用n.pattern == ""来判断路由规则是否匹配成功。例如，/p/python虽能成功匹配到:lang，但:lang的pattern值为空，因此匹配失败。

查询

同样也是递归查询每一层的节点，退出规则是，匹配到了*，匹配失败，或者匹配到了第len(parts)层节点。

3.1


路由分组控制(Route Group Control)

Group对象，还需要有访问Router的能力，为了方便，在Group中，保存一个指针，指向Engine，整个框架的所有资源都是由Engine统一协调的，那么就可以通过Engine间接地访问各种接口了
