# Uranus API Doc

For the convenience of develop something beyond **uranus**, we must be very familiar with the apis in uranus. We have some very simple steps to follow:

## Login

the first step, you should login to uranus, this is because we have to got the token in our local, if we don't have token, you can not get access to uranus server.
For login, just call this api:
```
/api/v1/user_login
```
and send your `user_name` and `user_password`:
```
{
    "user_name": "lucasjin",
    "user_password": "123123",
}
```
And you will got a msg like welcome, you are login success now (if you wanna new an account, just register.)

## Send Msg

you want develop a new client, then you should probably want send msg to other clients, there are some `MsgType` you **must know**.

- `hi`: Hi msg for let server knows your client type and your user address, you will get your token and address when login;
- `send`: If you start to send msg, then using send msg type;
- `add`: msg for add into a group;
- `del`: delete a msg or you want leave a group or delete a person;

OK, that's all, those are all the msg type in **uranus** now.

Now, this is the all msg structure:

**{hi}**

```
{
    "token": "ehurhuoetihry.8u08954hggh.pguwtyh",
    "user_addr": "usry7bvdgeug",
    "ua": "uranus-0.1/PiLinux",
    "location": "湖南长沙",
    "device": "RaspberryPi",
}
```

**{send}**

```
{
    "target_addr": "usry7u89ghutehu",
    "send_addr": "usry7bvdgeug",
    "content": "你他妈的还不回来吃午饭啊",
    "msg_type": 0,
}
```