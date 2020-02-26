/**
 */
//
import { el_1, el_2, el_3, create_element, dcTN, messageCache, output, getUserFromUUID, el_4 } from "./util.js";
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
    ]))
    messageCache.set(ch.uuid, []);

    await fetch(`./../api/channels/${ch.uuid}/messages`).then(x=>x.json()).then(x => {
        for (const item of x.message) {
            messageCache.get(ch.uuid).unshift([item.author, item.body]);
        }
    })
}

export async function addMessage(channel=volatile.activeChannel.dataset.uuid, from, message, raw_from=false, save=true, at=Date.now()) {
    const at_bottom = output.scrollTop === output.scrollTopMax;
    from.uuid = raw_from ? "" : from.uuid;
    const name = from.nickname || from.name;
    const time = new Date(at).toLocaleString();
    if (raw_from || output.dataset.active === channel) {
        output.appendChild(createMessage(from, {body:message,time:time}));
    }
    if (output.dataset.active !== channel) {
        const c = new Channel(channel);
        c.unread += 1;
    }
    if (at_bottom) output.scrollTop = output.scrollHeight;
    if (save===true) messageCache.get(channel).push([from.uuid, message]);
}

function createMessage(user, msg) {
    return create_element("div", [["class","msg"],["data-msg-uid",msg.uuid],["data-user-uid",user.uuid]], [
        create_element("div", [["class","ts"],["title",msg.time]], [dcTN(msg.time.substring(msg.time.indexOf(" ")))]),
        create_element("div", [["class","usr"]], [dcTN(user.name + ": ")]),
        create_element("div", [["class","dat"]], [dcTN(msg.body)])
    ])
}

export async function setActiveChannel(uid) {
    console.debug("channel-switch:", uid);
    let ac = el_1.querySelector(".active");
    if (ac !== null) ac.classList.remove("active");
    const c = new Channel(uid);
    volatile.activeChannel = c.el;
    volatile.activeChannel.classList.add("active");
    //
    output.dataset.active = uid;
    output.removeAllChildren();
    const new_message_history = messageCache.get(uid);
    for (const item of new_message_history) {
        addMessage(uid, await getUserFromUUID(item[0]), item[1], false, false);
    }
    //
    c.unread = 0;
}

export async function setMemberOnline(uid) {
    console.debug("user-ws-connect", uid);
    const ue = el_4.querySelector(`[data-user="${uid}"]`);
    if (ue === null) {
        const u = await getUserFromUUID(uid);
        el_4.appendChild(create_element("li", [["data-user",uid]], [
            dcTN(u.name)
        ]));
    }
}

export async function setMemberOffline(uid) {
    console.debug("user-ws-disconnect", uid);
    const ue = el_4.querySelector(`[data-user="${uid}"]`);
    if (ue !== null) {
        ue.remove();
    }
}
