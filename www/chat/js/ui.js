"use strict";
//
import { el_1, create_element, dcTN, messageCache, output, getUserFromUUID, el_4 } from "./util.js";
import { Channel } from "./ui.channel.js";

//

export const volatile = {
    activeChannel: null,
};

//

export async function addChannel(ch) {
    el_1.firstElementChild.appendChild(create_element("li", [["data-uuid",ch.uuid],["data-unread","0"]], [
        create_element("div", [], [dcTN(ch.name)]),
        create_element("div", [["class","unred"]], [dcTN("0")]),
    ]));
    messageCache.set(ch.uuid, []);

    await fetch(`./../api/channels/${ch.uuid}/messages`).then((x) => x.json()).then((x) => {
        for (const item of x.message) {
            messageCache.get(ch.uuid).unshift(item);
        }
    });
}

export function createMessage(user, msg) {
    const attrs = [
        ["class","msg"],
        ["data-msg-uid",msg.uuid],
        ["data-user-uid",user.uuid],
    ];
    return create_element("div", attrs, [
        create_element("div", [["class","ts"],["title",msg.time]], [dcTN(msg.time.substring(msg.time.indexOf(" ")))]),
        create_element("div", [["class","usr"]], [dcTN(user.name + ": ")]),
        create_element("div", [["class","dat"]], [dcTN(msg.body)]),
    ]);
}

export function addMessage(channel, from, message, save=true) {
    channel = channel ? channel : volatile.activeChannel.dataset.uuid;
    from.uuid = from.uuid ? from.uuid : "";
    const at_bottom = output.scrollTop === output.scrollTopMax;
    if (channel === null || output.dataset.active === channel) {
        const time = new Date(message.time ? message.time : Date.now()).toLocaleString();
        output.appendChild(createMessage(from, {...message, time}));
    }
    if (output.dataset.active !== channel) {
        const c = new Channel(channel);
        c.unread += 1;
    }
    if (at_bottom) output.scrollTop = output.scrollHeight;
    if (save === true) messageCache.get(channel).push(message);
}

export async function setActiveChannel(uid) {
    console.debug("channel-switch:", uid);
    const ac = el_1.querySelector(".active");
    if (ac !== null) ac.classList.remove("active");
    const c = new Channel(uid);
    if (c.el === null) return;
    volatile.activeChannel = c.el;
    volatile.activeChannel.classList.add("active");
    //
    output.dataset.active = uid;
    output.removeAllChildren();
    const new_message_history = messageCache.get(uid);
    for (const item of new_message_history) {
        addMessage(null, await getUserFromUUID(item.author), item, false, false);
    }
    //
    c.unread = 0;
    output.classList.remove("loading-done");
}

export async function setMemberOnline(uid) {
    console.debug("user-ws-connect", uid);
    const ue = el_4.querySelector(`[data-user="${uid}"]`);
    if (ue === null) {
        const u = await getUserFromUUID(uid);
        el_4.appendChild(create_element("li", [["data-user",uid]], [
            dcTN(u.name),
        ]));
    }
}

export function setMemberOffline(uid) {
    console.debug("user-ws-disconnect", uid);
    const ue = el_4.querySelector(`[data-user="${uid}"]`);
    if (ue !== null) {
        ue.remove();
    }
}
