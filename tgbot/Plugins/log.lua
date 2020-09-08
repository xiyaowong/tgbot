local log = require("log")
local json = require("json")
local http = require("http")
local mysql = require("mysql")

function ReceiveTGMsg(data)
    log.info("%s", "\nReceiveTGMsg")
    if data.MsgType == 'MessageText' then
        str = string.format(
            "ChatID %d\nChatTitle %s\nMessageID %d\nSenderUserID %d\nSenderUserFirstName %s\nSenderUserLastName %s\nSenderUserName %s\nSenderUserPhoneNumber %s\nMsgType %s\nContent %s",
            data.ChatID,
            data.ChatTitle,
            data.MessageID,
            data.SenderUserID,
            data.SenderUserFirstName,
            data.SenderUserLastName,
            data.SenderUserName,
            data.SendPhoneNumber,
            data.MsgType,
            data.Content
        )
    elseif data.MsgType == 'MessagePhoto' then
        str = string.format(
            "ChatID %d\nChatTitle %s\nMessageID %d\nSenderUserID %d\nSenderUserFirstName %s\nSenderUserLastName %s\nSenderUserName %s\nSenderUserPhoneNumber %s\nMsgType %s\nContent %s\nPhotoBase64 %s",
            data.ChatID,
            data.ChatTitle,
            data.MessageID,
            data.SenderUserID,
            data.SenderUserFirstName,
            data.SenderUserLastName,
            data.SenderUserName,
            data.SendPhoneNumber,
            data.MsgType,
            data.Content,
            '略'
    )
    end
    
    log.notice("From log.lua Log\n%s", str)
    return 1
end

function ReceiveEvents(data)
    return 1
end
