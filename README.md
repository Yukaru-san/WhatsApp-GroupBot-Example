# WhatsApp-GroupBot-Example [Outdated]

<b>With the release of WhatsApp's multi-device support this program stopped working. I may try to update it somewhere in the future!</b>

This project is an example on how to use my [WhatsApp-GroupBot](https://github.com/Yukaru-san/WhatsApp-GroupBot).<br><br>
I used this bot for myself a while ago and decided to make it public so people can get some inspiration on how to get started on programming their own.<br><br>
That being said I removed some of its features to protect the privacy of my contacts and users.

# Usage
Clone the program or download it as zip, adjust the parts of the source code that require a name of some group or contact and compile it using<br> 
```
go get
go build
```

# Lines to edit for yourself
The data that needs to be changed is always enclosed in { }.
<br> If one of your groups is named "Party Group" then replace "{Some Group Name}" with "Party Group".<br>
Here are the lines where you can find the names you need to alter:

```
main.go
        : 26  -> Fill the array with actual group names
        : 97  -> Replace {Some Group Name} with a group's name
        :102  -> Replace {Some Group Name} and {Another Group Name} with a group's name

stickerMethods.go 
        : 15  ->  Replace {somesticker} with the name of an actual sticker. Comment out if you don't know how to create one.
        : 23  ->  Change to the name of a contact saved on your phone
```
