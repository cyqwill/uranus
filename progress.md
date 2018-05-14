##  



## 2018-05-14

How to make every client by a **client**?
Simple, define a client, and every client is an user or an consumer.

**following things need to do**

1. 如何实现用户唯一凭证？token还是什么，建立几个基础信息结构体，很重要，比如pub：发布一个消息到指定的ID，如果是usr的话就是一个人，如果是grp的话就是一个群组，sub: 订阅一个消息，一个用户自动订阅它的所有好友的消息，以及他所在的群的消息，这个如何表示？值得思考。
2. 用户修改资料通过http来完成不要通过websocket。还有消息存储到数据库怎么操作？保存七天？一个月？永久保存？有个消息表，用户上线后还得把最新的消息推送给它，再它不在线的时候别人发过来的消息需要推送给他，怎么实现？每条消息都有一个状态：发送，已读，未送达。



Quick thought:

- 用户每次发消息都带一个user_token，首次登陆直接post login，发送用户名和密码获取token，token的目的是知道当前的用户是谁，token永久有效，通过get friends得到好友列表，chat history怎么写？也许可以用session来表示一个会话列表，每个user都有多个session，session包含sender和target，
- 我现在知道的只有一个hub，hub包含了很多clients，新增一个连接后，就把它加入clients，先解决点对点的问题再解决历史记录问题把。
- 首先定义几个标准消息，最重要的是pub和sub, 首先后台要区分一下是谁发过来的消息，然后根据你的target发送给指定的人