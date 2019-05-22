# Hermes
### What is Hermes
Hermes is a scalable, Database agonstic, GRPC based messaging service.
### Technology Stack
Nats-Streaming <br>
MongoDB <br>
Redis <br>
### How Hermes Works
Hermes uses Nats-Streaming as pub/sub service and uses nats concept of channels as 
a way of communication between two users. Each new event 
will be published to a channel and the eventHandler specific to the event (there are only two event handlers)
will deliver the event to the users using GRPC streaming.
#### Flow of new message
1.User creates a new session as a kind of registration of his/her device.<br>
2.User sends a join event<br>
3.User is now subscribed to a channel called user-discovery and will receive channelIds that he/she needs to subscribe to.<br>
4.User now subscribes to one of channels which has a new message for him/her self.<br>
5.NewMessage event handler recieves new event and push it to user.<br>
6.Done 
### Installation
```
   git clone github.com/hermes/hermes
   cd hermes
   docker-compose up .
```
### License 
Hermes uses MIT licenese for more information read licence file.