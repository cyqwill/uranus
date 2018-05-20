# Uranus

![](https://i.loli.net/2018/05/14/5af93f6aabc6a.jpeg)
**Uranus**: means *Sky* and *Hope*. The essential existence of my super universe. This project will serves all my client messages.
All clients exchanges messages and information through **Uranus**. And I will receive the final result of all devices comes in msg.

**Peer 2 peer chat now working very well!**

![](https://i.loli.net/2018/05/20/5b0196f93d267.png)

## Structure
Simply, **Uranus** contains a websocket for real-time communication and a long-pull progress for serving in-coming messages. So that every robot I built, every devices I make it connected into **Uranus**, I can access them. The main part of **Uranus** are:
- Websocket for normal devices such as mobile phones, websites etc.;
- MQTT for Iot devices, such as little spare sensors, embed in devices like Arduino and RaspberryPI.

With all those powered clients, I can receive all messages on **Uranus Central Hub**, which gives me full information things I want to know.


## Install
Pls make sure postgres were installed, and there is an database named **uranus**. Then you should change the database username and password in `config.toml` file.

There are quite some pre-requirements you should follow. Install all third-party packages inside *uranus*. and you should prepare your database for our *unranus* work properly.




## API Glance

Here is the api, you can also see `API.md` for detail.

#### 1. Login

the first step, you should login to uranus, this is because we have to got the token in our local, if we don't have token, you can not get access to uranus server.
For login, just call this api:
```
/api/v1/user_login
```
and send your `user_name` and `user_password`:
```json
{
    "user_name": "lucasjin",
    "user_password": "123123",
}
```
And you will got a msg like welcome, you are login success now (if you wanna new an account, just register.)

#### 2. Send Msg

you want develop a new client, then you should probably want send msg to other clients, there are some `MsgType` you **must know**.

- `hi`: Hi msg for let server knows your client type and your user address, you will get your token and address when login;
- `send`: If you start to send msg, then using send msg type;
- `add`: msg for add into a group;
- `del`: delete a msg or you want leave a group or delete a person;

OK, that's all, those are all the msg type in **uranus** now.

Now, this is the all msg structure:

**{hi}**

```json
{
    "token": "ehurhuoetihry.8u08954hggh.pguwtyh",
    "user_addr": "usry7bvdgeug",
    "ua": "uranus-0.1/PiLinux",
    "location": "湖南长沙",
    "device": "RaspberryPi",
}
```

**{send}**

```json
{
    "target_addr": "usry7u89ghutehu",
    "send_addr": "usry7bvdgeug",
    "content": "你他妈的还不回来吃午饭啊",
    "msg_type": 0,
}
```

## Current Functionality
- websocket 1 to 1 chat;
- websocket many to many chat;
- websocket 1 to many chat;
- websocket many to one chat;
- User register and identification;
- Needed user information;




## CopyRight
**Uranus** is a trademark, original implemented by *Lucas Jin*, you should using it under MIT license.