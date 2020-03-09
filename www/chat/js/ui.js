"use strict";
//
// jshint -W003
import { create_element, dcTN, numsBetween, ele_atBottom, deActivateChild, setDataBinding } from "./util.js";
import { Channel } from "./ui.channel.js";
import { SidebarRole } from "./ui.sidebar_role.js";
import { el_1, messageCache, output, getUserFromUUID, el_4, roleCache, userCache } from "./ui.util.js";

//

export const volatile = {
    activeChannel: null,
    /** @type {Element[]} */
    selectedMsgs: [],
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
    const attrs = [["class","msg"]];
    const attrsU = [["class","usr"]];
    if (msg.uuid) attrs.push(["data-msg-uid",msg.uuid]);
    if (user.uuid) attrs.push(["data-user-uid",user.uuid]);
    if (user.roles) {
        const a = user.roles.split(",").
            map((v) => roleCache.get(v)).
            sort((b,c) => b.position > c.position).
            filter((v) => v.color.length > 0);
        if (a.length > 0) {
            attrsU.push(["data-role",a[0].uuid]);
        }
    }
    //
    const time = new Date(msg.time||Date.now()).toLocaleString("en-GB");
    const el = create_element("div", attrs, [
        create_element("div", [["class","ts"],["title",time]], [dcTN(time.substring(time.indexOf(" ")))]),
        create_element("div", attrsU, [dcTN(user.name)]),
        create_element("div", [["class","dat"]], [dcTN(msg.body)]),
    ]);
    el.children[2].innerHTML = el.children[2].textContent.replace(/(https?:\/\/[^\s]+)/gu, (match) => `<a target="_blank" href="${match}">${match}</a>`);
    twemoji.parse(el.children[2]);
    //
    if (msg.uuid) {
        el.addEventListener("click", (e) => {
            /** @type {Element[]} */
            const fl = e.composedPath().filter((v) => v instanceof Element && v.matches("[data-msg-uid]"));
            if (fl.length === 0) return;
            const et = fl[0];
            if (e.ctrlKey) {
                if (et.classList.contains("selected")) {
                    volatile.selectedMsgs.splice(volatile.selectedMsgs.indexOf(et), 1);
                }
                else {
                    volatile.selectedMsgs.unshift(et);
                }
                et.classList.toggle("selected");
            }
            if (e.shiftKey) {
                if (volatile.selectedMsgs.length === 0) return;
                const p1 = Array.from(et.parentElement.children).indexOf(et);
                const p2 = Array.from(et.parentElement.children).indexOf(volatile.selectedMsgs[0]);
                const nb = numsBetween(p1, p2);
                for (let i = 0; i < nb.length; i++) {
                    const mc = et.parentElement.children[nb[i]];
                    if (mc === volatile.selectedMsgs[0]) continue;
                    mc.classList.add("selected");
                    volatile.selectedMsgs.unshift(mc);
                }
            }
        });
    }
    if (user.uuid) {
        el.querySelector(".usr").addEventListener("click", (e) => {
            const userN = userCache.get(user.uuid);
            setDataBinding("pp_user_name", userN.name);
            setDataBinding("pp_user_id", userN.id);
            setDataBinding("pp_user_uuid", userN.uuid);
            setDataBinding("pp_user_provider", userN.provider);
            setDataBinding("pp_user_snowflake", userN.snowflake);
            const pp = document.querySelector("dialog.popup.user");
            const ppr = pp.querySelector("ol");
            const rls = userN.roles.split(",").
                filter((v) => v.length > 0).
                map((v) => roleCache.get(v)).
                sort((a,b) => a.position > b.position);
            deActivateChild(ppr);
            const pps = pp.querySelector("div ol");
            deActivateChild(pps);
            for (const item of rls) {
                const ppra = ppr.querySelector(`[data-uid="${item.uuid}"]`);
                if (ppra === null) continue;
                ppra.classList.add("active");
                const ppsa = pps.querySelector(`[data-uid="${item.uuid}"]`);
                if (ppsa === null) continue;
                ppsa.classList.add("active");
            }
            pp.setAttribute("open","");
            pp.style.top = e.y+"px";
            pp.style.left = e.x+"px";
        });
    }
    //
    return el;
}

