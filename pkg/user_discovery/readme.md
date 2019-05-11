# read me plz

## user connects

    Authenticate
    create session
        add session to db
        subscribe to discovery channel
        add session to redis

## publish new message

    if field(channel) => goto 3
    if field(to):
        + query for existing channel
        + if not exist:
            + create channel
        pass channel_id goto 3
    3: ensure_channel(channel_id)
    publish_message_to(channel_id)

## ensure channel

    + get current session
    + if !session.is_subscribed_to(channel_id): //redis
        // subscribe_to_channel(channel_id)
        + user_discovery.PublishEvent(channel_id)

## publish message to

    go add_to_db
    go publish_to_nats

--------------------------------------

# subscribe
## (when user connects to server):

   amir+ subscribe to channel(user_discovery)


## (when recieve a message from any channel)
    if (channel_id == U_D ):
        recieve message from discovery() 
    else:
        recieve message handler()


## recieve message from discovery

   if ( user.id = discovery.userid ):
       subscribe to channel(discovery.channel_id)


## subscribe to channel

    subscribe_to_channel(channel_id)
    add_channel_to_user_session(channel_id)

## recieve message handler

    deliver to user with grpc


//user discovery event handler check if not already subscribed to channel subscribe
//new message event handler push kon be user ama che goone ??????