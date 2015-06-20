

This provides the channels used by ms as well as the message formats.


TODO - make idea of key scopes with different auth rules for some

# Channels

Things have both names and addresses. The names are things like @Cullen_Jennings or @fluffy or #coffee. Names are not globally unique but are evaluated ina given context. Names map to address which are statistically unique identifiers such as ChannelIDs or UserID.

### User Names

These start with an @

### Channel Names

Theses start with an #

### Special Channel Names

These start with and _

_null is never used and has channelID 0.

_bad is never used and has channelID 1. 

_userList is a single global channel that anyone can write their PublicUserInfo to. channelID is 2. 

_channelList is  single channel that anyone can write a PublicChannelInfo to. channelID is 3. 

_me is channel where the chID is the userID of the owner. It is used to save configuration information for that user. Only this user and can read and write to it.

_request:<userID> is a incoming request queue for request from others to this user. The channelID is found by looking at the uses PublicUserInfo. 


# Messages

## PublicUserInfo

This is a subset of PrivateUserInfo. Saved to _directory. Has channelID for this users  _request:<userID> channel. 

## PrivateUserInfo

Save to _me 


## PublicChannelInfo

Has name of channel and keyID. Subset of PrivateChannelInfo

## PrivateChannelInfo

saved to _me channel of each user of this channel

## reqestInvite

Written to _request:<userID>

## requestJoin


## TextMessage

Contains message in a channel. Formatted as ???. Unicode ???. Saved in normal #channel 


## bookmark

Need better name. Contains a channelID and sequence number indicating current scroll location in that channel. Saved in _me


