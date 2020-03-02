"use strict";
//
import { el_1, el_2, el_3, getUserFromUUID, output, messageCache, el_4, create_element, dcTN } from "./util.js";
import * as ui from "./ui.js";
import * as client from "./client.js";

//
let me = null;

//
(async function() {
    //
    await fetch("./../api/about").then((x) => x.json()).then((x) => {
        console.info(x);
        el_2.innerText = x.message.name;
    });

    //
    await fetch("./../api/users/@me").then((x) => x.json()).then((x) => {
        console.info(x);
        if (x.success === false) {
            location.assign("../");
            return;
        }
        console.info(x.message);
        me = x.message.me;
        const n = me.nickname || me.name;
        el_3.children[0].textContent = `@${n}`;
        const p = x.message.perms;
        for (const key in p) {
            if (!p[key]) {
                document.querySelectorAll(`[data-requires^="${key}"]`).forEach((el) => {
                    el.setAttribute("hidden", "");
                });
            }
        }
    });

    //
    await fetch("./../api/channels/@me").then((x) => x.json()).then(async (x) => {
        console.info(x);
        for (const item of x.message) {
            console.info(item);
            await ui.addChannel(item);
        }
        await ui.setActiveChannel(x.message[0].uuid);

        const el2 = document.getElementById("channel-name");
        el2.children[0].textContent = x.message[0].name;
        el2.children[1].textContent = x.message[0].description;

        el_1.querySelector("button").addEventListener("click", async () => {
            const {value: name} = await Swal({
                title: "Enter the new channel's name",
                input: "text",
                showCancelButton: true,
                inputValidator: (value) => !value && "You need to write something!",
            });
            if (name !== undefined) {
                const fd = new URLSearchParams();
                fd.append("name", name);
                return fetch("./../api/channels/create", {
                    method: "post",
                    credentials: "include",
                    headers: {
                        "content-type": "application/x-www-form-urlencoded",
                    },
                    body: fd,
                });
            }
        });
        el_1.querySelector("ol").addEventListener("click", (ev) => {
            const fl = ev.composedPath().filter((v) => v instanceof Element && v.matches("[data-uuid]"));
            if (fl.length === 0) return;
            ui.setActiveChannel(fl[0].dataset.uuid);
        });
        output.addEventListener("scroll", async (e) => {
            if (output.children.length === 0) return;
            if (e.target.scrollTop !== 0) return;
            if (output.classList.contains("loading")) return;
            if (output.classList.contains("loading-done")) return;
            //
            output.classList.add("loading");
            const fc = output.children[0];
            const lstm = output.children[0].dataset.msgUid;
            const chuid = ui.volatile.activeChannel.dataset.uuid;
            await fetch(`./../api/channels/${chuid}/messages?after=${lstm}`).then((y) => y.json()).then(async (y) => {
                if (y.message.length <= 1) {
                    output.classList.add("loading-done");
                    return;
                }
                for (let i = 1; i < y.message.length; i++) {
                    const item = y.message[i];
                    output.prepend(ui.createMessage(await getUserFromUUID(item.author), item));
                    messageCache.get(chuid).unshift(item);
                }
                output.scrollTop = fc.offsetTop-60;
            });
            output.classList.remove("loading");
        });
        document.addEventListener("keydown", async (e) => {
            if (e.key !== "Delete") return;
            if (document.activeElement !== document.body) return;
            if (!output.isInViewport()) return;
            if (ui.volatile.selectedMsgs.length === 0) return;
            await Swal.fire({
                title: "Are you sure you want to delete?",
                text: "You won't be able to revert this!",
                type: "warning",
                showCancelButton: true,
            }).then(async (r) => {
                if (!r.value) return;
                const m2d = ui.volatile.selectedMsgs.filter((v) => v.dataset.userUid === me.uuid).map((v) => v.dataset.msgUid);
                const fd = new FormData();
                m2d.forEach((v) => fd.append("ids", v));
                await fetch(`./../api/channels/${ui.volatile.activeChannel.dataset.uuid}/messages`, {
                    method: "DELETE",
                    body: fd,
                }).then((y) => y.json()).then(() => {
                    //
                });
            });
        });
    });

    await fetch("./../api/roles").then((x) => x.json()).then((x) => {
        const rls = Object.values(x.message);
        for (const item of rls) {
            if (!item.distinguish) continue;
            el_4.appendChild(create_element("div", [["data-count","0"]], [dcTN(item.name)]));
            el_4.appendChild(create_element("ul", [["data-uid",item.uuid]], []));
        }
        el_4.appendChild(create_element("div", [["data-count","0"]], [dcTN("Online")]));
        el_4.appendChild(create_element("ul", [["data-uid",""]], []));
    });

    await fetch("./../api/users/online").then((x) => x.json()).then((x) => {
        console.info(x);
        for (const item of x.message) {
            ui.setMemberOnline(item);
        }
    });

    //
    const input = document.getElementById("input").children[0];
    const socket = new WebSocket(`ws${location.protocol.substring(4)}//${location.host}/ws`);

    socket.addEventListener("open", () => {
        el_2.classList.remove("loading");
        el_2.classList.add("online");
    });
    socket.addEventListener("close", () => {
        el_2.classList.remove("online");
    });
    socket.addEventListener("message", async (e) => {
        const d = JSON.parse(e.data);
        switch (d.type) {
            case "pong": {
                // do nothing, keep connection alive
                console.debug("pong");
                break;
            }
            case "message": {
                const u = await getUserFromUUID(d.message.author);
                ui.addMessage(d.in, u, d.message, true);
                break;
            }
            case "new-channel": {
                ui.addChannel(d);
                break;
            }
            case "user-connect": {
                ui.setMemberOnline(d.user);
                break;
            }
            case "user-disconnect": {
                ui.setMemberOffline(d.user);
                if (d.user === me.uuid) socket.close();
                break;
            }
            case "message-delete": {
                const a1 = messageCache.get(d.channel);
                const a2 = a1.filter((v) => d.affected.includes(v.uuid));
                a2.forEach((v) => a1.splice(a1.indexOf(v), 1));
                if (output.dataset.active !== d.channel) break;
                for (const item of d.affected) {
                    const de = output.querySelector(`.msg[data-msg-uid="${item}"]`);
                    if (de !== null) de.remove();
                }
                break;
            }
            default: {
                console.log(d);
            }
        }
    });
    setInterval(() => {
        if (el_2.classList.contains("online")) {
            socket.send(JSON.stringify({
                type: "ping",
            }));
        }
    }, 30*1000);

    //
    input.addEventListener("keydown", (e) => {
        if (e.key === "Enter") {
            let msg_con = e.target.value;
            for (const item of client.commands) {
                if (msg_con.startsWith("/"+item[0])) {
                    msg_con = item[1](msg_con.replace("/"+item[0],"")).trim();
                }
            }
            e.target.value = "";
            if (msg_con===null) return;
            if (msg_con===undefined) return;
            if (msg_con.length === 0) return;
            socket.send(JSON.stringify({
                type: "message",
                in: ui.volatile.activeChannel.dataset.uuid,
                message: msg_con,
            }));
        }
    });
})();