export function addMessage(channel, from, message, save=true) {
    channel = channel ? channel : volatile.activeChannel.dataset.uuid;
    from.uuid = from.uuid ? from.uuid : "";
    const at_bottom = ele_atBottom(output);
    if (channel === null || output.dataset.active === channel) {
        output.appendChild(createMessage(from, message));
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
    volatile.selectedMsgs.splice(0, volatile.selectedMsgs.length);
}

export async function setMemberOnline(uid) {
    console.debug("user-ws-connect", uid);
    const ue = el_4.querySelector(`li[data-user="${uid}"]`);
    if (ue === null) {
        const u = await getUserFromUUID(uid);
        for (const item of el_4.querySelectorAll("ul")) {
            if (!u.roles.includes(item.dataset.uid)) continue;
            item.appendChild(create_element("li", [["data-user",uid]], [
                create_element("span", null, [dcTN(u.name)]),
                create_element("span", null, [dcTN("#"+u.id)]),
            ]));
            new SidebarRole(item).count += 1;
            break;
        }
    }
}

export function setMemberOffline(uid) {
    console.debug("user-ws-disconnect", uid);
    const ue = el_4.querySelector(`li[data-user="${uid}"]`);
    if (ue === null) return;
    new SidebarRole(ue.parentElement).count -= 1;
    ue.remove();
}

export function addRole(role) {
    roleCache.set(role.uuid, role);
    //
    const rlist = document.querySelector("x-settings[data-s-for=server] [data-s-section=roles] .selection nav");
    const oLen = rlist.children.length;
    const nEl = create_element("a", [["data-uid",role.uuid]], [dcTN(role.name)]);
    nEl.addEventListener("click", (e) => {
        const et = e.target;
        settingsRolesSetActive(Array.from(et.parentElement.children).indexOf(et));
    });
    rlist.insertBefore(
        nEl,
        rlist.querySelector(".div"),
    );
    if (oLen === 2) {
        rlist.parentElement.classList.add("active");
        settingsRolesSetActive(0);
    }
    //
    //
    const nEl2 = create_element("li", [["data-uid",role.uuid]], [dcTN(role.name)]);
    document.querySelector("dialog.popup.user ol").appendChild(nEl2.cloneNode(true));
    nEl2.addEventListener("click", (e) => {
        const et = e.target;
        const rid = et.dataset.uid;
        const uid = document.querySelector("[data-bind=pp_user_uuid]").textContent;
        const fd = new FormData();
        fd.append("p_name", et.classList.contains("active") ? "remove_role" : "add_role");
        fd.append("p_value", rid);
        fetch(`./../api/users/${uid}/update`, { method: "put", body: fd, });
        const ett = et.parentElement.parentElement.previousElementSibling.querySelector(`[data-uid="${rid}"]`);
        if (et.classList.contains("active")) ett.classList.remove("active"); else ett.classList.add("active");
        if (et.classList.contains("active")) et.classList.remove("active"); else et.classList.add("active");
    });
    document.querySelector("dialog.popup.user div ol").appendChild(nEl2);
}

export function settingsRolesSetActive(i) {
    const rlist = document.querySelector("x-settings[data-s-for=server] [data-s-section=roles] .selection nav");
    deActivateChild(rlist);
    rlist.children[i].classList.add("active");
    const r = roleCache.get(rlist.children[i].dataset.uid);
    const tin = rlist.parentElement.querySelectorAll("[fill]");
    for (const item of tin) {
        item.setAttribute("fill", r.uuid);
    }
    for (const item of ["name","color"]) {
        rlist.parentElement.querySelector(`[name="${item}"]`).setAttribute("value", r[item]);
    }
}
